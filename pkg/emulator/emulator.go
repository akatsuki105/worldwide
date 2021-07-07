package emulator

import (
	"fmt"
	"gbc/pkg/gbc"
	"io/ioutil"
	"os"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

var (
	second = time.NewTicker(time.Second)
)

type Emulator struct {
	GBC   *gbc.GBC
	Rom   string
	frame int
}

func New(romData []byte, j [8](func() bool), romDir string) *Emulator {
	return &Emulator{
		GBC: gbc.New(romData, j),
		Rom: romDir,
	}
}

// GameBoy save data is SRAM core dump
func (e *Emulator) WriteSav() {
	g := e.GBC
	savname := fmt.Sprintf("%s/%s.sav", e.Rom, g.Cartridge.Title)
	savfile, err := os.Create(savname)
	if err != nil {
		return
	}
	defer savfile.Close()

	var savdata []byte
	switch g.Cartridge.RAMSize {
	case 1:
		savdata = make([]byte, 0x800)
		for index := 0; index < 0x800; index++ {
			savdata[index] = g.RAMBank.Bank[0][index]
		}
	case 2:
		savdata = make([]byte, 0x2000*1)
		for index := 0; index < 0x2000; index++ {
			savdata[index] = g.RAMBank.Bank[0][index]
		}
	case 3:
		savdata = make([]byte, 0x2000*4)
		for i := 0; i < 4; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				savdata[index] = g.RAMBank.Bank[i][j]
			}
		}
	case 5:
		savdata = make([]byte, 0x2000*8)
		for i := 0; i < 8; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				savdata[index] = g.RAMBank.Bank[i][j]
			}
		}
	}

	if g.RTC.Enable {
		rtcData := g.RTC.Dump()
		for i := 0; i < 48; i++ {
			savdata = append(savdata, rtcData[i])
		}
	}

	_, err = savfile.Write(savdata)
	if err != nil {
		return
	}
}

func (e *Emulator) LoadSav() {
	g := e.GBC
	savname := fmt.Sprintf("%s/%s.sav", e.Rom, g.Cartridge.Title)
	savdata, err := ioutil.ReadFile(savname)
	if err != nil {
		return
	}
	switch g.Cartridge.RAMSize {
	case 1:
		for index := 0; index < 0x800; index++ {
			g.RAMBank.Bank[0][index] = savdata[index]
		}
	case 2:
		for index := 0; index < 0x2000; index++ {
			g.RAMBank.Bank[0][index] = savdata[index]
		}
	case 3:
		for i := 0; i < 4; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				g.RAMBank.Bank[i][j] = savdata[index]
			}
		}
	case 5:
		for i := 0; i < 8; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				g.RAMBank.Bank[i][j] = savdata[index]
			}
		}
	}

	if len(savdata) >= 0x1000 && len(savdata)%0x1000 == 48 {
		start := (len(savdata) / 0x1000) * 0x1000
		rtcData := savdata[start : start+48]
		g.RTC.Sync(rtcData)
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
