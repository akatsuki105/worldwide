package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pokemium/worldwide/pkg/emulator"

	"github.com/hajimehoshi/ebiten/v2"
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
		showVersion = flag.Bool("v", false, "show version")
		port        = flag.Int("p", -1, "server port (>1023)")
	)

	flag.Parse()

	if *showVersion {
		fmt.Println("worldwide:", getVersion())
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

	emu := emulator.New(romData, romDir)
	if *port > 0 {
		if *port < 1024 {
			fmt.Fprintf(os.Stderr, "Server Error: cannot use well-known port for server")
		} else {
			go emu.RunServer(*port)
		}
	}

	os.Chdir(cur)
	defer func() {
		os.Chdir(cur)
	}()

	if err := ebiten.RunGame(emu); err != nil {
		if err.Error() == "quit" {
			emu.Exit()
			return ExitCodeOK
		}
		return ExitCodeError
	}
	emu.Exit()
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
