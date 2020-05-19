package gpu

import (
	"image"
	"image/color"

	hq2x "github.com/Akatsuki-py/hq2xgo"
	"github.com/hajimehoshi/ebiten"
)

// GPU Graphic Processor Unit
type GPU struct {
	display       *ebiten.Image  // 160*144のイメージデータ
	original      *image.RGBA    // 160*144のイメージデータ
	hq2x          *ebiten.Image  // 320*288のイメージデータ(HQ2xかつ30fpsで使用)
	tileData      tileData       // タイルデータ
	LCDC          byte           // LCD Control
	LCDSTAT       byte           // LCD Status
	Scroll        [2]byte        // Scrollの座標
	displayColor  [144][160]byte // 160*144の色番号(背景色を記録)
	Palette       Palette
	BGPriorPixels [][5]byte
	// VRAM bank
	VRAMBankPtr     uint8
	VRAMBank        [2][0x2000]byte // 0x8000-0x9fff ゲームボーイカラーのみ
	HBlankDMALength int
	OAM             *ebiten.Image // OAMをまとめたもの
	debug           bool          // デバッグモードか
}

type tileData struct {
	overall *ebiten.Image         // タイルデータをいちまいの画像にまとめたもの
	tiles   [2][384]*ebiten.Image // 8*8のタイルデータの一覧
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
	g.display, _ = ebiten.NewImage(160, 144, ebiten.FilterDefault)
	g.original = image.NewRGBA(image.Rect(0, 0, 160, 144))
	g.hq2x, _ = ebiten.NewImage(320, 288, ebiten.FilterDefault)

	g.debug = debug
	if debug {
		g.initDebugTiles()
	}
}

// GetDisplay getter for display data
func (g *GPU) GetDisplay(hq2x bool) *ebiten.Image {
	if hq2x {
		return g.hq2x
	}
	return g.display
}

// HQ2x - scaling display data using HQ2x
func (g *GPU) HQ2x() *ebiten.Image {
	tmp, _ := hq2x.HQ2x(g.original)
	g.hq2x, _ = ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
	return g.hq2x
}

func (g *GPU) set(x, y int, c color.RGBA) {
	g.display.Set(x, y, c)
	g.original.SetRGBA(x, y, c)
}

func (g *GPU) fetchTileBaseAddr() uint16 {
	LCDC := g.LCDC
	if LCDC&0x10 != 0 {
		return 0x8000
	}
	return 0x8800
}
