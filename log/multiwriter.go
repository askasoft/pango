package log

// NewMultiWriter create a multi writer
func NewMultiWriter(ws ...Writer) *MultiWriter {
	return &MultiWriter{Writers: ws}
}

// MultiWriter write log to multiple writers.
type MultiWriter struct {
	Writers []Writer
}

// Write write log event to multiple writers.
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
