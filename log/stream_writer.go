package log

import (
	"bytes"
	"io"
	"os"

	"github.com/pandafw/pango/iox"
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
func (sw *StreamWriter) Write(le *Event) {
	if sw.Logfil != nil && sw.Logfil.Reject(le) {
		return
	}

	if sw.Output == nil {
		sw.Output = os.Stdout
	}
	if sw.Logfmt == nil {
		sw.Logfmt = le.Logger.GetFormatter()
		if sw.Logfmt == nil {
			sw.Logfmt = TextFmtDefault
		}
	}

	sw.bb.Reset()
	sw.Logfmt.Write(&sw.bb, le)
	if sw.Color {
		sw.Output.Write([]byte(colors[le.Level]))
		sw.Output.Write(sw.bb.Bytes())
		sw.Output.Write([]byte(colors[0]))
	} else {
		sw.Output.Write(sw.bb.Bytes())
	}
}

// Flush implementing method. empty.
func (sw *StreamWriter) Flush() {
}

// Close implementing method. empty.
func (sw *StreamWriter) Close() {
}

var colors = []string{
	iox.ConsoleColor.Reset,   // None
	iox.ConsoleColor.Red,     // Fatal
	iox.ConsoleColor.Magenta, // Error
	iox.ConsoleColor.Yellow,  // Warn
	iox.ConsoleColor.Blue,    // Info
	iox.ConsoleColor.White,   // Debug
	iox.ConsoleColor.Gray,    // Trace
}

func init() {
	RegisterWriter("console", func() Writer {
		return &StreamWriter{Output: os.Stdout, Color: true}
	})
	RegisterWriter("stdout", func() Writer {
		return &StreamWriter{Output: os.Stdout}
	})
	RegisterWriter("stderr", func() Writer {
		return &StreamWriter{Output: os.Stderr}
	})
}
