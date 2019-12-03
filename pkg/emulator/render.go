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
	width  = 160
	height = 144
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
		title = "GameBoy"
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
		second = time.Tick(time.Second)
	)

	var boost float64
	var salt float64
	if cpu.Cartridge.IsCGB {
		salt = 80
	} else {
		salt = 0
	}

	for !win.Closed() {
		if cpu.isBoosted {
			boost = 2
		} else {
			boost = 1
		}

		scrollX := uint(cpu.FetchMemory8(0xff43))
		scrollY := uint(cpu.FetchMemory8(0xff42))
		scrollTileX := scrollX / 8
		scrollPixelX := scrollX % 8
		scrollTileY := scrollY / 8
		scrollPixelY := scrollY % 8

		iterX := width / 8
		iterY := height / 8
		if scrollPixelX > 0 {
			iterX++
		}
		if scrollPixelY > 0 {
			iterY++
		}

		// 背景描画 + CPU稼働
		for y := 0; y < iterY; y++ {
			if y < height/8 {
				// CPU稼働
				wait.Add(1)
				go func() {
					for i := 0; i < 8; i++ {
						cpu.cycleLine = 0
						// HBlank mode0
						// OAM mode2
						cpu.setOAMRAMMode()
						for cpu.cycleLine < 17.25 {
							cpu.exec()
						}
						// LCD Driver mode3
						cpu.setLCDMode()
						for cpu.cycleLine < 40.25 {
							cpu.exec()
						}
						// HBlank mode0
						cpu.setHBlankMode()
						for cpu.cycleLine < 100.25*boost+salt {
							cpu.exec()
						}
						cpu.incrementLY()
					}
					wait.Done()
				}()
			}

			if !cpu.saving || frames&0x01 == 0 {
				// 背景(ウィンドウ)描画
				wait.Add(iterX)
				for x := 0; x < iterX; x++ {
					go func(x int) {
						var tileX, tileY uint
						var useWindow bool
						var entryX, entryY int

						LCDC := cpu.FetchMemory8(LCDCIO)
						WY := uint(cpu.FetchMemory8(WYIO))
						WX := uint(cpu.FetchMemory8(WXIO)) - 7
						if (LCDC>>5)%2 == 1 && (WY <= uint(y*8)) && (WX <= uint(x*8)) {
							tileX = (uint(x*8) - WX) / 8 % 32
							tileY = ((uint(y*8) - WY) / 8) % 32
							useWindow = true

							entryX = x * 8
							entryY = y * 8
						} else {
							tileX = (scrollTileX + uint(x)) % 32
							tileY = (scrollTileY + uint(y)) % 32
							useWindow = false

							entryX = x*8 - int(scrollPixelX)
							entryY = y*8 - int(scrollPixelY)
						}

						if LCDC>>7%2 == 1 {
							cpu.GPU.SetBGTile(entryX, entryY, tileX, tileY, useWindow, cpu.Cartridge.IsCGB)
						}
						wait.Done()
					}(x)
				}
			}
			wait.Wait()
		}

		if !cpu.saving || frames&0x01 == 0 {
			// スプライト描画
			var spriteWait sync.WaitGroup
			spriteWait.Add(8)
			for i := 0; i < 8; i++ {
				for j := i * 5; j < (i+1)*5; j++ {
					LCDC := cpu.FetchMemory8(LCDCIO)
					Y := int(cpu.FetchMemory8(0xfe00 + 4*uint16(j)))
					if LCDC>>1%2 == 1 && Y != 0 && Y < 160 {
						Y -= 16
						X := int(cpu.FetchMemory8(0xfe00+4*uint16(j)+1)) - 8
						tileIndex := uint(cpu.FetchMemory8(0xfe00 + 4*uint16(j) + 2))
						attr := cpu.FetchMemory8(0xfe00 + 4*uint16(j) + 3)
						cpu.GPU.SetSPRTile(int(X), Y, tileIndex, attr, cpu.Cartridge.IsCGB)
					}
				}
				spriteWait.Done()
			}
			spriteWait.Wait()

			// 背景優先のpixelを描画していく
			cpu.GPU.SetBGPriorPixels()
		}

		// VBlank
		for {
			cpu.cycleLine = 0

			if !cpu.saving || frames&0x01 == 0 {
				for cpu.cycleLine < 114*boost {
					cpu.exec()
				}
			}
			cpu.incrementLY()
			LY := cpu.FetchMemory8(LYIO)
			if LY == 0 {
				break
			}
		}

		if !cpu.saving || frames&0x01 == 0 {
			pic := cpu.GPU.Display
			matrix := pixel.IM.Moved(win.Bounds().Center())
			matrix = matrix.ScaledXY(win.Bounds().Center(), pixel.V(float64(cpu.expand), float64(cpu.expand)))
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, matrix)
		}

		win.Update()

		frames++
		select {
		case <-second:
			// fps := fmt.Sprintf("%s | FPS: %d", cfg.Title, frames)
			// win.SetTitle(fps)
			frames = 0
		default:
		}

		if frames&0x01 == 0 {
			cpu.handleJoypad(win)
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
