package gbc

const (
	A = iota
	B
	C
	D
	E
	H
	L
	F
)

// Register Z80
type Register struct {
	R   [8]byte
	SP  uint16
	PC  uint16
	IME bool
}

func (r *Register) AF() uint16 {
	return (uint16(r.R[A]) << 8) | uint16(r.R[F])
}
func (r *Register) setAF(value uint16) {
	r.R[A], r.R[F] = byte(value>>8), byte(value)
}

func (r *Register) BC() uint16 {
	return (uint16(r.R[B]) << 8) | uint16(r.R[C])
}
func (r *Register) setBC(value uint16) {
	r.R[B], r.R[C] = byte(value>>8), byte(value)
}

func (r *Register) DE() uint16 {
	return (uint16(r.R[D]) << 8) | uint16(r.R[E])
}
func (r *Register) setDE(value uint16) {
	r.R[D], r.R[E] = byte(value>>8), byte(value)
}

func (r *Register) HL() uint16 {
	return (uint16(r.R[H]) << 8) | uint16(r.R[L])
}
func (r *Register) setHL(value uint16) {
	r.R[H], r.R[L] = byte(value>>8), byte(value)
}

func (cpu *CPU) getRegister(s string) uint16 {
	switch s {
	case "A":
		return uint16(cpu.Reg.R[A])
	case "F":
		return uint16(cpu.Reg.R[F])
	case "B":
		return uint16(cpu.Reg.R[B])
	case "C":
		return uint16(cpu.Reg.R[C])
	case "D":
		return uint16(cpu.Reg.R[D])
	case "E":
		return uint16(cpu.Reg.R[E])
	case "H":
		return uint16(cpu.Reg.R[H])
	case "L":
		return uint16(cpu.Reg.R[L])
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
