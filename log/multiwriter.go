package log

// MultiWriter write messages to multiple writers.
type MultiWriter struct {
	Writers []Writer
}

// Write write message in console.
func (mw *MultiWriter) Write(le *Event) {
	for _, w := range mw.Writers {
		w.Write(le)
	}
}

// Close implementing method. empty.
func (mw *MultiWriter) Close() {
	for _, w := range mw.Writers {
		w.Close()
	}
}

// Flush implementing method. empty.
func (mw *MultiWriter) Flush() {
	for _, w := range mw.Writers {
		w.Flush()
	}
}
