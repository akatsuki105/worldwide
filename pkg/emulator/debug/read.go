package debug

import (
	"fmt"
	"net/http"
)

func (d *Debugger) Read1(w http.ResponseWriter, req *http.Request) {
	val := d.g.Load8(getAddr(w, req))

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("0x%02x", val)))
}

func (d *Debugger) Read2(w http.ResponseWriter, req *http.Request) {
	addr := getAddr(w, req)
	lower := uint16(d.g.Load8(addr))
	upper := uint16(d.g.Load8(addr + 1))
	val := upper<<8 | lower

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("0x%04x", val)))
}
