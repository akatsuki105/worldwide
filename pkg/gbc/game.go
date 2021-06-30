package gbc

import (
	"gbc/pkg/gbc/gpu"
	"gbc/pkg/util"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

func (g *GBC) Update() error {
	if frames == 0 {
		setIcon()
		g.debug.monitor.GBC.Reset()
	}
	if frames%3 == 0 {
		g.handleJoypad()
	}

	frames++
	g.debug.monitor.GBC.Reset()

	p, b := &g.debug.pause, &g.debug.Break
	if p.Delay() {
		p.DecrementDelay()
	}
	if p.On() || b.On() {
		return nil
	}

	skipRender = (g.Config.Display.FPS30) && (frames%2 == 1)

	LCDC := g.FetchMemory8(LCDCIO)
	scrollX, scrollY := uint(g.GPU.Scroll[0]), uint(g.GPU.Scroll[1])
	scrollPixelX := scrollX % 8

	iterX, iterY := width, height
	if scrollPixelX > 0 {
		iterX += 8
	}

	// render bg and run g
	LCDC1 := [144]bool{}
	for y := 0; y < iterY; y++ {
		scx, scy, ok := g.execScanline()
		if !ok {
			break
		}
		scrollX, scrollY = scx, scy
		scrollPixelX = scrollX % 8

		if y < height {
			LCDC1[y] = util.Bit(g.FetchMemory8(LCDCIO), 1)
		}

		// render background(or window)
		WY, WX := uint(g.FetchMemory8(WYIO)), uint(g.FetchMemory8(WXIO))-7
		if !skipRender {
			for x := 0; x < iterX; x += 8 {
				blockX, blockY := x/8, y/8

				var tileX, tileY uint
				var isWin bool
				var entryX int

				lineIdx := y % 8 // タイルの何行目を描画するか
				entryY := gpu.EntryY{}
				if util.Bit(LCDC, 5) && (WY <= uint(y)) && (WX <= uint(x)) {
					tileX, tileY = ((uint(x)-WX)/8)%32, ((uint(y)-WY)/8)%32
					isWin = true

					entryX = blockX * 8
					entryY.Block = blockY * 8
					entryY.Offset = y % 8
				} else {
					tileX, tileY = (scrollX+uint(x))/8%32, (scrollY+uint(y))/8%32
					isWin = false

					entryX = blockX*8 - int(scrollPixelX)
					entryY.Block = blockY * 8
					entryY.Offset = y % 8
					lineIdx = (int(scrollY) + y) % 8
				}

				if util.Bit(LCDC, 7) {
					g.GPU.SetBGLine(entryX, entryY, tileX, tileY, isWin, g.Cartridge.IsCGB, lineIdx)
				}
			}
		}
	}

	// save bgmap and tiledata on debug mode
	if g.debug.on {
		if !skipRender {
			bg := g.GPU.Display(false)
			g.GPU.Debug.SetBGMap(bg)
		}
		if frames%4 == 0 {
			go func() {
				g.GPU.UpdateTileData(g.Cartridge.IsCGB)
			}()
		}
	}

	if !skipRender {
		g.renderSprite(&LCDC1)   // render sprite
		g.GPU.SetBGPriorPixels() // render bg has higher priority
	}

	g.execVBlank()
	if g.debug.on {
		select {
		case <-second:
			fps = frames
			frames = 0
		default:
		}
	}
	return nil
}

func (g *GBC) Draw(screen *ebiten.Image) {
	g.renderScreen(screen)
}

func (g *GBC) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if g.debug.on {
		return 1270, 740
	}
	if g.Config.Display.HQ2x {
		return 160 * 2, 144 * 2
	}
	return 160, 144
}
