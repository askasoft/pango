package gormjm

import (
	"errors"
	"fmt"
	"time"

	"github.com/askasoft/pango/xjm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (gjm *gjm) CountJobLogs(jid int64, levels ...string) (int64, error) {
	tx := gjm.db.Table(gjm.lt).Where("jid = ?", jid)
	if len(levels) > 0 {
		tx.Where("level IN ?", levels)
	}

	var cnt int64
	r := tx.Count(&cnt)
	return cnt, r.Error
}

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
	tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: !asc})

	r := tx.Find(&jls)
	return jls, r.Error
}

func (gjm *gjm) AddJobLogs(jls []*xjm.JobLog) error {
	return gjm.db.Table(gjm.lt).Create(jls).Error
}

func (gjm *gjm) AddJobLog(jid int64, time time.Time, level string, message string) error {
	jlg := &xjm.JobLog{JID: jid, Time: time, Level: level, Message: message}
	return gjm.db.Table(gjm.lt).Create(jlg).Error
}

func (gjm *gjm) GetJob(jid int64) (*xjm.Job, error) {
	job := &xjm.Job{}
	r := gjm.db.Table(gjm.jt).Where("id = ?", jid).Take(job)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return job, nil
}

func (gjm *gjm) FindJob(name string, asc bool, status ...string) (*xjm.Job, error) {
	tx := gjm.db.Table(gjm.jt)
	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if len(status) > 0 {
		tx = tx.Where("status IN ?", status)
	}
	tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: !asc})

	job := &xjm.Job{}
	r := tx.Take(job)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return job, r.Error
}

func (gjm *gjm) findJobs(name string, start, limit int, asc bool, status ...string) *gorm.DB {
	tx := gjm.db.Table(gjm.jt)
	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if len(status) > 0 {
		tx = tx.Where("status IN ?", status)
	}
	tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: !asc})

	if start > 0 {
		tx = tx.Offset(start)
	}
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	return tx
}

func (gjm *gjm) FindJobs(name string, start, limit int, asc bool, status ...string) (jobs []*xjm.Job, err error) {
	tx := gjm.findJobs(name, start, limit, asc, status...)
	err = tx.Find(&jobs).Error
	return
}

func (gjm *gjm) IterJobs(it func(*xjm.Job) error, name string, start, limit int, asc bool, status ...string) error {
	tx := gjm.findJobs(name, start, limit, asc, status...)

	rows, err := tx.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		job := &xjm.Job{}

		if err := tx.ScanRows(rows, job); err != nil {
			return err
		}

		if err := it(job); err != nil {
			return err
		}
	}
	return nil
}

func (gjm *gjm) AppendJob(name, file, param string) (int64, error) {
	job := &xjm.Job{Name: name, File: file, Param: param, Status: xjm.JobStatusPending}
	r := gjm.db.Table(gjm.jt).Create(job)
	return job.ID, r.Error
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

func (gjm *gjm) CompleteJob(jid int64) error {
	job := &xjm.Job{ID: jid, Status: xjm.JobStatusCompleted}

	r := gjm.db.Table(gjm.jt).Select("status", "error").Updates(job)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (gjm *gjm) CheckoutJob(jid, rid int64) error {
	job := &xjm.Job{RID: rid, Status: xjm.JobStatusRunning}

	tx := gjm.db.Table(gjm.jt).Where("id = ? AND status <> ?", jid, xjm.JobStatusRunning)
	r := tx.Select("rid", "status", "error").Updates(job)
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

func (gjm *gjm) AddJobResult(jid, rid int64, result string) error {
	sql := fmt.Sprintf("UPDATE %s SET result = result || ?, updated_at = ? WHERE id = ? AND rid = ?", gjm.jt)
	r := gjm.db.Exec(sql, result, time.Now(), jid, rid)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (gjm *gjm) ReappendJobs(before time.Time) (int64, error) {
	job := &xjm.Job{Status: xjm.JobStatusPending}

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

func (gjm *gjm) DeleteJobs(jids ...int64) (jobs int64, logs int64, err error) {
	if len(jids) == 0 {
		return
	}

	r := gjm.db.Table(gjm.lt).Where("jid IN ?", jids).Delete(&xjm.JobLog{})
	logs, err = r.RowsAffected, r.Error
	if err != nil {
		return
	}

	r = gjm.db.Table(gjm.jt).Where("id IN ?", jids).Delete(&xjm.Job{})
	jobs, err = r.RowsAffected, r.Error
	return
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
