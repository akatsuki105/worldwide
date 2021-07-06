package scheduler

type EventName string

const (
	TimerUpdate EventName = "timerupdate"
	TimerIRQ    EventName = "timerirq"
	OAMDMA      EventName = "oamdma"
	EndMode0    EventName = "endMode0"
	EndMode1    EventName = "endMode1"
	EndMode2    EventName = "endMode2"
	EndMode3    EventName = "endMode3"
)

type Event struct {
	name     EventName
	callback func()
	when     uint64
	next     *Event
}
