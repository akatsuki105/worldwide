package emulator

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pokemium/worldwide/pkg/emulator/audio"
	"github.com/pokemium/worldwide/pkg/emulator/debug"
	"github.com/pokemium/worldwide/pkg/gbc"
)

var (
	second = time.NewTicker(time.Second)
	cache  []byte
)

type Emulator struct {
	GBC      *gbc.GBC
	Rom      string
	debugger *debug.Debugger
	frame    int
	pause    bool
}

func New(romData []byte, j [8](func() bool), romDir string) *Emulator {
	g := gbc.New(romData, j, audio.SetStream)
	audio.Init(&g.Sound.Enable)

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("60fps")
	ebiten.SetWindowSize(160*2, 144*2)

	e := &Emulator{
		GBC:      g,
		Rom:      romDir,
		debugger: debug.New(g),
	}
	e.setupCloseHandler()

	return e
}

func (e *Emulator) Update() error {
	if e.pause {
		return nil
	}

	defer e.GBC.PanicHandler("update", true)
	e.pause = e.GBC.Update(e.debugger.Breakpoints)
	if e.pause {
		return nil
	}

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

	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	if e.pause {
		screen.ReplacePixels(cache)
		return
	}

	defer e.GBC.PanicHandler("draw", true)
	cache = e.GBC.Draw()
	screen.ReplacePixels(cache)
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 160, 144
}

func (e *Emulator) Exit() {
	e.WriteSav()
}

func (e *Emulator) setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		e.Exit()
		os.Exit(0)
	}()
}
