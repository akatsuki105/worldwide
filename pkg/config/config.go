package config

import "os"

import "gopkg.in/ini.v1"

// Init init config
func Init() *ini.File {
	exist := checkConfigFileExist()
	if exist {
		cfg, _ := ini.Load("gb.ini")
		return cfg
	}

	// 設定ファイルが存在しないとき
	cfg := ini.Empty()
	// display
	cfg.Section("display").Key("expand").SetValue("2")
	cfg.Section("display").Key("saving").SetValue("yes")

	// Xbox 360 Controller
	cfg.Section("Xbox 360 Controller").Key("A").SetValue("1")
	cfg.Section("Xbox 360 Controller").Key("B").SetValue("0")
	cfg.Section("Xbox 360 Controller").Key("Start").SetValue("7")
	cfg.Section("Xbox 360 Controller").Key("Select").SetValue("6")
	cfg.Section("Xbox 360 Controller").Key("Horizontal").SetValue("0")
	cfg.Section("Xbox 360 Controller").Key("Vertical").SetValue("1")
	cfg.Section("Xbox 360 Controller").Key("VerticalNegative").SetValue("0")
	cfg.Section("Xbox 360 Controller").Key("Expand").SetValue("9")
	cfg.Section("Xbox 360 Controller").Key("Collapse").SetValue("8")

	// HORI CO.,LTD HORIPAD S
	cfg.Section("HORI CO.,LTD HORIPAD S").Key("A").SetValue("2")
	cfg.Section("HORI CO.,LTD HORIPAD S").Key("B").SetValue("1")
	cfg.Section("HORI CO.,LTD HORIPAD S").Key("Start").SetValue("3")
	cfg.Section("HORI CO.,LTD HORIPAD S").Key("Select").SetValue("0")
	cfg.Section("HORI CO.,LTD HORIPAD S").Key("Horizontal").SetValue("0")
	cfg.Section("HORI CO.,LTD HORIPAD S").Key("Vertical").SetValue("1")
	cfg.Section("HORI CO.,LTD HORIPAD S").Key("VerticalNegative").SetValue("1")
	cfg.Section("HORI CO.,LTD HORIPAD S").Key("Expand").SetValue("9")
	cfg.Section("HORI CO.,LTD HORIPAD S").Key("Collapse").SetValue("8")

	// save config
	cfg.SaveTo("gb.ini")

	return cfg
}

// 設定ファイルが存在するか確認 あるならtrue
func checkConfigFileExist() bool {
	filename := "gb.ini"
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}
