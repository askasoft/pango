package log

import (
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Event log event
type Event struct {
	Logger Logger
	Level  int
	Msg    string
	When   time.Time
	File   string
	Line   int
	Func   string
	Trace  string
}

// EventPool log event pool
var EventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{}
	},
}

// Caller get caller filename and line number
func (le *Event) Caller(depth int, trace bool) {
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
			sb := strings.Builder{}
			for ; next; frame, next = frames.Next() {
				sb.WriteString(frame.File)
				sb.WriteString(":")
				sb.WriteString(strconv.Itoa(frame.Line))
				sb.WriteString(" ")
				sb.WriteString(frame.Function)
				sb.WriteString("()")
				sb.WriteString(eol)
			}
			le.Trace = sb.String()
		}
	} else {
		le.File = "???"
		le.Line = 0
		le.Func = "???"
	}
}

// NewEvent create a new log event
func NewEvent(logger Logger, lvl int, msg string) *Event {
	le := EventPool.Get().(*Event)
	le.Logger = logger
	le.Level = lvl
	le.Msg = msg
	le.When = time.Now()
	le.File = ""
	le.Line = 0
	le.Trace = ""
	if logger.GetCallerDepth() > 0 {
		le.Caller(logger.GetCallerDepth(), logger.GetTraceLevel() >= lvl)
	}
	return le
}
