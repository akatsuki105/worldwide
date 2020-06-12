package emulator

type Cycle struct {
	tac      int // タイマー用
	div      int // DIVタイマー用
	scanline int // スキャンライン用
	serial   int
}

type TIMAReload struct {
	flag  bool
	value byte
}

type Timer struct {
	Cycle
	OAMDMA
	TIMAReload
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
	if cycle == 0 {
		return
	}

	TAC := cpu.RAM[TACIO]
	tickFlag := false

	// DI,EIの遅延処理
	if cpu.IMESwitch.Working {
		for i := 0; i < cycle; i++ {
			cpu.IMESwitch.Count--
			if cpu.IMESwitch.Count == 0 {
				cpu.Reg.IME = cpu.IMESwitch.Value
				cpu.IMESwitch.Working = false
				break
			}
		}
	}

	// シリアル通信のクロック管理
	if cpu.Config.Network.Network && cpu.Serial.TransferFlag > 0 {
		cpu.Cycle.serial += cycle
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
	cpu.Cycle.scanline += cycle

	// DIVレジスタ
	cpu.Cycle.div += cycle
	if cpu.Cycle.div >= 64 {
		cpu.RAM[DIVIO]++
		cpu.Cycle.div -= 64
	}

	if (TAC>>2)&0x01 == 1 {
		cpu.Cycle.tac += cycle
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

	if cpu.TIMAReload.flag {
		cpu.TIMAReload.flag = false
		cpu.RAM[TIMAIO] = cpu.TIMAReload.value
	}

	if tickFlag {
		TIMABefore := cpu.RAM[TIMAIO]
		TIMAAfter := TIMABefore + 1
		if TIMAAfter < TIMABefore {
			// overflow occurs
			cpu.TIMAReload = TIMAReload{
				flag:  true,
				value: uint8(cpu.RAM[TMAIO]),
			}
			cpu.RAM[TIMAIO] = 0
			cpu.setTimerFlag()
		} else {
			cpu.RAM[TIMAIO] = TIMAAfter
		}
	}

	// OAMDMA
	if cpu.OAMDMA.ptr > 0 {
		for i := 0; i < cycle; i++ {
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

			if cpu.OAMDMA.ptr == 0 {
				break
			}
		}
	}
}

func (cpu *CPU) resetTimer() {
	cpu.Cycle.div = 0
	cpu.RAM[DIVIO] = 0

	cpu.Cycle.tac = 0
}
