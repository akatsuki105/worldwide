package gbc

const (
	flagZ = 7
	flagN = 6
	flagH = 5
	flagC = 4
)

func (cpu *CPU) setCSub(dst, src byte) {
	cpu.setF(flagC, dst < uint8(dst-src))
}

func (cpu *CPU) f(idx int) bool {
	return cpu.Reg.F&(1<<idx) != 0
}

func (cpu *CPU) setF(idx int, flag bool) {
	if flag {
		cpu.Reg.F |= (1 << idx)
		return
	}
	cpu.Reg.F &= ^(1 << idx)
}
