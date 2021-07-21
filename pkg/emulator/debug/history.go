package debug

import (
	"net/http"

	"github.com/pokemium/worldwide/pkg/gbc"
	"github.com/pokemium/worldwide/pkg/util"
)

func (d *Debugger) Hisotry(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		d.getHistory(w, req)
	case "POST":
		d.postHistory(w, req)
	default:
		http.NotFound(w, req)
	}
}

func getHistory(w http.ResponseWriter, req *http.Request) uint16 {
	return getU16FromQuery(w, req, "history")
}

func (d *Debugger) getHistory(w http.ResponseWriter, req *http.Request) {
	result := ""
	for _, inst := range d.history {
		result += stringfyCurInst(inst) + "\n"
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

func (d *Debugger) postHistory(w http.ResponseWriter, req *http.Request) {
	length := getHistory(w, req)
	if length == 0 {
		d.history = []gbc.CurInst{}
		d.g.Callbacks = util.RemoveCallback(d.g.Callbacks, "history")
		return
	}
	if length > 100 {
		http.Error(w, "history's count cannot be greater than 100", http.StatusBadRequest)
		return
	}

	d.history = make([]gbc.CurInst, length)
	d.g.Callbacks, _ = util.SetCallback(d.g.Callbacks, "history", util.PRIO_HISTORY, d.putHistory)
}

func (d *Debugger) putHistory() bool {
	curInst := d.g.Inst
	if curInst.PC == d.history[len(d.history)-1].PC && curInst.Opcode == d.history[len(d.history)-1].Opcode {
		return false
	}

	for i, h := range d.history {
		if h.Opcode == 0 && h.PC == 0 {
			d.history[i] = curInst
			return false
		}

		if i == len(d.history)-1 {
			d.history[i] = curInst
		} else {
			d.history[i] = d.history[i+1]
		}
	}
	return false
}
