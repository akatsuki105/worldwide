package gbc

import (
	"gbc/pkg/gbc/scheduler"
	"gbc/pkg/util"
)

const (
	JOYPADIO  byte = 0x00
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
	model := g.video.Renderer.Model

	g.IO[DIVIO] = 0x1e
	g.IO[TIMAIO] = 0x00
	g.IO[TMAIO] = 0x00
	g.IO[TACIO] = 0xf8
	g.IO[IFIO] = 0xe1
	g.IO[0x10] = 0x80
	g.IO[0x11] = 0xbf
	g.IO[0x12] = 0xf3
	g.IO[0x14] = 0xbf
	g.IO[0x16] = 0x3f
	g.IO[0x17] = 0x00
	g.IO[0x19] = 0xbf
	g.IO[0x1a] = 0x7f
	g.IO[0x1b] = 0xff
	g.IO[0x1c] = 0x9f
	g.IO[0x1e] = 0xbf
	g.IO[0x20] = 0xff
	g.IO[0x21] = 0x00
	g.IO[0x22] = 0x00
	g.IO[0x23] = 0xbf
	g.IO[0x24] = 0x77
	g.IO[0x25] = 0xf3
	g.IO[0x26] = 0xf1
	g.storeIO(LCDCIO, 0x91)
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
	case JOYPADIO:
		value = g.joypad.Output()
	case LCDCIO:
		value = g.video.LCDC
	case LCDSTATIO:
		value = g.video.Stat
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
	case offset == JOYPADIO:
		g.joypad.P1 = value

	case offset == DIVIO:
		g.timer.divReset()
		return

	case offset == TIMAIO:
		if value > 0 && g.scheduler.Until(scheduler.TimerIRQ) > (2-util.Bool2U64(g.doubleSpeed)) {
			g.scheduler.DescheduleEvent(scheduler.TimerIRQ)
		}
		if g.scheduler.Until(scheduler.TimerIRQ) == util.Bool2U64(g.doubleSpeed)-2 {
			return
		}

	case offset == TACIO:
		value = g.timer.updateTAC(value)

	case offset == IFIO:
		g.IO[IFIO] = value | 0xe0 // IF[4-7] always set

	case offset == DMAIO: // dma transfer
		base := uint16(value) << 8
		if base >= 0xe000 {
			base &= 0xdfff
		}
		g.scheduler.DescheduleEvent(scheduler.OAMDMA)
		g.scheduler.ScheduleEvent(scheduler.OAMDMA, g.DMAService, 8*(2-util.Bool2U64(g.doubleSpeed)))
		g.dma.src = base
		g.dma.dest = 0xFE00
		g.dma.remaining = 0xa0

	case offset >= 0x10 && offset <= 0x26: // sound io
		g.Sound.Write(offset, value)
	case offset >= 0x30 && offset <= 0x3f: // sound io
		g.Sound.WriteWaveform(offset, value)

	case offset == LCDCIO:
		g.video.ProcessDots(0)
		g.video.Renderer.WriteVideoRegister(offset, value)
		g.video.WriteLCDC(value)

	case offset == LCDSTATIO:
		g.video.WriteSTAT(value)

	case offset == LYCIO:
		g.video.WriteLYC(value)

	case offset == SCYIO || offset == SCXIO || offset == WYIO || offset == WXIO:
		g.video.ProcessDots(0)
		value = g.video.Renderer.WriteVideoRegister(offset, value)
		g.IO[offset] = value

	case offset == BGPIO || offset == OBP0IO || offset == OBP1IO:
		g.video.ProcessDots(0)
		g.video.WritePalette(offset, value)

	// below case statements, gbc only
	case offset == VBKIO: // switch vram bank
		g.video.SwitchBank(value)

	case offset == HDMA5IO:
		HDMA5 := value
		mode := HDMA5 >> 7 // transfer mode
		if g.video.HBlankDMALength > 0 && mode == 0 {
			g.video.HBlankDMALength = 0
			g.IO[HDMA5IO] |= 0x80
		} else {
			length := (int(HDMA5&0x7f) + 1) * 16 // transfer size

			switch mode {
			case 0: // generic dma
				g.doVRAMDMATransfer(length)
				g.IO[HDMA5IO] = 0xff // complete
			case 1: // hblank dma
				g.video.HBlankDMALength = int(HDMA5&0x7f) + 1
				g.IO[HDMA5IO] &= 0x7f
			}
		}

	case offset == BCPSIO:
		g.video.BcpIndex = int(value & 0x3f)
		g.video.BcpIncrement = int(value & 0x80)
		g.IO[BCPDIO] = byte(g.video.Palette[g.video.BcpIndex>>1] >> (8 * (g.video.BcpIndex & 1)))

	case offset == OCPSIO:
		g.video.OcpIndex = int(value & 0x3f)
		g.video.OcpIncrement = int(value & 0x80)
		g.IO[OCPDIO] = byte(g.video.Palette[8*4+(g.video.OcpIndex>>1)] >> (8 * (g.video.OcpIndex & 1)))

	case offset == BCPDIO || offset == OCPDIO:
		if g.video.Mode() != 3 {
			g.video.ProcessDots(0)
		}
		g.video.WritePalette(offset, value)

	case offset == SVBKIO: // switch wram bank
		newWRAMBankPtr := value & 0x07
		if newWRAMBankPtr == 0 {
			newWRAMBankPtr++
		}
		g.WRAMBank.ptr = newWRAMBankPtr

	case offset == IEIO:
		g.IO[IEIO] = value
		g.updateIRQs()
	case offset == IFIO:
		g.IO[IFIO] = value | 0xE0
		g.updateIRQs()
	}

	g.IO[offset] = value
}
