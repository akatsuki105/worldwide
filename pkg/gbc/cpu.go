package gbc

import (
	"fmt"
	"math"
	"net"
	"sync"

	"gbc/pkg/apu"
	"gbc/pkg/cartridge"
	"gbc/pkg/config"
	"gbc/pkg/gpu"
	"gbc/pkg/joypad"
	"gbc/pkg/rtc"
	"gbc/pkg/serial"
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

// CPU Central Processing Unit
type CPU struct {
	Reg       Register
	RAM       [0x10000]byte
	Cartridge cartridge.Cartridge
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
func (cpu *CPU) TransferROM(rom []byte) {
	for i := 0x0000; i <= 0x7fff; i++ {
		cpu.RAM[i] = rom[i]
	}

	switch cpu.Cartridge.Type {
	case 0x00:
		cpu.Cartridge.MBC = cartridge.ROM
		cpu.transferROM(2, rom)
	case 0x01: // Type : 1 => MBC1
		cpu.Cartridge.MBC = cartridge.MBC1
		switch r := int(cpu.Cartridge.ROMSize); r {
		case 0, 1, 2, 3, 4, 5, 6:
			cpu.transferROM(int(math.Pow(2, float64(r+1))), rom)
		default:
			errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x02, 0x03: // Type : 2, 3 => MBC1+RAM
		cpu.Cartridge.MBC = cartridge.MBC1
		switch cpu.Cartridge.RAMSize {
		case 0, 1, 2:
			switch r := int(cpu.Cartridge.ROMSize); r {
			case 0, 1, 2, 3, 4, 5, 6:
				cpu.transferROM(int(math.Pow(2, float64(r+1))), rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
				panic(errorMsg)
			}
		case 3:
			cpu.bankMode = 1
			switch r := int(cpu.Cartridge.ROMSize); r {
			case 0:
			case 1, 2, 3, 4:
				cpu.transferROM(int(math.Pow(2, float64(r+1))), rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
				panic(errorMsg)
			}
		default:
			errorMsg := fmt.Sprintf("RAMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x05, 0x06: // Type : 5, 6 => MBC2
		cpu.Cartridge.MBC = cartridge.MBC2
		switch cpu.Cartridge.RAMSize {
		case 0, 1, 2:
			switch r := int(cpu.Cartridge.ROMSize); r {
			case 0, 1, 2, 3:
				cpu.transferROM(int(math.Pow(2, float64(r+1))), rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
				panic(errorMsg)
			}
		case 3:
			cpu.bankMode = 1
			switch r := int(cpu.Cartridge.ROMSize); r {
			case 0:
			case 1, 2, 3:
				cpu.transferROM(int(math.Pow(2, float64(r+1))), rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
				panic(errorMsg)
			}
		default:
			errorMsg := fmt.Sprintf("RAMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x0f, 0x10, 0x11, 0x12, 0x13: // Type : 0x0f, 0x10, 0x11, 0x12, 0x13 => MBC3
		cpu.Cartridge.MBC, cpu.RTC.Enable = cartridge.MBC3, true
		switch r := int(cpu.Cartridge.ROMSize); r {
		case 0, 1, 2, 3, 4, 5, 6:
			cpu.transferROM(int(math.Pow(2, float64(r+1))), rom)
		default:
			errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x19, 0x1a, 0x1b: // Type : 0x19, 0x1a, 0x1b => MBC5
		cpu.Cartridge.MBC = cartridge.MBC5
		switch r := int(cpu.Cartridge.ROMSize); r {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			cpu.transferROM(int(math.Pow(2, float64(r+1))), rom)
		default:
			errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
			panic(errorMsg)
		}
	default:
		errorMsg := fmt.Sprintf("Type is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
		panic(errorMsg)
	}
}

func (cpu *CPU) transferROM(bankNum int, rom []byte) {
	for bank := 0; bank < bankNum; bank++ {
		for i := 0x0000; i <= 0x3fff; i++ {
			cpu.ROMBank.bank[bank][i] = rom[bank*0x4000+i]
		}
	}
}

func (cpu *CPU) initRegister() {
	cpu.Reg.setAF(0x11b0) // A=01 => GB, A=11 => CGB
	cpu.Reg.setBC(0x0013)
	cpu.Reg.setDE(0x00d8)
	cpu.Reg.setHL(0x014d)
	cpu.Reg.PC, cpu.Reg.SP = 0x0100, 0xfffe
}

func (cpu *CPU) initIOMap() {
	cpu.RAM[0xff04] = 0x1e
	cpu.RAM[0xff05] = 0x00
	cpu.RAM[0xff06] = 0x00
	cpu.RAM[0xff07] = 0xf8
	cpu.RAM[0xff0f] = 0xe1
	cpu.RAM[0xff10] = 0x80
	cpu.RAM[0xff11] = 0xbf
	cpu.RAM[0xff12] = 0xf3
	cpu.RAM[0xff14] = 0xbf
	cpu.RAM[0xff16] = 0x3f
	cpu.RAM[0xff17] = 0x00
	cpu.RAM[0xff19] = 0xbf
	cpu.RAM[0xff1a] = 0x7f
	cpu.RAM[0xff1b] = 0xff
	cpu.RAM[0xff1c] = 0x9f
	cpu.RAM[0xff1e] = 0xbf
	cpu.RAM[0xff20] = 0xff
	cpu.RAM[0xff21] = 0x00
	cpu.RAM[0xff22] = 0x00
	cpu.RAM[0xff23] = 0xbf
	cpu.RAM[0xff24] = 0x77
	cpu.RAM[0xff25] = 0xf3
	cpu.RAM[0xff26] = 0xf1
	cpu.SetMemory8(LCDCIO, 0x91)
	cpu.SetMemory8(LCDSTATIO, 0x85)
	cpu.RAM[BGPIO] = 0xfc
	cpu.RAM[OBP0IO], cpu.RAM[OBP1IO] = 0xff, 0xff
}

func (cpu *CPU) initNetwork() {
	if cpu.Config.Network.Network {
		your, peer := cpu.Config.Network.Your, cpu.Config.Network.Peer
		myIP, myPort, _ := net.SplitHostPort(your)
		peerIP, peerPort, _ := net.SplitHostPort(peer)
		received := make(chan int)
		cpu.Serial.Init(myIP, myPort, peerIP, peerPort, received, &cpu.mutex)
		cpu.serialTick = make(chan int)

		go func() {
			for {
				<-received
				cpu.Serial.TransferFlag = 1
				<-cpu.serialTick
				cpu.Serial.Receive()
				cpu.Serial.ClearSC()
				cpu.setSerialFlag(true)
			}
		}()
	}
}

func (cpu *CPU) initDMGPalette() {
	c0, c1, c2, c3 := cpu.Config.Palette.Color0, cpu.Config.Palette.Color1, cpu.Config.Palette.Color2, cpu.Config.Palette.Color3
	gpu.InitPalette(c0, c1, c2, c3)
}

// Init cpu and ram
func (cpu *CPU) Init(romdir string, debug bool, test bool) {
	cpu.initRegister()
	cpu.initIOMap()

	cpu.ROMBank.ptr, cpu.WRAMBank.ptr = 1, 1

	cpu.GPU.Init(debug)
	cpu.Config = config.Init()
	cpu.boost = 1

	cpu.initNetwork()

	if !cpu.Cartridge.IsCGB {
		cpu.initDMGPalette()
	}

	// load save data
	cpu.romdir = romdir
	cpu.load()

	// Init APU
	cpu.Sound.Init(!test)

	// Init RTC
	go cpu.RTC.Init()

	cpu.debug.on = debug
	if debug {
		cpu.Config.Display.HQ2x, cpu.Config.Display.FPS30 = false, true
		cpu.debug.history.SetFlag(cpu.Config.Debug.History)
		cpu.debug.Break.ParseBreakpoints(cpu.Config.Debug.BreakPoints)
	}
}

// Exit gbc
func (cpu *CPU) Exit() {
	cpu.save()
	cpu.Serial.Exit()
}

// Exec 1cycle
func (cpu *CPU) exec(max int) {
	bank, PC := cpu.ROMBank.ptr, cpu.Reg.PC

	bytecode := cpu.FetchMemory8(PC)
	opcode := opcodes[bytecode]
	instruction, operand1, operand2, cycle, handler := opcode.Ins, opcode.Operand1, opcode.Operand2, opcode.Cycle1, opcode.Handler

	if !cpu.halt {
		if cpu.debug.on && cpu.debug.history.Flag() {
			cpu.debug.history.SetHistory(bank, PC, bytecode)
		}

		if handler != nil {
			handler(cpu, operand1, operand2)
		} else {
			switch instruction {
			case INS_LDH:
				LDH(cpu, operand1, operand2)
			case INS_AND:
				cpu.AND(operand1, operand2)
			case INS_XOR:
				cpu.XOR(operand1, operand2)
			case INS_CP:
				cpu.CP(operand1, operand2)
			case INS_OR:
				cpu.OR(operand1, operand2)
			case INS_ADD:
				cpu.ADD(operand1, operand2)
			case INS_SUB:
				cpu.SUB(operand1, operand2)
			case INS_ADC:
				cpu.ADC(operand1, operand2)
			case INS_SBC:
				cpu.SBC(operand1, operand2)
			default:
				errMsg := fmt.Sprintf("eip: 0x%04x opcode: 0x%02x", cpu.Reg.PC, bytecode)
				panic(errMsg)
			}
		}
	} else {
		cycle = (max - cpu.Cycle.scanline)
		if cycle > 16 { // make sound seemless
			cycle = 16
		}
		if cycle == 0 {
			cycle++
		}
		if !cpu.Reg.IME { // ref: https://rednex.github.io/rgbds/gbz80.7.html#HALT
			IE, IF := cpu.RAM[IEIO], cpu.RAM[IFIO]
			if pending := IE&IF > 0; pending {
				cpu.halt = false
			}
		}
		if pending {
			cpu.pend()
		}
	}

	cpu.timer(cycle)
	cpu.handleInterrupt()
}

func (cpu *CPU) execScanline() (scx uint, scy uint, ok bool) {
	// OAM mode2
	cpu.setOAMRAMMode()
	for cpu.Cycle.scanline <= 20*cpu.boost {
		cpu.exec(20 * cpu.boost)
	}

	// LCD Driver mode3
	cpu.Cycle.scanline -= 20 * cpu.boost
	cpu.setLCDMode()
	for cpu.Cycle.scanline <= 42*cpu.boost {
		cpu.exec(42 * cpu.boost)
	}

	scrollX, scrollY := uint(cpu.GPU.Scroll[0]), uint(cpu.GPU.Scroll[1])

	// HBlank mode0
	cpu.Cycle.scanline -= 42 * cpu.boost
	cpu.setHBlankMode()
	for cpu.Cycle.scanline <= (cyclePerLine-(20+42))*cpu.boost {
		cpu.exec((cyclePerLine - (20 + 42)) * cpu.boost)
	}
	cpu.Cycle.scanline -= (cyclePerLine - (20 + 42)) * cpu.boost

	cpu.incrementLY()
	return scrollX, scrollY, true
}

// VBlank
func (cpu *CPU) execVBlank() {
	for {
		for cpu.Cycle.scanline < cyclePerLine*cpu.boost {
			cpu.exec(cyclePerLine * cpu.boost)
		}
		cpu.incrementLY()
		LY := cpu.FetchMemory8(LYIO)
		if LY == 0 {
			break
		}
		cpu.Cycle.scanline = 0
	}
	cpu.Cycle.scanline = 0
}

func (cpu *CPU) isBoost() bool {
	return cpu.boost > 1
}
