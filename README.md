![logo](./logo.png)

# 🌏 worldwide
![Go](https://github.com/pokemium/worldwide/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/pokemium/worldwide)](https://goreportcard.com/report/github.com/pokemium/worldwide)
[![GitHub stars](https://img.shields.io/github/stars/pokemium/worldwide)](https://github.com/pokemium/worldwide/stargazers)
[![GitHub license](https://img.shields.io/github/license/pokemium/worldwide)](https://github.com/pokemium/worldwide/blob/master/LICENSE)

日本語のドキュメントは[こちら](./README.ja.md)

GameBoyColor emulator written in golang.  

This emulator can play a lot of ROMs work without problems and has many features.

<img src="https://imgur.com/ZlrXAW9.png" width="320px"> <img src="https://imgur.com/xVqjkrk.png" width="320px"><br/>
<img src="https://imgur.com/E7oob9c.png" width="320px"> <img src="https://imgur.com/nYpkH95.png" width="320px">

## 🚩 Features & TODO list
- [x] 60fps
- [x] Pass [cpu_instrs](https://github.com/retrio/gb-test-roms/tree/master/cpu_instrs) and [instr_timing](https://github.com/retrio/gb-test-roms/tree/master/instr_timing)
- [x] Low CPU consumption
- [x] Sound(ported from goboy)
- [x] GameBoy Color ROM support
- [x] Multi-platform support
- [x] MBC1, MBC2, MBC3, MBC5 support
- [x] RTC
- [x] SRAM save
- [x] Resizable window
- [x] HTTP server API
- [ ] Plugins support
- [ ] [Libretro](https://docs.libretro.com/) support
- [ ] Netplay in local network
- [ ] Netplay in global network
- [ ] SGB support
- [ ] Shader support

## 🎮 Usage

Download binary from [here](https://github.com/pokemium/worldwide/releases).

```sh
./worldwide "***.gb" # or ***.gbc
```

## 🐛 HTTP Server

`worldwide` contains an HTTP server, and the user can give various instructions to it through HTTP requests.

Please read [Server Document](./server/README.md).

## 🔨 Build

For those who want to build from source code.

Requirements
- Go 1.16
- make

```sh
make build                              # If you use Windows, `make build-windows`
./build/darwin-amd64/worldwide "***.gb" # If you use Windows, `./build/windows-amd64/worldwide.exe "***.gb"`
```

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
| <kbd>Backspace</kbd> | Select button |
