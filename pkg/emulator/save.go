package emulator

import (
	"fmt"
	"io/ioutil"
	"os"
)

func (cpu *CPU) save() {
	savname := fmt.Sprintf("%s/%s.sav", cpu.romdir, cpu.Cartridge.Title)
	savfile, err := os.Create(savname)
	if err != nil {
		return
	}
	defer savfile.Close()

	var savdata []byte
	switch cpu.Cartridge.RAMSize {
	case 1:
		savdata = make([]byte, 0x800)
		for index := 0; index < 0x800; index++ {
			savdata[index] = cpu.RAMBank[0][index]
		}
	case 2:
		savdata = make([]byte, 0x2000*1)
		for index := 0; index < 0x2000; index++ {
			savdata[index] = cpu.RAMBank[0][index]
		}
	case 3:
		savdata = make([]byte, 0x2000*4)
		for i := 0; i < 4; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				savdata[index] = cpu.RAMBank[i][j]
			}
		}
	case 5:
		savdata = make([]byte, 0x2000*8)
		for i := 0; i < 8; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				savdata[index] = cpu.RAMBank[i][j]
			}
		}
	}

	if cpu.RTC.Working {
		rtcData := cpu.RTC.Dump()
		for i := 0; i < 48; i++ {
			savdata = append(savdata, rtcData[i])
		}
	}

	_, err = savfile.Write(savdata)
	if err != nil {
		return
	}
}

func (cpu *CPU) load() {
	savname := fmt.Sprintf("%s/%s.sav", cpu.romdir, cpu.Cartridge.Title)
	savdata, err := ioutil.ReadFile(savname)
	if err != nil {
		return
	}

	switch cpu.Cartridge.RAMSize {
	case 1:
		for index := 0; index < 0x800; index++ {
			cpu.RAMBank[0][index] = savdata[index]
		}
	case 2:
		for index := 0; index < 0x2000; index++ {
			cpu.RAMBank[0][index] = savdata[index]
		}
	case 3:
		for i := 0; i < 4; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				cpu.RAMBank[i][j] = savdata[index]
			}
		}
	case 5:
		for i := 0; i < 5; i++ {
			for j := 0; j < 0x2000; j++ {
				index := i*0x2000 + j
				cpu.RAMBank[i][j] = savdata[index]
			}
		}
	}

	if len(savdata) >= 0x1000 && len(savdata)%0x1000 == 48 {
		start := (len(savdata) / 0x1000) * 0x1000
		rtcData := savdata[start : start+48]
		cpu.RTC.Sync(rtcData)
	}
}
