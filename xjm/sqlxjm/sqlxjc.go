package sqlxjm

import (
	"errors"
	"time"

	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xjm"
)

type sjc struct {
	db sqlx.Sqlx
	tb string // jc chain table
}

func JC(db sqlx.Sqlx, table string) xjm.JobChainer {
	return &sjc{
		db: db,
		tb: table,
	}
}

func (sjc *sjc) GetJobChain(cid int64) (*xjm.JobChain, error) {
	s := sjc.db.Rebind("SELECT * FROM " + sjc.tb + " WHERE id = ?")

	jc := &xjm.JobChain{}
	err := sjc.db.Get(jc, s, cid)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return jc, nil
}

func (sjc *sjc) FindJobChain(name string, asc bool, status ...string) (jc *xjm.JobChain, err error) {
	sqb := sqx.Builder{}

	sqb.Select("*").From(sjc.tb)
	if name != "" {
		sqb.Where("name = ?", name)
	}
	if len(status) > 0 {
		sqb.In("status", status)
	}
	sqb.Order("id " + str.If(asc, "ASC", "DESC"))
	sqb.Limit(1)

	sql, args := sqb.Build()
	sql = sjc.db.Rebind(sql)

	jc = &xjm.JobChain{}
	err = sjc.db.Get(jc, sql, args...)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil, nil
	}

	return jc, err
}

func (sjc *sjc) findJobChains(name string, start, limit int, asc bool, status ...string) *sqx.Builder {
	sqb := &sqx.Builder{}

	sqb.Select("*").From(sjc.tb)
	if name != "" {
		sqb.Where("name = ?", name)
	}
	if len(status) > 0 {
		sqb.In("status", status)
	}
	sqb.Order("id " + str.If(asc, "ASC", "DESC"))
	sqb.Offset(start).Limit(limit)

	return sqb
}

func (sjc *sjc) FindJobChains(name string, start, limit int, asc bool, status ...string) (jcs []*xjm.JobChain, err error) {
	sqb := sjc.findJobChains(name, start, limit, asc, status...)
	sql, args := sqb.Build()
	sql = sjc.db.Rebind(sql)

	err = sjc.db.Select(&jcs, sql, args...)
	if errors.Is(err, sqlx.ErrNoRows) {
		return nil, nil
	}
	return
}

func (sjc *sjc) IterJobChains(it func(*xjm.JobChain) error, name string, start, limit int, asc bool, status ...string) error {
	sqb := sjc.findJobChains(name, start, limit, asc, status...)
	sql, args := sqb.Build()
	sql = sjc.db.Rebind(sql)

	rows, err := sjc.db.Queryx(sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		jc := &xjm.JobChain{}

		if err := rows.StructScan(jc); err != nil {
			return err
		}

		if err := it(jc); err != nil {
			return err
		}
	}
	return nil
}

func (sjc *sjc) CreateJobChain(name, states string) (int64, error) {
	jc := &xjm.JobChain{Name: name, States: states, Status: xjm.JobStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	sqb := sqx.Builder{}

	sqb.Insert(sjc.tb)
	sqb.Columns("name", "status", "states", "created_at", "updated_at")
	sqb.Values(":name", ":status", ":states", ":created_at", ":updated_at")

	sql := sqb.SQL()
	if sjc.db.SupportLastInsertID() {
		r, err := sjc.db.NamedExec(sql, jc)
		if err != nil {
			return 0, err
		}
		return r.LastInsertId()
	}

	sql += " RETURNING id"
	err := sjc.db.NamedQueryRow(sql, jc).Scan(&jc.ID)
	return jc.ID, err
}

func (sjc *sjc) UpdateJobChain(cid int64, status string, states ...string) error {
	if status == "" && len(states) == 0 {
		return nil
	}

	sqb := sqx.Builder{}

	sqb.Update(sjc.tb)
	if status != "" {
		sqb.Set("status = ?", status)
	}
	if len(states) > 0 {
		sqb.Set("states = ?", states[0])
	}
	sqb.Set("updated_at = ?", time.Now())
	sqb.Where("id = ?", cid)

	sql, args := sqb.Build()
	sql = sjc.db.Rebind(sql)

	r, err := sjc.db.Exec(sql, args...)
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

func (sjc *sjc) CleanOutdatedJobChains(before time.Time) (cnt int64, err error) {
	sqb := sqx.Builder{}
	sqb.Delete(sjc.tb)
	sqb.Where("updated_at < ?", before)
	sqb.In("status", xjm.JobChainAbortedCompleted)

	sql := sjc.db.Rebind(sqb.SQL())

	var r sqlx.Result
	if r, err = sjc.db.Exec(sql, sqb.Params()...); err != nil {
		return
	}

	cnt, err = r.RowsAffected()
	return
}
