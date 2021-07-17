package debug

import (
	"fmt"
	"net/http"
	"strconv"
)

func getAddr(req *http.Request) uint16 {
	q := req.URL.Query()
	addr := uint16(0)
	for key, val := range q {
		if key == "addr" {
			a, _ := strconv.ParseUint(val[0][2:], 16, 16)
			addr = uint16(a)
		}
	}
	return addr
}

func (d *Debugger) Read1(w http.ResponseWriter, req *http.Request) {
	val := d.g.Load8(getAddr(req))

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("0x%02x", val)))
}

func (d *Debugger) Read2(w http.ResponseWriter, req *http.Request) {
	addr := getAddr(req)
	lower := uint16(d.g.Load8(addr))
	upper := uint16(d.g.Load8(addr + 1))
	val := upper<<8 | lower

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("0x%04x", val)))
}
