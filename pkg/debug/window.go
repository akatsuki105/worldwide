package debug

type Window struct {
	x int
	y int
}

func (w *Window) Size() (float64, float64) {
	return float64(w.x), float64(w.y)
}

func (w *Window) SetSize(x, y int) {
	w.x = x
	w.y = y
}
