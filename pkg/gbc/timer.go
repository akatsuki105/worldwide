package gbc

import (
	"gbc/pkg/util"
)

type TIMAReload struct {
	flag  bool
	value byte
	after bool // ref: [B] in https://gbdev.io/pandocs/#timer-overflow-behaviour
}

type Timer struct {
	tac int    // use in normal timer
	div int    // use in div timer
	sys uint16 // 16 bit system counter. ref: https://gbdev.io/pandocs/Timer_Obscure_Behaviour.html
	OAMDMA
	TIMAReload
	ResetAll bool
	TAC      struct {
		Change bool
		Old    byte
	}
}

type OAMDMA struct {
	start, ptr     uint16
	restart, reptr uint16 // OAM DMA is requested in OAM DMA
}

func (g *GBC) setTimerFlag() {
	g.storeIO(IFIO, g.loadIO(IFIO)|0x04)
	g.halt = false
}

func (g *GBC) clearTimerFlag() {
	g.storeIO(IFIO, g.loadIO(IFIO)&0xfb)
}

func (g *GBC) updateTimer(cycle int) {
	for i := 0; i < cycle; i++ {
		g.tick()
	}
	g.Sound.Buffer(4*cycle, g.boost)
}

// 0: 4096Hz (1024/4 cycle), 1: 262144Hz (16/4 cycle), 2: 65536Hz (64/4 cycle), 3: 16384Hz (256/4 cycle)
var clocks = [4]int{1024 / 4, 16 / 4, 64 / 4, 256 / 4}

func (g *GBC) tick() {
	tac := g.IO[TACIO-0xff00]
	tickFlag := false

	if g.timer.ResetAll {
		g.timer.ResetAll = false
		tickFlag = g.resetTimer()
	}
	if g.timer.TAC.Change && !tickFlag {
		g.timer.TAC.Change = false
		oldTAC, newTAC := g.timer.TAC.Old, g.IO[TACIO-0xff00]
		oldClock, newClock := uint16(clocks[oldTAC&0b11]), uint16(clocks[newTAC&0b11])
		oldEnable, newEnable := oldTAC&0b100 > 0, newTAC&0b100 > 0
		if oldEnable {
			if newEnable {
				tickFlag = g.timer.sys&(oldClock/2) > 0
			} else {
				tickFlag = g.timer.sys&(oldClock/2) > 0 && g.timer.sys&(newClock/2) == 0
			}
		}
	}

	// lag occurs in di, ei
	if g.IMESwitch.Working {
		g.IMESwitch.Count--
		if g.IMESwitch.Count == 0 {
			g.Reg.IME = g.IMESwitch.Value
			g.IMESwitch.Working = false
		}
	}

	g.cycles += 4
	g.timer.sys++ // 16 bit system counter
	g.timer.div++
	if g.timer.div >= 64 {
		g.IO[DIVIO-0xff00]++
		g.timer.div -= 64
	}

	if util.Bit(tac, 2) {
		g.timer.tac++
		if g.timer.tac >= clocks[tac&0b11] {
			g.timer.tac -= clocks[tac&0b11]
			tickFlag = true
		}
	}

	g.timer.TIMAReload.after = false
	if g.timer.TIMAReload.flag {
		g.timer.TIMAReload.flag = false
		g.IO[TIMAIO-0xff00] = g.timer.TIMAReload.value
		g.timer.TIMAReload.after = true
		g.setTimerFlag() // ref: https://gbdev.io/pandocs/#timer-overflow-behaviour
	}

	if tickFlag {
		TIMABefore := g.IO[TIMAIO-0xff00]
		TIMAAfter := TIMABefore + 1
		if TIMAAfter < TIMABefore { // overflow occurs
			g.timer.TIMAReload = TIMAReload{
				flag:  true,
				value: uint8(g.IO[TMAIO-0xff00]),
				after: false,
			}
			g.IO[TIMAIO-0xff00] = 0
		} else {
			g.IO[TIMAIO-0xff00] = TIMAAfter
		}
	}

	// OAMDMA
	if g.timer.OAMDMA.ptr > 0 {
		if g.timer.OAMDMA.ptr == 160 {
			g.Store8(0xfe00+uint16(g.timer.OAMDMA.ptr)-1, g.Load8(g.timer.OAMDMA.start+uint16(g.timer.OAMDMA.ptr)-1))
			g.Store8(OAM, 0xff)
		} else if g.timer.OAMDMA.ptr < 160 {
			g.Store8(0xfe00+uint16(g.timer.OAMDMA.ptr)-1, g.Load8(g.timer.OAMDMA.start+uint16(g.timer.OAMDMA.ptr)-1))
		}

		g.timer.OAMDMA.ptr--          // increment timer.OAMDMA count
		if g.timer.OAMDMA.reptr > 0 { // if next OAM is requested, increment that one too
			g.timer.OAMDMA.reptr--
			if g.timer.OAMDMA.reptr == 160 {
				g.timer.OAMDMA.start = g.timer.OAMDMA.restart
				g.timer.OAMDMA.ptr, g.timer.OAMDMA.reptr = 160, 0
			}
		}
	}
}

func (g *GBC) resetTimer() bool {
	g.timer.sys, g.timer.div, g.IO[DIVIO-0xff00] = 0, 0, 0

	old := g.timer.tac
	g.timer.tac = 0

	tickFlag := false
	tac := g.IO[TACIO-0xff00]
	if util.Bit(tac, 2) {
		tickFlag = old >= (clocks[tac&0b11] / 2) // ref: https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/timer/tim00_div_trigger.s
	}
	return tickFlag
}
