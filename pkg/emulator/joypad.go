package emulator

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gopkg.in/ini.v1"
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

// FormatJoypad Joypadの状態をbyteにして返す
func (pad *Joypad) FormatJoypad() byte {
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

// handleJoypad ジョイパッド入力処理
func (cpu *CPU) handleJoypad(win *pixelgl.Window) {
	joystickName := win.JoystickName(player0)

	// A
	if btnA(win, joystickName, cpu.config) {
		cpu.joypad.Button[0] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Button[0] = false
	}

	// B
	if btnB(win, joystickName, cpu.config) {
		cpu.joypad.Button[1] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Button[1] = false
	}

	// select
	if btnSelect(win, joystickName, cpu.config) {
		cpu.joypad.Button[2] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Button[2] = false
	}

	// start
	if btnStart(win, joystickName, cpu.config) {
		cpu.joypad.Button[3] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Button[3] = false
	}

	// 右
	if keyRight(win, joystickName, cpu.config) {
		cpu.joypad.Direction[0] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Direction[0] = false
	}

	// 左
	if keyLeft(win, joystickName, cpu.config) {
		cpu.joypad.Direction[1] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Direction[1] = false
	}

	// 上
	if keyUp(win, joystickName, cpu.config) {
		cpu.joypad.Direction[2] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Direction[2] = false
	}

	// 下
	if keyDown(win, joystickName, cpu.config) {
		cpu.joypad.Direction[3] = true
		if cpu.Reg.IME && cpu.getJoypadEnable() {
			cpu.triggerJoypad()
		}
	} else {
		cpu.joypad.Direction[3] = false
	}

	// coredump
	if btnSaveData(win, joystickName) {
		cpu.dumpData()
	}
	if btnLoadData(win, joystickName) {
		cpu.loadData()
	}

	// expand
	if btnExpandDisplay(win, joystickName, cpu.config) {
		cpu.expand *= 2
		time.Sleep(time.Millisecond * 400)
		win.SetBounds(pixel.R(0, 0, float64(width*cpu.expand), float64(height*cpu.expand)))
	}
	if btnCollapseDisplay(win, joystickName, cpu.config) {
		if cpu.expand >= 2 {
			cpu.expand /= 2
			time.Sleep(time.Millisecond * 400)
			win.SetBounds(pixel.R(0, 0, float64(width*cpu.expand), float64(height*cpu.expand)))
		}
	}
}

func btnA(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyX)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyX)
	}
	return win.JoystickPressed(player0, keyMap.Key("A").MustInt())
}

func btnB(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyZ)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyZ)
	}
	return win.JoystickPressed(player0, keyMap.Key("B").MustInt())
}

func btnStart(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyEnter)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyEnter)
	}
	return win.JoystickPressed(player0, keyMap.Key("Start").MustInt())
}

func btnSelect(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyRightShift)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyRightShift)
	}
	return win.JoystickPressed(player0, keyMap.Key("Select").MustInt())
}

func keyUp(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyUp)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyUp)
	}
	// joystick
	negative := keyMap.Key("VerticalNegative").MustInt()
	if negative == 0 {
		return win.JoystickAxis(player0, keyMap.Key("Vertical").MustInt()) >= 1
	}
	return win.JoystickAxis(player0, keyMap.Key("Vertical").MustInt()) <= -1
}

func keyDown(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyDown)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyDown)
	}
	// joystick
	negative := keyMap.Key("VerticalNegative").MustInt()
	if negative == 0 {
		return win.JoystickAxis(player0, keyMap.Key("Vertical").MustInt()) <= -1
	}
	return win.JoystickAxis(player0, keyMap.Key("Vertical").MustInt()) >= 1
}

func keyRight(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyRight)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyRight)
	}
	return win.JoystickAxis(player0, keyMap.Key("Horizontal").MustInt()) >= 1
}

func keyLeft(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyLeft)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyLeft)
	}
	return win.JoystickAxis(player0, keyMap.Key("Horizontal").MustInt()) <= -1
}

func btnExpandDisplay(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyE)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyE)
	}
	return win.JoystickPressed(player0, keyMap.Key("Expand").MustInt())
}

func btnCollapseDisplay(win *pixelgl.Window, joystickName string, config *ini.File) bool {
	if joystickName == "" {
		// keyboard
		return win.Pressed(pixelgl.KeyR)
	}

	keyMap, err := config.GetSection(joystickName)
	if err != nil {
		// use keyboard
		return win.Pressed(pixelgl.KeyR)
	}
	return win.JoystickPressed(player0, keyMap.Key("Collapse").MustInt())
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
