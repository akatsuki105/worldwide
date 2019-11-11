# GameBoy
[![Build Status](https://travis-ci.com/Akatsuki-py/gameboy.svg?branch=master)](https://travis-ci.com/Akatsuki-py/gameboy)

GameBoy emulator written in golang.

<img src="https://imgur.com/1wzyYLx.gif" width="176" height="158"> <img src="https://imgur.com/IxVW33Z.gif" width="176" height="158">

## Features

- 60fps
- cpu_instrs in retrio/gb-test-roms is all clear (STOP instruction doesn't work...)
- APU is implemented (thank to goboy)
- CGB software is OK
- Multi-platform (Win10 and Ubuntu18.04 is checked.)
- Joystick is OK (Xbox 360 Controller and Logitech Gamepad F310)
- ROM-only and MBC1, MBC2, MBC3 is OK
- coredump is enabled
- window expansion is enabled (E or R)

## Usage

- go 1.13

```
go run gb.go "xxx.gb"
```

or 

```
go build
./gameboy "xxx.gb"
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
