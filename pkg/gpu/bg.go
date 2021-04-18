package gpu

import (
	"gbc/pkg/util"
	"image/color"
)

type EntryY struct {
	Block  int
	Offset int
}

// SetBGLine 1タイルライン描画する
func (g *GPU) SetBGLine(entryX int, entryY EntryY, tileX, tileY uint, isWin, isCGB bool, lineIdx int) bool {
	index := tileX + tileY*32 // マップの何タイル目か

	// get map address
	mapAddr := 0x9800 + uint16(index)
	if util.Bit(g.LCDC, 3) {
		mapAddr = 0x9c00 + uint16(index)
	}
	if isWin {
		mapAddr = 0x9800 + uint16(index)
		if util.Bit(g.LCDC, 6) {
			mapAddr = 0x9c00 + uint16(index)
		}
	}

	tileIdx := g.VRAM.Bank[0][mapAddr-0x8000] // BG Mapから該当の画面の場所のタイル番号を取得
	baseAddr := g.fetchTileBaseAddr()
	if baseAddr == 0x8800 {
		tileIdx = uint8(int(int8(tileIdx)) + 128)
	}

	// bg attr
	attr := byte(0)
	if isCGB {
		attr = g.VRAM.Bank[1][mapAddr-0x8000]
	}

	tileDataOffset := uint16(tileIdx)*8 + uint16(lineIdx) // 何枚目のタイルか*8 + タイルの何行目か = 描画対象のタイルデータのオフセット
	tileDataAddr := baseAddr + 2*tileDataOffset           // タイルデータのアドレス
	return g.setBGLine(entryX, entryY.Block, entryY.Offset, tileDataAddr, attr, isCGB)
}

// SetBGPriorPixels 背景優先の背景を描画するための関数
func (g *GPU) SetBGPriorPixels() {
	for _, pixel := range g.BGPriorPixels {
		x, y := int(pixel[0]), int(pixel[1])
		R, G, B := pixel[2], pixel[3], pixel[4]
		c := color.RGBA{R, G, B, 0xff}
		if x < 160 && y < 144 {
			g.set(x, y, c)
		}
	}
	g.BGPriorPixels = [][5]byte{}
}

func (g *GPU) setBGLine(entryX, entryY, lineNumber int, addr uint16, attr byte, isCGB bool) bool {

	// entryX, entryY: 何Pixel目を基準として配置するか
	VRAMBankPtr := (attr >> 3) & 0x01
	if !isCGB {
		VRAMBankPtr = 0
	}

	lowerByte, upperByte := g.VRAM.Bank[VRAMBankPtr][addr-0x8000], g.VRAM.Bank[VRAMBankPtr][addr-0x8000+1]
	for i := 0; i < 8; i++ {
		bit := (7 - uint(i)) // upper bit
		upperColor, lowerColor := (upperByte>>bit)&0x01, (lowerByte>>bit)&0x01
		colorIdx := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

		var RGB, R, G, B byte
		var isTransparent bool

		// 色番号からRGB値を算出する
		if isCGB {
			palIdx := attr & 0x07 // パレット番号 OBPn
			R, G, B, isTransparent = g.parseCGBPallete(BGP, palIdx, colorIdx)
		} else {
			RGB, isTransparent = g.parsePallete(BGP, colorIdx)
			R, G, B = colors[RGB][0], colors[RGB][1], colors[RGB][2]
		}
		c := color.RGBA{R, G, B, 0xff}

		var deltaX, deltaY int
		if !isTransparent {
			if util.Bit(attr, 6) && util.Bit(attr, 5) { // xy flip
				deltaX, deltaY = 7-i, 7-lineNumber
			} else if util.Bit(attr, 6) { // y flip
				deltaX, deltaY = i, 7-lineNumber
			} else if util.Bit(attr, 5) { // x flip
				deltaX, deltaY = 7-i, lineNumber
			} else {
				deltaX, deltaY = i, lineNumber
			}

			x, y := entryX+deltaX, entryY+deltaY
			if (x >= 0 && x < 160) && (y >= 0 && y < 144) {
				g.displayColor[y][x] = colorIdx
				if util.Bit(attr, 7) {
					g.BGPriorPixels = append(g.BGPriorPixels, [5]byte{byte(x), byte(y), R, G, B})
				}
				g.set(x, y, c)
			} else if x >= 160 {
				break
			} else if y >= 144 {
				return false
			}
		}
	}

	return true
}
