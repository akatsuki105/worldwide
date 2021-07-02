package gpu

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

// GPU Graphic Processor Unit
type GPU struct {
	LCDC byte // LCD Control
	VRAM
	HBlankDMALength int
	Debug

	X, Ly int
	Stat  byte // LCD Status

	Renderer *Renderer
	oam      OAM

	// 0xff68
	BcpIndex, BcpIncrement int

	// 0xff6a
	OcpIndex, OcpIncrement int

	dmgPalette [12]uint16
	Palette    [64]Color

	frameCounter, frameskip, frameskipCounter int
}

var (
	// colors {R, G, B}
	colors [4][3]uint8 = [4][3]uint8{
		{175, 197, 160}, {93, 147, 66}, {22, 63, 48}, {0, 40, 0},
	}
)

const (
	BGP = iota
	OBP0
	OBP1
)

func New(debug bool) *GPU {
	g := &GPU{}

	g.Debug.On = debug
	if debug {
		g.initTileData()
		g.OAM = image.NewRGBA(image.Rect(0, 0, 16*8-1, 20*5-3))
	}
	g.Renderer = NewRenderer(g)
	g.dmgPalette = defaultDmgPalette
	return g
}

// Display returns gameboy display data
func (g *GPU) Display() *image.RGBA {
	i := image.NewRGBA(image.Rect(0, 0, 160, 144))
	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			p := g.Renderer.outputBuffer[y*160+x]
			r, g, b := byte((p&0b11111)*8), byte(((p>>5)&0b11111)*8), byte(((p>>10)&0b11111)*8)
			i.SetRGBA(x, y, color.RGBA{r, g, b, 0xff})
		}
	}
	return i
}

func (g *GPU) fetchTileBaseAddr() uint16 {
	if util.Bit(g.LCDC, 4) {
		return 0x8000
	}
	return 0x8800
}

// GBVideoWritePalette
// 0xff47, 0xff48, 0xff49, 0xff69, 0xff6b
func (g *GPU) WritePalette(address uint16, value byte) {
	if g.Renderer.model < util.GB_MODEL_SGB {
		switch address {
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
	} else if g.Renderer.model&util.GB_MODEL_SGB != 0 {
		g.Renderer.WriteVideoRegister(address&0xff, value)
	} else {
		switch address {
		// gameboy color
		case GB_REG_BCPD:
			if g.Mode() != 3 {
				if g.BcpIndex&1 == 1 {
					g.Palette[g.BcpIndex>>1] &= 0x00FF
					g.Palette[g.BcpIndex>>1] |= Color(value) << 8
				} else {
					g.Palette[g.BcpIndex>>1] &= 0xFF00
					g.Palette[g.BcpIndex>>1] |= Color(value)
				}
				g.Renderer.writePalette(g.BcpIndex>>1, g.Palette[g.BcpIndex>>1])
			}
			if g.BcpIncrement != 0 {
				g.BcpIndex++
				g.BcpIndex &= 0x3F
				// video->p->memory.io[GB_REG_BCPS] &= 0x80;
				// video->p->memory.io[GB_REG_BCPS] |= g.BcpIndex;
			}
			// video->p->memory.io[GB_REG_BCPD] = g.Palette[g.BcpIndex >> 1] >> (8 * (g.BcpIndex & 1));
		case GB_REG_OCPD:
			if g.Mode() != 3 {
				if g.OcpIndex&1 == 1 {
					g.Palette[8*4+(g.OcpIndex>>1)] &= 0x00FF
					g.Palette[8*4+(g.OcpIndex>>1)] |= Color(value) << 8
				} else {
					g.Palette[8*4+(g.OcpIndex>>1)] &= 0xFF00
					g.Palette[8*4+(g.OcpIndex>>1)] |= Color(value)
				}
				g.Renderer.writePalette(8*4+(g.OcpIndex>>1), g.Palette[8*4+(g.OcpIndex>>1)])
			}
			if g.OcpIncrement != 0 {
				g.OcpIndex++
				g.OcpIndex &= 0x3F
				// video->p->memory.io[GB_REG_OCPS] &= 0x80;
				// video->p->memory.io[GB_REG_OCPS] |= g.OcpIndex;
			}
			// video->p->memory.io[GB_REG_OCPD] = g.Palette[8 * 4 + (g.OcpIndex >> 1)] >> (8 * (g.OcpIndex & 1));
		}
	}
}

// GBVideoSwitchBank
func (g *GPU) SwitchBank(value byte) {
	value &= 1
	g.VRAM.Bank = uint16(value)
}

// GBVideoSetPalette
func (g *GPU) SetPalette(index uint, color uint32) {
	if index >= 12 {
		return
	}
	g.dmgPalette[index] = uint16((((color) & 0xF8) << 7) | (((color) & 0xF800) >> 6) | (((color) & 0xF80000) >> 19))
}

// GBVideoProcessDots
func (g *GPU) ProcessDots(cyclesLate uint32) {
	if g.Mode() != 3 {
		return
	}

	oldX := 0
	g.X = GB_VIDEO_HORIZONTAL_PIXELS
	g.Renderer.drawRange(oldX, g.X, g.Renderer.lastX)
}

// mode0 = HBlank
// 204 cycles
func (g *GPU) EndMode0() {
	if g.frameskipCounter <= 0 {
		g.Renderer.finishScanline(g.Ly)
	}

	// lyc := 0
	g.Ly++

	// oldStat := g.Stat
	if g.Ly < GB_VIDEO_VERTICAL_PIXELS {
		g.setMode(2)
	} else {
		g.setMode(1)
	}
}

// mode1 = VBlank
func (g *GPU) EndMode1() {
	if !util.Bit(g.LCDC, Enable) {
		return
	}

	g.Ly++
	switch g.Ly {
	case GB_VIDEO_VERTICAL_TOTAL_PIXELS + 1:
		g.Ly = 0
		g.setMode(2)
	case GB_VIDEO_VERTICAL_TOTAL_PIXELS:
	case GB_VIDEO_VERTICAL_TOTAL_PIXELS - 1:
	default:
	}

}

// mode2 = [mode0 -> mode2 -> mode3] -> [mode0 -> mode2 -> mode3] -> ...
// 80 cycles
func (g *GPU) EndMode2() {
	g.Renderer.cleanOAM(g.Ly)
	g.X = -(int(g.Renderer.scx) & 7)
	g.setMode(3)
}

// mode3 = [mode0 -> mode2 -> mode3] -> [mode0 -> mode2 -> mode3] -> ...
// 172 cycles
func (g *GPU) EndMode3(cyclesLate uint32) {
	g.ProcessDots(cyclesLate)
	g.setMode(0)
}

func (g *GPU) UpdateFrameCount() {
	g.frameskipCounter--
	if g.frameskipCounter < 0 {
		g.Renderer.finishFrame()
		g.frameskipCounter = g.frameskip
	}
	g.frameCounter++
}

func (g *GPU) Mode() byte {
	return g.Stat & 0x3
}

func (g *GPU) setMode(mode byte) {
	g.Stat = (g.Stat & 0xfc) | mode
}
