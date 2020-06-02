package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gbc/pkg/emulator"

	"github.com/hajimehoshi/ebiten"
)

var version string

func main() {
	os.Exit(Run())
}

// Run - エミュレータを実行する
func Run() int {
	var (
		showVersion  = flag.Bool("v", false, "show version")
		debug        = flag.Bool("debug", false, "enable debug mode")
		outputScreen = flag.String("test", "", "only CPU works and output screen map file")
	)

	flag.Parse()

	// バージョンオプションが指定されたときはバージョンを表示して終了する
	if *showVersion {
		fmt.Println("Worldwide: ", version)
		return 0
	}

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

	cpu.Cartridge.ParseCartridge(romData)
	cpu.TransferROM(romData)

	os.Chdir(cur)
	cpu.Init(romDir, *debug)
	defer func() {
		os.Chdir(cur)
		cpu.Exit()
	}()

	if *outputScreen != "" {
		sec := 60
		cpu.Sound.Off()
		cpu.DebugExec(30*sec, *outputScreen)
		return 0
	}

	ebiten.SetRunnableInBackground(true)
	if *debug {
		width, height := ebiten.MonitorSize()
		width, height = width*9/10, height*9/10
		if width >= 1280 && height >= 740 {
			width, height = 1280, 740
		}
		cpu.SetWindow(width, height)
		if err := ebiten.Run(cpu.Render, width, height, 1, "Worldwide(debug)"); err != nil {
			return 1
		}
	} else if cpu.Config.Display.HQ2x {
		if err := ebiten.Run(cpu.Render, 160*2, 144*2, 1, "Worldwide"); err != nil {
			return 1
		}
	} else {
		if err := ebiten.Run(cpu.Render, 160, 144, float64(cpu.Expand), "Worldwide"); err != nil {
			return 1
		}
	}
	return 0
}

func selectROM(p string) (string, error) {
	if p == "" {
		return p, fmt.Errorf("please input ROM file path")
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
