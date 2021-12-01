package log

var nopWriter = &NopWriter{}

// NopWriter implements Writer.
// Do nothing.
type NopWriter struct {
}

// Write do nothing.
func (nw *NopWriter) Write(le *Event) {
}

// Flush do nothing.
func (nw *NopWriter) Flush() {
}

// Close do nothing.
func (nw *NopWriter) Close() {
}
