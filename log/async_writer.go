package log

import (
	"sync"
)

// NewAsyncWriter create a async writer and start go routine
func NewAsyncWriter(w Writer, size int) *AsyncWriter {
	aw := &AsyncWriter{writer: w}
	aw.Start(size)
	return aw
}

// AsyncWriter write log to multiple writers.
type AsyncWriter struct {
	writer  Writer
	evtChan chan *Event
	sigChan chan string
	waitg   sync.WaitGroup
	mutex   *sync.Mutex
}

// Write async write the log event
func (aw *AsyncWriter) Write(le *Event) {
	if aw.mutex != nil {
		aw.mutex.Lock()
		defer aw.mutex.Unlock()
		aw.writer.Write(le)
		return
	}

	aw.evtChan <- le
}

// Flush async flush the underlying writer
func (aw *AsyncWriter) Flush() {
	if aw.mutex != nil {
		aw.mutex.Lock()
		defer aw.mutex.Unlock()
		aw.writer.Flush()
		return
	}

	aw.sigChan <- "flush"
}

// Close Close the underlying writer and wait it for done
func (aw *AsyncWriter) Close() {
	if aw.mutex != nil {
		aw.mutex.Lock()
		defer aw.mutex.Unlock()
		aw.writer.Close()
		return
	}

	aw.sigChan <- "close"
	aw.waitg.Wait()
}

// Start start the goroutine
func (aw *AsyncWriter) Start(size int) {
	aw.evtChan = make(chan *Event, size)
	aw.sigChan = make(chan string, 1)
	aw.waitg.Add(1)
	go aw.run()
}

// SetWriter set the log writer
func (aw *AsyncWriter) SetWriter(w Writer) {
	aw.writer = w
}

// SyncWriter switch to synchronized mode and replace the writer
func (aw *AsyncWriter) SyncWriter(w Writer) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	aw.mutex = mutex
	aw.writer = w

	// complete drain the event chan
	ec := aw.evtChan
	for len(ec) > 0 {
		le := <-ec
		w.Write(le)
	}

	// complete drain the signal chan
	sc := aw.sigChan
	for len(sc) > 0 {
		_ = <-sc
	}
}

// drain drain the event chan once (ignore after-coming event to prevent dead loop)
func (aw *AsyncWriter) drain() {
	ec := aw.evtChan
	for n := len(ec); n > 0; n-- {
		le := <-ec
		aw.writer.Write(le)
	}
}

// run start async log goroutine
func (aw *AsyncWriter) run() {
	done := false
	for {
		select {
		case sg := <-aw.sigChan:
			aw.drain()
			switch sg {
			case "flush":
				aw.writer.Flush()
			case "close":
				aw.writer.Close()
				done = true
			}
		case le := <-aw.evtChan:
			aw.writer.Write(le)
		}
		if done {
			break
		}
	}

	aw.waitg.Done()
}
