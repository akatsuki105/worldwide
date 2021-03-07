package emulator

import (
	"bytes"
	"fmt"
	"gbc/pkg/debug"
	"gbc/pkg/joypad"
	"image"
	"image/color"
	"image/png"
	"sync"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func setIcon() {
	buf := bytes.NewBuffer(icon)
	img, _ := png.Decode(buf)
	ebiten.SetWindowIcon([]image.Image{img})
}

func (cpu *CPU) renderScreen(screen *ebiten.Image) {
	display := cpu.GPU.GetDisplay(cpu.Config.Display.HQ2x)
	if cpu.debug.on {
		dScreen := ebiten.NewImage(int(debugWidth), int(debugHeight))
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
		screen.DrawImage(dScreen, op)
		return
	}

	if !skipRender && cpu.Config.Display.HQ2x {
		display = cpu.GPU.HQ2x()
	}
	screen.DrawImage(display, nil)
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
				} else {
					p.SetOn(30)
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
