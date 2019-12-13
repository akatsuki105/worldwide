package gpu

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
)

// GPU Graphic Processor Unit
type GPU struct {
	display       *pixel.PictureData // 160*144のイメージデータ
	LCDC          byte               // LCD Control
	LCDSTAT       byte               // LCD Status
	displayColor  [144][160]byte     // 160*144の色番号(背景色を記録)
	DMGPallte     [3]byte            // DMGのパレットデータ {BGP, OGP0, OGP1}
	CGBPallte     [2]byte            // CGBのパレットデータ {BCPSIO, OCPSIO}
	BGPallete     [64]byte
	SPRPallete    [64]byte
	BGPriorPixels [][5]byte
	// VRAM bank
	VRAMBankPtr     uint8
	VRAMBank        [2][0x2000]byte // 0x8000-0x9fff ゲームボーイカラーのみ
	HBlankDMALength int
}

var (
	// colors {R, G, B}
	colors [4][3]uint8 = [4][3]uint8{
		{175, 197, 160}, {93, 147, 66}, {22, 63, 48}, {0, 40, 0},
	}
)

// Init GPU
func (gpu *GPU) Init() {
	gpu.display = pixel.MakePictureData(pixel.R(0, 0, 160, 144))
}

// InitPallete init gameboy pallete color
func (gpu *GPU) InitPallete(color0, color1, color2, color3 []int) {
	colors[0] = [3]uint8{uint8(color0[0]), uint8(color0[1]), uint8(color0[2])}
	colors[1] = [3]uint8{uint8(color1[0]), uint8(color1[1]), uint8(color1[2])}
	colors[2] = [3]uint8{uint8(color2[0]), uint8(color2[1]), uint8(color2[2])}
	colors[3] = [3]uint8{uint8(color3[0]), uint8(color3[1]), uint8(color3[2])}
}

// GetDisplay getter for display data
func (gpu *GPU) GetDisplay() *pixel.PictureData {
	return gpu.display
}

func (gpu *GPU) set(x, y int, c color.RGBA) {
	gpu.display.Pix[160*144-(y*160+(160-x))] = c
}

// --------------------------------------------- Render -----------------------------------------------------

// SetBGLine 1タイルライン描画する
func (gpu *GPU) SetBGLine(entryX, entryY int, tileX, tileY uint, useWindow, isCGB bool, lineIndex int) {
	index := tileX + tileY*32 // マップの何タイル目か

	// タイル番号からタイルデータのあるアドレス取得
	var addr uint16
	LCDC := gpu.LCDC
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
	tileIndex := uint8(gpu.VRAMBank[0][addr-0x8000])
	baseAddr := gpu.fetchTileBaseAddr()
	if baseAddr == 0x8800 {
		tileIndex = uint8(int(int8(tileIndex)) + 128)
	}

	// 背景属性取得
	var attr byte
	if isCGB {
		attr = uint8(gpu.VRAMBank[1][addr-0x8000])
	} else {
		attr = 0
	}

	index16 := uint16(tileIndex)*8 + uint16(lineIndex) // 何枚目のタイルか*8 + タイルの何行目か
	addr = uint16(baseAddr + 2*index16)
	gpu.setTileLine(entryX, entryY, uint(lineIndex), addr, "BGP", attr, 8, isCGB)
}

// SetSPRTile スプライトを出力する
func (gpu *GPU) SetSPRTile(entryX, entryY int, tileIndex uint, attr byte, isCGB bool) {
	spriteYSize := gpu.fetchSPRYSize()
	if (attr>>4)%2 == 1 {
		for lineIndex := 0; lineIndex < spriteYSize; lineIndex++ {
			index := uint16(tileIndex)*8 + uint16(lineIndex) // 何枚目のタイルか*8 + タイルの何行目か
			addr := uint16(0x8000 + 2*index)                 // スプライトは0x8000のみ
			gpu.setTileLine(entryX, entryY, uint(lineIndex), addr, "OBP1", attr, spriteYSize, isCGB)
		}
	} else {
		for lineIndex := 0; lineIndex < spriteYSize; lineIndex++ {
			index := uint16(tileIndex)*8 + uint16(lineIndex) // 何枚目のタイルか*8 + タイルの何行目か
			addr := uint16(0x8000 + 2*index)                 // スプライトは0x8000のみ
			gpu.setTileLine(entryX, entryY, uint(lineIndex), addr, "OBP0", attr, spriteYSize, isCGB)
		}
	}
}

// SetBGPriorPixels 背景優先の背景を描画するための関数
func (gpu *GPU) SetBGPriorPixels() {
	for _, pixel := range gpu.BGPriorPixels {
		x, y := int(pixel[0]), int(pixel[1])
		R, G, B := pixel[2], pixel[3], pixel[4]
		c := color.RGBA{R, G, B, 0xff}
		if x < 160 && y < 144 {
			gpu.set(x, y, c)
		}
	}
	gpu.BGPriorPixels = [][5]byte{}
}

// --------------------------------------------- CGB pallete -----------------------------------------------------

// FetchBGPalleteIndex CGBのパレットインデックスを取得する
func (gpu *GPU) FetchBGPalleteIndex() byte {
	BCPS := gpu.CGBPallte[0]
	return BCPS & 0x3f
}

// FetchBGPalleteIncrement CGBのパレットインデックスが書き込み後にインクリメントするかを取得する
func (gpu *GPU) FetchBGPalleteIncrement() bool {
	BCPS := gpu.CGBPallte[0]
	return (BCPS >> 7) == 1
}

// FetchSPRPalleteIndex CGBのパレットインデックスを取得する
func (gpu *GPU) FetchSPRPalleteIndex() byte {
	OCPS := gpu.CGBPallte[1]
	return OCPS & 0x3f
}

// FetchSPRPalleteIncrement CGBのパレットインデックスが書き込み後にインクリメントするかを取得する
func (gpu *GPU) FetchSPRPalleteIncrement() bool {
	OCPS := gpu.CGBPallte[1]
	return (OCPS >> 7) == 1
}

// --------------------------------------------- internal method -----------------------------------------------------

func (gpu *GPU) fetchTileBaseAddr() uint16 {
	LCDC := gpu.LCDC
	if LCDC&0x10 != 0 {
		return 0x8000
	}
	return 0x8800
}

func (gpu *GPU) fetchSPRYSize() int {
	LCDC := gpu.LCDC
	if LCDC&0x04 != 0 {
		return 16
	}
	return 8
}

// ディスプレイにpixelデータをタイルの行単位でセットする
func (gpu *GPU) setTileLine(entryX, entryY int, lineIndex uint, addr uint16, tileType string, attr byte, spriteYSize int, isCGB bool) {
	// entryX, entryY: 何Pixel目を基準として配置するか
	var lowerByte, upperByte byte
	VRAMBankPtr := (attr >> 3) & 0x01
	lowerByte = gpu.VRAMBank[VRAMBankPtr][addr-0x8000]
	upperByte = gpu.VRAMBank[VRAMBankPtr][addr-0x8000+1]

	for j := 0; j < 8; j++ {
		bitCtr := (7 - uint(j)) // 上位何ビット目を取り出すか
		upperColor := (upperByte >> bitCtr) & 0x01
		lowerColor := (lowerByte >> bitCtr) & 0x01
		colorNumber := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

		var x, y int
		var c color.RGBA
		var RGB, R, G, B byte
		var isTransparent bool
		switch tileType {
		case "BGP":
			if isCGB {
				palleteNumber := attr & 0x07 // パレット番号 OBPn
				R, G, B, isTransparent = gpu.parseCGBPallete("BGP", palleteNumber, colorNumber)
			} else {
				RGB, isTransparent = gpu.parsePallete("BGP", colorNumber)
				R, G, B = colors[RGB][0], colors[RGB][1], colors[RGB][2]
			}
		case "OBP0":
			if isCGB {
				palleteNumber := attr & 0x07 // パレット番号 OBPn
				R, G, B, isTransparent = gpu.parseCGBPallete("OBP0", palleteNumber, colorNumber)
			} else {
				RGB, isTransparent = gpu.parsePallete("OBP0", colorNumber)
				R, G, B = colors[RGB][0], colors[RGB][1], colors[RGB][2]
			}
		case "OBP1":
			if isCGB {
				palleteNumber := attr & 0x07 // パレット番号 OBPn
				R, G, B, isTransparent = gpu.parseCGBPallete("OBP1", palleteNumber, colorNumber)
			} else {
				RGB, isTransparent = gpu.parsePallete("OBP1", colorNumber)
				R, G, B = colors[RGB][0], colors[RGB][1], colors[RGB][2]
			}
		}

		if !isTransparent {
			// 反転を考慮してpixelをセット
			if (attr>>6)&0x01 == 1 && (attr>>5)&0x01 == 1 {
				// 上下左右
				x = int(entryX + (7 - j))
				y = int(entryY + ((spriteYSize - 1) - int(lineIndex)))
			} else if (attr>>6)&0x01 == 1 {
				// 上下
				x = int(entryX + j)
				y = int(entryY + ((spriteYSize - 1) - int(lineIndex)))
			} else if (attr>>5)&0x01 == 1 {
				// 左右
				x = int(entryX + (7 - j))
				y = int(entryY + int(lineIndex))
			} else {
				// 反転無し
				x = int(entryX + j)
				y = int(entryY + int(lineIndex))
			}

			if (x >= 0 && x < 160) && (y >= 0 && y < 144) {
				if tileType == "BGP" {
					gpu.displayColor[y][x] = colorNumber

					if (attr>>7)&0x01 == 1 {
						gpu.BGPriorPixels = append(gpu.BGPriorPixels, [5]byte{byte(x), byte(y), R, G, B})
					} else {
						c = color.RGBA{R, G, B, 0xff}
						gpu.set(x, y, c)
					}
				} else {
					if (attr>>7)&0x01 == 0 && gpu.displayColor[y][x] != 0 {
						c = color.RGBA{R, G, B, 0xff}
						gpu.set(x, y, c)
					} else if gpu.displayColor[y][x] == 0 {
						c = color.RGBA{R, G, B, 0xff}
						gpu.set(x, y, c)
					}
				}
			}
		}
	}
}

func (gpu *GPU) parsePallete(name string, colorNumber byte) (RGB byte, transparent bool) {
	var pallete byte
	switch name {
	case "BGP":
		pallete = gpu.DMGPallte[0]
	case "OBP0":
		pallete = gpu.DMGPallte[1]
	case "OBP1":
		pallete = gpu.DMGPallte[2]
	default:
		errMsg := fmt.Sprintf("Error: BG Pallete name is invalid. %s", name)
		panic(errMsg)
	}

	transparent = false // 非透明

	switch colorNumber {
	case 0:
		RGB = (pallete >> 0) % 4
		if name == "OBP0" || name == "OBP1" {
			transparent = true
		}
	case 1:
		RGB = (pallete >> 2) % 4
	case 2:
		RGB = (pallete >> 4) % 4
	case 3:
		RGB = (pallete >> 6) % 4
	default:
		errMsg := fmt.Sprintf("Error: BG Pallete number is invalid. %d", colorNumber)
		panic(errMsg)
	}
	return RGB, transparent
}

func (gpu *GPU) parseCGBPallete(name string, palleteNumber, colorNumber byte) (R, G, B byte, transparent bool) {
	switch name {
	case "BGP":
		transparent = false
		i := palleteNumber*8 + colorNumber*2
		RGBLower, RGBUpper := uint16(gpu.BGPallete[i]), uint16(gpu.BGPallete[i+1])
		RGB := (RGBUpper << 8) | RGBLower
		R = byte(RGB & 0b11111)                 // bit 0-4
		G = byte((RGB & (0b11111 << 5)) >> 5)   // bit 5-9
		B = byte((RGB & (0b11111 << 10)) >> 10) // bit 10-14
	case "OBP0", "OBP1":
		if colorNumber == 0 {
			transparent = true
		} else {
			transparent = false
			i := palleteNumber*8 + colorNumber*2
			RGBLower, RGBUpper := uint16(gpu.SPRPallete[i]), uint16(gpu.SPRPallete[i+1])
			RGB := (RGBUpper << 8) | RGBLower
			R = byte(RGB & 0b11111)                 // bit 0-4
			G = byte((RGB & (0b11111 << 5)) >> 5)   // bit 5-9
			B = byte((RGB & (0b11111 << 10)) >> 10) // bit 10-14
		}
	}

	R *= 8
	G *= 8
	B *= 8
	return R, G, B, transparent
}
