package scheduler

import (
	"math"
)

// Scheduler manages system event
type Scheduler struct {
	cycles uint64
	root   *Event
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Cycle() uint64 {
	return s.cycles
}

func (s *Scheduler) Add(c uint64) {
	s.cycles += c
}

func (s *Scheduler) Next() uint64 {
	if s.root == nil {
		return math.MaxUint64
	}
	return s.root.when
}

func (s *Scheduler) ScheduleEvent(name EventName, callback func(), after uint64) {
	when := s.cycles + after
	var previous *Event = nil
	event := s.root
	for {
		if event == nil {
			s.root = &Event{
				name:     name,
				callback: callback,
				when:     when,
				next:     event,
			}
			return
		}
		if event.next == nil {
			// last executed
			event.next = &Event{
				name:     name,
				callback: callback,
				when:     when,
			}
			return
		}
		if when < event.when {
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
		previous = event
		event = event.next
	}
}

func (s *Scheduler) ScheduleEventAbsolute(name EventName, callback func(), when uint64) {
	var previous *Event = nil
	event := s.root
	for {
		if event.next == nil {
			// last executed
			event.next = &Event{
				name:     name,
				callback: callback,
				when:     when,
			}
			return
		}
		if when < event.when {
			// previous <- new <- event
			newEvent := &Event{
				name:     name,
				callback: callback,
				when:     when,
				next:     event,
			}
			if previous == nil {
				s.root = newEvent
			} else {
				previous.next = newEvent
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
				s.root = nil
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
	event.callback()
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
