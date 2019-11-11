package cartridge

// Cartridge ROMヘッダから得られたカードリッジ情報
type Cartridge struct {
	Title   string
	IsCGB   bool // 0x80 or 0xc0 => ゲームボーイカラーで true
	Type    uint8
	ROMSize uint8
	RAMSize uint8
	MBC     string
}

// ParseCartridge カートリッジ情報を読み取る
func (cart *Cartridge) ParseCartridge(rom *[]byte) {
	var titleBuf []byte
	for i := 0x0134; i < 0x0143; i++ {
		if (*rom)[i] == 0 {
			break
		}
		titleBuf = append(titleBuf, (*rom)[i])
	}
	cart.Title = string(titleBuf)
	cart.IsCGB = (uint8((*rom)[0x0143]) == 0x80 || uint8((*rom)[0x0143]) == 0xc0)
	cart.Type = uint8((*rom)[0x0147])
	cart.ROMSize = uint8((*rom)[0x0148])
	cart.RAMSize = uint8((*rom)[0x0149])
}
