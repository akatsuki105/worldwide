package gbc

import (
	"fmt"
	"gbc/pkg/emulator/debug"
	"gbc/pkg/util"
	"image/jpeg"
	"os"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Debug - Info used in debug mode
type Debug struct {
	Enable  bool
	Break   debug.Break
	history debug.History
	pause   debug.Pause
	Window  debug.Window
	monitor debug.Monitor
}

func (g *GBC) SetWindowSize(x, y int) {
	g.Debug.Window.SetSize(x, y)
}

func (g *GBC) debugRegister() string {
	A, F := g.Reg.R[A], g.Reg.R[F]
	B, C := g.Reg.R[B], g.Reg.R[C]
	D, E := g.Reg.R[D], g.Reg.R[E]
	H, L := g.Reg.R[H], g.Reg.R[L]

	bank := g.ROMBank.ptr
	PC := g.Reg.PC
	if PC < 0x4000 {
		bank = 0
	}

	return fmt.Sprintf(`Register
A: %02x       F: %02x
B: %02x       C: %02x
D: %02x       E: %02x
H: %02x       L: %02x
PC: %02x:%04x  SP: %04x`, A, F, B, C, D, E, H, L, bank, PC, g.Reg.SP)
}

func (g *GBC) debugIOMap() string {
	LCDC, STAT := g.Load8(LCDCIO), g.Load8(LCDSTATIO)
	DIV := g.Load8(DIVIO)
	LY, LYC := g.Load8(LYIO), g.Load8(LYCIO)
	IE, IF, IME := g.Load8(IEIO), g.Load8(IFIO), util.Bool2Int(g.Reg.IME)
	spd, rom := util.Bool2Int(g.doubleSpeed), g.ROMBank.ptr
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
func (g *GBC) DebugExec(frame int, output string) error {
	for i := 0; i < frame; i++ {
		for y := 0; y < 144; y++ {
			g.execScanline()
		}
		g.execVBlank()
	}

	// 最後の1frameは背景データを生成する
	for y := 0; y < 144; y++ {
		g.execScanline()
	}
	g.execVBlank()
	screen := g.video.Display()

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

func (g *GBC) checkBreakCond(breakpoint *debug.BreakPoint) bool {
	if !breakpoint.Cond.On {
		return true
	}

	lhs := uint16(0)
	switch breakpoint.Cond.LHS {
	case "A", "F", "B", "C", "D", "E", "H", "L", "AF", "BC", "DE", "HL", "SP":
		lhs = g.getRegister(breakpoint.Cond.LHS)
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

func (g *GBC) debugPrintOAM(screen *ebiten.Image) {
	// debug OAM
	ebitenutil.DebugPrintAt(screen, "OAM (Y, X, tile, attr)", 750, 320)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(float64(750), float64(340))
	OAMScreen := ebiten.NewImageFromImage(g.video.OAM)
	screen.DrawImage(OAMScreen, op)

	properties := [8]string{}
	for col := 0; col < 8; col++ {
		for row := 0; row < 5; row++ {
			i := row*8 + col
			Y, X, index, attr := g.video.OAMProperty(i)
			properties[col] += fmt.Sprintf("%02x\n%02x\n%02x\n%02x\n\n", Y, X, index, attr)
		}
	}

	for col, property := range properties {
		ebitenutil.DebugPrintAt(screen, property, 750+(col*64)+42, 340)
	}
}
