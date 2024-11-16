package log

import "sync"

// NewSyncWriter create a synchronized writer
func NewSyncWriter(w Writer) *SyncWriter {
	sw := &SyncWriter{writer: w}
	return sw
}

// SyncWriter synchronized log writer
type SyncWriter struct {
	writer Writer
	mutex  sync.Mutex
}

// Write synchronize write log event
func (sw *SyncWriter) Write(le *Event) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	sw.writer.Write(le)
}

// Flush synchronize flush the underlying writer
func (sw *SyncWriter) Flush() {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	sw.writer.Flush()
}

// Close synchronize close the underlying writer
func (sw *SyncWriter) Close() {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	sw.writer.Close()
	sw.writer = nopWriter
}

// SetWriter synchronize close the old log writer then set the new log writer
func (sw *SyncWriter) SetWriter(w Writer) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	sw.writer.Close()
	sw.writer = w
}
