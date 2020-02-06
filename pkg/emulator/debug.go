package emulator

import (
	"fmt"
)

var (
	maxHistory = 128

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
func (cpu *CPU) pushHistory(eip uint16, opcode byte) {
	instruction, operand1, operand2 := opcodeToString[opcode][0], opcodeToString[opcode][1], opcodeToString[opcode][2]
	log := fmt.Sprintf("eip:0x%04x   opcode:%02x	%s %s %s", eip, opcode, instruction, operand1, operand2)

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

func (cpu *CPU) dumpRegister() {
	A, F := byte(cpu.Reg.AF>>8), byte(cpu.Reg.AF)
	B, C := byte(cpu.Reg.BC>>8), byte(cpu.Reg.BC)
	D, E := byte(cpu.Reg.DE>>8), byte(cpu.Reg.DE)
	H, L := byte(cpu.Reg.HL>>8), byte(cpu.Reg.HL)

	fmt.Println("-- register --")
	fmt.Printf("A: %02x    F: %02x\n", A, F)
	fmt.Printf("B: %02x    C: %02x\n", B, C)
	fmt.Printf("D: %02x    E: %02x\n", D, E)
	fmt.Printf("H: %02x    L: %02x\n", H, L)
}
