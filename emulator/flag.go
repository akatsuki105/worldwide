package emulator

func (cpu *CPU) flagZ(b byte) {
	if b == 0 {
		cpu.setZFlag()
	} else {
		cpu.clearZFlag()
	}
}

func (cpu *CPU) flagN(isSub bool) {
	if isSub {
		cpu.setNFlag()
	} else {
		cpu.clearNFlag()
	}
}

func (cpu *CPU) flagH8(value4 byte) {
	if value4&0x10 != 0 {
		cpu.setHFlag()
	} else {
		cpu.clearHFlag()
	}
}

func (cpu *CPU) flagH16(value12 uint16) {
	if value12&0x1000 != 0 {
		cpu.setHFlag()
	} else {
		cpu.clearHFlag()
	}
}

func (cpu *CPU) flagC8(u16 uint16) {
	if u16&0x100 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
}

func (cpu *CPU) flagC8Sub(dst, src byte) {
	if dst < uint8(dst-src) {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
}

func (cpu *CPU) flagC16(u32 uint32) {
	if u32&0x10000 != 0 {
		cpu.setCFlag()
	} else {
		cpu.clearCFlag()
	}
}

func (cpu *CPU) getZFlag() bool {
	if cpu.Reg.AF&0x0080 != 0 {
		return true
	}
	return false
}

func (cpu *CPU) setZFlag() {
	cpu.Reg.AF |= 0x0080
}

func (cpu *CPU) clearZFlag() {
	cpu.Reg.AF &= 0xff7f
}

func (cpu *CPU) getNFlag() bool {
	if cpu.Reg.AF&0x0040 != 0 {
		return true
	}
	return false
}

func (cpu *CPU) setNFlag() {
	cpu.Reg.AF |= 0x0040
}

func (cpu *CPU) clearNFlag() {
	cpu.Reg.AF &= 0xffbf
}

func (cpu *CPU) getHFlag() bool {
	if cpu.Reg.AF&0x0020 != 0 {
		return true
	}
	return false
}

func (cpu *CPU) setHFlag() {
	cpu.Reg.AF |= 0x0020
}

func (cpu *CPU) clearHFlag() {
	cpu.Reg.AF &= 0xffdf
}

func (cpu *CPU) getCFlag() bool {
	if cpu.Reg.AF&0x0010 != 0 {
		return true
	}
	return false
}

func (cpu *CPU) setCFlag() {
	cpu.Reg.AF |= 0x0010
}

func (cpu *CPU) clearCFlag() {
	cpu.Reg.AF &= 0xffef
}
