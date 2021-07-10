package debug

import "gbc/pkg/gbc"

type Debugger struct {
	g    *gbc.GBC
	cart *cartridge
}

func New(g *gbc.GBC) *Debugger {
	return &Debugger{
		g:    g,
		cart: newCart(g),
	}
}
