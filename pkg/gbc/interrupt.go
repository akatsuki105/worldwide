package gbc

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

func (cpu *CPU) getIEIF() intrIEIF {
	ieif := intrIEIF{}
	IE, IF := cpu.RAM[IEIO], cpu.RAM[IFIO]

	VBlankEnable, VBlankFlag := IE&0x01 == 1, IF&0x01 == 1
	ieif.VBlank.IE, ieif.VBlank.IF = VBlankEnable, VBlankFlag

	LCDSTATEnable, LCDSTATFlag := (IE>>1)&0x01 == 1, (IF>>1)&0x01 == 1
	ieif.LCDSTAT.IE, ieif.LCDSTAT.IF = LCDSTATEnable, LCDSTATFlag

	TimerEnable, TimerFlag := (IE>>2)&0x01 == 1, (IF>>2)&0x01 == 1
	ieif.Timer.IE, ieif.Timer.IF = TimerEnable, TimerFlag

	SerialEnable, SerialFlag := (IE>>3)&0x01 == 1, (IF>>3)&0x01 == 1
	ieif.Serial.IE, ieif.Serial.IF = SerialEnable, SerialFlag

	JoypadEnable, JoypadFlag := (IE>>4)&0x01 == 1, (IF>>4)&0x01 == 1
	ieif.Joypad.IE, ieif.Joypad.IF = JoypadEnable, JoypadFlag

	return ieif
}

// ------------ VBlank --------------------

func (cpu *CPU) setVBlankFlag() {
	IF := cpu.fetchIO(IFIO) | 0x01
	cpu.setIO(IFIO, IF)
}

func (cpu *CPU) clearVBlankFlag() {
	IF := cpu.fetchIO(IFIO) & 0xfe
	cpu.setIO(IFIO, IF)
}

// ------------ LCD STAT ------------------

func (cpu *CPU) setLCDSTATFlag() {
	IF := cpu.fetchIO(IFIO) | 0x02
	cpu.setIO(IFIO, IF)
}

func (cpu *CPU) clearLCDSTATFlag() {
	IF := cpu.fetchIO(IFIO) & 0xfd
	cpu.setIO(IFIO, IF)
}

// ------------ timer --------------------
// timer.go

// ------------ Serial --------------------

func (cpu *CPU) setSerialFlag() {
	IF := cpu.fetchIO(IFIO) | 0x08
	cpu.setIO(IFIO, IF)
}

func (cpu *CPU) clearSerialFlag() {
	IF := cpu.fetchIO(IFIO) & 0xf7
	cpu.setIO(IFIO, IF)
}

// ------------ Joypad --------------------

func (cpu *CPU) getJoypadEnable() bool {
	IE := cpu.fetchIO(IEIO)
	return (IE>>4)%2 == 1
}

func (cpu *CPU) setJoypadFlag() {
	IF := cpu.fetchIO(IFIO) | 0x10
	cpu.setIO(IFIO, IF)
}

func (cpu *CPU) clearJoypadFlag() {
	IF := cpu.fetchIO(IFIO) & 0xef
	cpu.setIO(IFIO, IF)
}

// ------------ trigger --------------------

func (cpu *CPU) triggerInterrupt() {
	cpu.Reg.IME = false
	cpu.halt = false
	cpu.timer(5) // https://gbdev.gg8.se/wiki/articles/Interrupts#InterruptServiceRoutine
	cpu.pushPC()
}

func (cpu *CPU) triggerVBlank() {
	LCDActive := (cpu.fetchIO(LCDCIO) >> 7) == 1
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
		intr := cpu.getIEIF()

		if intr.VBlank.IE && intr.VBlank.IF {
			cpu.triggerVBlank()
			return
		}

		if intr.LCDSTAT.IE && intr.LCDSTAT.IF {
			cpu.triggerLCDC()
			return
		}

		if intr.Timer.IE && intr.Timer.IF {
			cpu.triggerTimer()
			return
		}

		if intr.Serial.IE && intr.Serial.IF {
			cpu.triggerSerial()
			return
		}

		if intr.Joypad.IE && intr.Joypad.IF {
			cpu.triggerJoypad()
			return
		}
	}
}
