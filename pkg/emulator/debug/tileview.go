package debug

import (
	"gbc/pkg/gbc"
)

const TILENUM = 384 // cgb -> 384*2 (bank1)

var (
	// colors {R, G, B}
	colors [4][3]byte = [4][3]byte{
		{0xff, 0xff, 0xff}, {0xad, 0xad, 0xad}, {0x52, 0x52, 0x52}, {0x00, 0x00, 0x00},
	}
)

func (d *Debugger) TileView() [2][]byte {
	buffer := [2][]byte{make([]byte, TILENUM*64*4), make([]byte, TILENUM*64*4)} // bank0, bank1

	for bank := 0; bank < 2; bank++ {
		if bank == 1 && !d.g.Cartridge.IsCGB {
			break
		}

		for i := 0; i < TILENUM; i++ {
			addr := 16 * i

			for y := 0; y < 8; y++ {
				tileAddr := addr + 2*y + 0x2000*bank
				tileDataLower, tileDataUpper := d.g.Video.VRAM.Buffer[tileAddr], d.g.Video.VRAM.Buffer[tileAddr+1]

				for x := 0; x < 8; x++ {
					b := (7 - uint(x))
					upperColor := (tileDataUpper >> b) & 0x01
					lowerColor := (tileDataLower >> b) & 0x01
					palIdx := (upperColor << 1) | lowerColor // 0 or 1 or 2 or 3
					colorIdx := (d.g.IO[gbc.BGPIO] >> (palIdx * 2)) & 0b11
					red, green, blue := colors[colorIdx][0], colors[colorIdx][1], colors[colorIdx][2]
					bufferIdx := i*64*4 + y*8*4 + x*4
					buffer[bank][bufferIdx] = red
					buffer[bank][bufferIdx+1] = green
					buffer[bank][bufferIdx+2] = blue
					buffer[bank][bufferIdx+3] = 0xff
				}
			}
		}
	}

	return buffer
}
