package debug

import (
	"net/http"
	"strconv"

	"github.com/pokemium/worldwide/pkg/util"
)

func (d *Debugger) Trace(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		if !*d.pause {
			http.Error(w, "trace API is available on pause state", http.StatusBadRequest)
			return
		}

		q := req.URL.Query()
		steps := uint16(0)
		for key, val := range q {
			if key == "step" {
				s, _ := strconv.ParseUint(val[0], 10, 16)
				steps = uint16(s)
			}
		}
		if steps == 0 {
			http.Error(w, "`step` is needed on query parameter(e.g. ?step=20)", http.StatusBadRequest)
			return
		}

		result := ""
		for s := uint16(0); s < steps; s++ {
			d.g.Step()
			for _, callback := range d.g.Callbacks {
				if callback.Priority == util.PRIO_BREAKPOINT {
					continue
				}
				if callback.Func() {
					break
				}
			}
			result += stringfyCurInst(d.g.Inst) + "\n"
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(result))
	default:
		http.NotFound(w, req)
		return
	}
}
