package rtc

import (
	"fmt"
	"time"

	"github.com/pokemium/worldwide/pkg/util"
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

func New(enable bool) *RTC { return &RTC{Enable: enable} }

func (rtc *RTC) IncrementSecond() {
	if rtc.Enable && rtc.isActive() {
		rtc.incrementSecond()
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
	if rtc.Ctr[S] == 60 {
		rtc.incrementMinute()
	}
}

func (rtc *RTC) incrementMinute() {
	rtc.Ctr[M]++
	rtc.Ctr[S] = 0
	if rtc.Ctr[M] == 60 {
		rtc.incrementHour()
	}
}

func (rtc *RTC) incrementHour() {
	rtc.Ctr[H]++
	rtc.Ctr[M] = 0
	if rtc.Ctr[H] == 24 {
		rtc.incrementDay()
	}
}

func (rtc *RTC) incrementDay() {
	old := rtc.Ctr[DL]
	rtc.Ctr[DL]++
	rtc.Ctr[H] = 0
	// pass 256 days
	if rtc.Ctr[DL] < old {
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

// Dump RTC on .sav format
//
// offset  size    desc
// 0       4       time seconds
// 4       4       time minutes
// 8       4       time hours
// 12      4       time days
// 16      4       time days high
// 20      4       latched time seconds
// 24      4       latched time minutes
// 28      4       latched time hours
// 32      4       latched time days
// 36      4       latched time days high
// 40      4       unix timestamp when saving
// 44      4       0   (probably the high dword of 64 bits time), absent in the 44 bytes version
func (rtc *RTC) Dump() []byte {
	result := make([]byte, 48)

	result[0], result[4], result[8], result[12], result[16] = rtc.Ctr[S], rtc.Ctr[M], rtc.Ctr[H], rtc.Ctr[DL], rtc.Ctr[DH]

	latch := rtc.LatchedRTC
	result[20], result[24], result[28], result[32], result[36] = latch.Ctr[S], latch.Ctr[M], latch.Ctr[H], latch.Ctr[DL], latch.Ctr[DH]

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

	savTime := (uint32(value[43]) << 24) | (uint32(value[42]) << 16) | (uint32(value[41]) << 8) | uint32(value[40])
	delta := uint32(time.Now().Unix()) - savTime
	for i := uint32(0); i < delta; i++ {
		rtc.incrementSecond()
	}
}
