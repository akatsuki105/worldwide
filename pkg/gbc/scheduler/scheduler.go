package scheduler

import (
	"fmt"
	"math"
)

// Scheduler manages system event
type Scheduler struct {
	cycles uint64
	root   *Event
}

func New() *Scheduler              { return &Scheduler{} }
func (s *Scheduler) Cycle() uint64 { return s.cycles }
func (s *Scheduler) Add(c uint64)  { s.cycles += c }

func (s *Scheduler) Next() uint64 {
	if s.root == nil {
		return math.MaxUint64
	}
	return s.root.when
}

func (s *Scheduler) ScheduleEvent(name EventName, callback func(cyclesLate uint64), after uint64) {
	when := s.cycles + after
	var previous *Event = nil
	event := s.root
	for {
		if event == nil {
			s.root = &Event{
				name:     name,
				callback: callback,
				when:     when,
			}
			return
		}

		if when < event.when {
			if previous == nil {
				// new <- event
				newEvent := &Event{
					name:     name,
					callback: callback,
					when:     when,
					next:     event,
				}
				s.root = newEvent
				return
			}
			if previous.when <= when {
				// previous <- new <- event
				newEvent := &Event{
					name:     name,
					callback: callback,
					when:     when,
					next:     event,
				}
				previous.next = newEvent
				return
			}
		}

		if event.next == nil && event.when <= when {
			// last executed
			event.next = &Event{
				name:     name,
				callback: callback,
				when:     when,
			}
			return
		}
		previous = event
		event = event.next
	}
}

func (s *Scheduler) DescheduleEvent(name EventName) {
	var previous *Event = nil
	event := s.root
	for {
		if event == nil {
			return
		}

		if event.name == name {
			if previous == nil {
				s.root = event.next
				return
			} else {
				previous.next = event.next
				return
			}
		}

		previous = event
		event = event.next
	}
}

func (s *Scheduler) DoEvent() {
	event := s.root
	if event == nil {
		return
	}
	s.root = event.next
	event.callback(s.cycles - event.when)
}

func (s *Scheduler) Until(name EventName) uint64 {
	event := s.root
	for {
		if event == nil {
			return math.MaxUint64
		}

		if event.name == name {
			return event.when - s.cycles
		}
		event = event.next
	}
}

func (s *Scheduler) String() string {
	result := ""
	event := s.root
	for event != nil {
		result += fmt.Sprintf("%s:%d->", event.name, event.when)
		event = event.next
	}
	return result
}
