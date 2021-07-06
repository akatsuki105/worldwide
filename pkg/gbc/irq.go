package gbc

import "gbc/pkg/util"

// ------------ trigger --------------------

func (g *GBC) triggerInterrupt() {
	g.pushPC()
}

func (g *GBC) triggerVBlank() {
	g.IO[IFIO-0xff00] = util.SetBit8(g.IO[IFIO-0xff00], 0, false)
	g.triggerInterrupt()
	g.Reg.PC = 0x0040
}

func (g *GBC) triggerLCDC() {
	g.IO[IFIO-0xff00] = util.SetBit8(g.IO[IFIO-0xff00], 1, false)
	g.triggerInterrupt()
	g.Reg.PC = 0x0048
}

func (g *GBC) triggerTimer() {
	g.IO[IFIO-0xff00] = util.SetBit8(g.IO[IFIO-0xff00], 2, false)
	g.triggerInterrupt()
	g.Reg.PC = 0x0050
}

func (g *GBC) triggerSerial() {
	g.IO[IFIO-0xff00] = util.SetBit8(g.IO[IFIO-0xff00], 3, false)
	g.triggerInterrupt()
	g.Reg.PC = 0x0058
}

func (g *GBC) triggerJoypad() {
	g.IO[IFIO-0xff00] = util.SetBit8(g.IO[IFIO-0xff00], 4, false)
	g.triggerInterrupt()
	g.Reg.PC = 0x0060
}
