package gbc

import (
	"gbc/pkg/cartridge"
)

var done = make(chan int)

// FetchMemory8 fetch value from ram
func (cpu *CPU) FetchMemory8(addr uint16) (value byte) {
	switch {
	case addr >= 0x4000 && addr < 0x8000: // rom bank
		value = cpu.ROMBank.bank[cpu.ROMBank.ptr][addr-0x4000]
	case addr >= 0x8000 && addr < 0xa000: // vram bank
		value = cpu.GPU.VRAM.Bank[cpu.GPU.VRAM.Ptr][addr-0x8000]
	case addr >= 0xa000 && addr < 0xc000: // rtc or ram bank
		if cpu.RTC.Mapped != 0 {
			value = cpu.RTC.Read(byte(cpu.RTC.Mapped))
		} else {
			value = cpu.RAMBank.bank[cpu.RAMBank.ptr][addr-0xa000]
		}
	case cpu.WRAMBank.ptr > 1 && addr >= 0xd000 && addr < 0xe000: // wram bank
		value = cpu.WRAMBank.bank[cpu.WRAMBank.ptr][addr-0xd000]
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
	case (addr >= 0xff10 && addr <= 0xff26) || (addr >= 0xff30 && addr <= 0xff3f): // sound IO
		value = cpu.Sound.Read(addr)
	case addr == LCDCIO:
		value = cpu.GPU.LCDC
	case addr == LCDSTATIO:
		value = cpu.GPU.LCDSTAT
	case addr == BCPDIO: // BG Palette
		value = cpu.GPU.Palette.BGPalette[cpu.GPU.BgPalIdx()]
	case addr == OCPDIO: // OAM Palette
		value = cpu.GPU.Palette.SPRPalette[cpu.GPU.SprPalIdx()]
	default:
		value = cpu.RAM[addr]
	}
	return value
}

// SetMemory8 set value into RAM
func (cpu *CPU) SetMemory8(addr uint16, value byte) {

	if addr <= 0x7fff { // rom
		if (addr >= 0x2000) && (addr <= 0x3fff) {
			switch cpu.Cartridge.MBC {
			case cartridge.MBC1: // lower 5bit in romptr
				if value == 0 {
					value++
				}
				upper2 := cpu.ROMBank.ptr >> 5
				lower5 := value
				newROMBankPtr := (upper2 << 5) | lower5
				cpu.switchROMBank(newROMBankPtr)
			case cartridge.MBC3:
				if cpu.GPU.HBlankDMALength == 0 {
					newROMBankPtr := value & 0x7f
					if newROMBankPtr == 0 {
						newROMBankPtr++
					}
					cpu.switchROMBank(newROMBankPtr)
				}
			case cartridge.MBC5:
				if addr < 0x3000 { // lower 8bit
					cpu.switchROMBank(value)
				}
			}
		} else if (addr >= 0x4000) && (addr <= 0x5fff) {
			switch cpu.Cartridge.MBC {
			case cartridge.MBC1:
				if cpu.bankMode == 0 { // switch upper 2bit in romptr
					upper2 := value
					lower5 := cpu.ROMBank.ptr & 0x1f
					newROMBankPtr := (upper2 << 5) | lower5
					cpu.switchROMBank(newROMBankPtr)
				} else if cpu.bankMode == 1 { // switch RAMptr
					newRAMBankPtr := value
					cpu.RAMBank.ptr = newRAMBankPtr
				}
			case cartridge.MBC3:
				switch {
				case value <= 0x07 && cpu.GPU.HBlankDMALength == 0:
					cpu.RTC.Mapped = 0
					cpu.RAMBank.ptr = value
				case value >= 0x08 && value <= 0x0c:
					cpu.RTC.Mapped = uint(value)
				}
			case cartridge.MBC5:
				// fmt.Println(value)
				cpu.RAMBank.ptr = value & 0x0f
			}
		} else if (addr >= 0x6000) && (addr <= 0x7fff) {
			switch cpu.Cartridge.MBC {
			case cartridge.MBC1:
				// ROM/RAM mode selection
				if value == 1 || value == 0 {
					cpu.bankMode = uint(value)
				}
			case cartridge.MBC3:
				if value == 1 {
					cpu.RTC.Latched = false
				} else if value == 0 {
					cpu.RTC.Latched = true
					cpu.RTC.Latch()
				}
			}
		}
	} else {

		if addr < 0xff80 || addr > 0xfffe { // only 0xff80-0xfffe can be accessed during OAMDMA
			if cpu.OAMDMA.ptr > 0 && cpu.OAMDMA.ptr <= 160 {
				return
			}
		}

		switch {
		case addr >= 0x8000 && addr < 0xa000: // vram
			cpu.GPU.VRAM.Bank[cpu.GPU.VRAM.Ptr][addr-0x8000] = value
		case addr >= 0xa000 && addr < 0xc000: // rtc or ram
			if cpu.RTC.Mapped == 0 {
				cpu.RAMBank.bank[cpu.RAMBank.ptr][addr-0xa000] = value
			} else {
				cpu.RTC.Write(byte(cpu.RTC.Mapped), value)
			}
		case cpu.WRAMBank.ptr > 1 && addr >= 0xd000 && addr < 0xe000: // wram
			cpu.WRAMBank.bank[cpu.WRAMBank.ptr][addr-0xd000] = value
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

		if cpu.Serial.TransferFlag == 0 {
			cpu.Serial.WriteSC(value)

			switch value {
			case 0x80:
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
							// if next transfer is requested, stop forcibly
						default:
						}
					}
					if success {
						cpu.Serial.TransferFlag = 1
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
						default:
						}
					}
					if success {
						cpu.Serial.TransferFlag = 1
						<-cpu.serialTick
						cpu.Serial.Receive()
						cpu.Serial.ClearSC()
						cpu.setSerialFlag()
					}
				}()
			}
		}

	case addr == DIVIO:
		cpu.Timer.ResetAll = true

	case addr == TIMAIO:
		if cpu.TIMAReload.flag {
			cpu.TIMAReload.flag = false
			cpu.RAM[TIMAIO] = value
		} else if cpu.TIMAReload.after {
			cpu.RAM[TIMAIO] = cpu.TIMAReload.value
		} else {
			cpu.RAM[TIMAIO] = value
		}

	case addr == TMAIO:
		if cpu.TIMAReload.flag {
			cpu.TIMAReload.value = value
		} else if cpu.TIMAReload.after {
			cpu.RAM[TIMAIO] = value
		}
		cpu.RAM[TMAIO] = value

	case addr == TACIO:
		cpu.Timer.TAC.Change = true
		cpu.Timer.TAC.Old = cpu.RAM[TACIO]
		cpu.RAM[TACIO] = value

	case addr == IFIO:
		cpu.RAM[IFIO] = value | 0xe0 // IF[4-7] always set

	case addr == DMAIO: // dma transfer
		start := uint16(cpu.Reg.A) << 8
		if cpu.OAMDMA.ptr > 0 {
			cpu.OAMDMA.restart = start
			cpu.OAMDMA.reptr = 160 + 2 // lag
		} else {
			cpu.OAMDMA.start = start
			cpu.OAMDMA.ptr = 160 + 2 // lag
		}

	case addr >= 0xff10 && addr <= 0xff26: // sound io
		cpu.Sound.Write(addr, value)
	case addr >= 0xff30 && addr <= 0xff3f: // sound io
		cpu.Sound.WriteWaveform(addr, value)

	case addr == LCDCIO:
		cpu.GPU.LCDC = value

	case addr == LCDSTATIO:
		cpu.GPU.LCDSTAT = value

	case addr == 0xff42:
		cpu.GPU.Scroll[1] = value
	case addr == 0xff43:
		cpu.GPU.Scroll[0] = value

	case addr == BGPIO:
		cpu.GPU.Palette.DMGPalette[0] = value
	case addr == OBP0IO:
		cpu.GPU.Palette.DMGPalette[1] = value
	case addr == OBP1IO:
		cpu.GPU.Palette.DMGPalette[2] = value

	// below case statements, gbc only
	case addr == VBKIO && cpu.GPU.HBlankDMALength == 0: // switch vram bank
		cpu.GPU.VRAM.Ptr = value & 0x01

	case addr == HDMA5IO:
		HDMA5 := value
		mode := HDMA5 >> 7 // transfer mode
		if cpu.GPU.HBlankDMALength > 0 && mode == 0 {
			cpu.GPU.HBlankDMALength = 0
			cpu.RAM[HDMA5IO] |= 0x80
		} else {
			length := (int(HDMA5&0x7f) + 1) * 16 // transfer size

			switch mode {
			case 0: // generic dma
				cpu.doVRAMDMATransfer(length)
				cpu.RAM[HDMA5IO] = 0xff // complete
			case 1: // hblank dma
				cpu.GPU.HBlankDMALength = int(HDMA5&0x7f) + 1
				cpu.RAM[HDMA5IO] &= 0x7f
			}
		}

	case addr == BCPSIO:
		cpu.GPU.Palette.CGBPalette[0] = value
	case addr == OCPSIO:
		cpu.GPU.Palette.CGBPalette[1] = value
	case addr == BCPDIO: // write bg palette
		cpu.GPU.Palette.BGPalette[cpu.GPU.BgPalIdx()] = value
		if cpu.GPU.IsBgPalIncrement() {
			cpu.GPU.Palette.CGBPalette[0]++
		}
	case addr == OCPDIO: // write oam palette
		cpu.GPU.Palette.SPRPalette[cpu.GPU.SprPalIdx()] = value
		if cpu.GPU.IsSprPalIncrement() {
			cpu.GPU.Palette.CGBPalette[1]++
		}

	case addr == SVBKIO: // switch wram bank
		newWRAMBankPtr := value & 0x07
		if newWRAMBankPtr == 0 {
			newWRAMBankPtr++
		}
		cpu.WRAMBank.ptr = newWRAMBankPtr
	}
}

func (cpu *CPU) switchROMBank(newROMBankPtr uint8) {
	switchFlag := (newROMBankPtr < (2 << cpu.Cartridge.ROMSize))
	if switchFlag {
		cpu.ROMBank.ptr = newROMBankPtr
	}
}

func (cpu *CPU) doVRAMDMATransfer(length int) {
	from := (uint16(cpu.RAM[HDMA1IO])<<8 | uint16(cpu.RAM[HDMA2IO])) & 0xfff0
	to := ((uint16(cpu.RAM[HDMA3IO])<<8 | uint16(cpu.RAM[HDMA4IO])) & 0x1ff0) + 0x8000

	for i := 0; i < length; i++ {
		value := cpu.FetchMemory8(from)
		cpu.SetMemory8(to, value)
		from++
		to++
	}

	cpu.RAM[HDMA1IO], cpu.RAM[HDMA2IO] = byte(from>>8), byte(from)
	cpu.RAM[HDMA3IO], cpu.RAM[HDMA4IO] = byte(to>>8), byte(to&0xf0)
}
