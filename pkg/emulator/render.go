package emulator

import (
	"bytes"
	"fmt"
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

type Pause struct {
	on    bool
	delay int
}

var (
	wait        sync.WaitGroup
	lineMutex   sync.Mutex
	frames      = 0
	second      = time.Tick(time.Second)
	skipRender  bool
	fps         = 0
	bgMap       *ebiten.Image
	OAMProperty = [40][4]byte{}
	pause       Pause
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

	if pause.delay > 0 {
		pause.delay--
	}
	if pause.on {
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
		scrollX, scrollY = cpu.execScanline()
		scrollPixelX = scrollX % 8

		LCDC = cpu.FetchMemory8(LCDCIO)
		if y < height {
			LCDC1[y] = util.Bit(LCDC, 1) == 1
		}

		WY := uint(cpu.FetchMemory8(WYIO))
		WX := uint(cpu.FetchMemory8(WXIO)) - 7

		// 背景(ウィンドウ)描画
		if !skipRender {
			for x := 0; x < iterX; x += 8 {
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
					if !cpu.GPU.SetBGLine(entryX, entryY, tileX, tileY, useWindow, cpu.Cartridge.IsCGB, lineNumber) {
						break
					}
				}
			}
		}
	}

	// デバッグモードのときはBGマップとタイルデータを保存
	if cpu.debug {
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
		cpu.setSprite(&LCDC1)

		// 背景優先のpixelを描画していく
		cpu.GPU.SetBGPriorPixels()
	}

	// VBlank
	cpu.execVBlank()

	if cpu.debug {
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
	if cpu.debug {
		screen.Fill(color.RGBA{35, 27, 167, 255})
		{
			// debug screen
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(10), float64(25))
			screen.DrawImage(display, op)
		}

		// debug FPS
		title := fmt.Sprintf("GameBoy FPS: %d", fps)
		ebitenutil.DebugPrintAt(screen, title, 10, 5)

		// debug register
		ebitenutil.DebugPrintAt(screen, cpu.debugRegister(), 340, 5)

		if bgMap != nil {
			// debug BG
			ebitenutil.DebugPrintAt(screen, "BG map", 10, 320)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(10), float64(340))
			screen.DrawImage(bgMap, op)
		}

		{
			// debug tiles
			ebitenutil.DebugPrintAt(screen, "Tiles", 200, 320)
			tile := cpu.GPU.GetTileData()
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(200), float64(340))
			screen.DrawImage(tile, op)
		}

		if cpu.GPU.OAM != nil {
			// debug OAM
			ebitenutil.DebugPrintAt(screen, "OAM (Y, X, tile, attr)", 750, 320)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(4, 4)
			op.GeoM.Translate(float64(750), float64(340))
			screen.DrawImage(cpu.GPU.OAM, op)

			for i := 0; i < 40; i++ {
				Y, X, index, attr := OAMProperty[i][0], OAMProperty[i][1], OAMProperty[i][2], OAMProperty[i][3]

				col := i % 8
				row := i / 8
				property := fmt.Sprintf("%02x\n%02x\n%02x\n%02x", Y, X, index, attr)
				ebitenutil.DebugPrintAt(screen, property, 750+(col*64)+42, 340+(row*80))
			}
		}
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
			if !cpu.Config.Display.HQ2x && !cpu.debug {
				cpu.Expand *= 2
				time.Sleep(time.Millisecond * 400)
				ebiten.SetScreenScale(float64(cpu.Expand))
			}
		case joypad.Collapse:
			if !cpu.Config.Display.HQ2x && cpu.Expand >= 2 && !cpu.debug {
				cpu.Expand /= 2
				time.Sleep(time.Millisecond * 400)
				ebiten.SetScreenScale(float64(cpu.Expand))
			}
		case joypad.Pause:
			if cpu.debug && pause.delay <= 0 {
				if pause.on {
					pause.on = false
					pause.delay = 30
					cpu.Sound.On()
				} else {
					pause.on = true
					pause.delay = 30
					cpu.Sound.Off()
				}
			}
		}
	}
}

func (cpu *CPU) setSprite(LCDC1 *[144]bool) {
	cpu.GPU.OAM, _ = ebiten.NewImage(16*8-1, 20*5-3, ebiten.FilterDefault)
	cpu.GPU.OAM.Fill(color.RGBA{0x8f, 0x8f, 0x8f, 0xff})

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

			if cpu.debug {
				OAMProperty[i] = [4]byte{byte(Y), byte(X), byte(tileIndex), attr}
			}
		}
	}
}
