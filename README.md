![logo](./logo.png)

# 🌏 Worldwide
![Go](https://github.com/Akatsuki-py/Worldwide/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Akatsuki-py/Worldwide)](https://goreportcard.com/report/github.com/Akatsuki-py/Worldwide)

日本語のドキュメントは[こちら](./README.ja.md)

GameBoyColor emulator written in golang.  

Almost all ROMs work without problems, and have a wide range of functions, including sound, save, and some communication functions. 


<img src="https://imgur.com/rCduRUc.gif">

## 🚩 Features & TODO list
- [x] 60fps
- [x] pass [cpu_instrs](https://github.com/retrio/gb-test-roms/tree/master/cpu_instrs) and [instr_timing](https://github.com/retrio/gb-test-roms/tree/master/instr_timing)
- [x] Low CPU consumption
- [x] Sound(ported from goboy)
- [x] GameBoy Color ROM support
- [x] Multi-platform support
- [x] Xbox 360 Controller support
- [x] [WebAssembly partial support](https://akatsuki-py.github.io/Worldwide/wasm.html)
- [x] MBC1
- [x] MBC2
- [x] MBC3
- [x] MBC5
- [x] RTC
- [x] System save
- [x] Quick save
- [x] Resizable window
- [x] Pallete color change in DMG
- [x] Serial DMG communication in local network
- [x] RaspberryPi support
- [ ] Serial CGB communication in local network
- [ ] Serial communication with remote network
- [ ] GUI Menu 
- [ ] WebAssembly Audio support
- [ ] SuperGameBoy support

## 🎮 Usage

Download worldwide.exe from [here](https://github.com/Akatsuki-py/Worldwide/releases).

```sh
./worldwide.exe "***.gb" # or ***.gbc
```

## 🔨 Build

For those who want to build from source code.

requirements
- go 1.13
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

## 💻 Web version

You can play [my emulator on website](https://akatsuki-py.github.io/Worldwide/).

This uses webAssembly and javascript.

<img src="https://imgur.com/7ZJxQIu.png">
