package gpu

import "gbc/pkg/util"

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

// GBVideoSoftwareRenderer
type Renderer struct {
	g                                 *GPU
	disableBG, disableOBJ, disableWIN bool
	highlightBG                       bool
	highlightOBJ                      [GB_VIDEO_MAX_OBJ]bool
	highlightWIN                      bool
	highlightColor                    uint16
	highlightAmount                   byte

	row     [168]uint16
	palette [192]uint16
	lookup  [192]byte

	scy, scx, wy, wx, currentWy, currentWx byte
	lastY, lastX                           int
	hasWindow                              bool

	lastHighlightAmount byte
	model               util.GBModel
	obj                 [GB_VIDEO_MAX_LINE_OBJ]Sprite
	objMax              int

	objOffsetX, objOffsetY, offsetScx, offsetScy, offsetWx, offsetWy int16
}

func NewRenderer() *Renderer {
	return &Renderer{
		lastY: GB_VIDEO_VERTICAL_PIXELS,
	}
}

func (r *Renderer) writePalette(index int, value uint16) {
	r.palette[index] = value
}

func (r *Renderer) writeVRAM(address uint16) {}
func (r *Renderer) writeOAM(oam uint16)      {}

// GBVideoSoftwareRendererDrawRange
func (r *Renderer) drawRange(startX, endX, y int) {
	r.lastY, r.lastX = y, endX
	if startX >= endX {
		return
	}

	mapIdx := GB_BASE_MAP // 0x9800
	if util.Bit(r.g.LCDC, TileMap) {
		mapIdx += GB_SIZE_MAP // 0x9c00
	}

	if r.disableBG {
		// TODO: memset(&softwareRenderer->row[startX], 0, (endX - startX) * sizeof(softwareRenderer->row[0]));
	}

	if util.Bit(r.g.LCDC, BgEnable) || r.model >= util.GB_MODEL_CGB {
		wy, wx := int(r.wy+r.currentWy), int(r.wx+r.currentWx)-7
		if util.Bit(r.g.LCDC, Window) && wy == y && wx <= endX {
			r.hasWindow = true
		}
		if util.Bit(r.g.LCDC, Window) && r.hasWindow && wx <= endX && !r.disableWIN {
			if wx > 0 && !r.disableBG {
				r.drawBackground(mapIdx, startX, wx, int(r.scx)-int(r.offsetScx), int(r.scy)+y-int(r.offsetScy), r.highlightBG)
			}

			mapIdx = GB_BASE_MAP
			if util.Bit(r.g.LCDC, TileMap) {
				mapIdx += GB_SIZE_MAP // 0x9c00
			}
			r.drawBackground(mapIdx, wx, endX, -wx-int(r.offsetWx), y-wy-int(r.offsetWy), r.highlightWIN)
		} else if !r.disableBG {
			r.drawBackground(mapIdx, startX, endX, int(r.scx)-int(r.offsetScx), int(r.scy)+y-int(r.offsetScy), r.highlightBG)
		}
	} else if !r.disableBG {
		// TODO: memset(&softwareRenderer->row[startX], 0, (endX - startX) * sizeof(softwareRenderer->row[0]));
	}

	if startX == 0 {
		r.cleanOAM(y)
	}
	if util.Bit(r.g.LCDC, ObjEnable) && !r.disableOBJ {
		for i := 0; i < r.objMax; i++ {
			r.drawObj(&Sprite{
				obj:   r.g.oam[i],
				index: int8(i),
			}, startX, endX, y)
		}
	}

	highlightAmount := (r.highlightAmount + 6) >> 4
	if r.lastHighlightAmount != highlightAmount {
		r.lastHighlightAmount = highlightAmount
	}
}

func (r *Renderer) drawBackground(mapIdx, startX, endX, sx, sy int, highlight bool) {

}

func (r *Renderer) drawObj(obj *Sprite, startX, endX, y int) {
	objX := int(obj.obj.x) + int(r.objOffsetX)
	ix := objX - 8
	if endX < ix || startX >= ix+8 {
		return
	}
	if objX < endX {
		endX = objX
	}
	if objX-8 > startX {
		startX = objX - 8
	}
	if startX < 0 {
		startX = 0
	}

	bank := 0
	tileOffset, bottomY := 0, 0
	objY := int(obj.obj.y) + int(r.objOffsetY)
	if util.Bit(obj.obj.attr, ObjAttrYFlip) {
		bottomY = 7 - ((y - objY - 16) & 7)
		if util.Bit(r.g.LCDC, ObjSize) && y-objY < -8 {
			tileOffset++
		}
	} else {
		bottomY = (y - objY - 16) & 7
		if util.Bit(r.g.LCDC, ObjSize) && y-objY >= -8 {
			tileOffset++
		}
	}
	if util.Bit(r.g.LCDC, ObjSize) && obj.obj.tile&1 == 1 {
		tileOffset--
	}

	mask, mask2 := uint(0x60), uint(0x100/3)
	if util.Bit(obj.obj.attr, ObjAttrPriority) {
		mask, mask2 = 0x63, 0
	}

	p := uint16(0x20)
	if r.highlightOBJ[obj.index] {
		p = 0x80 | 0x20
	}
	if r.model == util.GB_MODEL_CGB {
		p |= uint16(obj.obj.attr&0x07) * 4
		if util.Bit(obj.obj.attr, ObjAttrBank) {
			bank = 1
		}
		if !util.Bit(r.g.LCDC, BgEnable) {
			mask, mask2 = 0x60, 0x100/3
		}
	} else {
		p |= (uint16(obj.obj.attr&(1<<ObjAttrPalette)) + 8) * 4
	}

	bottomX, x, objTile := 0, startX, int(obj.obj.tile)+tileOffset
	if (x-objX)&7 != 0 {
		for ; x < endX; x++ {
			if util.Bit(obj.obj.attr, ObjAttrXFlip) {
				bottomX = (x - objX) & 7
			} else {
				bottomX = 7 - ((x - objX) & 7)
			}
			tileDataLower := r.g.VRAM.Bank[bank][(objTile*8+bottomY)*2]
			tileDataUpper := r.g.VRAM.Bank[bank][(objTile*8+bottomY)*2+1]
			tileDataUpper >>= bottomX
			tileDataLower >>= bottomX
			current := r.row[x]
			if ((tileDataUpper|tileDataLower)&1 > 0) && (uint(current)&mask == 0) && (uint(current)&mask2) <= OBJ_PRIORITY {
				r.row[x] = p | uint16((tileDataUpper&1)<<1) | uint16(tileDataLower&1)
			}
		}
	} else if util.Bit(obj.obj.attr, ObjAttrXFlip) {
		tileDataLower := r.g.VRAM.Bank[bank][(objTile*8+bottomY)*2]
		tileDataUpper := r.g.VRAM.Bank[bank][(objTile*8+bottomY)*2+1]
		current := r.row[x]
		if ((tileDataUpper|tileDataLower)&1) != 0 && (uint(current)&mask == 0) && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x] = p | uint16((tileDataUpper&1)<<1) | uint16(tileDataLower&1)
		}
		current = r.row[x+1]
		if ((tileDataUpper|tileDataLower)&2) != 0 && (uint(current)&mask == 0) && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+1] = p | uint16(tileDataUpper&2) | uint16((tileDataLower&2)>>1)
		}
		current = r.row[x+2]
		if ((tileDataUpper|tileDataLower)&4) != 0 && (uint(current)&mask == 0) && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+2] = p | uint16((tileDataUpper&4)>>1) | uint16((tileDataLower&4)>>2)
		}
		current = r.row[x+3]
		if ((tileDataUpper|tileDataLower)&8) != 0 && (uint(current)&mask == 0) && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+3] = p | uint16((tileDataUpper&8)>>2) | uint16((tileDataLower&8)>>3)
		}
		current = r.row[x+4]
		if ((tileDataUpper|tileDataLower)&16) != 0 && (uint(current)&mask == 0) && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+4] = p | uint16((tileDataUpper&16)>>3) | uint16((tileDataLower&16)>>4)
		}
		current = r.row[x+5]
		if ((tileDataUpper|tileDataLower)&32) != 0 && (uint(current)&mask == 0) && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+5] = p | uint16((tileDataUpper&32)>>4) | uint16((tileDataLower&32)>>5)
		}
		current = r.row[x+6]
		if ((tileDataUpper|tileDataLower)&64) != 0 && (uint(current)&mask == 0) && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+6] = p | uint16((tileDataUpper&64)>>5) | uint16((tileDataLower&64)>>6)
		}
		current = r.row[x+7]
		if ((tileDataUpper|tileDataLower)&128) != 0 && (uint(current)&mask == 0) && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+7] = p | uint16((tileDataUpper&128)>>6) | uint16((tileDataLower&128)>>7)
		}
	} else {
		tileDataLower := r.g.VRAM.Bank[bank][(objTile*8+bottomY)*2]
		tileDataUpper := r.g.VRAM.Bank[bank][(objTile*8+bottomY)*2+1]
		current := r.row[x+7]
		if ((tileDataUpper|tileDataLower)&1) != 0 && (uint(current)&mask) == 0 && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+7] = p | uint16((tileDataUpper&1)<<1) | uint16(tileDataLower&1)
		}
		current = r.row[x+6]
		if ((tileDataUpper|tileDataLower)&2) != 0 && (uint(current)&mask) == 0 && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+6] = p | uint16(tileDataUpper&2) | uint16((tileDataLower&2)>>1)
		}
		current = r.row[x+5]
		if ((tileDataUpper|tileDataLower)&4) != 0 && (uint(current)&mask) == 0 && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+5] = p | uint16((tileDataUpper&4)>>1) | uint16((tileDataLower&4)>>2)
		}
		current = r.row[x+4]
		if ((tileDataUpper|tileDataLower)&8) != 0 && (uint(current)&mask) == 0 && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+4] = p | uint16((tileDataUpper&8)>>2) | uint16((tileDataLower&8)>>3)
		}
		current = r.row[x+3]
		if ((tileDataUpper|tileDataLower)&16) != 0 && (uint(current)&mask) == 0 && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+3] = p | uint16((tileDataUpper&16)>>3) | uint16((tileDataLower&16)>>4)
		}
		current = r.row[x+2]
		if ((tileDataUpper|tileDataLower)&32) != 0 && (uint(current)&mask) == 0 && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+2] = p | uint16((tileDataUpper&32)>>4) | uint16((tileDataLower&32)>>5)
		}
		current = r.row[x+1]
		if ((tileDataUpper|tileDataLower)&64) != 0 && (uint(current)&mask) == 0 && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x+1] = p | uint16((tileDataUpper&64)>>5) | uint16((tileDataLower&64)>>6)
		}
		current = r.row[x]
		if ((tileDataUpper|tileDataLower)&128) != 0 && (uint(current)&mask) == 0 && (uint(current)&mask2) <= OBJ_PRIORITY {
			r.row[x] = p | uint16((tileDataUpper&128)>>6) | uint16((tileDataLower&128)>>7)
		}
	}
}

func (r *Renderer) cleanOAM(y int) {
	spriteHeight := 8
	if util.Bit(r.g.LCDC, ObjSize) {
		spriteHeight = 16
	}

	o := 0
	for i := 0; i < GB_VIDEO_MAX_OBJ && o < GB_VIDEO_MAX_LINE_OBJ; i++ {
		oy := int(r.g.oam[i].y)
		if y < oy-16 || y >= oy-16+spriteHeight {
			continue
		}

		r.obj[o].obj = r.g.oam[i]
		r.obj[o].index = int8(i)
		o++
		if o == 10 {
			break
		}
	}
	r.objMax = o
}
