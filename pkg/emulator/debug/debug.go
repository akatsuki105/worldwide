package debug

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pokemium/worldwide/pkg/gbc"
)

const (
	TILE_PER_ROW = 16
)

type Debugger struct {
	g           *gbc.GBC
	pause       *bool
	Breakpoints []uint16
	history     []gbc.CurInst
}

func New(g *gbc.GBC, pause *bool) *Debugger {
	return &Debugger{
		g:     g,
		pause: pause,
	}
}

func (d *Debugger) Reset(g *gbc.GBC) {
	d.g = g
	d.history = make([]gbc.CurInst, len(d.history))
}

func getAddr(w http.ResponseWriter, req *http.Request) uint16 {
	return getU16FromQuery(w, req, "addr")
}

func getU16FromQuery(w http.ResponseWriter, req *http.Request, queryKey string) uint16 {
	switch req.Method {
	case "GET", "DELETE":
		q := req.URL.Query()
		addr := uint16(0)
		for key, val := range q {
			if key == queryKey {
				hexString := val[0]
				if !strings.HasPrefix(hexString, "0x") {
					http.Error(w, fmt.Sprintf("query parameter('%s') must be hex value(e.g. 0x486)", queryKey), http.StatusBadRequest)
					return 0
				}
				a, _ := strconv.ParseUint(hexString[2:], 16, 16)
				addr = uint16(a)
			}
		}
		return addr
	case "POST":
		body, _ := io.ReadAll(req.Body)
		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal)
		hexString := keyVal[queryKey]
		if !strings.HasPrefix(hexString, "0x") {
			http.Error(w, fmt.Sprintf("request parameter('%s') must be hex value(e.g. 0x486)", queryKey), http.StatusBadRequest)
			return 0
		}

		addr := uint16(0)
		a, _ := strconv.ParseUint(hexString[2:], 16, 16)
		addr = uint16(a)
		return addr
	}
	return 0
}

func stringfyCurInst(i gbc.CurInst) string {
	return fmt.Sprintf("0x%04x: %s", i.PC, disasm(i.Opcode))
}
