package gpu

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

func (g *GPU) UpdateTileData(isCGB bool) {
	itr := 1
	if isCGB {
		itr = 2
	}

	for bank := 0; bank < itr; bank++ {
		for i := 0; i < 384; i++ {

			tileAddr := 0x8000 + 16*i
			for y := 0; y < 8; y++ {
				addr := uint16(tileAddr + 2*y)
				lowerByte, upperByte := g.VRAM.Buffer[addr-0x8000+0x2000*uint16(bank)], g.VRAM.Buffer[addr-0x8000+1+0x2000*uint16(bank)]

				for x := 0; x < 8; x++ {
					bitCtr := (7 - uint(x)) // 上位何ビット目を取り出すか
					upperColor := (upperByte >> bitCtr) & 0x01
					lowerColor := (lowerByte >> bitCtr) & 0x01
					colorNumber := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

					// 色番号からRGB値を算出する
					RGB, _ := g.parsePallete(OBP0, colorNumber)
					R, G, B := DmgColor[RGB][0], DmgColor[RGB][1], DmgColor[RGB][2]
					c := color.RGBA{R, G, B, 0xff}

					// overall と 各タイルに対して
					overallX := bank*(16*8+gridWidthY) + (i%16)*8
					overallY := (i/16)*8 + (i/16)*gridWidthX/16
					g.Debug.tileData.overall.Set(overallX+x, overallY+y, c)
					g.Debug.tileData.tiles[bank][i].Set(x, y, c)
				}
			}
		}
	}
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
