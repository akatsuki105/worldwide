package debug

import (
	"fmt"
	"net/http"
	"strings"
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
	addr := getAddr(req)

	alreadyExists := false
	for _, bk := range d.Breakpoints {
		if addr == bk {
			alreadyExists = true
			break
		}
	}

	if !alreadyExists {
		d.Breakpoints = append(d.Breakpoints, addr)
	}
}

func (d *Debugger) deleteBreakpoint(w http.ResponseWriter, req *http.Request) {
	addr := getAddr(req)

	newBreakpoints := make([]uint16, 0)
	for _, bk := range d.Breakpoints {
		if addr != bk {
			newBreakpoints = append(newBreakpoints, bk)
		}
	}

	d.Breakpoints = newBreakpoints
}
