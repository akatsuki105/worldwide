package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

const (
	tomlName = "worldwide.toml"
)

// Config for emulator
type Config struct {
	Display Display `toml:"display"`
	Pallete Pallete `toml:"pallete"`
	Network Network `toml:"network"`
	Joypad  Joypad  `toml:"joypad"`
	Debug   Debug   `toml:"debug"`
}

// Display config
type Display struct {
	HQ2x  bool `toml:"hq2x"`  // エミュレータのハイレゾ化が有効かどうか
	FPS30 bool `toml:"fps30"` // fpsを30に下げるモードかどうか
}

// Pallete for DMG
type Pallete struct {
	Color0 [3]int `toml:"color0"`
	Color1 [3]int `toml:"color1"`
	Color2 [3]int `toml:"color2"`
	Color3 [3]int `toml:"color3"`
}

// Network config
type Network struct {
	Network bool   `toml:"network"`
	Your    string `toml:"your"`
	Peer    string `toml:"peer"`
}

// Joypad config
type Joypad struct {
	A         uint    `toml:"A"`
	B         uint    `toml:"B"`
	Start     uint    `toml:"Start"`
	Select    uint    `toml:"Select"`
	Threshold float64 `toml:"threshold"`
}

// Debug config
type Debug struct {
	BreakPoints []string `toml:"breakpoints"`
	History     bool     `toml:"history"`
}

func Init() *Config {
	cfg := &Config{}

	// load config
	if _, err := toml.DecodeFile(tomlName, cfg); err == nil {
		return cfg
	}

	// create config
	cfgText := `[display]
hq2x = false # use HQ2x scaling mode
fps30 = false # reduce fps 30

[pallete]
# DMG Color Pallete [R, G, B]
color0 = [175, 197, 160]
color1 = [93, 147, 66]
color2 = [22, 63, 48]
color3 = [0, 40, 0]

[network]
network = false
your = "127.0.0.1:8888"
peer = "127.0.0.1:9999"

[joypad]
A = 1
B = 0
Start = 7
Select = 6
threshold = 0.7 # How reactive axis is

[debug]
# "BANK:PC;Cond" e.g. "00:0460;SP==c0f3", "01:ffff;"
breakpoints = []
history = false # history uses a lot of CPU resource
`
	ioutil.WriteFile(tomlName, []byte(cfgText), 0666)
	toml.Decode(cfgText, cfg)
	return cfg
}
