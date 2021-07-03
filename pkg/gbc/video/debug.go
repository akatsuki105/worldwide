package video

import (
	"image"
	"image/color"
	"image/draw"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type tileData struct {
	overall *image.RGBA         // タイルデータをいちまいの画像にまとめたもの
	tiles   [2][384]*image.RGBA // 8*8のタイルデータの一覧
}

type Debug struct {
	On          bool
	tileData    tileData
	OAM         *image.RGBA // OAMをまとめたもの
	oamProperty [40][4]byte
	bgMap       *image.RGBA // 背景のみ
}

const (
	gridWidthX = 2
	gridWidthY = 3
)

func (d *Debug) initTileData() {
	d.tileData.overall = image.NewRGBA(image.Rect(0, 0, 32*8+gridWidthY, 24*8+gridWidthX))

	// draw grid
	gridColor := color.RGBA{0x8f, 0x8f, 0x8f, 0xff}
	for y := 0; y < 24*8+gridWidthX; y++ {
		for i := 0; i < gridWidthY; i++ {
			d.tileData.overall.Set(16*8+i, y, gridColor)
		}
	}
	for x := 0; x < 32*8+gridWidthY; x++ {
		for i := 0; i < gridWidthX; i++ {
			// horizontal grid has two line
			d.tileData.overall.Set(x, 8*8+i, gridColor)
			d.tileData.overall.Set(x, 16*8+i, gridColor)
		}
	}

	for bank := 0; bank < 2; bank++ {
		for i := 0; i < 384; i++ {
			d.tileData.tiles[bank][i] = image.NewRGBA(image.Rect(0, 0, 8, 8))
		}
	}
}

func (d *Debug) GetTileData() *ebiten.Image {
	result := ebiten.NewImageFromImage(d.tileData.overall)
	return result
}

func (d *Debug) BGMap() *image.RGBA      { return d.bgMap }
func (d *Debug) SetBGMap(bg *image.RGBA) { d.bgMap = bg }

func (d *Debug) FillOAM() {
	c := color.RGBA{0x8f, 0x8f, 0x8f, 0xff}
	draw.Draw(d.OAM, d.OAM.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
}

func (d *Debug) OAMProperty(index int) (byte, byte, byte, byte) {
	Y, X := d.oamProperty[index][0], d.oamProperty[index][1]
	tileIndex := d.oamProperty[index][2]
	attr := d.oamProperty[index][3]
	return Y, X, tileIndex, attr
}
func (d *Debug) SetOAMProperty(index int, X, Y, tileIndex, attr byte) {
	d.oamProperty[index] = [4]byte{Y, X, tileIndex, attr}
}
