package joypad

import (
	"github.com/hajimehoshi/ebiten"
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
	Pressed = iota + 1
	Save
	Load
	Pause
)

func gamepad(n uint) ebiten.GamepadButton {
	switch n {
	case 0:
		return ebiten.GamepadButton0
	case 1:
		return ebiten.GamepadButton1
	case 2:
		return ebiten.GamepadButton2
	case 3:
		return ebiten.GamepadButton3
	case 4:
		return ebiten.GamepadButton4
	case 5:
		return ebiten.GamepadButton5
	case 6:
		return ebiten.GamepadButton6
	case 7:
		return ebiten.GamepadButton7
	case 8:
		return ebiten.GamepadButton8
	}

	return ebiten.GamepadButton0
}

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
func (pad *Joypad) Input(padA, padB, padStart, padSelect uint, threshold float64) (result int) {

	// A
	if btnA(padA) {
		pad.Button[0] = true
		result = Pressed
	} else {
		pad.Button[0] = false
	}

	// B
	if btnB(padB) {
		pad.Button[1] = true
		result = Pressed
	} else {
		pad.Button[1] = false
	}

	// select
	if btnSelect(padSelect) {
		pad.Button[2] = true
		result = Pressed
	} else {
		pad.Button[2] = false
	}

	// start
	if btnStart(padStart) {
		pad.Button[3] = true
		result = Pressed
	} else {
		pad.Button[3] = false
	}

	// 右
	if keyRight(threshold) {
		pad.Direction[0] = true
		result = Pressed
	} else {
		pad.Direction[0] = false
	}

	// 左
	if keyLeft(threshold) {
		pad.Direction[1] = true
		result = Pressed
	} else {
		pad.Direction[1] = false
	}

	// 上
	if keyUp(threshold) {
		pad.Direction[2] = true
		result = Pressed
	} else {
		pad.Direction[2] = false
	}

	// 下
	if keyDown(threshold) {
		pad.Direction[3] = true
		result = Pressed
	} else {
		pad.Direction[3] = false
	}

	if btnSaveData() {
		result = Save
	}
	if btnLoadData() {
		result = Load
	}

	if btnPause() {
		result = Pause
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

func btnA(pad uint) bool {
	return ebiten.IsGamepadButtonPressed(0, gamepad(pad)) || ebiten.IsKeyPressed(ebiten.KeyX) || ebiten.IsKeyPressed(ebiten.KeyS)
}

func btnB(pad uint) bool {
	return ebiten.IsGamepadButtonPressed(0, gamepad(pad)) || ebiten.IsKeyPressed(ebiten.KeyZ) || ebiten.IsKeyPressed(ebiten.KeyA)
}

func btnStart(pad uint) bool {
	return ebiten.IsGamepadButtonPressed(0, gamepad(pad)) || ebiten.IsKeyPressed(ebiten.KeyEnter)
}

func btnSelect(pad uint) bool {
	return ebiten.IsGamepadButtonPressed(0, gamepad(pad)) || ebiten.IsKeyPressed(ebiten.KeyShift)
}

func keyUp(threshold float64) bool {
	if ebiten.GamepadAxis(0, 1) > threshold {
		return true
	}

	return ebiten.IsKeyPressed(ebiten.KeyUp)
}

func keyDown(threshold float64) bool {
	if ebiten.GamepadAxis(0, 1) < -threshold {
		return true
	}

	return ebiten.IsKeyPressed(ebiten.KeyDown)
}

func keyRight(threshold float64) bool {
	if ebiten.GamepadAxis(0, 0) > threshold {
		return true
	}

	return ebiten.IsKeyPressed(ebiten.KeyRight)
}

func keyLeft(threshold float64) bool {
	if ebiten.GamepadAxis(0, 0) < -threshold {
		return true
	}

	return ebiten.IsKeyPressed(ebiten.KeyLeft)
}

func btnExpandDisplay() bool {
	return ebiten.IsKeyPressed(ebiten.KeyE)
}

func btnCollapseDisplay() bool {
	return ebiten.IsKeyPressed(ebiten.KeyR)
}

func btnSaveData() bool {
	return ebiten.IsKeyPressed(ebiten.KeyD) && ebiten.IsKeyPressed(ebiten.KeyS)
}

func btnLoadData() bool {
	return ebiten.IsKeyPressed(ebiten.KeyL)
}

func btnPause() bool {
	return ebiten.IsKeyPressed(ebiten.KeyP)
}
