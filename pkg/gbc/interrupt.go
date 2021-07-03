package gbc

import "gbc/pkg/util"

// IMESwitch - ${Count}サイクル後にIMEを${Value}の値に切り替える ${Working}=falseのときは無効
type IMESwitch struct {
	Count   uint
	Value   bool
	Working bool
}

type intrIEIF struct {
	VBlank, LCDSTAT, Timer, Serial, Joypad struct {
		IE, IF bool
	}
}

func (g *GBC) ieif() intrIEIF {
	ieif := intrIEIF{}
	ieio, ifio := g.IO[IEIO-0xff00], g.IO[IFIO-0xff00]

	ieif.VBlank.IE, ieif.VBlank.IF = util.Bit(ieio, 0), util.Bit(ifio, 0)
	ieif.LCDSTAT.IE, ieif.LCDSTAT.IF = util.Bit(ieio, 1), util.Bit(ifio, 1)
	ieif.Timer.IE, ieif.Timer.IF = util.Bit(ieio, 2), util.Bit(ifio, 2)
	ieif.Serial.IE, ieif.Serial.IF = util.Bit(ieio, 3), util.Bit(ifio, 3)
	ieif.Joypad.IE, ieif.Joypad.IF = util.Bit(ieio, 4), util.Bit(ifio, 4)

	return ieif
}

// ------------ VBlank --------------------

func (g *GBC) setVBlankFlag(b bool) {
	if b {
		g.storeIO(IFIO, g.loadIO(IFIO)|0x01)
		return
	}
	g.storeIO(IFIO, g.loadIO(IFIO)&0xfe)
}

func (g *GBC) setLCDSTATFlag(b bool) {
	if b {
		g.storeIO(IFIO, g.loadIO(IFIO)|0x02)
		return
	}
	g.storeIO(IFIO, g.loadIO(IFIO)&0xfd)
}

func (g *GBC) setSerialFlag(b bool) {
	if b {
		g.storeIO(IFIO, g.loadIO(IFIO)|0x08)
		return
	}
	g.storeIO(IFIO, g.loadIO(IFIO)&0xf7)
}

func (g *GBC) getJoypadEnable() bool {
	return util.Bit(g.loadIO(IEIO), 4)
}

func (g *GBC) setJoypadFlag(b bool) {
	if b {
		g.storeIO(IFIO, g.loadIO(IFIO)|0x10)
		return
	}
	g.storeIO(IFIO, g.loadIO(IFIO)&0xef)
}

// ------------ trigger --------------------

func (g *GBC) triggerInterrupt() {
	g.Reg.IME, g.halt = false, false
	g.updateTimer(5) // https://gbdev.gg8.se/wiki/articles/Interrupts#InterruptServiceRoutine
	g.pushPC()
}

func (g *GBC) triggerVBlank() {
	if util.Bit(g.loadIO(LCDCIO), 7) {
		g.setVBlankFlag(false)
		g.triggerInterrupt()
		g.Reg.PC = 0x0040
	}
}

func (g *GBC) triggerLCDC() {
	g.setLCDSTATFlag(false)
	g.triggerInterrupt()
	g.Reg.PC = 0x0048
}

func (g *GBC) triggerTimer() {
	g.clearTimerFlag()
	g.triggerInterrupt()
	g.Reg.PC = 0x0050
}

func (g *GBC) triggerSerial() {
	g.setSerialFlag(false)
	g.triggerInterrupt()
	g.Reg.PC = 0x0058
}

func (g *GBC) triggerJoypad() {
	g.setJoypadFlag(false)
	g.triggerInterrupt()
	g.Reg.PC = 0x0060
}
