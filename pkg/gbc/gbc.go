package gbc

import (
	"fmt"
	"math"
	"net"
	"sync"

	"gbc/pkg/emulator/config"
	"gbc/pkg/emulator/joypad"
	"gbc/pkg/gbc/apu"
	"gbc/pkg/gbc/cart"
	"gbc/pkg/gbc/gpu"
	"gbc/pkg/gbc/rtc"
	"gbc/pkg/gbc/serial"
)

const (
	HBlankMode = iota
	VBlankMode
	OAMRAMMode
	LCDMode
)

// ROMBank - 0x4000-0x7fff
type ROMBank struct {
	ptr  uint8
	bank [256][0x4000]byte
}

// RAMBank - 0xa000-0xbfff
type RAMBank struct {
	ptr  uint8
	bank [16][0x2000]byte
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
	Cartridge cart.Cartridge
	mutex     sync.Mutex
	joypad    joypad.Joypad
	halt      bool
	Config    *config.Config
	mode      int
	Timer
	serialTick chan int
	ROMBank
	RAMBank
	WRAMBank
	bankMode uint
	Sound    apu.APU
	GPU      gpu.GPU
	RTC      rtc.RTC
	boost    int // 1 or 2
	Serial   serial.Serial
	romdir   string // dir path where rom exists
	IMESwitch
	debug Debug
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
	g.RAM[0xff04] = 0x1e
	g.RAM[0xff05] = 0x00
	g.RAM[0xff06] = 0x00
	g.RAM[0xff07] = 0xf8
	g.RAM[0xff0f] = 0xe1
	g.RAM[0xff10] = 0x80
	g.RAM[0xff11] = 0xbf
	g.RAM[0xff12] = 0xf3
	g.RAM[0xff14] = 0xbf
	g.RAM[0xff16] = 0x3f
	g.RAM[0xff17] = 0x00
	g.RAM[0xff19] = 0xbf
	g.RAM[0xff1a] = 0x7f
	g.RAM[0xff1b] = 0xff
	g.RAM[0xff1c] = 0x9f
	g.RAM[0xff1e] = 0xbf
	g.RAM[0xff20] = 0xff
	g.RAM[0xff21] = 0x00
	g.RAM[0xff22] = 0x00
	g.RAM[0xff23] = 0xbf
	g.RAM[0xff24] = 0x77
	g.RAM[0xff25] = 0xf3
	g.RAM[0xff26] = 0xf1
	g.SetMemory8(LCDCIO, 0x91)
	g.SetMemory8(LCDSTATIO, 0x85)
	g.RAM[BGPIO] = 0xfc
	g.RAM[OBP0IO], g.RAM[OBP1IO] = 0xff, 0xff
}

func (g *GBC) initNetwork() {
	if g.Config.Network.Network {
		your, peer := g.Config.Network.Your, g.Config.Network.Peer
		myIP, myPort, _ := net.SplitHostPort(your)
		peerIP, peerPort, _ := net.SplitHostPort(peer)
		received := make(chan int)
		g.Serial.Init(myIP, myPort, peerIP, peerPort, received, &g.mutex)
		g.serialTick = make(chan int)

		go func() {
			for {
				<-received
				g.Serial.TransferFlag = 1
				<-g.serialTick
				g.Serial.Receive()
				g.Serial.ClearSC()
				g.setSerialFlag(true)
			}
		}()
	}
}

func (g *GBC) initDMGPalette() {
	c0, c1, c2, c3 := g.Config.Palette.Color0, g.Config.Palette.Color1, g.Config.Palette.Color2, g.Config.Palette.Color3
	gpu.InitPalette(c0, c1, c2, c3)
}

// Init g and ram
func (g *GBC) Init(romdir string, debug bool, test bool) {
	g.initRegister()
	g.initIOMap()

	g.ROMBank.ptr, g.WRAMBank.ptr = 1, 1

	g.GPU.Init(debug)
	g.Config = config.Init()
	g.boost = 1

	g.initNetwork()

	if !g.Cartridge.IsCGB {
		g.initDMGPalette()
	}

	// load save data
	g.romdir = romdir
	g.load()

	// Init APU
	g.Sound.Init(!test)

	// Init RTC
	go g.RTC.Init()

	g.debug.on = debug
	if debug {
		g.Config.Display.HQ2x, g.Config.Display.FPS30 = false, true
		g.debug.history.SetFlag(g.Config.Debug.History)
		g.debug.Break.ParseBreakpoints(g.Config.Debug.BreakPoints)
	}
}

// Exit gbc
func (g *GBC) Exit() {
	g.save()
	g.Serial.Exit()
}

// Exec 1cycle
func (g *GBC) exec(max int) {
	bank, PC := g.ROMBank.ptr, g.Reg.PC

	bytecode := g.FetchMemory8(PC)
	opcode := opcodes[bytecode]
	instruction, operand1, operand2, cycle, handler := opcode.Ins, opcode.Operand1, opcode.Operand2, opcode.Cycle1, opcode.Handler

	if !g.halt {
		if g.debug.on && g.debug.history.Flag() {
			g.debug.history.SetHistory(bank, PC, bytecode)
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
			IE, IF := g.RAM[IEIO], g.RAM[IFIO]
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
	g.setOAMRAMMode()
	for g.Cycle.scanline <= 20*g.boost {
		g.exec(20 * g.boost)
	}

	// LCD Driver mode3
	g.Cycle.scanline -= 20 * g.boost
	g.setLCDMode()
	for g.Cycle.scanline <= 42*g.boost {
		g.exec(42 * g.boost)
	}

	scrollX, scrollY := uint(g.GPU.Scroll[0]), uint(g.GPU.Scroll[1])

	// HBlank mode0
	g.Cycle.scanline -= 42 * g.boost
	g.setHBlankMode()
	for g.Cycle.scanline <= (cyclePerLine-(20+42))*g.boost {
		g.exec((cyclePerLine - (20 + 42)) * g.boost)
	}
	g.Cycle.scanline -= (cyclePerLine - (20 + 42)) * g.boost

	g.incrementLY()
	return scrollX, scrollY, true
}

// VBlank
func (g *GBC) execVBlank() {
	for {
		for g.Cycle.scanline < cyclePerLine*g.boost {
			g.exec(cyclePerLine * g.boost)
		}
		g.incrementLY()
		LY := g.FetchMemory8(LYIO)
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
