package mygormsm

import (
	"errors"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/xsm"
	"github.com/askasoft/pango/xsm/mysm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gsm struct {
	db *gorm.DB
}

func SM(db *gorm.DB) xsm.SchemaManager {
	return &gsm{db}
}

func (gsm *gsm) ExistsSchema(s string) (bool, error) {
	if asg.Contains(mysm.SysDBs, s) {
		return false, nil
	}

	var sn string
	err := gsm.db.Table("information_schema.schemata").Where("schema_name = ?", s).Select("schema_name").Take(&sn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (gsm *gsm) ListSchemas() ([]string, error) {
	tx := gsm.db.Table("information_schema.schemata")
	tx = tx.Where("schema_name NOT IN ?", mysm.SysDBs)
	tx = tx.Select("schema_name")
	tx = tx.Order("schema_name asc")

	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sn string

	var ss []string
	for rows.Next() {
		if err = rows.Scan(&sn); err != nil {
			return nil, err
		}
		ss = append(ss, sn)
	}
	return ss, nil
}

func (gsm *gsm) CreateSchema(name, comment string) error {
	err := gsm.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(mysm.SQLCreateSchema(name)).Error; err != nil {
			return err
		}
		if comment != "" {
			if err := tx.Exec(mysm.SQLCommentSchema(name, comment)).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (gsm *gsm) CommentSchema(name string, comment string) error {
	return gsm.db.Exec(mysm.SQLCommentSchema(name, comment)).Error
}

func (gsm *gsm) RenameSchema(old string, new string) error {
	return errors.New("unsupport")
}

func (gsm *gsm) DeleteSchema(name string) error {
	return gsm.db.Exec(mysm.SQLDeleteSchema(name)).Error
}

func (gsm *gsm) buildQuery(sq *xsm.SchemaQuery) *gorm.DB {
	tx := gsm.db.Table("information_schema.schemata")

	tx = tx.Where("schema_name NOT IN ?", mysm.SysDBs)
	if sq.Name != "" {
		tx = tx.Where("schema_name LIKE ?", sqx.StringLike(sq.Name))
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
		"schema_name AS name",
		"(SELECT SUM(data_length + index_length) FROM information_schema.tables WHERE table_schema = schema_name) AS size",
		"schema_comment AS comment",
	)

	tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: sq.Col}, Desc: sq.IsDesc()})
	if sq.Col != "name" {
		tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: sq.IsDesc()})
	}
	tx = tx.Offset(sq.Start()).Limit(sq.Limit)

	err = tx.Find(&schemas).Error
	return
}
