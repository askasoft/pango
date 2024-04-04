package xjm

import (
	"fmt"
	"os"
	"time"

	"github.com/askasoft/pango/log"
)

// JobLogWriter implements log Writer Interface and writes messages to terminal.
type JobLogWriter struct {
	log.LogFilter
	log.BatchWriter

	jmr JobManager
	jid int64
}

func NewJobLogWriter(jmr JobManager, jid int64) *JobLogWriter {
	jlw := &JobLogWriter{jmr: jmr, jid: jid}

	jlw.Filter = log.NewLevelFilter(log.LevelDebug)
	jlw.BatchCount = 100
	jlw.CacheCount = 200
	jlw.FlushLevel = log.LevelWarn
	jlw.FlushDelta = time.Second

	return jlw
}

// Write write message in console.
func (jlw *JobLogWriter) Write(le *log.Event) (err error) {
	if jlw.Reject(le) {
		return
	}

	jlw.InitBuffer()
	jlw.EventBuffer.Push(le)

	if jlw.ShouldFlush(le) {
		jlw.Flush()
	}

	return nil
}

// Flush implementing method. empty.
func (jlw *JobLogWriter) Flush() {
	if jlw.EventBuffer == nil || jlw.EventBuffer.IsEmpty() {
		return
	}

	if err := jlw.flush(); err == nil {
		jlw.EventBuffer.Clear()
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func (jlw *JobLogWriter) flush() error {
	var jls []*JobLog

	for it := jlw.EventBuffer.Iterator(); it.Next(); {
		le := it.Value()
		jl := &JobLog{
			JID:     jlw.jid,
			Time:    le.Time,
			Level:   le.Level.Prefix(),
			Message: le.Msg,
		}
		jls = append(jls, jl)
	}

	return jlw.jmr.AddJobLogs(jls)
}

// Close implementing method. empty.
func (jlw *JobLogWriter) Close() {
	jlw.Flush()
}
