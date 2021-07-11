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

const (
	AF = iota
	BC
	DE
	HL
	HLI
	HLD
	SP
	PC
)
const (
	flagZ, flagN, flagH, flagC = 7, 6, 5, 4
)

// Register Z80
type Register struct {
	R   [8]byte
	SP  uint16
	PC  uint16
	IME bool
}

func (r *Register) R16(i int) uint16 {
	switch i {
	case AF:
		return r.AF()
	case BC:
		return r.BC()
	case DE:
		return r.DE()
	case HL:
		return r.HL()
	case HLD:
		hl := r.HL()
		r.setHL(hl - 1)
		return hl
	case HLI:
		hl := r.HL()
		r.setHL(hl + 1)
		return hl
	case SP:
		return r.SP
	case PC:
		return r.PC
	}
	panic("invalid register16")
}

func (r *Register) setR16(i int, val uint16) {
	switch i {
	case AF:
		r.setAF(val)
	case BC:
		r.setBC(val)
	case DE:
		r.setDE(val)
	case HL:
		r.setHL(val)
	case SP:
		r.SP = val
	case PC:
		r.PC = val
	}
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

func (g *GBC) getRegister(s string) uint16 {
	switch s {
	case "A":
		return uint16(g.Reg.R[A])
	case "F":
		return uint16(g.Reg.R[F])
	case "B":
		return uint16(g.Reg.R[B])
	case "C":
		return uint16(g.Reg.R[C])
	case "D":
		return uint16(g.Reg.R[D])
	case "E":
		return uint16(g.Reg.R[E])
	case "H":
		return uint16(g.Reg.R[H])
	case "L":
		return uint16(g.Reg.R[L])
	case "AF":
		return g.Reg.AF()
	case "BC":
		return g.Reg.BC()
	case "DE":
		return g.Reg.DE()
	case "HL":
		return g.Reg.HL()
	case "SP":
		return g.Reg.SP
	}

	return 0
}

// flag

func (g *GBC) setCSub(dst, src byte) {
	g.setF(flagC, dst < uint8(dst-src))
}

func (g *GBC) f(idx int) bool {
	return g.Reg.R[F]&(1<<idx) != 0
}

func (g *GBC) setF(idx int, flag bool) {
	if flag {
		g.Reg.R[F] |= (1 << idx)
		return
	}
	g.Reg.R[F] &= ^(1 << idx)
}
