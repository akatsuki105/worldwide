package gpu

const (
	GB_VIDEO_HORIZONTAL_PIXELS = 160
	GB_VIDEO_VERTICAL_PIXELS   = 144
	GB_VIDEO_MAX_OBJ           = 40
	GB_VIDEO_MAX_LINE_OBJ      = 10
)

const (
	GB_BASE_MAP        = 0x1800
	GB_SIZE_MAP        = 0x400
	GB_SIZE_VRAM_BANK0 = 0x2000
)

const (
	OBJ_PRIORITY  = 0x100
	OBJ_PRIO_MASK = 0xff
)

const (
	PAL_BG            = 0x0
	PAL_OBJ           = 0x20
	PAL_HIGHLIGHT     = 0x80
	PAL_HIGHLIGHT_BG  = PAL_HIGHLIGHT | PAL_BG
	PAL_HIGHLIGHT_OBJ = PAL_HIGHLIGHT | PAL_OBJ
	PAL_SGB_BORDER    = 0x40
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

const (
	GB_REG_LCDC = 0x40
	GB_REG_STAT = 0x41
	GB_REG_SCY  = 0x42
	GB_REG_SCX  = 0x43
	GB_REG_LY   = 0x44
	GB_REG_LYC  = 0x45
	GB_REG_DMA  = 0x46
	GB_REG_BGP  = 0x47
	GB_REG_OBP0 = 0x48
	GB_REG_OBP1 = 0x49
	GB_REG_WY   = 0x4A
	GB_REG_WX   = 0x4B
)
