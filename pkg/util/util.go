package util

// Bit check val's idx bit
func Bit(val interface{}, idx int) bool {
	switch val := val.(type) {
	case uint64:
		if idx < 0 || idx > 63 {
			return false
		}
		return (val & (1 << idx)) != 0

	case uint32:
		if idx < 0 || idx > 31 {
			return false
		}
		return (val & (1 << idx)) != 0

	case uint:
		if idx < 0 || idx > 31 {
			return false
		}
		return (val & (1 << idx)) != 0

	case int:
		if idx < 0 || idx > 31 {
			return false
		}
		return (val & (1 << idx)) != 0

	case uint16:
		if idx < 0 || idx > 15 {
			return false
		}
		return (val & (1 << idx)) != 0

	case byte:
		if idx < 0 || idx > 7 {
			return false
		}
		return (val & (1 << idx)) != 0
	}
	return false
}

func SetBit16(val uint16, idx int, b bool) uint16 {
	if b {
		val = val | (1 << idx)
	} else {
		val = val & ^(1 << idx)
	}
	return val
}

func SetBit8(val byte, idx int, b bool) byte {
	if b {
		val = val | (1 << idx)
	} else {
		val = val & ^(1 << idx)
	}
	return val
}

func Bool2Int(b bool) int {
	if b {
		return 1
	}
	return 0
}
func Bool2U8(b bool) byte {
	if b {
		return 1
	}
	return 0
}
func Bool2U16(b bool) uint16 {
	if b {
		return 1
	}
	return 0
}
func Bool2U32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}
func Bool2U64(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

func SetMSB(val byte, b bool) byte {
	if b {
		val |= 0x80
	} else {
		val &= 0x7f
	}
	return val
}

func SetLSB(val byte, b bool) byte {
	if b {
		val |= 1
	} else {
		val &= 0xfe
	}
	return val
}
