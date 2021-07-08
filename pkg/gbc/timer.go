package gbc

import (
	"gbc/pkg/gbc/scheduler"
	"gbc/pkg/util"
	"math"
)

const (
	GB_DMG_DIV_PERIOD = 16 // 16cycle = 1/16384 sec, 8cycle = 1/32768 sec
)

type Timer struct {
	p           *GBC
	internalDiv uint32
	nextDiv     uint32
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
	t.p.sound.Buffer(int(cycles))
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
// 1/16384 sec or 1/32768 sec
func (t *Timer) divIncrement() {
	tMultiplier := util.Bool2U32(t.p.doubleSpeed)
	for t.nextDiv >= GB_DMG_DIV_PERIOD>>tMultiplier {
		t.nextDiv -= GB_DMG_DIV_PERIOD >> tMultiplier

		if t.timaPeriod > 0 && (t.internalDiv&(t.timaPeriod-1)) == (t.timaPeriod-1) {
			t.p.IO[TIMAIO]++
			if t.p.IO[TIMAIO] == 0 {
				// overflow(4 cycles delay https://github.com/Gekkio/mooneye-gb/blob/master/tests/acceptance/timer/tima_reload.s)
				t.p.scheduler.ScheduleEvent(scheduler.TimerIRQ, t.irq, 4)
			}
		}

		t.internalDiv++
		t.p.IO[DIVIO] = byte(t.internalDiv >> 4)
	}
}

// _GBTimerUpdate (system count)
// 1/16384 sec or 1/32768 sec
func (t *Timer) update(cyclesLate uint64) {
	t.divIncrement()
	t.nextDiv += uint32(cyclesLate)

	// Batch div increments
	divsToGo := 16 - (t.internalDiv & 15)
	timaToGo := uint32(math.MaxUint32)
	if t.timaPeriod > 0 {
		timaToGo = t.timaPeriod - (t.internalDiv & (t.timaPeriod - 1))
	}
	if timaToGo < divsToGo {
		divsToGo = timaToGo
	}
	t.nextDiv = (GB_DMG_DIV_PERIOD * divsToGo) >> util.Bool2U32(t.p.doubleSpeed)
	t.p.scheduler.ScheduleEvent(scheduler.TimerUpdate, t.update, uint64(t.nextDiv)-cyclesLate)
}

// GBTimerDivReset
// triggered on writing DIV
func (t *Timer) divReset() {
	t.nextDiv -= uint32(t.p.scheduler.Until(scheduler.TimerUpdate))
	t.p.scheduler.DescheduleEvent(scheduler.TimerUpdate)
	t.divIncrement()

	t.p.IO[DIVIO] = 0
	t.internalDiv = 0
	t.nextDiv = GB_DMG_DIV_PERIOD >> util.Bool2U32(t.p.doubleSpeed) // 16 or 8 -> 1/16384 sec or 1/32768 sec
	t.p.scheduler.ScheduleEvent(scheduler.TimerUpdate, t.update, uint64(t.nextDiv))
}

// triggerd on writing TAC
func (t *Timer) updateTAC(tac byte) byte {
	if util.Bit(tac, 2) {
		t.nextDiv -= uint32(t.p.scheduler.Until(scheduler.TimerUpdate))
		t.p.scheduler.DescheduleEvent(scheduler.TimerUpdate)
		t.divIncrement()

		timaLt := [4]uint32{1024 >> 4, 16 >> 4, 64 >> 4, 256 >> 4}
		t.timaPeriod = timaLt[tac&0x3]

		t.nextDiv += GB_DMG_DIV_PERIOD >> util.Bool2U32(t.p.doubleSpeed)
		t.p.scheduler.ScheduleEvent(scheduler.TimerUpdate, t.update, uint64(t.nextDiv))
	} else {
		t.timaPeriod = 0
	}
	return tac
}
