package log

import (
	"os"
)

// ConsoleWriter implements LogWriter Interface and writes messages to terminal.
type ConsoleWriter struct {
	Color  bool      `json:"color"` //this filed is useful only when system's terminal supports color
	Logfmt Formatter // log formatter
	Logfil Filter    // log filter
}

// SetFormat set a log formatter
func (cw *ConsoleWriter) SetFormat(format string) {
	cw.Logfmt = NewTextFormatter(format)
}

// Write write message in console.
func (cw *ConsoleWriter) Write(le *Event) {
	if cw.Logfil != nil && cw.Logfil.Reject(le) {
		return
	}
	if cw.Logfmt == nil {
		cw.Logfmt = le.Logger.GetFormatter()
	}
	msg := cw.Logfmt.Format(le)
	if cw.Color {
		msg = colors[le.Level](msg)
	}
	os.Stdout.WriteString(msg)
	return
}

// Flush implementing method. empty.
func (cw *ConsoleWriter) Flush() {
}

// Close implementing method. empty.
func (cw *ConsoleWriter) Close() {
}

// brush is a color join function
type brush func(string) string

// newBrush return a fix color Brush
func newBrush(color string) brush {
	pre := "\x1b[" + color + "m"
	reset := "\x1b[0m"
	return func(text string) string {
		return pre + text + reset
	}
}

var colors = []brush{
	newBrush("0"),  // None		reset
	newBrush("91"), // Fatal	red
	newBrush("95"), // Error	magenta
	newBrush("93"), // Warn		yellow
	newBrush("94"), // Info		blue
	newBrush("97"), // Debug	white
	newBrush("90"), // Trace	grey
}
