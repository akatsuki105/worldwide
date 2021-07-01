package util

type GBModel byte

const (
	GB_MODEL_AUTODETECT GBModel = 0xFF
	GB_MODEL_DMG        GBModel = 0x00
	GB_MODEL_SGB        GBModel = 0x20
	GB_MODEL_MGB        GBModel = 0x40
	GB_MODEL_SGB2       GBModel = 0x60
	GB_MODEL_CGB        GBModel = 0x80
	GB_MODEL_AGB        GBModel = 0xC0
)
