package rtc

import "time"

// RTC Real Time Clock
type RTC struct {
	Mapped     uint
	S          byte
	M          byte
	H          byte
	DL         byte
	DH         byte
	Latched    bool
	LatchedRTC LatchedRTC
}

// LatchedRTC Latched RTC
type LatchedRTC struct {
	S  byte
	M  byte
	H  byte
	DL byte
	DH byte
}

// Init クロックを開始する
func (rtc *RTC) Init() {
	for range time.Tick(time.Second) {
		if rtc.isActive() {
			rtc.incrementSecond()
		}
	}
}

// Read fetch clock register
func (rtc *RTC) Read(target byte) byte {
	if rtc.Latched {
		switch target {
		case 0x08:
			return rtc.LatchedRTC.S
		case 0x09:
			return rtc.LatchedRTC.M
		case 0x0a:
			return rtc.LatchedRTC.H
		case 0x0b:
			return rtc.LatchedRTC.DL
		case 0x0c:
			return rtc.LatchedRTC.DH
		}
	} else {
		switch target {
		case 0x08:
			return rtc.S
		case 0x09:
			return rtc.M
		case 0x0a:
			return rtc.H
		case 0x0b:
			return rtc.DL
		case 0x0c:
			return rtc.DH
		}
	}
	return 0
}

// Latch ラッチする
func (rtc *RTC) Latch() {
	rtc.LatchedRTC.S = rtc.S
	rtc.LatchedRTC.M = rtc.M
	rtc.LatchedRTC.H = rtc.H
	rtc.LatchedRTC.DL = rtc.DL
	rtc.LatchedRTC.DH = rtc.DH
}

// Write set clock register
func (rtc *RTC) Write(target, value byte) {
	switch target {
	case 0x08:
		rtc.S = value
	case 0x09:
		rtc.M = value
	case 0x0a:
		rtc.H = value
	case 0x0b:
		rtc.DL = value
	case 0x0c:
		rtc.DH = value
	}
}

func (rtc *RTC) incrementSecond() {
	rtc.S++
	// 1分
	if rtc.S == 60 {
		rtc.incrementMinute()
	}
}

func (rtc *RTC) incrementMinute() {
	rtc.M++
	rtc.S = 0
	// 1時間
	if rtc.M == 60 {
		rtc.incrementHour()
	}
}

func (rtc *RTC) incrementHour() {
	rtc.H++
	rtc.M = 0
	// 1日
	if rtc.H == 24 {
		rtc.incrementDay()
	}
}

func (rtc *RTC) incrementDay() {
	previousDL := rtc.DL
	rtc.DL++
	rtc.H = 0
	// 256日経過
	if rtc.DL < previousDL {
		rtc.DL = 0
		if rtc.DH&0x01 == 1 {
			// 日付最上位bitが1
			rtc.DH |= 0x80
			rtc.DH &= 0x7f
		} else {
			// 日付最上位bitが0
			rtc.DH |= 0x01
		}
	}
}

func (rtc *RTC) isActive() bool {
	return (rtc.DH&0b01000000 == 0)
}

// エミュレータ専用関数 現在のPC時刻とゲームボーイの時刻(とりあえず時刻だけ)を同期する
func (rtc *RTC) syncRealTime() {
	t := time.Now()
	rtc.H = byte(t.Hour())
	rtc.M = byte(t.Minute())
	rtc.S = byte(t.Second())
}
