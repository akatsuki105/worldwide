package gpu

import "fmt"

type Palette struct {
	DMGPallte  [3]byte // DMGのパレットデータ {BGP, OGP0, OGP1}
	CGBPallte  [2]byte // CGBのパレットデータ {BCPSIO, OCPSIO}
	BGPallete  [64]byte
	SPRPallete [64]byte
}

// InitPalette init gameboy pallete color
func InitPalette(color0, color1, color2, color3 [3]int) {
	colors[0] = [3]uint8{uint8(color0[0]), uint8(color0[1]), uint8(color0[2])}
	colors[1] = [3]uint8{uint8(color1[0]), uint8(color1[1]), uint8(color1[2])}
	colors[2] = [3]uint8{uint8(color2[0]), uint8(color2[1]), uint8(color2[2])}
	colors[3] = [3]uint8{uint8(color3[0]), uint8(color3[1]), uint8(color3[2])}
}

// FetchBGPalleteIndex CGBのパレットインデックスを取得する
func (g *GPU) FetchBGPalleteIndex() byte {
	BCPS := g.Palette.CGBPallte[0]
	return BCPS & 0x3f
}

// FetchBGPalleteIncrement CGBのパレットインデックスが書き込み後にインクリメントするかを取得する
func (g *GPU) FetchBGPalleteIncrement() bool {
	BCPS := g.Palette.CGBPallte[0]
	return (BCPS >> 7) == 1
}

// FetchSPRPalleteIndex CGBのパレットインデックスを取得する
func (g *GPU) FetchSPRPalleteIndex() byte {
	OCPS := g.Palette.CGBPallte[1]
	return OCPS & 0x3f
}

// FetchSPRPalleteIncrement CGBのパレットインデックスが書き込み後にインクリメントするかを取得する
func (g *GPU) FetchSPRPalleteIncrement() bool {
	OCPS := g.Palette.CGBPallte[1]
	return (OCPS >> 7) == 1
}

func (g *GPU) parsePallete(tileType int, colorNumber byte) (RGB byte, transparent bool) {
	pallete := g.Palette.DMGPallte[tileType]

	transparent = false // 非透明

	switch colorNumber {
	case 0:
		RGB = (pallete >> 0) % 4
		if tileType == OBP0 || tileType == OBP1 {
			transparent = true
		}
	case 1, 2, 3:
		RGB = (pallete >> (2 * colorNumber)) % 4
	default:
		errMsg := fmt.Sprintf("parsePallete Error: BG Pallete number is invalid. %d", colorNumber)
		panic(errMsg)
	}
	return RGB, transparent
}

func (g *GPU) parseCGBPallete(tileType int, palleteNumber, colorNumber byte) (R, G, B byte, transparent bool) {
	transparent = false
	switch tileType {
	case BGP:
		i := palleteNumber*8 + colorNumber*2
		RGBLower, RGBUpper := uint16(g.Palette.BGPallete[i]), uint16(g.Palette.BGPallete[i+1])
		RGB := (RGBUpper << 8) | RGBLower
		R = byte(RGB & 0b11111)                 // bit 0-4
		G = byte((RGB & (0b11111 << 5)) >> 5)   // bit 5-9
		B = byte((RGB & (0b11111 << 10)) >> 10) // bit 10-14
	case OBP0, OBP1:
		if colorNumber == 0 {
			transparent = true
		} else {
			i := palleteNumber*8 + colorNumber*2
			RGBLower, RGBUpper := uint16(g.Palette.SPRPallete[i]), uint16(g.Palette.SPRPallete[i+1])
			RGB := (RGBUpper << 8) | RGBLower
			R = byte(RGB & 0b11111)                 // bit 0-4
			G = byte((RGB & (0b11111 << 5)) >> 5)   // bit 5-9
			B = byte((RGB & (0b11111 << 10)) >> 10) // bit 10-14
		}
	}

	// 内部の色番号をRGB値に変換する
	R = R * 8
	G = G * 8
	B = B * 8
	return R, G, B, transparent
}
