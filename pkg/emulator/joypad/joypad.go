package joypad

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var Handler = [8](func() bool){
	btnA, btnB, btnSelect, btnStart, keyRight, keyLeft, keyUp, keyDown,
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
