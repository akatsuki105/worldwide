package gpu

import "image/color"

type EntryY struct {
	Block  int
	Offset int
}

// SetBGLine 1タイルライン描画する
func (g *GPU) SetBGLine(entryX int, entryY EntryY, tileX, tileY uint, useWindow, isCGB bool, lineNumber int) bool {
	index := tileX + tileY*32 // マップの何タイル目か

	// タイル番号からタイルデータのあるアドレス取得
	var BGMapAddr uint16
	LCDC := g.LCDC
	if useWindow {
		BGMapAddr = 0x9800 + uint16(index)
		if LCDC&0x40 != 0 {
			BGMapAddr = 0x9c00 + uint16(index)
		}
	} else {
		BGMapAddr = 0x9800 + uint16(index)
		if LCDC&0x08 != 0 {
			BGMapAddr = 0x9c00 + uint16(index)
		}
	}
	tileNumber := g.VRAM.Bank[0][BGMapAddr-0x8000] // BG Mapから該当の画面の場所のタイル番号を取得
	baseAddr := g.fetchTileBaseAddr()
	if baseAddr == 0x8800 {
		tileNumber = uint8(int(int8(tileNumber)) + 128)
	}

	// 背景属性取得
	attr := byte(0)
	if isCGB {
		attr = g.VRAM.Bank[1][BGMapAddr-0x8000]
	}

	tileDataOffset := uint16(tileNumber)*8 + uint16(lineNumber) // 何枚目のタイルか*8 + タイルの何行目か = 描画対象のタイルデータのオフセット
	tileDataAddr := baseAddr + 2*tileDataOffset                 // タイルデータのアドレス
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
		bitCtr := (7 - uint(i)) // 上位何ビット目を取り出すか
		upperColor := (upperByte >> bitCtr) & 0x01
		lowerColor := (lowerByte >> bitCtr) & 0x01
		colorNumber := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

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
				deltaX = 7 - i
				deltaY = 7 - lineNumber
			} else if (attr>>6)&0x01 == 1 {
				// 上下
				deltaX = i
				deltaY = 7 - lineNumber
			} else if (attr>>5)&0x01 == 1 {
				// 左右
				deltaX = 7 - i
				deltaY = lineNumber
			} else {
				// 反転無し
				deltaX = i
				deltaY = lineNumber
			}
			x := entryX + deltaX
			y := entryY + deltaY

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
