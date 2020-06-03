package debug

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type CPU struct {
	all  int
	halt int
}

type Monitor struct {
	CPU
}

func (c *CPU) DrawUsage(screen *ebiten.Image, x, y int) {
	all, halt := c.all, c.halt
	usage := (all - halt) * 100 / all

	width, height := 20, 100
	gauge := image.NewRGBA(image.Rect(0, 0, width, height))
	rgba := color.RGBA{0x8f, 0x8f, 0x8f, 0xff}
	for h := 0; h < usage; h++ {
		for w := 0; w < width; w++ {
			gauge.Set(w, (height - h), rgba)
		}
	}

	gaugeEbiten, _ := ebiten.NewImageFromImage(gauge, ebiten.FilterDefault)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(gaugeEbiten, op)
}

func (c *CPU) Add(halt bool, count int) {
	c.all += count
	if halt {
		c.halt += count
	}
}

func (c *CPU) Reset() {
	c.all, c.halt = 1, 1
}
