package log

import (
	"time"
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
func (wbw *BatchWriter) SetFlushLevel(lvl string) {
	wbw.FlushLevel = ParseLevel(lvl)
}

func (wbw *BatchWriter) InitBuffer() {
	if wbw.BatchCount < 1 {
		wbw.BatchCount = 10
	}
	if wbw.CacheCount < wbw.BatchCount {
		wbw.CacheCount = wbw.BatchCount * 2
	}

	if wbw.EventBuffer == nil {
		wbw.EventBuffer = NewEventBuffer(wbw.CacheCount)
	}
}

func (wbw *BatchWriter) ShouldFlush(le *Event) bool {
	if wbw.EventBuffer.Len() >= wbw.BatchCount {
		return true
	}
	if le.Level <= wbw.FlushLevel {
		return true
	}
	if wbw.FlushDelta > 0 && wbw.EventBuffer.Len() > 1 {
		if fle, ok := wbw.EventBuffer.Peek(); ok {
			if le.When.Sub(fle.When) >= wbw.FlushDelta {
				return true
			}
		}
	}
	return false
}
