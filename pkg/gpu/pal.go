package gpu

import "fmt"

// InitPallete init gameboy pallete color
func InitPallete(color0, color1, color2, color3 [3]int) {
	colors[0] = [3]uint8{uint8(color0[0]), uint8(color0[1]), uint8(color0[2])}
	colors[1] = [3]uint8{uint8(color1[0]), uint8(color1[1]), uint8(color1[2])}
	colors[2] = [3]uint8{uint8(color2[0]), uint8(color2[1]), uint8(color2[2])}
	colors[3] = [3]uint8{uint8(color3[0]), uint8(color3[1]), uint8(color3[2])}
}

// FetchBGPalleteIndex CGBのパレットインデックスを取得する
func (g *GPU) FetchBGPalleteIndex() byte {
	BCPS := g.CGBPallte[0]
	return BCPS & 0x3f
}

// FetchBGPalleteIncrement CGBのパレットインデックスが書き込み後にインクリメントするかを取得する
func (g *GPU) FetchBGPalleteIncrement() bool {
	BCPS := g.CGBPallte[0]
	return (BCPS >> 7) == 1
}

// FetchSPRPalleteIndex CGBのパレットインデックスを取得する
func (g *GPU) FetchSPRPalleteIndex() byte {
	OCPS := g.CGBPallte[1]
	return OCPS & 0x3f
}

// FetchSPRPalleteIncrement CGBのパレットインデックスが書き込み後にインクリメントするかを取得する
func (g *GPU) FetchSPRPalleteIncrement() bool {
	OCPS := g.CGBPallte[1]
	return (OCPS >> 7) == 1
}

func (g *GPU) parsePallete(tileType int, colorNumber byte) (RGB byte, transparent bool) {
	var pallete byte
	switch tileType {
	case BGP:
		pallete = g.DMGPallte[0]
	case OBP0:
		pallete = g.DMGPallte[1]
	case OBP1:
		pallete = g.DMGPallte[2]
	default:
		errMsg := fmt.Sprintf("parsePallete Error: BG Pallete tile type is invalid. %d", tileType)
		panic(errMsg)
	}

	transparent = false // 非透明

	switch colorNumber {
	case 0:
		RGB = (pallete >> 0) % 4
		if tileType == OBP0 || tileType == OBP1 {
			transparent = true
		}
	case 1:
		RGB = (pallete >> 2) % 4
	case 2:
		RGB = (pallete >> 4) % 4
	case 3:
		RGB = (pallete >> 6) % 4
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
		RGBLower, RGBUpper := uint16(g.BGPallete[i]), uint16(g.BGPallete[i+1])
		RGB := (RGBUpper << 8) | RGBLower
		R = byte(RGB & 0b11111)                 // bit 0-4
		G = byte((RGB & (0b11111 << 5)) >> 5)   // bit 5-9
		B = byte((RGB & (0b11111 << 10)) >> 10) // bit 10-14
	case OBP0, OBP1:
		if colorNumber == 0 {
			transparent = true
		} else {
			i := palleteNumber*8 + colorNumber*2
			RGBLower, RGBUpper := uint16(g.SPRPallete[i]), uint16(g.SPRPallete[i+1])
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
