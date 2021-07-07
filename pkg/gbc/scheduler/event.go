package scheduler

type EventName string

const (
	TimerUpdate EventName = "TimerUpdate"
	TimerIRQ    EventName = "TimerIRQ"
	OAMDMA      EventName = "Oamdma"
	HDMA        EventName = "Hdma"
	EndMode0    EventName = "EndMode0"
	EndMode1    EventName = "EndMode1"
	EndMode2    EventName = "EndMode2"
	EndMode3    EventName = "EndMode3"
	UpdateFrame EventName = "UpdateFrame"
)

type Event struct {
	name     EventName
	callback func()
	when     uint64
	next     *Event
}
