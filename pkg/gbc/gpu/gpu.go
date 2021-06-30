package gpu

import (
	"gbc/pkg/util"
	"image"
	"image/color"

	hq2x "github.com/pokemium/hq2xgo"
)

type VRAM struct {
	Ptr  uint8
	Bank [2][0x2000]byte // 0x8000-0x9fff ゲームボーイカラーのみ
}

// GPU Graphic Processor Unit
type GPU struct {
	display       *image.RGBA    // 160*144のイメージデータ
	hq2x          *image.RGBA    // 320*288のイメージデータ(HQ2xかつ30fpsで使用)
	LCDC          byte           // LCD Control
	LCDSTAT       byte           // LCD Status
	Scroll        [2]byte        // Scrollの座標
	displayColor  [144][160]byte // 160*144の色番号(背景色を記録)
	Palette       Palette
	BGPriorPixels [][5]byte
	VRAM
	HBlankDMALength int
	Debug
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
	g.display, g.hq2x = image.NewRGBA(image.Rect(0, 0, 160, 144)), image.NewRGBA(image.Rect(0, 0, 320, 288))
	g.Debug.On = debug
	if debug {
		g.initTileData()
		g.OAM = image.NewRGBA(image.Rect(0, 0, 16*8-1, 20*5-3))
	}
}

// Display returns gameboy display data
func (g *GPU) Display(hq2x bool) *image.RGBA {
	if hq2x {
		return g.hq2x
	}
	return g.display
}

// GetOriginal - getter for display data in image.RGBA format. Function for debug.
func (g *GPU) GetOriginal() *image.RGBA {
	return g.display
}

// HQ2x - scaling display data using HQ2x
func (g *GPU) HQ2x() *image.RGBA {
	g.hq2x, _ = hq2x.HQ2x(g.display)
	return g.hq2x
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
