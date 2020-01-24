package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gbc/pkg/emulator"

	"github.com/faiface/pixel/pixelgl"
	"github.com/sqweek/dialog"
)

func main() {
	flag.Parse()
	fp := flag.Arg(0)

	cur, _ := os.Getwd()

	cpu := &emulator.CPU{}
	romPath := selectROM(fp)
	romDir := filepath.Dir(romPath)
	romData := readROM(romPath)

	cpu.Cartridge.ParseCartridge(&romData)
	cpu.TransferROM(&romData)

	os.Chdir(cur)
	cpu.Init(romDir)
	defer func() {
		os.Chdir(cur)
		cpu.Exit()
	}()

	// go cpu.Debug(2)

	pixelgl.Run(cpu.Render)
}

func selectROM(p string) string {
	if p == "" {
		switch runtime.GOOS {
		case "windows":
			cd, _ := os.Getwd()
			tmp, err := dialog.File().Filter("GameBoy ROM File", "gb*").Load()
			if err != nil {
				os.Exit(0)
			}
			p = tmp
			os.Chdir(cd)
		default:
			os.Exit(0)
		}
	}
	return p
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
