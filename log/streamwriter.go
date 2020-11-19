package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

// StreamWriter implements log Writer Interface and writes messages to terminal.
type StreamWriter struct {
	Color  bool         //this filed is useful only when system's terminal supports color
	Output io.Writer    // log output
	Logfmt Formatter    // log formatter
	Logfil Filter       // log filter
	bb     bytes.Buffer // log buffer
}

// SetFormat set a log formatter
func (sw *StreamWriter) SetFormat(format string) {
	sw.Logfmt = NewTextFormatter(format)
}

// SetColor set a log formatter
func (sw *StreamWriter) SetColor(color string) error {
	clr, err := strconv.ParseBool(color)
	if err != nil {
		return fmt.Errorf("Invalid Color: %v", err)
	}
	sw.Color = clr
	return nil
}

// Write write message in console.
func (sw *StreamWriter) Write(le *Event) {
	if sw.Logfil != nil && sw.Logfil.Reject(le) {
		return
	}

	le.Logger.Lock()
	defer le.Logger.Unlock()

	if sw.Output == nil {
		sw.Output = os.Stdout
	}
	if sw.Logfmt == nil {
		sw.Logfmt = le.Logger.GetFormatter()
	}

	sw.bb.Reset()
	sw.Logfmt.Write(&sw.bb, le)
	if sw.Color {
		sw.Output.Write([]byte("\x1b["))
		sw.Output.Write([]byte(colors[le.Level]))
		sw.Output.Write([]byte("m"))
		sw.Output.Write(sw.bb.Bytes())
		sw.Output.Write([]byte("\x1b[0m"))
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
	"0",  // None		reset
	"91", // Fatal	red
	"95", // Error	magenta
	"93", // Warn		yellow
	"94", // Info		blue
	"97", // Debug	white
	"90", // Trace	grey
}
