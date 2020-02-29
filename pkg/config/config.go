package config

import (
	"os"

	"gopkg.in/ini.v1"
)

// Init init config
func Init() *ini.File {
	exist := checkConfigFileExist()
	if exist {
		cfg, _ := ini.Load("worldwide.ini")
		return cfg
	}

	// create new config file
	cfg := ini.Empty()

	// display config
	cfg.Section("display").Key("expand").SetValue("2")
	cfg.Section("display").Key("smooth").SetValue("true")

	// DMG pallete color
	cfg.Section("pallete").Key("color0").SetValue("175,197,160")
	cfg.Section("pallete").Key("color1").SetValue("93,147,66")
	cfg.Section("pallete").Key("color2").SetValue("22,63,48")
	cfg.Section("pallete").Key("color3").SetValue("0,40,0")

	// network config
	cfg.Section("network").Key("network").SetValue("false")
	cfg.Section("network").Key("your").SetValue("localhost:8888")
	cfg.Section("network").Key("peer").SetValue("localhost:9999")

	// save config
	cfg.SaveTo("worldwide.ini")

	return cfg
}

// check ini file exists
func checkConfigFileExist() bool {
	filename := "worldwide.ini"
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}
