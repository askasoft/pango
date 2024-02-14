package xjm

import (
	"time"

	"github.com/askasoft/pango/log"
)

type JobRunner struct {
	Log *log.Log

	jmr JobManager
	jlw *JobLogWriter

	jid int64 // Job ID
	rid int64 // Runner ID

	pingAt    time.Time
	PingAfter time.Duration // Ping after duration
}

// NewJobRunner create a JobRunner
func NewJobRunner(jmr JobManager, jid, rid int64, logger ...log.Logger) *JobRunner {
	jr := &JobRunner{
		Log:       log.NewLog(),
		jmr:       jmr,
		jid:       jid,
		rid:       rid,
		PingAfter: time.Second,
	}

	jr.Log.SetFormat("%t{2006-01-02 15:04:05} [%p] - %m")
	jr.jlw = NewJobLogWriter(jmr, jid)
	if len(logger) > 0 {
		bw := log.NewBridgeWriter(logger[0])
		mw := log.NewMultiWriter(jr.jlw, bw)
		jr.Log.SetWriter(mw)
	} else {
		jr.Log.SetWriter(jr.jlw)
	}

	return jr
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

func (jr *JobRunner) Ping() error {
	if jr.pingAt.Add(jr.PingAfter).After(time.Now()) {
		return nil
	}

	if err := jr.jmr.PingJob(jr.jid, jr.rid); err != nil {
		return err
	}

	jr.pingAt = time.Now()
	return nil
}

func (jr *JobRunner) PingAborted() bool {
	return jr.Ping() != nil
}

func (jr *JobRunner) Running(result string) error {
	if jr.pingAt.Add(jr.PingAfter).After(time.Now()) {
		return nil
	}

	if err := jr.jmr.RunningJob(jr.jid, jr.rid, result); err != nil {
		return err
	}

	jr.pingAt = time.Now()
	return nil
}

func (jr *JobRunner) Abort(reason string) error {
	err := jr.jmr.AbortJob(jr.jid, reason)
	jr.Log.Flush()
	return err
}

func (jr *JobRunner) Complete(result string) error {
	err := jr.jmr.CompleteJob(jr.jid, result)
	jr.Log.Flush()
	return err
}
