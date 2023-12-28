//go:build !go1.18
// +build !go1.18

package log

import "github.com/askasoft/pango/col"

// EventBuffer a event buffer
type EventBuffer struct {
	col.RingBuffer
	BufSize int
}

func NewEventBuffer(size int) *EventBuffer {
	return &EventBuffer{
		BufSize: size,
	}
}

// Len returns the number of elements currently stored in the buffer.
func (eb *EventBuffer) Len() int {
	return eb.RingBuffer.Len()
}

// Push push event to the buffer
func (eb *EventBuffer) Push(le *Event) {
	eb.RingBuffer.Push(le)
	for eb.RingBuffer.Len() > eb.BufSize {
		eb.RingBuffer.Poll()
	}
}

// Poll poll event from the buffer
func (eb *EventBuffer) Poll() (*Event, bool) {
	if le, ok := eb.RingBuffer.Poll(); ok {
		return le.(*Event), true
	}
	return nil, false
}

// Peek peek event from the buffer
func (eb *EventBuffer) Peek() (*Event, bool) {
	if le, ok := eb.RingBuffer.Peek(); ok {
		return le.(*Event), true
	}
	return nil, false
}
