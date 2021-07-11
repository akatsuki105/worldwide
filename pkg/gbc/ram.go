package gbc

import (
	"github.com/pokemium/Worldwide/pkg/gbc/cart"
)

// Load8 fetch value from ram
func (g *GBC) Load8(addr uint16) (value byte) {
	switch {

	case addr < 0x4000:
		// ROM bank0
		value = g.ROM.buffer[0][addr]
	case addr >= 0x4000 && addr < 0x8000:
		// ROM bank1..256
		value = g.ROM.buffer[g.ROM.bank][addr-0x4000]

	case addr >= 0x8000 && addr < 0xa000:
		// VRAM
		value = g.Video.VRAM.Buffer[addr-0x8000+(0x2000*g.Video.VRAM.Bank)]

	case addr >= 0xa000 && addr < 0xc000:
		// RTC or RAM
		if g.RTC.Mapped != 0 {
			value = g.RTC.Read(byte(g.RTC.Mapped))
		} else {
			value = g.RAM.Buffer[g.RAM.bank][addr-0xa000]
		}

	case addr >= 0xc000 && addr < 0xd000:
		// WRAM bank0
		value = g.WRAM.buffer[0][addr-0xc000]
	case addr >= 0xd000 && addr < 0xe000:
		// WRAM bank1..7
		value = g.WRAM.buffer[g.WRAM.bank][addr-0xd000]

	case addr >= 0xfe00 && addr < 0xfea0:
		// OAM
		value = g.Video.Oam.Get(addr - 0xfe00)

	case addr >= 0xff00:
		// IO, HRAM, IE
		value = g.loadIO(byte(addr))
	}
	return value
}

// Store8 set value into RAM
func (g *GBC) Store8(addr uint16, value byte) {
	if addr <= 0x7fff {
		g.mbcWrite(addr, value)
		return
	}

	switch {
	case addr >= 0x8000 && addr < 0xa000: // vram
		g.Video.VRAM.Buffer[addr-0x8000+(0x2000*g.Video.VRAM.Bank)] = value

	case addr >= 0xa000 && addr < 0xc000:
		// RTC or RAM
		if g.RTC.Mapped == 0 {
			g.RAM.Buffer[g.RAM.bank][addr-0xa000] = value
			return
		}
		g.RTC.Write(byte(g.RTC.Mapped), value)

	case addr >= 0xc000 && addr < 0xd000:
		// WRAM bank0
		g.WRAM.buffer[0][addr-0xc000] = value
	case addr >= 0xd000 && addr < 0xe000:
		// WRAM bank1 or 7
		g.WRAM.buffer[g.WRAM.bank][addr-0xd000] = value

	case addr >= 0xfe00 && addr <= 0xfe9f:
		// OAM
		g.Video.Oam.Set(addr-0xfe00, value)

	case addr >= 0xff00:
		// IO, HRAM, IE
		g.storeIO(byte(addr), value)
	}
}

func (g *GBC) mbcWrite(addr uint16, value byte) {
	if (addr >= 0x2000) && (addr <= 0x3fff) {
		switch g.Cartridge.MBC {
		case cart.MBC1: // lower 5bit in romptr
			if value == 0 {
				value++
			}
			upper2 := g.ROM.bank >> 5
			lower5 := value
			bank := (upper2 << 5) | lower5
			g.switchROMBank(bank)
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
				lower5 := g.ROM.bank & 0x1f
				bank := (upper2 << 5) | lower5
				g.switchROMBank(bank)
			} else if g.bankMode == 1 { // switch RAMptr
				g.RAM.bank = value
			}
		case cart.MBC3:
			switch {
			case value <= 0x07:
				g.RTC.Mapped = 0
				g.RAM.bank = value
			case value >= 0x08 && value <= 0x0c:
				g.RTC.Mapped = uint(value)
			}
		case cart.MBC5:
			g.RAM.bank = value & 0x0f
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
}

func (g *GBC) switchROMBank(bank uint8) {
	switchFlag := (bank < (2 << g.Cartridge.ROMSize))
	if switchFlag {
		g.ROM.bank = bank
	}
}
