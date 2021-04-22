package joypad

import (
	"gbc/pkg/util"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// Joypad state
type Joypad struct {
	P1                byte
	Button, Direction [4]bool // start, select, b, a, down, up, left, right
}

const (
	Pressed = iota + 1
	Pause
)

// Output returns joypad state in bitfield format
func (pad *Joypad) Output() byte {
	joypad := byte(0x00)
	if p15 := !util.Bit(pad.P1, 5); p15 {
		for i := 0; i < 4; i++ {
			if pad.Button[i] {
				joypad |= (1 << uint(i))
			}
		}
	}
	if p14 := !util.Bit(pad.P1, 4); p14 {
		for i := 0; i < 4; i++ {
			if pad.Direction[i] {
				joypad |= (1 << uint(i))
			}
		}
	}
	return ^joypad
}

// Input joypad
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

	// right
	if keyRight(threshold) {
		pad.Direction[0] = true
		result = Pressed
	} else {
		pad.Direction[0] = false
	}

	// left
	if keyLeft(threshold) {
		pad.Direction[1] = true
		result = Pressed
	} else {
		pad.Direction[1] = false
	}

	// up
	if keyUp(threshold) {
		pad.Direction[2] = true
		result = Pressed
	} else {
		pad.Direction[2] = false
	}

	// down
	if keyDown(threshold) {
		pad.Direction[3] = true
		result = Pressed
	} else {
		pad.Direction[3] = false
	}

	if btnPause() {
		result = Pause
	}

	return result
}

func btnA(pad uint) bool {
	return ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton(pad)) || ebiten.IsKeyPressed(ebiten.KeyX) || ebiten.IsKeyPressed(ebiten.KeyS)
}

func btnB(pad uint) bool {
	return ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton(pad)) || ebiten.IsKeyPressed(ebiten.KeyZ) || ebiten.IsKeyPressed(ebiten.KeyA)
}

func btnStart(pad uint) bool {
	return ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton(pad)) || ebiten.IsKeyPressed(ebiten.KeyEnter)
}

func btnSelect(pad uint) bool {
	return ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton(pad)) || ebiten.IsKeyPressed(ebiten.KeyShift)
}

func keyUp(threshold float64) bool {
	if threshold > 0 && ebiten.GamepadAxis(0, 1) > threshold {
		return true
	}
	if threshold < 0 && ebiten.GamepadAxis(0, 1) < threshold {
		return true
	}

	return ebiten.IsKeyPressed(ebiten.KeyUp)
}

func keyDown(threshold float64) bool {
	if threshold > 0 && -ebiten.GamepadAxis(0, 1) > threshold {
		return true
	}
	if threshold < 0 && -ebiten.GamepadAxis(0, 1) < threshold {
		return true
	}

	return ebiten.IsKeyPressed(ebiten.KeyDown)
}

func keyRight(threshold float64) bool {
	if threshold > 0 && ebiten.GamepadAxis(0, 0) > threshold {
		return true
	}
	if threshold < 0 && ebiten.GamepadAxis(0, 0) > -threshold {
		return true
	}

	return ebiten.IsKeyPressed(ebiten.KeyRight)
}

func keyLeft(threshold float64) bool {
	if threshold > 0 && ebiten.GamepadAxis(0, 0) < -threshold {
		return true
	}
	if threshold < 0 && ebiten.GamepadAxis(0, 0) < threshold {
		return true
	}

	return ebiten.IsKeyPressed(ebiten.KeyLeft)
}

func btnPause() bool {
	return ebiten.IsKeyPressed(ebiten.KeyP)
}
