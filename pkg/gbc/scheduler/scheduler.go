package scheduler

import "math"

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

func (s *Scheduler) ScheduleEvent(name string, callback func(), after uint64) {
	when := s.cycles + after
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

func (s *Scheduler) ScheduleEventAbsolute(name string, callback func(), when uint64) {
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
