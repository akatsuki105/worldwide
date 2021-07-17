package debug

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pokemium/worldwide/pkg/gbc"
)

type Register struct {
	A string
	F string
	B string
	C string
	D string
	E string
	H string
	L string
}

func (d *Debugger) Register(w http.ResponseWriter, req *http.Request) {
	A, F := d.g.Reg.R[gbc.A], d.g.Reg.R[gbc.F]
	B, C := d.g.Reg.R[gbc.B], d.g.Reg.R[gbc.C]
	D, E := d.g.Reg.R[gbc.D], d.g.Reg.R[gbc.E]
	H, L := d.g.Reg.R[gbc.H], d.g.Reg.R[gbc.L]

	r := Register{
		A: fmt.Sprintf("0x%02x", A),
		F: fmt.Sprintf("0x%02x", F),
		B: fmt.Sprintf("0x%02x", B),
		C: fmt.Sprintf("0x%02x", C),
		D: fmt.Sprintf("0x%02x", D),
		E: fmt.Sprintf("0x%02x", E),
		H: fmt.Sprintf("0x%02x", H),
		L: fmt.Sprintf("0x%02x", L),
	}

	res, err := json.Marshal(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
