package video

import (
	"gbc/pkg/util"
	"image"
	"image/color"
)

// DMG -> Palette[i] = 0~3
//
// CGB -> Palette[i] = Bit0-4(R) | Bit5-9(G) | Bit10-14(B)
type Color uint16

var defaultDmgPalette = [12]uint16{0x7fff, 0x56b5, 0x294a, 0x0000, 0x7fff, 0x56b5, 0x294a, 0x0000, 0x7fff, 0x56b5, 0x294a, 0x0000}

type VRAM struct {
	Bank   uint16       // 0 or 1
	Buffer [0x4000]byte // (0x8000-0x9fff)x2 (using bank on CGB)
}

// Video processes graphics
type Video struct {
	LCDC byte // LCD Control
	VRAM
	io              *[0x100]byte
	HBlankDMALength int
	Debug

	X, Ly int
	Stat  byte // LCD Status

	Renderer *Renderer
	Oam      *OAM

	// 0xff68
	BcpIndex, BcpIncrement int

	// 0xff6a
	OcpIndex, OcpIncrement int

	dmgPalette [12]uint16
	Palette    [64]Color

	FrameCounter, frameskip, frameskipCounter int
	updateIRQs                                func()

	NextLength int
}

var (
	// colors {R, G, B}
	DmgColor [4][3]uint8 = [4][3]uint8{
		{175, 197, 160}, {93, 147, 66}, {22, 63, 48}, {0, 40, 0},
	}
)

const (
	BGP = iota
	OBP0
	OBP1
)

func New(io *[0x100]byte, updateIRQs func()) *Video {
	g := &Video{
		io:         io,
		Oam:        NewOAM(),
		dmgPalette: defaultDmgPalette,
		updateIRQs: updateIRQs,
	}

	g.Debug.On = false
	g.Renderer = NewRenderer(g)
	g.Reset()
	return g
}

func (g *Video) Reset() {
	g.Ly, g.X = 0, 0
	g.Stat = 1
	g.FrameCounter, g.frameskipCounter = 0, 0

	g.SwitchBank(0)
	for i := 0; i < len(g.VRAM.Buffer); i++ {
		g.VRAM.Buffer[i] = 0
	}

	g.Palette[0] = Color(g.dmgPalette[0])
	g.Palette[1] = Color(g.dmgPalette[1])
	g.Palette[2] = Color(g.dmgPalette[2])
	g.Palette[3] = Color(g.dmgPalette[3])
	g.Palette[8*4+0] = Color(g.dmgPalette[4])
	g.Palette[8*4+1] = Color(g.dmgPalette[5])
	g.Palette[8*4+2] = Color(g.dmgPalette[6])
	g.Palette[8*4+3] = Color(g.dmgPalette[7])
	g.Palette[9*4+0] = Color(g.dmgPalette[8])
	g.Palette[9*4+1] = Color(g.dmgPalette[9])
	g.Palette[9*4+2] = Color(g.dmgPalette[10])
	g.Palette[9*4+3] = Color(g.dmgPalette[11])

	g.Renderer.writePalette(0, g.Palette[0])
	g.Renderer.writePalette(1, g.Palette[1])
	g.Renderer.writePalette(2, g.Palette[2])
	g.Renderer.writePalette(3, g.Palette[3])
	g.Renderer.writePalette(8*4+0, g.Palette[8*4+0])
	g.Renderer.writePalette(8*4+1, g.Palette[8*4+1])
	g.Renderer.writePalette(8*4+2, g.Palette[8*4+2])
	g.Renderer.writePalette(8*4+3, g.Palette[8*4+3])
	g.Renderer.writePalette(9*4+0, g.Palette[9*4+0])
	g.Renderer.writePalette(9*4+1, g.Palette[9*4+1])
	g.Renderer.writePalette(9*4+2, g.Palette[9*4+2])
	g.Renderer.writePalette(9*4+3, g.Palette[9*4+3])
}

// Display returns gameboy display data
func (g *Video) Display() *image.RGBA {
	i := image.NewRGBA(image.Rect(0, 0, HORIZONTAL_PIXELS, VERTICAL_PIXELS))
	for y := 0; y < VERTICAL_PIXELS; y++ {
		for x := 0; x < HORIZONTAL_PIXELS; x++ {
			p := g.Renderer.outputBuffer[y*HORIZONTAL_PIXELS+x]
			red, green, blue := byte((p&0b11111)*8), byte(((p>>5)&0b11111)*8), byte(((p>>10)&0b11111)*8)

			i.SetRGBA(x, y, color.RGBA{red, green, blue, 0xff})
		}
	}
	return i
}

// GBVideoWritePalette
// 0xff47, 0xff48, 0xff49, 0xff69, 0xff6b
func (g *Video) WritePalette(offset byte, value byte) {
	if g.Renderer.Model < util.GB_MODEL_SGB {
		switch offset {
		case GB_REG_BGP:
			// Palette = 0(white) or 1(light gray) or 2(dark gray) or 3(black)
			g.Palette[0] = Color(g.dmgPalette[value&3])
			g.Palette[1] = Color(g.dmgPalette[(value>>2)&3])
			g.Palette[2] = Color(g.dmgPalette[(value>>4)&3])
			g.Palette[3] = Color(g.dmgPalette[(value>>6)&3])
			g.Renderer.writePalette(0, g.Palette[0])
			g.Renderer.writePalette(1, g.Palette[1])
			g.Renderer.writePalette(2, g.Palette[2])
			g.Renderer.writePalette(3, g.Palette[3])
		case GB_REG_OBP0:
			g.Palette[8*4+0] = Color(g.dmgPalette[(value&3)+4])
			g.Palette[8*4+1] = Color(g.dmgPalette[((value>>2)&3)+4])
			g.Palette[8*4+2] = Color(g.dmgPalette[((value>>4)&3)+4])
			g.Palette[8*4+3] = Color(g.dmgPalette[((value>>6)&3)+4])
			g.Renderer.writePalette(8*4+0, g.Palette[8*4+0])
			g.Renderer.writePalette(8*4+1, g.Palette[8*4+1])
			g.Renderer.writePalette(8*4+2, g.Palette[8*4+2])
			g.Renderer.writePalette(8*4+3, g.Palette[8*4+3])
		case GB_REG_OBP1:
			g.Palette[9*4+0] = Color(g.dmgPalette[(value&3)+8])
			g.Palette[9*4+1] = Color(g.dmgPalette[((value>>2)&3)+8])
			g.Palette[9*4+2] = Color(g.dmgPalette[((value>>4)&3)+8])
			g.Palette[9*4+3] = Color(g.dmgPalette[((value>>6)&3)+8])
			g.Renderer.writePalette(9*4+0, g.Palette[9*4+0])
			g.Renderer.writePalette(9*4+1, g.Palette[9*4+1])
			g.Renderer.writePalette(9*4+2, g.Palette[9*4+2])
			g.Renderer.writePalette(9*4+3, g.Palette[9*4+3])
		}
	} else if g.Renderer.Model&util.GB_MODEL_SGB != 0 {
		g.Renderer.WriteVideoRegister(offset&0xff, value)
	} else {
		switch offset {
		// gameboy color
		case GB_REG_BCPD:
			if g.Mode() != 3 {
				if g.BcpIndex&1 == 1 {
					// update upper
					g.Palette[g.BcpIndex>>1] &= 0x00FF
					g.Palette[g.BcpIndex>>1] |= Color(uint16(value) << 8)
				} else {
					// update lower
					g.Palette[g.BcpIndex>>1] &= 0xFF00
					g.Palette[g.BcpIndex>>1] |= Color(value)
				}
				g.Renderer.writePalette(g.BcpIndex>>1, g.Palette[g.BcpIndex>>1])
			}
			if g.BcpIncrement != 0 {
				g.BcpIndex++
				g.BcpIndex &= 0x3F
				g.io[GB_REG_BCPS] &= 0x80
				g.io[GB_REG_BCPS] |= byte(g.BcpIndex)
			}
			g.io[GB_REG_BCPD] = byte(g.Palette[g.BcpIndex>>1] >> (8 * (g.BcpIndex & 1)))
		case GB_REG_OCPD:
			if g.Mode() != 3 {
				if g.OcpIndex&1 == 1 {
					g.Palette[8*4+(g.OcpIndex>>1)] &= 0x00FF
					g.Palette[8*4+(g.OcpIndex>>1)] |= Color(uint16(value) << 8)
				} else {
					g.Palette[8*4+(g.OcpIndex>>1)] &= 0xFF00
					g.Palette[8*4+(g.OcpIndex>>1)] |= Color(value)
				}
				g.Renderer.writePalette(8*4+(g.OcpIndex>>1), g.Palette[8*4+(g.OcpIndex>>1)])
			}
			if g.OcpIncrement != 0 {
				g.OcpIndex++
				g.OcpIndex &= 0x3F
				g.io[GB_REG_OCPS] &= 0x80
				g.io[GB_REG_OCPS] |= byte(g.OcpIndex)
			}
			g.io[GB_REG_OCPD] = byte(g.Palette[8*4+(g.OcpIndex>>1)] >> (8 * (g.OcpIndex & 1)))
		}
	}
}

// GBVideoSwitchBank
func (g *Video) SwitchBank(value byte) {
	value &= 1
	g.VRAM.Bank = uint16(value)
}

// GBVideoProcessDots
func (g *Video) ProcessDots(cyclesLate uint32) {
	if g.Mode() != 3 {
		return
	}

	oldX := 0
	g.X = HORIZONTAL_PIXELS
	g.Renderer.drawRange(oldX, g.X, g.Ly)
}

func (g *Video) EndMode() {
	switch g.Mode() {
	case 0:
		g.endMode0()
	case 1:
		g.endMode1()
	case 2:
		g.endMode2()
	case 3:
		g.endMode3()
	}
}

// mode0 = HBlank
// 204 cycles
func (g *Video) endMode0() {
	if g.frameskipCounter <= 0 {
		g.Renderer.finishScanline(g.Ly)
	}

	lyc := g.io[GB_REG_LYC]
	g.Ly++
	g.io[GB_REG_LY] = byte(g.Ly)

	oldStat := g.Stat
	if g.Ly < VERTICAL_PIXELS {
		g.setMode(2)
	} else {
		g.setMode(1)
		if !statIRQAsserted(oldStat) && statIRQAsserted(g.Stat) {
			g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 1, true)
		}
		g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 0, true)
	}

	if !statIRQAsserted(oldStat) && statIRQAsserted(g.Stat) {
		g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 1, true)
	}

	// LYC stat is delayed 1 T-cycle
	oldStat = g.Stat
	g.Stat = util.SetBit8(g.Stat, 2, lyc == g.io[GB_REG_LY])
	if !statIRQAsserted(oldStat) && statIRQAsserted(g.Stat) {
		g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 1, true)
	}

	g.updateIRQs()
}

// mode1 = VBlank
func (g *Video) endMode1() {
	if !util.Bit(g.LCDC, Enable) {
		return
	}

	lyc := g.io[GB_REG_LYC]
	g.Ly++
	switch g.Ly {
	case VERTICAL_TOTAL_PIXELS + 1:
		g.Ly = 0
		g.io[GB_REG_LY] = byte(g.Ly)
		g.setMode(2)
	case VERTICAL_TOTAL_PIXELS:
		g.io[GB_REG_LY] = 0
		g.NextLength = HORIZONTAL_LENGTH - 8
	case VERTICAL_TOTAL_PIXELS - 1:
		g.io[GB_REG_LY] = byte(g.Ly)
		g.NextLength = 8
	default:
		g.io[GB_REG_LY] = byte(g.Ly)
		g.NextLength = HORIZONTAL_LENGTH
	}

	oldStat := g.Stat
	g.Stat = util.SetBit8(g.Stat, 2, lyc == g.io[GB_REG_LY])
	if !statIRQAsserted(oldStat) && statIRQAsserted(g.Stat) {
		g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 1, true)
		g.updateIRQs()
	}
}

// mode2 = [mode0 -> mode2 -> mode3] -> [mode0 -> mode2 -> mode3] -> ...
// 80 cycles
func (g *Video) endMode2() {
	oldStat := g.Stat
	g.X = -(int(g.Renderer.scx) & 7)
	g.setMode(3)
	if !statIRQAsserted(oldStat) && statIRQAsserted(g.Stat) {
		g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 1, true)
		g.updateIRQs()
	}
}

// mode3 = [mode0 -> mode2 -> mode3] -> [mode0 -> mode2 -> mode3] -> ...
// 172 cycles
func (g *Video) endMode3() {
	oldStat := g.Stat
	g.ProcessDots(0)
	g.setMode(0)
	if !statIRQAsserted(oldStat) && statIRQAsserted(g.Stat) {
		g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 1, true)
		g.updateIRQs()
	}
}

func (g *Video) UpdateFrameCount() {
	g.frameskipCounter--
	if g.frameskipCounter < 0 {
		g.Renderer.finishFrame()
		g.frameskipCounter = g.frameskip
	}
	g.FrameCounter++
}

func (g *Video) Mode() byte {
	return g.Stat & 0x3
}

func (g *Video) setMode(mode byte) {
	g.Stat = (g.Stat & 0xfc) | mode
}

// GBVideoWriteLCDC
func (g *Video) WriteLCDC(value byte) {
	if !util.Bit(g.LCDC, Enable) && util.Bit(value, Enable) {
		g.setMode(2)
		g.Ly = 0
		g.io[GB_REG_LY] = 0
		oldStat := g.Stat
		g.setMode(0)
		g.Stat = util.SetBit8(g.Stat, 2, byte(g.Ly) == g.io[GB_REG_LYC])
		if !statIRQAsserted(oldStat) && statIRQAsserted(g.Stat) {
			g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 1, true)
			g.updateIRQs()
		}
		g.Renderer.writePalette(0, g.Palette[0])
	}
	if util.Bit(g.LCDC, Enable) && !util.Bit(value, Enable) {
		g.setMode(0)
		g.Ly = 0
		g.io[GB_REG_LY] = 0
		g.Renderer.writePalette(0, Color(g.dmgPalette[0]))
	}
}

// GBVideoWriteSTAT
func (g *Video) WriteSTAT(value byte) {
	oldStat := g.Stat
	g.Stat = (g.Stat & 0x7) | (value & 0x78)
	if !util.Bit(g.LCDC, Enable) || g.Renderer.Model >= util.GB_MODEL_CGB {
		return
	}
	if !statIRQAsserted(oldStat) && g.Mode() < 3 {
		g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 1, true)
		g.updateIRQs()
	}
}

// GBVideoWriteLYC
func (g *Video) WriteLYC(value byte) {
	oldStat := g.Stat
	if util.Bit(g.LCDC, Enable) {
		g.Stat = util.SetBit8(g.Stat, 2, value == byte(g.Ly))
		if !statIRQAsserted(oldStat) && statIRQAsserted(g.Stat) {
			g.io[GB_REG_IF] = util.SetBit8(g.io[GB_REG_IF], 1, true)
			g.updateIRQs()
		}
	}
}

func statIRQAsserted(stat byte) bool {
	if util.Bit(stat, 6) && util.Bit(stat, 2) {
		return true
	}
	switch stat & 0x3 {
	case 0:
		return util.Bit(stat, 3)
	case 1:
		return util.Bit(stat, 4)
	case 2:
		return util.Bit(stat, 5)
	}
	return false
}
