package emulator

import (
	"fmt"
	"net"
	"sync"

	"gbc/pkg/apu"
	"gbc/pkg/cartridge"
	"gbc/pkg/config"
	"gbc/pkg/debug"
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
	halt      bool // Halt状態か
	Config    *config.Config
	mode      int
	// timer関連
	Timer
	serialTick chan int
	ROMBank
	RAMBank
	WRAMBank
	bankMode uint
	// サウンド
	Sound apu.APU
	// 画面
	GPU gpu.GPU
	// RTC
	RTC   rtc.RTC
	boost int // 倍速か
	// シリアル通信
	Serial serial.Serial

	romdir string // ロムがあるところのディレクトリパス

	IMESwitch
	debug Debug
}

// TransferROM Transfer ROM from cartridge to Memory
func (cpu *CPU) TransferROM(rom []byte) {
	for i := 0x0000; i <= 0x7fff; i++ {
		cpu.RAM[i] = rom[i]
	}

	// カードリッジタイプで場合分け
	switch cpu.Cartridge.Type {
	case 0x00:
		// Type : 0
		cpu.Cartridge.MBC = cartridge.ROM
		cpu.transferROM(2, rom)
	case 0x01:
		// Type : 1 => MBC1
		cpu.Cartridge.MBC = cartridge.MBC1
		switch cpu.Cartridge.ROMSize {
		case 0:
			cpu.transferROM(2, rom)
		case 1:
			cpu.transferROM(4, rom)
		case 2:
			cpu.transferROM(8, rom)
		case 3:
			cpu.transferROM(16, rom)
		case 4:
			cpu.transferROM(32, rom)
		case 5:
			cpu.transferROM(64, rom)
		case 6:
			cpu.transferROM(128, rom)
		default:
			errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x02, 0x03:
		// Type : 2, 3 => MBC1+RAM
		cpu.Cartridge.MBC = cartridge.MBC1
		switch cpu.Cartridge.RAMSize {
		case 0, 1, 2:
			switch cpu.Cartridge.ROMSize {
			case 0:
				cpu.transferROM(2, rom)
			case 1:
				cpu.transferROM(4, rom)
			case 2:
				cpu.transferROM(8, rom)
			case 3:
				cpu.transferROM(16, rom)
			case 4:
				cpu.transferROM(32, rom)
			case 5:
				cpu.transferROM(64, rom)
			case 6:
				cpu.transferROM(128, rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
				panic(errorMsg)
			}
		case 3:
			cpu.bankMode = 1
			switch cpu.Cartridge.ROMSize {
			case 0:
			case 1:
				cpu.transferROM(4, rom)
			case 2:
				cpu.transferROM(8, rom)
			case 3:
				cpu.transferROM(16, rom)
			case 4:
				cpu.transferROM(32, rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
				panic(errorMsg)
			}
		default:
			errorMsg := fmt.Sprintf("RAMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x05, 0x06:
		// Type : 5, 6 => MBC2
		cpu.Cartridge.MBC = cartridge.MBC2
		switch cpu.Cartridge.RAMSize {
		case 0, 1, 2:
			switch cpu.Cartridge.ROMSize {
			case 0:
				cpu.transferROM(2, rom)
			case 1:
				cpu.transferROM(4, rom)
			case 2:
				cpu.transferROM(8, rom)
			case 3:
				cpu.transferROM(16, rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
				panic(errorMsg)
			}
		case 3:
			cpu.bankMode = 1
			switch cpu.Cartridge.ROMSize {
			case 0:
			case 1:
				cpu.transferROM(4, rom)
			case 2:
				cpu.transferROM(8, rom)
			case 3:
				cpu.transferROM(16, rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
				panic(errorMsg)
			}
		default:
			errorMsg := fmt.Sprintf("RAMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x0f, 0x10, 0x11, 0x12, 0x13:
		// Type : 0x0f, 0x10, 0x11, 0x12, 0x13 => MBC3
		cpu.Cartridge.MBC = cartridge.MBC3

		cpu.RTC.Working = true

		switch cpu.Cartridge.ROMSize {
		case 0:
			cpu.transferROM(2, rom)
		case 1:
			cpu.transferROM(4, rom)
		case 2:
			cpu.transferROM(8, rom)
		case 3:
			cpu.transferROM(16, rom)
		case 4:
			cpu.transferROM(32, rom)
		case 5:
			cpu.transferROM(64, rom)
		case 6:
			cpu.transferROM(128, rom)
		default:
			errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
			panic(errorMsg)
		}
	case 0x19, 0x1a, 0x1b:
		// Type : 0x19, 0x1a, 0x1b => MBC5
		cpu.Cartridge.MBC = cartridge.MBC5
		switch cpu.Cartridge.ROMSize {
		case 0:
			cpu.transferROM(2, rom)
		case 1:
			cpu.transferROM(4, rom)
		case 2:
			cpu.transferROM(8, rom)
		case 3:
			cpu.transferROM(16, rom)
		case 4:
			cpu.transferROM(32, rom)
		case 5:
			cpu.transferROM(64, rom)
		case 6:
			cpu.transferROM(128, rom)
		case 7:
			cpu.transferROM(256, rom)
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
	cpu.Reg.AF = 0x11b0 // A=01 => GB, A=11 => CGB
	cpu.Reg.BC = 0x0013
	cpu.Reg.DE = 0x00d8
	cpu.Reg.HL = 0x014d
	cpu.Reg.PC = 0x0100
	cpu.Reg.SP = 0xfffe
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
	cpu.RAM[OBP0IO] = 0xff
	cpu.RAM[OBP1IO] = 0xff
}

func (cpu *CPU) initNetwork() {
	if cpu.Config.Network.Network {
		your := cpu.Config.Network.Your
		peer := cpu.Config.Network.Peer
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
				cpu.setSerialFlag()
			}
		}()
	}
}

func (cpu *CPU) initDMGPalette() {
	color0 := cpu.Config.Pallete.Color0
	color1 := cpu.Config.Pallete.Color1
	color2 := cpu.Config.Pallete.Color2
	color3 := cpu.Config.Pallete.Color3
	gpu.InitPalette(color0, color1, color2, color3)
}

// Init CPU・メモリの初期化
func (cpu *CPU) Init(romdir string, debug bool, test bool) {
	cpu.initRegister()
	cpu.initIOMap()

	cpu.ROMBank.ptr = 1
	cpu.WRAMBank.ptr = 1

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
		cpu.Config.Display.HQ2x = false
		cpu.Config.Display.FPS30 = true
		cpu.debug.history.SetFlag(cpu.Config.Debug.History)
		cpu.debug.Break.ParseBreakpoints(cpu.Config.Debug.BreakPoints)
	}
}

// Exit 後始末を行う
func (cpu *CPU) Exit() {
	cpu.save()
	cpu.Serial.Exit()
}

// Exec 1サイクル
func (cpu *CPU) exec() bool {
	bank, PC := cpu.ROMBank.ptr, cpu.Reg.PC

	bytecode := cpu.FetchMemory8(PC)
	opcode := opcodes[bytecode]
	instruction, operand1, operand2, cycle1, cycle2, handler := opcode.Ins, opcode.Operand1, opcode.Operand2, opcode.Cycle1, opcode.Cycle2, opcode.Handler
	cycle := cycle1

	// in breakpoint
	b := &cpu.debug.Break
	if cpu.debug.on {
		switch b.Flag() {
		case debug.BreakOn:
			return true
		case debug.BreakOff:
			for _, breakpoint := range b.BreakPoints() {
				if PC != breakpoint.PC {
					continue
				}

				if (PC > 0x4000 && bank == breakpoint.Bank) || breakpoint.Bank == 0 {
					if cpu.checkBreakCond(&breakpoint) {
						b.SetFlag(debug.BreakOn)
						cpu.debug.history.SetHistory(bank, PC, bytecode)
						return true
					}
				}
			}
		case debug.BreakDelay:
			b.SetFlag(debug.BreakOff)
		}
	}

	halt := cpu.halt
	if !cpu.halt {
		if cpu.debug.on && cpu.debug.history.Flag() {
			cpu.debug.history.SetHistory(bank, PC, bytecode)
		}

		if handler != nil {
			handler(cpu, operand1, operand2)
		} else {
			switch instruction {
			case INS_HALT:
				cpu.HALT(operand1, operand2)
			case INS_LD:
				LD(cpu, operand1, operand2)
			case INS_LDH:
				LDH(cpu, operand1, operand2)
			case INS_JR:
				JR(cpu, operand1, operand2)
			case INS_NOP:
				cpu.NOP(operand1, operand2)
			case INS_AND:
				cpu.AND(operand1, operand2)
			case INS_INC:
				cpu.INC(operand1, operand2)
			case INS_DEC:
				cpu.DEC(operand1, operand2)
			case INS_PUSH:
				cpu.PUSH(operand1, operand2)
				cycle = 0 // PUSH内部でサイクルのインクリメントを行う
			case INS_POP:
				cpu.POP(operand1, operand2)
				cycle = 0 // POP内部でサイクルのインクリメントを行う
			case INS_XOR:
				cpu.XOR(operand1, operand2)
			case INS_JP:
				JP(cpu, operand1, operand2) // JP内部でサイクルのインクリメントを行う
			case INS_CALL:
				CALL(cpu, operand1, operand2) // CALL内部でサイクルのインクリメントを行う
			case INS_RET:
				if !cpu.RET(operand1, operand2) {
					cycle = cycle2
				}
			case INS_RETI:
				cpu.RETI(operand1, operand2)
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
			case INS_CPL:
				cpu.CPL(operand1, operand2)
			case INS_PREFIX:
				cpu.PREFIXCB(operand1, operand2)
				cycle = 0 // PREFIXCB内部でサイクルのインクリメントを行う
			case INS_RRA:
				cpu.RRA(operand1, operand2)
			case INS_DAA:
				cpu.DAA(operand1, operand2)
			case INS_RST:
				cpu.RST(operand1, operand2)
			case INS_SCF:
				cpu.SCF(operand1, operand2)
			case INS_CCF:
				cpu.CCF(operand1, operand2)
			case INS_RLCA:
				cpu.RLCA(operand1, operand2)
			case INS_RLA:
				cpu.RLA(operand1, operand2)
			case INS_RRCA:
				cpu.RRCA(operand1, operand2)
			case INS_DI:
				cpu.DI(operand1, operand2)
			case INS_EI:
				cpu.EI(operand1, operand2)
			case INS_STOP:
				cpu.STOP(operand1, operand2)
			default:
				errMsg := fmt.Sprintf("eip: 0x%04x opcode: 0x%02x", cpu.Reg.PC, bytecode)
				panic(errMsg)
			}
		}
	} else {
		cycle = 1

		// ref: https://rednex.github.io/rgbds/gbz80.7.html#HALT
		if !cpu.Reg.IME {
			IE, IF := cpu.RAM[IEIO], cpu.RAM[IFIO]
			pending := IE&IF != 0
			if pending {
				cpu.halt = false
			}
		}
	}

	cpu.debug.monitor.Add(halt, cycle)
	cpu.timer(cycle)

	cpu.handleInterrupt()

	return false
}

func (cpu *CPU) execScanline() (scx uint, scy uint, ok bool) {
	// OAM mode2
	cpu.setOAMRAMMode()
	for cpu.Cycle.scanline <= 20*cpu.boost {
		if inBreak := cpu.exec(); inBreak {
			return 0, 0, false
		}
	}

	// LCD Driver mode3
	cpu.Cycle.scanline -= 20 * cpu.boost
	cpu.setLCDMode()
	for cpu.Cycle.scanline <= 42*cpu.boost {
		if inBreak := cpu.exec(); inBreak {
			return 0, 0, false
		}
	}

	scrollX, scrollY := cpu.GPU.GetScroll()

	// HBlank mode0
	cpu.Cycle.scanline -= 42 * cpu.boost
	cpu.setHBlankMode()
	for cpu.Cycle.scanline <= (cyclePerLine-(20+42))*cpu.boost {
		if inBreak := cpu.exec(); inBreak {
			return 0, 0, false
		}
	}
	cpu.Cycle.scanline -= (cyclePerLine - (20 + 42)) * cpu.boost

	cpu.incrementLY()
	return scrollX, scrollY, true
}

// VBlank
func (cpu *CPU) execVBlank() {
	for {
		for cpu.Cycle.scanline < cyclePerLine*cpu.boost {
			if inBreak := cpu.exec(); inBreak {
				return
			}
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
