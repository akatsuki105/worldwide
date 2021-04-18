package gbc

import (
	"gbc/pkg/gpu"
	"gbc/pkg/util"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

func (cpu *CPU) Update() error {
	if frames == 0 {
		setIcon()
		cpu.debug.monitor.CPU.Reset()
	}
	if frames%3 == 0 {
		cpu.handleJoypad()
	}

	frames++
	cpu.debug.monitor.CPU.Reset()

	p, b := &cpu.debug.pause, &cpu.debug.Break
	if p.Delay() {
		p.DecrementDelay()
	}
	if p.On() || b.On() {
		return nil
	}

	skipRender = (cpu.Config.Display.FPS30) && (frames%2 == 1)

	LCDC := cpu.FetchMemory8(LCDCIO)
	scrollX, scrollY := uint(cpu.GPU.Scroll[0]), uint(cpu.GPU.Scroll[1])
	scrollPixelX := scrollX % 8

	iterX, iterY := width, height
	if scrollPixelX > 0 {
		iterX += 8
	}

	// render bg and run cpu
	LCDC1 := [144]bool{}
	for y := 0; y < iterY; y++ {
		scx, scy, ok := cpu.execScanline()
		if !ok {
			break
		}
		scrollX, scrollY = scx, scy
		scrollPixelX = scrollX % 8

		if y < height {
			LCDC1[y] = util.Bit(cpu.FetchMemory8(LCDCIO), 1)
		}

		// render background(or window)
		WY, WX := uint(cpu.FetchMemory8(WYIO)), uint(cpu.FetchMemory8(WXIO))-7
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
					cpu.GPU.SetBGLine(entryX, entryY, tileX, tileY, isWin, cpu.Cartridge.IsCGB, lineIdx)
				}
			}
		}
	}

	// save bgmap and tiledata on debug mode
	if cpu.debug.on {
		if !skipRender {
			bg := cpu.GPU.Display(false)
			cpu.GPU.Debug.SetBGMap(bg)
		}
		if frames%4 == 0 {
			go func() {
				cpu.GPU.UpdateTileData(cpu.Cartridge.IsCGB)
			}()
		}
	}

	if !skipRender {
		cpu.renderSprite(&LCDC1)   // render sprite
		cpu.GPU.SetBGPriorPixels() // render bg has higher priority
	}

	cpu.execVBlank()
	if cpu.debug.on {
		select {
		case <-second:
			fps = frames
			frames = 0
		default:
		}
	}
	return nil
}

func (cpu *CPU) Draw(screen *ebiten.Image) {
	cpu.renderScreen(screen)
}

func (cpu *CPU) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if cpu.debug.on {
		return 1270, 740
	}
	if cpu.Config.Display.HQ2x {
		return 160 * 2, 144 * 2
	}
	return 160, 144
}
