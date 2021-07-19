package debug

import (
	"log"
	"time"

	"golang.org/x/net/websocket"
)

func (d *Debugger) IO(ws *websocket.Conn) {
	err := websocket.Message.Send(ws, d.g.IO[:])
	if err != nil {
		log.Printf("error sending data: %v\n", err)
		return
	}

	for range time.NewTicker(time.Millisecond * 100).C {
		err := websocket.Message.Send(ws, d.g.IO[:])
		if err != nil {
			log.Printf("error sending data: %v\n", err)
			return
		}
	}
}
