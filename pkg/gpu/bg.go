package gpu

import "image/color"

// SetBGLine 1タイルライン描画する
func (g *GPU) SetBGLine(entryX, entryY int, tileX, tileY uint, useWindow, isCGB bool, lineIndex int) bool {
	index := tileX + tileY*32 // マップの何タイル目か

	// タイル番号からタイルデータのあるアドレス取得
	var addr uint16
	LCDC := g.LCDC
	if useWindow {
		if LCDC&0x40 != 0 {
			addr = 0x9c00 + uint16(index)
		} else {
			addr = 0x9800 + uint16(index)
		}
	} else {
		if LCDC&0x08 != 0 {
			addr = 0x9c00 + uint16(index)
		} else {
			addr = 0x9800 + uint16(index)
		}
	}
	tileIndex := uint8(g.VRAMBank[0][addr-0x8000])
	baseAddr := g.fetchTileBaseAddr()
	if baseAddr == 0x8800 {
		tileIndex = uint8(int(int8(tileIndex)) + 128)
	}

	// 背景属性取得
	var attr byte
	if isCGB {
		attr = uint8(g.VRAMBank[1][addr-0x8000])
	} else {
		attr = 0
	}

	index16 := uint16(tileIndex)*8 + uint16(lineIndex) // 何枚目のタイルか*8 + タイルの何行目か
	addr = uint16(baseAddr + 2*index16)
	return g.setTileLine(entryX, entryY, uint(lineIndex), addr, BGP, attr, 8, isCGB, -1)
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
