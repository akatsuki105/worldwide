package emulator

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width    = 160
	height   = 144
	overload = 18
)

var (
	cpuLineWait sync.WaitGroup
	gpuLineWait sync.WaitGroup
	lineMutex   sync.Mutex
)

// Render レンダリングを行う
func (cpu *CPU) Render() {

	var title string
	if cpu.Header.Title != "" {
		title = cpu.Header.Title
	} else {
		title = "GB"
	}

	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	cpu.VRAMModified = true
	cpu.cacheTile()
	cpu.tileCache = cpu.newTileCache
	cpu.tileModified = false
	go func() {
		for range time.Tick(16 * time.Millisecond) {
			cpu.cacheTile()
		}
	}()
	TileBatch := pixel.NewBatch(&pixel.TrianglesData{}, cpu.tileCache)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	for !win.Closed() {
		TileBatch.Clear()

		scrollX := uint(cpu.FetchMemory8(0xff43))
		scrollY := uint(cpu.FetchMemory8(0xff42))
		scrollTileX := scrollX / 8
		scrollPixelX := scrollX % 8
		scrollTileY := scrollY / 8
		scrollPixelY := scrollY % 8

		iterX := width / 8
		iterY := height / 8
		if scrollX+width < 256 {
			iterX++
		}
		if scrollY+height < 256 {
			iterY++
		}

		// 背景描画 + CPU稼働
		for y := 0; y < iterY; y++ {
			// CPU稼働
			cpuLineWait.Add(1)
			go func() {
				for i := 0; i < 8; i++ {
					// OAM
					cpu.setOAMRAMMode()
					for j := 0; j < int(math.Ceil(100/overload)); j++ {
						cpu.exec()
					}
					// LCD Driver
					cpu.setLCDMode()
					for j := 0; j < int(math.Ceil(150/overload)); j++ {
						cpu.exec()
					}
					// HBlank
					cpu.setHBlankMode()
					for j := 0; j < int(math.Ceil(200/overload)); j++ {
						cpu.exec()
					}
					cpu.incrementLY()
				}
				cpuLineWait.Done()
			}()

			// 背景(ウィンドウ)描画
			gpuLineWait.Add(iterX)
			for x := 0; x < iterX; x++ {
				go func(x int) {
					var tileX, tileY uint
					var useWindow bool

					LCDC := cpu.FetchMemory8(LCDCIO)
					WY := uint(cpu.FetchMemory8(WYIO))
					if (LCDC>>5)%2 == 1 && (WY <= uint(y*8)) {
						tileX = (scrollTileX + uint(x)) % 32
						tileY = ((uint(y*8) - WY) / 8) % 32
						WX := uint(cpu.FetchMemory8(WXIO)) - 7
						if uint(x*8) >= WX {
							tileX = (uint(x*8) - WX) / 8
						}
						useWindow = true
					} else {
						tileX = (scrollTileX + uint(x)) % 32
						tileY = (scrollTileY + uint(y)) % 32
						useWindow = false
					}

					if LCDC>>7%2 == 1 {
						rect := cpu.outputBGTile(tileX, tileY, useWindow)
						BGTile := pixel.NewSprite(cpu.tileCache, rect)
						matrix := pixel.IM.Moved(pixel.V(float64(uint(x*8)-scrollPixelX+4), float64(uint(height-y*8)+scrollPixelY-4)))
						lineMutex.Lock()
						BGTile.Draw(TileBatch, matrix)
						lineMutex.Unlock()
					}
					gpuLineWait.Done()
				}(x)
			}
			gpuLineWait.Wait()
			cpuLineWait.Wait()
		}

		// スプライト描画
		var spriteWait sync.WaitGroup
		spriteWait.Add(40)
		for i := 0; i < 40; i++ {
			go func(i int) {
				LCDC := cpu.FetchMemory8(LCDCIO)
				Y := uint8(cpu.FetchMemory8(0xfe00 + 4*uint16(i)))
				if LCDC>>1%2 == 1 && Y != 0 && Y < 160 {
					Y -= 16
					X := uint8(cpu.FetchMemory8(0xfe00+4*uint16(i)+1)) - 8
					tileNum := uint8(cpu.FetchMemory8(0xfe00 + 4*uint16(i) + 2))
					attr := cpu.FetchMemory8(0xfe00 + 4*uint16(i) + 3)

					rect := cpu.outputSPRTile(tileNum, attr)
					SPRTile := pixel.NewSprite(cpu.tileCache, rect)
					matrix := pixel.IM.Moved(pixel.V(float64(X+4), float64(height-Y-4)))
					lineMutex.Lock()
					SPRTile.Draw(TileBatch, matrix)
					lineMutex.Unlock()

					spriteYSize := cpu.fetchSPRYSize()
					if spriteYSize == 16 {
						rect := cpu.outputSPRTile(tileNum+1, attr)
						SPRTile := pixel.NewSprite(cpu.tileCache, rect)
						matrix := pixel.IM.Moved(pixel.V(float64(X+4), float64(height-(Y+8)-4)))
						lineMutex.Lock()
						SPRTile.Draw(TileBatch, matrix)
						lineMutex.Unlock()
					}
				}
				spriteWait.Done()
			}(i)
		}
		spriteWait.Wait()

		// VBlank
		for {
			for j := 0; j < int(math.Ceil(456/overload)); j++ {
				cpu.exec()
			}
			cpu.incrementLY()
			LY := cpu.FetchMemory8(LYIO)
			if LY == 0 {
				break
			}
		}

		TileBatch.Draw(win)

		if cpu.tileModified {
			cpu.tileCache = cpu.newTileCache
			cpu.tileModified = false
			TileBatch = pixel.NewBatch(&pixel.TrianglesData{}, cpu.tileCache)
		}

		win.Update()

		frames++
		select {
		case <-second:
			fmt.Printf("%s | FPS: %d\n", cfg.Title, frames)
			frames = 0
		default:
		}

		go cpu.handleJoypad(win)

		// coredump
		if win.Pressed(pixelgl.KeyD) && win.Pressed(pixelgl.KeyS) {
			cpu.dumpData()
		}
		if win.Pressed(pixelgl.KeyL) {
			cpu.loadData()
		}
	}
}
