package job

import (
	"github.com/askasoft/pango/log"
)

type JobMessage struct {
	Lvl string
	Msg string
}

// JobLogWriter implements log Writer Interface and writes messages to terminal.
type JobLogWriter struct {
	log.LogFilter
	log.LogFormatter

	Output []JobMessage // log output
}

// Write write message in console.
func (jlw *JobLogWriter) Write(le *log.Event) (err error) {
	if jlw.Reject(le) {
		return
	}

	bs := jlw.Format(le)

	jlw.Output = append(jlw.Output, JobMessage{
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
	jlw.Output = jlw.Output[:0]
}
