package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gbc/pkg/emulator"
	"gbc/pkg/emulator/joypad"
	"gbc/pkg/gbc"

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

// Run program
func Run() int {
	var (
		showVersion  = flag.Bool("v", false, "show version")
		debug        = flag.Bool("debug", false, "enable debug mode")
		outputScreen = flag.String("test", "", "only CPU works and output screen map file")
	)

	flag.Parse()

	if *showVersion {
		fmt.Println("Worldwide:", getVersion())
		return ExitCodeOK
	}

	romPath := flag.Arg(0)
	cur, _ := os.Getwd()

	romDir := filepath.Dir(romPath)
	romData, err := readROM(romPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ROM Error: %s\n", err)
		return ExitCodeError
	}

	emu := &emulator.Emulator{
		GBC: gbc.New(romData, joypad.Handler),
		Rom: romDir,
	}

	test := *outputScreen != ""
	os.Chdir(cur)
	defer func() {
		os.Chdir(cur)
	}()

	if test {
		sec := 60
		emu.GBC.DebugExec(30*sec, *outputScreen)
		return ExitCodeOK
	}

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("60fps")

	if *debug {
		ebiten.SetWindowSize(1270, 740)
	} else {
		ebiten.SetWindowSize(160*2, 144*2)
	}

	emu.LoadSav()
	if err := ebiten.RunGame(emu); err != nil {
		return ExitCodeError
	}
	return ExitCodeOK
}

func getVersion() string {
	if version == "" {
		return "Develop"
	}
	return version
}

func readROM(path string) ([]byte, error) {
	if path == "" {
		return []byte{}, errors.New("please type .gb or .gbc file path")
	}
	if filepath.Ext(path) != ".gb" && filepath.Ext(path) != ".gbc" {
		return []byte{}, errors.New("please type .gb or .gbc file")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, errors.New("fail to read file")
	}
	return bytes, nil
}
