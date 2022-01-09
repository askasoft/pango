package log

import "github.com/pandafw/pango/col"

// EventBuffer a event buffer
type EventBuffer struct {
	buffer  col.RingBuffer
	BufSize int
}

// Len returns the number of elements currently stored in the buffer.
func (eb *EventBuffer) Len() int {
	return eb.buffer.Len()
}

// Push push event to the buffer
func (eb *EventBuffer) Push(le *Event) {
	eb.buffer.Push(le)
	for eb.buffer.Len() > eb.BufSize {
		eb.buffer.Poll()
	}
}

// Poll poll event from the buffer
func (eb *EventBuffer) Poll() *Event {
	if le, ok := eb.buffer.Poll(); ok {
		return le.(*Event)
	}
	return nil
}

// Peek peek event from the buffer
func (eb *EventBuffer) Peek() *Event {
	if le, ok := eb.buffer.Peek(); ok {
		return le.(*Event)
	}
	return nil
}
