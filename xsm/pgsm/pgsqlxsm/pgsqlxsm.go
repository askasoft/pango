package pgsqlxsm

import (
	"errors"

	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xsm"
	"github.com/askasoft/pango/xsm/pgsm"
)

type ssm struct {
	db *sqlx.DB
}

func SM(db *sqlx.DB) xsm.SchemaManager {
	return &ssm{
		db: db,
	}
}

func (ssm *ssm) ExistsSchema(s string) (bool, error) {
	if str.ContainsByte(s, '_') {
		return false, nil
	}

	sqb := ssm.db.Builder()
	sqb.Select("nspname").From(pgsm.TablePgNamespace).Where("nspname = ?", s)
	sql, args := sqb.Build()

	pn := &pgsm.PgNamesapce{}
	err := ssm.db.Get(pn, sql, args...)
	if err != nil {
		if errors.Is(err, sqlx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ssm *ssm) ListSchemas() ([]string, error) {
	sqb := ssm.db.Builder()
	sqb.Select("nspname").From(pgsm.TablePgNamespace).Where("nspname NOT LIKE ?", sqx.StringLike("_")).Order("nspname", false)
	sql, args := sqb.Build()

	rows, err := ssm.db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ss []string

	pn := &pgsm.PgNamesapce{}
	for rows.Next() {
		if err = rows.StructScan(pn); err != nil {
			return nil, err
		}
		ss = append(ss, pn.Nspname)
	}

	return ss, nil
}

func (ssm *ssm) CreateSchema(name, comment string) error {
	err := ssm.db.Transaction(func(tx *sqlx.Tx) error {
		if _, err := tx.Exec(pgsm.SQLCreateSchema(name)); err != nil {
			return err
		}
		if comment != "" {
			if _, err := tx.Exec(pgsm.SQLCommentSchema(name, comment)); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (ssm *ssm) CommentSchema(name string, comment string) error {
	_, err := ssm.db.Exec(pgsm.SQLCommentSchema(name, comment))
	return err
}

func (ssm *ssm) RenameSchema(old string, new string) error {
	_, err := ssm.db.Exec(pgsm.SQLRenameSchema(old, new))
	return err
}

func (ssm *ssm) DeleteSchema(name string) error {
	_, err := ssm.db.Exec(pgsm.SQLDeleteSchema(name))
	return err
}

func (ssm *ssm) buildQuery(sq *xsm.SchemaQuery) *sqlx.Builder {
	sqb := ssm.db.Builder()

	sqb.From(pgsm.TablePgNamespace)
	sqb.Where("nspname NOT LIKE ?", sqx.StringLike("_"))

	if sq.Name != "" {
		sqb.Where("nspname LIKE ?", sqx.StringLike(sq.Name))
	}
	return sqb
}

func (ssm *ssm) CountSchemas(sq *xsm.SchemaQuery) (total int, err error) {
	sqb := ssm.buildQuery(sq).Select("COUNT(1)")
	sql, args := sqb.Build()

	err = ssm.db.Get(&total, sql, args...)
	return
}

func (ssm *ssm) FindSchemas(sq *xsm.SchemaQuery) (schemas []*xsm.SchemaInfo, err error) {
	sqb := ssm.buildQuery(sq)
	sqb.Select(
		"nspname AS name",
		"COALESCE((SELECT SUM(pg_relation_size(oid)) FROM pg_catalog.pg_class WHERE relnamespace = pg_namespace.oid), 0) AS size",
		"COALESCE(obj_description(oid, 'pg_namespace'), '') AS comment",
	)

	sqb.Order(sq.Col, sq.IsDesc())
	if sq.Col != "name" {
		sqb.Order("name", sq.IsDesc())
	}
	sqb.Offset(sq.Start()).Limit(sq.Limit)

	sql, args := sqb.Build()

	err = ssm.db.Select(&schemas, sql, args...)
	return
}
