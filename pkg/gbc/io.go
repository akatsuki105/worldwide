package gbc

import (
	"github.com/pokemium/worldwide/pkg/gbc/scheduler"
	"github.com/pokemium/worldwide/pkg/gbc/video"
	"github.com/pokemium/worldwide/pkg/util"
)

const (
	JOYPIO    byte = 0x00
	SBIO      byte = 0x01
	SCIO      byte = 0x02
	DIVIO     byte = 0x04
	TIMAIO    byte = 0x05
	TMAIO     byte = 0x06
	TACIO     byte = 0x07
	IFIO      byte = 0x0f
	LCDCIO    byte = 0x40
	LCDSTATIO byte = 0x41
	SCYIO     byte = 0x42
	SCXIO     byte = 0x43
	LYIO      byte = 0x44
	LYCIO     byte = 0x45
	DMAIO     byte = 0x46
	BGPIO     byte = 0x47
	OBP0IO    byte = 0x48
	OBP1IO    byte = 0x49
	WYIO      byte = 0x4a
	WXIO      byte = 0x4b
	KEY1IO    byte = 0x4d
	VBKIO     byte = 0x4f
	HDMA1IO   byte = 0x51
	HDMA2IO   byte = 0x52
	HDMA3IO   byte = 0x53
	HDMA4IO   byte = 0x54
	HDMA5IO   byte = 0x55
	BCPSIO    byte = 0x68
	BCPDIO    byte = 0x69
	OCPSIO    byte = 0x6a
	OCPDIO    byte = 0x6b
	SVBKIO    byte = 0x70
	IEIO      byte = 0xff
)

func (g *GBC) resetIO() {
	model := g.Video.Renderer.Model

	g.storeIO(DIVIO, 0x1e)
	g.IO[TIMAIO] = 0x00
	g.storeIO(TMAIO, 0x00)
	g.storeIO(TACIO, 0xf8)
	g.storeIO(IFIO, 0xe1)
	g.IO[0x10], g.IO[0x11], g.IO[0x12], g.IO[0x14] = 0x80, 0xbf, 0xf3, 0xbf // sound1
	g.IO[0x16], g.IO[0x19] = 0x3f, 0xbf                                     // sound2
	g.IO[0x1a], g.IO[0x1b], g.IO[0x1c], g.IO[0x1e] = 0x7f, 0xff, 0x9f, 0xbf // sound3
	g.IO[0x20], g.IO[0x23] = 0xff, 0xbf                                     // sound4
	g.IO[0x24], g.IO[0x25], g.IO[0x26] = 0x77, 0xf3, 0xf1                   // sound control
	g.storeIO(LCDSTATIO, 0x85)
	g.storeIO(BGPIO, 0xfc)
	if model < util.GB_MODEL_CGB {
		g.storeIO(OBP0IO, 0xff)
		g.storeIO(OBP1IO, 0xff)
	}
	if model&util.GB_MODEL_CGB != 0 {
		g.storeIO(VBKIO, 0)
		g.storeIO(BCPSIO, 0x80)
		g.storeIO(OCPSIO, 0)
		g.storeIO(SVBKIO, 1)
	}
}

func (g *GBC) loadIO(offset byte) (value byte) {
	switch offset {
	case JOYPIO:
		value = g.joypad.Output()
	case LCDCIO:
		value = g.Video.LCDC
	case LCDSTATIO:
		value = g.Video.Stat
	default:
		if (offset >= 0x10 && offset <= 0x26) || (offset >= 0x30 && offset <= 0x3f) {
			value = g.Sound.Read(offset)
		} else {
			value = g.IO[offset]
		}
	}
	return value
}

func (g *GBC) storeIO(offset byte, value byte) {
	switch {
	case offset == JOYPIO:
		g.joypad.P1 = value

	case offset == DIVIO:
		g.timer.divReset()
		return

	case offset == TIMAIO:
		if value > 0 && g.scheduler.Until(scheduler.TimerIRQ) > (2-util.Bool2U64(g.DoubleSpeed)) {
			g.scheduler.DescheduleEvent(scheduler.TimerIRQ)
		}
		if g.scheduler.Until(scheduler.TimerIRQ) == util.Bool2U64(g.DoubleSpeed)-2 {
			return
		}

	case offset == TACIO:
		value = g.timer.updateTAC(value)

	case offset == IFIO:
		g.IO[IFIO] = value | 0xe0 // IF[4-7] always set
		g.updateIRQs()
		return

	case offset == DMAIO: // dma transfer
		base := uint16(value) << 8
		if base >= 0xe000 {
			base &= 0xdfff
		}
		g.scheduler.DescheduleEvent(scheduler.OAMDMA)
		g.scheduler.ScheduleEvent(scheduler.OAMDMA, g.dmaService, 8*(2-util.Bool2U64(g.DoubleSpeed)))
		g.dma.src = base
		g.dma.dest = 0xFE00
		g.dma.remaining = 0xa0

	case offset >= 0x10 && offset <= 0x26: // sound io
		g.Sound.Write(offset, value)
	case offset >= 0x30 && offset <= 0x3f: // sound io
		g.Sound.WriteWaveform(offset, value)

	case offset == LCDCIO:
		g.Video.ProcessDots(0)
		old := g.Video.LCDC
		g.Video.Renderer.WriteVideoRegister(offset, value)
		g.Video.WriteLCDC(old, value)

	case offset == LCDSTATIO:
		g.Video.WriteSTAT(value)

	case offset == LYCIO:
		g.Video.WriteLYC(value)

	case offset == SCYIO || offset == SCXIO || offset == WYIO || offset == WXIO:
		g.Video.ProcessDots(0)
		value = g.Video.Renderer.WriteVideoRegister(offset, value)

	case offset == BGPIO || offset == OBP0IO || offset == OBP1IO:
		g.Video.ProcessDots(0)
		g.Video.WritePalette(offset, value)

	// below case statements, gbc only
	case offset == VBKIO: // switch vram bank
		g.Video.SwitchBank(value)

	case offset == HDMA5IO:
		value = g.writeHDMA5(value)

	case offset == BCPSIO:
		g.Video.BcpIndex = int(value & 0x3f)
		g.Video.BcpIncrement = int(value & 0x80)
		g.IO[BCPDIO] = byte(g.Video.Palette[g.Video.BcpIndex>>1] >> (8 * (g.Video.BcpIndex & 1)))

	case offset == OCPSIO:
		g.Video.OcpIndex = int(value & 0x3f)
		g.Video.OcpIncrement = int(value & 0x80)
		g.IO[OCPDIO] = byte(g.Video.Palette[8*4+(g.Video.OcpIndex>>1)] >> (8 * (g.Video.OcpIndex & 1)))

	case offset == BCPDIO || offset == OCPDIO:
		if g.Video.Mode() != 3 {
			g.Video.ProcessDots(0)
		}
		g.Video.WritePalette(offset, value)

	case offset == SVBKIO: // switch wram bank
		bank := value & 0x07
		if bank == 0 {
			bank = 1
		}
		g.WRAM.bank = bank

	case offset == IEIO:
		g.IO[IEIO] = value
		g.updateIRQs()
		return
	}

	g.IO[offset] = value
}

// GBMemoryWriteHDMA5
func (g *GBC) writeHDMA5(value byte) byte {
	g.hdma.src = uint16(g.IO[HDMA1IO]) << 8
	g.hdma.src |= uint16(g.IO[HDMA2IO])
	g.hdma.dest = uint16(g.IO[HDMA3IO]) << 8
	g.hdma.dest |= uint16(g.IO[HDMA4IO])
	g.hdma.src &= 0xfff0

	g.hdma.dest &= 0x1ff0
	g.hdma.dest |= 0x8000
	wasHdma := g.hdma.enable
	g.hdma.enable = value&0x80 > 0

	if (!wasHdma && !g.hdma.enable) || (util.Bit(g.Video.LCDC, video.Enable) && g.Video.Mode() == 0) {
		if g.hdma.enable {
			g.hdma.remaining = 0x10
		} else {
			g.hdma.remaining = int(((value & 0x7F) + 1) * 0x10)
		}
		g.cpuBlocked = true
		g.scheduler.ScheduleEvent(scheduler.HDMA, g.hdmaService, 0)
	} else if g.hdma.enable && !util.Bit(g.Video.LCDC, video.Enable) {
		return 0x80 | byte((value+1)&0x7f)
	}

	return value & 0x7f
}
