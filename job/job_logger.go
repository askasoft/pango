package job

import (
	"bytes"

	"github.com/askasoft/pango/log"
)

type JobMessage struct {
	Lvl string
	Msg string
}

// JobLogWriter implements log Writer Interface and writes messages to terminal.
type JobLogWriter struct {
	Output []JobMessage // log output
	bb     bytes.Buffer
}

// Write write message in console.
func (jlw *JobLogWriter) Write(le *log.Event) (err error) {
	lf := le.Logger.GetFormatter()
	if lf == nil {
		lf = log.TextFmtDefault
	}

	jlw.bb.Reset()
	lf.Write(&jlw.bb, le)
	jlw.Output = append(jlw.Output, JobMessage{
		Lvl: le.Level.Prefix(),
		Msg: jlw.bb.String(),
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
