package debug

import (
	"fmt"
	"gbc/pkg/gbc"
	"gbc/pkg/util"
)

type Debugger struct {
	Enable bool
	g      *gbc.GBC
	cart   *Cartridge
}

func New(enable bool, g *gbc.GBC) *Debugger {
	return &Debugger{
		Enable: enable,
		g:      g,
		cart:   newCart(g),
	}
}

func (d *Debugger) Register() string {
	A, F := d.g.Reg.R[gbc.A], d.g.Reg.R[gbc.F]
	B, C := d.g.Reg.R[gbc.B], d.g.Reg.R[gbc.C]
	D, E := d.g.Reg.R[gbc.D], d.g.Reg.R[gbc.E]
	H, L := d.g.Reg.R[gbc.H], d.g.Reg.R[gbc.L]

	return fmt.Sprintf(`Register
A: %02x       F: %02x
B: %02x       C: %02x
D: %02x       E: %02x
H: %02x       L: %02x
PC: %04x  SP: %04x`, A, F, B, C, D, E, H, L, d.g.Reg.PC, d.g.Reg.SP)
}

func (d *Debugger) IOMap() string {
	LCDC, STAT := d.g.Load8(0xff00+uint16(gbc.LCDCIO)), d.g.Load8(0xff00+uint16(gbc.LCDSTATIO))
	DIV := d.g.Load8(0xff00 + uint16(gbc.DIVIO))
	LY, LYC := d.g.Load8(0xff00+uint16(gbc.LYIO)), d.g.Load8(0xff00+uint16(gbc.LYCIO))
	IE, IF, IME := d.g.Load8(0xff00+uint16(gbc.IEIO)), d.g.Load8(0xff00+uint16(gbc.IFIO)), util.Bool2Int(d.g.Reg.IME)
	spd := util.Bool2Int(d.g.DoubleSpeed)
	return fmt.Sprintf(`IO
LCDC: %02x   STAT: %02x
DIV: %02x
LY: %02x     LYC: %02x
IE: %02x     IF: %02x    IME: %02x
SPD: %02x`, LCDC, STAT, DIV, LY, LYC, IE, IF, IME, spd)
}

func (d *Debugger) Cartridge() string {
	return fmt.Sprintf("%v", d.cart)
}
