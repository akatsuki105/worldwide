package gbc

func (cpu *CPU) push(b byte) {
	cpu.SetMemory8(cpu.Reg.SP-1, b)
	cpu.Reg.SP--
}

func (cpu *CPU) pop() byte {
	value := cpu.FetchMemory8(cpu.Reg.SP)
	cpu.Reg.SP++
	return value
}

func (cpu *CPU) pushPC() {
	upper, lower := byte(cpu.Reg.PC>>8), byte(cpu.Reg.PC)
	cpu.push(upper)
	cpu.push(lower)
}

func (cpu *CPU) pushPCCALL() {
	upper := byte(cpu.Reg.PC >> 8)
	cpu.push(upper)
	cpu.timer(1) // M = 4: PC push: memory access for high byte
	lower := byte(cpu.Reg.PC & 0x00ff)
	cpu.push(lower)
	cpu.timer(1) // M = 5: PC push: memory access for low byte
}

func (cpu *CPU) popPC() {
	lower := uint16(cpu.pop())
	upper := uint16(cpu.pop())
	cpu.Reg.PC = (upper << 8) | lower
}
