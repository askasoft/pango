package log

import (
	"bytes"
	"io"
	"os"

	"github.com/askasoft/pango/iox"
)

// StreamWriter implements log Writer Interface and writes messages to terminal.
type StreamWriter struct {
	Color  bool         //this filed is useful only when system's terminal supports color
	Output io.Writer    // log output
	Logfmt Formatter    // log formatter
	Logfil Filter       // log filter
	bb     bytes.Buffer // log buffer
}

// SetFormat set the log formatter
func (sw *StreamWriter) SetFormat(format string) {
	sw.Logfmt = NewLogFormatter(format)
}

// SetFilter set the log filter
func (sw *StreamWriter) SetFilter(filter string) {
	sw.Logfil = NewLogFilter(filter)
}

// Write write message in console.
func (sw *StreamWriter) Write(le *Event) (err error) {
	if sw.Logfil != nil && sw.Logfil.Reject(le) {
		return
	}

	if sw.Output == nil {
		sw.Output = os.Stdout
	}

	lf := sw.Logfmt
	if lf == nil {
		lf = le.Logger().GetFormatter()
		if lf == nil {
			lf = TextFmtDefault
		}
	}

	sw.bb.Reset()
	lf.Write(&sw.bb, le)
	if sw.Color {
		_, err = sw.Output.Write(colors[le.Level()])
		if err != nil {
			return
		}
		_, err = sw.Output.Write(sw.bb.Bytes())
		if err != nil {
			return
		}
		_, err = sw.Output.Write(colors[0])
	} else {
		_, err = sw.Output.Write(sw.bb.Bytes())
	}
	return
}

// Flush implementing method. empty.
func (sw *StreamWriter) Flush() {
}

// Close implementing method. empty.
func (sw *StreamWriter) Close() {
}

var colors = [][]byte{
	[]byte(iox.ConsoleColor.Reset),   // None
	[]byte(iox.ConsoleColor.Red),     // Fatal
	[]byte(iox.ConsoleColor.Magenta), // Error
	[]byte(iox.ConsoleColor.Yellow),  // Warn
	[]byte(iox.ConsoleColor.Blue),    // Info
	[]byte(iox.ConsoleColor.White),   // Debug
	[]byte(iox.ConsoleColor.Gray),    // Trace
}

// NewConsoleWriter create a color console log writer
func NewConsoleWriter() Writer {
	return &StreamWriter{Output: os.Stdout, Color: true}
}

// NewStdoutWriter create a stdout log writer
func NewStdoutWriter() Writer {
	return &StreamWriter{Output: os.Stdout}
}

// NewStderrWriter create a stderr writer
func NewStderrWriter() Writer {
	return &StreamWriter{Output: os.Stderr}
}

func init() {
	RegisterWriter("console", NewConsoleWriter)
	RegisterWriter("stdout", NewStdoutWriter)
	RegisterWriter("stderr", NewStderrWriter)
}
