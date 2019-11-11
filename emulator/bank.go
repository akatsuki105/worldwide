package emulator

import (
	"fmt"
	"sync"
)

// 指定したROMバンクへ現在のメモリを書き込む
func (cpu *CPU) writeROMBank(ROMBankPtr uint8) {
	currentRAM := cpu.RAM[0x4000:0x8000]
	var wait sync.WaitGroup
	wait.Add(0x4000)
	for i := 0x0000; i <= 0x3fff; i++ {
		go func(i int) {
			cpu.ROMBank[ROMBankPtr][i] = currentRAM[i]
			wait.Done()
		}(i)
	}
	wait.Wait()
}

// ROMバンクの切り替え
func (cpu *CPU) switchROMBank(newROMBankPtr uint8) {
	switchFlag := false

	switch cpu.Header.ROMSize {
	case 0:
	case 1:
		switchFlag = (newROMBankPtr < 4)
	case 2:
		switchFlag = (newROMBankPtr < 8)
	case 3:
		switchFlag = (newROMBankPtr < 16)
	case 4:
		switchFlag = (newROMBankPtr < 32)
	case 5:
		switchFlag = (newROMBankPtr < 64)
	case 6:
		switchFlag = (newROMBankPtr < 128)
	default:
		errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Header.CartridgeType, cpu.Header.ROMSize, cpu.Header.RAMSize)
		panic(errorMsg)
	}

	if switchFlag {
		cpu.ROMBankPtr = newROMBankPtr
	}
}

// RAMバンクの切り替え
func (cpu *CPU) switchRAMBank(newRAMBankPtr uint8) {
	cpu.RAMBankPtr = newRAMBankPtr
}
