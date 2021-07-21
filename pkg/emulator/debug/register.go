package debug

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pokemium/worldwide/pkg/gbc"
	"github.com/pokemium/worldwide/pkg/util"
)

type Register struct {
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
	IE, IF      string
	IME         string
	Halt        string
	DoubleSpeed string
}

var regs = map[string]int{
	"a": 0,
	"b": 1,
	"c": 2,
	"d": 3,
	"e": 4,
	"h": 5,
	"l": 6,
	"f": 7,
}

func (d *Debugger) Register(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		r := Register{
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
			IE:          fmt.Sprintf("0x%02x", d.g.IO[gbc.IEIO]),
			IF:          fmt.Sprintf("0x%02x", d.g.IO[gbc.IFIO]),
			IME:         fmt.Sprintf("0x%02x", util.Bool2U8(d.g.Reg.IME)),
			Halt:        fmt.Sprintf("%+v", d.g.Halt),
			DoubleSpeed: fmt.Sprintf("%+v", d.g.DoubleSpeed),
		}

		res, err := json.Marshal(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	case "POST":
		body, _ := io.ReadAll(req.Body)
		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal)

		target := strings.ToLower(keyVal["target"])
		val := strings.ToLower(keyVal["value"])
		if !strings.HasPrefix(val, "0x") {
			http.Error(w, "value must be hexadecimal (e.g. 0x04)", http.StatusBadRequest)
			return
		}

		switch target {
		case "a", "b", "c", "d", "e", "h", "l", "f":
			a, _ := strconv.ParseUint(val[2:], 16, 8)
			val := byte(a)
			d.g.Reg.R[regs[target]] = val
		case "pc":
			a, _ := strconv.ParseUint(val[2:], 16, 16)
			val := uint16(a)
			d.g.Reg.PC = val
		case "sp":
			a, _ := strconv.ParseUint(val[2:], 16, 16)
			val := uint16(a)
			d.g.Reg.SP = val
		case "ime":
			val, _ := strconv.ParseUint(val[2:], 16, 1) // 0x0 or 0x1
			d.g.Reg.IME = val > 0
		default:
			http.Error(w, "invalid target", http.StatusBadRequest)
		}
	default:
		http.NotFound(w, req)
		return
	}
}
