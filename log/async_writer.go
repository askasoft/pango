package log

import "sync"

// AsyncWriter async log writer interface
type AsyncWriter interface {
	Writer

	SyncClose()
}

// NewAsyncWriter create a async writer
func NewAsyncWriter(w Writer, size int) AsyncWriter {
	aw := &asyncWriter{writer: w}
	aw.startAsync(size)
	return aw
}

// asyncWriter write log to multiple writers.
type asyncWriter struct {
	writer Writer

	evtChan chan *Event
	sigChan chan string
	waitg   sync.WaitGroup
}

func (aw *asyncWriter) Write(le *Event) {
	ec := aw.evtChan
	if ec != nil {
		ec <- le
	}
}

func (aw *asyncWriter) Flush() {
	sc := aw.sigChan
	if sc != nil {
		sc <- "flush"
	}
}

func (aw *asyncWriter) Close() {
	sc := aw.sigChan
	if sc != nil {
		sc <- "close"
	}
}

// SyncClose Close the writer and wait for done
func (aw *asyncWriter) SyncClose() {
	aw.Close()
	aw.waitg.Wait()
}

// startAsync set the log to asynchronous and start the goroutine
func (aw *asyncWriter) startAsync(size int) {
	aw.evtChan = make(chan *Event, size)
	aw.sigChan = make(chan string, 1)
	aw.waitg.Add(1)
	go aw.run()
}

func (aw *asyncWriter) drain() {
	for len(aw.evtChan) > 0 {
		le := <-aw.evtChan
		aw.writer.Write(le)
	}
}

// run start async log goroutine
func (aw *asyncWriter) run() {
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
			case "done":
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
