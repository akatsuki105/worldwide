package debug

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pokemium/worldwide/pkg/gbc"
)

const (
	TILE_PER_ROW = 16
)

type Debugger struct {
	g           *gbc.GBC
	Breakpoints []uint16
}

func New(g *gbc.GBC) *Debugger {
	return &Debugger{
		g: g,
	}
}

func getAddr(req *http.Request) uint16 {
	switch req.Method {
	case "GET", "DELETE":
		q := req.URL.Query()
		addr := uint16(0)
		for key, val := range q {
			if key == "addr" {
				a, _ := strconv.ParseUint(val[0][2:], 16, 16)
				addr = uint16(a)
			}
		}
		return addr
	case "POST":
		body, _ := io.ReadAll(req.Body)
		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal)
		p := keyVal["addr"]

		addr := uint16(0)
		fmt.Println(p)
		a, _ := strconv.ParseUint(p[2:], 16, 16)
		addr = uint16(a)
		return addr
	}
	return 0
}
