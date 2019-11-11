package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"./emulator"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	flag.Parse()
	romData := readGB(flag.Arg(0))

	cpu := &emulator.CPU{}
	cpu.ParseROMHeader(romData)
	cpu.LoadROM(romData)
	cpu.InitCPU()
	cpu.InitAPU()

	go cpu.Debug()

	pixelgl.Run(cpu.Render)
}

func readGB(path string) []byte {
	if path == "" {
		fmt.Println("please enter gb file path")
		os.Exit(0)
	}
	if filepath.Ext(path) != ".gb" {
		fmt.Println("please select .gb file")
		os.Exit(0)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}
