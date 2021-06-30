package debug

import (
	"image"
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type GBC struct {
	all  int
	halt int
}

type Monitor struct {
	GBC
}

func (c *GBC) DrawUsage(screen *ebiten.Image, x, y int, boost bool) {
	all, halt := c.all, c.halt
	usage := (all - halt) * 100 / all

	width, height := 10+2, 100+2
	gauge := image.NewRGBA(image.Rect(0, 0, width, height))
	// Gauge border
	border := color.White
	for h := 0; h < height; h++ {
		gauge.Set(0, h, border)
		gauge.Set(width-1, h, border)
	}
	for w := 0; w < width; w++ {
		gauge.Set(w, 0, border)
		gauge.Set(w, height-1, border)
	}

	// Usage bar
	rgba := color.RGBA{0x00, 0xe9, 0x21, 0xff}
	if boost {
		rgba = color.RGBA{0xff, 0xd7, 0x00, 0xff}
	}
	for h := 0; h < usage; h++ {
		for w := 0; w < width-2; w++ {
			gauge.Set(w+1, (height - (h + 2)), rgba)
		}
	}

	gaugeEbiten := ebiten.NewImageFromImage(gauge)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(gaugeEbiten, op)
}

func (c *GBC) Add(halt bool, count int) {
	c.all += count
	if halt {
		c.halt += count
	}
}

func (c *GBC) Reset() {
	c.all, c.halt = 1, 1
}
