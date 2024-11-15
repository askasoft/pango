package log

import (
	"time"

	"github.com/askasoft/pango/log/internal"
)

// BatchWriter implements log Writer Interface and batch send log messages to webhook.
type BatchWriter struct {
	BatchCount int           // event flush batch count
	CacheCount int           // max cacheable event count
	FlushLevel Level         // flush events if event <= FlushLevel
	FlushDelta time.Duration // flush events if [current log event time] - [first log event time] >= FlushDelta

	EventBuffer *EventBuffer
}

// SetFlushLevel set the flush level
func (bw *BatchWriter) SetFlushLevel(lvl string) {
	bw.FlushLevel = ParseLevel(lvl)
}

func (bw *BatchWriter) InitBuffer() {
	if bw.BatchCount < 1 {
		bw.BatchCount = 10
	}
	if bw.CacheCount < bw.BatchCount {
		bw.CacheCount = bw.BatchCount * 2
	}

	if bw.EventBuffer == nil {
		bw.EventBuffer = NewEventBuffer(bw.CacheCount)
	}
}

func (bw *BatchWriter) ShouldFlush(le *Event) bool {
	if bw.EventBuffer == nil {
		return false
	}

	if bw.EventBuffer.Len() >= bw.BatchCount {
		return true
	}
	if le.Level <= bw.FlushLevel {
		return true
	}

	if bw.FlushDelta > 0 && bw.EventBuffer.Len() > 1 {
		if fle, ok := bw.EventBuffer.Peek(); ok {
			if le.Time.Sub(fle.Time) >= bw.FlushDelta {
				return true
			}
		}
	}
	return false
}

func (bw *BatchWriter) BatchWrite(le *Event, flush func() error) error {
	bw.InitBuffer()
	bw.EventBuffer.Push(le)

	if bw.ShouldFlush(le) {
		if err := flush(); err != nil {
			return err
		}
	}
	return nil
}

func (bw *BatchWriter) BatchFlush(flush func() error) {
	if bw.EventBuffer == nil || bw.EventBuffer.IsEmpty() {
		return
	}

	if err := flush(); err != nil {
		internal.Perror(err)
	}
}
