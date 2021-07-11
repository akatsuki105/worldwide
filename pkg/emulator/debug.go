package emulator

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (e *Emulator) drawDebugScreen() *ebiten.Image {
	screen := ebiten.NewImage(DEBUG_BG_X, DEBUG_BG_Y)
	screen.Fill(color.RGBA{DEBUG_BG[0], DEBUG_BG[1], DEBUG_BG[2], 0xff})

	// game screen
	e.drawDebugGameScreen(screen)

	// debug title
	ebitenutil.DebugPrintAt(screen, "GameBoy screen", 10, 5)

	// debug register
	ebitenutil.DebugPrintAt(screen, e.debugger.Register(), 340, 5)
	ebitenutil.DebugPrintAt(screen, e.debugger.IOMap(), 490, 5)

	// debug cartridge
	ebitenutil.DebugPrintAt(screen, e.debugger.Cartridge(), 680, 5)

	// tile view
	e.drawDebugTileView(screen)

	return screen
}

func (e *Emulator) drawDebugGameScreen(screen *ebiten.Image) {
	gbcScreen := ebiten.NewImage(160, 144)
	gbcScreen.ReplacePixels(e.GBC.Draw())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(float64(10), float64(25))
	screen.DrawImage(gbcScreen, op)
}

const (
	TILE_PER_ROW = 16
)

func (e *Emulator) drawDebugTileView(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Tile View(Bank0, Bank1)", 10, 320)
	banks := e.debugger.TileView()

	// each row has 16 tiles
	for row := 0; row < 384/TILE_PER_ROW; row++ {
		tileStart, tileEnd := row*TILE_PER_ROW, (row+1)*TILE_PER_ROW
		tileView := ebiten.NewImage(8*TILE_PER_ROW, 8)
		tileView.ReplacePixels(banks[0][tileStart*64*4 : tileEnd*64*4])
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(float64(10), float64(340+row*16))
		screen.DrawImage(tileView, op)
	}

	for row := 0; row < 384/TILE_PER_ROW; row++ {
		tileStart, tileEnd := row*TILE_PER_ROW, (row+1)*TILE_PER_ROW
		tileView := ebiten.NewImage(8*TILE_PER_ROW, 8)
		tileView.ReplacePixels(banks[1][tileStart*64*4 : tileEnd*64*4])
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(float64(290), float64(340+row*16))
		screen.DrawImage(tileView, op)
	}
}
