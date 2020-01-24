package boot

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// ブートアニメーションを管理する
type Boot struct {
	Table  map[uint]uint // 何フレーム目に何枚目の画像を表示するか
	Frame  uint          // 現在何フレーム目か
	Images [31]image.Image
	Length uint // このブートに何フレーム使うか
	valid  bool // ブートROMを実行するか
}

func New(valid bool) *Boot {
	boot := &Boot{
		valid: valid,
	}

	if valid {
		// imgをInit
		for i := 0; i < 31; i++ {
			filename := fmt.Sprintf("asset/%d.png", i)
			b, err := Asset(filename)
			if err != nil {
				panic(err)
			}
			img, _, _ := image.Decode(bytes.NewReader(b))
			boot.Images[i] = img
		}

		// Tableをinit
		boot.Table = map[uint]uint{}
		for i := 0; i < 11*3; i++ {
			boot.Table[uint(i)] = 0
		}
		for i := 11 * 3; i < 18*3; i++ {
			boot.Table[uint(i)] = 1
		}
		for i := 18; i < 36; i++ {
			for j := 0; j < 3; j++ {
				index := i*3 + j
				boot.Table[uint(index)] = 2 + uint(i) - 18
			}
		}
		for i := 36 * 3; i < 48*3; i++ {
			boot.Table[uint(i)] = 20
		}
		for i := 48; i < 57; i++ {
			for j := 0; j < 3; j++ {
				index := i*3 + j
				boot.Table[uint(index)] = 21 + uint(i) - 48
			}
		}
		for i := 57 * 3; i <= 60*3; i++ {
			boot.Table[uint(i)] = 30
		}

		boot.Length = 180
	}

	return boot
}

func (boot *Boot) PlaySE() {
	if boot.valid {
		data, err := Asset("asset/se.mp3")
		if err != nil {
			panic(err)
		}
		r := ioutil.NopCloser(bytes.NewReader(data))

		streamer, format, err := mp3.Decode(r)
		if err != nil {
			panic(err)
		}
		defer streamer.Close()
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		speaker.Play(streamer)
	}
}

func (boot *Boot) ExitSE() {
	if boot.valid {
		speaker.Close()
	}
}

func (boot *Boot) Valid() bool {
	return boot.valid
}
