package gbc

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
	g.IO[byte(addr)] = value

	switch {
	case addr == JOYPADIO:
		g.joypad.P1 = value

	case addr == DIVIO:
		g.Timer.ResetAll = true

	case addr == TIMAIO:
		if g.TIMAReload.flag {
			g.TIMAReload.flag = false
			g.IO[TIMAIO-0xff00] = value
		} else if g.TIMAReload.after {
			g.IO[TIMAIO-0xff00] = g.TIMAReload.value
		} else {
			g.IO[TIMAIO-0xff00] = value
		}

	case addr == TMAIO:
		if g.TIMAReload.flag {
			g.TIMAReload.value = value
		} else if g.TIMAReload.after {
			g.IO[TIMAIO-0xff00] = value
		}
		g.IO[TMAIO-0xff00] = value

	case addr == TACIO:
		g.Timer.TAC.Change = true
		g.Timer.TAC.Old = g.IO[TACIO-0xff00]
		g.IO[TACIO-0xff00] = value

	case addr == IFIO:
		g.IO[IFIO-0xff00] = value | 0xe0 // IF[4-7] always set

	case addr == DMAIO: // dma transfer
		start := uint16(g.Reg.R[A]) << 8
		if g.OAMDMA.ptr > 0 {
			g.OAMDMA.restart = start
			g.OAMDMA.reptr = 160 + 2 // lag
		} else {
			g.OAMDMA.start = start
			g.OAMDMA.ptr = 160 + 2 // lag
		}

	case addr >= 0xff10 && addr <= 0xff26: // sound io
		g.Sound.Write(addr, value)
	case addr >= 0xff30 && addr <= 0xff3f: // sound io
		g.Sound.WriteWaveform(addr, value)

	case addr == LCDCIO:
		g.video.ProcessDots(0)
		g.video.Renderer.WriteVideoRegister(byte(addr), value)
		g.video.WriteLCDC(value)

	case addr == LCDSTATIO:
		g.video.Stat = value

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
	}
}
