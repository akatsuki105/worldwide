package emulator

import (
	"bytes"
	"image/png"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width        = 160
	height       = 144
	cyclePerLine = 114
)

var (
	wait      sync.WaitGroup
	lineMutex sync.Mutex
)

// Render レンダリングを行う
func (cpu *CPU) Render() {

	var title string
	if cpu.Cartridge.Title != "" {
		title = cpu.Cartridge.Title
	} else {
		title = "Worldwide"
	}

	icon, err := loadIcon()
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  title,
		Icon:   []pixel.Picture{icon},
		Bounds: pixel.R(0, 0, float64(width*cpu.expand), float64(height*cpu.expand)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		frames = 0
	)

	var boost float64

	win.SetSmooth(cpu.smooth)

	for !win.Closed() {

		if !win.Focused() && !cpu.network {
			cpu.Sound.Off()
			frames++
			win.Update()
			continue
		}
		cpu.Sound.On()

		if cpu.isBoosted {
			boost = 2
		} else {
			boost = 1
		}

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
			scrollTileX = scrollX / 8
			scrollPixelX = scrollX % 8
			scrollTileY = scrollY / 8
			scrollPixelY = scrollY % 8

			if y < height {
				cpu.cycleLine = 0
				// HBlank mode0
				// OAM mode2
				cpu.setOAMRAMMode()
				for cpu.cycleLine < 20 {
					cpu.exec()
				}
				// LCD Driver mode3
				cpu.setLCDMode()
				for cpu.cycleLine < 20+43 {
					cpu.exec()
				}
				// HBlank mode0
				cpu.setHBlankMode()
				for cpu.cycleLine < cyclePerLine*boost {
					cpu.exec()
				}
				cpu.incrementLY()
			}

			LCDC = cpu.FetchMemory8(LCDCIO)
			WY := uint(cpu.FetchMemory8(WYIO))
			WX := uint(cpu.FetchMemory8(WXIO)) - 7

			// 背景(ウィンドウ)描画
			for x := 0; x < iterX; x += 8 {
				blockX := x / 8
				blockY := y / 8

				var tileX, tileY uint
				var useWindow bool
				var entryX, entryY int

				if (LCDC>>5)%2 == 1 && (WY <= uint(y)) && (WX <= uint(x)) {
					tileX = (uint(x) - WX) / 8 % 32
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

		// VBlank
		for {
			cpu.cycleLine = 0

			for cpu.cycleLine < cyclePerLine*boost {
				cpu.exec()
			}
			cpu.incrementLY()
			LY := cpu.FetchMemory8(LYIO)
			if LY == 0 {
				break
			}
		}

		pic := cpu.GPU.GetDisplay()
		matrix := pixel.IM.Moved(win.Bounds().Center())
		matrix = matrix.ScaledXY(win.Bounds().Center(), pixel.V(float64(cpu.expand), float64(cpu.expand)))
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, matrix)

		win.Update()

		frames++

		if frames%3 == 0 {
			result := cpu.joypad.Input(win)
			if result != "" {
				switch result {
				case "pressed":
					// Joypad Interrupt
					if cpu.Reg.IME && cpu.getJoypadEnable() {
						cpu.triggerJoypad()
					}
				case "save":
					cpu.Sound.Off()
					cpu.dumpData()
					cpu.Sound.On()
				case "load":
					cpu.Sound.Off()
					cpu.loadData()
					cpu.Sound.On()
				case "expand":
					cpu.expand *= 2
					time.Sleep(time.Millisecond * 400)
					win.SetBounds(pixel.R(0, 0, float64(width*cpu.expand), float64(height*cpu.expand)))
				case "collapse":
					if cpu.expand >= 2 {
						cpu.expand /= 2
						time.Sleep(time.Millisecond * 400)
						win.SetBounds(pixel.R(0, 0, float64(width*cpu.expand), float64(height*cpu.expand)))
					}
				}
			}
		}
	}
}

func loadIcon() (pixel.Picture, error) {
	buf := bytes.NewBuffer(icon)
	img, err := png.Decode(buf)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}
