package log

import (
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/str"
)

// Event log event
type Event struct {
	Logger Logger `json:"-"`
	Level  Level
	Msg    string
	Time   time.Time
	File   string
	Line   int
	Func   string
	Trace  string
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
		Logger: logger,
		Level:  lvl,
		Msg:    msg,
		Time:   time.Now(),
		File:   "???",
		Func:   "???",
		Line:   0,
	}
	if logger.GetCallerDepth() > 0 {
		le.CallerDepth(logger.GetCallerDepth(), logger.GetTraceLevel() >= lvl)
	}
	return le
}
