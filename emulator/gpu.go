package emulator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"

	"github.com/faiface/pixel"
)

var (
	// colors {R, G, B}
	colors [64][3]uint8 = [64][3]uint8{
		{0xff, 0xff, 0xff}, {0xcc, 0xcc, 0xcc}, {0x66, 0x66, 0x66}, {0x00, 0x00, 0x00},
	}
)

// PalleteModified パレットが変更されたか
type PalleteModified struct {
	BGP  bool
	OBP0 bool
	OBP1 bool
}

func (cpu *CPU) fetchTileBaseAddr() uint16 {
	LCDC := cpu.FetchMemory8(LCDCIO)
	if LCDC&0x10 != 0 {
		return 0x8000
	}
	return 0x8800
}

func (cpu *CPU) fetchSPRYSize() int {
	LCDC := cpu.FetchMemory8(LCDCIO)
	if LCDC&0x04 != 0 {
		return 16
	}
	return 8
}

func (cpu *CPU) setHBlankMode() {
	STAT := cpu.FetchMemory8(LCDCSTATIO)
	STAT &= 0xfc // bit0-1を00にする
	cpu.SetMemory8(LCDCSTATIO, STAT)
	if (STAT & 0x08) != 0 {
		cpu.setLCDSTATFlag()
	}
}

func (cpu *CPU) setVBlankMode() {
	STAT := cpu.FetchMemory8(LCDCSTATIO)
	STAT = (STAT | 0x01) & 0xfd // bit0-1を01にする
	cpu.SetMemory8(LCDCSTATIO, STAT)
}

func (cpu *CPU) setOAMRAMMode() {
	STAT := cpu.FetchMemory8(LCDCSTATIO)
	STAT = (STAT | 0x02) & 0xfe // bit0-1を10にする
	cpu.SetMemory8(LCDCSTATIO, STAT)
}

func (cpu *CPU) setLCDMode() {
	STAT := cpu.FetchMemory8(LCDCSTATIO)
	STAT |= 0x03 // bit0-1を11にする
	cpu.SetMemory8(LCDCSTATIO, STAT)
}

func (cpu *CPU) incrementLY() {
	LY := uint8(cpu.FetchMemory8(LYIO))
	LY++
	if LY == 144 {
		// VBlank期間フラグを立てる
		cpu.setVBlankMode()

		if cpu.Reg.IME && cpu.getVBlankEnable() {
			cpu.triggerVBlank()
		}
	}
	if LY > 153 {
		LY = 0
	}
	cpu.SetMemory8(LYIO, byte(LY))
	cpu.compareLYC(LY)
}

func (cpu *CPU) cacheOneLine(img *image.RGBA, i int) {
	tileNum := i / 8 // 何枚目のタイルか
	tileY := i % 8   // タイルのYビット目
	addr := uint16(0x8000 + 2*i)
	lowerByte := cpu.FetchMemory8(addr)     // 10010101
	upperByte := cpu.FetchMemory8(addr + 1) // 00110111

	lowerCache := cpu.VRAMCache[addr-0x8000]
	upperCache := cpu.VRAMCache[addr+1-0x8000]

	if lowerByte != lowerCache || upperByte != upperCache {
		cpu.VRAMCache[addr-0x8000] = lowerByte
		cpu.VRAMCache[addr+1-0x8000] = upperByte

		var lineWait sync.WaitGroup
		lineWait.Add(8)
		for j := 0; j < 8; j++ {
			go func(j int) {
				bitCtr := (7 - uint(j)) // 上位何ビット目を取り出すか
				upperColor := (upperByte >> bitCtr) % 2
				lowerColor := (lowerByte >> bitCtr) % 2
				pallete := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

				RGB, _ := cpu.parsePallete("BGP", pallete)
				R, G, B := colors[RGB][0], colors[RGB][1], colors[RGB][2]
				var x, y int
				if addr <= 0x8fff {
					x = tileNum*8 + j
					y = tileY
					img.Set(x, y, color.RGBA{R, G, B, 0xff})
					if addr >= 0x8800 {
						x = (tileNum-128)*8 + j
						y = 8 + tileY
						img.Set(x, y, color.RGBA{R, G, B, 0xff})
					}

					// SPR0
					OBP0RGB, transparent0 := cpu.parsePallete("OBP0", pallete)
					OBP0R, OBP0G, OBP0B := colors[OBP0RGB][0], colors[OBP0RGB][1], colors[OBP0RGB][2]

					// SPR0 反転なし
					x = tileNum*8 + j
					y = 16 + tileY
					img.Set(x, y, color.RGBA{OBP0R, OBP0G, OBP0B, transparent0})
					// SPR0 上下反転
					x = tileNum*8 + j
					y = (24 + 7) - tileY
					img.Set(x, y, color.RGBA{OBP0R, OBP0G, OBP0B, transparent0})
					// SPR0 左右反転
					x = (tileNum*8 + 7) - j
					y = 32 + tileY
					img.Set(x, y, color.RGBA{OBP0R, OBP0G, OBP0B, transparent0})
					// SPR0 上下左右反転
					x = (tileNum*8 + 7) - j
					y = (40 + 7) - tileY
					img.Set(x, y, color.RGBA{OBP0R, OBP0G, OBP0B, transparent0})

					// SPR1
					OBP1RGB, transparent1 := cpu.parsePallete("OBP1", pallete)
					OBP1R, OBP1G, OBP1B := colors[OBP1RGB][0], colors[OBP1RGB][1], colors[OBP1RGB][2]

					// SPR1 反転なし
					x = tileNum*8 + j
					y = 48 + tileY
					img.Set(x, y, color.RGBA{OBP1R, OBP1G, OBP1B, transparent1})
					// SPR1 上下反転
					x = tileNum*8 + j
					y = (56 + 7) - tileY
					img.Set(x, y, color.RGBA{OBP1R, OBP1G, OBP1B, transparent1})
					// SPR1 左右反転
					x = (tileNum*8 + 7) - j
					y = 64 + tileY
					img.Set(x, y, color.RGBA{OBP1R, OBP1G, OBP1B, transparent1})
					// SPR1 上下左右反転
					x = (tileNum*8 + 7) - j
					y = (72 + 7) - tileY
					img.Set(x, y, color.RGBA{OBP1R, OBP1G, OBP1B, transparent1})
				} else {
					x = (tileNum-128)*8 + j
					y = 8 + tileY
					img.Set(x, y, color.RGBA{R, G, B, 0xff})
				}
				lineWait.Done()
			}(j)
		}
		lineWait.Wait()
	} else if cpu.PalleteModified.BGP || cpu.PalleteModified.OBP0 || cpu.PalleteModified.OBP1 {
		cpu.VRAMCache[addr-0x8000] = lowerByte
		cpu.VRAMCache[addr+1-0x8000] = upperByte

		var lineWait sync.WaitGroup
		lineWait.Add(8)
		for j := 0; j < 8; j++ {
			go func(j int) {
				bitCtr := (7 - uint(j)) // 上位何ビット目を取り出すか
				upperColor := (upperByte >> bitCtr) % 2
				lowerColor := (lowerByte >> bitCtr) % 2
				pallete := (upperColor << 1) + lowerColor // 0 or 1 or 2 or 3

				RGB, _ := cpu.parsePallete("BGP", pallete)
				R, G, B := colors[RGB][0], colors[RGB][1], colors[RGB][2]
				var x, y int
				if addr <= 0x8fff {
					if cpu.PalleteModified.BGP {
						x = tileNum*8 + j
						y = tileY
						img.Set(x, y, color.RGBA{R, G, B, 0xff})
						if addr >= 0x8800 {
							x = (tileNum-128)*8 + j
							y = 8 + tileY
							img.Set(x, y, color.RGBA{R, G, B, 0xff})
						}
					}

					if cpu.PalleteModified.OBP0 {
						// SPR0
						OBP0RGB, transparent0 := cpu.parsePallete("OBP0", pallete)
						OBP0R, OBP0G, OBP0B := colors[OBP0RGB][0], colors[OBP0RGB][1], colors[OBP0RGB][2]

						// SPR0 反転なし
						x = tileNum*8 + j
						y = 16 + tileY
						img.Set(x, y, color.RGBA{OBP0R, OBP0G, OBP0B, transparent0})
						// SPR0 上下反転
						x = tileNum*8 + j
						y = (24 + 7) - tileY
						img.Set(x, y, color.RGBA{OBP0R, OBP0G, OBP0B, transparent0})
						// SPR0 左右反転
						x = (tileNum*8 + 7) - j
						y = 32 + tileY
						img.Set(x, y, color.RGBA{OBP0R, OBP0G, OBP0B, transparent0})
						// SPR0 上下左右反転
						x = (tileNum*8 + 7) - j
						y = (40 + 7) - tileY
						img.Set(x, y, color.RGBA{OBP0R, OBP0G, OBP0B, transparent0})
					}

					if cpu.PalleteModified.OBP1 {
						// SPR1
						OBP1RGB, transparent1 := cpu.parsePallete("OBP1", pallete)
						OBP1R, OBP1G, OBP1B := colors[OBP1RGB][0], colors[OBP1RGB][1], colors[OBP1RGB][2]

						// SPR1 反転なし
						x = tileNum*8 + j
						y = 48 + tileY
						img.Set(x, y, color.RGBA{OBP1R, OBP1G, OBP1B, transparent1})
						// SPR1 上下反転
						x = tileNum*8 + j
						y = (56 + 7) - tileY
						img.Set(x, y, color.RGBA{OBP1R, OBP1G, OBP1B, transparent1})
						// SPR1 左右反転
						x = (tileNum*8 + 7) - j
						y = 64 + tileY
						img.Set(x, y, color.RGBA{OBP1R, OBP1G, OBP1B, transparent1})
						// SPR1 上下左右反転
						x = (tileNum*8 + 7) - j
						y = (72 + 7) - tileY
						img.Set(x, y, color.RGBA{OBP1R, OBP1G, OBP1B, transparent1})
					}
				} else {
					if cpu.PalleteModified.BGP {
						x = (tileNum-128)*8 + j
						y = 8 + tileY
						img.Set(x, y, color.RGBA{R, G, B, 0xff})
					}
				}
				lineWait.Done()
			}(j)
		}
		lineWait.Wait()
	}
}

func (cpu *CPU) cacheTile() {
	// VRAMに変更があるならキャッシュを更新
	if cpu.VRAMModified || cpu.PalleteModified.BGP || cpu.PalleteModified.OBP0 || cpu.PalleteModified.OBP1 {
		img := cpu.mapCache
		// 2byteとって初めて1pxのデータがそろうので0x1800/2回のイテレーション
		var wait sync.WaitGroup
		wait.Add(0x1800 / 2)
		for i := 0; i < (0x1800 / 2); i++ {
			go func(i int) {
				cpu.cacheOneLine(img, i)
				wait.Done()
			}(i)
		}
		wait.Wait()

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			panic("error/png")
		}

		tmp, _, _ := image.Decode(buf)
		pic := pixel.PictureDataFromImage(tmp)

		cpu.newTileCache = pic
		cpu.VRAMModified = false
		cpu.PalleteModified = PalleteModified{
			BGP:  false,
			OBP0: false,
			OBP1: false,
		}
		cpu.tileModified = true
	}
}

func (cpu *CPU) outputTile(name string, tileNum uint8, attr byte) (rect pixel.Rect) {
	if name == "BG" {
		baseAddr := cpu.fetchTileBaseAddr()
		if baseAddr == 0x8000 {
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-0), float64((uint(tileNum)+1)*8), float64(80-8))
		} else if baseAddr == 0x8800 {
			tileNum = uint8(int(int8(tileNum)) + 128)
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-8), float64((uint(tileNum)+1)*8), float64(80-16))
		}
	} else if name == "SPR0" {
		if (attr>>6)%2 == 1 && (attr>>5)%2 == 1 {
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-40), float64((uint(tileNum)+1)*8), float64(80-48)) // 上下左右
		} else if (attr>>6)%2 == 1 {
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-24), float64((uint(tileNum)+1)*8), float64(80-32)) // 上下
		} else if (attr>>5)%2 == 1 {
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-32), float64((uint(tileNum)+1)*8), float64(80-40)) // 左右
		} else {
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-16), float64((uint(tileNum)+1)*8), float64(80-24)) // 反転なし
		}
	} else if name == "SPR1" {
		if (attr>>6)%2 == 1 && (attr>>5)%2 == 1 {
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-72), float64((uint(tileNum)+1)*8), float64(80-80)) // 上下左右
		} else if (attr>>6)%2 == 1 {
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-56), float64((uint(tileNum)+1)*8), float64(80-64)) // 上下
		} else if (attr>>5)%2 == 1 {
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-64), float64((uint(tileNum)+1)*8), float64(80-72)) // 左右
		} else {
			rect = pixel.R(float64(uint(tileNum)*8), float64(80-48), float64((uint(tileNum)+1)*8), float64(80-56)) // 反転なし
		}
	}
	return rect
}

// outputBGTile x, yはスクリーンデータ全体(32*32) not 20*18
func (cpu *CPU) outputBGTile(x, y uint, useWindow bool) (rect pixel.Rect) {
	index := x + y*32

	var addr uint16
	LCDC := cpu.FetchMemory8(LCDCIO)
	if useWindow {
		if LCDC&0x40 != 0 {
			addr = 0x9c00 + uint16(index)
		} else {
			addr = 0x9800 + uint16(index)
		}
	} else {
		if LCDC&0x08 != 0 {
			addr = 0x9c00 + uint16(index)
		} else {
			addr = 0x9800 + uint16(index)
		}
	}
	tileNum := uint8(cpu.FetchMemory8(addr))
	rect = cpu.outputTile("BG", tileNum, 0)
	return rect
}

// outputSPRTile スプライトを出力する
func (cpu *CPU) outputSPRTile(tileNum uint8, attr byte) (rect pixel.Rect) {
	if (attr>>4)%2 == 1 {
		rect = cpu.outputTile("SPR1", tileNum, attr)
	} else {
		rect = cpu.outputTile("SPR0", tileNum, attr)
	}
	return rect
}

func (cpu *CPU) dumpVRAM(filename string) {
	dumpname := fmt.Sprintf("./dump/%s_VRAM.dmp", filename)
	dumpfile, err := os.Create(dumpname)
	if err != nil {
		fmt.Println("dump failed.")
	}
	defer dumpfile.Close()

	data := cpu.RAM[0x8000:0xa000]
	_, err = dumpfile.Write(data)
	if err != nil {
		fmt.Println("dump failed.")
	}
}

func (cpu *CPU) compareLYC(LY uint8) {
	LYC := cpu.FetchMemory8(0xff45)
	if LYC == LY {
		// LCDC STAT IOポートの一致フラグをセットする
		STAT := cpu.FetchMemory8(LCDCSTATIO) | 0x04
		cpu.SetMemory8(LCDCSTATIO, STAT)

		enable := cpu.getLCDSTATEnable()
		if enable {
			cpu.triggerLCDC()
		}
	}
}

func (cpu *CPU) parsePallete(name string, colorNumber byte) (RGB, transparent byte) {
	var pallete byte
	switch name {
	case "BGP":
		pallete = cpu.FetchMemory8(BGPIO)
	case "OBP0":
		pallete = cpu.FetchMemory8(OBP0IO)
	case "OBP1":
		pallete = cpu.FetchMemory8(OBP1IO)
	default:
		errMsg := fmt.Sprintf("Error: BG Pallete name is invalid. %s", name)
		panic(errMsg)
	}

	transparent = 0xff // 非透明

	switch colorNumber {
	case 0:
		RGB = (pallete >> 0) % 4
		if name == "OBP0" || name == "OBP1" {
			transparent = 0x00
		}
	case 1:
		RGB = (pallete >> 2) % 4
	case 2:
		RGB = (pallete >> 4) % 4
	case 3:
		RGB = (pallete >> 6) % 4
	default:
		errMsg := fmt.Sprintf("Error: BG Pallete number is invalid. %d", colorNumber)
		panic(errMsg)
	}
	return RGB, transparent
}
