package gpu

type OAM [40]Obj

// GBObj
type Obj struct {
	y, x, tile, attr byte
}

type Sprite struct {
	obj   Obj
	index int8
}
