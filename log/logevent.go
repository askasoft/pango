package log

import (
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/cog/ringbuffer"
	"github.com/askasoft/pango/str"
)

// EventBuffer a event buffer
type EventBuffer = ringbuffer.RingBuffer[*Event]

// Event log event
type Event struct {
	Name    string
	Props   map[string]any
	Level   Level
	Time    time.Time
	Message string
	File    string
	Line    int
	Func    string
	Trace   string
}

// CallerDepth get caller filename and line number
func (le *Event) CallerDepth(depth int, trace bool) {
	dep := 1
	if trace {
		dep = 30
	}

	rpc := make([]uintptr, dep)
	n := runtime.Callers(depth, rpc)
	if n > 0 {
		frames := runtime.CallersFrames(rpc)
		frame, next := frames.Next()
		_, le.Func = path.Split(frame.Function)
		_, le.File = path.Split(frame.File)
		le.Line = frame.Line
		if next {
			var sb strings.Builder
			for ; next; frame, next = frames.Next() {
				sb.WriteString(frame.File)
				sb.WriteString(":")
				sb.WriteString(strconv.Itoa(frame.Line))
				sb.WriteString(" ")
				sb.WriteString(frame.Function)
				sb.WriteString("()")
				sb.WriteString(EOL)
			}
			le.Trace = sb.String()
		}
	}
}

// CallerStop get caller filename and line number
func (le *Event) CallerStop(stop string, trace bool) {
	rpc := make([]uintptr, 30)
	n := runtime.Callers(2, rpc)
	if n > 0 {
		found := false
		frames := runtime.CallersFrames(rpc)
		for frame, next := frames.Next(); next; frame, next = frames.Next() {
			if str.Contains(frame.File, stop) {
				found = true
				continue
			}

			if found {
				_, le.Func = path.Split(frame.Function)
				_, le.File = path.Split(frame.File)
				le.Line = frame.Line
				if trace && next {
					var sb strings.Builder
					for ; next; frame, next = frames.Next() {
						sb.WriteString(frame.File)
						sb.WriteString(":")
						sb.WriteString(strconv.Itoa(frame.Line))
						sb.WriteString(" ")
						sb.WriteString(frame.Function)
						sb.WriteString("()")
						sb.WriteString(EOL)
					}
					le.Trace = sb.String()
				}
				break
			}
		}
	}
}

func NewEvent(logger Logger, lvl Level, msg string) *Event {
	le := &Event{
		Name:    logger.GetName(),
		Props:   logger.GetProps(),
		Level:   lvl,
		Time:    time.Now(),
		File:    "???",
		Func:    "???",
		Message: msg,
	}
	if logger.GetCallerDepth() > 0 {
		le.CallerDepth(logger.GetCallerDepth(), logger.GetTraceLevel() >= lvl)
	}
	return le
}
