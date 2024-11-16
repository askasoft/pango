package log

import "github.com/askasoft/pango/log/internal"

// RetrySupport cache event if write failed, and retry write when the next log event come.
type RetrySupport struct {
	Retries     int
	RetryBuffer EventBuffer
}

func (rs *RetrySupport) RetryWrite(le *Event, write func(*Event) error) {
	err := rs.retry(write)
	if err == nil {
		err = write(le)
		if err != nil {
			internal.Perror(err)
		}
	}

	if err != nil {
		rs.RetryBuffer.Push(le)
		if rs.RetryBuffer.Len() > rs.Retries {
			rs.RetryBuffer.Poll()
		}
	}
}

func (rs *RetrySupport) RetryFlush(write func(*Event) error) {
	_ = rs.retry(write)
}

func (rs *RetrySupport) retry(write func(*Event) error) error {
	for le, ok := rs.RetryBuffer.Peek(); ok; le, ok = rs.RetryBuffer.Peek() {
		if err := write(le); err != nil {
			internal.Perror(err)
			return err
		}
		rs.RetryBuffer.Poll()
	}
	return nil
}
