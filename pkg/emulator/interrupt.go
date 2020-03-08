package emulator

// IMESwitch - ${Count}サイクル後にIMEを${Value}の値に切り替える ${Working}=falseのときは無効
type IMESwitch struct {
	Count   uint
	Value   bool
	Working bool
}

// ------------ VBlank --------------------

func (cpu *CPU) getVBlankEnable() bool {
	IE := cpu.FetchMemory8(IEIO)
	VBlankEnable := IE % 2
	if VBlankEnable == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setVBlankEnable() {
	IE := cpu.FetchMemory8(IEIO) | 0x01
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) clearVBlankEnable() {
	IE := cpu.FetchMemory8(IEIO) & 0xfe
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) getVBlankFlag() bool {
	IF := cpu.FetchMemory8(IFIO)
	VBlankFlag := IF % 2
	if VBlankFlag == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setVBlankFlag() {
	IF := cpu.FetchMemory8(IFIO) | 0x01
	cpu.SetMemory8(IFIO, IF)
	cpu.halt = false
}

func (cpu *CPU) clearVBlankFlag() {
	IF := cpu.FetchMemory8(IFIO) & 0xfe
	cpu.SetMemory8(IFIO, IF)
}

// ------------ LCD STAT ------------------

func (cpu *CPU) getLCDSTATEnable() bool {
	IE := cpu.FetchMemory8(IEIO)
	LCDSTATEnable := (IE >> 1) % 2
	if LCDSTATEnable == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setLCDSTATEnable() {
	IE := cpu.FetchMemory8(IEIO) | 0x02
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) clearLCDSTATEnable() {
	IE := cpu.FetchMemory8(IEIO) & 0xfd
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) getLCDSTATFlag() bool {
	IF := cpu.FetchMemory8(IFIO)
	LCDSTATFlag := (IF >> 1) % 2
	if LCDSTATFlag == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setLCDSTATFlag() {
	IF := cpu.FetchMemory8(IFIO) | 0x02
	cpu.SetMemory8(IFIO, IF)
	cpu.halt = false
}

func (cpu *CPU) clearLCDSTATFlag() {
	IF := cpu.FetchMemory8(IFIO) & 0xfd
	cpu.SetMemory8(IFIO, IF)
}

// ------------ timer --------------------

func (cpu *CPU) getTimerEnable() bool {
	IE := cpu.FetchMemory8(IEIO)
	TimerEnable := (IE >> 2) % 2
	if TimerEnable == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setTimerEnable() {
	IE := cpu.FetchMemory8(IEIO) | 0x04
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) clearTimerEnable() {
	IE := cpu.FetchMemory8(IEIO) & 0xfb
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) getTimerFlag() bool {
	IF := cpu.FetchMemory8(IFIO)
	LCDSTATFlag := (IF >> 2) % 2
	if LCDSTATFlag == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setTimerFlag() {
	IF := cpu.FetchMemory8(IFIO) | 0x04
	cpu.SetMemory8(IFIO, IF)
	cpu.halt = false
}

func (cpu *CPU) clearTimerFlag() {
	IF := cpu.FetchMemory8(IFIO) & 0xfb
	cpu.SetMemory8(IFIO, IF)
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
	if cpu.network && cpu.Serial.TransferFlag > 0 {
		cpu.cycleSerial += cycle
		if cpu.cycleSerial > 128*8 {
			cpu.Serial.TransferFlag = 0
			close(cpu.serialTick)
			cpu.cycleSerial = 0
			cpu.serialTick = make(chan int)
		}
	} else {
		cpu.cycleSerial = 0
	}

	// スキャンライン
	cpu.cycleLine += cycle

	// DIVレジスタ
	cpu.cycleDIV += cycle
	if cpu.cycleDIV >= 64 {
		cpu.RAM[DIVIO]++
		cpu.cycleDIV -= 64
	}

	if (TAC>>2)%2 == 1 {
		cpu.cycle += cycle
		switch TAC % 4 {
		case 0:
			if cpu.cycle >= 256 {
				cpu.cycle -= 256
				tickFlag = true
			}
		case 1:
			if cpu.cycle >= 4 {
				cpu.cycle -= 4
				tickFlag = true
			}
		case 2:
			if cpu.cycle >= 12 {
				cpu.cycle -= 12
				tickFlag = true
			}
		case 3:
			if cpu.cycle >= 64 {
				cpu.cycle -= 64
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
			}

			if cpu.ptrOAMDMA == 0 {
				if cpu.reptrOAMDMA > 0 {
					cpu.startOAMDMA = cpu.restartOAMDMA
					cpu.ptrOAMDMA = cpu.reptrOAMDMA
					cpu.reptrOAMDMA = 0
				}
				break
			}
		}
	}
}

// ------------ Serial --------------------

func (cpu *CPU) getSerialEnable() bool {
	IE := cpu.FetchMemory8(IEIO)
	SerialEnable := (IE >> 3) % 2
	if SerialEnable == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setSerialEnable() {
	IE := cpu.FetchMemory8(IEIO) | 0x08
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) clearSerialEnable() {
	IE := cpu.FetchMemory8(IEIO) & 0xf7
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) getSerialFlag() bool {
	IF := cpu.FetchMemory8(IFIO)
	serialFlag := (IF >> 3) % 2
	if serialFlag == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setSerialFlag() {
	IF := cpu.FetchMemory8(IFIO) | 0x08
	cpu.SetMemory8(IFIO, IF)
	cpu.halt = false
}

func (cpu *CPU) clearSerialFlag() {
	IF := cpu.FetchMemory8(IFIO) & 0xf7
	cpu.SetMemory8(IFIO, IF)
}

// ------------ Joypad --------------------

func (cpu *CPU) getJoypadEnable() bool {
	IE := cpu.FetchMemory8(IEIO)
	JoypadEnable := (IE >> 4) % 2
	if JoypadEnable == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setJoypadEnable() {
	IE := cpu.FetchMemory8(IEIO) | 0x10
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) clearJoypadEnable() {
	IE := cpu.FetchMemory8(IEIO) & 0xef
	cpu.SetMemory8(IEIO, IE)
}

func (cpu *CPU) getJoypadFlag() bool {
	IF := cpu.FetchMemory8(IFIO)
	JoypadFlag := (IF >> 4) % 2
	if JoypadFlag == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setJoypadFlag() {
	IF := cpu.FetchMemory8(IFIO) | 0x10
	cpu.SetMemory8(IFIO, IF)
	cpu.halt = false
}

func (cpu *CPU) clearJoypadFlag() {
	IF := cpu.FetchMemory8(IFIO) & 0xef
	cpu.SetMemory8(IFIO, IF)
}

// ------------ trigger --------------------

func (cpu *CPU) triggerInterrupt() {
	cpu.Reg.IME = false
	cpu.halt = false
	cpu.timer(5) // https://gbdev.gg8.se/wiki/articles/Interrupts#InterruptServiceRoutine
	cpu.pushPC()
}

func (cpu *CPU) triggerVBlank() {
	LCDActive := (cpu.FetchMemory8(LCDCIO) >> 7) == 1
	if LCDActive {
		cpu.clearVBlankFlag()
		cpu.triggerInterrupt()
		cpu.Reg.PC = 0x0040
	}
}

func (cpu *CPU) triggerLCDC() {
	cpu.clearLCDSTATFlag()
	cpu.triggerInterrupt()
	cpu.Reg.PC = 0x0048
}

func (cpu *CPU) triggerTimer() {
	cpu.clearTimerFlag()
	cpu.triggerInterrupt()
	cpu.Reg.PC = 0x0050
}

func (cpu *CPU) triggerSerial() {
	cpu.clearSerialFlag()
	cpu.triggerInterrupt()
	cpu.Reg.PC = 0x0058
}

func (cpu *CPU) triggerJoypad() {
	cpu.clearJoypadFlag()
	cpu.triggerInterrupt()
	cpu.Reg.PC = 0x0060
}

// ------------ handler --------------------

// 能動的な割り込みに対処する
func (cpu *CPU) handleInterrupt() {

	if cpu.Reg.IME {
		if cpu.getVBlankEnable() && cpu.getVBlankFlag() {
			cpu.triggerVBlank()
		}

		if cpu.getLCDSTATEnable() && cpu.getLCDSTATFlag() {
			cpu.triggerLCDC()
		}

		if cpu.getTimerEnable() && cpu.getTimerFlag() {
			cpu.triggerTimer()
		}

		if cpu.getSerialEnable() && cpu.getSerialFlag() {
			cpu.triggerSerial()
		}

		if cpu.getJoypadEnable() && cpu.getJoypadFlag() {
			cpu.triggerJoypad()
		}
	}
}
