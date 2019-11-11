package emulator

import (
	"github.com/faiface/pixel/pixelgl"
)

// Joypad Joypadの入力を管理する
type Joypad struct {
	Button    [4]bool // 0bスタートセレクトBA
	Direction [4]bool // 0b下上左右
}

const (
	Player0 = 0
)

var Xbox360Controller map[string]int = map[string]int{"A": 1, "B": 0, "Start": 8, "Select": 7, "Horizontal": 0, "Vertical": 1}
var LogitechGamepadF310 map[string]int = map[string]int{"A": 1, "B": 0, "Start": 8, "Select": 7, "Horizontal": 0, "Vertical": 1}

func (cpu *CPU) formatJoypad() byte {
	joypad := byte(0x00)
	if cpu.getP15() {
		for i := 0; i < 4; i++ {
			if cpu.joypad.Button[i] {
				joypad |= (1 << uint(i))
			}
		}
	}
	if cpu.getP14() {
		for i := 0; i < 4; i++ {
			if cpu.joypad.Direction[i] {
				joypad |= (1 << uint(i))
			}
		}
	}
	return ^joypad
}

func (cpu *CPU) getP14() bool {
	JOYPAD := cpu.FetchMemory8(0xff00)
	if JOYPAD&0x10 == 0 {
		return true
	}
	return false
}

func (cpu *CPU) getP15() bool {
	JOYPAD := cpu.FetchMemory8(0xff00)
	if JOYPAD&0x20 == 0 {
		return true
	}
	return false
}

// handleJoypad ジョイパッド入力処理
func (cpu *CPU) handleJoypad(win *pixelgl.Window) {
	// A
	if btnA(win) {
		cpu.joypad.Button[0] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Button[0] = false
	}

	// B
	if btnB(win) {
		cpu.joypad.Button[1] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Button[1] = false
	}

	// select
	if btnSelect(win) {
		cpu.joypad.Button[2] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Button[2] = false
	}

	// start
	if btnStart(win) {
		cpu.joypad.Button[3] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Button[3] = false
	}

	// 右
	if keyRight(win) {
		cpu.joypad.Direction[0] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Direction[0] = false
	}

	// 左
	if keyLeft(win) {
		cpu.joypad.Direction[1] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Direction[1] = false
	}

	// 上
	if keyUp(win) {
		cpu.joypad.Direction[2] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Direction[2] = false
	}

	// 下
	if keyDown(win) {
		cpu.joypad.Direction[3] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Direction[3] = false
	}
}

func btnA(win *pixelgl.Window) bool {
	switch win.JoystickName(Player0) {
	case "Xbox 360 Controller":
		return win.JoystickPressed(Player0, Xbox360Controller["A"])
	case "Logitech Gamepad F310":
		return win.JoystickPressed(Player0, LogitechGamepadF310["A"])
	default:
		return win.Pressed(pixelgl.KeyX)
	}
}

func btnB(win *pixelgl.Window) bool {
	switch win.JoystickName(Player0) {
	case "Xbox 360 Controller":
		return win.JoystickPressed(Player0, Xbox360Controller["B"])
	case "Logitech Gamepad F310":
		return win.JoystickPressed(Player0, LogitechGamepadF310["B"])
	default:
		return win.Pressed(pixelgl.KeyZ)
	}
}

func btnStart(win *pixelgl.Window) bool {
	switch win.JoystickName(Player0) {
	case "Xbox 360 Controller":
		return win.JoystickPressed(Player0, Xbox360Controller["Start"])
	case "Logitech Gamepad F310":
		return win.JoystickPressed(Player0, LogitechGamepadF310["Start"])
	default:
		return win.Pressed(pixelgl.KeyEnter)
	}
}

func btnSelect(win *pixelgl.Window) bool {
	switch win.JoystickName(Player0) {
	case "Xbox 360 Controller":
		return win.JoystickPressed(Player0, Xbox360Controller["Select"])
	case "Logitech Gamepad F310":
		return win.JoystickPressed(Player0, LogitechGamepadF310["Select"])
	default:
		return win.Pressed(pixelgl.KeyRightShift)
	}
}

func keyUp(win *pixelgl.Window) bool {
	switch win.JoystickName(Player0) {
	case "Xbox 360 Controller":
		return win.JoystickAxis(Player0, Xbox360Controller["Vertical"]) >= 1
	case "Logitech Gamepad F310":
		return win.JoystickAxis(Player0, LogitechGamepadF310["Vertical"]) <= -1
	default:
		return win.Pressed(pixelgl.KeyDown)
	}
}

func keyDown(win *pixelgl.Window) bool {
	switch win.JoystickName(Player0) {
	case "Xbox 360 Controller":
		return win.JoystickAxis(Player0, Xbox360Controller["Vertical"]) <= -1
	case "Logitech Gamepad F310":
		return win.JoystickAxis(Player0, LogitechGamepadF310["Vertical"]) >= 1
	default:
		return win.Pressed(pixelgl.KeyDown)
	}
}

func keyRight(win *pixelgl.Window) bool {
	switch win.JoystickName(Player0) {
	case "Xbox 360 Controller":
		return win.JoystickAxis(Player0, Xbox360Controller["Horizontal"]) >= 1
	case "Logitech Gamepad F310":
		return win.JoystickAxis(Player0, LogitechGamepadF310["Horizontal"]) >= 1
	default:
		return win.Pressed(pixelgl.KeyRight)
	}
}

func keyLeft(win *pixelgl.Window) bool {
	switch win.JoystickName(Player0) {
	case "Xbox 360 Controller":
		return win.JoystickAxis(Player0, Xbox360Controller["Horizontal"]) <= -1
	case "Logitech Gamepad F310":
		return win.JoystickAxis(Player0, LogitechGamepadF310["Horizontal"]) <= -1
	default:
		return win.Pressed(pixelgl.KeyLeft)
	}
}
