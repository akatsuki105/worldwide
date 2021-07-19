package gbc

import (
	"math"

	"github.com/pokemium/worldwide/pkg/gbc/scheduler"
	"github.com/pokemium/worldwide/pkg/util"
)

const (
	GB_DMG_DIV_PERIOD = 16 // this is INTERNAL div interval, 16cycle or 8cycles
)

type Timer struct {
	p           *GBC
	internalDiv uint32 // INTERNAL div counter, real div is `internalDiv >> 4`
	nextDiv     uint32 // next INTERNAL div
	timaPeriod  uint32
}

func NewTimer(p *GBC) *Timer {
	t := &Timer{
		p: p,
	}
	t.reset()
	return t
}

// GBTimerReset
func (t *Timer) reset() {
	t.nextDiv = GB_DMG_DIV_PERIOD
	t.timaPeriod = 1024 >> 4
}

// mTimingTick
func (t *Timer) tick(cycles uint32) {
	t.p.Sound.Buffer(int(cycles))
	t.p.scheduler.Add(uint64(cycles))
	for {
		if t.p.scheduler.Next() > t.p.scheduler.Cycle() {
			break
		}
		t.p.scheduler.DoEvent()
	}
}

// _GBTimerIRQ
func (t *Timer) irq(cyclesLate uint64) {
	t.p.IO[TIMAIO] = t.p.IO[TMAIO]
	t.p.IO[IFIO] = util.SetBit8(t.p.IO[IFIO], 2, true)
	t.p.updateIRQs()
}

// _GBTimerDivIncrement
// oneloop equals to every n cycles (n=16cycles or 8cycles, p.s. 16384Hz=256cycles or 128cycles)
func (t *Timer) internalDivIncrement() {
	tMultiplier := util.Bool2U32(t.p.DoubleSpeed)
	interval := uint32(GB_DMG_DIV_PERIOD >> tMultiplier) // 16 or 8

	// normally, t.nextDiv is greater than 256 or 128, so real div increment should occur.
	for t.nextDiv >= interval {
		t.nextDiv -= interval

		if t.timaPeriod > 0 && (t.internalDiv&(t.timaPeriod-1)) == (t.timaPeriod-1) {
			t.p.IO[TIMAIO]++
			if t.p.IO[TIMAIO] == 0 {
				// overflow(4 cycles delay https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/timer/tima_reload.s)
				t.p.scheduler.ScheduleEvent(scheduler.TimerIRQ, t.irq, 4<<util.Bool2U32(t.p.DoubleSpeed))
			}
		}

		t.internalDiv++
		t.p.IO[DIVIO] = byte(t.internalDiv >> 4)
	}
}

// _GBTimerUpdate (system count)
// 1/16384sec(256cycles) or 1/32768sec(128cycles)
func (t *Timer) update(cyclesLate uint64) {
	t.nextDiv += uint32(cyclesLate)
	t.internalDivIncrement()

	// Batch div increments into real div increment
	divsToGo := 16 - (t.internalDiv & 15)
	timaToGo := uint32(math.MaxUint32)
	if t.timaPeriod > 0 {
		timaToGo = t.timaPeriod - (t.internalDiv & (t.timaPeriod - 1))
	}
	if timaToGo < divsToGo {
		divsToGo = timaToGo
	}
	t.nextDiv = (GB_DMG_DIV_PERIOD * divsToGo) >> util.Bool2U32(t.p.DoubleSpeed) // 256 or 128

	t.p.scheduler.ScheduleEvent(scheduler.TimerUpdate, t.update, uint64(t.nextDiv)-cyclesLate)
}

// GBTimerDivReset
// triggered on writing DIV
func (t *Timer) divReset() {
	t.nextDiv -= uint32(t.p.scheduler.Until(scheduler.TimerUpdate))
	t.p.scheduler.DescheduleEvent(scheduler.TimerUpdate)
	t.internalDivIncrement()

	t.p.IO[DIVIO] = 0
	t.internalDiv = 0
	t.nextDiv = GB_DMG_DIV_PERIOD >> util.Bool2U32(t.p.DoubleSpeed) // 16 or 8 -> 1/16384 sec or 1/32768 sec
	t.p.scheduler.ScheduleEvent(scheduler.TimerUpdate, t.update, uint64(t.nextDiv))
}

// triggerd on writing TAC
func (t *Timer) updateTAC(tac byte) byte {
	if util.Bit(tac, 2) {
		t.nextDiv -= uint32(t.p.scheduler.Until(scheduler.TimerUpdate))
		t.p.scheduler.DescheduleEvent(scheduler.TimerUpdate)
		t.internalDivIncrement()

		timaLt := [4]uint32{1024 >> 4, 16 >> 4, 64 >> 4, 256 >> 4}
		t.timaPeriod = timaLt[tac&0x3]

		t.nextDiv += GB_DMG_DIV_PERIOD >> util.Bool2U32(t.p.DoubleSpeed)
		t.p.scheduler.ScheduleEvent(scheduler.TimerUpdate, t.update, uint64(t.nextDiv))
	} else {
		t.timaPeriod = 0
	}
	return tac
}
