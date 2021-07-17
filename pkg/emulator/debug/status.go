package debug

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pokemium/worldwide/pkg/gbc"
	"github.com/pokemium/worldwide/pkg/util"
)

type Status struct {
	A           string
	F           string
	B           string
	C           string
	D           string
	E           string
	H           string
	L           string
	PC          string
	SP          string
	IE          string
	IF          string
	IME         string
	LCDC        string
	STAT        string
	LY          string
	LYC         string
	DoubleSpeed string
}

func (d *Debugger) Status(w http.ResponseWriter, req *http.Request) {
	r := Status{
		A:           fmt.Sprintf("0x%02x", d.g.Reg.R[gbc.A]),
		F:           fmt.Sprintf("0x%02x", d.g.Reg.R[gbc.F]),
		B:           fmt.Sprintf("0x%02x", d.g.Reg.R[gbc.B]),
		C:           fmt.Sprintf("0x%02x", d.g.Reg.R[gbc.C]),
		D:           fmt.Sprintf("0x%02x", d.g.Reg.R[gbc.D]),
		E:           fmt.Sprintf("0x%02x", d.g.Reg.R[gbc.E]),
		H:           fmt.Sprintf("0x%02x", d.g.Reg.R[gbc.H]),
		L:           fmt.Sprintf("0x%02x", d.g.Reg.R[gbc.L]),
		PC:          fmt.Sprintf("0x%04x", d.g.Reg.PC),
		SP:          fmt.Sprintf("0x%04x", d.g.Reg.SP),
		IE:          fmt.Sprintf("0x%04x", d.g.IO[gbc.IEIO]),
		IF:          fmt.Sprintf("0x%04x", d.g.IO[gbc.IFIO]),
		IME:         fmt.Sprintf("0x%02x", util.Bool2U8(d.g.Reg.IME)),
		LCDC:        fmt.Sprintf("0x%02x", d.g.IO[gbc.LCDCIO]),
		STAT:        fmt.Sprintf("0x%02x", d.g.IO[gbc.LCDSTATIO]),
		LY:          fmt.Sprintf("0x%02x", d.g.Video.Ly),
		LYC:         fmt.Sprintf("0x%02x", d.g.IO[gbc.LYCIO]),
		DoubleSpeed: fmt.Sprintf("%+v", d.g.DoubleSpeed),
	}

	res, err := json.Marshal(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
