package sqlxjm

import (
	"errors"
	"time"

	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/xjm"
)

type sjm struct {
	db sqlx.Sqlx
	jt string // job table
	lt string // log table
}

func JM(db sqlx.Sqlx, jobTable, logTable string) xjm.JobManager {
	return &sjm{
		db: db,
		jt: jobTable,
		lt: logTable,
	}
}

func (sjm *sjm) CountJobLogs(jid int64, levels ...string) (cnt int64, err error) {
	sqb := sjm.db.Builder()

	sqb.Count().From(sjm.lt).Where("jid = ?", jid)
	if len(levels) > 0 {
		sqb.In("level", levels)
	}

	sql, args := sqb.Build()

	err = sjm.db.Get(&cnt, sql, args...)
	return
}

func (sjm *sjm) GetJobLogs(jid int64, minLid, maxLid int64, asc bool, limit int, levels ...string) (jls []*xjm.JobLog, err error) {
	sqb := sjm.db.Builder()

	sqb.Select().From(sjm.lt).Where("jid = ?", jid)
	if len(levels) > 0 {
		sqb.In("level", levels)
	}
	if minLid > 0 {
		sqb.Where("id >= ?", minLid)
	}
	if maxLid > 0 {
		sqb.Where("id <= ?", maxLid)
	}
	sqb.Limit(limit)
	sqb.Order("id", !asc)

	sql, args := sqb.Build()

	err = sjm.db.Select(&jls, sql, args...)
	return
}

func (sjm *sjm) AddJobLogs(jls []*xjm.JobLog) error {
	s := "INSERT INTO " + sjm.lt + " (jid, time, level, message) VALUES (:jid, :time, :level, :message)"
	_, err := sjm.db.NamedExec(s, jls)
	return err
}

func (sjm *sjm) AddJobLog(jid int64, time time.Time, level string, message string) error {
	sqb := sjm.db.Builder()

	sqb.Insert(sjm.lt)
	sqb.Setc("jid", jid)
	sqb.Setc("time", time)
	sqb.Setc("level", level)
	sqb.Setc("message", message)

	sql, args := sqb.Build()

	_, err := sjm.db.Exec(sql, args...)
	return err
}

func (sjm *sjm) GetJob(jid int64) (*xjm.Job, error) {
	sqb := sjm.db.Builder()
	sqb.Select().From(sjm.jt).Where("id = ?", jid)
	sql, args := sqb.Build()

	job := &xjm.Job{}
	err := sjm.db.Get(job, sql, args...)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (sjm *sjm) FindJob(name string, asc bool, status ...string) (job *xjm.Job, err error) {
	sqb := sjm.db.Builder()

	sqb.Select().From(sjm.jt)
	if name != "" {
		sqb.Where("name = ?", name)
	}
	if len(status) > 0 {
		sqb.In("status", status)
	}
	sqb.Order("id", !asc)
	sqb.Limit(1)

	sql, args := sqb.Build()

	job = &xjm.Job{}
	err = sjm.db.Get(job, sql, args...)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil, nil
	}

	return job, err
}

func (sjm *sjm) findJobs(name string, start, limit int, asc bool, status ...string) *sqlx.Builder {
	sqb := sjm.db.Builder()

	sqb.Select().From(sjm.jt)
	if name != "" {
		sqb.Where("name = ?", name)
	}
	if len(status) > 0 {
		sqb.In("status", status)
	}
	sqb.Order("id", !asc)
	sqb.Offset(start).Limit(limit)

	return sqb
}

func (sjm *sjm) FindJobs(name string, start, limit int, asc bool, status ...string) (jobs []*xjm.Job, err error) {
	sqb := sjm.findJobs(name, start, limit, asc, status...)
	sql, args := sqb.Build()

	err = sjm.db.Select(&jobs, sql, args...)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil, nil
	}
	return
}

func (sjm *sjm) IterJobs(it func(*xjm.Job) error, name string, start, limit int, asc bool, status ...string) error {
	sqb := sjm.findJobs(name, start, limit, asc, status...)
	sql, args := sqb.Build()

	rows, err := sjm.db.Queryx(sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		job := &xjm.Job{}

		if err := rows.StructScan(job); err != nil {
			return err
		}

		if err := it(job); err != nil {
			return err
		}
	}
	return nil
}

func (sjm *sjm) AppendJob(name, file, param string) (int64, error) {
	now := time.Now()

	sqb := sjm.db.Builder()
	sqb.Insert(sjm.jt)
	sqb.Setc("rid", 0)
	sqb.Setc("name", name)
	sqb.Setc("status", xjm.JobStatusPending)
	sqb.Setc("file", file)
	sqb.Setc("param", param)
	sqb.Setc("state", "")
	sqb.Setc("result", "")
	sqb.Setc("error", "")
	sqb.Setc("created_at", now)
	sqb.Setc("updated_at", now)

	if !sjm.db.SupportLastInsertID() {
		sqb.Returns("id")
	}

	sql, args := sqb.Build()
	return sjm.db.Create(sql, args...)
}

func (sjm *sjm) AbortJob(jid int64, reason string) error {
	sqb := sjm.db.Builder()

	sqb.Update(sjm.jt)
	sqb.Setc("status", xjm.JobStatusAborted)
	sqb.Setc("error", reason)
	sqb.Setc("updated_at", time.Now())
	sqb.Where("id = ?", jid)
	sqb.In("status", xjm.JobPendingRunning)

	sql, args := sqb.Build()

	r, err := sjm.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (sjm *sjm) CompleteJob(jid int64) error {
	sqb := sjm.db.Builder()

	sqb.Update(sjm.jt)
	sqb.Setc("status", xjm.JobStatusCompleted)
	sqb.Setc("error", "")
	sqb.Setc("updated_at", time.Now())
	sqb.Where("id = ?", jid)

	sql, args := sqb.Build()

	r, err := sjm.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (sjm *sjm) CheckoutJob(jid, rid int64) error {
	sqb := sjm.db.Builder()

	sqb.Update(sjm.jt)
	sqb.Setc("rid", rid)
	sqb.Setc("status", xjm.JobStatusRunning)
	sqb.Setc("error", "")
	sqb.Setc("updated_at", time.Now())
	sqb.Where("id = ?", jid)
	sqb.Where("status = ?", xjm.JobStatusPending)

	sql, args := sqb.Build()

	r, err := sjm.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return xjm.ErrJobCheckout
	}
	return nil
}

func (sjm *sjm) PingJob(jid, rid int64) error {
	sqb := sjm.db.Builder()

	sqb.Update(sjm.jt)
	sqb.Setc("updated_at", time.Now())
	sqb.Where("id = ?", jid)
	sqb.Where("rid = ?", rid)
	sqb.Where("status = ?", xjm.JobStatusRunning)

	sql, args := sqb.Build()

	r, err := sjm.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return xjm.ErrJobAborted
	}
	return nil
}

func (sjm *sjm) SetJobState(jid, rid int64, state string) error {
	sqb := sjm.db.Builder()

	sqb.Update(sjm.jt)
	sqb.Setc("state", state)
	sqb.Setc("updated_at", time.Now())
	sqb.Where("id = ?", jid)
	sqb.Where("rid = ?", rid)

	sql, args := sqb.Build()

	r, err := sjm.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (sjm *sjm) AddJobResult(jid, rid int64, result string) error {
	sqb := sjm.db.Builder()

	sqb.Update(sjm.jt)
	sqb.Setx("result", "result || ?", result)
	sqb.Setc("updated_at", time.Now())
	sqb.Where("id = ?", jid)
	sqb.Where("rid = ?", rid)

	sql, args := sqb.Build()

	r, err := sjm.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (sjm *sjm) ReappendJobs(before time.Time) (int64, error) {
	sqb := sjm.db.Builder()

	sqb.Update(sjm.jt)
	sqb.Setc("rid", 0)
	sqb.Setc("state", xjm.JobStatusPending)
	sqb.Setc("error", "")
	sqb.Setc("updated_at", time.Now())
	sqb.Where("status = ?", xjm.JobStatusRunning)
	sqb.Where("updated_at < ?", before)

	sql, args := sqb.Build()

	r, err := sjm.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

func (sjm *sjm) StartJobs(limit int, start func(*xjm.Job)) error {
	sqb := sjm.db.Builder()

	sqb.Select().From(sjm.jt)
	sqb.Where("status = ?", xjm.JobStatusPending)
	sqb.Order("id", false)
	sqb.Limit(limit)

	sql, args := sqb.Build()

	var jobs []*xjm.Job
	err := sjm.db.Select(&jobs, sql, args...)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil
	}

	if err != nil {
		return err
	}

	for _, job := range jobs {
		start(job)
	}

	return nil
}

func (sjm *sjm) DeleteJobs(jids ...int64) (jobs int64, logs int64, err error) {
	if len(jids) == 0 {
		return
	}

	sqa := sjm.db.Builder()
	sqa.Delete(sjm.lt)
	sqa.In("jid", jids)

	sql, args := sqa.Build()

	var r sqlx.Result
	if r, err = sjm.db.Exec(sql, args...); err != nil {
		return
	}
	if logs, err = r.RowsAffected(); err != nil {
		return
	}

	sqb := sjm.db.Builder()
	sqb.Delete(sjm.jt)
	sqb.In("id", jids)

	sql, args = sqb.Build()
	if r, err = sjm.db.Exec(sql, args...); err != nil {
		return
	}

	jobs, err = r.RowsAffected()
	return
}

func (sjm *sjm) CleanOutdatedJobs(before time.Time) (jobs int64, logs int64, err error) {
	sqb := sjm.db.Builder()
	sqb.Select("id").From(sjm.jt)
	sqb.Where("updated_at < ?", before)
	sqb.In("status", xjm.JobAbortedCompleted)

	sqa := sjm.db.Builder()
	sqa.Delete(sjm.lt)
	sqa.Where("jid IN ("+sqb.SQL()+")", sqb.Params()...)

	sql, args := sqa.Build()

	var r sqlx.Result
	if r, err = sjm.db.Exec(sql, args...); err != nil {
		return
	}
	if logs, err = r.RowsAffected(); err != nil {
		return
	}

	sqb.Delete(sjm.jt)

	sql, args = sqb.Build()
	if r, err = sjm.db.Exec(sql, args...); err != nil {
		return
	}

	jobs, err = r.RowsAffected()
	return
}
