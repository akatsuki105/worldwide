package video

const (
	HORIZONTAL_PIXELS     = 160
	VERTICAL_PIXELS       = 144
	MODE_2_LENGTH         = 80
	MODE_3_LENGTH         = 172
	MODE_0_LENGTH         = 204
	HORIZONTAL_LENGTH     = 456
	TOTAL_LENGTH          = 70224
	MAX_OBJ               = 40
	MAX_LINE_OBJ          = 10
	VERTICAL_TOTAL_PIXELS = 154
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
	// Interrupts
	GB_REG_IF = 0x0F
	GB_REG_IE = 0xFF

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

	GB_REG_KEY0  = 0x4C
	GB_REG_KEY1  = 0x4D
	GB_REG_VBK   = 0x4F
	GB_REG_BANK  = 0x50
	GB_REG_HDMA1 = 0x51
	GB_REG_HDMA2 = 0x52
	GB_REG_HDMA3 = 0x53
	GB_REG_HDMA4 = 0x54
	GB_REG_HDMA5 = 0x55
	GB_REG_RP    = 0x56
	GB_REG_BCPS  = 0x68
	GB_REG_BCPD  = 0x69
	GB_REG_OCPS  = 0x6A
	GB_REG_OCPD  = 0x6B
	GB_REG_OPRI  = 0x6C
	GB_REG_SVBK  = 0x70
	GB_REG_UNK72 = 0x72
	GB_REG_UNK73 = 0x73
	GB_REG_UNK74 = 0x74
	GB_REG_UNK75 = 0x75
	GB_REG_PCM12 = 0x76
	GB_REG_PCM34 = 0x77
	GB_REG_MAX   = 0x100
)
