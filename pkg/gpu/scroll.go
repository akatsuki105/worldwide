package gpu

// GetScroll - スクロール値を得る
func (g *GPU) GetScroll() (x, y uint) {
	x, y = uint(g.Scroll[0]), uint(g.Scroll[1])
	return x, y
}

// SetScrollX - スクロールのX座標を書き込む
func (g *GPU) SetScrollX(x byte) {
	g.Scroll[0] = x
}

// SetScrollY - スクロールのY座標を書き込む
func (g *GPU) SetScrollY(y byte) {
	g.Scroll[1] = y
}
