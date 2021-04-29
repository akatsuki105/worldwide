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

func (cpu *CPU) setTimerFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)|0x04)
	cpu.halt = false
}

func (cpu *CPU) clearTimerFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)&0xfb)
}

func (cpu *CPU) timer(cycle int) {
	for i := 0; i < cycle; i++ {
		cpu.tick()
	}
	cpu.Sound.Buffer(4*cycle, cpu.boost)
}

// 0: 4096Hz (1024/4 cycle), 1: 262144Hz (16/4 cycle), 2: 65536Hz (64/4 cycle), 3: 16384Hz (256/4 cycle)
var clocks = [4]int{1024 / 4, 16 / 4, 64 / 4, 256 / 4}

func (cpu *CPU) tick() {
	tac := cpu.RAM[TACIO]
	tickFlag := false

	if cpu.Timer.ResetAll {
		cpu.Timer.ResetAll = false
		tickFlag = cpu.resetTimer()
	}
	if cpu.Timer.TAC.Change && !tickFlag {
		cpu.Timer.TAC.Change = false
		oldTAC, newTAC := cpu.Timer.TAC.Old, cpu.RAM[TACIO]
		oldClock, newClock := uint16(clocks[oldTAC&0b11]), uint16(clocks[newTAC&0b11])
		oldEnable, newEnable := oldTAC&0b100 > 0, newTAC&0b100 > 0
		if oldEnable {
			if newEnable {
				tickFlag = cpu.Cycle.sys&(oldClock/2) > 0
			} else {
				tickFlag = cpu.Cycle.sys&(oldClock/2) > 0 && cpu.Cycle.sys&(newClock/2) == 0
			}
		}
	}

	// lag occurs in di, ei
	if cpu.IMESwitch.Working {
		cpu.IMESwitch.Count--
		if cpu.IMESwitch.Count == 0 {
			cpu.Reg.IME = cpu.IMESwitch.Value
			cpu.IMESwitch.Working = false
		}
	}

	// clock management in serial communication
	if cpu.Config.Network.Network && cpu.Serial.TransferFlag > 0 {
		cpu.Cycle.serial++
		if cpu.Cycle.serial > 128*8 {
			cpu.Serial.TransferFlag = 0
			close(cpu.serialTick)
			cpu.Cycle.serial = 0
			cpu.serialTick = make(chan int)
		}
	} else {
		cpu.Cycle.serial = 0
	}

	cpu.Cycle.scanline++
	cpu.Cycle.sys++ // 16 bit system counter
	cpu.Cycle.div++
	if cpu.Cycle.div >= 64 {
		cpu.RAM[DIVIO]++
		cpu.Cycle.div -= 64
	}

	if util.Bit(tac, 2) {
		cpu.Cycle.tac++
		if cpu.Cycle.tac >= clocks[tac&0b11] {
			cpu.Cycle.tac -= clocks[tac&0b11]
			tickFlag = true
		}
	}

	cpu.TIMAReload.after = false
	if cpu.TIMAReload.flag {
		cpu.TIMAReload.flag = false
		cpu.RAM[TIMAIO] = cpu.TIMAReload.value
		cpu.TIMAReload.after = true
		cpu.setTimerFlag() // ref: https://gbdev.io/pandocs/#timer-overflow-behaviour
	}

	if tickFlag {
		TIMABefore := cpu.RAM[TIMAIO]
		TIMAAfter := TIMABefore + 1
		if TIMAAfter < TIMABefore { // overflow occurs
			cpu.TIMAReload = TIMAReload{
				flag:  true,
				value: uint8(cpu.RAM[TMAIO]),
				after: false,
			}
			cpu.RAM[TIMAIO] = 0
		} else {
			cpu.RAM[TIMAIO] = TIMAAfter
		}
	}

	// OAMDMA
	if cpu.OAMDMA.ptr > 0 {
		if cpu.OAMDMA.ptr == 160 {
			cpu.RAM[0xfe00+uint16(cpu.OAMDMA.ptr)-1] = cpu.FetchMemory8(cpu.OAMDMA.start + uint16(cpu.OAMDMA.ptr) - 1)
			cpu.RAM[OAM] = 0xff
		} else if cpu.OAMDMA.ptr < 160 {
			cpu.RAM[0xfe00+uint16(cpu.OAMDMA.ptr)-1] = cpu.FetchMemory8(cpu.OAMDMA.start + uint16(cpu.OAMDMA.ptr) - 1)
		}

		cpu.OAMDMA.ptr--          // increment OAMDMA count
		if cpu.OAMDMA.reptr > 0 { // if next OAM is requested, increment that one too
			cpu.OAMDMA.reptr--
			if cpu.OAMDMA.reptr == 160 {
				cpu.OAMDMA.start = cpu.OAMDMA.restart
				cpu.OAMDMA.ptr, cpu.OAMDMA.reptr = 160, 0
			}
		}
	}
}

func (cpu *CPU) resetTimer() bool {
	cpu.Cycle.sys, cpu.Cycle.div, cpu.RAM[DIVIO] = 0, 0, 0

	old := cpu.Cycle.tac
	cpu.Cycle.tac = 0

	tickFlag := false
	tac := cpu.RAM[TACIO]
	if util.Bit(tac, 2) {
		tickFlag = old >= (clocks[tac&0b11] / 2) // ref: https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/timer/tim00_div_trigger.s
	}
	return tickFlag
}
