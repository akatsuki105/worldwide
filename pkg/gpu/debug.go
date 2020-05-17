package gpu

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

// --------------------------------------------- debug tiles -----------------------------------------------------

const (
	gridWidthX = 2
	gridWidthY = 3
)

func (g *GPU) InitTiles() {
	g.tileData.overall, _ = ebiten.NewImage(32*8+gridWidthY, 24*8+gridWidthX, ebiten.FilterDefault)
	g.tileData.overall.Fill(color.RGBA{255, 255, 255, 255})

	// gridを引く
	gridColor := color.RGBA{0x8f, 0x8f, 0x8f, 0xff}
	for y := 0; y < 24*8+gridWidthX; y++ {
		for i := 0; i < gridWidthY; i++ {
			g.tileData.overall.Set(16*8+i, y, gridColor)
		}
	}
	for x := 0; x < 32*8+gridWidthY; x++ {
		for i := 0; i < gridWidthX; i++ {
			// 横グリッドは2本
			g.tileData.overall.Set(x, 8*8+i, gridColor)
			g.tileData.overall.Set(x, 16*8+i, gridColor)
		}
	}

	for bank := 0; bank < 2; bank++ {
		for i := 0; i < 384; i++ {
			g.tileData.tiles[bank][i], _ = ebiten.NewImage(8, 8, ebiten.FilterDefault)
		}
	}
}

func (g *GPU) GetTileData() *ebiten.Image {
	return g.tileData.overall
}

func (g *GPU) UpdateTiles(isCGB bool) {
	itr := 1
	if isCGB {
		itr = 2
	}

	for bank := 0; bank < itr; bank++ {
		for i := 0; i < 384; i++ {

			tileAddr := 0x8000 + 16*i
			for y := 0; y < 8; y++ {
				addr := tileAddr + 2*y
				lowerByte, upperByte := g.VRAMBank[bank][addr-0x8000], g.VRAMBank[bank][addr-0x8000+1]

				for x := 0; x < 8; x++ {
					bitCtr := (7 - uint(x)) // 上位何ビット目を取り出すか
					upperColor := (upperByte >> bitCtr) & 0x01
					lowerColor := (lowerByte >> bitCtr) & 0x01
					colorNumber := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

					// 色番号からRGB値を算出する
					RGB, _ := g.parsePallete(OBP0, colorNumber)
					R, G, B := colors[RGB][0], colors[RGB][1], colors[RGB][2]
					c := color.RGBA{R, G, B, 0xff}

					// overall と 各タイルに対して
					overallX := bank*(16*8+gridWidthY) + (i%16)*8
					overallY := (i/16)*8 + (i/16)*gridWidthX/16
					g.tileData.overall.Set(overallX+x, overallY+y, c)
					g.tileData.tiles[bank][i].Set(x, y, c)
				}
			}
		}
	}
}
