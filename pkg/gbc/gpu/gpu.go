package gpu

import (
	"gbc/pkg/util"
	"image"
	"image/color"
)

type VRAM struct {
	Bank   uint16       // 0 or 1
	Buffer [0x4000]byte // (0x8000-0x9fff)x2 (using bank on CGB)
}

// GPU Graphic Processor Unit
type GPU struct {
	display       *image.RGBA    // 160*144
	LCDC          byte           // LCD Control
	LCDSTAT       byte           // LCD Status
	Scroll        [2]byte        // Scroll coord
	displayColor  [144][160]byte // 160*144 color data
	Palette       Palette
	BGPriorPixels [][5]byte
	VRAM
	HBlankDMALength int
	Debug
	renderer *Renderer
	oam      *OAM
}

var (
	// colors {R, G, B}
	colors [4][3]uint8 = [4][3]uint8{
		{175, 197, 160}, {93, 147, 66}, {22, 63, 48}, {0, 40, 0},
	}
)

const (
	BGP = iota
	OBP0
	OBP1
)

// Init GPU
func (g *GPU) Init(debug bool) {
	g.display = image.NewRGBA(image.Rect(0, 0, 160, 144))
	g.Debug.On = debug
	if debug {
		g.initTileData()
		g.OAM = image.NewRGBA(image.Rect(0, 0, 16*8-1, 20*5-3))
	}
}

// Display returns gameboy display data
func (g *GPU) Display() *image.RGBA {
	return g.display
}

// GetOriginal - getter for display data in image.RGBA format. Function for debug.
func (g *GPU) GetOriginal() *image.RGBA {
	return g.display
}

func (g *GPU) set(x, y int, c color.RGBA) {
	g.display.SetRGBA(x, y, c)
}

func (g *GPU) fetchTileBaseAddr() uint16 {
	if util.Bit(g.LCDC, 4) {
		return 0x8000
	}
	return 0x8800
}
