package cartridge

import "fmt"

type Debug struct {
	title         string
	cartridgeType string
	rom           string
	ram           string
}

func (cart *Cartridge) newDebug() *Debug {
	title := cart.Title
	cartridgeType := cart.cartridgeType()
	rom := cart.rom()
	ram := cart.ram()
	return &Debug{
		title,
		cartridgeType,
		rom,
		ram,
	}
}

func (debug *Debug) String() string {
	result := fmt.Sprintf(`Cartridge
Title: %s
Cartridge Type: %s
ROM Size: %s
RAM Size: %s`, debug.title, debug.cartridgeType, debug.rom, debug.ram)
	return result
}

func (cart *Cartridge) cartridgeType() string {
	switch cart.Type {
	case 0x00:
		return "ROM ONLY"
	case 0x01:
		return "MBC1"
	case 0x02:
		return "MBC1+RAM"
	case 0x03:
		return "MBC1+RAM+BATTERY"
	case 0x05:
		return "MBC2"
	case 0x06:
		return "MBC2+BATTERY"
	case 0x08:
		return "ROM+RAM"
	case 0x09:
		return "ROM+RAM+BATTERY"
	case 0x0b:
		return "MMM01"
	case 0x0c:
		return "MMM01+RAM"
	case 0x0d:
		return "MMM01+RAM+BATTERY"
	case 0x0f:
		return "MBC3+TIMER+BATTERY"
	case 0x10:
		return "MBC3+TIMER+RAM+BATTERY"
	case 0x11:
		return "MBC3"
	case 0x12:
		return "MBC3+RAM"
	case 0x13:
		return "MBC3+RAM+BATTERY"
	case 0x19:
		return "MBC5"
	case 0x1a:
		return "MBC5+RAM"
	case 0x1b:
		return "MBC5+RAM+BATTERY"
	case 0x1c:
		return "MBC5+RUMBLE"
	case 0x1d:
		return "MBC5+RUMBLE+RAM"
	case 0x1e:
		return "MBC5+RUMBLE+RAM+BATTERY"
	default:
		return "UNKNOWN"
	}
}

func (cart *Cartridge) rom() string {
	switch cart.ROMSize {
	case 0x00:
		return "32KB"
	case 0x01:
		return "64KB"
	case 0x02:
		return "128KB"
	case 0x03:
		return "256KB"
	case 0x04:
		return "512KB"
	case 0x05:
		return "1MB"
	case 0x06:
		return "2MB"
	case 0x07:
		return "4MB"
	case 0x08:
		return "8MB"
	case 0x52:
		return "1.1MB"
	case 0x53:
		return "1.2MB"
	case 0x54:
		return "1.5MB"
	default:
		return "UNKNOWN"
	}
}

func (cart *Cartridge) ram() string {
	switch cart.RAMSize {
	case 0x00:
		return "None"
	case 0x01:
		return "2KB"
	case 0x02:
		return "8KB"
	case 0x03:
		return "32KB"
	case 0x04:
		return "128KB"
	case 0x05:
		return "64KB"
	default:
		return "UNKNOWN"
	}
}
