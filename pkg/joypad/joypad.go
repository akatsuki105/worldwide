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
func (pad *Joypad) Input() (result int) {

	// A
	if btnA() {
		pad.Button[0] = true
		result = Pressed
	} else {
		pad.Button[0] = false
	}

	// B
	if btnB() {
		pad.Button[1] = true
		result = Pressed
	} else {
		pad.Button[1] = false
	}

	// select
	if btnSelect() {
		pad.Button[2] = true
		result = Pressed
	} else {
		pad.Button[2] = false
	}

	// start
	if btnStart() {
		pad.Button[3] = true
		result = Pressed
	} else {
		pad.Button[3] = false
	}

	// 右
	if keyRight() {
		pad.Direction[0] = true
		result = Pressed
	} else {
		pad.Direction[0] = false
	}

	// 左
	if keyLeft() {
		pad.Direction[1] = true
		result = Pressed
	} else {
		pad.Direction[1] = false
	}

	// 上
	if keyUp() {
		pad.Direction[2] = true
		result = Pressed
	} else {
		pad.Direction[2] = false
	}

	// 下
	if keyDown() {
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

	// expand
	if btnExpandDisplay() {
		result = Expand
	}
	if btnCollapseDisplay() {
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

func btnA() bool {
	return ebiten.IsKeyPressed(ebiten.KeyX) || ebiten.IsKeyPressed(ebiten.KeyS)
}

func btnB() bool {
	return ebiten.IsKeyPressed(ebiten.KeyZ) || ebiten.IsKeyPressed(ebiten.KeyA)
}

func btnStart() bool {
	return ebiten.IsKeyPressed(ebiten.KeyEnter)
}

func btnSelect() bool {
	return ebiten.IsKeyPressed(ebiten.KeyShift)
}

func keyUp() bool {
	return ebiten.IsKeyPressed(ebiten.KeyUp)
}

func keyDown() bool {
	return ebiten.IsKeyPressed(ebiten.KeyDown)
}

func keyRight() bool {
	return ebiten.IsKeyPressed(ebiten.KeyRight)
}

func keyLeft() bool {
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
