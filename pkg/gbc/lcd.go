package gbc

import "gbc/pkg/util"

func (cpu *CPU) setHBlankMode() {
	cpu.mode = HBlankMode
	STAT := cpu.FetchMemory8(LCDSTATIO)
	STAT &= 0b1111_1100
	cpu.SetMemory8(LCDSTATIO, STAT)

	if cpu.GPU.HBlankDMALength > 0 {
		cpu.doVRAMDMATransfer(0x10)
		if cpu.GPU.HBlankDMALength == 1 {
			cpu.GPU.HBlankDMALength--
			cpu.RAM[HDMA5IO] = 0xff
		} else {
			cpu.GPU.HBlankDMALength--
			cpu.RAM[HDMA5IO] = byte(cpu.GPU.HBlankDMALength)
		}
	}

	if util.Bit(STAT, 3) {
		cpu.setLCDSTATFlag(true)
	}
}

func (cpu *CPU) setVBlankMode() {
	cpu.mode = VBlankMode
	STAT := cpu.FetchMemory8(LCDSTATIO)
	STAT = (STAT | 0x01) & 0xfd // bit0-1: 01
	cpu.SetMemory8(LCDSTATIO, STAT)
}

func (cpu *CPU) setOAMRAMMode() {
	cpu.mode = OAMRAMMode
	STAT := cpu.FetchMemory8(LCDSTATIO)
	STAT = (STAT | 0x02) & 0xfe // bit0-1: 10
	cpu.SetMemory8(LCDSTATIO, STAT)
	if util.Bit(STAT, 5) {
		cpu.setLCDSTATFlag(true)
	}
}

func (cpu *CPU) setLCDMode() {
	cpu.mode = LCDMode
	STAT := cpu.FetchMemory8(LCDSTATIO)
	STAT |= 0b11
	cpu.SetMemory8(LCDSTATIO, STAT)
}

func (cpu *CPU) incrementLY() {
	LY := cpu.FetchMemory8(LYIO)
	LY++
	if LY == 144 { // set vblank flag
		cpu.setVBlankMode()
		cpu.setVBlankFlag(true)
	}
	LY %= 154 // LY = LY >= 154 ? 0 : LY
	cpu.RAM[LYIO] = LY
	cpu.compareLYC(LY)
}

func (cpu *CPU) compareLYC(LY uint8) {
	LYC := cpu.FetchMemory8(LYCIO)
	if LYC == LY {
		// LCDC STAT IOポートの一致フラグをセットする
		STAT := cpu.FetchMemory8(LCDSTATIO) | 0x04
		cpu.SetMemory8(LCDSTATIO, STAT)

		if util.Bit(STAT, 6) {
			cpu.setLCDSTATFlag(true)
		}
		return
	}

	// LCDC STAT IOポートの一致フラグをクリアする
	STAT := cpu.FetchMemory8(LCDSTATIO) & 0b11111011
	cpu.SetMemory8(LCDSTATIO, STAT)
}
