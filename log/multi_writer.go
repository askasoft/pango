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
func (mw *MultiWriter) Write(le *Event) error {
	for _, w := range mw.Writers {
		safeWrite(w, le)
	}
	return nil
}

// Close close multiple writers.
func (mw *MultiWriter) Close() {
	for _, w := range mw.Writers {
		safeClose(w)
	}
}

// Flush flush multiple writers.
func (mw *MultiWriter) Flush() {
	for _, w := range mw.Writers {
		safeFlush(w)
	}
}
