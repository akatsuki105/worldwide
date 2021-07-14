package emulator

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var DEBUG_BG = color.RGBA{35, 27, 167, 0xff}

const (
	DEBUG_BG_X = 1080
	DEBUG_BG_Y = 740
)

func (e *Emulator) drawDebugScreen() *ebiten.Image {
	screen := ebiten.NewImage(DEBUG_BG_X, DEBUG_BG_Y)
	screen.Fill(DEBUG_BG)

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

	// sprite view
	e.drawDebugSpriteView(screen)

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
	buffer := make([]byte, len(banks[0]))

	// each row has 16 tiles
	for b := 0; b < 2; b++ {
		if b == 1 && !e.GBC.Cartridge.IsCGB {
			tileView := ebiten.NewImage(8*TILE_PER_ROW, 8*384/TILE_PER_ROW)
			tileView.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(10+280), float64(340))
			screen.DrawImage(tileView, op)
		}

		var wg sync.WaitGroup
		wg.Add(384 / TILE_PER_ROW)
		for row := 0; row < 384/TILE_PER_ROW; row++ {
			// 0..63, 0..63, 0..63, .. -> 0..7, 0..7, ... 8..15, 8..15,
			go func(row int) {
				rowStart, rowEnd := row*TILE_PER_ROW, (row+1)*TILE_PER_ROW
				bufferIdx := (TILE_PER_ROW * 64 * 4) * row
				rowBuffer := banks[b][rowStart*64*4 : rowEnd*64*4]

				for t := 0; t < TILE_PER_ROW; t++ {
					rowBufferBase := t * 64 * 4
					for y := 0; y < 8; y++ {
						tileRowBuffer := rowBuffer[rowBufferBase+y*8*4 : rowBufferBase+(y+1)*8*4] // (y*8)..((y*8)+7)
						for x := 0; x < 8; x++ {
							idx := bufferIdx + y*8*TILE_PER_ROW*4 + t*8*4 + x*4
							buffer[idx] = tileRowBuffer[x*4]
							buffer[idx+1] = tileRowBuffer[x*4+1]
							buffer[idx+2] = tileRowBuffer[x*4+2]
							buffer[idx+3] = tileRowBuffer[x*4+3]
						}
					}
				}
				wg.Done()
			}(row)
		}
		wg.Wait()

		tileView := ebiten.NewImage(8*TILE_PER_ROW, 8*384/TILE_PER_ROW) // [8x16, 8]
		tileView.ReplacePixels(buffer)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(float64(10+280*b), float64(340))
		screen.DrawImage(tileView, op)
	}
}

func (e *Emulator) drawDebugSpriteView(screen *ebiten.Image) {
	startX := 570
	ebitenutil.DebugPrintAt(screen, "OAM (Y, X, tile, attr)", startX, 320)

	buffers := e.debugger.SprView()
	var wg sync.WaitGroup
	wg.Add(40)
	for i := 0; i < 40; i++ {
		go func(i int) {
			x, y := i&0x7, i/8
			sprView := ebiten.NewImage(8, 8)
			sprView.ReplacePixels(buffers[i][:])
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(4, 4)
			op.GeoM.Translate(float64(startX+x*48), float64(340+y*48))
			screen.DrawImage(sprView, op)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
