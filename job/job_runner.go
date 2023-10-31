package job

import (
	"errors"
	"sync/atomic"

	"github.com/askasoft/pango/log"
)

var ErrJobAborted = errors.New("Aborted")

func NewJobRunner(run func(), logger ...log.Logger) *JobRunner {
	jr := &JobRunner{
		Run: run,
		Log: log.NewLog(),
	}

	jr.Out.SetFilter("level:DEBUG")
	jr.Out.SetFormat("%t{2006-01-02 15:04:05} [%p] - %m%n")

	if len(logger) > 0 {
		bw := log.NewBridgeWriter(logger[0])
		mw := log.NewMultiWriter(&jr.Out, bw)
		jr.Log.SetWriter(mw)
	} else {
		jr.Log.SetWriter(&jr.Out)
	}

	return jr
}

type JobRunner struct {
	Run     func()
	Log     *log.Log
	Out     JobLogWriter
	aborted int32
	running int32
}

func (jr *JobRunner) IsRunning() bool {
	return atomic.LoadInt32(&jr.running) != 0
}

func (jr *JobRunner) IsAborted() bool {
	return atomic.LoadInt32(&jr.aborted) != 0
}

func (jr *JobRunner) Abort() {
	atomic.StoreInt32(&jr.aborted, 1)
}

func (jr *JobRunner) Start() {
	if jr.IsRunning() {
		return
	}

	jr.Out.Clear()

	atomic.StoreInt32(&jr.aborted, 0)
	atomic.StoreInt32(&jr.running, 1)

	go jr.safeRun()
}

func (jr *JobRunner) safeRun() {
	defer func() {
		if err := recover(); err != nil {
			jr.Log.Errorf("Panic: %v", err)
		}
		atomic.StoreInt32(&jr.running, 0)
	}()

	jr.Run()
}
