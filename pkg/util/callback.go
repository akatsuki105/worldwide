package util

import (
	"errors"
	"sort"
)

type Priority uint

const (
	PRIO_HISTORY    = 0
	PRIO_BREAKPOINT = 1
)

type Callback struct {
	Name     string
	Priority Priority
	Func     func() bool
}

func SetCallback(callbacks []*Callback, name string, priority Priority, callback func() bool) ([]*Callback, error) {
	for _, c := range callbacks {
		if name == c.Name {
			return callbacks, errors.New("callback name is already used")
		}
		if priority == c.Priority {
			return callbacks, errors.New("callback priority is already used")
		}
	}
	callbacks = append(callbacks, &Callback{name, priority, callback})
	sortCallbacks(callbacks)
	return callbacks, nil
}

func RemoveCallback(callbacks []*Callback, name string) []*Callback {
	for i, c := range callbacks {
		if name == c.Name {
			sortCallbacks(callbacks)
			callbacks = append(callbacks[:i], callbacks[i+1:]...)
			return callbacks
		}
	}
	return callbacks
}

func sortCallbacks(callbacks []*Callback) {
	sort.Slice(callbacks, func(i, j int) bool {
		return callbacks[i].Priority < callbacks[j].Priority
	})
}
