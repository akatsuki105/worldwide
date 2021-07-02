package gpu

import (
	"fmt"
)

type Palette struct {
	DMGPalette            [3]byte // DMG's pal data {BGP, OGP0, OGP1}
	CGBPalette            [2]byte // CGB's pal data {BCPSIO, OCPSIO}
	BGPalette, SPRPalette [64]byte
}

// InitPalette init gameboy palette color
func InitPalette(color0, color1, color2, color3 [3]int) {
	DmgColor[0] = [3]uint8{uint8(color0[0]), uint8(color0[1]), uint8(color0[2])}
	DmgColor[1] = [3]uint8{uint8(color1[0]), uint8(color1[1]), uint8(color1[2])}
	DmgColor[2] = [3]uint8{uint8(color2[0]), uint8(color2[1]), uint8(color2[2])}
	DmgColor[3] = [3]uint8{uint8(color3[0]), uint8(color3[1]), uint8(color3[2])}
}

func (g *GPU) parsePallete(tileType int, colorIdx byte) (rgb byte, transparent bool) {
	pal := byte(0)
	transparent = false
	switch colorIdx {
	case 0:
		rgb, transparent = pal&0b11, tileType == OBP0 || tileType == OBP1
	case 1, 2, 3:
		rgb = (pal >> (2 * colorIdx)) & 0b11
	default:
		panic(fmt.Errorf("parsePallete Error: BG Pallete number is invalid. %d", colorIdx))
	}
	return rgb, transparent
}
