package debug

import (
	"fmt"
	"gbc/pkg/gbc"
)

var rom = map[byte]string{0x00: "32KB", 0x01: "64KB", 0x02: "128KB", 0x03: "256KB", 0x04: "512KB", 0x05: "1MB", 0x06: "2MB", 0x07: "4MB", 0x08: "8MB", 0x52: "1.1MB", 0x53: "1.2MB", 0x54: "1.5MB"}
var ram = [6]string{"None", "2KB", "8KB", "32KB", "128KB", "64KB"}
var cartType = map[byte]string{0x00: "ROM ONLY", 0x01: "MBC1", 0x02: "MBC1+RAM", 0x03: "MBC1+RAM+BATTERY", 0x05: "MBC2", 0x06: "MBC2+BATTERY", 0x08: "ROM+RAM", 0x09: "ROM+RAM+BATTERY", 0x0b: "MBC1", 0x0c: "MBC1+RAM", 0x0d: "MBC1+RAM+BATTERY", 0x0f: "MBC3+TIMER+BATTERY", 0x10: "MBC3+TIMER+RAM+BATTERY", 0x11: "MBC3", 0x12: "MBC3+RAM", 0x13: "MBC3+RAM+BATTERY", 0x19: "MBC5", 0x1a: "MBC5+RAM", 0x1b: "MBC5+RAM+BATTERY", 0x1c: "MBC5+RUMBLE", 0x1d: "MBC5+RUMBLE+RAM", 0x1e: "MBC5+RUMBLE+RAM+BATTERY"}

type Cartridge struct {
	title, cartType, rom, ram string
}

func newCart(g *gbc.GBC) *Cartridge {
	cart := g.Cartridge
	return &Cartridge{cart.Title, cartType[cart.Type], rom[cart.ROMSize], ram[cart.RAMSize]}
}

func (c *Cartridge) String() string {
	result := fmt.Sprintf(`Cartridge
Title: %s
Cartridge Type: %s
ROM Size: %s
RAM Size: %s`, c.title, c.cartType, c.rom, c.ram)
	return result
}
