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
func (g *GPU) Init() {
	g.display, _ = ebiten.NewImage(160, 144, ebiten.FilterDefault)
	g.original = image.NewRGBA(image.Rect(0, 0, 160, 144))
	g.hq2x, _ = ebiten.NewImage(320, 288, ebiten.FilterDefault)
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

// --------------------------------------------- internal method -----------------------------------------------------

func (g *GPU) fetchTileBaseAddr() uint16 {
	LCDC := g.LCDC
	if LCDC&0x10 != 0 {
		return 0x8000
	}
	return 0x8800
}

// ディスプレイにpixelデータをタイルの行単位でセットする 描画範囲外に出たときはfalseを返して描画を切り上げるべき旨を呼び出し元に伝える
func (g *GPU) setTileLine(entryX, entryY int, lineIndex uint, addr uint16, tileType int, attr byte, spriteYSize int, isCGB bool, OAMindex int) bool {

	// entryX, entryY: 何Pixel目を基準として配置するか
	VRAMBankPtr := (attr >> 3) & 0x01
	if !isCGB {
		VRAMBankPtr = 0
	}

	lowerByte, upperByte := g.VRAMBank[VRAMBankPtr][addr-0x8000], g.VRAMBank[VRAMBankPtr][addr-0x8000+1]

	for j := 0; j < 8; j++ {
		bitCtr := (7 - uint(j)) // 上位何ビット目を取り出すか
		upperColor := (upperByte >> bitCtr) & 0x01
		lowerColor := (lowerByte >> bitCtr) & 0x01
		colorNumber := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

		var x, y int
		var RGB, R, G, B byte
		var isTransparent bool

		// 色番号からRGB値を算出する
		if isCGB {
			palleteNumber := attr & 0x07 // パレット番号 OBPn
			R, G, B, isTransparent = g.parseCGBPallete(tileType, palleteNumber, colorNumber)
		} else {
			RGB, isTransparent = g.parsePallete(tileType, colorNumber)
			R, G, B = colors[RGB][0], colors[RGB][1], colors[RGB][2]
		}
		c := color.RGBA{R, G, B, 0xff}

		var deltaX, deltaY int
		if !isTransparent {
			// 反転を考慮してpixelをセット
			if (attr>>6)&0x01 == 1 && (attr>>5)&0x01 == 1 {
				// 上下左右
				deltaX = int((7 - j))
				deltaY = int(((spriteYSize - 1) - int(lineIndex)))
			} else if (attr>>6)&0x01 == 1 {
				// 上下
				deltaX = int(j)
				deltaY = int(((spriteYSize - 1) - int(lineIndex)))
			} else if (attr>>5)&0x01 == 1 {
				// 左右
				deltaX = int((7 - j))
				deltaY = int(int(lineIndex))
			} else {
				// 反転無し
				deltaX = int(j)
				deltaY = int(int(lineIndex))
			}
			x = entryX + deltaX
			y = entryY + deltaY

			// debug OAM
			if OAMindex >= 0 {
				col := OAMindex % 8
				row := OAMindex / 8
				g.OAM.Set(col*16+deltaX+2, row*20+deltaY, c)
			}

			if (x >= 0 && x < 160) && (y >= 0 && y < 144) {
				if tileType == BGP {
					g.displayColor[y][x] = colorNumber
					if (attr>>7)&0x01 == 1 {
						g.BGPriorPixels = append(g.BGPriorPixels, [5]byte{byte(x), byte(y), R, G, B})
					}
					g.set(x, y, c)
				} else {
					// スプライト
					if (attr>>7)&0x01 == 0 && g.displayColor[y][x] != 0 {
						g.set(x, y, c)
					} else if g.displayColor[y][x] == 0 {
						g.set(x, y, c)
					}
				}
			} else if x >= 160 {
				break
			} else if y >= 144 {
				return false
			}
		}
	}

	return true
}
