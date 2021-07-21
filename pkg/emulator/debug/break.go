package debug

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pokemium/worldwide/pkg/util"
)

func (d *Debugger) Break(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		d.getBreakpoints(w, req)
	case "POST":
		d.postBreakpoint(w, req)
	case "DELETE":
		d.deleteBreakpoint(w, req)
	}
}

func (d *Debugger) getBreakpoints(w http.ResponseWriter, req *http.Request) {
	result := "["
	for _, bk := range d.Breakpoints {
		result += fmt.Sprintf("0x%04x, ", bk)
	}
	result = strings.TrimRight(result, ", ")
	result += "]"

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

func (d *Debugger) postBreakpoint(w http.ResponseWriter, req *http.Request) {
	addr := getAddr(w, req)

	alreadyExists := false
	for _, bk := range d.Breakpoints {
		if addr == bk {
			alreadyExists = true
			break
		}
	}

	if !alreadyExists {
		d.Breakpoints = append(d.Breakpoints, addr)
		d.g.Callbacks = util.RemoveCallback(d.g.Callbacks, "break")
		d.g.Callbacks, _ = util.SetCallback(d.g.Callbacks, "break", util.PRIO_BREAKPOINT, d.checkBreakpoint)
	}
}

func (d *Debugger) deleteBreakpoint(w http.ResponseWriter, req *http.Request) {
	addr := getAddr(w, req)

	newBreakpoints := make([]uint16, 0)
	for _, bk := range d.Breakpoints {
		if addr != bk {
			newBreakpoints = append(newBreakpoints, bk)
		}
	}

	d.Breakpoints = newBreakpoints

	d.g.Callbacks = util.RemoveCallback(d.g.Callbacks, "break")
	d.g.Callbacks, _ = util.SetCallback(d.g.Callbacks, "break", util.PRIO_BREAKPOINT, d.checkBreakpoint)
}

func (d *Debugger) checkBreakpoint() bool {
	for _, bk := range d.Breakpoints {
		if d.g.Reg.PC == bk {
			*d.pause = true
			return true
		}
	}
	return false
}
