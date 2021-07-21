package emulator

import (
	"os"
	"path/filepath"

	"github.com/pokemium/worldwide/pkg/gbc/cart"
)

// GameBoy save data is SRAM core dump
func (e *Emulator) writeSav() {
	savname := filepath.Join(e.RomDir, e.GBC.Cartridge.Title+".sav")

	savfile, err := os.Create(savname)
	if err != nil {
		return
	}
	defer savfile.Close()

	var buffer []byte
	switch e.GBC.Cartridge.RAMSize {
	case cart.RAM_UNUSED:
		buffer = make([]byte, 0x800)
		for index := 0; index < 0x800; index++ {
			buffer[index] = e.GBC.RAM.Buffer[0][index]
		}
	case cart.RAM_8KB:
		buffer = make([]byte, 0x2000*1)
		for index := 0; index < 0x2000; index++ {
			buffer[index] = e.GBC.RAM.Buffer[0][index]
		}
	case cart.RAM_32KB:
		buffer = make([]byte, 0x2000*4)
		for i := 0; i < 4; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				buffer[index] = e.GBC.RAM.Buffer[i][j]
			}
		}
	case cart.RAM_64KB:
		buffer = make([]byte, 0x2000*8)
		for i := 0; i < 8; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				buffer[index] = e.GBC.RAM.Buffer[i][j]
			}
		}
	}

	if e.GBC.Cartridge.HasRTC() {
		rtcData := e.GBC.RTC.Dump()
		for i := 0; i < 48; i++ {
			buffer = append(buffer, rtcData[i])
		}
	}

	_, err = savfile.Write(buffer)
	if err != nil {
		return
	}
}

func (e *Emulator) loadSav() {
	savname := filepath.Join(e.RomDir, e.GBC.Cartridge.Title+".sav")

	savdata, err := os.ReadFile(savname)
	if err != nil {
		return
	}

	switch e.GBC.Cartridge.RAMSize {
	case cart.RAM_UNUSED:
		for index := 0; index < 0x800; index++ {
			e.GBC.RAM.Buffer[0][index] = savdata[index]
		}
	case cart.RAM_8KB:
		for index := 0; index < 0x2000; index++ {
			e.GBC.RAM.Buffer[0][index] = savdata[index]
		}
	case cart.RAM_32KB:
		for i := 0; i < 4; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				e.GBC.RAM.Buffer[i][j] = savdata[index]
			}
		}
	case cart.RAM_64KB:
		for i := 0; i < 8; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				e.GBC.RAM.Buffer[i][j] = savdata[index]
			}
		}
	}

	if e.GBC.Cartridge.HasRTC() {
		start := (len(savdata) / 0x1000) * 0x1000
		rtcData := savdata[start : start+48]
		e.GBC.RTC.Sync(rtcData)
	}
}
