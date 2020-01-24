package emulator

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"gopkg.in/ini.v1"

	"gbc/pkg/apu"
	"gbc/pkg/boot"
	"gbc/pkg/cartridge"
	"gbc/pkg/config"
	"gbc/pkg/gpu"
	"gbc/pkg/joypad"
	"gbc/pkg/rtc"
	"gbc/pkg/serial"
)

// CPU Central Processing Unit
type CPU struct {
	Reg       Register
	RAM       [0x10000]byte
	Cartridge cartridge.Cartridge
	mutex     sync.Mutex
	history   []string
	joypad    joypad.Joypad
	halt      bool // Halt状態か
	config    *ini.File
	// timer関連
	cycle       float64 // タイマー用
	cycleDIV    float64 // DIVタイマー用
	cycleLine   float64 // スキャンライン用
	cycleSerial float64
	serialTick  chan int
	// ROM bank
	ROMBankPtr uint8
	ROMBank    [256][0x4000]byte // 0x4000-0x7fff
	// RAM bank
	RAMBankPtr uint8
	RAMBank    [16][0x2000]byte // 0xa000-0xbfff
	// WRAM bank
	WRAMBankPtr uint8
	WRAMBank    [8][0x1000]byte // 0xd000-0xdfff ゲームボーイカラーのみ
	bankMode    uint
	// サウンド
	Sound apu.APU
	// 画面
	GPU    gpu.GPU
	expand uint
	saving bool // trueだとリフレッシュレートが30に落ちるがCPU使用率は下がる
	smooth bool // pixelglのsmoothモードの有無
	// RTC
	RTC       rtc.RTC
	isBoosted bool // 倍速か
	// シリアル通信
	Serial  serial.Serial
	network bool

	romdir string    // ロムがあるところのディレクトリパス
	Boot   boot.Boot // ブートデータ
}

// TransferROM Transfer ROM from cartridge to Memory
func (cpu *CPU) TransferROM(rom *[]byte) {
	for i := 0x0000; i <= 0x7fff; i++ {
		cpu.RAM[i] = (*rom)[i]
	}

	// カードリッジタイプで場合分け
	switch cpu.Cartridge.Type {
	case 0x00:
		// Type : 0
		cpu.Cartridge.MBC = "ROM"
		cpu.transferROM(2, rom)
	case 0x01:
		// Type : 1 => MBC1
		cpu.Cartridge.MBC = "MBC1"
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
		cpu.Cartridge.MBC = "MBC1"
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
		cpu.Cartridge.MBC = "MBC2"
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
		cpu.Cartridge.MBC = "MBC3"

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
		cpu.Cartridge.MBC = "MBC5"
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

func (cpu *CPU) transferROM(bankNum int, rom *[]byte) {
	for bank := 0; bank < bankNum; bank++ {
		for i := 0x0000; i <= 0x3fff; i++ {
			cpu.ROMBank[bank][i] = (*rom)[bank*0x4000+i]
		}
	}
}

// Init CPU・メモリの初期化
func (cpu *CPU) Init(romdir string) {
	cpu.Reg.AF = 0x11b0 // A=01 => GB, A=11 => CGB
	cpu.Reg.BC = 0x0013
	cpu.Reg.DE = 0x00d8
	cpu.Reg.HL = 0x014d
	cpu.Reg.PC = 0x0100
	cpu.Reg.SP = 0xfffe

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

	cpu.ROMBankPtr = 1
	cpu.WRAMBankPtr = 1

	cpu.GPU.Init()
	cpu.config = config.Init()

	expand, err := cpu.config.Section("display").Key("expand").Uint()
	if err != nil {
		cpu.expand = 1
	} else {
		cpu.expand = expand
	}

	saving, err := cpu.config.Section("display").Key("saving").Bool()
	if err != nil {
		saving = false
	}
	cpu.saving = saving

	smooth, err := cpu.config.Section("display").Key("smooth").Bool()
	if err != nil {
		smooth = false
	}
	cpu.smooth = smooth

	network, err := cpu.config.Section("network").Key("network").Bool()
	if err != nil {
		network = false
	}
	cpu.network = network
	if network {
		your := cpu.config.Section("network").Key("your").MustString("127.0.0.1:8888")
		peer := cpu.config.Section("network").Key("peer").MustString("127.0.0.1:9999")
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

	if !cpu.Cartridge.IsCGB {
		color0 := cpu.config.Section("pallete").Key("color0").Ints(",")
		color1 := cpu.config.Section("pallete").Key("color1").Ints(",")
		color2 := cpu.config.Section("pallete").Key("color2").Ints(",")
		color3 := cpu.config.Section("pallete").Key("color3").Ints(",")
		cpu.GPU.InitPallete(color0, color1, color2, color3)
	}

	// load save data
	cpu.romdir = romdir
	cpu.load()

	// Init BootData
	{
		valid, _ := cpu.config.Section("display").Key("boot").Bool()
		cpu.Boot = *boot.New(valid)
	}

	// Init RTC
	go cpu.RTC.Init()
}

// Exit 後始末を行う
func (cpu *CPU) Exit() {
	cpu.save()
	cpu.Serial.Exit()
}

// Exec 1サイクル
func (cpu *CPU) exec() {
	cpu.mutex.Lock()
	opcode := cpu.FetchMemory8(cpu.Reg.PC)
	instruction, operand1, operand2 := instructions[opcode][0], instructions[opcode][1], instructions[opcode][2]
	cycle, _ := strconv.ParseFloat(instructions[opcode][3], 64)

	// if instruction != "HALT" {
	// 	cpu.pushHistory(cpu.Reg.PC, opcode, instruction, operand1, operand2)
	// }

	switch instruction {
	case "HALT":
		cpu.HALT(operand1, operand2)
	case "LD":
		cpu.LD(operand1, operand2)
	case "LDH":
		cpu.LDH(operand1, operand2)
	case "JR":
		cpu.JR(operand1, operand2)
	case "NOP":
		cpu.NOP(operand1, operand2)
	case "AND":
		cpu.AND(operand1, operand2)
	case "INC":
		cpu.INC(operand1, operand2)
	case "DEC":
		cpu.DEC(operand1, operand2)
	case "PUSH":
		cpu.PUSH(operand1, operand2)
	case "POP":
		cpu.POP(operand1, operand2)
	case "XOR":
		cpu.XOR(operand1, operand2)
	case "JP":
		cpu.JP(operand1, operand2)
	case "CALL":
		cpu.CALL(operand1, operand2)
	case "RET":
		cpu.RET(operand1, operand2)
	case "RETI":
		cpu.RETI(operand1, operand2)
	case "CP":
		cpu.CP(operand1, operand2)
	case "OR":
		cpu.OR(operand1, operand2)
	case "ADD":
		cpu.ADD(operand1, operand2)
	case "SUB":
		cpu.SUB(operand1, operand2)
	case "ADC":
		cpu.ADC(operand1, operand2)
	case "SBC":
		cpu.SBC(operand1, operand2)
	case "CPL":
		cpu.CPL(operand1, operand2)
	case "PREFIX CB":
		cpu.PREFIXCB(operand1, operand2)
	case "RRA":
		cpu.RRA(operand1, operand2)
	case "DAA":
		cpu.DAA(operand1, operand2)
	case "RST":
		cpu.RST(operand1, operand2)
	case "SCF":
		cpu.SCF(operand1, operand2)
	case "CCF":
		cpu.CCF(operand1, operand2)
	case "RLCA":
		cpu.RLCA(operand1, operand2)
	case "RLA":
		cpu.RLA(operand1, operand2)
	case "RRCA":
		cpu.RRCA(operand1, operand2)
	case "DI":
		cpu.DI(operand1, operand2)
	case "EI":
		cpu.EI(operand1, operand2)
	case "STOP":
		cpu.STOP(operand1, operand2)
	default:
		cpu.writeHistory()

		errMsg := fmt.Sprintf("eip: 0x%04x opcode: 0x%02x", cpu.Reg.PC, opcode)
		panic(errMsg)
	}

	// incrementDebugCounter(instruction, operand1, operand2)

	cpu.mutex.Unlock()

	cpu.timer(instruction, cycle)

	cpu.handleInterrupt()
}
