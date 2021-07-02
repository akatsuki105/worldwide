package gpu

type OAM [40]*Obj

func NewOAM() *OAM {
	o := &OAM{}
	for i := 0; i < 40; i++ {
		o[i] = &Obj{}
	}
	return o
}

// GBObj
type Obj struct {
	y, x, tile, attr byte
}

func (o *Obj) Set(idx int, value byte) {
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
