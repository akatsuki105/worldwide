package gbc

import "github.com/pokemium/worldwide/pkg/util"

func (g *GBC) push(b byte) {
	g.Store8(g.Reg.SP-1, b)
	g.Reg.SP--
}

func (g *GBC) pop() byte {
	value := g.Load8(g.Reg.SP)
	g.Reg.SP++
	return value
}

func (g *GBC) pushPC() {
	upper, lower := byte(g.Reg.PC>>8), byte(g.Reg.PC)
	g.push(upper)
	g.push(lower)
}

func (g *GBC) pushPCCALL() {
	upper, lower := byte(g.Reg.PC>>8), byte(g.Reg.PC&0x00ff)
	g.push(upper)
	g.timer.tick(1 * 4 >> util.Bool2U32(g.DoubleSpeed)) // M = 4: PC push: memory access for high byte
	g.push(lower)
	g.timer.tick(1 * 4 >> util.Bool2U32(g.DoubleSpeed)) // M = 5: PC push: memory access for low byte
}

func (g *GBC) popPC() {
	lower := uint16(g.pop())
	upper := uint16(g.pop())
	g.Reg.PC = (upper << 8) | lower
}
