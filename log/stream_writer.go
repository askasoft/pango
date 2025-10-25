package log

import (
	"io"
	"os"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/iox"
)

// StreamWriter implements log Writer Interface and writes messages to terminal.
type StreamWriter struct {
	FilterSupport
	FormatSupport

	Color  bool      // this field is useful only when system's terminal supports color.
	Output io.Writer // output writer. if nil, use os.Stdout as default or os.Stderr at Error level.
}

// Write write message to output writer.
// If Output is nil, use os.Stdout as default or os.Stderr at Error level.
func (sw *StreamWriter) Write(le *Event) {
	if sw.Reject(le) {
		return
	}

	out := sw.Output
	if out == nil {
		out = gog.If(le.Level > LevelError, os.Stdout, os.Stderr)
	}

	bs := sw.Format(le)
	if sw.Color {
		_, _ = out.Write(colors[le.Level])
		_, _ = out.Write(bs)
		_, _ = out.Write(colors[0])
	} else {
		_, _ = out.Write(bs)
	}
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

// NewConsoleWriter create a console log writer
func NewConsoleWriter(color ...bool) *StreamWriter {
	return &StreamWriter{Color: asg.First(color)}
}

// NewStdoutWriter create a stdout log writer
func NewStdoutWriter() *StreamWriter {
	return &StreamWriter{Output: os.Stdout}
}

// NewStderrWriter create a stderr writer
func NewStderrWriter() *StreamWriter {
	return &StreamWriter{Output: os.Stderr}
}

func init() {
	RegisterWriter("console", func() Writer {
		return NewConsoleWriter()
	})
	RegisterWriter("stdout", func() Writer {
		return NewStdoutWriter()
	})
	RegisterWriter("stderr", func() Writer {
		return NewStderrWriter()
	})
}
