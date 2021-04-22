package gbc

import "gbc/pkg/util"

func (cpu *CPU) setHBlankMode() {
	cpu.mode = HBlankMode
	stat := cpu.FetchMemory8(LCDSTATIO) & 0b1111_1100
	cpu.SetMemory8(LCDSTATIO, stat)

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

	if util.Bit(stat, 3) {
		cpu.setLCDSTATFlag(true)
	}
}

func (cpu *CPU) setVBlankMode() {
	cpu.mode = VBlankMode
	stat := (cpu.FetchMemory8(LCDSTATIO) | 0x01) & 0xfd // bit0-1: 01
	cpu.SetMemory8(LCDSTATIO, stat)
}

func (cpu *CPU) setOAMRAMMode() {
	cpu.mode = OAMRAMMode
	stat := (cpu.FetchMemory8(LCDSTATIO) | 0x02) & 0xfe // bit0-1: 10
	cpu.SetMemory8(LCDSTATIO, stat)
	if util.Bit(stat, 5) {
		cpu.setLCDSTATFlag(true)
	}
}

func (cpu *CPU) setLCDMode() {
	cpu.mode = LCDMode
	stat := cpu.FetchMemory8(LCDSTATIO) | 0b11
	cpu.SetMemory8(LCDSTATIO, stat)
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
	cpu.checkLYC(LY)
}

func (cpu *CPU) checkLYC(LY uint8) {
	LYC := cpu.FetchMemory8(LYCIO)
	if LYC == LY {
		stat := cpu.FetchMemory8(LCDSTATIO) | 0x04 // set lyc flag
		cpu.SetMemory8(LCDSTATIO, stat)

		if util.Bit(stat, 6) { // trigger LYC=LY interrupt
			cpu.setLCDSTATFlag(true)
		}
		return
	}

	stat := cpu.FetchMemory8(LCDSTATIO) & 0b11111011 // clear lyc flag
	cpu.SetMemory8(LCDSTATIO, stat)
}
