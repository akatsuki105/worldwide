package cart

const (
	ROM = iota
	MBC1
	MBC2
	MBC3
	MBC5
)

// Cartridge - Cartridge info from ROM Header
type Cartridge struct {
	Title                  string
	IsCGB                  bool // gameboy color ROM is true
	Type, ROMSize, RAMSize uint8
	MBC                    int
}

// Parse - load cartridge info from byte slice
func New(rom []byte) *Cartridge {
	cart := &Cartridge{}
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
	return cart
}
