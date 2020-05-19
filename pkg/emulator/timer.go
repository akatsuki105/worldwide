package emulator

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
	LCDSTATFlag := (IF >> 2) % 2
	if LCDSTATFlag == 1 {
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
		cpu.cycle.serial += cycle
		if cpu.cycle.serial > 128*8 {
			cpu.Serial.TransferFlag = 0
			close(cpu.serialTick)
			cpu.cycle.serial = 0
			cpu.serialTick = make(chan int)
		}
	} else {
		cpu.cycle.serial = 0
	}

	// スキャンライン
	cpu.cycle.scanline += cycle

	// DIVレジスタ
	cpu.cycle.div += cycle
	if cpu.cycle.div >= 64 {
		cpu.RAM[DIVIO]++
		cpu.cycle.div -= 64
	}

	if (TAC>>2)%2 == 1 {
		cpu.cycle.tac += cycle
		switch TAC % 4 {
		case 0:
			if cpu.cycle.tac >= 256 {
				cpu.cycle.tac -= 256
				tickFlag = true
			}
		case 1:
			if cpu.cycle.tac >= 4 {
				cpu.cycle.tac -= 4
				tickFlag = true
			}
		case 2:
			if cpu.cycle.tac >= 12 {
				cpu.cycle.tac -= 12
				tickFlag = true
			}
		case 3:
			if cpu.cycle.tac >= 64 {
				cpu.cycle.tac -= 64
				tickFlag = true
			}
		}
	}

	if tickFlag {
		TIMABefore := cpu.RAM[TIMAIO]
		TIMAAfter := TIMABefore + 1
		if TIMAAfter < TIMABefore {
			// オーバーフローしたとき
			TIMAAfter = uint8(cpu.RAM[TMAIO])
			cpu.RAM[TIMAIO] = TIMAAfter
			cpu.setTimerFlag()
		} else {
			cpu.RAM[TIMAIO] = TIMAAfter
		}
	}

	// OAMDMA
	if cpu.ptrOAMDMA > 0 {
		for i := 0; i < cycle; i++ {
			if cpu.ptrOAMDMA == 160 {
				cpu.RAM[0xfe00+uint16(cpu.ptrOAMDMA)-1] = cpu.FetchMemory8(cpu.startOAMDMA + uint16(cpu.ptrOAMDMA) - 1)
				cpu.RAM[OAM] = 0xff
			} else if cpu.ptrOAMDMA < 160 {
				cpu.RAM[0xfe00+uint16(cpu.ptrOAMDMA)-1] = cpu.FetchMemory8(cpu.startOAMDMA + uint16(cpu.ptrOAMDMA) - 1)
			}

			// OAMDMAを1カウント進める(重複しているときはそっちのカウントも進める)
			cpu.ptrOAMDMA--
			if cpu.reptrOAMDMA > 0 {
				cpu.reptrOAMDMA--

				if cpu.reptrOAMDMA == 160 {
					cpu.startOAMDMA = cpu.restartOAMDMA
					cpu.ptrOAMDMA = 160
					cpu.reptrOAMDMA = 0
				}
			}

			if cpu.ptrOAMDMA == 0 {
				break
			}
		}
	}
}
