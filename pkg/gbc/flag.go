package gbc

const (
	flagZ, flagN, flagH, flagC = 7, 6, 5, 4
)

func (g *GBC) setCSub(dst, src byte) {
	g.setF(flagC, dst < uint8(dst-src))
}

func (g *GBC) f(idx int) bool {
	return g.Reg.R[F]&(1<<idx) != 0
}

func (g *GBC) setF(idx int, flag bool) {
	if flag {
		g.Reg.R[F] |= (1 << idx)
		return
	}
	g.Reg.R[F] &= ^(1 << idx)
}
