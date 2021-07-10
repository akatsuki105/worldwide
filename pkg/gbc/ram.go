package gbc

import (
	"gbc/pkg/gbc/cart"
)

// Load8 fetch value from ram
func (g *GBC) Load8(addr uint16) (value byte) {
	switch {
	case addr >= 0x4000 && addr < 0x8000: // rom bank
		value = g.ROMBank.bank[g.ROMBank.ptr][addr-0x4000]
	case addr >= 0x8000 && addr < 0xa000: // vram bank
		value = g.Video.VRAM.Buffer[addr-0x8000+(0x2000*g.Video.VRAM.Bank)]
	case addr >= 0xa000 && addr < 0xc000: // rtc or ram bank
		if g.RTC.Mapped != 0 {
			value = g.RTC.Read(byte(g.RTC.Mapped))
		} else {
			value = g.RAMBank.Bank[g.RAMBank.ptr][addr-0xa000]
		}
	case g.WRAMBank.ptr > 1 && addr >= 0xd000 && addr < 0xe000: // wram bank
		value = g.WRAMBank.bank[g.WRAMBank.ptr][addr-0xd000]
	case addr >= 0xfe00 && addr <= 0xfe9f:
		value = g.Video.Oam.Get(addr - 0xfe00)
	case addr >= 0xff00:
		value = g.loadIO(byte(addr))
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
				newROMBankPtr := value & 0x7f
				if newROMBankPtr == 0 {
					newROMBankPtr++
				}
				g.switchROMBank(newROMBankPtr)
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
				case value <= 0x07:
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
		switch {
		case addr >= 0x8000 && addr < 0xa000: // vram
			g.Video.VRAM.Buffer[addr-0x8000+(0x2000*g.Video.VRAM.Bank)] = value
		case addr >= 0xa000 && addr < 0xc000: // rtc or ram
			if g.RTC.Mapped == 0 {
				g.RAMBank.Bank[g.RAMBank.ptr][addr-0xa000] = value
			} else {
				g.RTC.Write(byte(g.RTC.Mapped), value)
			}
		case g.WRAMBank.ptr > 1 && addr >= 0xd000 && addr < 0xe000: // wram
			g.WRAMBank.bank[g.WRAMBank.ptr][addr-0xd000] = value
		case addr >= 0xfe00 && addr <= 0xfe9f:
			g.Video.Oam.Set(addr-0xfe00, value)
		case addr >= 0xff00:
			g.storeIO(byte(addr), value)
		default:
			g.RAM[addr] = value
		}
	}
}

func (g *GBC) switchROMBank(newROMBankPtr uint8) {
	switchFlag := (newROMBankPtr < (2 << g.Cartridge.ROMSize))
	if switchFlag {
		g.ROMBank.ptr = newROMBankPtr
	}
}

func (g *GBC) doVRAMDMATransfer(length int) {
	from := (uint16(g.IO[HDMA1IO])<<8 | uint16(g.IO[HDMA2IO])) & 0xfff0
	to := ((uint16(g.IO[HDMA3IO])<<8 | uint16(g.IO[HDMA4IO])) & 0x1ff0) + 0x8000

	for i := 0; i < length; i++ {
		value := g.Load8(from)
		g.Store8(to, value)
		from++
		to++
	}

	g.IO[HDMA1IO], g.IO[HDMA2IO] = byte(from>>8), byte(from)
	g.IO[HDMA3IO], g.IO[HDMA4IO] = byte(to>>8), byte(to&0xf0)
}
