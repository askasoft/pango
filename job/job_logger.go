package job

import (
	"github.com/askasoft/pango/log"
)

type JobLogMsg struct {
	Lvl string
	Msg string
}

// JobLogWriter implements log Writer Interface and writes messages to terminal.
type JobLogWriter struct {
	log.LogFilter
	log.LogFormatter

	msgs []JobLogMsg
}

// Write write message in console.
func (jlw *JobLogWriter) Write(le *log.Event) (err error) {
	if jlw.Reject(le) {
		return
	}

	bs := jlw.Format(le)

	jlw.msgs = append(jlw.msgs, JobLogMsg{
		Lvl: le.Level.Prefix(),
		Msg: string(bs),
	})
	return
}

// Flush implementing method. empty.
func (jlw *JobLogWriter) Flush() {
}

// Close implementing method. empty.
func (jlw *JobLogWriter) Close() {
}

// Clear clear the output
func (jlw *JobLogWriter) Clear() {
	jlw.msgs = jlw.msgs[:0]
}

// GetMessage get job log messages from start to start+limit.
// GetMessage(0, -1): get all message
func (jlw *JobLogWriter) GetMessages(start, limit int) []JobLogMsg {
	msgs := jlw.msgs
	mlen := len(msgs)

	if start > mlen {
		return nil
	}

	end := mlen
	if limit >= 0 {
		end = start + limit
		if end > mlen {
			end = mlen
		}
	}

	return msgs[start:end]
}
