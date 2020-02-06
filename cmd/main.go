package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gbc/pkg/emulator"

	"github.com/faiface/pixel/pixelgl"
	"github.com/sqweek/dialog"
)

func main() {
	os.Exit(Run())
}

// Run - エミュレータを実行する
func Run() int {
	var (
		debug = flag.Bool("debug", false, "enable debug mode")
	)

	flag.Parse()
	fp := flag.Arg(0)
	cur, _ := os.Getwd()

	cpu := &emulator.CPU{}

	// ROMファイルのパスを取得する
	romPath, err := selectROM(fp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	romDir := filepath.Dir(romPath)

	// ROMファイルを読み込む
	romData, err := readROM(romPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	cpu.Cartridge.ParseCartridge(&romData)
	cpu.TransferROM(&romData)

	os.Chdir(cur)
	cpu.Init(romDir, *debug)
	defer func() {
		os.Chdir(cur)
		cpu.Exit()
	}()

	pixelgl.Run(cpu.Render)
	return 0
}

func selectROM(p string) (string, error) {
	if p == "" {
		switch runtime.GOOS {
		case "windows":
			cd, _ := os.Getwd()
			tmp, err := dialog.File().Filter("GameBoy ROM File", "gb*").Load()
			if err != nil {
				return p, fmt.Errorf("failed to read ROM file: %s", err)
			}
			p = tmp
			os.Chdir(cd)
		default:
			return p, fmt.Errorf("ROM file is nil")
		}
	}
	return p, nil
}

func readROM(path string) ([]byte, error) {
	if path == "" {
		return []byte{}, errors.New("please select gb or gbc file path")
	}
	if filepath.Ext(path) != ".gb" && filepath.Ext(path) != ".gbc" {
		return []byte{}, errors.New("please select .gb or .gbc file")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes, nil
}
