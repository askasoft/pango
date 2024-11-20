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

func (jr *JobRunner) GetJob() (*Job, error) {
	return jr.jmr.GetJob(jr.jid)
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

func (jr *JobRunner) Complete() error {
	return jr.jmr.CompleteJob(jr.jid)
}

func (jr *JobRunner) Ping() error {
	return jr.jmr.PingJob(jr.jid, jr.rid)
}

func (jr *JobRunner) Running(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if err := jr.Ping(); err != nil {
				return err
			}
			timer.Reset(interval)
		}
	}
}
