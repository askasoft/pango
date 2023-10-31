package job

import (
	"testing"
	"time"
)

type myJobRunner struct {
	*JobRunner

	count int
}

func (mjr *myJobRunner) Run() {
	mjr.count++

	mjr.Log.Debugf("run %d", mjr.count)
}

func TestNewJobRunner(t *testing.T) {
	mjr := &myJobRunner{}
	mjr.JobRunner = NewJobRunner(mjr.Run)

	var job Job = mjr
	job.Start()

	for job.IsRunning() {
		time.Sleep(time.Millisecond * 100)
	}

	if mjr.count != 1 {
		t.Error("my job not run")
	}
}
