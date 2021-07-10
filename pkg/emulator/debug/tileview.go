package debug

const TILENUM = 192 // cgb -> 192*2 (bank1)

func (d *Debugger) TileView() [2][]byte {
	buffer := [2][]byte{make([]byte, (0x0800*4)*3), make([]byte, (0x0800*4)*3)} // bank0, bank1

	for bank := 0; bank < 2; bank++ {
		if bank == 1 && !d.g.Cartridge.IsCGB {
			break
		}

		for i := 0; i < TILENUM; i++ {
			base := 16 * i

			for y := 0; y < 8; y++ {
				addr := base + 2*y
				tileDataLower, tileDataUpper := d.g.Video.VRAM.Buffer[addr+0x2000*bank], d.g.Video.VRAM.Buffer[addr+0x2000*bank+1]

				for x := 0; x < 8; x++ {
					b := (7 - uint(x))
					upperColor := (tileDataUpper >> b) & 0x01
					lowerColor := (tileDataLower >> b) & 0x01
					colorIdx := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3
					red, green, blue := byte(0xff), byte(0xff), byte(0xff)
					switch colorIdx {
					case 1:
						red, green, blue = 0xad, 0xad, 0xad
					case 2:
						red, green, blue = 0x52, 0x52, 0x52
					case 3:
						red, green, blue = 0x00, 0x00, 0x00
					}

					buffer[bank][i*64*4+y*8*4+x*4] = red
					buffer[bank][i*64*4+y*8*4+x*4+1] = green
					buffer[bank][i*64*4+y*8*4+x*4+2] = blue
					buffer[bank][i*64*4+y*8*4+x*4+3] = 0xff
				}
			}
		}
	}

	return buffer
}
