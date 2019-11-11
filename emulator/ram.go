package emulator

// FetchMemory8 引数で指定したアドレスから値を取得する
func (cpu *CPU) FetchMemory8(addr uint16) byte {
	var value byte
	if addr >= 0x4000 && addr < 0x8000 {
		// ROMバンク
		value = cpu.ROMBank[cpu.ROMBankPtr][addr-0x4000]
	} else if addr >= 0xa000 && addr < 0xc000 {
		// RAMバンク
		value = cpu.RAMBank[cpu.RAMBankPtr][addr-0xa000]
	} else if (addr >= 0xff10 && addr <= 0xff26) || (addr >= 0xff30 && addr <= 0xff3f) {
		value = cpu.Sound.Read(addr)
	} else {
		value = cpu.RAM[addr]
	}
	return value
}

// SetMemory8 引数で指定したアドレスにvalueを書き込む
func (cpu *CPU) SetMemory8(addr uint16, value byte) {
	if addr <= 0x7fff {
		// ROM領域
		if (addr >= 0x2000) && (addr <= 0x3fff) {
			// ROMバンク下位5bit
			if value == 0 {
				value = 1
			}
			upper2 := cpu.ROMBankPtr >> 5
			lower5 := value
			newROMBankPtr := (upper2 << 5) | lower5
			cpu.switchROMBank(newROMBankPtr)
		} else if (addr >= 0x4000) && (addr <= 0x5fff) {
			// RAM バンク番号または、 ROM バンク番号の上位ビット
			if cpu.bankMode == 0 {
				// ROMptrの上位2bitの切り替え
				upper2 := value
				lower5 := cpu.ROMBankPtr & 0x1f
				newROMBankPtr := (upper2 << 5) | lower5
				cpu.switchROMBank(newROMBankPtr)
			} else if cpu.bankMode == 1 {
				// RAMptrの切り替え
				newRAMBankPtr := value
				cpu.switchRAMBank(newRAMBankPtr)
			}
		} else if (addr >= 0x6000) && (addr <= 0x7fff) {
			// ROM/RAM モード選択
			if value == 1 || value == 0 {
				cpu.bankMode = uint(value)
			}
		}
	} else {
		// RAM領域
		if addr >= 0xa000 && addr < 0xc000 {
			cpu.RAMBank[cpu.RAMBankPtr][addr-0xa000] = value
		} else {
			cpu.RAM[addr] = value
		}

		// DMA転送
		if addr == DMAIO {
			start := uint16(cpu.getAReg()) << 8
			for i := 0; i <= 0x9f; i++ {
				cpu.SetMemory8(0xfe00+uint16(i), cpu.FetchMemory8(start+uint16(i)))
			}
		}

		// VRAMアクセス
		if (addr >= 0x8000) && (addr <= 0x97ff) {
			cpu.VRAMModified = true
		}

		// パレットアクセス
		if addr == BGPIO {
			cpu.PalleteModified.BGP = true
		}
		if addr == OBP0IO {
			cpu.PalleteModified.OBP0 = true
		}
		if addr == OBP1IO {
			cpu.PalleteModified.OBP1 = true
		}

		// サウンドアクセス
		if addr >= 0xff10 && addr <= 0xff26 {
			cpu.Sound.Write(addr, value)
		}
		if addr >= 0xff30 && addr <= 0xff3f {
			cpu.Sound.WriteWaveform(addr, value)
		}
	}
}
