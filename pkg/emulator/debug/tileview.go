package debug

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

const TILENUM = 384 // cgb -> 384*2 (bank1)

func (d *Debugger) getRawTileView(bank int) []byte {
	buffer := make([]byte, TILENUM*64*4)

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
				p := d.g.Video.Palette[d.g.Video.Renderer.Lookup[palIdx]]
				red, green, blue := byte((p&0b11111)*8), byte(((p>>5)&0b11111)*8), byte(((p>>10)&0b11111)*8)
				bufferIdx := i*64*4 + y*8*4 + x*4
				buffer[bufferIdx] = red
				buffer[bufferIdx+1] = green
				buffer[bufferIdx+2] = blue
				buffer[bufferIdx+3] = 0xff
			}
		}
	}

	return buffer
}

func (d *Debugger) getTileView(bank int) []byte {
	rawBuffer := d.getRawTileView(bank)
	m := image.NewRGBA(image.Rect(0, 0, 8*16, 8*384/TILE_PER_ROW))
	var wg sync.WaitGroup
	wg.Add(384 / TILE_PER_ROW)
	for row := 0; row < 384/TILE_PER_ROW; row++ {
		// 0..63, 0..63, 0..63, .. -> 0..7, 0..7, ... 8..15, 8..15,
		go func(row int) {
			rowStart, rowEnd := row*TILE_PER_ROW, (row+1)*TILE_PER_ROW
			rowBuffer := rawBuffer[rowStart*64*4 : rowEnd*64*4]

			for t := 0; t < TILE_PER_ROW; t++ {
				rowBufferBase := t * 64 * 4
				for y := 0; y < 8; y++ {
					tileRowBuffer := rowBuffer[rowBufferBase+y*8*4 : rowBufferBase+(y+1)*8*4] // (y*8)..((y*8)+7)
					for x := 0; x < 8; x++ {
						m.SetRGBA(t*8+x, row*8+y, color.RGBA{tileRowBuffer[x*4], tileRowBuffer[x*4+1], tileRowBuffer[x*4+2], tileRowBuffer[x*4+3]})
					}
				}
			}
			wg.Done()
		}(row)
	}
	wg.Wait()

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, m, nil); err != nil {
		log.Println("unable to encode image.")
	}

	return buffer.Bytes()
}

func (d *Debugger) tileView(ws *websocket.Conn, bank int) {
	err := websocket.Message.Send(ws, d.getTileView(bank))
	if err != nil {
		log.Printf("error sending data: %v\n", err)
		return
	}

	for range time.NewTicker(time.Millisecond * 100).C {
		err := websocket.Message.Send(ws, d.getTileView(bank))
		if err != nil {
			log.Printf("error sending data: %v\n", err)
			return
		}
	}
}

func (d *Debugger) TileView0(ws *websocket.Conn) {
	d.tileView(ws, 0)
}

func (d *Debugger) TileView1(ws *websocket.Conn) {
	d.tileView(ws, 1)
}
