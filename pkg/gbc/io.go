package gbc

import (
	"gbc/pkg/gbc/scheduler"
	"gbc/pkg/util"
)

func (g *GBC) resetIO() {
	model := g.video.Renderer.Model

	g.IO[0x04] = 0x1e
	g.IO[0x05] = 0x00
	g.IO[0x06] = 0x00
	g.IO[0x07] = 0xf8
	g.IO[0x0f] = 0xe1
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
	g.Store8(LCDCIO, 0x91)
	g.Store8(LCDSTATIO, 0x85)
	g.Store8(BGPIO, 0xfc)
	if model < util.GB_MODEL_CGB {
		g.Store8(OBP0IO, 0xff)
		g.Store8(OBP1IO, 0xff)
	}
	if model&util.GB_MODEL_CGB != 0 {
		g.Store8(VBKIO, 0)
		g.Store8(BCPSIO, 0x80)
		g.Store8(OCPSIO, 0)
		g.Store8(SVBKIO, 1)
	}
}

func (g *GBC) loadIO(addr uint16) (value byte) {
	switch {
	case addr == JOYPADIO:
		value = g.joypad.Output()
	case (addr >= 0xff10 && addr <= 0xff26) || (addr >= 0xff30 && addr <= 0xff3f): // sound IO
		value = g.Sound.Read(addr)
	case addr == LCDCIO:
		value = g.video.LCDC
	case addr == LCDSTATIO:
		value = g.video.Stat
	default:
		value = g.IO[byte(addr)]
	}
	return value
}

func (g *GBC) storeIO(addr uint16, value byte) {
	switch {
	case addr == JOYPADIO:
		g.joypad.P1 = value

	case addr == DIVIO:
		g.timer.divReset()
		return

	case addr == TIMAIO:
		if value > 0 && g.scheduler.Until(scheduler.TimerIRQ) > (2-util.Bool2U64(g.doubleSpeed)) {
			g.scheduler.DescheduleEvent(scheduler.TimerIRQ)
		}
		if g.scheduler.Until(scheduler.TimerIRQ) == util.Bool2U64(g.doubleSpeed)-2 {
			return
		}

	case addr == TACIO:
		value = g.timer.updateTAC(value)

	case addr == IFIO:
		g.IO[IFIO-0xff00] = value | 0xe0 // IF[4-7] always set

	case addr == DMAIO: // dma transfer
		base := uint16(value) << 8
		if base >= 0xe000 {
			base &= 0xdfff
		}
		g.scheduler.DescheduleEvent(scheduler.OAMDMA)
		g.scheduler.ScheduleEvent(scheduler.OAMDMA, g.DMAService, 8*(2-util.Bool2U64(g.doubleSpeed)))
		g.dma.src = base
		g.dma.dest = 0xFE00
		g.dma.remaining = 0xa0

	case addr >= 0xff10 && addr <= 0xff26: // sound io
		g.Sound.Write(addr, value)
	case addr >= 0xff30 && addr <= 0xff3f: // sound io
		g.Sound.WriteWaveform(addr, value)

	case addr == LCDCIO:
		g.video.ProcessDots(0)
		g.video.Renderer.WriteVideoRegister(byte(addr), value)
		g.video.WriteLCDC(value)

	case addr == LCDSTATIO:
		g.video.WriteSTAT(value)

	case addr == LYCIO:
		g.video.WriteLYC(value)

	case addr == SCYIO || addr == SCXIO || addr == WYIO || addr == WXIO:
		g.video.ProcessDots(0)
		value = g.video.Renderer.WriteVideoRegister(byte(addr), value)
		g.IO[byte(addr)] = value

	case addr == BGPIO || addr == OBP0IO || addr == OBP1IO:
		g.video.ProcessDots(0)
		g.video.WritePalette(byte(addr), value)

	// below case statements, gbc only
	case addr == VBKIO: // switch vram bank
		g.video.SwitchBank(value)

	case addr == HDMA5IO:
		HDMA5 := value
		mode := HDMA5 >> 7 // transfer mode
		if g.video.HBlankDMALength > 0 && mode == 0 {
			g.video.HBlankDMALength = 0
			g.IO[HDMA5IO-0xff00] |= 0x80
		} else {
			length := (int(HDMA5&0x7f) + 1) * 16 // transfer size

			switch mode {
			case 0: // generic dma
				g.doVRAMDMATransfer(length)
				g.IO[HDMA5IO-0xff00] = 0xff // complete
			case 1: // hblank dma
				g.video.HBlankDMALength = int(HDMA5&0x7f) + 1
				g.IO[HDMA5IO-0xff00] &= 0x7f
			}
		}

	case addr == BCPSIO:
		g.video.BcpIndex = int(value & 0x3f)
		g.video.BcpIncrement = int(value & 0x80)
		g.IO[BCPDIO-0xff00] = byte(g.video.Palette[g.video.BcpIndex>>1] >> (8 * (g.video.BcpIndex & 1)))

	case addr == OCPSIO:
		g.video.OcpIndex = int(value & 0x3f)
		g.video.OcpIncrement = int(value & 0x80)
		g.IO[OCPDIO-0xff00] = byte(g.video.Palette[8*4+(g.video.OcpIndex>>1)] >> (8 * (g.video.OcpIndex & 1)))

	case addr == BCPDIO || addr == OCPDIO:
		if g.video.Mode() != 3 {
			g.video.ProcessDots(0)
		}
		g.video.WritePalette(byte(addr), value)

	case addr == SVBKIO: // switch wram bank
		newWRAMBankPtr := value & 0x07
		if newWRAMBankPtr == 0 {
			newWRAMBankPtr++
		}
		g.WRAMBank.ptr = newWRAMBankPtr

	case addr == IEIO:
		g.IO[IEIO-0xff00] = value
		g.updateIRQs()
	case addr == IFIO:
		g.IO[IFIO-0xff00] = value | 0xE0
		g.updateIRQs()
	}

	g.IO[byte(addr)] = value
}
