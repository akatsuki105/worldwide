package util

func Bit(b byte, n int) int {
	switch {
	case n < 0:
		return 0
	case n > 7:
		return 0
	default:
		return int((b >> n) & 0x01)
	}
}

func Bool2Int(b bool) int {
	if b {
		return 1
	}
	return 0
}
