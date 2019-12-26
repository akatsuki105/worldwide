package rtc

import (
	"fmt"
	"time"
)

// RTC Real Time Clock
type RTC struct {
	Working    bool
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
	rtc.Working = true
	for range time.Tick(time.Second) {
		if rtc.Working {
			if rtc.isActive() {
				rtc.incrementSecond()
			}
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
	// 1分経過
	if rtc.S == 60 {
		rtc.incrementMinute()
	}
}

func (rtc *RTC) incrementMinute() {
	rtc.M++
	rtc.S = 0
	// 1時間経過
	if rtc.M == 60 {
		rtc.incrementHour()
	}
}

func (rtc *RTC) incrementHour() {
	rtc.H++
	rtc.M = 0
	// 1日経過
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

// Dump RTC on common format
func (rtc *RTC) Dump() []byte {
	result := make([]byte, 48)

	result[0] = rtc.S
	result[4] = rtc.M
	result[8] = rtc.H
	result[12] = rtc.DL
	result[16] = rtc.DH
	result[20] = rtc.LatchedRTC.S
	result[24] = rtc.LatchedRTC.M
	result[28] = rtc.LatchedRTC.H
	result[32] = rtc.LatchedRTC.DL
	result[36] = rtc.LatchedRTC.DH

	now := time.Now().Unix()
	now0_8 := byte(now)
	now8_16 := byte(now >> 8)
	now16_24 := byte(now >> 16)
	now24_32 := byte(now >> 24)
	result[40] = now0_8
	result[41] = now8_16
	result[42] = now16_24
	result[43] = now24_32
	return result
}

// Sync RTC data
func (rtc *RTC) Sync(value []byte) {
	if len(value) != 44 && len(value) != 48 {
		fmt.Println("invalid RTC format")
		return
	}

	rtc.S = value[0]
	rtc.M = value[4]
	rtc.H = value[8]
	rtc.DL = value[12]
	rtc.DH = value[16]
	rtc.LatchedRTC.S = value[20]
	rtc.LatchedRTC.M = value[24]
	rtc.LatchedRTC.H = value[28]
	rtc.LatchedRTC.DL = value[32]
	rtc.LatchedRTC.DH = value[36]

	lastSaveTime := (int64(value[43]) << 24) | (int64(value[42]) << 16) | (int64(value[41]) << 8) | int64(value[40])
	// fmt.Println(time.Now(), time.Unix(lastSaveTime, 0))
	delta := int(time.Now().Unix() - lastSaveTime)

	for i := 0; i < delta; i++ {
		rtc.Working = false
		rtc.incrementSecond()
	}
	rtc.Working = true
}
