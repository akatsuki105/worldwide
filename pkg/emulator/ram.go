package emulator

import (
	"fmt"
)

var done = make(chan int)

// FetchMemory8 引数で指定したアドレスから値を取得する
func (cpu *CPU) FetchMemory8(addr uint16) (value byte) {
	switch {
	case addr >= 0x4000 && addr < 0x8000:
		// ROMバンク
		value = cpu.ROMBank[cpu.ROMBankPtr][addr-0x4000]
	case addr >= 0x8000 && addr < 0xa000:
		// VRAMバンク
		value = cpu.GPU.VRAMBank[cpu.GPU.VRAMBankPtr][addr-0x8000]
	case addr >= 0xa000 && addr < 0xc000:
		if cpu.RTC.Mapped != 0 {
			value = cpu.RTC.Read(byte(cpu.RTC.Mapped))
		} else {
			// RAMバンク
			value = cpu.RAMBank[cpu.RAMBankPtr][addr-0xa000]
		}
	case cpu.WRAMBankPtr > 1 && addr >= 0xd000 && addr < 0xe000:
		// WRAMバンク
		value = cpu.WRAMBank[cpu.WRAMBankPtr][addr-0xd000]
	case addr >= 0xff00:
		value = cpu.fetchIO(addr)
	default:
		value = cpu.RAM[addr]
	}
	return value
}

func (cpu *CPU) fetchIO(addr uint16) (value byte) {
	switch {
	case addr == JOYPADIO:
		value = cpu.joypad.Output()
	case addr == SBIO:
		value = cpu.Serial.ReadSB()
	case addr == SCIO:
		value = cpu.Serial.ReadSC()
	case (addr >= 0xff10 && addr <= 0xff26) || (addr >= 0xff30 && addr <= 0xff3f):
		// サウンドアクセス
		value = cpu.Sound.Read(addr)
	case addr == LCDCIO:
		value = cpu.GPU.LCDC
	case addr == LCDSTATIO:
		value = cpu.GPU.LCDSTAT
	case addr == BCPDIO:
		// 背景パレットデータ読み込み
		index := cpu.GPU.FetchBGPalleteIndex()
		value = cpu.GPU.BGPallete[index]
	case addr == OCPDIO:
		// スプライトパレットデータ読み込み
		index := cpu.GPU.FetchSPRPalleteIndex()
		value = cpu.GPU.SPRPallete[index]
	default:
		value = cpu.RAM[addr]
	}
	return value
}

// SetMemory8 引数で指定したアドレスにvalueを書き込む
func (cpu *CPU) SetMemory8(addr uint16, value byte) {
	if addr <= 0x7fff {
		// ROM領域
		if (addr >= 0x2000) && (addr <= 0x3fff) {
			switch cpu.Cartridge.MBC {
			case "MBC1":
				// ROMバンク下位5bit
				if value == 0 {
					value = 1
				}
				upper2 := cpu.ROMBankPtr >> 5
				lower5 := value
				newROMBankPtr := (upper2 << 5) | lower5
				cpu.switchROMBank(newROMBankPtr)
			case "MBC3":
				if cpu.GPU.HBlankDMALength == 0 {
					newROMBankPtr := value & 0x7f
					if newROMBankPtr == 0 {
						newROMBankPtr = 1
					}
					cpu.switchROMBank(newROMBankPtr)
				}
			case "MBC5":
				if addr < 0x3000 {
					// 下位8bit
					newROMBankPtr := value
					cpu.switchROMBank(newROMBankPtr)
				} else {
					// 上位1bit
					// fmt.Println(value)
				}
			}
		} else if (addr >= 0x4000) && (addr <= 0x5fff) {
			switch cpu.Cartridge.MBC {
			case "MBC1":
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
					cpu.RAMBankPtr = newRAMBankPtr
				}
			case "MBC3":
				switch {
				case value <= 0x07 && cpu.GPU.HBlankDMALength == 0:
					cpu.RTC.Mapped = 0
					cpu.RAMBankPtr = value
				case value >= 0x08 && value <= 0x0c:
					cpu.RTC.Mapped = uint(value)
				}
			case "MBC5":
				// fmt.Println(value)
				cpu.RAMBankPtr = value & 0x0f
			}
		} else if (addr >= 0x6000) && (addr <= 0x7fff) {
			switch cpu.Cartridge.MBC {
			case "MBC1":
				// ROM/RAM モード選択
				if value == 1 || value == 0 {
					cpu.bankMode = uint(value)
				}
			case "MBC3":
				if value == 1 {
					cpu.RTC.Latched = false
				} else if value == 0 {
					cpu.RTC.Latched = true
					cpu.RTC.Latch()
				}
			}
		}
	} else {
		switch {
		case addr >= 0x8000 && addr < 0xa000:
			// VRAM
			cpu.GPU.VRAMBank[cpu.GPU.VRAMBankPtr][addr-0x8000] = value
		case addr >= 0xa000 && addr < 0xc000:
			if cpu.RTC.Mapped == 0 {
				// RAM
				cpu.RAMBank[cpu.RAMBankPtr][addr-0xa000] = value
			} else {
				cpu.RTC.Write(byte(cpu.RTC.Mapped), value)
			}
		case cpu.WRAMBankPtr > 1 && addr >= 0xd000 && addr < 0xe000:
			// WRAM
			cpu.WRAMBank[cpu.WRAMBankPtr][addr-0xd000] = value
		case addr >= 0xff00:
			cpu.setIO(addr, value)
		default:
			cpu.RAM[addr] = value
		}
	}
}

func (cpu *CPU) setIO(addr uint16, value byte) {
	cpu.RAM[addr] = value

	switch {
	case addr == JOYPADIO:
		cpu.joypad.P1 = value

	case addr == SBIO:
		cpu.Serial.WriteSB(value)
	case addr == SCIO:

		if !cpu.Serial.InTransfer {
			cpu.Serial.WriteSC(value)

			switch value {
			case 0x80:
				// TODO: 本来は通信相手がいるいないにかかわらず8クロック後に問答無用で割り込み
				if cpu.Serial.WaitCtr > 0 {
					cpu.Serial.Wait.Done()
					cpu.Serial.WaitCtr--
				}
			case 0x81:
				close(done)
				done = make(chan int)
				go func() {
					success := false
					for !success {
						success = cpu.Serial.Transfer(0)
						select {
						case <-done:
							// 次の通信が催促されたら強制的に打ち切る
							break
						default:
						}
					}
					if success {
						cpu.Serial.InTransfer = true
						<-cpu.serialTick
						cpu.Serial.Receive()
						cpu.Serial.ClearSC()
						cpu.setSerialFlag()
					}
				}()
			case 0x83:
				close(done)
				done = make(chan int)
				go func() {
					success := false
					for !success {
						success = cpu.Serial.Transfer(0)
						select {
						case <-done:
							break
						default:
						}
					}
					if success {
						cpu.Serial.InTransfer = true
						<-cpu.serialTick
						cpu.Serial.Receive()
						cpu.Serial.ClearSC()
						cpu.setSerialFlag()
					}
				}()
			}
		}

	case addr == DMAIO:
		// DMA転送
		start := uint16(cpu.getAReg()) << 8
		for i := 0; i <= 0x9f; i++ {
			cpu.SetMemory8(0xfe00+uint16(i), cpu.FetchMemory8(start+uint16(i)))
		}
		cpu.cycleLine += 150

	case addr >= 0xff10 && addr <= 0xff26:
		// サウンドアクセス
		cpu.Sound.Write(addr, value)
	case addr >= 0xff30 && addr <= 0xff3f:
		// サウンドアクセス
		cpu.Sound.WriteWaveform(addr, value)

	case addr == LCDCIO:
		cpu.GPU.LCDC = value

	case addr == LCDSTATIO:
		cpu.GPU.LCDSTAT = value

	case addr == 0xff42:
		cpu.GPU.WriteScrollY(value)
	case addr == 0xff43:
		cpu.GPU.WriteScrollX(value)

	case addr == BGPIO:
		cpu.GPU.DMGPallte[0] = value
	case addr == OBP0IO:
		cpu.GPU.DMGPallte[1] = value
	case addr == OBP1IO:
		cpu.GPU.DMGPallte[2] = value

	// 以降はゲームボーイカラーのみ
	case addr == VBKIO && cpu.GPU.HBlankDMALength == 0:
		// VRAMバンク切り替え
		newVRAMBankPtr := value & 0x01
		cpu.GPU.VRAMBankPtr = newVRAMBankPtr

	case addr == HDMA5IO:
		HDMA5 := value
		mode := HDMA5 >> 7 // 転送モード
		if cpu.GPU.HBlankDMALength > 0 && mode == 0 {
			cpu.GPU.HBlankDMALength = 0
			cpu.RAM[HDMA5IO] |= 0x80
		} else {
			length := (int(HDMA5&0x7f) + 1) * 16 // 転送するデータ長

			switch mode {
			case 0:
				// 汎用DMA
				cpu.doVRAMDMATransfer(length)
				cpu.RAM[HDMA5IO] = 0xff // 完了
			case 1:
				// H-Blank DMA
				cpu.GPU.HBlankDMALength = int(HDMA5&0x7f) + 1
				cpu.RAM[HDMA5IO] &= 0x7f
			}
		}

	case addr == BCPSIO:
		cpu.GPU.CGBPallte[0] = value
	case addr == OCPSIO:
		cpu.GPU.CGBPallte[1] = value
	case addr == BCPDIO:
		// 背景パレットデータ書き込み
		index := cpu.GPU.FetchBGPalleteIndex()
		cpu.GPU.BGPallete[index] = value
		if cpu.GPU.FetchBGPalleteIncrement() {
			cpu.GPU.CGBPallte[0]++
		}
	case addr == OCPDIO:
		// スプライトパレットデータ書き込み
		index := cpu.GPU.FetchSPRPalleteIndex()
		cpu.GPU.SPRPallete[index] = value
		if cpu.GPU.FetchSPRPalleteIncrement() {
			cpu.GPU.CGBPallte[1]++
		}

	case addr == SVBKIO:
		// WRAMバンク切り替え
		newWRAMBankPtr := value & 0x07
		if newWRAMBankPtr == 0 {
			newWRAMBankPtr = 1
		}
		cpu.WRAMBankPtr = newWRAMBankPtr
	}
}

// ROMバンクの切り替え
func (cpu *CPU) switchROMBank(newROMBankPtr uint8) {
	switchFlag := false

	switch cpu.Cartridge.ROMSize {
	case 0:
	case 1:
		switchFlag = (newROMBankPtr < 4)
	case 2:
		switchFlag = (newROMBankPtr < 8)
	case 3:
		switchFlag = (newROMBankPtr < 16)
	case 4:
		switchFlag = (newROMBankPtr < 32)
	case 5:
		switchFlag = (newROMBankPtr < 64)
	case 6:
		switchFlag = (newROMBankPtr < 128)
	case 7:
		switchFlag = (newROMBankPtr <= 255)
	default:
		errorMsg := fmt.Sprintf("ROMSize is invalid => type:%x rom:%x ram:%x\n", cpu.Cartridge.Type, cpu.Cartridge.ROMSize, cpu.Cartridge.RAMSize)
		panic(errorMsg)
	}

	if switchFlag {
		cpu.ROMBankPtr = newROMBankPtr
	}
}

func (cpu *CPU) doVRAMDMATransfer(length int) {
	from := (uint16(cpu.FetchMemory8(HDMA1IO))<<8 | uint16(cpu.FetchMemory8(HDMA2IO))) & 0xfff0
	to := ((uint16(cpu.FetchMemory8(HDMA3IO))<<8 | uint16(cpu.FetchMemory8(HDMA4IO))) & 0x1ff0) + 0x8000

	for i := 0; i < length; i++ {
		value := cpu.FetchMemory8(from)
		cpu.SetMemory8(to, value)
		from++
		to++
	}

	cpu.RAM[HDMA1IO] = byte(from >> 8)
	cpu.RAM[HDMA2IO] = byte((from & 0xff))
	cpu.RAM[HDMA3IO] = byte(to >> 8)
	cpu.RAM[HDMA4IO] = byte(to & 0xf0)
}
