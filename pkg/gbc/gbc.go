package gbc

import (
	"fmt"
	"math"
	"os"
	"runtime"

	"gbc/pkg/emulator/config"
	"gbc/pkg/gbc/apu"
	"gbc/pkg/gbc/cart"
	"gbc/pkg/gbc/joypad"
	"gbc/pkg/gbc/rtc"
	"gbc/pkg/gbc/scheduler"
	"gbc/pkg/gbc/video"
	"gbc/pkg/util"
)

var irqVec = [5]uint16{0x0040, 0x0048, 0x0050, 0x0058, 0x0060}

// ROM
//
// 0x0000-0x3fff: bank0
//
// 0x4000-0x7fff: bank1-256
type ROM struct {
	bank   byte
	buffer [256][0x4000]byte
}

// RAM - 0xa000-0xbfff
type RAM struct {
	bank   byte
	Buffer [16][0x2000]byte // num of banks changes depending on ROM
}

// WRAM
//
// 0xc000-0xcfff: bank0
//
// 0xd000-0xdfff: bank1-7
type WRAM struct {
	// fixed at 1 on DMG, changes from 1 to 7 on CGB
	bank   byte
	buffer [8][0x1000]byte
}

type Dma struct {
	src, dest uint16
	remaining int
}

type Hdma struct {
	enable    bool
	src, dest uint16
	remaining int
}

const (
	NoIRQ = iota
	VBlankIRQ
	LCDCIRQ
	TimerIRQ
	SerialIRQ
	JoypadIRQ
)

// GBC core structure
type GBC struct {
	Reg         Register
	IO          [0x100]byte // 0xff00-0xffff
	Cartridge   *cart.Cartridge
	joypad      *joypad.Joypad
	halt        bool
	Config      *config.Config
	timer       *Timer
	ROM         ROM
	RAM         RAM
	WRAM        WRAM
	bankMode    uint
	sound       *apu.APU
	Video       *video.Video
	RTC         rtc.RTC
	DoubleSpeed bool
	model       util.GBModel
	irqPending  int
	scheduler   *scheduler.Scheduler
	dma         Dma
	hdma        Hdma
	cpuBlocked  bool
}

// TransferROM Transfer ROM from cartridge to Memory
func (g *GBC) TransferROM(rom []byte) {
	switch g.Cartridge.Type {
	case 0x00:
		g.Cartridge.MBC = cart.ROM
		g.transferROM(2, rom)
	case 0x01: // Type : 1 => MBC1
		g.Cartridge.MBC = cart.MBC1
		switch r := int(g.Cartridge.ROMSize); r {
		case 0, 1, 2, 3, 4, 5, 6:
			g.transferROM(int(math.Pow(2, float64(r+1))), rom)
		default:
			errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x02, 0x03: // Type : 2, 3 => MBC1+RAM
		g.Cartridge.MBC = cart.MBC1
		switch g.Cartridge.RAMSize {
		case 0, 1, 2:
			switch r := int(g.Cartridge.ROMSize); r {
			case 0, 1, 2, 3, 4, 5, 6:
				g.transferROM(int(math.Pow(2, float64(r+1))), rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
				panic(errorMsg)
			}
		case 3:
			g.bankMode = 1
			switch r := int(g.Cartridge.ROMSize); r {
			case 0:
			case 1, 2, 3, 4:
				g.transferROM(int(math.Pow(2, float64(r+1))), rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
				panic(errorMsg)
			}
		default:
			errorMsg := fmt.Sprintf("RAMSize is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x05, 0x06: // Type : 5, 6 => MBC2
		g.Cartridge.MBC = cart.MBC2
		switch g.Cartridge.RAMSize {
		case 0, 1, 2:
			switch r := int(g.Cartridge.ROMSize); r {
			case 0, 1, 2, 3:
				g.transferROM(int(math.Pow(2, float64(r+1))), rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
				panic(errorMsg)
			}
		case 3:
			g.bankMode = 1
			switch r := int(g.Cartridge.ROMSize); r {
			case 0:
			case 1, 2, 3:
				g.transferROM(int(math.Pow(2, float64(r+1))), rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
				panic(errorMsg)
			}
		default:
			errorMsg := fmt.Sprintf("RAMSize is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x0f, 0x10, 0x11, 0x12, 0x13: // Type : 0x0f, 0x10, 0x11, 0x12, 0x13 => MBC3
		g.Cartridge.MBC, g.RTC.Enable = cart.MBC3, true
		switch r := int(g.Cartridge.ROMSize); r {
		case 0, 1, 2, 3, 4, 5, 6:
			g.transferROM(int(math.Pow(2, float64(r+1))), rom)
		default:
			errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x19, 0x1a, 0x1b: // Type : 0x19, 0x1a, 0x1b => MBC5
		g.Cartridge.MBC = cart.MBC5
		switch r := int(g.Cartridge.ROMSize); r {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			g.transferROM(int(math.Pow(2, float64(r+1))), rom)
		default:
			errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
			panic(errorMsg)
		}
	default:
		errorMsg := fmt.Sprintf("Type is invalid => type:%x rom:%x ram:%x\n", g.Cartridge.Type, g.Cartridge.ROMSize, g.Cartridge.RAMSize)
		panic(errorMsg)
	}
}

func (g *GBC) transferROM(bankNum int, rom []byte) {
	for bank := 0; bank < bankNum; bank++ {
		for i := 0x0000; i <= 0x3fff; i++ {
			g.ROM.buffer[bank][i] = rom[bank*0x4000+i]
		}
	}
}

func (g *GBC) resetRegister() {
	g.Reg.setAF(0x11b0) // A=01 => GB, A=11 => CGB
	g.Reg.setBC(0x0013)
	g.Reg.setDE(0x00d8)
	g.Reg.setHL(0x014d)
	g.Reg.PC, g.Reg.SP = 0x0100, 0xfffe
}

func New(romData []byte, j [8](func() bool)) *GBC {
	g := &GBC{
		Cartridge: cart.New(romData),
		scheduler: scheduler.New(),
		joypad:    joypad.New(j),
	}

	g.Video = video.New(&g.IO, g.updateIRQs, g.hdmaMode3, g.scheduler.ScheduleEvent, g.scheduler.DescheduleEvent)
	if g.Cartridge.IsCGB {
		g.setModel(util.GB_MODEL_CGB)
	}

	g.timer = NewTimer(g)
	g.scheduler.ScheduleEvent(scheduler.TimerUpdate, g.timer.update, 0)

	g.resetRegister()
	g.resetIO()

	g.ROM.bank, g.WRAM.bank = 1, 1

	g.Config = config.New()

	// Init APU
	g.sound = apu.New(true)

	// Init RTC
	go g.RTC.Init()

	g.scheduler.ScheduleEvent(scheduler.EndMode2, g.Video.EndMode2, video.MODE_2_LENGTH)

	g.TransferROM(romData)
	return g
}

// Exec 1cycle
func (g *GBC) step() {
	PC := g.Reg.PC
	bytecode := g.Load8(PC)
	opcode := opcodes[bytecode]
	instruction, operand1, operand2, cycle, handler := opcode.Ins, opcode.Operand1, opcode.Operand2, opcode.Cycle1, opcode.Handler

	if !g.halt {
		if g.irqPending > 0 {
			oldIrqPending := g.irqPending
			g.irqPending = 0
			g.Reg.IME = false
			g.updateIRQs()
			g.triggerIRQ(int(oldIrqPending - 1))
			return
		} else if handler != nil {
			handler(g, operand1, operand2)
		} else {
			switch instruction {
			case INS_SBC:
				g.SBC(operand1, operand2)
			default:
				errMsg := fmt.Sprintf("eip: 0x%04x opcode: 0x%02x", g.Reg.PC, bytecode)
				panic(errMsg)
			}
		}
		cycle *= (4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	} else {
		cycle = int(g.scheduler.Next() - g.scheduler.Cycle())
	}

	g.timer.tick(uint32(cycle))
}

// 1 frame
func (g *GBC) Update() error {
	frame := g.Frame()
	if frame%3 == 0 {
		g.handleJoypad()
	}

	for frame == g.Video.FrameCounter {
		g.step()
	}
	return nil
}

func (g *GBC) PanicHandler(place string, stack bool) {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "%s emulation error: %s in 0x%04x\n", place, err, g.Reg.PC)
		for depth := 0; ; depth++ {
			_, file, line, ok := runtime.Caller(depth)
			if !ok {
				break
			}
			fmt.Printf("======> %d: %v:%d\n", depth, file, line)
		}
		os.Exit(0)
	}
}

func (g *GBC) setModel(m util.GBModel) {
	g.model = m
	g.Video.Renderer.Model = m
}

func (g *GBC) updateIRQs() {
	irqs := g.IO[IEIO] & g.IO[IFIO] & 0x1f
	if irqs == 0 {
		g.irqPending = 0
		return
	}

	g.halt = false
	if !g.Reg.IME {
		g.irqPending = 0
		return
	}
	if g.irqPending > 0 {
		return
	}

	for i := 0; i < 4; i++ {
		if util.Bit(irqs, i) {
			g.irqPending = i + 1
			return
		}
	}
}

func (g *GBC) Draw() []uint8 { return g.Video.Display().Pix }

func (g *GBC) handleJoypad() {
	pressed := g.joypad.Input()
	if pressed {
		g.IO[IFIO] = util.SetBit8(g.IO[IFIO], 4, true)
		g.updateIRQs()
	}
}

func (g *GBC) Frame() int { return g.Video.FrameCounter }

// _GBMemoryDMAService
func (g *GBC) dmaService(cyclesLate uint64) {
	remaining := g.dma.remaining
	g.dma.remaining = 0
	b := g.Load8(g.dma.src)
	g.Store8(g.dma.dest, b)
	g.dma.src++
	g.dma.dest++
	g.dma.remaining = remaining - 1
	if g.dma.remaining > 0 {
		after := 4 * (2 - util.Bool2U64(g.DoubleSpeed))
		g.scheduler.ScheduleEvent(scheduler.OAMDMA, g.dmaService, after)
	}
}

// _GBMemoryHDMAService
func (g *GBC) hdmaService(cyclesLate uint64) {
	g.cpuBlocked = true

	b := g.Load8(g.hdma.src)
	g.Store8(g.hdma.dest, b)

	g.hdma.src++
	g.hdma.dest++
	g.hdma.remaining--

	if g.hdma.remaining > 0 {
		g.scheduler.DescheduleEvent(scheduler.HDMA)
		g.scheduler.ScheduleEvent(scheduler.HDMA, g.hdmaService, 4)
		return
	}

	g.cpuBlocked = false
	g.IO[HDMA1IO] = byte(g.hdma.src >> 8)
	g.IO[HDMA2IO] = byte(g.hdma.src)
	g.IO[HDMA3IO] = byte(g.hdma.dest >> 8)
	g.IO[HDMA4IO] = byte(g.hdma.dest)
	if g.hdma.enable {
		g.IO[HDMA5IO]--
		if g.IO[HDMA5IO] == 0xff {
			g.hdma.enable = false
		}
	} else {
		g.IO[HDMA5IO] = 0xff
	}
}

func (g *GBC) triggerIRQ(idx int) {
	g.IO[IFIO] = util.SetBit8(g.IO[IFIO], idx, false)
	g.timer.tick(20)
	g.pushPC()
	g.Reg.PC = irqVec[idx]
}

func (g *GBC) hdmaMode3() {
	if g.Video.Ly < video.VERTICAL_PIXELS && g.hdma.enable && g.IO[HDMA5IO] != 0xff {
		g.hdma.remaining = 0x10
		g.cpuBlocked = true
		g.scheduler.DescheduleEvent(scheduler.HDMA)
		g.scheduler.ScheduleEvent(scheduler.HDMA, g.hdmaService, 0)
	}
}
