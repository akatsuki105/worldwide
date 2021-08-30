package gbc

import (
	"fmt"
	"math"
	"os"
	"runtime"

	"github.com/pokemium/worldwide/pkg/gbc/apu"
	"github.com/pokemium/worldwide/pkg/gbc/cart"
	"github.com/pokemium/worldwide/pkg/gbc/joypad"
	"github.com/pokemium/worldwide/pkg/gbc/rtc"
	"github.com/pokemium/worldwide/pkg/gbc/scheduler"
	"github.com/pokemium/worldwide/pkg/gbc/video"
	"github.com/pokemium/worldwide/pkg/util"
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

type CurInst struct {
	Opcode byte
	PC     uint16
}

// GBC core structure
type GBC struct {
	Reg  Register
	Inst CurInst

	// memory
	ROM  ROM
	RAM  RAM
	WRAM WRAM
	IO   [0x100]byte // 0xff00-0xffff

	Cartridge   *cart.Cartridge
	joypad      *joypad.Joypad
	Halt        bool
	timer       *Timer
	bankMode    uint
	Sound       *apu.APU
	Video       *video.Video
	RTC         *rtc.RTC
	DoubleSpeed bool
	model       util.GBModel
	irqPending  int
	scheduler   *scheduler.Scheduler
	dma         Dma
	hdma        Hdma
	cpuBlocked  bool

	// plugins
	Callbacks []*util.Callback
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

func New(romData []byte, j [8](func() bool), setAudioStream func([]byte)) *GBC {
	c := cart.New(romData)
	g := &GBC{
		Cartridge: c,
		scheduler: scheduler.New(),
		joypad:    joypad.New(j),
		RTC:       rtc.New(c.HasRTC()),
		Sound:     apu.New(true, setAudioStream),
	}

	// init graphics
	g.Video = video.New(&g.IO, g.updateIRQs, g.hdmaMode3, g.scheduler)
	if g.Cartridge.IsCGB {
		g.setModel(util.GB_MODEL_CGB)
	}

	// init timer
	g.timer = NewTimer(g)
	g.skipBIOS()
	g.TransferROM(romData)
	return g
}

// Exec 1cycle
func (g *GBC) Step() {
	pc := g.Reg.PC
	opcode := g.Load8(pc)
	g.Inst.Opcode, g.Inst.PC = opcode, pc

	inst := gbz80insts[opcode]
	operand1, operand2, cycle, handler := inst.Operand1, inst.Operand2, inst.Cycle1, inst.Handler

	if g.Halt || g.cpuBlocked {
		cycle = int(g.scheduler.Next() - g.scheduler.Cycle())
	} else {
		if g.irqPending > 0 {
			oldIrqPending := g.irqPending
			g.irqPending = 0
			g.setInterrupts(false)
			g.triggerIRQ(int(oldIrqPending - 1))
			return
		}

		g.Reg.PC++
		handler(g, operand1, operand2)
		cycle *= (4 >> uint32(util.Bool2U64(g.DoubleSpeed)))
	}

	g.timer.tick(uint32(cycle))
}

// 1 frame
func (g *GBC) Update() {
	frame := g.Frame()
	if frame%3 == 0 {
		g.handleJoypad()
	}

	for frame == g.Video.FrameCounter {
		g.Step()

		// trigger callbacks
		for _, callback := range g.Callbacks {
			if callback.Func() {
				break
			}
		}
	}

	g.Sound.Update()
}

func (g *GBC) PanicHandler(place string, stack bool) {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "%s emulation error: %s in 0x%04x\n", place, err, g.Reg.PC)
		for depth := 0; ; depth++ {
			_, file, line, ok := runtime.Caller(depth)
			if !ok {
				break
			}
			fmt.Fprintf(os.Stderr, "======> %d: %v:%d\n", depth, file, line)
		}
		os.Exit(1)
	}
}

func (g *GBC) setModel(m util.GBModel) {
	g.model = m
	g.Video.Renderer.Model = m
}

// GBUpdateIRQs
func (g *GBC) updateIRQs() {
	irqs := g.IO[IEIO] & g.IO[IFIO] & 0x1f
	if irqs == 0 {
		g.irqPending = 0
		return
	}

	g.Halt = false
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

func (g *GBC) Draw() []byte { return g.Video.Display().Pix }

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
		g.scheduler.ScheduleEvent(scheduler.OAMDMA, g.dmaService, (4>>util.Bool2U64(g.DoubleSpeed))-cyclesLate) // 4 * 40 = 160cycle
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
		g.scheduler.ScheduleEvent(scheduler.HDMA, g.hdmaService, 4-cyclesLate)
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

func (g *GBC) skipBIOS() {
	g.resetRegister()
	g.resetIO()
	g.ROM.bank, g.WRAM.bank = 1, 1

	g.storeIO(LCDCIO, 0x91)
	g.Video.SkipBIOS()
}

// GBSetInterrupts
func (g *GBC) setInterrupts(enable bool) {
	g.scheduler.DescheduleEvent(scheduler.EiPending)
	if enable {
		g.scheduler.ScheduleEvent(scheduler.EiPending, func(_ uint64) {
			g.Reg.IME = true
			g.updateIRQs()
		}, 8>>util.Bool2Int(g.DoubleSpeed))
		return
	}

	g.Reg.IME = false
	g.updateIRQs()
}
