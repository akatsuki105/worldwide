package gbc

import (
	"gbc/pkg/util"
)

func (g *GBC) incrementLY() {
	LY := g.Load8(LYIO)
	LY++
	if LY == 144 { // set vblank flag
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
