package joypad

import "gbc/pkg/util"

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

const (
	Pressed = iota + 1
	Pause
)

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
			if pad.button[i] {
				joypad |= (1 << uint(i))
			}
		}
	}
	if p14 := !util.Bit(pad.P1, 4); p14 {
		for i := 0; i < 4; i++ {
			if pad.direction[i] {
				joypad |= (1 << uint(i))
			}
		}
	}
	return ^joypad
}

// Input joypad
func (j *Joypad) Input() bool {
	pressed := false

	// A
	if j.handler[A]() {
		old := j.button[0]
		j.button[0] = true
		if !old && j.button[0] {
			pressed = true
		}
	} else {
		j.button[0] = false
	}

	// B
	if j.handler[B]() {
		old := j.button[1]
		j.button[1] = true
		if !old && j.button[1] {
			pressed = true
		}
	} else {
		j.button[1] = false
	}

	// select
	if j.handler[Select]() {
		old := j.button[2]
		j.button[2] = true
		if !old && j.button[2] {
			pressed = true
		}
	} else {
		j.button[2] = false
	}

	// start
	if j.handler[Start]() {
		old := j.button[3]
		j.button[3] = true
		if !old && j.button[3] {
			pressed = true
		}
	} else {
		j.button[3] = false
	}

	// right
	if j.handler[Right]() {
		old := j.direction[0]
		j.direction[0] = true
		if !old && j.direction[0] {
			pressed = true
		}
	} else {
		j.direction[0] = false
	}

	// left
	if j.handler[Left]() {
		old := j.direction[1]
		j.direction[1] = true
		if !old && j.direction[1] {
			pressed = true
		}
	} else {
		j.direction[1] = false
	}

	// up
	if j.handler[Up]() {
		old := j.direction[2]
		j.direction[2] = true
		if !old && j.direction[2] {
			pressed = true
		}
	} else {
		j.direction[2] = false
	}

	// down
	if j.handler[Down]() {
		old := j.direction[3]
		j.direction[3] = true
		if !old && j.direction[3] {
			pressed = true
		}
	} else {
		j.direction[3] = false
	}

	return pressed
}
