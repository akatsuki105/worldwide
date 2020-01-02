package emulator

import (
	"fmt"
	"os"
	"os/signal"
)

var (
	maxHistory = 256

	insCounter map[string]uint = map[string]uint{}
	opCounter  map[string]uint = map[string]uint{}
)

func incrementDebugCounter(opcode, operand1, operand2 string) {
	ins := fmt.Sprintf("%s %s,%s", opcode, operand1, operand2)

	opctr := opCounter[opcode]
	opCounter[opcode] = opctr + 1

	insctr := insCounter[ins]
	insCounter[ins] = insctr + 1
}

func writeOpCounter() {
	sum := uint(0)
	for _, counter := range opCounter {
		sum += counter
	}

	for opcode, counter := range opCounter {
		percent := float64(counter*100) / float64(sum)
		if percent > 2.0 {
			fmt.Println(opcode, " => ", percent, "%")
		}
	}
}

func writeInsCounter() {
	sum := uint(0)
	for _, counter := range insCounter {
		sum += counter
	}

	for instruction, counter := range insCounter {
		percent := float64(counter*100) / float64(sum)
		if percent > 1.0 {
			fmt.Println(instruction, " => ", percent, "%")
		}
	}
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
func (cpu *CPU) Debug(mode int) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	switch mode {
	case 0:
	case 1:
		println("\n ============== Debug mode ==============\n")
		cpu.writeHistory()
	case 2:
		cpu.mutex.Lock()
		println("\n ============== Opcode counter ==============\n")
		writeOpCounter()
		println("\n ============== Instruction counter ==============\n")
		writeInsCounter()
		cpu.mutex.Unlock()
	}

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
