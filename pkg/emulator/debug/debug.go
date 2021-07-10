package debug

import "gbc/pkg/gbc"

type Debugger struct {
	Enable bool
	g      *gbc.GBC
	cart   *cartridge
}

func New(enable bool, g *gbc.GBC) *Debugger {
	return &Debugger{
		Enable: enable,
		g:      g,
		cart:   newCart(g),
	}
}
