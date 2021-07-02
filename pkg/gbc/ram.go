package gbc

import (
	"gbc/pkg/gbc/cart"
)

var done = make(chan int)

// Load8 fetch value from ram
func (g *GBC) Load8(addr uint16) (value byte) {
	switch {
	case addr >= 0x4000 && addr < 0x8000: // rom bank
		value = g.ROMBank.bank[g.ROMBank.ptr][addr-0x4000]
	case addr >= 0x8000 && addr < 0xa000: // vram bank
		value = g.GPU.VRAM.Buffer[addr-0x8000+0x2000*g.GPU.VRAM.Bank]
	case addr >= 0xa000 && addr < 0xc000: // rtc or ram bank
		if g.RTC.Mapped != 0 {
			value = g.RTC.Read(byte(g.RTC.Mapped))
		} else {
			value = g.RAMBank.Bank[g.RAMBank.ptr][addr-0xa000]
		}
	case g.WRAMBank.ptr > 1 && addr >= 0xd000 && addr < 0xe000: // wram bank
		value = g.WRAMBank.bank[g.WRAMBank.ptr][addr-0xd000]
	case addr >= 0xff00:
		value = g.loadIO(addr)
	default:
		value = g.RAM[addr]
	}
	return value
}

func (g *GBC) loadIO(addr uint16) (value byte) {
	switch {
	case addr == JOYPADIO:
		value = g.joypad.Output()
	case addr == SBIO:
	case addr == SCIO:
	case (addr >= 0xff10 && addr <= 0xff26) || (addr >= 0xff30 && addr <= 0xff3f): // sound IO
		value = g.Sound.Read(addr)
	case addr == LCDCIO:
		value = g.GPU.LCDC
	case addr == LCDSTATIO:
		value = g.GPU.Stat
	case addr == LYIO:
		value = byte(g.GPU.Ly)
	default:
		value = g.RAM[addr]
	}
	return value
}

// Store8 set value into RAM
func (g *GBC) Store8(addr uint16, value byte) {

	if addr <= 0x7fff { // rom
		if (addr >= 0x2000) && (addr <= 0x3fff) {
			switch g.Cartridge.MBC {
			case cart.MBC1: // lower 5bit in romptr
				if value == 0 {
					value++
				}
				upper2 := g.ROMBank.ptr >> 5
				lower5 := value
				newROMBankPtr := (upper2 << 5) | lower5
				g.switchROMBank(newROMBankPtr)
			case cart.MBC3:
				if g.GPU.HBlankDMALength == 0 {
					newROMBankPtr := value & 0x7f
					if newROMBankPtr == 0 {
						newROMBankPtr++
					}
					g.switchROMBank(newROMBankPtr)
				}
			case cart.MBC5:
				if addr < 0x3000 { // lower 8bit
					g.switchROMBank(value)
				}
			}
		} else if (addr >= 0x4000) && (addr <= 0x5fff) {
			switch g.Cartridge.MBC {
			case cart.MBC1:
				if g.bankMode == 0 { // switch upper 2bit in romptr
					upper2 := value
					lower5 := g.ROMBank.ptr & 0x1f
					newROMBankPtr := (upper2 << 5) | lower5
					g.switchROMBank(newROMBankPtr)
				} else if g.bankMode == 1 { // switch RAMptr
					newRAMBankPtr := value
					g.RAMBank.ptr = newRAMBankPtr
				}
			case cart.MBC3:
				switch {
				case value <= 0x07 && g.GPU.HBlankDMALength == 0:
					g.RTC.Mapped = 0
					g.RAMBank.ptr = value
				case value >= 0x08 && value <= 0x0c:
					g.RTC.Mapped = uint(value)
				}
			case cart.MBC5:
				// fmt.Println(value)
				g.RAMBank.ptr = value & 0x0f
			}
		} else if (addr >= 0x6000) && (addr <= 0x7fff) {
			switch g.Cartridge.MBC {
			case cart.MBC1:
				// ROM/RAM mode selection
				if value == 1 || value == 0 {
					g.bankMode = uint(value)
				}
			case cart.MBC3:
				if value == 1 {
					g.RTC.Latched = false
				} else if value == 0 {
					g.RTC.Latched = true
					g.RTC.Latch()
				}
			}
		}
	} else {

		if addr < 0xff80 || addr > 0xfffe { // only 0xff80-0xfffe can be accessed during OAMDMA
			if g.OAMDMA.ptr > 0 && g.OAMDMA.ptr <= 160 {
				return
			}
		}

		switch {
		case addr >= 0x8000 && addr < 0xa000: // vram
			g.GPU.VRAM.Buffer[addr-0x8000+(0x2000*g.GPU.VRAM.Bank)] = value
		case addr >= 0xa000 && addr < 0xc000: // rtc or ram
			if g.RTC.Mapped == 0 {
				g.RAMBank.Bank[g.RAMBank.ptr][addr-0xa000] = value
			} else {
				g.RTC.Write(byte(g.RTC.Mapped), value)
			}
		case g.WRAMBank.ptr > 1 && addr >= 0xd000 && addr < 0xe000: // wram
			g.WRAMBank.bank[g.WRAMBank.ptr][addr-0xd000] = value
		case addr >= 0xff00:
			g.storeIO(addr, value)
		default:
			g.RAM[addr] = value
		}
	}
}

func (g *GBC) storeIO(addr uint16, value byte) {
	g.RAM[addr] = value

	switch {
	case addr == JOYPADIO:
		g.joypad.P1 = value

	case addr == DIVIO:
		g.Timer.ResetAll = true

	case addr == TIMAIO:
		if g.TIMAReload.flag {
			g.TIMAReload.flag = false
			g.RAM[TIMAIO] = value
		} else if g.TIMAReload.after {
			g.RAM[TIMAIO] = g.TIMAReload.value
		} else {
			g.RAM[TIMAIO] = value
		}

	case addr == TMAIO:
		if g.TIMAReload.flag {
			g.TIMAReload.value = value
		} else if g.TIMAReload.after {
			g.RAM[TIMAIO] = value
		}
		g.RAM[TMAIO] = value

	case addr == TACIO:
		g.Timer.TAC.Change = true
		g.Timer.TAC.Old = g.RAM[TACIO]
		g.RAM[TACIO] = value

	case addr == IFIO:
		g.RAM[IFIO] = value | 0xe0 // IF[4-7] always set

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
		g.GPU.ProcessDots(0)
		g.GPU.Renderer.WriteVideoRegister(addr&0xff, value)
		g.GPU.WriteLCDC(value)

	case addr == LCDSTATIO:
		g.GPU.Stat = value

	case addr == SCYIO || addr == SCXIO || addr == WYIO || addr == WXIO:
		g.GPU.ProcessDots(0)
		value = g.GPU.Renderer.WriteVideoRegister(addr, value)
		g.RAM[addr] = value

	case addr == BGPIO || addr == OBP0IO || addr == OBP1IO:
		g.GPU.ProcessDots(0)
		g.GPU.WritePalette(addr, value)

	// below case statements, gbc only
	case addr == VBKIO: // switch vram bank
		g.GPU.SwitchBank(value)

	case addr == HDMA5IO:
		HDMA5 := value
		mode := HDMA5 >> 7 // transfer mode
		if g.GPU.HBlankDMALength > 0 && mode == 0 {
			g.GPU.HBlankDMALength = 0
			g.RAM[HDMA5IO] |= 0x80
		} else {
			length := (int(HDMA5&0x7f) + 1) * 16 // transfer size

			switch mode {
			case 0: // generic dma
				g.doVRAMDMATransfer(length)
				g.RAM[HDMA5IO] = 0xff // complete
			case 1: // hblank dma
				g.GPU.HBlankDMALength = int(HDMA5&0x7f) + 1
				g.RAM[HDMA5IO] &= 0x7f
			}
		}

	case addr == BCPSIO:
		g.GPU.BcpIndex = int(value & 0x3f)
		g.GPU.BcpIncrement = int(value & 0x80)
		g.RAM[BCPDIO] = byte(g.GPU.Palette[g.GPU.BcpIndex>>1] >> (8 * (g.GPU.BcpIndex & 1)))

	case addr == OCPSIO:
		g.GPU.OcpIndex = int(value & 0x3f)
		g.GPU.OcpIncrement = int(value & 0x80)
		g.RAM[OCPDIO] = byte(g.GPU.Palette[8*4+(g.GPU.OcpIndex>>1)] >> (8 * (g.GPU.OcpIndex & 1)))

	case addr == BCPDIO || addr == OCPDIO:
		if g.GPU.Mode() != 3 {
			g.GPU.ProcessDots(0)
		}
		g.GPU.WritePalette(addr, value)

	case addr == SVBKIO: // switch wram bank
		newWRAMBankPtr := value & 0x07
		if newWRAMBankPtr == 0 {
			newWRAMBankPtr++
		}
		g.WRAMBank.ptr = newWRAMBankPtr
	}
}

func (g *GBC) switchROMBank(newROMBankPtr uint8) {
	switchFlag := (newROMBankPtr < (2 << g.Cartridge.ROMSize))
	if switchFlag {
		g.ROMBank.ptr = newROMBankPtr
	}
}

func (g *GBC) doVRAMDMATransfer(length int) {
	from := (uint16(g.RAM[HDMA1IO])<<8 | uint16(g.RAM[HDMA2IO])) & 0xfff0
	to := ((uint16(g.RAM[HDMA3IO])<<8 | uint16(g.RAM[HDMA4IO])) & 0x1ff0) + 0x8000

	for i := 0; i < length; i++ {
		value := g.Load8(from)
		g.Store8(to, value)
		from++
		to++
	}

	g.RAM[HDMA1IO], g.RAM[HDMA2IO] = byte(from>>8), byte(from)
	g.RAM[HDMA3IO], g.RAM[HDMA4IO] = byte(to>>8), byte(to&0xf0)
}
