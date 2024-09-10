package pggormsm

import (
	"errors"

	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xsm"
	"github.com/askasoft/pango/xsm/pgsm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gsm struct {
	db *gorm.DB
}

func SM(db *gorm.DB) xsm.SchemaManager {
	return &gsm{
		db: db,
	}
}

func (gsm *gsm) ExistsSchema(s string) (bool, error) {
	if str.ContainsByte(s, '_') {
		return false, nil
	}

	pn := &pgsm.PgNamesapce{}
	err := gsm.db.Table(pgsm.TablePgNamespace).Where("nspname = ?", s).Select("nspname").Take(pn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (gsm *gsm) ListSchemas() ([]string, error) {
	tx := gsm.db.Table(pgsm.TablePgNamespace).Where("nspname NOT LIKE ?", sqx.StringLike("_")).Select("nspname").Order("nspname asc")
	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ss []string

	pn := &pgsm.PgNamesapce{}
	for rows.Next() {
		if err = tx.ScanRows(rows, pn); err != nil {
			return nil, err
		}
		ss = append(ss, pn.Nspname)
	}

	return ss, nil
}

func (gsm *gsm) CreateSchema(name, comment string) error {
	err := gsm.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(pgsm.SQLCreateSchema(name)).Error; err != nil {
			return err
		}
		if comment != "" {
			if err := tx.Exec(pgsm.SQLCommentSchema(name, comment)).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (gsm *gsm) CommentSchema(name string, comment string) error {
	return gsm.db.Exec(pgsm.SQLCommentSchema(name, comment)).Error
}

func (gsm *gsm) RenameSchema(old string, new string) error {
	return gsm.db.Exec(pgsm.SQLRenameSchema(old, new)).Error
}

func (gsm *gsm) DeleteSchema(name string) error {
	return gsm.db.Exec(pgsm.SQLDeleteSchema(name)).Error
}

func (gsm *gsm) buildQuery(sq *xsm.SchemaQuery) *gorm.DB {
	tx := gsm.db.Table(pgsm.TablePgNamespace)

	tx = tx.Where("nspname NOT LIKE ?", sqx.StringLike("_"))
	if sq.Name != "" {
		tx = tx.Where("nspname LIKE ?", sqx.StringLike(sq.Name))
	}
	return tx
}

func (gsm *gsm) CountSchemas(sq *xsm.SchemaQuery) (total int, err error) {
	var cnt int64
	err = gsm.buildQuery(sq).Count(&cnt).Error
	total = int(cnt)
	return
}

func (gsm *gsm) FindSchemas(sq *xsm.SchemaQuery) (schemas []*xsm.SchemaInfo, err error) {
	tx := gsm.buildQuery(sq)
	tx = tx.Select(
		"nspname AS name",
		"COALESCE((SELECT SUM(pg_relation_size(oid)) FROM pg_catalog.pg_class WHERE relnamespace = pg_namespace.oid), 0) AS size",
		"COALESCE(obj_description(oid, 'pg_namespace'), '') AS comment",
	)

	tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: sq.Col}, Desc: sq.IsDesc()})
	if sq.Col != "name" {
		tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: sq.IsDesc()})
	}
	tx = tx.Offset(sq.Start()).Limit(sq.Limit)

	err = tx.Find(&schemas).Error
	return
}
