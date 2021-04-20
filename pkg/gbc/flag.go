package gbc

const (
	flagZ, flagN, flagH, flagC = 7, 6, 5, 4
)

func (cpu *CPU) setCSub(dst, src byte) {
	cpu.setF(flagC, dst < uint8(dst-src))
}

func (cpu *CPU) f(idx int) bool {
	return cpu.Reg.R[F]&(1<<idx) != 0
}

func (cpu *CPU) setF(idx int, flag bool) {
	if flag {
		cpu.Reg.R[F] |= (1 << idx)
		return
	}
	cpu.Reg.R[F] &= ^(1 << idx)
}
