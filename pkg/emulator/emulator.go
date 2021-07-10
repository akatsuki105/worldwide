package emulator

import (
	"fmt"
	"gbc/pkg/emulator/debug"
	"gbc/pkg/gbc"
	"image/color"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var DEBUG_BG = [3]byte{35, 27, 167}

const (
	DEBUG_BG_X = 1270
	DEBUG_BG_Y = 740
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
	g := gbc.New(romData, j)

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
	if e.debugger.Enable {
		screen.DrawImage(e.debugScreen(), nil)
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

func (e *Emulator) debugScreen() *ebiten.Image {
	screen := ebiten.NewImage(DEBUG_BG_X, DEBUG_BG_Y)
	screen.Fill(color.RGBA{DEBUG_BG[0], DEBUG_BG[1], DEBUG_BG[2], 0xff})

	gbcScreen := ebiten.NewImage(160, 144)
	gbcScreen.ReplacePixels(e.GBC.Draw())
	{
		// debug screen
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(float64(10), float64(25))
		screen.DrawImage(ebiten.NewImageFromImage(gbcScreen), op)
	}

	// debug title
	ebitenutil.DebugPrintAt(screen, "GameBoy screen", 10, 5)

	// debug register
	ebitenutil.DebugPrintAt(screen, e.debugger.Register(), 340, 5)
	ebitenutil.DebugPrintAt(screen, e.debugger.IOMap(), 490, 5)

	// debug cartridge
	ebitenutil.DebugPrintAt(screen, e.debugger.Cartridge(), 680, 5)

	return screen
}
