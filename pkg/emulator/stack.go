package emulator

func (cpu *CPU) push(b byte) {
	cpu.SetMemory8(cpu.Reg.SP-1, b)
	cpu.Reg.SP--
}

func (cpu *CPU) pop() byte {
	value := cpu.FetchMemory8(cpu.Reg.SP)
	cpu.Reg.SP++
	return value
}

// ------------ AF --------------------

func (cpu *CPU) pushAF() {
	upper := byte(cpu.Reg.AF >> 8)
	lower := byte(cpu.Reg.AF & 0x00f0)
	cpu.push(upper)
	cpu.push(lower)
}

func (cpu *CPU) popAF() {
	lower := uint16(cpu.pop() & 0xf0)
	upper := uint16(cpu.pop())
	AF := (upper << 8) | lower
	cpu.Reg.AF = AF
}

// ------------ BC --------------------

func (cpu *CPU) pushBC() {
	upper := byte(cpu.Reg.BC >> 8)
	lower := byte(cpu.Reg.BC & 0x00ff)
	cpu.push(upper)
	cpu.push(lower)
}

func (cpu *CPU) popBC() {
	lower := uint16(cpu.pop())
	upper := uint16(cpu.pop())
	BC := (upper << 8) | lower
	cpu.Reg.BC = BC
}

// ------------ DE --------------------

func (cpu *CPU) pushDE() {
	upper := byte(cpu.Reg.DE >> 8)
	lower := byte(cpu.Reg.DE & 0x00ff)
	cpu.push(upper)
	cpu.push(lower)
}

func (cpu *CPU) popDE() {
	lower := uint16(cpu.pop())
	upper := uint16(cpu.pop())
	DE := (upper << 8) | lower
	cpu.Reg.DE = DE
}

// ------------ HL --------------------

func (cpu *CPU) pushHL() {
	upper := byte(cpu.Reg.HL >> 8)
	lower := byte(cpu.Reg.HL & 0x00ff)
	cpu.push(upper)
	cpu.push(lower)
}

func (cpu *CPU) popHL() {
	lower := uint16(cpu.pop())
	upper := uint16(cpu.pop())
	HL := (upper << 8) | lower
	cpu.Reg.HL = HL
}

// ------------ PC --------------------

func (cpu *CPU) pushPC() {
	upper := byte(cpu.Reg.PC >> 8)
	lower := byte(cpu.Reg.PC & 0x00ff)
	cpu.push(upper)
	cpu.push(lower)
}

func (cpu *CPU) popPC() {
	lower := uint16(cpu.pop())
	upper := uint16(cpu.pop())
	PC := (upper << 8) | lower
	cpu.Reg.PC = PC
}
