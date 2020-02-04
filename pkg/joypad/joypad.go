package joypad

import (
	"github.com/faiface/pixel/pixelgl"
)

// Joypad Joypadの入力を管理する
type Joypad struct {
	P1        byte
	Button    [4]bool // 0bスタートセレクトBA
	Direction [4]bool // 0b下上左右
}

type keyList struct {
	A, B, Start, Select, Horizontal, Vertical, Expand, Collapse uint
}

const (
	player0 = 0
)

const (
	Pressed = iota + 1
	Save
	Load
	Expand
	Collapse
)

var Xbox360Controller map[string]int = map[string]int{"A": 1, "B": 0, "Start": 7, "Select": 6, "Horizontal": 0, "Vertical": 1}
var LogitechGamepadF310 map[string]int = map[string]int{"A": 1, "B": 0, "Start": 8, "Select": 7, "Horizontal": 0, "Vertical": 1}
var HORIPAD map[string]int = map[string]int{"A": 2, "B": 1, "Start": 3, "Select": 0, "Horizontal": 0, "Vertical": 1}

// Output Joypadの状態をbyteにして返す
func (pad *Joypad) Output() byte {
	joypad := byte(0x00)
	if pad.getP15() {
		for i := 0; i < 4; i++ {
			if pad.Button[i] {
				joypad |= (1 << uint(i))
			}
		}
	}
	if pad.getP14() {
		for i := 0; i < 4; i++ {
			if pad.Direction[i] {
				joypad |= (1 << uint(i))
			}
		}
	}
	return ^joypad
}

// Input ジョイパッド入力処理
func (pad *Joypad) Input(win *pixelgl.Window) (result int) {
	joystickName := win.JoystickName(player0)

	// A
	if btnA(win, joystickName) {
		pad.Button[0] = true
		result = Pressed
	} else {
		pad.Button[0] = false
	}

	// B
	if btnB(win, joystickName) {
		pad.Button[1] = true
		result = Pressed
	} else {
		pad.Button[1] = false
	}

	// select
	if btnSelect(win, joystickName) {
		pad.Button[2] = true
		result = Pressed
	} else {
		pad.Button[2] = false
	}

	// start
	if btnStart(win, joystickName) {
		pad.Button[3] = true
		result = Pressed
	} else {
		pad.Button[3] = false
	}

	// 右
	if keyRight(win, joystickName) {
		pad.Direction[0] = true
		result = Pressed
	} else {
		pad.Direction[0] = false
	}

	// 左
	if keyLeft(win, joystickName) {
		pad.Direction[1] = true
		result = Pressed
	} else {
		pad.Direction[1] = false
	}

	// 上
	if keyUp(win, joystickName) {
		pad.Direction[2] = true
		result = Pressed
	} else {
		pad.Direction[2] = false
	}

	// 下
	if keyDown(win, joystickName) {
		pad.Direction[3] = true
		result = Pressed
	} else {
		pad.Direction[3] = false
	}

	if btnSaveData(win, joystickName) {
		result = Save
	}
	if btnLoadData(win, joystickName) {
		result = Load
	}

	// expand
	if btnExpandDisplay(win, joystickName) {
		result = Expand
	}
	if btnCollapseDisplay(win, joystickName) {
		result = Collapse
	}

	return result
}

func (pad *Joypad) getP14() bool {
	JOYPAD := pad.P1
	if JOYPAD&0x10 == 0 {
		return true
	}
	return false
}

func (pad *Joypad) getP15() bool {
	JOYPAD := pad.P1
	if JOYPAD&0x20 == 0 {
		return true
	}
	return false
}

func btnA(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	case "Xbox 360 Controller":
		return win.JoystickPressed(player0, Xbox360Controller["A"])
	case "Logitech Gamepad F310":
		return win.JoystickPressed(player0, LogitechGamepadF310["A"])
	case "HORI CO.,LTD HORIPAD S":
		return win.JoystickPressed(player0, HORIPAD["A"])
	default:
		return win.Pressed(pixelgl.KeyX) || win.Pressed(pixelgl.KeyS)
	}
}

func btnB(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	case "Xbox 360 Controller":
		return win.JoystickPressed(player0, Xbox360Controller["B"])
	case "Logitech Gamepad F310":
		return win.JoystickPressed(player0, LogitechGamepadF310["B"])
	case "HORI CO.,LTD HORIPAD S":
		return win.JoystickPressed(player0, HORIPAD["B"])
	default:
		return win.Pressed(pixelgl.KeyZ) || win.Pressed(pixelgl.KeyA)
	}
}

func btnStart(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	case "Xbox 360 Controller":
		return win.JoystickPressed(player0, Xbox360Controller["Start"])
	case "Logitech Gamepad F310":
		return win.JoystickPressed(player0, LogitechGamepadF310["Start"])
	case "HORI CO.,LTD HORIPAD S":
		return win.JoystickPressed(player0, HORIPAD["Start"])
	default:
		return win.Pressed(pixelgl.KeyEnter)
	}
}

func btnSelect(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	case "Xbox 360 Controller":
		return win.JoystickPressed(player0, Xbox360Controller["Select"])
	case "Logitech Gamepad F310":
		return win.JoystickPressed(player0, LogitechGamepadF310["Select"])
	case "HORI CO.,LTD HORIPAD S":
		return win.JoystickPressed(player0, HORIPAD["Select"])
	default:
		return win.Pressed(pixelgl.KeyRightShift)
	}
}

func keyUp(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	case "Xbox 360 Controller":
		return win.JoystickAxis(player0, Xbox360Controller["Vertical"]) >= 1
	case "Logitech Gamepad F310":
		return win.JoystickAxis(player0, LogitechGamepadF310["Vertical"]) <= -1
	case "HORI CO.,LTD HORIPAD S":
		return win.JoystickAxis(player0, HORIPAD["Vertical"]) <= -1
	default:
		return win.Pressed(pixelgl.KeyUp)
	}
}

func keyDown(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	case "Xbox 360 Controller":
		return win.JoystickAxis(player0, Xbox360Controller["Vertical"]) <= -1
	case "Logitech Gamepad F310":
		return win.JoystickAxis(player0, LogitechGamepadF310["Vertical"]) >= 1
	case "HORI CO.,LTD HORIPAD S":
		return win.JoystickAxis(player0, HORIPAD["Vertical"]) >= 1
	default:
		return win.Pressed(pixelgl.KeyDown)
	}
}

func keyRight(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	case "Xbox 360 Controller":
		return win.JoystickAxis(player0, Xbox360Controller["Horizontal"]) >= 1
	case "Logitech Gamepad F310":
		return win.JoystickAxis(player0, LogitechGamepadF310["Horizontal"]) >= 1
	case "HORI CO.,LTD HORIPAD S":
		return win.JoystickAxis(player0, HORIPAD["Horizontal"]) >= 1
	default:
		return win.Pressed(pixelgl.KeyRight)
	}
}

func keyLeft(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	case "Xbox 360 Controller":
		return win.JoystickAxis(player0, Xbox360Controller["Horizontal"]) <= -1
	case "Logitech Gamepad F310":
		return win.JoystickAxis(player0, LogitechGamepadF310["Horizontal"]) <= -1
	case "HORI CO.,LTD HORIPAD S":
		return win.JoystickAxis(player0, HORIPAD["Horizontal"]) <= -1
	default:
		return win.Pressed(pixelgl.KeyLeft)
	}
}

func btnExpandDisplay(win *pixelgl.Window, joystickName string) bool {
	return win.Pressed(pixelgl.KeyE)
}

func btnCollapseDisplay(win *pixelgl.Window, joystickName string) bool {
	return win.Pressed(pixelgl.KeyR)
}

func btnSaveData(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	default:
		return win.Pressed(pixelgl.KeyD) && win.Pressed(pixelgl.KeyS)
	}
}

func btnLoadData(win *pixelgl.Window, joystickName string) bool {
	switch joystickName {
	default:
		return win.Pressed(pixelgl.KeyL)
	}
}
