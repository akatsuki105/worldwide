# Worldwide
[![Build Status](https://travis-ci.com/Akatsuki-py/Worldwide.svg?branch=master)](https://travis-ci.com/Akatsuki-py/Worldwide)

GameBoyColor emulator written in golang.

<img src="https://imgur.com/UnmQnVE.gif" width="320" height="288">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://imgur.com/cFugCTA.gif" width="320" height="288">


<img src="https://imgur.com/8YR987D.png" width="320" height="288">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://imgur.com/2zwsb84.png" width="320" height="288">

## Features & TODO list
- [x] 60fps
- [x] [cpu_instrs](https://github.com/retrio/gb-test-roms/tree/master/cpu_instrs) is clear
- [x] Low CPU consumption
- [x] Sound(thank to goboy)
- [x] GameBoy Color ROM support
- [x] Multi-platform support
- [x] Xbox 360 Controller support
- [x] [WebAssembly partial support](https://akatsuki-py.github.io/Worldwide/wasm.html)
- [x] MBC1
- [x] MBC2
- [x] MBC3
- [x] MBC5
- [x] RTC
- [x] Save game data
- [x] Quicksave
- [x] Resizable window
- [x] Pallete color change in DMG
- [x] Pokemon Crystal JPN version
- [x] Serial DMG communication in local network
- [ ] Serial CGB communication in local network
- [ ] Serial communication with remote network
- [ ] GUI Menu 
- [ ] WebAssembly Audio support
- [ ] RaspberryPi support
- [ ] SuperGameBoy support

## Usage

- go 1.13

```
go run ./cmd/main.go "xxx.gb"
```

## Download

Please download [here](https://github.com/Akatsuki-py/Worldwide/releases).

## Command 

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

## web version

You can play [my emulator on website](https://akatsuki-py.github.io/Worldwide/).

This uses webAssembly and javascript.

<img src="https://imgur.com/7ZJxQIu.png">
