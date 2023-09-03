package job

import (
	"errors"
	"sync/atomic"

	"github.com/askasoft/pango/log"
)

var ErrJobAborted = errors.New("Aborted")

func NewJobRunner(run func(*JobRunner), logger ...log.Logger) *JobRunner {
	jr := &JobRunner{
		run: run,
		Log: log.NewLog(),
	}

	jr.outputs.SetFilter("level:DEBUG")
	jr.outputs.SetFormat("%t{2006-01-02 15:04:05} [%p] - %m%n%T")

	if len(logger) > 0 {
		bw := log.NewBridgeWriter(logger[0])
		mw := log.NewMultiWriter(&jr.outputs, bw)
		jr.Log.SetWriter(mw)
	} else {
		jr.Log.SetWriter(&jr.outputs)
	}

	return jr
}

type JobRunner struct {
	Log     *log.Log
	run     func(*JobRunner)
	aborted int32
	running int32
	outputs JobLogWriter
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
	go jr.Run()
}

func (jr *JobRunner) Run() {
	defer atomic.StoreInt32(&jr.running, 0)

	atomic.StoreInt32(&jr.aborted, 0)
	atomic.StoreInt32(&jr.running, 1)

	jr.ClearOutput()

	jr.run(jr)
}

func (jr *JobRunner) ClearOutput() {
	jr.outputs.Clear()
}

func (jr *JobRunner) GetOutputs(skip int) []JobMessage {
	if skip > len(jr.outputs.Output) {
		return nil
	}

	return jr.outputs.Output[skip:]
}
