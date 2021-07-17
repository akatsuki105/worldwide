package debug

import (
	"github.com/pokemium/worldwide/pkg/gbc"
)

const (
	TILE_PER_ROW = 16
)

type Debugger struct {
	g *gbc.GBC
}

func New(g *gbc.GBC) *Debugger {
	return &Debugger{
		g: g,
	}
}
