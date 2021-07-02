package gpu

import (
	"fmt"
)

type Palette struct {
	DMGPalette            [3]byte // DMG's pal data {BGP, OGP0, OGP1}
	CGBPalette            [2]byte // CGB's pal data {BCPSIO, OCPSIO}
	BGPalette, SPRPalette [64]byte
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
