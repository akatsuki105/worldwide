package audio

import (
	"github.com/hajimehoshi/oto"
	"github.com/pokemium/worldwide/pkg/gbc/apu"
)

var context *oto.Context
var player *oto.Player
var Stream []byte
var enable *bool

func Reset(enablePtr *bool) {
	enable = enablePtr

	var err error
	context, err = oto.NewContext(apu.SAMPLE_RATE, 2, 1, apu.SAMPLE_RATE/apu.BUF_SEC)
	if err != nil {
		panic(err)
	}

	player = context.NewPlayer()
}

func Play() {
	if player == nil || !*enable {
		return
	}
	player.Write(Stream)
}

func SetStream(b []byte) { Stream = b }
