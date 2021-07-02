package gpu

import (
	"gbc/pkg/util"
)

// SetSPRTile スプライトを出力する
func (g *GPU) SetSPRTile(OAMindex, entryX, entryY int, tileIndex uint, attr byte, isCGB bool) {
	spriteYSize := g.fetchSPRYSize()
	if util.Bit(attr, 4) {
		for lineNumber := 0; lineNumber < spriteYSize; lineNumber++ {
			index := uint16(tileIndex)*8 + uint16(lineNumber) // 何枚目のタイルか*8 + タイルの何行目か
			addr := uint16(0x8000 + 2*index)                  // スプライトは0x8000のみ
			continueFlag := g.setSPRLine(entryX, entryY, lineNumber, addr, OBP1, attr, isCGB, OAMindex)
			if !continueFlag {
				break
			}
		}
	} else {
		for lineNumber := 0; lineNumber < spriteYSize; lineNumber++ {
			index := uint16(tileIndex)*8 + uint16(lineNumber) // 何枚目のタイルか*8 + タイルの何行目か
			addr := uint16(0x8000 + 2*index)                  // スプライトは0x8000のみ
			continueFlag := g.setSPRLine(entryX, entryY, lineNumber, addr, OBP0, attr, isCGB, OAMindex)
			if !continueFlag {
				break
			}
		}
	}
}

func (g *GPU) fetchSPRYSize() int {
	LCDC := g.LCDC
	if LCDC&0x04 != 0 {
		return 16
	}
	return 8
}

func (g *GPU) setSPRLine(entryX, entryY, lineNumber int, addr uint16, tileType int, attr byte, isCGB bool, OAMindex int) bool {
	return true
}
