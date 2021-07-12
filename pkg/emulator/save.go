package emulator

import (
	"fmt"
	"io/ioutil"
	"os"
)

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
			savdata[index] = g.RAM.Buffer[0][index]
		}
	case 2:
		savdata = make([]byte, 0x2000*1)
		for index := 0; index < 0x2000; index++ {
			savdata[index] = g.RAM.Buffer[0][index]
		}
	case 3:
		savdata = make([]byte, 0x2000*4)
		for i := 0; i < 4; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				savdata[index] = g.RAM.Buffer[i][j]
			}
		}
	case 5:
		savdata = make([]byte, 0x2000*8)
		for i := 0; i < 8; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				savdata[index] = g.RAM.Buffer[i][j]
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
			g.RAM.Buffer[0][index] = savdata[index]
		}
	case 2:
		for index := 0; index < 0x2000; index++ {
			g.RAM.Buffer[0][index] = savdata[index]
		}
	case 3:
		for i := 0; i < 4; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				g.RAM.Buffer[i][j] = savdata[index]
			}
		}
	case 5:
		for i := 0; i < 8; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				g.RAM.Buffer[i][j] = savdata[index]
			}
		}
	}

	if len(savdata) >= 0x1000 && len(savdata)%0x1000 == 48 {
		start := (len(savdata) / 0x1000) * 0x1000
		rtcData := savdata[start : start+48]
		g.RTC.Sync(rtcData)
	}
}
