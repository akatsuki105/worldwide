package emulator

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

func (e *Emulator) RunServer(port int) {
	http.HandleFunc("/pause", e.Pause)
	http.HandleFunc("/continue", e.Continue)
	http.HandleFunc("/mute", e.toggleSound)
	http.HandleFunc("/debug/status", e.debugger.Status)
	http.HandleFunc("/debug/cartridge", e.debugger.Cartridge)
	http.HandleFunc("/debug/read1", e.debugger.Read1)
	http.HandleFunc("/debug/read2", e.debugger.Read2)
	http.Handle("/debug/tileview/bank0", websocket.Handler(e.debugger.TileView0))
	http.Handle("/debug/tileview/bank1", websocket.Handler(e.debugger.TileView1))
	http.Handle("/debug/sprview", websocket.Handler(e.debugger.SprView))
	http.Handle("/debug/io", websocket.Handler(e.debugger.IO))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (e *Emulator) Pause(w http.ResponseWriter, req *http.Request)    { e.pause = true }
func (e *Emulator) Continue(w http.ResponseWriter, req *http.Request) { e.pause = false }

func (e *Emulator) toggleSound(w http.ResponseWriter, req *http.Request) {
	e.GBC.Sound.Enable = !e.GBC.Sound.Enable
}
