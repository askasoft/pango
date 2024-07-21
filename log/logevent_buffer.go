//go:build go1.18
// +build go1.18

package log

import "github.com/askasoft/pango/cog/ringbuffer"

// EventBuffer a event buffer
type EventBuffer struct {
	ringbuffer.RingBuffer[*Event]
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
