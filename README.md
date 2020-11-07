![logo](./logo.png)

# 🌏 Worldwide
![Go](https://github.com/pokemium/Worldwide/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/pokemium/Worldwide)](https://goreportcard.com/report/github.com/pokemium/Worldwide)
[![GitHub stars](https://img.shields.io/github/stars/pokemium/Worldwide)](https://github.com/pokemium/Worldwide/stargazers)
[![GitHub license](https://img.shields.io/github/license/pokemium/Worldwide)](https://github.com/pokemium/Worldwide/blob/master/LICENSE)

日本語のドキュメントは[こちら](./README.ja.md)

GameBoyColor emulator written in golang.  

This emulator can play almost all ROMs work without problems and has many features.


<img src="https://imgur.com/RrOKzJB.png" width="320px"> <img src="https://imgur.com/yIIlkKq.png" width="320px"><br/>
<img src="https://imgur.com/02YAzow.png" width="320px"> <img src="https://imgur.com/QCXeV3B.png" width="320px">


## 🚩 Features & TODO list
- [x] 60fps
- [x] Pass [cpu_instrs](https://github.com/retrio/gb-test-roms/tree/master/cpu_instrs) and [instr_timing](https://github.com/retrio/gb-test-roms/tree/master/instr_timing)
- [x] Low CPU consumption
- [x] Sound(ported from goboy)
- [x] GameBoy Color ROM support
- [x] Multi-platform support
- [x] Joypad support
- [x] MBC1, MBC2, MBC3, MBC5 support
- [x] RTC
- [x] SRAM save
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

Download the prebuilt executable files from [here](https://github.com/pokemium/Worldwide/releases).

```sh
./worldwide.exe [filename]
```

Only files with `.gb` or `.gbc` suffix are allowed.
The configuration file `worldwide.toml` will be created when you run the emulator for the first time.

## 🐛 Debug

You can play this emulator in debug mode.

```sh
./worldwide.exe --debug "***.gb"
```

<img src="https://imgur.com/YxQF9AF.png">

## ✨ HQ2x

You can play games in HQ2x(high-resolution) mode.

HQ2x can be enabled in config file.

<img src="https://imgur.com/bu6WanY.png" width="320px"> <img src="https://imgur.com/OntekWj.png" width="320px">

## 🔨 Build

For those who want to build from source code.

Requirements
- Go 1.15
- make

You should get one of the following executable file:
* `worldwide.exe` for Microsoft Windows
* `worldwide.app` for macOS
* `worldwide` for GNU/Linux and \*BSD

You can specify the ROM to execute:
```sh
make run ROM="***.gb"
```

## 📥 Download

Please download [here](https://github.com/pokemium/Worldwide/releases).

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
| <kbd>D + S</kbd>     | Memory Dump  |
| <kbd>L</kbd>         | Memory Load |
