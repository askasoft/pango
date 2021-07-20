package log

import "sync"

// NewSyncWriter create a sync writer
func NewSyncWriter(w Writer) Writer {
	sw := &syncWriter{writer: w}
	return sw
}

// syncWriter synchronized log writer
type syncWriter struct {
	writer Writer
	mutex  sync.Mutex
}

func (sw *syncWriter) Write(le *Event) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	sw.writer.Write(le)
}

func (sw *syncWriter) Flush() {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	sw.writer.Flush()
}

func (sw *syncWriter) Close() {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	sw.writer.Close()
}
