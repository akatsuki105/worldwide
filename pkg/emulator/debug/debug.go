package debug

import "gbc/pkg/gbc"

type Debugger struct {
	g *gbc.GBC
}

func New(g *gbc.GBC) *Debugger {
	return &Debugger{g}
}
