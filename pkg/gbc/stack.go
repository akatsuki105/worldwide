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

// ------------ AF --------------------

func (cpu *CPU) pushAF() {
	cpu.push(cpu.Reg.A)
	cpu.timer(1)
	cpu.push(cpu.Reg.F & 0x00f0)
}

func (cpu *CPU) popAF() {
	cpu.Reg.F = cpu.pop() & 0xf0
	cpu.timer(1)
	cpu.Reg.A = cpu.pop()
}

// ------------ BC --------------------

func (cpu *CPU) pushBC() {
	cpu.push(cpu.Reg.B)
	cpu.timer(1)
	cpu.push(cpu.Reg.C)
}

func (cpu *CPU) popBC() {
	cpu.Reg.C = cpu.pop()
	cpu.timer(1)
	cpu.Reg.B = cpu.pop()
}

// ------------ DE --------------------

func (cpu *CPU) pushDE() {
	cpu.push(cpu.Reg.D) // まだOAMDMA中なのでここでのアクセスは弾かれる https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/push_timing.s
	cpu.timer(1)        // OAMDMAが終わる
	cpu.push(cpu.Reg.E)
}

func (cpu *CPU) popDE() {
	cpu.Reg.E = cpu.pop()
	cpu.timer(1)
	cpu.Reg.D = cpu.pop()
}

// ------------ HL --------------------

func (cpu *CPU) pushHL() {
	cpu.push(cpu.Reg.H)
	cpu.timer(1)
	cpu.push(cpu.Reg.L)
}

func (cpu *CPU) popHL() {
	cpu.Reg.L = cpu.pop()
	cpu.timer(1)
	cpu.Reg.H = cpu.pop()
}

// ------------ PC --------------------

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
