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
func (pad *Joypad) Input() bool {
	pressed := false

	// A
	if btnA() {
		old := pad.Button[0]
		pad.Button[0] = true
		if !old && pad.Button[0] {
			pressed = true
		}
	} else {
		pad.Button[0] = false
	}

	// B
	if btnB() {
		old := pad.Button[1]
		pad.Button[1] = true
		if !old && pad.Button[1] {
			pressed = true
		}
	} else {
		pad.Button[1] = false
	}

	// select
	if btnSelect() {
		old := pad.Button[2]
		pad.Button[2] = true
		if !old && pad.Button[2] {
			pressed = true
		}
	} else {
		pad.Button[2] = false
	}

	// start
	if btnStart() {
		old := pad.Button[3]
		pad.Button[3] = true
		if !old && pad.Button[3] {
			pressed = true
		}
	} else {
		pad.Button[3] = false
	}

	// right
	if keyRight() {
		old := pad.Direction[0]
		pad.Direction[0] = true
		if !old && pad.Direction[0] {
			pressed = true
		}
	} else {
		pad.Direction[0] = false
	}

	// left
	if keyLeft() {
		old := pad.Direction[1]
		pad.Direction[1] = true
		if !old && pad.Direction[1] {
			pressed = true
		}
	} else {
		pad.Direction[1] = false
	}

	// up
	if keyUp() {
		old := pad.Direction[2]
		pad.Direction[2] = true
		if !old && pad.Direction[2] {
			pressed = true
		}
	} else {
		pad.Direction[2] = false
	}

	// down
	if keyDown() {
		old := pad.Direction[3]
		pad.Direction[3] = true
		if !old && pad.Direction[3] {
			pressed = true
		}
	} else {
		pad.Direction[3] = false
	}

	return pressed
}

func btnA() bool {
	return ebiten.IsKeyPressed(ebiten.KeyX)
}

func btnB() bool {
	return ebiten.IsKeyPressed(ebiten.KeyZ)
}

func btnStart() bool {
	return ebiten.IsKeyPressed(ebiten.KeyEnter)
}

func btnSelect() bool {
	return ebiten.IsKeyPressed(ebiten.KeyBackspace)
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
