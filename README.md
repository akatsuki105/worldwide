![logo](./logo.png)

# 🌏 Worldwide
![Go](https://github.com/Akatsuki-py/Worldwide/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Akatsuki-py/Worldwide)](https://goreportcard.com/report/github.com/Akatsuki-py/Worldwide)
[![GitHub stars](https://img.shields.io/github/stars/Akatsuki-py/Worldwide)](https://github.com/Akatsuki-py/Worldwide/stargazers)
[![GitHub license](https://img.shields.io/github/license/Akatsuki-py/Worldwide)](https://github.com/Akatsuki-py/Worldwide/blob/master/LICENSE)

日本語のドキュメントは[こちら](./README.ja.md)

GameBoyColor emulator written in golang.  

This emulator can play almost all ROMs work without problems and has many features.


<img src="https://imgur.com/rCduRUc.gif">

## 🚩 Features & TODO list
- [x] 60fps
- [x] Pass [cpu_instrs](https://github.com/retrio/gb-test-roms/tree/master/cpu_instrs) and [instr_timing](https://github.com/retrio/gb-test-roms/tree/master/instr_timing)
- [x] Low CPU consumption
- [x] Sound(ported from goboy)
- [x] GameBoy Color ROM support
- [x] Multi-platform support
- [x] Joypad support
- [x] [WebAssembly partial support](https://akatsuki-py.github.io/Worldwide/wasm.html)
- [x] MBC1, MBC2, MBC3, MBC5 support
- [x] RTC
- [x] System save
- [x] Quick save
- [x] Resizable window
- [x] Pallete color change in DMG
- [x] Serial DMG communication in local network
- [x] RaspberryPi support
- [x] Debugger
- [x] HQ2x mode 
- [ ] Serial CGB communication in local network
- [ ] Serial communication with global network
- [ ] SuperGameBoy support

## 🎮 Usage

Download worldwide.exe from [here](https://github.com/Akatsuki-py/Worldwide/releases).

```sh
./worldwide.exe "***.gb" # or ***.gbc
```

## 🐛 Debug

You can play this emulator in debug mode.

```sh
./worldwide.exe --debug "***.gb"
```

<img src="https://user-images.githubusercontent.com/37920078/81895677-c8fc7f00-95ed-11ea-8377-8f83d68f191e.PNG">

## 🔨 Build

For those who want to build from source code.

requirements
- go 1.14
- make

```sh
make
./worldwide "***.gb" # ./worldwide.exe on Windows
```

## 📥 Download

Please download [here](https://github.com/Akatsuki-py/Worldwide/releases).

## 📄 Command 

| keyboard             | game pad      |
| -------------------- | ------------- |
| <kbd>&larr;</kbd>    | &larr; button |
| <kbd>&uarr;</kbd>    | &uarr; button |
| <kbd>&darr;</kbd>    | &darr; button |
| <kbd>&rarr;</kbd>    | &rarr; button |
| <kbd>X</kbd>         | A button      |
| <kbd>Z</kbd>         | B button      |
| <kbd>Enter</kbd>     | Start button  |
| <kbd>Right shift</kbd> | Select button |
| <kbd>E</kbd>         | Expand display  |
| <kbd>R</kbd>         | Collapse display |
| <kbd>D + S</kbd>     | Memory Dump  |
| <kbd>L</kbd>         | Memory Load |