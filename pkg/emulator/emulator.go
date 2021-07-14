package emulator

import (
	"fmt"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/pokemium/Worldwide/pkg/emulator/audio"
	"github.com/pokemium/Worldwide/pkg/emulator/debug"
	"github.com/pokemium/Worldwide/pkg/gbc"
)

var (
	second = time.NewTicker(time.Second)
)

type Emulator struct {
	GBC      *gbc.GBC
	Rom      string
	debugger *debug.Debugger
	frame    int
}

func New(romData []byte, j [8](func() bool), romDir string, isDebugMode bool) *Emulator {
	audio.Init()
	g := gbc.New(romData, j, audio.SetStream)

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("60fps")
	if isDebugMode {
		ebiten.SetWindowSize(DEBUG_BG_X, DEBUG_BG_Y)
	} else {
		ebiten.SetWindowSize(160*2, 144*2)
	}
	return &Emulator{
		GBC:      g,
		Rom:      romDir,
		debugger: debug.New(isDebugMode, g),
	}
}

func (e *Emulator) Update() error {
	defer e.GBC.PanicHandler("update", true)
	err := e.GBC.Update()
	audio.Play()

	select {
	case <-second.C:
		e.GBC.RTC.IncrementSecond()
		oldFrame := e.frame
		e.frame = e.GBC.Frame()
		fps := e.frame - oldFrame
		ebiten.SetWindowTitle(fmt.Sprintf("%dfps", fps))
	default:
	}

	return err
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	if e.debugger.Enable {
		screen.DrawImage(e.drawDebugScreen(), nil)
		return
	}

	defer e.GBC.PanicHandler("draw", true)
	screen.ReplacePixels(e.GBC.Draw())
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if e.debugger.Enable {
		return DEBUG_BG_X, DEBUG_BG_Y
	}
	return 160, 144
}
