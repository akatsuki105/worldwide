package debug

import (
	"gbc/pkg/gbc/video"
	"gbc/pkg/util"
)

func (d *Debugger) SprView() [40][64 * 4]byte {
	buffer := [40][64 * 4]byte{}

	for i := 0; i < 40; i++ {
		for j := 0; j < 64; j++ {
			buffer[i][j*4] = 0xff
			buffer[i][j*4+1] = 0xff
			buffer[i][j*4+2] = 0xff
			buffer[i][j*4+3] = 0xff
		}

		y := int(d.g.Video.Oam.Get(4 * uint16(i)))
		if y <= 0 || y >= 160 {
			continue
		}

		objTile := int(d.g.Video.Oam.Get(4*uint16(i) + 2))
		attr := d.g.Video.Oam.Get(4*uint16(i) + 3)

		for y := 0; y < 8; y++ {
			tileDataLower := d.g.Video.VRAM.Buffer[(objTile*8+y)*2]
			tileDataUpper := d.g.Video.VRAM.Buffer[(objTile*8+y)*2+1]

			for x := 0; x < 8; x++ {
				b := 7 - x
				palIdx := uint16(((tileDataUpper>>b)&0b1)<<1) | uint16((tileDataLower>>b)&1) // 0 or 1 or 2 or 3
				base := video.PAL_OBJ + 4*util.Bool2U16(util.Bit(attr, 4))                   // 8*4 or 9*4
				p := d.g.Video.Renderer.Palette[d.g.Video.Renderer.Lookup[base+palIdx]]
				buffer[i][(y*8+x)*4], buffer[i][(y*8+x)*4+1], buffer[i][(y*8+x)*4+2] = byte((p&0b11111)*8), byte(((p>>5)&0b11111)*8), byte(((p>>10)&0b11111)*8)
			}
		}
	}

	return buffer
}
