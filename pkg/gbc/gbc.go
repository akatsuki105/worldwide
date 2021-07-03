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
	"gbc/pkg/gbc/gpu"
	"gbc/pkg/gbc/rtc"
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

// GBC core structure
type GBC struct {
	Reg       Register
	RAM       [0x10000]byte
	IO        [0x100]byte // 0xff00-0xffff
	Cartridge cart.Cartridge
	joypad    joypad.Joypad
	halt      bool
	Config    *config.Config
	Timer
	ROMBank
	RAMBank
	WRAMBank
	bankMode uint
	Sound    apu.APU
	GPU      *gpu.GPU
	RTC      rtc.RTC
	boost    int // 1 or 2
	IMESwitch
	Debug Debug
	model util.GBModel
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

func (g *GBC) initRegister() {
	g.Reg.setAF(0x11b0) // A=01 => GB, A=11 => CGB
	g.Reg.setBC(0x0013)
	g.Reg.setDE(0x00d8)
	g.Reg.setHL(0x014d)
	g.Reg.PC, g.Reg.SP = 0x0100, 0xfffe
}

func (g *GBC) initIOMap() {
	model := g.GPU.Renderer.Model

	g.IO[0x04] = 0x1e
	g.IO[0x05] = 0x00
	g.IO[0x06] = 0x00
	g.IO[0x07] = 0xf8
	g.IO[0x0f] = 0xe1
	g.IO[0x10] = 0x80
	g.IO[0x11] = 0xbf
	g.IO[0x12] = 0xf3
	g.IO[0x14] = 0xbf
	g.IO[0x16] = 0x3f
	g.IO[0x17] = 0x00
	g.IO[0x19] = 0xbf
	g.IO[0x1a] = 0x7f
	g.IO[0x1b] = 0xff
	g.IO[0x1c] = 0x9f
	g.IO[0x1e] = 0xbf
	g.IO[0x20] = 0xff
	g.IO[0x21] = 0x00
	g.IO[0x22] = 0x00
	g.IO[0x23] = 0xbf
	g.IO[0x24] = 0x77
	g.IO[0x25] = 0xf3
	g.IO[0x26] = 0xf1
	g.Store8(LCDCIO, 0x91)
	g.Store8(LCDSTATIO, 0x85)
	g.Store8(BGPIO, 0xfc)
	if model < util.GB_MODEL_CGB {
		g.Store8(OBP0IO, 0xff)
		g.Store8(OBP1IO, 0xff)
	}
	if model&util.GB_MODEL_CGB != 0 {
		g.Store8(VBKIO, 0)
		g.Store8(BCPSIO, 0x80)
		g.Store8(OCPSIO, 0)
		g.Store8(SVBKIO, 1)
	}
}

// Init g and ram
func (g *GBC) Init(debug bool, test bool) {
	g.GPU = gpu.New(&g.IO)
	if g.Cartridge.IsCGB {
		g.setModel(util.GB_MODEL_CGB)
	}

	g.initRegister()
	g.initIOMap()

	g.ROMBank.ptr, g.WRAMBank.ptr = 1, 1

	g.Config = config.Init()
	g.boost = 1

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
}

// Exit gbc
func (g *GBC) Exit() {
}

// Exec 1cycle
func (g *GBC) exec(max int) {
	bank, PC := g.ROMBank.ptr, g.Reg.PC

	bytecode := g.Load8(PC)
	opcode := opcodes[bytecode]
	instruction, operand1, operand2, cycle, handler := opcode.Ins, opcode.Operand1, opcode.Operand2, opcode.Cycle1, opcode.Handler

	if !g.halt {
		if g.Debug.Enable && g.Debug.history.Flag() {
			g.Debug.history.SetHistory(bank, PC, bytecode)
		}

		if handler != nil {
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
		cycle = (max - g.Cycle.scanline)
		if cycle > 16 { // make sound seemless
			cycle = 16
		}
		if cycle == 0 {
			cycle++
		}
		if !g.Reg.IME { // ref: https://rednex.github.io/rgbds/gbz80.7.html#HALT
			IE, IF := g.IO[IEIO-0xff00], g.IO[IFIO-0xff00]
			if pending := IE&IF > 0; pending {
				g.halt = false
			}
		}
		if pending {
			g.pend()
		}
	}

	g.timer(cycle)
	g.handleInterrupt()
}

func (g *GBC) execScanline() (scx uint, scy uint, ok bool) {
	// OAM mode2
	for g.Cycle.scanline <= 20*g.boost {
		g.exec(20 * g.boost)
	}
	g.GPU.EndMode2()

	// LCD Driver mode3
	g.Cycle.scanline -= 20 * g.boost
	for g.Cycle.scanline <= 42*g.boost {
		g.exec(42 * g.boost)
	}
	g.GPU.EndMode3(0)

	scrollX, scrollY := uint(g.IO[SCXIO-0xff00]), uint(g.IO[SCYIO-0xff00])

	// HBlank mode0
	g.Cycle.scanline -= 42 * g.boost
	for g.Cycle.scanline <= (cyclePerLine-(20+42))*g.boost {
		g.exec((cyclePerLine - (20 + 42)) * g.boost)
	}
	g.GPU.EndMode0()
	g.Cycle.scanline -= (cyclePerLine - (20 + 42)) * g.boost

	LY := g.Load8(LYIO)
	if LY == 144 { // set vblank flag
		g.setVBlankFlag(true)
	}
	g.checkLYC(LY)
	return scrollX, scrollY, true
}

// VBlank
func (g *GBC) execVBlank() {
	for {
		for g.Cycle.scanline < cyclePerLine*g.boost {
			g.exec(cyclePerLine * g.boost)
		}
		g.GPU.EndMode1()
		LY := g.Load8(LYIO)
		if LY == 144 { // set vblank flag
			g.setVBlankFlag(true)
		}
		g.checkLYC(LY)

		if LY == 0 {
			break
		}
		g.Cycle.scanline = 0
	}
	g.Cycle.scanline = 0
}

func (g *GBC) isBoost() bool {
	return g.boost > 1
}

func (g *GBC) Update() error {
	if frames == 0 {
		g.Debug.monitor.GBC.Reset()
	}
	if frames%3 == 0 {
		g.handleJoypad()
	}

	frames++
	g.Debug.monitor.GBC.Reset()

	p, b := &g.Debug.pause, &g.Debug.Break
	if p.Delay() {
		p.DecrementDelay()
	}
	if p.On() || b.On() {
		return nil
	}

	skipRender = (g.Config.Display.FPS30) && (frames%2 == 1)

	scrollX := uint(g.IO[SCXIO-0xff00])
	scrollPixelX := scrollX % 8

	iterX, iterY := width, height
	if scrollPixelX > 0 {
		iterX += 8
	}

	// render bg and run g
	LCDC1 := [144]bool{}
	for y := 0; y < iterY; y++ {
		scx, _, ok := g.execScanline()
		if !ok {
			break
		}
		scrollX = scx

		if y < height {
			LCDC1[y] = util.Bit(g.Load8(LCDCIO), 1)
		}
	}

	g.execVBlank()
	g.GPU.UpdateFrameCount()
	if g.Debug.Enable {
		select {
		case <-second:
			fps = frames
			frames = 0
		default:
		}
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
	g.GPU.Renderer.Model = m
}

func (g *GBC) checkLYC(LY uint8) {
	LYC := g.Load8(LYCIO)
	if LYC == LY {
		stat := g.Load8(LCDSTATIO) | 0x04 // set lyc flag
		g.Store8(LCDSTATIO, stat)

		if util.Bit(stat, 6) { // trigger LYC=LY interrupt
			g.setLCDSTATFlag(true)
		}
		return
	}

	stat := g.Load8(LCDSTATIO) & 0b11111011 // clear lyc flag
	g.Store8(LCDSTATIO, stat)
}
