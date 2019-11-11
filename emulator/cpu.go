package emulator

import (
	"fmt"
	"image"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/faiface/pixel"
)

var (
	maxHistory = 256
)

// CPU Central Processing Unit
type CPU struct {
	Reg              Register
	RAM              [0x10000]byte
	Header           ROMHeader
	mutex            sync.Mutex
	history          []string
	joypad           Joypad
	interruptTrigger bool
	// timer関連
	cycle    float64
	cycleDIV float64
	// MBC関連
	ROMBankPtr uint8
	ROMBank    [128][0x4000]byte
	RAMBankPtr uint8
	RAMBank    [4][0x2000]byte
	bankMode   uint
	// 画面関連
	tileCache       *pixel.PictureData // タイルデータのキャッシュ
	newTileCache    *pixel.PictureData
	tileModified    bool
	mapCache        *image.RGBA
	VRAMCache       [0x2000]byte    // VRAMデータのキャッシュ
	VRAMModified    bool            // VRAMを変更したか(キャッシュを更新する必要があるか)
	PalleteModified PalleteModified // Palleteを変更したか
	// サウンド
	Sound APU
}

// ROMHeader ROMヘッダから得られたカードリッジ情報
type ROMHeader struct {
	Title         string
	CartridgeType uint8
	ROMSize       uint8
	RAMSize       uint8
}

// ParseROMHeader ROM情報をヘッダ構造体に読み込む
func (cpu *CPU) ParseROMHeader(rom []byte) {
	var titleBuf []byte
	for i := 0x0134; i < 0x0143; i++ {
		if rom[i] == 0 {
			break
		}
		titleBuf = append(titleBuf, rom[i])
	}
	cpu.Header.Title = string(titleBuf)

	cpu.Header.CartridgeType = uint8(rom[0x0147])
	cpu.Header.ROMSize = uint8(rom[0x0148])
	cpu.Header.RAMSize = uint8(rom[0x0149])
}

// LoadROM ROM情報をメモリに読み込む
func (cpu *CPU) LoadROM(rom []byte) {
	for i := 0x0000; i <= 0x7fff; i++ {
		cpu.RAM[i] = rom[i]
	}

	// カードリッジタイプで場合分け
	switch cpu.Header.CartridgeType {
	case 0:
		// CartridgeType : 0
		cpu.transferROM(2, &rom)
	case 1:
		// CartridgeType : 1
		switch cpu.Header.ROMSize {
		case 0:
			cpu.transferROM(2, &rom)
		case 1:
			cpu.transferROM(4, &rom)
		case 2:
			cpu.transferROM(8, &rom)
		case 3:
			cpu.transferROM(16, &rom)
		case 4:
			cpu.transferROM(32, &rom)
		case 5:
			cpu.transferROM(64, &rom)
		case 6:
			cpu.transferROM(128, &rom)
		default:
			errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Header.CartridgeType, cpu.Header.ROMSize, cpu.Header.RAMSize)
			panic(errorMsg)
		}
	case 2, 3:
		// CartridgeType : 2, 3
		switch cpu.Header.RAMSize {
		case 0, 1, 2:
			switch cpu.Header.ROMSize {
			case 0:
				cpu.transferROM(2, &rom)
			case 1:
				cpu.transferROM(4, &rom)
			case 2:
				cpu.transferROM(8, &rom)
			case 3:
				cpu.transferROM(16, &rom)
			case 4:
				cpu.transferROM(32, &rom)
			case 5:
				cpu.transferROM(64, &rom)
			case 6:
				cpu.transferROM(128, &rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Header.CartridgeType, cpu.Header.ROMSize, cpu.Header.RAMSize)
				panic(errorMsg)
			}
		case 3:
			cpu.bankMode = 1
			switch cpu.Header.ROMSize {
			case 0:
			case 1:
				cpu.transferROM(4, &rom)
			case 2:
				cpu.transferROM(8, &rom)
			case 3:
				cpu.transferROM(16, &rom)
			case 4:
				cpu.transferROM(32, &rom)
			default:
				errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Header.CartridgeType, cpu.Header.ROMSize, cpu.Header.RAMSize)
				panic(errorMsg)
			}
		default:
			errorMsg := fmt.Sprintf("RAMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Header.CartridgeType, cpu.Header.ROMSize, cpu.Header.RAMSize)
			panic(errorMsg)
		}
	default:
		errorMsg := fmt.Sprintf("CartridgeType is invalid => type:%x rom:%x ram:%x\n", cpu.Header.CartridgeType, cpu.Header.ROMSize, cpu.Header.RAMSize)
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

// InitCPU CPU・メモリの初期化
func (cpu *CPU) InitCPU() {
	cpu.Reg.AF = 0x01b0
	cpu.Reg.BC = 0x0013
	cpu.Reg.DE = 0x00d8
	cpu.Reg.HL = 0x014d
	cpu.Reg.PC = 0x0100
	cpu.Reg.SP = 0xfffe

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
	cpu.RAM[LCDCIO] = 0x91
	cpu.RAM[BGPIO] = 0xfc
	cpu.RAM[OBP0IO] = 0xff
	cpu.RAM[OBP1IO] = 0xff

	cpu.ROMBankPtr = 1

	cpu.mapCache = image.NewRGBA(image.Rect(0, 0, 8*256, 8*10)) // BGタイル1=256 BGタイル2=256 SPRタイル(OBP0)=256*4 SPRタイル(OBP1)=256*4 => 10行分のタイル (8*10)
}

// InitAPU init apu
func (cpu *CPU) InitAPU() {
	cpu.Sound.Init()
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
	case "LD":
		cpu.LD(operand1, operand2)
	case "LDH":
		cpu.LDH(operand1, operand2)
	case "NOP":
		cpu.NOP(operand1, operand2)
	case "INC":
		cpu.INC(operand1, operand2)
	case "DEC":
		cpu.DEC(operand1, operand2)
	case "JR":
		cpu.JR(operand1, operand2)
	case "HALT":
		cpu.HALT(operand1, operand2)
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
	case "DI":
		cpu.DI(operand1, operand2)
	case "EI":
		cpu.EI(operand1, operand2)
	case "CP":
		cpu.CP(operand1, operand2)
	case "AND":
		cpu.AND(operand1, operand2)
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
	case "PUSH":
		cpu.PUSH(operand1, operand2)
	case "POP":
		cpu.POP(operand1, operand2)
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
	case "STOP":
		cpu.STOP(operand1, operand2)
	default:
		cpu.writeHistory()

		errMsg := fmt.Sprintf("eip: 0x%04x opcode: 0x%02x", cpu.Reg.PC, opcode)
		panic(errMsg)
	}

	cpu.mutex.Unlock()

	cpu.timer(cycle)

	cpu.handleInterrupt()
}

// pushHistory CPUのログを追加する
func (cpu *CPU) pushHistory(eip uint16, opcode byte, instruction, operand1, operand2 string) {
	log := fmt.Sprintf("eip:0x%04x   opcode:%02x   %s %s,%s", eip, opcode, instruction, operand1, operand2)

	cpu.history = append(cpu.history, log)

	if len(cpu.history) > maxHistory {
		cpu.history = cpu.history[1:]
	}
}

// writeHistory CPUのログを書き出す
func (cpu *CPU) writeHistory() {
	for i, log := range cpu.history {
		fmt.Printf("%d: %s\n", i, log)
	}
}

// Debug Ctrl + C handler
func (cpu *CPU) Debug() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
	cpu.dumpVRAM("VRAM")
	println("\n ============== Debug mode ==============\n")
	cpu.writeHistory()
	os.Exit(1)
}

func hexPrintf(label string, hex int) {
	fmt.Printf("%s: 0x%x\n", label, hex)
}

func (cpu *CPU) exit(message string, breakPoint uint16) {
	if breakPoint == 0 {
		cpu.writeHistory()
		panic(message)
	} else if cpu.Reg.PC == breakPoint {
		cpu.writeHistory()
		panic(message)
	}
}

func (cpu *CPU) debugPC() {
	fmt.Printf("PC: 0x%04x\n", cpu.Reg.PC)
	fmt.Printf("next: %02x %02x %02x %02x\n", cpu.RAM[cpu.Reg.PC+1], cpu.RAM[cpu.Reg.PC+2], cpu.RAM[cpu.Reg.PC+3], cpu.RAM[cpu.Reg.PC+4])
}
