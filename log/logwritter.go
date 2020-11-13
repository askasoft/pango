package log

// Writer defines the behavior of a log writer.
type Writer interface {
	Write(le *Event)
	Close()
	Flush()
}
