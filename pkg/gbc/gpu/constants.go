package gpu

const (
	GB_VIDEO_VERTICAL_PIXELS = 144
	GB_VIDEO_MAX_OBJ         = 40
	GB_VIDEO_MAX_LINE_OBJ    = 10
)

const (
	GB_BASE_MAP        = 0x1800
	GB_SIZE_MAP        = 0x400
	GB_SIZE_VRAM_BANK0 = 0x2000
)

const (
	OBJ_PRIORITY   = 0x100
	OBJ_PRIO_MASK  = 0xff
	PAL_SGB_BORDER = 0x40
)

const (
	BgEnable = iota
	ObjEnable
	ObjSize
	TileMap
	TileData
	Window
	WindowTileMap
	Enable
)

const (
	ObjAttrBank = iota + 3
	ObjAttrPalette
	ObjAttrXFlip
	ObjAttrYFlip
	ObjAttrPriority
)
