package gbc

// Register Z80
type Register struct {
	A, F byte
	B, C byte
	D, E byte
	H, L byte
	SP   uint16
	PC   uint16
	IME  bool
}

func (r *Register) AF() uint16 {
	return (uint16(r.A) << 8) | uint16(r.F)
}
func (r *Register) setAF(value uint16) {
	r.A, r.F = byte(value>>8), byte(value)
}

func (r *Register) BC() uint16 {
	return (uint16(r.B) << 8) | uint16(r.C)
}
func (r *Register) setBC(value uint16) {
	r.B, r.C = byte(value>>8), byte(value)
}

func (r *Register) DE() uint16 {
	return (uint16(r.D) << 8) | uint16(r.E)
}
func (r *Register) setDE(value uint16) {
	r.D, r.E = byte(value>>8), byte(value)
}

func (r *Register) HL() uint16 {
	return (uint16(r.H) << 8) | uint16(r.L)
}
func (r *Register) setHL(value uint16) {
	r.H, r.L = byte(value>>8), byte(value)
}

func (cpu *CPU) getRegister(s string) uint16 {
	switch s {
	case "A":
		return uint16(cpu.Reg.A)
	case "F":
		return uint16(cpu.Reg.F)
	case "B":
		return uint16(cpu.Reg.B)
	case "C":
		return uint16(cpu.Reg.C)
	case "D":
		return uint16(cpu.Reg.D)
	case "E":
		return uint16(cpu.Reg.E)
	case "H":
		return uint16(cpu.Reg.H)
	case "L":
		return uint16(cpu.Reg.L)
	case "AF":
		return cpu.Reg.AF()
	case "BC":
		return cpu.Reg.BC()
	case "DE":
		return cpu.Reg.DE()
	case "HL":
		return cpu.Reg.HL()
	case "SP":
		return cpu.Reg.SP
	}

	return 0
}
