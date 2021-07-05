package scheduler

type Event struct {
	name     string
	callback func()
	when     uint64
	next     *Event
}
