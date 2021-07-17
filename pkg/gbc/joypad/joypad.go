package joypad

import "github.com/pokemium/worldwide/pkg/util"

const (
	A      = 0
	B      = 1
	Select = 2
	Start  = 3
	Right  = 4
	Left   = 5
	Up     = 6
	Down   = 7
)

// Joypad state
type Joypad struct {
	P1                byte
	button, direction [4]bool // start, select, b, a, down, up, left, right
	handler           [8](func() bool)
}

func New(h [8](func() bool)) *Joypad {
	return &Joypad{
		handler: h,
	}
}

// Output returns joypad state in bitfield format
func (pad *Joypad) Output() byte {
	joypad := byte(0x00)
	if p15 := !util.Bit(pad.P1, 5); p15 {
		for i := 0; i < 4; i++ {
			joypad = util.SetBit8(joypad, i, pad.button[i])
		}
	}
	if p14 := !util.Bit(pad.P1, 4); p14 {
		for i := 0; i < 4; i++ {
			joypad = util.SetBit8(joypad, i, pad.direction[i])
		}
	}
	return ^joypad
}

// Input joypad
func (j *Joypad) Input() bool {
	pressed := false

	// A,B,Start,Select
	for i := 0; i < 4; i++ {
		if j.handler[i]() {
			old := j.button[i]
			j.button[i] = true
			if !old && j.button[i] {
				pressed = true
			}
		} else {
			j.button[i] = false
		}
	}

	// Right, Left, Up, Down
	for i := 0; i < 4; i++ {
		if j.handler[i+4]() {
			old := j.direction[i]
			j.direction[i] = true
			if !old && j.direction[i] {
				pressed = true
			}
		} else {
			j.direction[i] = false
		}
	}

	return pressed
}
