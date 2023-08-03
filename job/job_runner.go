package job

import (
	"errors"
	"sync/atomic"

	"github.com/askasoft/pango/log"
)

var ErrJobAborted = errors.New("Aborted")

func NewJobRunner(run func(*JobRunner)) *JobRunner {
	jr := &JobRunner{
		run: run,
		Log: log.NewLog(),
	}

	jr.Log.SetLevel(log.LevelInfo)
	jr.Log.SetWriter(&jr.outputs)
	jr.Log.SetFormatter(log.NewTextFormatter("%t{2006-01-02 15:04:05} [%p] - %m%n%T"))
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
	jr.outputs.Clear()

	jr.run(jr)
}
func (jr *JobRunner) GetOutputs(skip int) []JobMessage {
	if skip > len(jr.outputs.Output) {
		return nil
	}

	return jr.outputs.Output[skip:]
}
