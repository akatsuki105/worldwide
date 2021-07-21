package emulator

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

func (e *Emulator) RunServer(port int) {
	http.HandleFunc("/pause", e.Pause)
	http.HandleFunc("/continue", e.Continue)
	http.HandleFunc("/reset", e.Reset)
	http.HandleFunc("/quit", e.Quit)
	http.HandleFunc("/mute", e.toggleSound)
	http.HandleFunc("/debug/register", e.debugger.Register)
	http.HandleFunc("/debug/break", e.debugger.Break)
	http.HandleFunc("/debug/cartridge", e.debugger.Cartridge)
	http.HandleFunc("/debug/read1", e.debugger.Read1)
	http.HandleFunc("/debug/read2", e.debugger.Read2)
	http.HandleFunc("/debug/disasm", e.debugger.Disasm)
	http.HandleFunc("/debug/trace", e.debugger.Trace)
	http.HandleFunc("/debug/history", e.debugger.Hisotry)
	http.Handle("/debug/tileview/bank0", websocket.Handler(e.debugger.TileView0))
	http.Handle("/debug/tileview/bank1", websocket.Handler(e.debugger.TileView1))
	http.Handle("/debug/sprview", websocket.Handler(e.debugger.SprView))
	http.Handle("/debug/io", websocket.Handler(e.debugger.IO))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// these requests accept get method, but have side effect ummmm....
func (e *Emulator) Pause(w http.ResponseWriter, req *http.Request)    { e.pause = true }
func (e *Emulator) Continue(w http.ResponseWriter, req *http.Request) { e.pause = false }
func (e *Emulator) Reset(w http.ResponseWriter, req *http.Request)    { e.reset = true }
func (e *Emulator) Quit(w http.ResponseWriter, req *http.Request)     { e.quit = true }
func (e *Emulator) toggleSound(w http.ResponseWriter, req *http.Request) {
	e.GBC.Sound.Enable = !e.GBC.Sound.Enable
}
