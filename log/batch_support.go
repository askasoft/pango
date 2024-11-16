package log

import (
	"time"

	"github.com/askasoft/pango/log/internal"
)

// BatchSupport support event cache and flush events on FlushLevel or BatchCount reached.
type BatchSupport struct {
	BatchCount  int           // flush events if events count >= BatchCount
	CacheCount  int           // the maximun cacheable event count
	FlushLevel  Level         // flush events if event <= FlushLevel
	FlushDelta  time.Duration // flush events if [current log event time] - [first log event time] >= FlushDelta
	BatchBuffer EventBuffer
}

// SetFlushLevel set the flush level
func (bs *BatchSupport) SetFlushLevel(lvl string) {
	bs.FlushLevel = ParseLevel(lvl)
}

func (bs *BatchSupport) shouldFlush(le *Event) bool {
	if bs.BatchBuffer.Len() >= bs.BatchCount {
		return true
	}
	if le.Level <= bs.FlushLevel {
		return true
	}

	if bs.FlushDelta > 0 && bs.BatchBuffer.Len() > 1 {
		if fle, ok := bs.BatchBuffer.Peek(); ok {
			if le.Time.Sub(fle.Time) >= bs.FlushDelta {
				return true
			}
		}
	}
	return false
}

func (bs *BatchSupport) BatchWrite(le *Event, flush func(*EventBuffer) error) {
	bs.BatchBuffer.Push(le)

	if bs.shouldFlush(le) {
		if err := flush(&bs.BatchBuffer); err != nil {
			internal.Perror(err)
			if bs.BatchBuffer.Len() > bs.CacheCount {
				bs.BatchBuffer.Poll()
			}
		} else {
			bs.BatchBuffer.Clear()
		}
	}
}

func (bs *BatchSupport) BatchFlush(flush func(*EventBuffer) error) {
	if bs.BatchBuffer.IsEmpty() {
		return
	}

	if err := flush(&bs.BatchBuffer); err != nil {
		internal.Perror(err)
	}
}
