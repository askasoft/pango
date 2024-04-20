package gormjm

import (
	"errors"
	"time"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xjm"
	"gorm.io/gorm"
)

type gjm struct {
	db *gorm.DB
	jt string // job table
	lt string // log table
}

func JM(db *gorm.DB, jobTable, logTable string) xjm.JobManager {
	return &gjm{
		db: db,
		jt: jobTable,
		lt: logTable,
	}
}

// CountJobLogs count job logs
func (gjm *gjm) CountJobLogs(jid int64, levels ...string) (int64, error) {
	tx := gjm.db.Table(gjm.lt).Where("jid = ?", jid)
	if len(levels) > 0 {
		tx.Where("level IN ?", levels)
	}

	var cnt int64
	r := tx.Count(&cnt)
	return cnt, r.Error
}

// GetJobLogs get job logs
// set levels to ("I", "W", "E", "F") to filter DEBUG/TRACE logs
func (gjm *gjm) GetJobLogs(jid int64, min, max int64, asc bool, limit int, levels ...string) ([]*xjm.JobLog, error) {
	var jls []*xjm.JobLog

	tx := gjm.db.Table(gjm.lt).Where("jid = ?", jid)
	if len(levels) > 0 {
		tx.Where("level IN ?", levels)
	}
	if min > 0 {
		tx = tx.Where("id >= ?", min)
	}
	if max > 0 {
		tx = tx.Where("id <= ?", max)
	}
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	tx = tx.Order("id " + str.If(asc, "ASC", "DESC"))

	r := tx.Find(&jls)
	return jls, r.Error
}

func (gjm *gjm) AddJobLogs(jls []*xjm.JobLog) error {
	r := gjm.db.Table(gjm.lt).Create(jls)
	return r.Error
}

func (gjm *gjm) GetJob(jid int64) (*xjm.Job, error) {
	job := &xjm.Job{}
	r := gjm.db.Table(gjm.jt).Where("id = ?", jid).First(job)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return job, nil
}

// FindJob find the latest job by name, default select all columns.
// cols: columns to select.
func (gjm *gjm) FindJob(name string, cols ...string) (*xjm.Job, error) {
	tx := gjm.db.Table(gjm.jt).Where("name = ?", name).Order("id DESC")
	if len(cols) > 0 {
		tx = tx.Select(cols)
	}

	job := &xjm.Job{}
	r := tx.First(job)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return job, r.Error
}

// FindJobs find jobs by name, default select all columns.
// cols: columns to select.
func (gjm *gjm) FindJobs(name string, start, limit int, cols ...string) ([]*xjm.Job, error) {
	jobs := []*xjm.Job{}

	tx := gjm.db.Table(gjm.jt).Where("name = ?", name).Order("id DESC")
	if len(cols) > 0 {
		tx = tx.Select(cols)
	}
	if start > 0 {
		tx = tx.Offset(start)
	}
	if limit > 0 {
		tx = tx.Limit(limit)
	}

	r := tx.Find(&jobs)
	return jobs, r.Error
}

func (gjm *gjm) AppendJob(name, file, param string) (int64, error) {
	job := &xjm.Job{Name: name, File: file, Param: param, Status: xjm.JobStatusPending}
	r := gjm.db.Table(gjm.jt).Create(job)
	return job.ID, r.Error
}

func (gjm *gjm) FindAndAbortJob(name, reason string) error {
	job, err := gjm.FindJob(name)
	if err != nil {
		return err
	}

	return gjm.AbortJob(job.ID, reason)
}

func (gjm *gjm) AbortJob(jid int64, reason string) error {
	job := &xjm.Job{Status: xjm.JobStatusAborted, Error: reason}
	jss := xjm.JobPendingRunning

	tx := gjm.db.Table(gjm.jt).Where("id = ? AND status IN ?", jid, jss)
	r := tx.Select("status", "error").Updates(job)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (gjm *gjm) CompleteJob(jid int64, result string) error {
	job := &xjm.Job{Status: xjm.JobStatusCompleted, Result: result}

	tx := gjm.db.Table(gjm.jt).Where("id = ?", jid)
	r := tx.Select("status", "result", "error").Updates(job)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (gjm *gjm) CheckoutJob(jid, rid int64) error {
	job := &xjm.Job{RID: rid, Status: xjm.JobStatusRunning, Error: ""}

	r := gjm.db.Table(gjm.jt).Select("rid", "status", "error").Where("id = ? AND status <> ?", jid, xjm.JobStatusRunning).Updates(job)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return xjm.ErrJobCheckout
	}
	return nil
}

func (gjm *gjm) PingJob(jid, rid int64) error {
	tx := gjm.db.Table(gjm.jt)
	r := tx.Where("id = ? AND rid = ? AND status = ?", jid, rid, xjm.JobStatusRunning).Update("updated_at", time.Now())
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return xjm.ErrJobPing
	}
	return nil
}

func (gjm *gjm) RunningJob(jid, rid int64, state string) error {
	tx := gjm.db.Table(gjm.jt).Where("id = ? AND rid = ?", jid, rid)
	r := tx.Update("state", state)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (gjm *gjm) ReappendJobs(before time.Time) (int64, error) {
	job := &xjm.Job{RID: 0, Status: xjm.JobStatusPending, Error: ""}

	tx := gjm.db.Table(gjm.jt).Where("status = ? AND updated_at < ?", xjm.JobStatusRunning, before)
	r := tx.Select("rid", "status", "error").Updates(job)
	return r.RowsAffected, r.Error
}

func (gjm *gjm) StartJobs(limit int, run func(*xjm.Job)) error {
	var jobs []*xjm.Job

	r := gjm.db.Table(gjm.jt).Where("status = ?", xjm.JobStatusPending).Order("id asc").Limit(limit).Find(&jobs)
	if r.Error != nil {
		return r.Error
	}

	for _, job := range jobs {
		go run(job)
	}

	return nil
}

func (gjm *gjm) CleanOutdatedJobs(before time.Time) (jobs int64, logs int64, err error) {
	jss := xjm.JobAbortedCompleted
	where := "jid IN (SELECT id FROM " + gjm.jt + " WHERE status IN ? AND updated_at < ?)"

	r := gjm.db.Table(gjm.lt).Where(where, jss, before).Delete(&xjm.JobLog{})
	logs, err = r.RowsAffected, r.Error
	if err != nil {
		return
	}

	r = gjm.db.Table(gjm.jt).Where("status IN ? AND updated_at < ?", jss, before).Delete(&xjm.Job{})
	jobs, err = r.RowsAffected, r.Error
	return
}
