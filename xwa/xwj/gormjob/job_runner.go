package gormjob

import (
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
	"gorm.io/gorm"
)

type JobRunner struct {
	Log *log.Log
	jlw *JobLogWriter

	DB       *gorm.DB
	JobTable string
	LogTable string

	jid int64 // Job ID
	rid int64 // Runner ID

	pingAt  time.Time
	Timeout time.Duration
}

func NewJobRunner(db *gorm.DB, jobTable, logTable string, jid, rid int64, logger ...log.Logger) *JobRunner {
	jr := &JobRunner{
		Log:      log.NewLog(),
		DB:       db,
		JobTable: jobTable,
		LogTable: logTable,
		jid:      jid,
		rid:      rid,
		Timeout:  time.Second,
	}

	jr.jlw = NewJobLogWriter(db, logTable, jid)
	if len(logger) > 0 {
		bw := log.NewBridgeWriter(logger[0])
		mw := log.NewMultiWriter(jr.jlw, bw)
		jr.Log.SetWriter(mw)
	} else {
		jr.Log.SetWriter(jr.jlw)
	}

	return jr
}

func (jr *JobRunner) SetLogFormat(format string) {
	jr.jlw.SetFormat(format)
}

func (jr *JobRunner) JobID() int64 {
	return jr.jid
}

func (jr *JobRunner) RunnerID() int64 {
	return jr.rid
}

func (jr *JobRunner) DecodeParams(p string, v any) error {
	err := Decode(p, v)
	if err != nil {
		err = fmt.Errorf("Failed to decode parameters for job #%d : %w", jr.jid, err)
	}
	return err
}

func (jr *JobRunner) GetJob() (*Job, error) {
	return GetJob(jr.DB, jr.JobTable, jr.jid)
}

func (jr *JobRunner) Checkout() error {
	return CheckoutJob(jr.DB, jr.JobTable, jr.jid, jr.rid, jr.Log)
}

func (jr *JobRunner) Ping() error {
	if jr.pingAt.Add(jr.Timeout).After(time.Now()) {
		return nil
	}

	if err := PingJob(jr.DB, jr.JobTable, jr.jid, jr.rid, jr.Log); err != nil {
		return err
	}

	jr.pingAt = time.Now()
	return nil
}

func (jr *JobRunner) PingAborted() bool {
	err := jr.Ping()

	return err != nil
}

func (jr *JobRunner) Running(result any) error {
	if jr.pingAt.Add(jr.Timeout).After(time.Now()) {
		return nil
	}

	if err := UpdateJob(jr.DB, jr.JobTable, jr.jid, jr.rid, JobStatusRunning, Encode(result), jr.Log); err != nil {
		return err
	}

	jr.pingAt = time.Now()
	return nil
}

func (jr *JobRunner) Abort(reason string) error {
	err := AbortJob(jr.DB, jr.JobTable, jr.jid, reason, jr.Log)
	jr.Log.Flush()
	return err
}

func (jr *JobRunner) Complete(result any) error {
	err := CompleteJob(jr.DB, jr.JobTable, jr.jid, Encode(result), jr.Log)
	jr.Log.Flush()
	return err
}
