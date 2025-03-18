package xjm

import (
	"context"
	"time"

	"github.com/askasoft/pango/log"
)

type JobRunner struct {
	job *Job
	jmr JobManager
	jlw *JobLogWriter
	log *log.Log
}

// NewJobRunner create a JobRunner
func NewJobRunner(job *Job, jmr JobManager, logger ...log.Logger) *JobRunner {
	jr := &JobRunner{
		job: job,
		jmr: jmr,
		log: log.NewLog(),
	}

	jr.jlw = NewJobLogWriter(jmr, job.ID)

	var lw log.Writer = jr.jlw
	if len(logger) > 0 {
		lw = log.NewMultiWriter(jr.jlw, log.NewBridgeWriter(logger[0]))
	}

	jr.log.SetWriter(log.NewAsyncWriter(lw, 100))
	return jr
}

func (jr *JobRunner) Log() *log.Log {
	return jr.log
}

func (jr *JobRunner) JobManager() JobManager {
	return jr.jmr
}

func (jr *JobRunner) JobLogWriter() *JobLogWriter {
	return jr.jlw
}

func (jr *JobRunner) JobID() int64 {
	return jr.job.ID
}

func (jr *JobRunner) ChainID() int64 {
	return jr.job.CID
}

func (jr *JobRunner) RunnerID() int64 {
	return jr.job.RID
}

func (jr *JobRunner) Locale() string {
	return jr.job.Locale
}

func (jr *JobRunner) JobName() string {
	return jr.job.Name
}

func (jr *JobRunner) JobParam() string {
	return jr.job.Param
}

func (jr *JobRunner) GetJob(cols ...string) (*Job, error) {
	return jr.jmr.GetJob(jr.job.ID, cols...)
}

func (jr *JobRunner) Checkout() error {
	return jr.jmr.CheckoutJob(jr.job.ID, jr.job.RID)
}

func (jr *JobRunner) SetState(state string) error {
	return jr.jmr.SetJobState(jr.job.ID, jr.job.RID, state)
}

func (jr *JobRunner) AddResult(result string) error {
	return jr.jmr.AddJobResult(jr.job.ID, jr.job.RID, result)
}

func (jr *JobRunner) Abort(reason string) error {
	return jr.jmr.AbortJob(jr.job.ID, reason)
}

func (jr *JobRunner) Cancel(reason string) error {
	return jr.jmr.CancelJob(jr.job.ID, reason)
}

func (jr *JobRunner) Finish() error {
	return jr.jmr.FinishJob(jr.job.ID)
}

func (jr *JobRunner) Pin() error {
	return jr.jmr.PinJob(jr.job.ID, jr.job.RID)
}

func (jr *JobRunner) Running(ctx context.Context, getTimeout, pinTimeout time.Duration) error {
	gettm, pintm := time.NewTimer(getTimeout), time.NewTimer(pinTimeout)

	defer func() {
		gettm.Stop()
		pintm.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-gettm.C:
			job, err := jr.GetJob("id", "rid", "status")
			if err != nil {
				return err
			}
			if job.RID != jr.job.RID || job.Status != JobStatusRunning {
				return ErrJobPin
			}
			gettm.Reset(getTimeout)
		case <-pintm.C:
			if err := jr.Pin(); err != nil {
				return err
			}
			pintm.Reset(pinTimeout)
		}
	}
}
