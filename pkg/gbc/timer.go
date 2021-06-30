package gbc

import "gbc/pkg/util"

type Cycle struct {
	tac      int // use in normal timer
	div      int // use in div timer
	scanline int // use in scanline counter
	serial   int
	sys      uint16 // 16 bit system counter. ref: https://gbdev.io/pandocs/Timer_Obscure_Behaviour.html
}

type TIMAReload struct {
	flag  bool
	value byte
	after bool // ref: [B] in https://gbdev.io/pandocs/#timer-overflow-behaviour
}

type Timer struct {
	Cycle
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
	g.setIO(IFIO, g.fetchIO(IFIO)|0x04)
	g.halt = false
}

func (g *GBC) clearTimerFlag() {
	g.setIO(IFIO, g.fetchIO(IFIO)&0xfb)
}

func (g *GBC) timer(cycle int) {
	for i := 0; i < cycle; i++ {
		g.tick()
	}
	g.Sound.Buffer(4*cycle, g.boost)
}

// 0: 4096Hz (1024/4 cycle), 1: 262144Hz (16/4 cycle), 2: 65536Hz (64/4 cycle), 3: 16384Hz (256/4 cycle)
var clocks = [4]int{1024 / 4, 16 / 4, 64 / 4, 256 / 4}

func (g *GBC) tick() {
	tac := g.RAM[TACIO]
	tickFlag := false

	if g.Timer.ResetAll {
		g.Timer.ResetAll = false
		tickFlag = g.resetTimer()
	}
	if g.Timer.TAC.Change && !tickFlag {
		g.Timer.TAC.Change = false
		oldTAC, newTAC := g.Timer.TAC.Old, g.RAM[TACIO]
		oldClock, newClock := uint16(clocks[oldTAC&0b11]), uint16(clocks[newTAC&0b11])
		oldEnable, newEnable := oldTAC&0b100 > 0, newTAC&0b100 > 0
		if oldEnable {
			if newEnable {
				tickFlag = g.Cycle.sys&(oldClock/2) > 0
			} else {
				tickFlag = g.Cycle.sys&(oldClock/2) > 0 && g.Cycle.sys&(newClock/2) == 0
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

	// clock management in serial communication
	if g.Config.Network.Network && g.Serial.TransferFlag > 0 {
		g.Cycle.serial++
		if g.Cycle.serial > 128*8 {
			g.Serial.TransferFlag = 0
			close(g.serialTick)
			g.Cycle.serial = 0
			g.serialTick = make(chan int)
		}
	} else {
		g.Cycle.serial = 0
	}

	g.Cycle.scanline++
	g.Cycle.sys++ // 16 bit system counter
	g.Cycle.div++
	if g.Cycle.div >= 64 {
		g.RAM[DIVIO]++
		g.Cycle.div -= 64
	}

	if util.Bit(tac, 2) {
		g.Cycle.tac++
		if g.Cycle.tac >= clocks[tac&0b11] {
			g.Cycle.tac -= clocks[tac&0b11]
			tickFlag = true
		}
	}

	g.TIMAReload.after = false
	if g.TIMAReload.flag {
		g.TIMAReload.flag = false
		g.RAM[TIMAIO] = g.TIMAReload.value
		g.TIMAReload.after = true
		g.setTimerFlag() // ref: https://gbdev.io/pandocs/#timer-overflow-behaviour
	}

	if tickFlag {
		TIMABefore := g.RAM[TIMAIO]
		TIMAAfter := TIMABefore + 1
		if TIMAAfter < TIMABefore { // overflow occurs
			g.TIMAReload = TIMAReload{
				flag:  true,
				value: uint8(g.RAM[TMAIO]),
				after: false,
			}
			g.RAM[TIMAIO] = 0
		} else {
			g.RAM[TIMAIO] = TIMAAfter
		}
	}

	// OAMDMA
	if g.OAMDMA.ptr > 0 {
		if g.OAMDMA.ptr == 160 {
			g.RAM[0xfe00+uint16(g.OAMDMA.ptr)-1] = g.FetchMemory8(g.OAMDMA.start + uint16(g.OAMDMA.ptr) - 1)
			g.RAM[OAM] = 0xff
		} else if g.OAMDMA.ptr < 160 {
			g.RAM[0xfe00+uint16(g.OAMDMA.ptr)-1] = g.FetchMemory8(g.OAMDMA.start + uint16(g.OAMDMA.ptr) - 1)
		}

		g.OAMDMA.ptr--          // increment OAMDMA count
		if g.OAMDMA.reptr > 0 { // if next OAM is requested, increment that one too
			g.OAMDMA.reptr--
			if g.OAMDMA.reptr == 160 {
				g.OAMDMA.start = g.OAMDMA.restart
				g.OAMDMA.ptr, g.OAMDMA.reptr = 160, 0
			}
		}
	}
}

func (g *GBC) resetTimer() bool {
	g.Cycle.sys, g.Cycle.div, g.RAM[DIVIO] = 0, 0, 0

	old := g.Cycle.tac
	g.Cycle.tac = 0

	tickFlag := false
	tac := g.RAM[TACIO]
	if util.Bit(tac, 2) {
		tickFlag = old >= (clocks[tac&0b11] / 2) // ref: https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/timer/tim00_div_trigger.s
	}
	return tickFlag
}
