package cartridge

const (
	ROM = iota
	MBC1
	MBC2
	MBC3
	MBC5
)

// Cartridge - Cartridge info from ROM Header
type Cartridge struct {
	Title   string
	IsCGB   bool // gameboy color ROM is true
	Type    uint8
	ROMSize uint8
	RAMSize uint8
	MBC     int
	Debug   *Debug
}

// ParseCartridge - read cartridge info from byte slice
func (cart *Cartridge) ParseCartridge(rom []byte) {
	var titleBuf []byte
	for i := 0x0134; i < 0x0143; i++ {
		if rom[i] == 0 {
			break
		}
		titleBuf = append(titleBuf, rom[i])
	}
	cart.Title = string(titleBuf)
	cart.IsCGB = (uint8(rom[0x0143]) == 0x80 || uint8(rom[0x0143]) == 0xc0)
	cart.Type = uint8(rom[0x0147])
	cart.ROMSize = uint8(rom[0x0148])
	cart.RAMSize = uint8(rom[0x0149])

	cart.Debug = cart.newDebug()
}
