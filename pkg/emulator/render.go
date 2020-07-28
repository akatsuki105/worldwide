package emulator

import (
	"bytes"
	"fmt"
	"gbc/pkg/debug"
	"gbc/pkg/gpu"
	"gbc/pkg/joypad"
	"gbc/pkg/util"
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

const (
	debugWidth  = 1270.
	debugHeight = 740.
)

var (
	wait       sync.WaitGroup
	frames     = 0
	second     = time.Tick(time.Second)
	skipRender bool
	fps        = 0
)

// Render レンダリングを行う
func (cpu *CPU) Render(screen *ebiten.Image) error {

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

func setIcon() {
	buf := bytes.NewBuffer(icon)
	img, _ := png.Decode(buf)
	ebiten.SetWindowIcon([]image.Image{img})
}

func (cpu *CPU) renderScreen(screen *ebiten.Image) {
	display := cpu.GPU.GetDisplay(cpu.Config.Display.HQ2x)
	if cpu.debug.on {
		dScreen, _ := ebiten.NewImage(int(debugWidth), int(debugHeight), ebiten.FilterDefault)
		dScreen.Fill(color.RGBA{35, 27, 167, 255})
		{
			// debug screen
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(10), float64(25))
			dScreen.DrawImage(display, op)
		}

		// debug FPS
		title := fmt.Sprintf("GameBoy FPS: %d", fps)
		ebitenutil.DebugPrintAt(dScreen, title, 10, 5)

		// debug register
		ebitenutil.DebugPrintAt(dScreen, cpu.debugRegister(), 340, 5)
		ebitenutil.DebugPrintAt(dScreen, cpu.debugIOMap(), 490, 5)

		// debug Cartridge
		ebitenutil.DebugPrintAt(dScreen, cpu.Cartridge.Debug.String(), 680, 5)

		cpuUsageX := 340
		// debug history (optional)
		if cpu.debug.history.Flag() {
			ebitenutil.DebugPrintAt(dScreen, cpu.debug.history.History(), 340, 120)
			cpuUsageX = 540
		}
		// debug CPU Usage
		ebitenutil.DebugPrintAt(dScreen, "CPU", cpuUsageX, 120)
		cpu.debug.monitor.CPU.DrawUsage(dScreen, cpuUsageX+2, 140, cpu.isBoost())

		bgMap := cpu.GPU.Debug.BGMap()
		if bgMap != nil {
			// debug BG
			ebitenutil.DebugPrintAt(dScreen, "BG map", 10, 320)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(10), float64(340))
			dScreen.DrawImage(bgMap, op)
		}

		{
			// debug tiles
			ebitenutil.DebugPrintAt(dScreen, "Tiles", 200, 320)
			tile := cpu.GPU.GetTileData()
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(200), float64(340))
			dScreen.DrawImage(tile, op)
		}

		// debug OAM
		if cpu.GPU.OAM != nil {
			cpu.debugPrintOAM(dScreen)
		}

		op := &ebiten.DrawImageOptions{}
		windowX, windowY := cpu.debug.Window.Size()
		op.GeoM.Scale(windowX/debugWidth, windowY/debugHeight)
		screen.DrawImage(dScreen, op)
	} else {
		if !skipRender && cpu.Config.Display.HQ2x {
			display = cpu.GPU.HQ2x()
		}
		screen.DrawImage(display, nil)
	}
}

func (cpu *CPU) handleJoypad() {
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
		case joypad.Pause:
			p := &cpu.debug.pause
			b := &cpu.debug.Break

			if !cpu.debug.on {
				return
			}

			if b.On() {
				b.SetFlag(debug.BreakDelay)
				p.SetOff(30)
				return
			}

			if !p.Delay() {
				if p.On() {
					p.SetOff(30)
					cpu.Sound.On()
				} else {
					p.SetOn(30)
					cpu.Sound.Off()
				}
			}
		}
	}
}

func (cpu *CPU) renderSprite(LCDC1 *[144]bool) {
	if cpu.debug.on {
		cpu.GPU.FillOAM()
	}

	for i := 0; i < 40; i++ {
		Y := int(cpu.FetchMemory8(0xfe00 + 4*uint16(i)))
		if Y != 0 && Y < 160 {
			Y -= 16
			X := int(cpu.FetchMemory8(0xfe00+4*uint16(i)+1)) - 8
			tileIndex := uint(cpu.FetchMemory8(0xfe00 + 4*uint16(i) + 2))
			attr := cpu.FetchMemory8(0xfe00 + 4*uint16(i) + 3)
			if Y >= 0 && LCDC1[Y] {
				cpu.GPU.SetSPRTile(i, int(X), Y, tileIndex, attr, cpu.Cartridge.IsCGB)
			}

			if cpu.debug.on {
				cpu.GPU.SetOAMProperty(i, byte(X+8), byte(Y+16), byte(tileIndex), attr)
			}
		}
	}
}
