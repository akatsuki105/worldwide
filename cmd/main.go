package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gbc/pkg/emulator"
	"github.com/akatsuki-py/tfd"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sqweek/dialog"
)

func main() {

	cpu := &emulator.CPU{}
	romPath := selectROM()
	romData := readROM(romPath)

	cpu.Cartridge.ParseCartridge(&romData)
	cpu.TransferROM(&romData)

	cpu.Init()

	// go cpu.Debug()

	pixelgl.Run(cpu.Render)
}

func selectROM() string {
	var filepath string
	flag.Parse()
	filepath = flag.Arg(0)
	if filepath == "" {
		switch runtime.GOOS {
		case "windows":
			cd, _ := os.Getwd()
			tmp, err := dialog.File().Filter("GameBoy ROM File", "gb*").Load()
			if err != nil {
				os.Exit(0)
			}
			filepath = tmp
			os.Chdir(cd)
		case "linux":
			tmp, err := tfd.CreateSelectDialog([]string{"gb", "gbc"}, false)
			if err != nil {
				os.Exit(0)
			}
			filepath = tmp
		default:
			os.Exit(0)
		}
	}
	return filepath
}

func readROM(path string) []byte {
	if path == "" {
		dialog.Message("%s", "please select gb or gbc file path").Title("Error").Error()
		os.Exit(0)
	}
	if filepath.Ext(path) != ".gb" && filepath.Ext(path) != ".gbc" {
		dialog.Message("%s", "please select .gb or .gbc file").Title("Error").Error()
		os.Exit(0)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}
