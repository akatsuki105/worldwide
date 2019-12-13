package emulator

import (
	"fmt"
	"os"
	"os/signal"
)

var (
	maxHistory = 256
)

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
	println("\n ============== Debug mode ==============\n")
	cpu.writeHistory()
	os.Exit(1)
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

func (cpu *CPU) debugPC(delta int) {
	fmt.Printf("PC: 0x%04x\n", cpu.Reg.PC)
	for i := 1; i < delta; i++ {
		fmt.Printf("%02x ", cpu.RAM[cpu.Reg.PC+uint16(i)])
	}
	fmt.Println()
}
