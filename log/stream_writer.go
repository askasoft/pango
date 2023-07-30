package log

import (
	"io"
	"os"

	"github.com/askasoft/pango/iox"
)

// StreamWriter implements log Writer Interface and writes messages to terminal.
type StreamWriter struct {
	LogFilter
	LogFormatter

	Color  bool      //this filed is useful only when system's terminal supports color
	Output io.Writer // log output
}

// Write write message in console.
func (sw *StreamWriter) Write(le *Event) (err error) {
	if sw.Reject(le) {
		return
	}

	if sw.Output == nil {
		sw.Output = os.Stdout
	}

	bs := sw.Format(le)
	if sw.Color {
		_, err = sw.Output.Write(colors[le.Level])
		if err != nil {
			return
		}
		_, err = sw.Output.Write(bs)
		if err != nil {
			return
		}
		_, err = sw.Output.Write(colors[0])
	} else {
		_, err = sw.Output.Write(bs)
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
