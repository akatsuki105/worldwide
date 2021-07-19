package cart

const (
	ROM = iota
	MBC1
	MBC2
	MBC3
	MBC5
)

// ram size
const (
	NO_RAM = iota
	RAM_UNUSED
	RAM_8KB
	RAM_32KB
	RAM_128KB
	RAM_64KB
)

// Cartridge - Cartridge info from ROM Header
type Cartridge struct {
	Title                  string
	IsCGB                  bool // gameboy color ROM is true
	Type, ROMSize, RAMSize byte
	MBC                    int
}

// Parse - load cartridge info from byte slice
func New(rom []byte) *Cartridge {
	var buf []byte
	for i := 0x0134; i < 0x0143; i++ {
		if rom[i] == 0 {
			break
		}
		buf = append(buf, rom[i])
	}

	return &Cartridge{
		Title:   string(buf),
		IsCGB:   rom[0x0143] == 0x80 || rom[0x0143] == 0xc0,
		Type:    rom[0x0147],
		ROMSize: rom[0x0148],
		RAMSize: rom[0x0149],
	}
}

func (c *Cartridge) HasRTC() bool {
	return c.Type == 0x0f || c.Type == 0x10
}
