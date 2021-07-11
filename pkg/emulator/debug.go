package emulator

import (
	"image/color"
	"sync"

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
	for b := 0; b < 2; b++ {
		var wg sync.WaitGroup
		wg.Add(384 / TILE_PER_ROW)
		for row := 0; row < 384/TILE_PER_ROW; row++ {
			go func(row int) {
				rowStart, rowEnd := row*TILE_PER_ROW, (row+1)*TILE_PER_ROW

				rowView := ebiten.NewImage(8*TILE_PER_ROW, 8)      // [8x16, 8]
				rowBuffer := banks[b][rowStart*64*4 : rowEnd*64*4] // 0..63, 0..63, 0..63, ..
				newRowBuffer := make([]byte, TILE_PER_ROW*64*4)    // 0..7, 0..7, ... 8..15, 8..15,

				for t := 0; t < TILE_PER_ROW; t++ {
					rowBufferBase := t * 64 * 4
					for y := 0; y < 8; y++ {
						tileRowBuffer := rowBuffer[rowBufferBase+y*8*4 : rowBufferBase+(y+1)*8*4] // (y*8)..((y*8)+7)
						for x := 0; x < 8; x++ {
							idx := y*8*TILE_PER_ROW*4 + t*8*4 + x*4
							newRowBuffer[idx] = tileRowBuffer[x*4]
							newRowBuffer[idx+1] = tileRowBuffer[x*4+1]
							newRowBuffer[idx+2] = tileRowBuffer[x*4+2]
							newRowBuffer[idx+3] = tileRowBuffer[x*4+3]
						}
					}
				}
				rowView.ReplacePixels(newRowBuffer)

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(2, 2)
				op.GeoM.Translate(float64(10+280*b), float64(340+row*16))
				screen.DrawImage(rowView, op)
				wg.Done()
			}(row)
		}
		wg.Wait()
	}
}
