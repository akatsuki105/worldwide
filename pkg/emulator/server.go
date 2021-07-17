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
	http.HandleFunc("/debug/register", e.debugger.Register)
	http.HandleFunc("/debug/cartridge", e.debugger.Cartridge)
	http.Handle("/debug/tileview", websocket.Handler(e.debugger.TileView))
	http.Handle("/debug/sprview", websocket.Handler(e.debugger.SprView))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (e *Emulator) Pause(w http.ResponseWriter, req *http.Request)    { e.pause = true }
func (e *Emulator) Continue(w http.ResponseWriter, req *http.Request) { e.pause = false }

func (e *Emulator) toggleSound(w http.ResponseWriter, req *http.Request) {
	e.GBC.Sound.Enable = !e.GBC.Sound.Enable
}
