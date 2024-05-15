package sqlxjm

import (
	"errors"
	"time"

	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
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
	sqb := &sqx.Builder{}

	sqb.Select("COUNT(1)").From(sjm.lt).Where("jid = ?", jid)
	if len(levels) > 0 {
		sqb.In("level", levels)
	}

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

	err = sjm.db.Get(&cnt, sql, args...)
	return
}

func (sjm *sjm) GetJobLogs(jid int64, min, max int64, asc bool, limit int, levels ...string) (jls []*xjm.JobLog, err error) {
	sqb := &sqx.Builder{}

	sqb.Select("*").From(sjm.lt).Where("jid = ?", jid)
	if len(levels) > 0 {
		sqb.In("level", levels)
	}
	if min > 0 {
		sqb.Where("id >= ?", min)
	}
	if max > 0 {
		sqb.Where("id <= ?", max)
	}
	sqb.Limit(limit)
	sqb.Order("id " + str.If(asc, "ASC", "DESC"))

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

	err = sjm.db.Select(&jls, sql, args...)
	return
}

func (sjm *sjm) AddJobLogs(jls []*xjm.JobLog) error {
	s := "INSERT INTO " + sjm.lt + " (jid, time, level, message) VALUES (:jid, :time, :level, :message)"
	_, err := sjm.db.NamedExec(s, jls)
	return err
}

func (sjm *sjm) GetJob(jid int64) (*xjm.Job, error) {
	s := sjm.db.Rebind("SELECT * FROM " + sjm.jt + " WHERE id = ?")

	job := &xjm.Job{}
	err := sjm.db.Get(job, s, jid)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (sjm *sjm) FindJob(name string, asc bool, status ...string) (job *xjm.Job, err error) {
	sqb := &sqx.Builder{}

	sqb.Select("*").From(sjm.jt)
	sqb.Where("name = ?", name)
	if len(status) > 0 {
		sqb.In("status", status)
	}
	sqb.Order("id " + str.If(asc, "ASC", "DESC"))
	sqb.Limit(1)

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

	job = &xjm.Job{}
	err = sjm.db.Get(job, sql, args...)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil, nil
	}

	return job, err
}

func (sjm *sjm) findJobs(name string, start, limit int, asc bool, status ...string) *sqx.Builder {
	sqb := &sqx.Builder{}

	sqb.Select("*").From(sjm.jt).Where("name = ?", name)
	if len(status) > 0 {
		sqb.In("status", status)
	}
	sqb.Order("id " + str.If(asc, "ASC", "DESC"))
	sqb.Offset(start).Limit(limit)

	return sqb
}

func (sjm *sjm) FindJobs(name string, start, limit int, asc bool, status ...string) (jobs []*xjm.Job, err error) {
	sqb := sjm.findJobs(name, start, limit, asc, status...)
	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

	err = sjm.db.Select(&jobs, sql, args...)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil, nil
	}
	return
}

func (sjm *sjm) IterJobs(it func(*xjm.Job) error, name string, start, limit int, asc bool, status ...string) error {
	sqb := sjm.findJobs(name, start, limit, asc, status...)
	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

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
	job := &xjm.Job{Name: name, File: file, Param: param, Status: xjm.JobStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	sqb := &sqx.Builder{}

	sqb.Insert(sjm.jt)
	sqb.Columns("rid", "name", "status", "file", "param", "state", "result", "error", "created_at", "updated_at")
	sqb.Values(":rid", ":name", ":status", ":file", ":param", ":state", ":result", ":error", ":created_at", ":updated_at")

	sql := sqb.SQL()
	if sjm.db.SupportLastInsertID() {
		r, err := sjm.db.NamedExec(sql, job)
		if err != nil {
			return 0, err
		}
		return r.LastInsertId()
	}

	sql += " RETURNING id"
	err := sjm.db.NamedQueryRow(sql, job).Scan(&job.ID)
	return job.ID, err
}

func (sjm *sjm) AbortJob(jid int64, reason string) error {
	sqb := &sqx.Builder{}

	sqb.Update(sjm.jt)
	sqb.Set("status = ?", xjm.JobStatusAborted)
	sqb.Set("error = ?", reason)
	sqb.Set("updated_at = ?", time.Now())
	sqb.Where("id = ?", jid)
	sqb.In("status", xjm.JobPendingRunning)

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

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
	sqb := &sqx.Builder{}

	sqb.Update(sjm.jt)
	sqb.Set("status = ?", xjm.JobStatusCompleted)
	sqb.Set("error = ?", "")
	sqb.Set("updated_at = ?", time.Now())
	sqb.Where("id = ?", jid)

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

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
	sqb := &sqx.Builder{}

	sqb.Update(sjm.jt)
	sqb.Set("rid = ?", rid)
	sqb.Set("status = ?", xjm.JobStatusRunning)
	sqb.Set("error = ?", "")
	sqb.Set("updated_at = ?", time.Now())
	sqb.Where("id = ?", jid)
	sqb.Where("status = ?", xjm.JobStatusPending)

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

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
	sqb := &sqx.Builder{}

	sqb.Update(sjm.jt)
	sqb.Set("updated_at = ?", time.Now())
	sqb.Where("id = ?", jid)
	sqb.Where("rid = ?", rid)
	sqb.Where("status = ?", xjm.JobStatusRunning)

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

	r, err := sjm.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return xjm.ErrJobPing
	}
	return nil
}

func (sjm *sjm) RunningJob(jid, rid int64, state string) error {
	sqb := &sqx.Builder{}

	sqb.Update(sjm.jt)
	sqb.Set("state = ?", state)
	sqb.Set("updated_at = ?", time.Now())
	sqb.Where("id = ?", jid)
	sqb.Where("rid = ?", rid)

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

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
	sqb := &sqx.Builder{}

	sqb.Update(sjm.jt)
	sqb.Set("result = result || ?", result)
	sqb.Set("updated_at = ?", time.Now())
	sqb.Where("id = ?", jid)
	sqb.Where("rid = ?", rid)

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

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
	sqb := &sqx.Builder{}

	sqb.Update(sjm.jt)
	sqb.Set("rid = ?", 0)
	sqb.Set("state = ?", xjm.JobStatusPending)
	sqb.Set("error = ?", "")
	sqb.Set("updated_at = ?", time.Now())
	sqb.Where("status = ?", xjm.JobStatusRunning)
	sqb.Where("updated_at < ?", before)

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

	r, err := sjm.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

func (sjm *sjm) StartJobs(limit int, run func(*xjm.Job)) error {
	sqb := &sqx.Builder{}

	sqb.Select("*")
	sqb.From(sjm.jt)
	sqb.Where("status = ?", xjm.JobStatusPending)
	sqb.Order("id ASC")
	sqb.Limit(limit)

	sql, args := sqb.Build()
	sql = sjm.db.Rebind(sql)

	var jobs []*xjm.Job
	err := sjm.db.Select(&jobs, sql, args...)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil
	}

	if err != nil {
		return err
	}

	for _, job := range jobs {
		go run(job)
	}

	return nil
}

func (sjm *sjm) CleanOutdatedJobs(before time.Time) (jobs int64, logs int64, err error) {
	sqb := &sqx.Builder{}
	sqb.Select("id").From(sjm.jt)
	sqb.Where("updated_at < ?", before)
	sqb.In("status", xjm.JobAbortedCompleted)

	sqa := &sqx.Builder{}
	sqa.Delete(sjm.lt)
	sqa.Where("jid IN ("+sqb.SQL()+")", sqb.Params()...)

	sql := sjm.db.Rebind(sqa.SQL())

	var r sqlx.Result
	if r, err = sjm.db.Exec(sql, sqa.Params()...); err != nil {
		return
	}
	if logs, err = r.RowsAffected(); err != nil {
		return
	}

	sqb.Delete(sjm.jt)

	sql = sjm.db.Rebind(sqb.SQL())
	if r, err = sjm.db.Exec(sql, sqb.Params()...); err != nil {
		return
	}
	jobs, err = r.RowsAffected()
	return
}
