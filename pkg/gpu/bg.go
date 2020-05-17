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
	return g.setBGLine(entryX, entryY, uint(lineIndex), addr, attr, isCGB)
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

func (g *GPU) setBGLine(entryX, entryY int, lineIndex uint, addr uint16, attr byte, isCGB bool) bool {

	// entryX, entryY: 何Pixel目を基準として配置するか
	VRAMBankPtr := (attr >> 3) & 0x01
	if !isCGB {
		VRAMBankPtr = 0
	}

	lowerByte, upperByte := g.VRAMBank[VRAMBankPtr][addr-0x8000], g.VRAMBank[VRAMBankPtr][addr-0x8000+1]

	for j := 0; j < 8; j++ {
		bitCtr := (7 - uint(j)) // 上位何ビット目を取り出すか
		upperColor := (upperByte >> bitCtr) & 0x01
		lowerColor := (lowerByte >> bitCtr) & 0x01
		colorNumber := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

		var x, y int
		var RGB, R, G, B byte
		var isTransparent bool

		// 色番号からRGB値を算出する
		if isCGB {
			palleteNumber := attr & 0x07 // パレット番号 OBPn
			R, G, B, isTransparent = g.parseCGBPallete(BGP, palleteNumber, colorNumber)
		} else {
			RGB, isTransparent = g.parsePallete(BGP, colorNumber)
			R, G, B = colors[RGB][0], colors[RGB][1], colors[RGB][2]
		}
		c := color.RGBA{R, G, B, 0xff}

		var deltaX, deltaY int
		if !isTransparent {
			// 反転を考慮してpixelをセット
			if (attr>>6)&0x01 == 1 && (attr>>5)&0x01 == 1 {
				// 上下左右
				deltaX = int((7 - j))
				deltaY = int((7 - int(lineIndex)))
			} else if (attr>>6)&0x01 == 1 {
				// 上下
				deltaX = int(j)
				deltaY = int((7 - int(lineIndex)))
			} else if (attr>>5)&0x01 == 1 {
				// 左右
				deltaX = int((7 - j))
				deltaY = int(lineIndex)
			} else {
				// 反転無し
				deltaX = int(j)
				deltaY = int(lineIndex)
			}
			x = entryX + deltaX
			y = entryY + deltaY

			if (x >= 0 && x < 160) && (y >= 0 && y < 144) {
				g.displayColor[y][x] = colorNumber
				if (attr>>7)&0x01 == 1 {
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
