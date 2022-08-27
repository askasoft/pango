package log

import (
	"sync"
	"time"
)

// NewAsyncWriter create a async writer and start go routine
func NewAsyncWriter(w Writer, size int) *AsyncWriter {
	aw := &AsyncWriter{writer: w}
	aw.Start(size)
	return aw
}

type signal struct {
	signal string
	option any
}

// AsyncWriter wrapper a log writer to implement asynchrous write
type AsyncWriter struct {
	writer  Writer
	evtChan chan *Event
	sigChan chan signal
	waitg   sync.WaitGroup
}

// Write async write the log event
func (aw *AsyncWriter) Write(le *Event) error {
	aw.evtChan <- le
	return nil
}

// Flush async flush the underlying writer
func (aw *AsyncWriter) Flush() {
	aw.sigChan <- signal{"flush", nil}
}

// Close Close the underlying writer and wait it for done
func (aw *AsyncWriter) Close() {
	aw.sigChan <- signal{"close", nil}
	aw.waitg.Wait()
}

// SetWriter close the old writer and set the new writer
func (aw *AsyncWriter) SetWriter(w Writer) {
	aw.sigChan <- signal{"switch", w}
}

// Start start the goroutine
func (aw *AsyncWriter) Start(size int) {
	aw.evtChan = make(chan *Event, size)
	aw.sigChan = make(chan signal, 1)
	aw.waitg.Add(1)
	go aw.run()
}

// Stop stop the run() go-routine
func (aw *AsyncWriter) Stop() {
	aw.sigChan <- signal{"stop", nil}
	aw.waitg.Wait()
}

func (aw *AsyncWriter) write(w Writer, le *Event) {
	err := w.Write(le)
	if err != nil {
		perror(err)
	}
}

// drainOnce drain the event chan once (ignore after-coming event to prevent dead loop)
func (aw *AsyncWriter) drainOnce() {
	w := aw.writer

	ec := aw.evtChan
	for n := len(ec); n > 0; n-- {
		le := <-ec
		aw.write(w, le)
	}
}

// drainFull complete drain the event chan
func (aw *AsyncWriter) drainFull() {
	w := aw.writer

	// complete drain the event chan
	ec := aw.evtChan
	for len(ec) > 0 {
		le := <-ec
		aw.write(w, le)
	}

	// complete drain the signal chan
	sc := aw.sigChan
	for len(sc) > 0 {
		<-sc
	}
}

// run start async log goroutine
func (aw *AsyncWriter) run() {
	done := false
	for {
		select {
		case sg := <-aw.sigChan:
			switch sg.signal {
			case "switch":
				aw.drainFull()
				ow := aw.writer
				aw.writer = sg.option.(Writer)
				ow.Close()
			case "flush":
				aw.drainOnce()
				aw.writer.Flush()
			case "close":
				aw.drainFull()
				aw.writer.Close()
				done = true
			case "stop":
				aw.drainFull()
				done = true
			}
		case le := <-aw.evtChan:
			aw.write(aw.writer, le)
		}
		if done {
			break
		}
	}

	close(aw.evtChan)
	close(aw.sigChan)
	aw.waitg.Done()
}

// StopAfter auto stop the run() go-routine when the evtChan is empty and after duration d.
func (aw *AsyncWriter) StopAfter(d time.Duration) {
	timer := time.NewTimer(d)
	go func() {
		for {
			<-timer.C
			if len(aw.evtChan) == 0 && len(aw.sigChan) == 0 {
				aw.Stop()
				return
			}

			timer.Reset(d)
		}
	}()
}
