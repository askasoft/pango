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

// eventPool log event pool
var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{}
	},
}

// Clear clear event values
func (le *Event) Clear() {
	le.Logger = nil
	le.Level = LevelNone
	le.Msg = ""
	le.When = 0
	le.File = ""
	le.Line = 0
	le.Func = ""
	le.Trace = ""
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

// newEvent get a log event from pool
func newEvent(logger Logger, lvl int, msg string) *Event {
	le := eventPool.Get().(*Event)
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

// putEvent put event back to pool
func putEvent(le *Event) {
	le.Clear()
	eventPool.Put(le)
}
