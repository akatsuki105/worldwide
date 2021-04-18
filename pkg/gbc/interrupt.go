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

func (cpu *CPU) ieif() intrIEIF {
	ieif := intrIEIF{}
	ieio, ifio := cpu.RAM[IEIO], cpu.RAM[IFIO]

	ieif.VBlank.IE, ieif.VBlank.IF = util.Bit(ieio, 0), util.Bit(ifio, 0)
	ieif.LCDSTAT.IE, ieif.LCDSTAT.IF = util.Bit(ieio, 1), util.Bit(ifio, 1)
	ieif.Timer.IE, ieif.Timer.IF = util.Bit(ieio, 2), util.Bit(ifio, 2)
	ieif.Serial.IE, ieif.Serial.IF = util.Bit(ieio, 3), util.Bit(ifio, 3)
	ieif.Joypad.IE, ieif.Joypad.IF = util.Bit(ieio, 4), util.Bit(ifio, 4)

	return ieif
}

// ------------ VBlank --------------------

func (cpu *CPU) setVBlankFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)|0x01)
}

func (cpu *CPU) clearVBlankFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)&0xfe)
}

// ------------ LCD STAT ------------------

func (cpu *CPU) setLCDSTATFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)|0x02)
}

func (cpu *CPU) clearLCDSTATFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)&0xfd)
}

// ------------ timer --------------------
// timer.go

// ------------ Serial --------------------

func (cpu *CPU) setSerialFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)|0x08)
}

func (cpu *CPU) clearSerialFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)&0xf7)
}

// ------------ Joypad --------------------

func (cpu *CPU) getJoypadEnable() bool {
	return util.Bit(cpu.fetchIO(IEIO), 4)
}

func (cpu *CPU) setJoypadFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)|0x10)
}

func (cpu *CPU) clearJoypadFlag() {
	cpu.setIO(IFIO, cpu.fetchIO(IFIO)&0xef)
}

// ------------ trigger --------------------

func (cpu *CPU) triggerInterrupt() {
	cpu.Reg.IME, cpu.halt = false, false
	cpu.timer(5) // https://gbdev.gg8.se/wiki/articles/Interrupts#InterruptServiceRoutine
	cpu.pushPC()
}

func (cpu *CPU) triggerVBlank() {
	if util.Bit(cpu.fetchIO(LCDCIO), 7) {
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
		intr := cpu.ieif()

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
