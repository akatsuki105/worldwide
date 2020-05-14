package gpu

import (
	"fmt"
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
	DMGPallte     [3]byte        // DMGのパレットデータ {BGP, OGP0, OGP1}
	CGBPallte     [2]byte        // CGBのパレットデータ {BCPSIO, OCPSIO}
	BGPallete     [64]byte
	SPRPallete    [64]byte
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

// InitPallete init gameboy pallete color
func InitPallete(color0, color1, color2, color3 [3]int) {
	colors[0] = [3]uint8{uint8(color0[0]), uint8(color0[1]), uint8(color0[2])}
	colors[1] = [3]uint8{uint8(color1[0]), uint8(color1[1]), uint8(color1[2])}
	colors[2] = [3]uint8{uint8(color2[0]), uint8(color2[1]), uint8(color2[2])}
	colors[3] = [3]uint8{uint8(color3[0]), uint8(color3[1]), uint8(color3[2])}
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

// --------------------------------------------- Render -----------------------------------------------------

// SetBGLine 1タイルライン描画する
func (g *GPU) SetBGLine(entryX, entryY int, tileX, tileY uint, useWindow, isCGB bool, lineIndex int) bool {
	index := tileX + tileY*32 // マップの何タイル目か

	// タイル番号からタイルデータのあるアドレス取得
	var addr uint16
	LCDC := g.LCDC
	if useWindow {
		if LCDC&0x40 != 0 {
			addr = 0x9c00 + uint16(index)
		} else {
			addr = 0x9800 + uint16(index)
		}
	} else {
		if LCDC&0x08 != 0 {
			addr = 0x9c00 + uint16(index)
		} else {
			addr = 0x9800 + uint16(index)
		}
	}
	tileIndex := uint8(g.VRAMBank[0][addr-0x8000])
	baseAddr := g.fetchTileBaseAddr()
	if baseAddr == 0x8800 {
		tileIndex = uint8(int(int8(tileIndex)) + 128)
	}

	// 背景属性取得
	var attr byte
	if isCGB {
		attr = uint8(g.VRAMBank[1][addr-0x8000])
	} else {
		attr = 0
	}

	index16 := uint16(tileIndex)*8 + uint16(lineIndex) // 何枚目のタイルか*8 + タイルの何行目か
	addr = uint16(baseAddr + 2*index16)
	return g.setTileLine(entryX, entryY, uint(lineIndex), addr, BGP, attr, 8, isCGB, -1)
}

// SetSPRTile スプライトを出力する
func (g *GPU) SetSPRTile(OAMindex, entryX, entryY int, tileIndex uint, attr byte, isCGB bool) {
	spriteYSize := g.fetchSPRYSize()
	if (attr>>4)%2 == 1 {
		for lineIndex := 0; lineIndex < spriteYSize; lineIndex++ {
			index := uint16(tileIndex)*8 + uint16(lineIndex) // 何枚目のタイルか*8 + タイルの何行目か
			addr := uint16(0x8000 + 2*index)                 // スプライトは0x8000のみ
			continueFlag := g.setTileLine(entryX, entryY, uint(lineIndex), addr, OBP1, attr, spriteYSize, isCGB, OAMindex)
			if !continueFlag {
				break
			}
		}
	} else {
		for lineIndex := 0; lineIndex < spriteYSize; lineIndex++ {
			index := uint16(tileIndex)*8 + uint16(lineIndex) // 何枚目のタイルか*8 + タイルの何行目か
			addr := uint16(0x8000 + 2*index)                 // スプライトは0x8000のみ
			continueFlag := g.setTileLine(entryX, entryY, uint(lineIndex), addr, OBP0, attr, spriteYSize, isCGB, OAMindex)
			if !continueFlag {
				break
			}
		}
	}
}

// SetBGPriorPixels 背景優先の背景を描画するための関数
func (g *GPU) SetBGPriorPixels() {
	for _, pixel := range g.BGPriorPixels {
		x, y := int(pixel[0]), int(pixel[1])
		R, G, B := pixel[2], pixel[3], pixel[4]
		c := color.RGBA{R, G, B, 0xff}
		if x < 160 && y < 144 {
			g.set(x, y, c)
		}
	}
	g.BGPriorPixels = [][5]byte{}
}

// --------------------------------------------- CGB pallete -----------------------------------------------------

// FetchBGPalleteIndex CGBのパレットインデックスを取得する
func (g *GPU) FetchBGPalleteIndex() byte {
	BCPS := g.CGBPallte[0]
	return BCPS & 0x3f
}

// FetchBGPalleteIncrement CGBのパレットインデックスが書き込み後にインクリメントするかを取得する
func (g *GPU) FetchBGPalleteIncrement() bool {
	BCPS := g.CGBPallte[0]
	return (BCPS >> 7) == 1
}

// FetchSPRPalleteIndex CGBのパレットインデックスを取得する
func (g *GPU) FetchSPRPalleteIndex() byte {
	OCPS := g.CGBPallte[1]
	return OCPS & 0x3f
}

// FetchSPRPalleteIncrement CGBのパレットインデックスが書き込み後にインクリメントするかを取得する
func (g *GPU) FetchSPRPalleteIncrement() bool {
	OCPS := g.CGBPallte[1]
	return (OCPS >> 7) == 1
}

// --------------------------------------------- scroll method -----------------------------------------------------

// ReadScroll - スクロール値を得る
func (g *GPU) ReadScroll() (x, y uint) {
	x, y = uint(g.Scroll[0]), uint(g.Scroll[1])
	return x, y
}

// WriteScrollX - スクロールのX座標を書き込む
func (g *GPU) WriteScrollX(x byte) {
	g.Scroll[0] = x
}

// WriteScrollY - スクロールのY座標を書き込む
func (g *GPU) WriteScrollY(y byte) {
	g.Scroll[1] = y
}

// --------------------------------------------- internal method -----------------------------------------------------

func (g *GPU) fetchTileBaseAddr() uint16 {
	LCDC := g.LCDC
	if LCDC&0x10 != 0 {
		return 0x8000
	}
	return 0x8800
}

func (g *GPU) fetchSPRYSize() int {
	LCDC := g.LCDC
	if LCDC&0x04 != 0 {
		return 16
	}
	return 8
}

// ディスプレイにpixelデータをタイルの行単位でセットする 描画範囲外に出たときはfalseを返して描画を切り上げるべき旨を呼び出し元に伝える
func (g *GPU) setTileLine(entryX, entryY int, lineIndex uint, addr uint16, tileType int, attr byte, spriteYSize int, isCGB bool, OAMindex int) bool {

	// entryX, entryY: 何Pixel目を基準として配置するか
	VRAMBankPtr := (attr >> 3) & 0x01
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
				g.OAM.Set(col*16+deltaX, row*20+deltaY, c)
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

func (g *GPU) parsePallete(tileType int, colorNumber byte) (RGB byte, transparent bool) {
	var pallete byte
	switch tileType {
	case BGP:
		pallete = g.DMGPallte[0]
	case OBP0:
		pallete = g.DMGPallte[1]
	case OBP1:
		pallete = g.DMGPallte[2]
	default:
		errMsg := fmt.Sprintf("parsePallete Error: BG Pallete tile type is invalid. %d", tileType)
		panic(errMsg)
	}

	transparent = false // 非透明

	switch colorNumber {
	case 0:
		RGB = (pallete >> 0) % 4
		if tileType == OBP0 || tileType == OBP1 {
			transparent = true
		}
	case 1:
		RGB = (pallete >> 2) % 4
	case 2:
		RGB = (pallete >> 4) % 4
	case 3:
		RGB = (pallete >> 6) % 4
	default:
		errMsg := fmt.Sprintf("parsePallete Error: BG Pallete number is invalid. %d", colorNumber)
		panic(errMsg)
	}
	return RGB, transparent
}

func (g *GPU) parseCGBPallete(tileType int, palleteNumber, colorNumber byte) (R, G, B byte, transparent bool) {
	transparent = false
	switch tileType {
	case BGP:
		i := palleteNumber*8 + colorNumber*2
		RGBLower, RGBUpper := uint16(g.BGPallete[i]), uint16(g.BGPallete[i+1])
		RGB := (RGBUpper << 8) | RGBLower
		R = byte(RGB & 0b11111)                 // bit 0-4
		G = byte((RGB & (0b11111 << 5)) >> 5)   // bit 5-9
		B = byte((RGB & (0b11111 << 10)) >> 10) // bit 10-14
	case OBP0, OBP1:
		if colorNumber == 0 {
			transparent = true
		} else {
			i := palleteNumber*8 + colorNumber*2
			RGBLower, RGBUpper := uint16(g.SPRPallete[i]), uint16(g.SPRPallete[i+1])
			RGB := (RGBUpper << 8) | RGBLower
			R = byte(RGB & 0b11111)                 // bit 0-4
			G = byte((RGB & (0b11111 << 5)) >> 5)   // bit 5-9
			B = byte((RGB & (0b11111 << 10)) >> 10) // bit 10-14
		}
	}

	// 内部の色番号をRGB値に変換する
	R = R * 8
	G = G * 8
	B = B * 8
	return R, G, B, transparent
}

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
