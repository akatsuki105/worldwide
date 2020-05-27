package emulator

// Register CPUのレジスタ
type Register struct {
	AF  uint16
	BC  uint16
	DE  uint16
	HL  uint16
	SP  uint16
	PC  uint16
	IME bool
}

func (cpu *CPU) getAReg() byte {
	value := byte(cpu.Reg.AF >> 8)
	return value
}

func (cpu *CPU) getARegLower4() byte {
	A := cpu.getAReg()
	value := A & 0x0f
	return value
}

func (cpu *CPU) setAReg(A byte) {
	F := cpu.Reg.AF & 0x00f0
	AF := (uint16(A) << 8) | F
	cpu.Reg.AF = AF
}

func (cpu *CPU) getFReg() byte {
	value := byte(cpu.Reg.AF)
	return value
}

func (cpu *CPU) getBReg() byte {
	value := byte(cpu.Reg.BC >> 8)
	return value
}

func (cpu *CPU) getBRegLower4() byte {
	B := cpu.getBReg()
	value := B & 0x0f // 0000_1111
	return value
}

func (cpu *CPU) setBReg(B byte) {
	C := cpu.Reg.BC & 0x00ff
	BC := (uint16(B) << 8) | C
	cpu.Reg.BC = BC
}

func (cpu *CPU) getCReg() byte {
	value := byte(cpu.Reg.BC & 0x00ff)
	return value
}

func (cpu *CPU) getCRegLower4() byte {
	C := cpu.getCReg()
	value := C & 0x0f // 0000_1111
	return value
}

func (cpu *CPU) setCReg(C byte) {
	B := cpu.Reg.BC >> 8
	BC := (B << 8) | uint16(C)
	cpu.Reg.BC = BC
}

func (cpu *CPU) getBCRegLower12() uint16 {
	value := cpu.Reg.BC & 0x0fff
	return value
}

func (cpu *CPU) getDReg() byte {
	value := byte(cpu.Reg.DE >> 8)
	return value
}

func (cpu *CPU) getDRegLower4() byte {
	D := cpu.getDReg()
	value := D & 0x0f // 0000_1111
	return value
}

func (cpu *CPU) setDReg(D byte) {
	E := cpu.Reg.DE & 0x00ff
	DE := (uint16(D) << 8) | E
	cpu.Reg.DE = DE
}

func (cpu *CPU) getEReg() byte {
	value := byte(cpu.Reg.DE & 0x00ff)
	return value
}

func (cpu *CPU) getERegLower4() byte {
	E := cpu.getEReg()
	value := E & 0x0f // 0000_1111
	return value
}

func (cpu *CPU) setEReg(E byte) {
	D := cpu.Reg.DE >> 8
	DE := (D << 8) | uint16(E)
	cpu.Reg.DE = DE
}

func (cpu *CPU) getDERegLower12() uint16 {
	value := cpu.Reg.DE & 0x0fff
	return value
}

func (cpu *CPU) getHReg() byte {
	value := byte(cpu.Reg.HL >> 8)
	return value
}

func (cpu *CPU) getHRegLower4() byte {
	H := cpu.getHReg()
	value := H & 0x0f // 0000_1111
	return value
}

func (cpu *CPU) setHReg(H byte) {
	L := cpu.Reg.HL & 0x00ff
	HL := (uint16(H) << 8) | L
	cpu.Reg.HL = HL
}

func (cpu *CPU) getLReg() byte {
	value := byte(cpu.Reg.HL & 0x00ff)
	return value
}

func (cpu *CPU) getLRegLower4() byte {
	L := cpu.getLReg()
	value := L & 0x0f // 0000_1111
	return value
}

func (cpu *CPU) setLReg(L byte) {
	H := cpu.Reg.HL >> 8
	HL := (H << 8) | uint16(L)
	cpu.Reg.HL = HL
}

func (cpu *CPU) getHLRegLower12() uint16 {
	value := cpu.Reg.HL & 0x0fff
	return value
}

func (cpu *CPU) getSPRegLower4() byte {
	value := byte(cpu.Reg.SP & 0x000f)
	return value
}

func (cpu *CPU) getSPRegLower12() uint16 {
	value := cpu.Reg.SP & 0x0fff
	return value
}

func (cpu *CPU) getRegister(s string) uint16 {
	switch s {
	case "A":
		return uint16(cpu.getAReg())
	case "F":
		return uint16(cpu.getFReg())
	case "B":
		return uint16(cpu.getBReg())
	case "C":
		return uint16(cpu.getCReg())
	case "D":
		return uint16(cpu.getDReg())
	case "E":
		return uint16(cpu.getEReg())
	case "H":
		return uint16(cpu.getHReg())
	case "L":
		return uint16(cpu.getLReg())
	case "AF":
		return cpu.Reg.AF
	case "BC":
		return cpu.Reg.BC
	case "DE":
		return cpu.Reg.DE
	case "HL":
		return cpu.Reg.HL
	case "SP":
		return cpu.Reg.SP
	}

	return 0
}
