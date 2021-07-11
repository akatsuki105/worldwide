package video

type OAM struct {
	Objs   [40]*Obj
	Buffer [0xa0]byte
}

func NewOAM() *OAM {
	o := &OAM{}
	for i := 0; i < 40; i++ {
		o.Objs[i] = &Obj{}
	}
	return o
}

func (o *OAM) Get(offset uint16) byte {
	obj := offset >> 2
	idx := offset & 3
	return o.Objs[obj].get(idx)
}

func (o *OAM) Set(offset uint16, value byte) {
	obj := offset >> 2
	idx := offset & 3
	o.Objs[obj].set(idx, value)
}

// GBObj
type Obj struct {
	y, x, tile, attr byte
}

func (o *Obj) get(idx uint16) byte {
	switch idx {
	case 0:
		return o.y
	case 1:
		return o.x
	case 2:
		return o.tile
	case 3:
		return o.attr
	}
	return 0xff
}

func (o *Obj) set(idx uint16, value byte) {
	switch idx {
	case 0:
		o.y = value
	case 1:
		o.x = value
	case 2:
		o.tile = value
	case 3:
		o.attr = value
	}
}

type Sprite struct {
	obj   Obj
	index int8
}
