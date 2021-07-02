package gpu

type EntryY struct {
	Block  int
	Offset int
}

// SetBGLine 1タイルライン描画する
func (g *GPU) SetBGLine(entryX int, entryY EntryY, tileX, tileY uint, isWin, isCGB bool, lineIdx int) bool {
	return false
}
