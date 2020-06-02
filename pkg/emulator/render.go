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
	"image/draw"
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
	wait        sync.WaitGroup
	frames      = 0
	second      = time.Tick(time.Second)
	skipRender  bool
	fps         = 0
	bgMap       *ebiten.Image
	OAMProperty = [40][4]byte{}
)

// Render レンダリングを行う
func (cpu *CPU) Render(screen *ebiten.Image) error {

	if frames == 0 {
		setIcon()
	}

	cpu.renderScreen(screen)

	if frames%3 == 0 {
		cpu.handleJoypad()
	}

	frames++

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
			bgMap, _ = ebiten.NewImageFromImage(bg, ebiten.FilterDefault)
		}

		if frames%4 == 0 {
			go func() {
				cpu.GPU.UpdateTiles(cpu.Cartridge.IsCGB)
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
		debugScreen, _ := ebiten.NewImage(int(debugWidth), int(debugHeight), ebiten.FilterDefault)
		debugScreen.Fill(color.RGBA{35, 27, 167, 255})
		{
			// debug screen
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(10), float64(25))
			debugScreen.DrawImage(display, op)
		}

		// debug FPS
		title := fmt.Sprintf("GameBoy FPS: %d", fps)
		ebitenutil.DebugPrintAt(debugScreen, title, 10, 5)

		// debug register
		ebitenutil.DebugPrintAt(debugScreen, cpu.debugRegister(), 340, 5)
		ebitenutil.DebugPrintAt(debugScreen, cpu.debugIOMap(), 490, 5)
		ebitenutil.DebugPrintAt(debugScreen, cpu.debug.history.History(), 340, 120)

		if bgMap != nil {
			// debug BG
			ebitenutil.DebugPrintAt(debugScreen, "BG map", 10, 320)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(10), float64(340))
			debugScreen.DrawImage(bgMap, op)
		}

		{
			// debug tiles
			ebitenutil.DebugPrintAt(debugScreen, "Tiles", 200, 320)
			tile := cpu.GPU.GetTileData()
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(200), float64(340))
			debugScreen.DrawImage(tile, op)
		}

		if cpu.GPU.OAM != nil {
			// debug OAM
			ebitenutil.DebugPrintAt(debugScreen, "OAM (Y, X, tile, attr)", 750, 320)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(4, 4)
			op.GeoM.Translate(float64(750), float64(340))
			OAMScreen, _ := ebiten.NewImageFromImage(cpu.GPU.OAM, ebiten.FilterDefault)
			debugScreen.DrawImage(OAMScreen, op)

			for i := 0; i < 40; i++ {
				Y, X, index, attr := OAMProperty[i][0], OAMProperty[i][1], OAMProperty[i][2], OAMProperty[i][3]

				col := i % 8
				row := i / 8
				property := fmt.Sprintf("%02x\n%02x\n%02x\n%02x", Y, X, index, attr)
				ebitenutil.DebugPrintAt(debugScreen, property, 750+(col*64)+42, 340+(row*80))
			}
		}
		op := &ebiten.DrawImageOptions{}
		monitorX, monitorY := cpu.Monitor()
		op.GeoM.Scale(monitorX/debugWidth, monitorY/debugHeight)
		screen.DrawImage(debugScreen, op)
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
		case joypad.Expand:
			if !cpu.Config.Display.HQ2x && !cpu.debug.on {
				cpu.Expand *= 2
				time.Sleep(time.Millisecond * 400)
				ebiten.SetScreenScale(float64(cpu.Expand))
			}
		case joypad.Collapse:
			if !cpu.Config.Display.HQ2x && cpu.Expand >= 2 && !cpu.debug.on {
				cpu.Expand /= 2
				time.Sleep(time.Millisecond * 400)
				ebiten.SetScreenScale(float64(cpu.Expand))
			}
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
		OAMScreen := cpu.GPU.OAM
		c := color.RGBA{0x8f, 0x8f, 0x8f, 0xff}
		draw.Draw(OAMScreen, OAMScreen.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
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
				OAMProperty[i] = [4]byte{byte(Y), byte(X), byte(tileIndex), attr}
			}
		}
	}
}
