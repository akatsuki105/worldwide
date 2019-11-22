package emulator

func (cpu *CPU) setHBlankMode() {
	STAT := cpu.FetchMemory8(LCDSTATIO)
	STAT &= 0xfc // bit0-1を00にする
	cpu.SetMemory8(LCDSTATIO, STAT)

	if cpu.GPU.HBlankDMALength > 0 {
		cpu.doVRAMDMATransfer(0x10)
		if cpu.GPU.HBlankDMALength == 1 {
			cpu.RAM[HDMA5IO] = 0xff
		} else {
			cpu.GPU.HBlankDMALength--
			cpu.RAM[HDMA5IO] = byte(cpu.GPU.HBlankDMALength)
		}
	}

	if (STAT & 0x08) != 0 {
		cpu.setLCDSTATFlag()
	}
}

func (cpu *CPU) getVBlankMode() bool {
	STAT := cpu.FetchMemory8(LCDSTATIO)
	if STAT&0x03 == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setVBlankMode() {
	STAT := cpu.FetchMemory8(LCDSTATIO)
	STAT = (STAT | 0x01) & 0xfd // bit0-1を01にする
	cpu.SetMemory8(LCDSTATIO, STAT)
}

func (cpu *CPU) setOAMRAMMode() {
	STAT := cpu.FetchMemory8(LCDSTATIO)
	STAT = (STAT | 0x02) & 0xfe // bit0-1を10にする
	cpu.SetMemory8(LCDSTATIO, STAT)
	if (STAT & 0x20) != 0 {
		cpu.setLCDSTATFlag()
	}
}

func (cpu *CPU) setLCDMode() {
	STAT := cpu.FetchMemory8(LCDSTATIO)
	STAT |= 0x03 // bit0-1を11にする
	cpu.SetMemory8(LCDSTATIO, STAT)
}

func (cpu *CPU) incrementLY() {
	LY := uint8(cpu.FetchMemory8(LYIO))
	LY++
	if LY == 144 {
		// VBlank期間フラグを立てる
		cpu.setVBlankMode()

		if cpu.Reg.IME && cpu.getVBlankEnable() {
			cpu.triggerVBlank()
		}
	}
	if LY > 153 {
		LY = 0
	}
	cpu.RAM[LYIO] = byte(LY)
	cpu.compareLYC(LY)
}

func (cpu *CPU) compareLYC(LY uint8) {
	LYC := cpu.FetchMemory8(0xff45)
	if LYC == LY {
		// LCDC STAT IOポートの一致フラグをセットする
		STAT := cpu.FetchMemory8(LCDSTATIO) | 0x04
		cpu.SetMemory8(LCDSTATIO, STAT)

		enable := cpu.getLCDSTATEnable()
		if enable && (STAT>>6)&0x01 == 1 {
			cpu.triggerLCDC()
		}
	}
}
