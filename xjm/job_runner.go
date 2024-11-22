package xjm

import (
	"context"
	"time"

	"github.com/askasoft/pango/log"
)

type JobRunner struct {
	log *log.Log

	jmr JobManager
	jlw *JobLogWriter

	jnm string // Job Name
	jid int64  // Job ID
	rid int64  // Runner ID
}

// NewJobRunner create a JobRunner
func NewJobRunner(jmr JobManager, jnm string, jid, rid int64, logger ...log.Logger) *JobRunner {
	jr := &JobRunner{
		log: log.NewLog(),
		jmr: jmr,
		jnm: jnm,
		jid: jid,
		rid: rid,
	}

	jr.jlw = NewJobLogWriter(jmr, jid)

	var lw log.Writer = jr.jlw
	if len(logger) > 0 {
		lw = log.NewMultiWriter(jr.jlw, log.NewBridgeWriter(logger[0]))
	}

	jr.log.SetWriter(log.NewSyncWriter(lw))
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

func (jr *JobRunner) JobName() string {
	return jr.jnm
}

func (jr *JobRunner) JobID() int64 {
	return jr.jid
}

func (jr *JobRunner) RunnerID() int64 {
	return jr.rid
}

func (jr *JobRunner) GetJob(cols ...string) (*Job, error) {
	return jr.jmr.GetJob(jr.jid, cols...)
}

func (jr *JobRunner) Checkout() error {
	return jr.jmr.CheckoutJob(jr.jid, jr.rid)
}

func (jr *JobRunner) SetState(state string) error {
	return jr.jmr.SetJobState(jr.jid, jr.rid, state)
}

func (jr *JobRunner) AddResult(result string) error {
	return jr.jmr.AddJobResult(jr.jid, jr.rid, result)
}

func (jr *JobRunner) Abort(reason string) error {
	return jr.jmr.AbortJob(jr.jid, reason)
}

func (jr *JobRunner) Cancel(reason string) error {
	return jr.jmr.CancelJob(jr.jid, reason)
}

func (jr *JobRunner) Finish() error {
	return jr.jmr.FinishJob(jr.jid)
}

func (jr *JobRunner) Pin() error {
	return jr.jmr.PinJob(jr.jid, jr.rid)
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
			if job.RID != jr.rid || job.Status != JobStatusRunning {
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
