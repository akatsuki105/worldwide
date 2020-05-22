package emulator

import (
	"fmt"
)

var (
	maxHistory = 128
)

func (cpu *CPU) debugRegister() string {
	A, F := byte(cpu.Reg.AF>>8), byte(cpu.Reg.AF)
	B, C := byte(cpu.Reg.BC>>8), byte(cpu.Reg.BC)
	D, E := byte(cpu.Reg.DE>>8), byte(cpu.Reg.DE)
	H, L := byte(cpu.Reg.HL>>8), byte(cpu.Reg.HL)
	return fmt.Sprintf(`Register
A: %02x       F: %02x
B: %02x       C: %02x
D: %02x       E: %02x
H: %02x       L: %02x
PC: 0x%04x  SP: 0x%04x`, A, F, B, C, D, E, H, L, cpu.Reg.PC, cpu.Reg.SP)
}

func (cpu *CPU) debugLCD() string {
	LCDC := cpu.FetchMemory8(LCDCIO)
	STAT := cpu.FetchMemory8(LCDSTATIO)
	SCX, SCY := cpu.GPU.GetScroll()
	WY := cpu.FetchMemory8(WYIO)
	WX := cpu.FetchMemory8(WXIO) - 7
	return fmt.Sprintf(`-- LCD --
LCDC: %02x
STAT: %02x
SCX: %02x    SCY: %02x
WX: %02x     WY: %02x`, LCDC, STAT, SCX, SCY, WX, WY)
}

func (cpu *CPU) DebugExec(frame int) {
	for i := 0; i < frame; i++ {
		for y := 0; y < 144; y++ {
			cpu.execScanline()
		}
		cpu.execVBlank()
	}
}
