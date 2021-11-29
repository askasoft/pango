package log

import "sync"

// NewAsyncWriter create a async writer
func NewAsyncWriter(w Writer, size int) *AsyncWriter {
	aw := &AsyncWriter{writer: w}
	aw.Start(size)
	return aw
}

// AsyncWriter write log to multiple writers.
type AsyncWriter struct {
	writer Writer

	evtChan chan *Event
	sigChan chan string
	waitg   sync.WaitGroup
}

// Write async write the log event
func (aw *AsyncWriter) Write(le *Event) {
	ec := aw.evtChan
	if ec != nil {
		ec <- le
	}
}

// Flush async flush the underlying writer
func (aw *AsyncWriter) Flush() {
	sc := aw.sigChan
	if sc != nil {
		sc <- "flush"
	}
}

// Close Close the underlying writer and wait it for done
func (aw *AsyncWriter) Close() {
	sc := aw.sigChan
	if sc != nil {
		sc <- "close"
	}
	aw.waitg.Wait()
}

// Start start the goroutine
func (aw *AsyncWriter) Start(size int) {
	aw.evtChan = make(chan *Event, size)
	aw.sigChan = make(chan string, 1)
	aw.waitg.Add(1)
	go aw.run()
}

func (aw *AsyncWriter) drain() {
	for len(aw.evtChan) > 0 {
		le := <-aw.evtChan
		aw.writer.Write(le)
	}
}

// run start async log goroutine
func (aw *AsyncWriter) run() {
	done := false
	for {
		select {
		case le := <-aw.evtChan:
			aw.writer.Write(le)
		case sg := <-aw.sigChan:
			aw.drain()
			switch sg {
			case "flush":
				aw.writer.Flush()
			case "close":
				aw.writer.Close()
				done = true
			}
		}
		if done {
			break
		}
	}

	ec, sc := aw.evtChan, aw.sigChan
	aw.evtChan, aw.sigChan = nil, nil
	close(ec)
	close(sc)

	aw.waitg.Done()
}
