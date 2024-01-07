package xwj

import (
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/xwa/xwm"
	"gorm.io/gorm"
)

type JobRunner struct {
	Log *log.Log
	jlw *JobLogWriter

	db  *gorm.DB
	jid int64
	rid int64

	job     *xwm.Job
	takeAt  time.Time
	pingAt  time.Time
	timeout time.Duration
}

func NewJobRunner(db *gorm.DB, jid, rid int64, logger ...log.Logger) *JobRunner {
	jr := &JobRunner{Log: log.NewLog(), db: db, jid: jid, rid: rid, timeout: time.Second}

	jr.jlw = NewJobLogWriter(db, jid)
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

func (jr *JobRunner) SetRunnerID(rid int64) {
	jr.rid = rid
}

func (jr *JobRunner) SetTimeout(timeout time.Duration) {
	jr.timeout = timeout
}

func (jr *JobRunner) expired() {
	jr.takeAt = time.Time{}
}

func (jr *JobRunner) take() (*xwm.Job, error) {
	if jr.takeAt.Add(jr.timeout).After(time.Now()) {
		return jr.job, nil
	}

	job := &xwm.Job{}
	r := jr.db.Select("id", "rid", "name", "status").Take(job, jr.jid)
	if r.Error != nil {
		log.Errorf("Failed to fetch job #%d: %v", jr.jid, r.Error)
	} else {
		jr.takeAt = time.Now()
		jr.job = job
	}
	return jr.job, r.Error
}

func (jr *JobRunner) IsAborted() bool {
	job, err := jr.take()
	if err != nil || job == nil {
		return true
	}

	return job.IsAborted()
}

func (jr *JobRunner) DecodeParams(p string, v any) error {
	err := Decode(p, v)
	if err != nil {
		log.Errorf("Failed to decode parameters for job #%d : %v", jr.jid, err)
	}
	return err
}

func (jr *JobRunner) Checkout() error {
	log.Infof("Checkout job #%d", jr.jid)

	job := &xwm.Job{RID: jr.rid, Status: xwm.JobStatusRunning, Error: ""}
	r := jr.db.Select("rid", "status", "error").Where("id = ? AND status <> ?", jr.jid, xwm.JobStatusRunning).Updates(job)
	if r.Error != nil {
		log.Errorf("Failed to checkout job #%d: %v", jr.jid, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		log.Warnf("Unable to checkout job #%d: %v", jr.jid, ErrJobCheckout)
		return ErrJobCheckout
	}

	jr.expired()
	return nil
}

func (jr *JobRunner) Running(result any) error {
	if jr.pingAt.Add(jr.timeout).After(time.Now()) {
		return nil
	}

	if err := jr.update(xwm.JobStatusRunning, result); err != nil {
		return err
	}

	jr.pingAt = time.Now()
	return nil
}

func (jr *JobRunner) Abort(reason string) error {
	log.Infof("Abort job #%d: %s", jr.jid, reason)

	jr.expired()

	job := &xwm.Job{Status: xwm.JobStatusAborted, Error: reason}
	r := jr.db.Where("id = ?", jr.jid).Select("status", "error").Updates(job)
	if r.Error != nil {
		log.Errorf("Failed to abort job #%d: %v", jr.jid, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		log.Warnf("Unable to abort job #%d: %v", jr.jid, ErrJobMissing)
		return ErrJobMissing
	}

	jr.Log.Flush()
	return nil
}

func (jr *JobRunner) Complete(result any) error {
	log.Infof("Complete job #%d: %v", jr.jid, result)

	err := jr.update(xwm.JobStatusCompleted, result)

	jr.Log.Flush()
	return err
}

func (jr *JobRunner) update(status string, result any) error {
	jr.expired()

	job := &xwm.Job{Status: status, Result: Encode(result)}
	r := jr.db.Where("id = ? AND rid = ?", jr.jid, jr.rid).Select("status", "result").Updates(job)
	if r.Error != nil {
		log.Errorf("Failed to update job #%d (%s): %v", jr.jid, status, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		log.Warnf("Unable to update job #%d (%s): %v", jr.jid, status, ErrJobMissing)
		return ErrJobMissing
	}

	return nil
}
