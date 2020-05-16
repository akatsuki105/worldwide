package gpu

// SetSPRTile スプライトを出力する
func (g *GPU) SetSPRTile(OAMindex, entryX, entryY int, tileIndex uint, attr byte, isCGB bool) {
	spriteYSize := g.fetchSPRYSize()
	if (attr>>4)%2 == 1 {
		for lineIndex := 0; lineIndex < spriteYSize; lineIndex++ {
			index := uint16(tileIndex)*8 + uint16(lineIndex) // 何枚目のタイルか*8 + タイルの何行目か
			addr := uint16(0x8000 + 2*index)                 // スプライトは0x8000のみ
			continueFlag := g.setTileLine(entryX, entryY, uint(lineIndex), addr, OBP1, attr, spriteYSize, isCGB, OAMindex)
			if !continueFlag {
				break
			}
		}
	} else {
		for lineIndex := 0; lineIndex < spriteYSize; lineIndex++ {
			index := uint16(tileIndex)*8 + uint16(lineIndex) // 何枚目のタイルか*8 + タイルの何行目か
			addr := uint16(0x8000 + 2*index)                 // スプライトは0x8000のみ
			continueFlag := g.setTileLine(entryX, entryY, uint(lineIndex), addr, OBP0, attr, spriteYSize, isCGB, OAMindex)
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
