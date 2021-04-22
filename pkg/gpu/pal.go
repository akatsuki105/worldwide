package gpu

import (
	"fmt"
	"gbc/pkg/util"
)

type Palette struct {
	DMGPalette            [3]byte // DMG's pal data {BGP, OGP0, OGP1}
	CGBPalette            [2]byte // CGB's pal data {BCPSIO, OCPSIO}
	BGPalette, SPRPalette [64]byte
}

// InitPalette init gameboy palette color
func InitPalette(color0, color1, color2, color3 [3]int) {
	colors[0] = [3]uint8{uint8(color0[0]), uint8(color0[1]), uint8(color0[2])}
	colors[1] = [3]uint8{uint8(color1[0]), uint8(color1[1]), uint8(color1[2])}
	colors[2] = [3]uint8{uint8(color2[0]), uint8(color2[1]), uint8(color2[2])}
	colors[3] = [3]uint8{uint8(color3[0]), uint8(color3[1]), uint8(color3[2])}
}

// BgPalIdx returns bg palette index for CGB
func (g *GPU) BgPalIdx() byte {
	BCPS := g.Palette.CGBPalette[0]
	return BCPS & 0x3f
}

// isBGPalIncrement returns whether bg palette index is incremented after write
func (g *GPU) IsBgPalIncrement() bool {
	return util.Bit(g.Palette.CGBPalette[0], 7)
}

// SprPalIdx returns spr palette index for CGB
func (g *GPU) SprPalIdx() byte {
	OCPS := g.Palette.CGBPalette[1]
	return OCPS & 0x3f
}

// isBGPalIncrement returns whether spr palette index is incremented after write
func (g *GPU) IsSprPalIncrement() bool {
	return util.Bit(g.Palette.CGBPalette[1], 7)
}

func (g *GPU) parsePallete(tileType int, colorIdx byte) (rgb byte, transparent bool) {
	pal := g.Palette.DMGPalette[tileType]
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

func (g *GPU) parseCGBPallete(tileType int, palIdx, colorIdx byte) (R, G, B byte, transparent bool) {
	transparent = false
	switch tileType {
	case BGP:
		i := palIdx*8 + colorIdx*2
		RGBLower, RGBUpper := uint16(g.Palette.BGPalette[i]), uint16(g.Palette.BGPalette[i+1])
		RGB := (RGBUpper << 8) | RGBLower
		R = byte(RGB & 0b11111)                 // bit 0-4
		G = byte((RGB & (0b11111 << 5)) >> 5)   // bit 5-9
		B = byte((RGB & (0b11111 << 10)) >> 10) // bit 10-14
	case OBP0, OBP1:
		if colorIdx == 0 {
			transparent = true
		} else {
			i := palIdx*8 + colorIdx*2
			RGBLower, RGBUpper := uint16(g.Palette.SPRPalette[i]), uint16(g.Palette.SPRPalette[i+1])
			RGB := (RGBUpper << 8) | RGBLower
			R = byte(RGB & 0b11111)                 // bit 0-4
			G = byte((RGB & (0b11111 << 5)) >> 5)   // bit 5-9
			B = byte((RGB & (0b11111 << 10)) >> 10) // bit 10-14
		}
	}

	R, G, B = R*8, G*8, B*8 // color idx -> RGB value
	return R, G, B, transparent
}
