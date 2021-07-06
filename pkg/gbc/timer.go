package gbc

import (
	"gbc/pkg/gbc/scheduler"
	"gbc/pkg/util"
	"math"
)

const (
	GB_DMG_DIV_PERIOD = 16
)

type GBTimer struct {
	p           *GBC
	internalDiv uint32
	nextDiv     uint32
	timaPeriod  uint32
}

func NewTimer(p *GBC) *GBTimer {
	t := &GBTimer{
		p: p,
	}
	t.reset()
	return t
}

// GBTimerReset
func (t *GBTimer) reset() {
	t.nextDiv = GB_DMG_DIV_PERIOD * 2
	t.timaPeriod = 1024 >> 4
}

// mTimingTick
func (t *GBTimer) tick(cycles uint32) {
	t.p.scheduler.Add(uint64(cycles))

	for {
		if t.p.scheduler.Next() > t.p.scheduler.Cycle() {
			break
		}
		t.p.scheduler.DoEvent()
	}
}

// _GBTimerIRQ
func (t *GBTimer) irq() {
	t.p.IO[TIMAIO-0xff00] = t.p.IO[TMAIO-0xff00]
	t.p.IO[IFIO-0xff00] = util.SetBit8(t.p.IO[IFIO-0xff00], 2, true)
	t.p.updateIRQs()
}

// _GBTimerDivIncrement
// 1/16384 sec or 1/32768 sec
func (t *GBTimer) divIncrement() {
	tMultiplier := 2 - util.Bool2U32(t.p.doubleSpeed)
	for t.nextDiv >= GB_DMG_DIV_PERIOD*tMultiplier {
		t.nextDiv -= GB_DMG_DIV_PERIOD * tMultiplier

		if t.timaPeriod > 0 && (t.internalDiv&(t.timaPeriod-1)) == (t.timaPeriod-1) {
			t.p.IO[TIMAIO-0xff00]++
			if t.p.IO[TIMAIO-0xff00] == 0 {
				// overflow
				t.p.scheduler.ScheduleEvent(scheduler.TimerIRQ, t.irq, uint64(7*tMultiplier))
			}
		}

		t.internalDiv++
		t.p.IO[DIVIO-0xff00] = byte(t.internalDiv >> 4)
	}
}

// _GBTimerUpdate (system count)
// 1/16384 sec or 1/32768 sec
func (t *GBTimer) update() {
	t.divIncrement()

	// Batch div increments
	divsToGo := 16 - (t.internalDiv & 15)
	timaToGo := uint32(math.MaxUint32)
	if t.timaPeriod > 0 {
		timaToGo = t.timaPeriod - (t.internalDiv & (t.timaPeriod - 1))
	}
	if timaToGo < divsToGo {
		divsToGo = timaToGo
	}
	t.nextDiv = GB_DMG_DIV_PERIOD * divsToGo * (2 - util.Bool2U32(t.p.doubleSpeed))
	t.p.scheduler.ScheduleEvent(scheduler.TimerUpdate, t.update, uint64(t.nextDiv))
}

// GBTimerDivReset
// triggered on writing DIV
func (t *GBTimer) divReset() {
	t.nextDiv -= uint32(t.p.scheduler.Until(scheduler.TimerUpdate))
	t.p.scheduler.DescheduleEvent(scheduler.TimerUpdate)
	t.divIncrement()
	tMultiplier := 2 - util.Bool2U64(t.p.doubleSpeed)
	if ((t.internalDiv << 1) | (t.nextDiv>>((4-util.Bool2U32(t.p.doubleSpeed))&1))&t.timaPeriod) > 0 {
		t.p.IO[TIMAIO-0xff00]++
		if t.p.IO[TIMAIO-0xff00] == 0 {
			t.p.scheduler.ScheduleEvent(scheduler.TimerIRQ, t.irq, 7*tMultiplier)
		}
	}

	t.p.IO[DIVIO-0xff00] = 0
	t.internalDiv = 0
	t.nextDiv = GB_DMG_DIV_PERIOD * (2 - util.Bool2U32(t.p.doubleSpeed)) // 16 or 32 -> 1/16384 sec or 1/32768 sec
	t.p.scheduler.ScheduleEvent(scheduler.TimerUpdate, t.update, uint64(t.nextDiv))
}

// triggerd on writing TAC
func (t *GBTimer) updateTAC(tac byte) byte {
	if util.Bit(tac, 2) {
		t.nextDiv -= uint32(t.p.scheduler.Until(scheduler.TimerUpdate))
		t.p.scheduler.DescheduleEvent(scheduler.TimerUpdate)
		t.divIncrement()

		timaLt := [4]uint32{1024 >> 4, 16 >> 4, 64 >> 4, 256 >> 4}
		t.timaPeriod = timaLt[tac&0x3]

		t.nextDiv += GB_DMG_DIV_PERIOD * (2 - util.Bool2U32(t.p.doubleSpeed))
		t.p.scheduler.ScheduleEvent(scheduler.TimerUpdate, t.update, uint64(t.nextDiv))
	} else {
		t.timaPeriod = 0
	}
	return tac
}
