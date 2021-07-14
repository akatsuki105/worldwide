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
}

// Display config
type Display struct {
	FPS30 bool `toml:"fps30"` // reduce fps to 30
}

func New() *Config {
	cfg := &Config{}

	// load config
	if _, err := toml.DecodeFile(tomlName, cfg); err == nil {
		return cfg
	}

	// create config
	cfgText := `[display]
fps30 = false # reduce fps 30
`
	ioutil.WriteFile(tomlName, []byte(cfgText), 0666)
	toml.Decode(cfgText, cfg)
	return cfg
}
