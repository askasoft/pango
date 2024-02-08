package xwj

import (
	"time"

	"github.com/askasoft/pango/log"
	"gorm.io/gorm"
)

type JobRunner struct {
	Log *log.Log
	jlw *JobLogWriter

	DB  *gorm.DB
	jid int64
	rid int64

	pingAt  time.Time
	timeout time.Duration
}

func NewJobRunner(db *gorm.DB, jid, rid int64, logger ...log.Logger) *JobRunner {
	jr := &JobRunner{Log: log.NewLog(), DB: db, jid: jid, rid: rid, timeout: time.Second}

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

func (jr *JobRunner) DecodeParams(p string, v any) error {
	err := Decode(p, v)
	if err != nil {
		log.Errorf("Failed to decode parameters for job #%d : %v", jr.jid, err)
	}
	return err
}

func (jr *JobRunner) GetJob(details ...bool) (*Job, error) {
	return GetJob(jr.DB, jr.jid, details...)
}

func (jr *JobRunner) Checkout() error {
	log.Infof("Checkout job #%d", jr.jid)

	job := &Job{RID: jr.rid, Status: JobStatusRunning, Error: ""}
	r := jr.DB.Select("rid", "status", "error").Where("id = ? AND status <> ?", jr.jid, JobStatusRunning).Updates(job)
	if r.Error != nil {
		log.Errorf("Failed to checkout job #%d: %v", jr.jid, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		log.Warnf("Unable to checkout job #%d: %v", jr.jid, ErrJobCheckout)
		return ErrJobCheckout
	}

	return nil
}

func (jr *JobRunner) Ping() error {
	if jr.pingAt.Add(jr.timeout).After(time.Now()) {
		return nil
	}

	tx := jr.DB.Model(&Job{})
	r := tx.Where("id = ? AND rid = ? AND status = ?", jr.jid, jr.rid, JobStatusRunning).Update("updated_at", time.Now())
	if r.Error != nil {
		log.Errorf("Failed to ping job #%d: %v", jr.jid, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		log.Warnf("Unable to ping job #%d: %d", jr.jid, r.RowsAffected)
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
	if jr.pingAt.Add(jr.timeout).After(time.Now()) {
		return nil
	}

	if err := jr.update(JobStatusRunning, result); err != nil {
		return err
	}

	jr.pingAt = time.Now()
	return nil
}

func (jr *JobRunner) Abort(reason string) error {
	log.Infof("Abort job #%d: %s", jr.jid, reason)

	job := &Job{Status: JobStatusAborted, Error: reason}
	r := jr.DB.Where("id = ?", jr.jid).Select("status", "error").Updates(job)
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

	err := jr.update(JobStatusCompleted, result)

	jr.Log.Flush()
	return err
}

func (jr *JobRunner) update(status string, result any) error {
	job := &Job{Status: status, Result: Encode(result)}
	r := jr.DB.Where("id = ? AND rid = ?", jr.jid, jr.rid).Select("status", "result").Updates(job)
	if r.Error != nil {
		log.Errorf("Failed to update job #%d (%s): %v", jr.jid, status, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		log.Errorf("Unable to update job #%d (%s): %d, %v", jr.jid, status, r.RowsAffected, ErrJobMissing)
		return ErrJobMissing
	}

	return nil
}
