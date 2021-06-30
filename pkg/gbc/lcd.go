package gbc

import "gbc/pkg/util"

func (g *GBC) setHBlankMode() {
	g.mode = HBlankMode
	stat := g.Load8(LCDSTATIO) & 0b1111_1100
	g.Store8(LCDSTATIO, stat)

	if g.GPU.HBlankDMALength > 0 {
		g.doVRAMDMATransfer(0x10)
		if g.GPU.HBlankDMALength == 1 {
			g.GPU.HBlankDMALength--
			g.RAM[HDMA5IO] = 0xff
		} else {
			g.GPU.HBlankDMALength--
			g.RAM[HDMA5IO] = byte(g.GPU.HBlankDMALength)
		}
	}

	if util.Bit(stat, 3) {
		g.setLCDSTATFlag(true)
	}
}

func (g *GBC) setVBlankMode() {
	g.mode = VBlankMode
	stat := (g.Load8(LCDSTATIO) | 0x01) & 0xfd // bit0-1: 01
	g.Store8(LCDSTATIO, stat)
}

func (g *GBC) setOAMRAMMode() {
	g.mode = OAMRAMMode
	stat := (g.Load8(LCDSTATIO) | 0x02) & 0xfe // bit0-1: 10
	g.Store8(LCDSTATIO, stat)
	if util.Bit(stat, 5) {
		g.setLCDSTATFlag(true)
	}
}

func (g *GBC) setLCDMode() {
	g.mode = LCDMode
	stat := g.Load8(LCDSTATIO) | 0b11
	g.Store8(LCDSTATIO, stat)
}

func (g *GBC) incrementLY() {
	LY := g.Load8(LYIO)
	LY++
	if LY == 144 { // set vblank flag
		g.setVBlankMode()
		g.setVBlankFlag(true)
	}
	LY %= 154 // LY = LY >= 154 ? 0 : LY
	g.RAM[LYIO] = LY
	g.checkLYC(LY)
}

func (g *GBC) checkLYC(LY uint8) {
	LYC := g.Load8(LYCIO)
	if LYC == LY {
		stat := g.Load8(LCDSTATIO) | 0x04 // set lyc flag
		g.Store8(LCDSTATIO, stat)

		if util.Bit(stat, 6) { // trigger LYC=LY interrupt
			g.setLCDSTATFlag(true)
		}
		return
	}

	stat := g.Load8(LCDSTATIO) & 0b11111011 // clear lyc flag
	g.Store8(LCDSTATIO, stat)
}
