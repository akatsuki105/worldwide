package gbc

import (
	"fmt"
	"gbc/pkg/debug"
	"gbc/pkg/gpu"
	"gbc/pkg/util"
	"image/jpeg"
	"os"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Debug - Info used in debug mode
type Debug struct {
	on      bool
	Break   debug.Break
	history debug.History
	pause   debug.Pause
	Window  debug.Window
	monitor debug.Monitor
}

func (cpu *CPU) SetWindowSize(x, y int) {
	cpu.debug.Window.SetSize(x, y)
}

func (cpu *CPU) debugRegister() string {
	A, F := cpu.Reg.A, cpu.Reg.F
	B, C := cpu.Reg.B, cpu.Reg.C
	D, E := cpu.Reg.D, cpu.Reg.E
	H, L := cpu.Reg.H, cpu.Reg.L

	bank := cpu.ROMBank.ptr
	PC := cpu.Reg.PC
	if PC < 0x4000 {
		bank = 0
	}

	return fmt.Sprintf(`Register
A: %02x       F: %02x
B: %02x       C: %02x
D: %02x       E: %02x
H: %02x       L: %02x
PC: %02x:%04x  SP: %04x`, A, F, B, C, D, E, H, L, bank, PC, cpu.Reg.SP)
}

func (cpu *CPU) debugIOMap() string {
	LCDC, STAT := cpu.FetchMemory8(LCDCIO), cpu.FetchMemory8(LCDSTATIO)
	DIV := cpu.FetchMemory8(DIVIO)
	LY, LYC := cpu.FetchMemory8(LYIO), cpu.FetchMemory8(LYCIO)
	IE, IF, IME := cpu.FetchMemory8(IEIO), cpu.FetchMemory8(IFIO), util.Bool2Int(cpu.Reg.IME)
	spd, rom := cpu.boost/2, cpu.ROMBank.ptr
	return fmt.Sprintf(`IO
LCDC: %02x   STAT: %02x
DIV: %02x
LY: %02x     LYC: %02x
IE: %02x     IF: %02x    IME: %02x
SPD: %02x    ROM: %02x`, LCDC, STAT, DIV, LY, LYC, IE, IF, IME, spd, rom)
}

const (
	WX, WY, scrollX, scrollY, scrollPixelX = 0, 0, 0, 0, 0
)

// DebugExec - used in test
func (cpu *CPU) DebugExec(frame int, output string) error {
	for i := 0; i < frame; i++ {
		for y := 0; y < 144; y++ {
			cpu.execScanline()
		}
		cpu.execVBlank()
	}

	// 最後の1frameは背景データを生成する
	for y := 0; y < 144; y++ {
		cpu.execScanline()

		LCDC := cpu.FetchMemory8(LCDCIO)
		for x := 0; x < 160; x += 8 {
			blockX, blockY := x/8, y/8

			var tileX, tileY uint
			var isWin bool
			var entryX int

			lineIdx := y % 8 // タイルの何行目を描画するか
			entryY := gpu.EntryY{}
			if util.Bit(LCDC, 5) && (WY <= uint(y)) && (WX <= uint(x)) {
				tileX, tileY = ((uint(x)-WX)/8)%32, ((uint(y)-WY)/8)%32
				isWin = true

				entryX = blockX * 8
				entryY.Block, entryY.Offset = blockY*8, y%8
			} else {
				tileX, tileY = (scrollX+uint(x))/8%32, (scrollY+uint(y))/8%32
				isWin = false

				entryX = blockX*8 - int(scrollPixelX)
				entryY.Block = blockY * 8
				entryY.Offset = y % 8
				lineIdx = (int(scrollY) + y) % 8
			}

			if util.Bit(LCDC, 7) {
				if !cpu.GPU.SetBGLine(entryX, entryY, tileX, tileY, isWin, cpu.Cartridge.IsCGB, lineIdx) {
					break
				}
			}
		}
	}
	cpu.execVBlank()
	screen := cpu.GPU.GetOriginal()

	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	opt := jpeg.Options{
		Quality: 100,
	}
	if err = jpeg.Encode(file, screen, &opt); err != nil {
		return err
	}
	return nil
}

func (cpu *CPU) checkBreakCond(breakpoint *debug.BreakPoint) bool {
	if !breakpoint.Cond.On {
		return true
	}

	lhs := uint16(0)
	switch breakpoint.Cond.LHS {
	case "A", "F", "B", "C", "D", "E", "H", "L", "AF", "BC", "DE", "HL", "SP":
		lhs = cpu.getRegister(breakpoint.Cond.LHS)
	default:
		return false
	}

	rhs := breakpoint.Cond.RHS
	switch breakpoint.Cond.Operand {
	case debug.Equal:
		return lhs == rhs
	case debug.NEqual:
		return lhs != rhs
	case debug.Gte:
		return lhs >= rhs
	case debug.Lte:
		return lhs <= rhs
	case debug.Gt:
		return lhs > rhs
	case debug.Lt:
		return lhs < rhs
	default:
		return false
	}
}

func (cpu *CPU) debugPrintOAM(screen *ebiten.Image) {
	// debug OAM
	ebitenutil.DebugPrintAt(screen, "OAM (Y, X, tile, attr)", 750, 320)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(float64(750), float64(340))
	OAMScreen := ebiten.NewImageFromImage(cpu.GPU.OAM)
	screen.DrawImage(OAMScreen, op)

	properties := [8]string{}
	for col := 0; col < 8; col++ {
		for row := 0; row < 5; row++ {
			i := row*8 + col
			Y, X, index, attr := cpu.GPU.OAMProperty(i)
			properties[col] += fmt.Sprintf("%02x\n%02x\n%02x\n%02x\n\n", Y, X, index, attr)
		}
	}

	for col, property := range properties {
		ebitenutil.DebugPrintAt(screen, property, 750+(col*64)+42, 340)
	}
}
