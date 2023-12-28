package log

// NewFailoverWriter create a failover writer
func NewFailoverWriter(w Writer, bufSize int) *FailoverWriter {
	fw := &FailoverWriter{writer: w}
	fw.evtbuf = NewEventBuffer(bufSize)
	return fw
}

// FailoverWriter wraps a log writer, cache log event if underlying log writer write failed.
// Resend the cached log events to the underlying log writer when the next log event come.
type FailoverWriter struct {
	writer Writer
	evtbuf *EventBuffer
}

// Write write event to underlying writer
func (fw *FailoverWriter) Write(le *Event) error {
	err := fw.flush()
	if err == nil {
		err = fw.writer.Write(le)
		if err != nil {
			perror(err)
		}
	}

	if err != nil {
		fw.evtbuf.Push(le)
	}

	return nil
}

// flush flush buffered event
func (fw *FailoverWriter) flush() error {
	for le, ok := fw.evtbuf.Peek(); ok; le, ok = fw.evtbuf.Peek() {
		if err := fw.writer.Write(le); err != nil {
			perror(err)
			return err
		}
		fw.evtbuf.Poll()
	}
	return nil
}

// Flush implementing method. empty.
func (fw *FailoverWriter) Flush() {
	fw.flush()
	fw.writer.Flush()
}

// Close implementing method. empty.
func (fw *FailoverWriter) Close() {
	fw.flush()
	fw.writer.Close()
}
