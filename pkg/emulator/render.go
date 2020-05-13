package emulator

import (
	"bytes"
	"fmt"
	"gbc/pkg/joypad"
	"image"
	"image/color"
	"image/png"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	width        = 160
	height       = 144
	cyclePerLine = 114
)

var (
	wait       sync.WaitGroup
	lineMutex  sync.Mutex
	frames     = 0
	second     = time.Tick(time.Second)
	skipRender bool
	fps        = 0
	bgMap      *ebiten.Image
)

// Render レンダリングを行う
func (cpu *CPU) Render(screen *ebiten.Image) error {

	if frames == 0 {
		setIcon()
	}

	skipRender = (cpu.Config.Display.FPS30) && (frames%2 == 1)

	LCDC := cpu.FetchMemory8(LCDCIO)
	scrollX, scrollY := cpu.GPU.ReadScroll()
	scrollTileX := scrollX / 8
	scrollPixelX := scrollX % 8
	scrollTileY := scrollY / 8
	scrollPixelY := scrollY % 8

	iterX := width
	iterY := height
	if scrollPixelX > 0 {
		iterX += 8
	}
	if scrollPixelY > 0 {
		iterY += 8
	}

	// 背景描画 + CPU稼働
	for y := 0; y < iterY; y++ {

		scrollX, scrollY = cpu.GPU.ReadScroll()
		scrollTileX, scrollPixelX = scrollX/8, scrollX%8
		scrollTileY, scrollPixelY = scrollY/8, scrollY%8

		if y < height {

			// OAM mode2
			cpu.cycleLine = 0
			cpu.setOAMRAMMode()
			for cpu.cycleLine <= 20*cpu.boost {
				cpu.exec()
			}

			// LCD Driver mode3
			cpu.cycleLine = 0
			cpu.setLCDMode()
			for cpu.cycleLine <= 42*cpu.boost {
				cpu.exec()
			}

			// HBlank mode0
			cpu.cycleLine = 0
			cpu.setHBlankMode()
			for cpu.cycleLine <= (cyclePerLine-(20+42))*cpu.boost {
				cpu.exec()
			}
			cpu.incrementLY()
		}

		LCDC = cpu.FetchMemory8(LCDCIO)
		WY := uint(cpu.FetchMemory8(WYIO))
		WX := uint(cpu.FetchMemory8(WXIO)) - 7

		if !skipRender {
			// 背景(ウィンドウ)描画
			for x := 0; x < iterX; x += 8 {
				blockX := x / 8
				blockY := y / 8

				var tileX, tileY uint
				var useWindow bool
				var entryX, entryY int

				if (LCDC>>5)%2 == 1 && (WY <= uint(y)) && (WX <= uint(x)) {
					tileX = ((uint(x) - WX) / 8) % 32
					tileY = ((uint(y) - WY) / 8) % 32
					useWindow = true

					entryX = blockX * 8
					entryY = blockY * 8
				} else {
					tileX = (scrollTileX + uint(x/8)) % 32
					tileY = (scrollTileY + uint(y/8)) % 32
					useWindow = false

					entryX = blockX*8 - int(scrollPixelX)
					entryY = blockY*8 - int(scrollPixelY)
				}

				if LCDC>>7%2 == 1 {
					if !cpu.GPU.SetBGLine(entryX, entryY, tileX, tileY, useWindow, cpu.Cartridge.IsCGB, y%8) {
						break
					}
				}
			}
		}
	}

	// デバッグモードのときはBGマップを保存
	if cpu.debug {
		bg := cpu.GPU.GetDisplay(false)
		bgMap, _ = ebiten.NewImageFromImage(bg, ebiten.FilterDefault)
	}

	if !skipRender {
		// スプライト描画
		for i := 0; i < 40; i++ {
			Y := int(cpu.FetchMemory8(0xfe00 + 4*uint16(i)))
			if LCDC>>1%2 == 1 && Y != 0 && Y < 160 {
				Y -= 16
				X := int(cpu.FetchMemory8(0xfe00+4*uint16(i)+1)) - 8
				tileIndex := uint(cpu.FetchMemory8(0xfe00 + 4*uint16(i) + 2))
				attr := cpu.FetchMemory8(0xfe00 + 4*uint16(i) + 3)
				cpu.GPU.SetSPRTile(int(X), Y, tileIndex, attr, cpu.Cartridge.IsCGB)
			}
		}

		// 背景優先のpixelを描画していく
		cpu.GPU.SetBGPriorPixels()
	}

	// VBlank
	wait.Add(1)
	go func() {
		for {
			cpu.cycleLine = 0

			for cpu.cycleLine < cyclePerLine*cpu.boost {
				cpu.exec()
			}
			cpu.incrementLY()
			LY := cpu.FetchMemory8(LYIO)
			if LY == 0 {
				break
			}
		}
		wait.Done()
	}()

	display := cpu.GPU.GetDisplay(cpu.Config.Display.HQ2x)
	if cpu.debug {
		screen.Fill(color.RGBA{35, 27, 167, 255})
		{
			// debug screen
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(1.5, 1.5)
			op.GeoM.Translate(float64(10), float64(25))
			screen.DrawImage(display, op)
		}

		// debug FPS
		title := fmt.Sprintf("GameBoy FPS: %d", fps)
		ebitenutil.DebugPrintAt(screen, title, 10, 5)

		// debug register
		ebitenutil.DebugPrintAt(screen, cpu.debugRegister(), 270, 5)

		{
			// debug BG
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(10), float64(270))
			screen.DrawImage(bgMap, op)
			ebitenutil.DebugPrintAt(screen, "BG map", 10, 250)
		}

	} else {
		if !skipRender && cpu.Config.Display.HQ2x {
			display = cpu.GPU.HQ2x()
		}
		screen.DrawImage(display, nil)
	}

	frames++

	if frames%3 == 0 {
		pad := cpu.Config.Joypad
		result := cpu.joypad.Input(pad.A, pad.B, pad.Start, pad.Select, pad.Threshold)
		if result != 0 {
			switch result {
			case joypad.Pressed:
				// Joypad Interrupt
				if cpu.Reg.IME && cpu.getJoypadEnable() {
					cpu.setJoypadFlag()
				}
			case joypad.Save:
				cpu.Sound.Off()
				cpu.dumpData()
				cpu.Sound.On()
			case joypad.Load:
				cpu.Sound.Off()
				cpu.loadData()
				cpu.Sound.On()
			case joypad.Expand:
				if !cpu.Config.Display.HQ2x {
					cpu.Expand *= 2
					time.Sleep(time.Millisecond * 400)
					ebiten.SetScreenScale(float64(cpu.Expand))
				}
			case joypad.Collapse:
				if !cpu.Config.Display.HQ2x && cpu.Expand >= 2 {
					cpu.Expand /= 2
					time.Sleep(time.Millisecond * 400)
					ebiten.SetScreenScale(float64(cpu.Expand))
				}
			}
		}
	}

	if cpu.debug {
		select {
		case <-second:
			fps = frames
			frames = 0
		default:
		}
	}

	wait.Wait()
	return nil
}

func setIcon() {
	buf := bytes.NewBuffer(icon)
	img, _ := png.Decode(buf)
	ebiten.SetWindowIcon([]image.Image{img})
}
