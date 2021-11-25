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

// Close close multiple writers.
func (mw *MultiWriter) Close() {
	for _, w := range mw.Writers {
		w.Close()
	}
}

// SyncClose Close the multiple writers and wait them for done
func (mw *MultiWriter) SyncClose() {
	for _, w := range mw.Writers {
		if aw, ok := w.(AsyncWriter); ok {
			aw.SyncClose()
		} else {
			w.Close()
		}
	}
}

// Flush flush multiple writers.
func (mw *MultiWriter) Flush() {
	for _, w := range mw.Writers {
		w.Flush()
	}
}
