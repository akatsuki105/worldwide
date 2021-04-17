package gbc

const (
	flagZ = 7
	flagN = 6
	flagH = 5
	flagC = 4
)

func (cpu *CPU) setFlagH8(value4 byte) {
	cpu.setF(flagH, value4&0x10 != 0)
}

func (cpu *CPU) setFlagH16(value12 uint16) {
	cpu.setF(flagH, value12&0x1000 != 0)
}

func (cpu *CPU) setFlagC8(u16 uint16) {
	cpu.setF(flagC, u16&0x100 != 0)
}

func (cpu *CPU) setFlagC8Sub(dst, src byte) {
	cpu.setF(flagC, dst < uint8(dst-src))
}

func (cpu *CPU) setFlagC16(u32 uint32) {
	cpu.setF(flagC, u32&0x10000 != 0)
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
