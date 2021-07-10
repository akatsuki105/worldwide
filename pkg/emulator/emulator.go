package emulator

import (
	"fmt"
	"gbc/pkg/emulator/debug"
	"gbc/pkg/gbc"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
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

func New(romData []byte, j [8](func() bool), romDir string, debugEnable bool) *Emulator {
	g := gbc.New(romData, j)
	return &Emulator{
		GBC:      g,
		Rom:      romDir,
		debugger: debug.New(debugEnable, g),
	}
}

func (e *Emulator) Update() error {
	defer e.GBC.PanicHandler("update", true)
	err := e.GBC.Update()

	select {
	case <-second.C:
		oldFrame := e.frame
		e.frame = e.GBC.Frame()
		fps := e.frame - oldFrame
		ebiten.SetWindowTitle(fmt.Sprintf("%dfps", fps))
	default:
	}

	return err
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	defer e.GBC.PanicHandler("draw", true)
	screen.ReplacePixels(e.GBC.Draw())
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 160, 144
}
