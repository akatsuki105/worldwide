# Worldwide
[![Build Status](https://travis-ci.com/Akatsuki-py/Worldwide.svg?branch=master)](https://travis-ci.com/Akatsuki-py/Worldwide)

GameBoyColor emulator written in golang.

<img src="https://imgur.com/cFugCTA.gif" width="320" height="288">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://imgur.com/8YR987D.png" width="320" height="288">


<img src="https://imgur.com/8eDP0un.png" width="320" height="288">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://imgur.com/2zwsb84.png" width="320" height="288">

## Features

- 60fps
- cpu_instrs in retrio/gb-test-roms is all clear
- APU is implemented (thank to goboy)
- CGB ROM is OK
- Multi-platform (Win10 and Ubuntu18.04 is checked.)
- Joystick is OK (Xbox 360 Controller)
- ROM-only and MBC1, MBC2, MBC3, MBC5 is OK
- coredump is enabled
- window expansion is enabled (E or R)

## Usage

- go 1.13

```
go run ./cmd/main.go "xxx.gb"
```

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
