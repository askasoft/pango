package xwj

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

func (jr *JobRunner) GetJob(details ...bool) (*Job, error) {
	return GetJob(jr.DB, jr.JobTable, jr.jid, details...)
}

func (jr *JobRunner) Checkout() error {
	jr.Log.Debugf("Checkout job #%d", jr.jid)

	job := &Job{RID: jr.rid, Status: JobStatusRunning, Error: ""}
	r := jr.DB.Table(jr.JobTable).Select("rid", "status", "error").Where("id = ? AND status <> ?", jr.jid, JobStatusRunning).Updates(job)
	if r.Error != nil {
		jr.Log.Errorf("Failed to checkout job #%d: %v", jr.jid, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		jr.Log.Errorf("Unable to checkout job #%d: %v", jr.jid, ErrJobCheckout)
		return ErrJobCheckout
	}

	return nil
}

func (jr *JobRunner) Ping() error {
	if jr.pingAt.Add(jr.Timeout).After(time.Now()) {
		return nil
	}

	tx := jr.DB.Table(jr.JobTable)
	r := tx.Where("id = ? AND rid = ? AND status = ?", jr.jid, jr.rid, JobStatusRunning).Update("updated_at", time.Now())
	if r.Error != nil {
		jr.Log.Errorf("Failed to ping job #%d: %v", jr.jid, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		return ErrJobPing
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

	if err := jr.update(JobStatusRunning, result); err != nil {
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
	jr.Log.Debugf("Complete job #%d: %v", jr.jid, result)

	err := jr.update(JobStatusCompleted, result)

	jr.Log.Flush()
	return err
}

func (jr *JobRunner) update(status string, result any) error {
	job := &Job{Status: status, Result: Encode(result)}

	r := jr.DB.Table(jr.JobTable).Where("id = ? AND rid = ?", jr.jid, jr.rid).Select("status", "result").Updates(job)
	if r.Error != nil {
		jr.Log.Errorf("Failed to update job #%d (%s): %v", jr.jid, status, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		jr.Log.Errorf("Unable to update job #%d (%s): %d, %v", jr.jid, status, r.RowsAffected, ErrJobMissing)
		return ErrJobMissing
	}

	return nil
}
