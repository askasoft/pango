package log

import (
	"sync"
)

// NewAsyncWriter create a async writer and start go routine
func NewAsyncWriter(w Writer, size int) *AsyncWriter {
	aw := &AsyncWriter{
		writer:  w,
		evtChan: make(chan *Event, size),
		sigChan: make(chan signal, 4),
	}

	aw.Start()
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
func (aw *AsyncWriter) Write(le *Event) {
	aw.evtChan <- le
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

// SetWriter send a "switch" signal to switch the writer to `w` and close the old writer
func (aw *AsyncWriter) SetWriter(w Writer) {
	aw.sigChan <- signal{"switch", w}
}

// Start start async log go-routine
func (aw *AsyncWriter) Start() {
	aw.waitg.Add(1)
	go aw.run()
}

// Stop send a "stop" signal to the run() go-routine
func (aw *AsyncWriter) Stop() {
	aw.sigChan <- signal{"stop", nil}
}

// Wait wait for the run() go-routine exit
func (aw *AsyncWriter) Wait() {
	aw.waitg.Wait()
}

func (aw *AsyncWriter) run() {
	stop, done := false, false

	defer func() {
		if done {
			aw.writer.Close()
		}

		aw.waitg.Done()

		// Send to a closed channel will cause panic, and it's safe to keep channels open.
		// GC will collect the unreachable channels.
		// close(aw.evtChan)
		// close(aw.sigChan)
	}()

	for {
		select {
		case sg := <-aw.sigChan:
			switch sg.signal {
			case "switch":
				ow := aw.writer
				aw.writer = sg.option.(Writer)
				ow.Close()
			case "flush":
				aw.writer.Flush()
			case "close":
				stop = true
				done = true
			case "stop":
				stop = true
			}
		case le := <-aw.evtChan:
			safeWrite(aw.writer, le)
		default:
			if stop && len(aw.evtChan) == 0 && len(aw.sigChan) == 0 {
				return
			}
		}
	}
}
