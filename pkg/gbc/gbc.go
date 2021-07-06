package gbc

import (
	"fmt"
	"math"
	"os"
	"runtime"

	"gbc/pkg/emulator/config"
	"gbc/pkg/emulator/joypad"
	"gbc/pkg/gbc/apu"
	"gbc/pkg/gbc/cart"
	"gbc/pkg/gbc/rtc"
	"gbc/pkg/gbc/scheduler"
	"gbc/pkg/gbc/video"
	"gbc/pkg/util"
)

// ROMBank - 0x4000-0x7fff
type ROMBank struct {
	ptr  uint8
	bank [256][0x4000]byte
}

// RAMBank - 0xa000-0xbfff
type RAMBank struct {
	ptr  uint8
	Bank [16][0x2000]byte
}

// WRAMBank - 0xd000-0xdfff ゲームボーイカラーのみ
type WRAMBank struct {
	ptr  uint8
	bank [8][0x1000]byte
}

type Dma struct {
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
	Reg       Register
	RAM       [0x10000]byte
	IO        [0x100]byte // 0xff00-0xffff
	Cartridge cart.Cartridge
	joypad    joypad.Joypad
	halt      bool
	Config    *config.Config
	cycles    int
	timer     Timer
	ROMBank
	RAMBank
	WRAMBank
	bankMode    uint
	Sound       apu.APU
	video       *video.Video
	RTC         rtc.RTC
	doubleSpeed bool
	IMESwitch
	Debug      Debug
	model      util.GBModel
	irqPending int
	scheduler  *scheduler.Scheduler
	dma        Dma
}

// TransferROM Transfer ROM from cartridge to Memory
func (g *GBC) TransferROM(rom []byte) {
	for i := 0x0000; i <= 0x7fff; i++ {
		g.RAM[i] = rom[i]
	}

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
			g.ROMBank.bank[bank][i] = rom[bank*0x4000+i]
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

// Init g and ram
func (g *GBC) Init(debug bool, test bool) {
	g.scheduler = scheduler.New()
	g.video = video.New(&g.IO, g.updateIRQs, g.scheduler.ScheduleEvent)
	if g.Cartridge.IsCGB {
		g.setModel(util.GB_MODEL_CGB)
	}

	g.resetRegister()
	g.resetIO()

	g.ROMBank.ptr, g.WRAMBank.ptr = 1, 1

	g.Config = config.Init()
	g.doubleSpeed = false

	// Init APU
	g.Sound.Init(!test)

	// Init RTC
	go g.RTC.Init()

	g.Debug.Enable = debug
	if debug {
		g.Config.Display.HQ2x, g.Config.Display.FPS30 = false, true
		g.Debug.history.SetFlag(g.Config.Debug.History)
		g.Debug.Break.ParseBreakpoints(g.Config.Debug.BreakPoints)
	}

	g.scheduler.ScheduleEvent("endMode2", g.video.EndMode2, video.MODE_2_LENGTH)
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
			switch oldIrqPending {
			case VBlankIRQ:
				g.triggerVBlank()
			case LCDCIRQ:
				g.triggerLCDC()
			case TimerIRQ:
				g.triggerTimer()
			case SerialIRQ:
				g.triggerSerial()
			case JoypadIRQ:
				g.triggerJoypad()
			}
			return
		} else if handler != nil {
			handler(g, operand1, operand2)
		} else {
			switch instruction {
			case INS_LDH:
				LDH(g, operand1, operand2)
			case INS_AND:
				g.AND(operand1, operand2)
			case INS_XOR:
				g.XOR(operand1, operand2)
			case INS_CP:
				g.CP(operand1, operand2)
			case INS_OR:
				g.OR(operand1, operand2)
			case INS_ADD:
				g.ADD(operand1, operand2)
			case INS_SUB:
				g.SUB(operand1, operand2)
			case INS_ADC:
				g.ADC(operand1, operand2)
			case INS_SBC:
				g.SBC(operand1, operand2)
			default:
				errMsg := fmt.Sprintf("eip: 0x%04x opcode: 0x%02x", g.Reg.PC, bytecode)
				panic(errMsg)
			}
		}
	} else {
		cycle = int(g.scheduler.Next()-g.scheduler.Cycle()) / 4
		if cycle > 16 { // make sound seemless
			cycle = 16
		}
		if cycle == 0 {
			cycle++
		}
		if !g.Reg.IME { // ref: https://rednex.github.io/rgbds/gbz80.7.html#HALT
			if pending := g.IO[IEIO-0xff00]&g.IO[IFIO-0xff00] > 0; pending {
				g.halt = false
			}
		}
		if pending {
			g.pend()
		}
	}

	g.updateTimer(cycle)
}

func (g *GBC) execScanline() {
	for {
		if g.scheduler.Cycle() < g.scheduler.Next() {
			g.step()
		} else {
			g.scheduler.DoEvent()
			mode := g.video.Mode()
			if mode == 1 || mode == 2 {
				return
			}
		}
	}
}

// VBlank
func (g *GBC) execVBlank() {
	for {
		if g.scheduler.Cycle() < g.scheduler.Next() {
			g.step()
		} else {
			g.scheduler.DoEvent()
			mode := g.video.Mode()
			if mode == 2 {
				return
			}
		}
	}
}

// 1 frame
func (g *GBC) Update() error {
	if g.Frame()%3 == 0 {
		g.handleJoypad()
	}

	p, b := &g.Debug.pause, &g.Debug.Break
	if p.Delay() {
		p.DecrementDelay()
	}
	if p.On() || b.On() {
		return nil
	}

	// 0-143
	for y := 0; y < video.VERTICAL_PIXELS; y++ {
		g.execScanline()
	}

	// 143-154
	g.execVBlank()

	g.video.UpdateFrameCount()
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
	g.video.Renderer.Model = m
}

func (g *GBC) updateIRQs() {
	irqs := g.IO[IEIO-0xff00] & g.IO[IFIO-0xff00] & 0x1f
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

func (g *GBC) Draw() []uint8 {
	display := g.video.Display()
	return display.Pix
}

func (g *GBC) handleJoypad() {
	pad := g.Config.Joypad
	result := g.joypad.Input(pad.A, pad.B, pad.Start, pad.Select, pad.Threshold)
	if result != 0 {
		switch result {
		case joypad.Pressed: // Joypad Interrupt
			if g.Reg.IME && g.getJoypadEnable() {
				g.setJoypadFlag(true)
			}
		}
	}
}

func (g *GBC) Frame() int {
	return g.video.FrameCounter
}

// _GBMemoryDMAService
func (g *GBC) DMAService() {
	remaining := g.dma.remaining
	g.dma.remaining = 0
	b := g.Load8(g.dma.src)
	g.Store8(g.dma.dest, b)
	g.dma.src++
	g.dma.dest++
	g.dma.remaining = remaining - 1
	if g.dma.remaining > 0 {
		after := 4 * (2 - util.Bool2U64(g.doubleSpeed))
		g.scheduler.ScheduleEvent("oamdma", g.DMAService, after)
	}
}
