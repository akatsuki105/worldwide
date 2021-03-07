package emulator

type Cycle struct {
	tac      int // タイマー用
	div      int // DIVタイマー用
	scanline int // スキャンライン用
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
	start   uint16
	ptr     uint16
	restart uint16 // OAMDMA中に再びOAMDMAをリクエストしたとき
	reptr   uint16 // OAMDMA中に再びOAMDMAをリクエストしたとき
}

func (cpu *CPU) getTimerEnable() bool {
	IE := cpu.fetchIO(IEIO)
	TimerEnable := (IE >> 2) % 2
	if TimerEnable == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setTimerEnable() {
	IE := cpu.fetchIO(IEIO) | 0x04
	cpu.setIO(IEIO, IE)
}

func (cpu *CPU) clearTimerEnable() {
	IE := cpu.fetchIO(IEIO) & 0xfb
	cpu.setIO(IEIO, IE)
}

func (cpu *CPU) getTimerFlag() bool {
	IF := cpu.fetchIO(IFIO)
	TimerFlag := (IF >> 2) % 2
	if TimerFlag == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setTimerFlag() {
	IF := cpu.fetchIO(IFIO) | 0x04
	cpu.setIO(IFIO, IF)
	cpu.halt = false
}

func (cpu *CPU) clearTimerFlag() {
	IF := cpu.fetchIO(IFIO) & 0xfb
	cpu.setIO(IFIO, IF)
}

func (cpu *CPU) timer(cycle int) {
	for i := 0; i < cycle; i++ {
		cpu.tick()
		cpu.Sound.Buffer(4, cpu.boost)
	}
}

func (cpu *CPU) tick() {
	TAC := cpu.RAM[TACIO]
	tickFlag := false

	if cpu.Timer.ResetAll {
		cpu.Timer.ResetAll = false
		tickFlag = cpu.resetTimer()
	}
	if cpu.Timer.TAC.Change && !tickFlag {
		cpu.Timer.TAC.Change = false
		clocks := [4]uint16{1024 / 4, 16 / 4, 64 / 4, 256 / 4}
		oldTAC, newTAC := cpu.Timer.TAC.Old, cpu.RAM[TACIO]
		oldClock, newClock := clocks[oldTAC&0b11], clocks[newTAC&0b11]
		oldEnable, newEnable := oldTAC&0b100 > 0, newTAC&0b100 > 0
		if oldEnable {
			if newEnable {
				tickFlag = cpu.Cycle.sys&(oldClock/2) > 0
			} else {
				tickFlag = cpu.Cycle.sys&(oldClock/2) > 0 && cpu.Cycle.sys&(newClock/2) == 0
			}
		}
	}

	// DI,EIの遅延処理
	if cpu.IMESwitch.Working {
		cpu.IMESwitch.Count--
		if cpu.IMESwitch.Count == 0 {
			cpu.Reg.IME = cpu.IMESwitch.Value
			cpu.IMESwitch.Working = false
		}
	}

	// シリアル通信のクロック管理
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

	// スキャンライン
	cpu.Cycle.scanline++

	// 16 bit system counter
	cpu.Cycle.sys++

	// DIVレジスタ
	cpu.Cycle.div++
	if cpu.Cycle.div >= 64 {
		cpu.RAM[DIVIO]++
		cpu.Cycle.div -= 64
	}

	if (TAC>>2)&0x01 == 1 {
		cpu.Cycle.tac++
		switch TAC % 4 {
		case 0:
			// 4096Hz (1024/4 cycle)
			if cpu.Cycle.tac >= 256 {
				cpu.Cycle.tac -= 256
				tickFlag = true
			}
		case 1:
			// 262144Hz (16/4 cycle)
			if cpu.Cycle.tac >= 4 {
				cpu.Cycle.tac -= 4
				tickFlag = true
			}
		case 2:
			// 65536Hz (64/4 cycle)
			if cpu.Cycle.tac >= 16 {
				cpu.Cycle.tac -= 16
				tickFlag = true
			}
		case 3:
			// 16384Hz (256/4 cycle)
			if cpu.Cycle.tac >= 64 {
				cpu.Cycle.tac -= 64
				tickFlag = true
			}
		}
	}

	if cpu.TIMAReload.after {
		cpu.TIMAReload.after = false
	}
	if cpu.TIMAReload.flag {
		cpu.TIMAReload.flag = false
		cpu.RAM[TIMAIO] = cpu.TIMAReload.value
		cpu.TIMAReload.after = true
		cpu.setTimerFlag() // ref: https://gbdev.io/pandocs/#timer-overflow-behaviour
	}

	if tickFlag {
		TIMABefore := cpu.RAM[TIMAIO]
		TIMAAfter := TIMABefore + 1
		if TIMAAfter < TIMABefore {
			// overflow occurs
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

		// OAMDMAを1カウント進める(重複しているときはそっちのカウントも進める)
		cpu.OAMDMA.ptr--
		if cpu.OAMDMA.reptr > 0 {
			cpu.OAMDMA.reptr--

			if cpu.OAMDMA.reptr == 160 {
				cpu.OAMDMA.start = cpu.OAMDMA.restart
				cpu.OAMDMA.ptr = 160
				cpu.OAMDMA.reptr = 0
			}
		}
	}
}

func (cpu *CPU) resetTimer() bool {
	cpu.Cycle.sys = 0
	cpu.Cycle.div = 0
	cpu.RAM[DIVIO] = 0

	old := cpu.Cycle.tac
	cpu.Cycle.tac = 0

	tickFlag := false
	TAC := cpu.RAM[TACIO]
	if (TAC>>2)&0x01 == 1 {
		switch TAC % 4 {
		case 0:
			// 4096Hz (1024/4 cycle)
			// ref: https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/timer/tim00_div_trigger.s
			tickFlag = old >= 512/4
		case 1:
			// 262144Hz (16/4 cycle)
			// ref: https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/timer/tim01_div_trigger.s
			tickFlag = old >= 8/4
		case 2:
			// 65536Hz (64/4 cycle)
			// ref: https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/timer/tim10_div_trigger.s
			tickFlag = old >= 32/4
		case 3:
			// 16384Hz (256/4 cycle)
			// ref: https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/timer/tim11_div_trigger.s
			tickFlag = old >= 128/4
		}
	}
	return tickFlag
}
