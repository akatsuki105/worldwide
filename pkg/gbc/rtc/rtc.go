package rtc

import (
	"fmt"
	"gbc/pkg/util"
	"time"
)

const (
	S = iota
	M
	H
	DL
	DH
)

// RTC Real Time Clock
type RTC struct {
	Enable     bool
	Mapped     uint
	Ctr        [5]byte
	Latched    bool
	LatchedRTC LatchedRTC
}

// LatchedRTC Latched RTC
type LatchedRTC struct{ Ctr [5]byte }

// Init rtc clock
func (rtc *RTC) Init() {
	rtc.Enable = true
	for range time.Tick(time.Second) {
		if rtc.Enable {
			if rtc.isActive() {
				rtc.incrementSecond()
			}
		}
	}
}

// Read fetch clock register
func (rtc *RTC) Read(target byte) byte {
	if rtc.Latched {
		return rtc.LatchedRTC.Ctr[target-0x08]
	}
	return rtc.Ctr[target-0x08]
}

// Latch rtc
func (rtc *RTC) Latch() {
	for i := 0; i < 5; i++ {
		rtc.LatchedRTC.Ctr[i] = rtc.Ctr[i]
	}
}

// Write set clock register
func (rtc *RTC) Write(target, value byte) {
	rtc.Ctr[target-0x08] = value
}

func (rtc *RTC) incrementSecond() {
	rtc.Ctr[S]++
	// pass one minute
	if rtc.Ctr[S] == 60 {
		rtc.incrementMinute()
	}
}

func (rtc *RTC) incrementMinute() {
	rtc.Ctr[M]++
	rtc.Ctr[S] = 0
	// pass 1 hour
	if rtc.Ctr[M] == 60 {
		rtc.incrementHour()
	}
}

func (rtc *RTC) incrementHour() {
	rtc.Ctr[H]++
	rtc.Ctr[M] = 0
	// pass a day
	if rtc.Ctr[H] == 24 {
		rtc.incrementDay()
	}
}

func (rtc *RTC) incrementDay() {
	previousDL := rtc.Ctr[DL]
	rtc.Ctr[DL]++
	rtc.Ctr[H] = 0
	// pass 256 days
	if rtc.Ctr[DL] < previousDL {
		rtc.Ctr[DL] = 0
		if rtc.Ctr[DH]&0x01 == 1 {
			// msb on day is set
			rtc.Ctr[DH] |= 0x80
			rtc.Ctr[DH] &= 0x7f
		} else {
			// msb on day is clear
			rtc.Ctr[DH] |= 0x01
		}
	}
}

func (rtc *RTC) isActive() bool { return !util.Bit(rtc.Ctr[DH], 6) }

// Dump RTC on common format
func (rtc *RTC) Dump() []byte {
	result := make([]byte, 48)

	result[0], result[4], result[8], result[12], result[16] = rtc.Ctr[S], rtc.Ctr[M], rtc.Ctr[H], rtc.Ctr[DL], rtc.Ctr[DH]
	result[20], result[24], result[28], result[32], result[36] = rtc.LatchedRTC.Ctr[S], rtc.LatchedRTC.Ctr[M], rtc.LatchedRTC.Ctr[H], rtc.LatchedRTC.Ctr[DL], rtc.LatchedRTC.Ctr[DH]

	now := time.Now().Unix()
	result[40], result[41], result[42], result[43] = byte(now), byte(now>>8), byte(now>>16), byte(now>>24)
	return result
}

// Sync RTC data
func (rtc *RTC) Sync(value []byte) {
	if len(value) != 44 && len(value) != 48 {
		fmt.Println("invalid RTC format")
		return
	}
	rtc.Ctr = [5]byte{value[0], value[4], value[8], value[12], value[16]}
	rtc.LatchedRTC.Ctr = [5]byte{value[20], value[24], value[28], value[32], value[36]}
	lastSaveTime := (int64(value[43]) << 24) | (int64(value[42]) << 16) | (int64(value[41]) << 8) | int64(value[40])
	delta := int(time.Now().Unix()-lastSaveTime) / 60
	rtc.advance(delta)
}

// Advance rtc clock
func (rtc *RTC) advance(minutes int) {
	for i := 0; i < minutes; i++ {
		rtc.incrementMinute()
	}
}
