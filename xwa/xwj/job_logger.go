package xwj

import (
	"fmt"
	"os"
	"time"

	"github.com/askasoft/pango/log"
	"gorm.io/gorm"
)

// JobLogWriter implements log Writer Interface and writes messages to terminal.
type JobLogWriter struct {
	log.LogFilter
	log.LogFormatter
	log.BatchWriter

	db  *gorm.DB
	jid int64
}

func NewJobLogWriter(db *gorm.DB, jid int64) *JobLogWriter {
	jlw := &JobLogWriter{db: db, jid: jid}

	jlw.Filter = log.NewLevelFilter(log.LevelDebug)
	jlw.Formatter = log.NewTextFormatter("%t{2006-01-02 15:04:05} [%p] - %m")
	jlw.BatchCount = 100
	jlw.CacheCount = 200
	jlw.FlushLevel = log.LevelError
	jlw.FlushDelta = time.Second * 2

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
	if jlw.EventBuffer.IsEmpty() {
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
		bs := jlw.Format(le)
		jl := &JobLog{
			JID:     jlw.jid,
			When:    le.When,
			Level:   le.Level.Prefix(),
			Message: string(bs),
		}
		jls = append(jls, jl)
	}

	r := jlw.db.Create(jls)
	return r.Error
}

// Close implementing method. empty.
func (jlw *JobLogWriter) Close() {
	jlw.Flush()
}
