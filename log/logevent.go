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
	logger Logger
	level  Level
	msg    string
	when   time.Time
	file   string
	line   int
	_func  string
	trace  string
}

// Logger returns log event's Logger
func (le *Event) Logger() Logger {
	return le.logger
}

// Level returns log event's level
func (le *Event) Level() Level {
	return le.level
}

// Msg returns log event's message
func (le *Event) Msg() string {
	return le.msg
}

// When returns log event's time
func (le *Event) When() time.Time {
	return le.when
}

// File returns log event's file
func (le *Event) File() string {
	return le.file
}

// Line returns log event's line
func (le *Event) Line() int {
	return le.line
}

// Func returns log event's function
func (le *Event) Func() string {
	return le._func
}

// Trace returns log event's stack trace
func (le *Event) Trace() string {
	return le.trace
}

// eventPool log event pool
var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{}
	},
}

// clear clear event values
func (le *Event) clear() {
	le.logger = nil
	le.level = LevelNone
	le.msg = ""
	le.file = ""
	le.line = 0
	le._func = ""
	le.trace = ""
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
		_, le._func = path.Split(frame.Function)
		_, le.file = path.Split(frame.File)
		le.line = frame.Line
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
			le.trace = sb.String()
		}
	} else {
		le.file = "???"
		le.line = 0
		le._func = "???"
	}
}

// newEvent get a log event from pool
func newEvent(logger Logger, lvl Level, msg string) *Event {
	le := eventPool.Get().(*Event)
	le.logger = logger
	le.level = lvl
	le.msg = msg
	le.when = time.Now()
	if logger.GetCallerDepth() > 0 {
		le.Caller(logger.GetCallerDepth(), logger.GetTraceLevel() >= lvl)
	}
	return le
}

// putEvent put event back to pool
func putEvent(le *Event) {
	le.clear()
	eventPool.Put(le)
}
