package gpu

import (
	"gbc/pkg/util"
	"image/color"
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
	spriteYSize := g.fetchSPRYSize()

	// entryX, entryY: 何Pixel目を基準として配置するか
	VRAMBankPtr := (attr >> 3) & 0x01
	if !isCGB {
		VRAMBankPtr = 0
	}

	lowerByte, upperByte := g.VRAM.Bank[VRAMBankPtr][addr-0x8000], g.VRAM.Bank[VRAMBankPtr][addr-0x8000+1]
	for i := 0; i < 8; i++ {
		bit := (7 - uint(i)) // 上位何ビット目を取り出すか
		upperColor, lowerColor := (upperByte>>bit)&0x01, (lowerByte>>bit)&0x01
		colorIdx := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

		// 色番号からRGB値を算出する
		RGB, isTransparent := g.parsePallete(tileType, colorIdx)
		R, G, B := colors[RGB][0], colors[RGB][1], colors[RGB][2]
		if isCGB {
			palIdx := attr & 0x07 // パレット番号 OBPn
			R, G, B, isTransparent = g.parseCGBPallete(tileType, palIdx, colorIdx)
		}
		c := color.RGBA{R, G, B, 0xff}

		var deltaX, deltaY int
		if !isTransparent {
			switch {
			case util.Bit(attr, 6) && util.Bit(attr, 5): // xy flip
				deltaX, deltaY = int(7-i), (spriteYSize-1)-lineNumber
			case util.Bit(attr, 6): // y flip
				deltaX, deltaY = int(i), (spriteYSize-1)-lineNumber
			case util.Bit(attr, 5): // x flip
				deltaX, deltaY = int(7-i), lineNumber
			default:
				deltaX, deltaY = int(i), lineNumber
			}

			x, y := entryX+deltaX, entryY+deltaY
			if g.Debug.On { // debug OAM
				col, row := OAMindex%8, OAMindex/8
				g.OAM.Set(col*16+deltaX+2, row*20+deltaY, c)
			}

			if (x >= 0 && x < 160) && (y >= 0 && y < 144) {
				if (attr>>7)&0x01 == 0 && g.displayColor[y][x] != 0 {
					g.set(x, y, c)
				} else if g.displayColor[y][x] == 0 {
					g.set(x, y, c)
				}
			} else if x >= 160 {
				break
			} else if y >= 144 {
				return false
			}
		}
	}

	return true
}
