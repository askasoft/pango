package gormjm

import (
	"errors"
	"time"

	"github.com/askasoft/pango/xjm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gjc struct {
	db *gorm.DB
	tb string // jc chain table
}

func JC(db *gorm.DB, table string) xjm.JobChainer {
	return &gjc{
		db: db,
		tb: table,
	}
}

func (gjc *gjc) GetJobChain(cid int64) (*xjm.JobChain, error) {
	jc := &xjm.JobChain{}
	r := gjc.db.Table(gjc.tb).Where("id = ?", cid).Take(jc)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return jc, nil
}

func (gjc *gjc) FindJobChain(name string, asc bool, status ...string) (*xjm.JobChain, error) {
	tx := gjc.db.Table(gjc.tb)
	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if len(status) > 0 {
		tx = tx.Where("status IN ?", status)
	}
	tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: !asc})

	jc := &xjm.JobChain{}
	r := tx.Take(jc)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return jc, r.Error
}

func (gjc *gjc) findJobChains(name string, start, limit int, asc bool, status ...string) *gorm.DB {
	tx := gjc.db.Table(gjc.tb)
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

func (gjc *gjc) FindJobChains(name string, start, limit int, asc bool, status ...string) (jcs []*xjm.JobChain, err error) {
	tx := gjc.findJobChains(name, start, limit, asc, status...)
	err = tx.Find(&jcs).Error
	return
}

func (gjc *gjc) IterJobChains(it func(*xjm.JobChain) error, name string, start, limit int, asc bool, status ...string) error {
	tx := gjc.findJobChains(name, start, limit, asc, status...)

	rows, err := tx.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		jc := &xjm.JobChain{}

		if err := tx.ScanRows(rows, jc); err != nil {
			return err
		}

		if err := it(jc); err != nil {
			return err
		}
	}
	return nil
}

func (gjc *gjc) CreateJobChain(name, states string) (int64, error) {
	jc := &xjm.JobChain{Name: name, States: states, Status: xjm.JobChainPending}
	r := gjc.db.Table(gjc.tb).Create(jc)
	return jc.ID, r.Error
}

func (gjc *gjc) UpdateJobChain(cid int64, status string, states ...string) error {
	jc := &xjm.JobChain{ID: cid}

	cols := make([]string, 0, 2)
	if status != "" {
		jc.Status = status
		cols = append(cols, "status")
	}
	if len(states) > 0 {
		jc.States = states[0]
		cols = append(cols, "states")
	}

	if len(cols) == 0 {
		return nil
	}

	r := gjc.db.Table(gjc.tb).Select(cols).Updates(jc)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return xjm.ErrJobMissing
	}
	return nil
}

func (gjc *gjc) DeleteJobChains(cids ...int64) (cnt int64, err error) {
	if len(cids) == 0 {
		return
	}

	r := gjc.db.Table(gjc.tb).Where("id IN ?", cids).Delete(&xjm.JobChain{})
	cnt, err = r.RowsAffected, r.Error
	return
}

func (gjc *gjc) CleanOutdatedJobChains(before time.Time) (cnt int64, err error) {
	jss := xjm.JobChainAbortedCompleted

	r := gjc.db.Table(gjc.tb).Where("status IN ? AND updated_at < ?", jss, before).Delete(&xjm.JobChain{})
	cnt, err = r.RowsAffected, r.Error
	return
}
