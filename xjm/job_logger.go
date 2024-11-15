package xjm

import (
	"time"

	"github.com/askasoft/pango/log"
)

// JobLogWriter implements log Writer Interface and writes messages to terminal.
type JobLogWriter struct {
	log.LogFilter
	log.BatchWriter

	jmr JobManager
	jid int64
	jls []*JobLog // buffer
}

func NewJobLogWriter(jmr JobManager, jid int64) *JobLogWriter {
	jw := &JobLogWriter{jmr: jmr, jid: jid}

	jw.Filter = log.NewLevelFilter(log.LevelDebug)
	jw.BatchCount = 100
	jw.CacheCount = 200
	jw.FlushLevel = log.LevelWarn
	jw.FlushDelta = time.Second

	return jw
}

// Write write log event.
func (jw *JobLogWriter) Write(le *log.Event) (err error) {
	if jw.Reject(le) {
		return
	}

	return jw.BatchWrite(le, jw.flush)
}

// Flush flush cached log events
func (jw *JobLogWriter) Flush() {
	jw.BatchFlush(jw.flush)
}

// Close flush cached log events
func (jw *JobLogWriter) Close() {
	jw.Flush()
}

func (jw *JobLogWriter) flush() error {
	for len(jw.jls) < jw.EventBuffer.Len() {
		jw.jls = append(jw.jls, &JobLog{})
	}

	jls := jw.jls[:jw.EventBuffer.Len()]

	for n, it := 0, jw.EventBuffer.Iterator(); it.Next(); {
		le := it.Value()

		jl := jls[n]
		jl.ID = 0
		jl.JID = jw.jid
		jl.Time = le.Time
		jl.Level = le.Level.Prefix()
		jl.Message = le.Msg
		n++
	}

	if err := jw.jmr.AddJobLogs(jls); err != nil {
		return err
	}

	jw.EventBuffer.Clear()
	return nil
}
