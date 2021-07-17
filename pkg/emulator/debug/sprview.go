package debug

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"sync"
	"time"

	"github.com/pokemium/worldwide/pkg/gbc/video"
	"github.com/pokemium/worldwide/pkg/util"
	"golang.org/x/net/websocket"
)

func (d *Debugger) getRawSprView() [40][64 * 4]byte {
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

func (d *Debugger) getSprView() []byte {
	rawBuffer := d.getRawSprView()
	m := image.NewRGBA(image.Rect(0, 0, 8*8, 8*5))

	var wg sync.WaitGroup
	wg.Add(40)
	for i := 0; i < 40; i++ {
		go func(i int) {
			col, row := i&0x7, i/8
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					idx := (y*8 + x) * 4
					m.Set(col*8+x, row*8+y, color.RGBA{rawBuffer[i][idx], rawBuffer[i][idx+1], rawBuffer[i][idx+2], rawBuffer[i][idx+3]})
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, m, nil); err != nil {
		log.Println("unable to encode image.")
	}

	return buffer.Bytes()
}

func (d *Debugger) SprView(ws *websocket.Conn) {
	err := websocket.Message.Send(ws, d.getSprView())
	if err != nil {
		log.Printf("error sending data: %v\n", err)
		return
	}

	for range time.NewTicker(time.Second).C {
		err := websocket.Message.Send(ws, d.getSprView())
		if err != nil {
			log.Printf("error sending data: %v\n", err)
			return
		}
	}
}
