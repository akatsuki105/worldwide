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

	for range time.NewTicker(time.Second).C {
		err := websocket.Message.Send(ws, d.g.IO[:])
		if err != nil {
			log.Printf("error sending data: %v\n", err)
			return
		}
	}
}
