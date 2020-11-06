package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gbc/pkg/emulator"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

var version string

const (
	ExitCodeOK int = iota
	ExitCodeError
)

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
		fmt.Println("Worldwide:", getVersion())
		return ExitCodeOK
	}

	romPath := flag.Arg(0)
	cur, _ := os.Getwd()

	cpu := &emulator.CPU{}

	romDir := filepath.Dir(romPath)

	// ROMファイルを読み込む
	romData, err := readROM(romPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ROM Error: %s\n", err)
		return ExitCodeError
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
		return ExitCodeOK
	}

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Worldwide")

	if *debug {
		ebiten.SetWindowSize(1270, 740)
	} else {
		ebiten.SetWindowSize(160*2, 144*2)
	}

	if err := ebiten.RunGame(cpu); err != nil {
		return ExitCodeError
	}
	return ExitCodeOK
}

func getVersion() string {
	if version == "" {
		return "develop"
	}
	return version
}

func readROM(path string) ([]byte, error) {
	if path == "" {
		return []byte{}, errors.New("please select .gb or .gbc file path")
	}
	if filepath.Ext(path) != ".gb" && filepath.Ext(path) != ".gbc" {
		return []byte{}, errors.New("please select .gb or .gbc file")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, errors.New("fail to read file")
	}
	return bytes, nil
}
