package emulator

import (
	"gbc/pkg/gpu"
	"gbc/pkg/util"

	"github.com/hajimehoshi/ebiten"
)

func (cpu *CPU) Update(screen *ebiten.Image) error {

	if frames == 0 {
		setIcon()
		cpu.debug.monitor.CPU.Reset()
	}

	cpu.renderScreen(screen)

	if frames%3 == 0 {
		cpu.handleJoypad()
	}

	frames++
	cpu.debug.monitor.CPU.Reset()

	p := &cpu.debug.pause
	b := &cpu.debug.Break
	if p.Delay() {
		p.DecrementDelay()
	}
	if p.On() || b.On() {
		return nil
	}

	skipRender = (cpu.Config.Display.FPS30) && (frames%2 == 1)

	LCDC := cpu.FetchMemory8(LCDCIO)
	scrollX, scrollY := cpu.GPU.GetScroll()
	scrollPixelX := scrollX % 8

	iterX := width
	iterY := height
	if scrollPixelX > 0 {
		iterX += 8
	}

	// 背景描画 + CPU稼働
	LCDC1 := [144]bool{}
	for y := 0; y < iterY; y++ {

		// CPU works
		scx, scy, ok := cpu.execScanline()
		if !ok {
			break
		}
		scrollX, scrollY = scx, scy

		scrollPixelX = scrollX % 8

		LCDC = cpu.FetchMemory8(LCDCIO)
		if y < height {
			LCDC1[y] = util.Bit(LCDC, 1) == 1
		}

		WY := uint(cpu.FetchMemory8(WYIO))
		WX := uint(cpu.FetchMemory8(WXIO)) - 7

		// 背景(ウィンドウ)描画
		if !skipRender {
			wait.Add(iterX / 8)
			for x := 0; x < iterX; x += 8 {
				go func(x int) {
					blockX := x / 8
					blockY := y / 8

					var tileX, tileY uint
					var useWindow bool
					var entryX int

					lineNumber := y % 8 // タイルの何行目を描画するか
					entryY := gpu.EntryY{}
					if util.Bit(LCDC, 5) == 1 && (WY <= uint(y)) && (WX <= uint(x)) {
						tileX = ((uint(x) - WX) / 8) % 32
						tileY = ((uint(y) - WY) / 8) % 32
						useWindow = true

						entryX = blockX * 8
						entryY.Block = blockY * 8
						entryY.Offset = y % 8
					} else {
						tileX = (scrollX + uint(x)) / 8 % 32
						tileY = (scrollY + uint(y)) / 8 % 32
						useWindow = false

						entryX = blockX*8 - int(scrollPixelX)
						entryY.Block = blockY * 8
						entryY.Offset = y % 8
						lineNumber = (int(scrollY) + y) % 8
					}

					if util.Bit(LCDC, 7) == 1 {
						cpu.GPU.SetBGLine(entryX, entryY, tileX, tileY, useWindow, cpu.Cartridge.IsCGB, lineNumber)
					}
					wait.Done()
				}(x)
			}
			wait.Wait()
		}
	}

	// デバッグモードのときはBGマップとタイルデータを保存
	if cpu.debug.on {
		if !skipRender {
			bg := cpu.GPU.GetDisplay(false)
			cpu.GPU.Debug.SetBGMap(bg)
		}

		if frames%4 == 0 {
			go func() {
				cpu.GPU.UpdateTileData(cpu.Cartridge.IsCGB)
			}()
		}
	}

	if !skipRender {
		// スプライト描画
		cpu.renderSprite(&LCDC1)

		// 背景優先のpixelを描画していく
		cpu.GPU.SetBGPriorPixels()
	}

	// VBlank
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

func (cpu *CPU) Draw(screen *ebiten.Image) error {
	return nil
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
