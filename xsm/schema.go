package xsm

import (
	"github.com/askasoft/pango/xvw/args"
)

type SchemaInfo struct {
	Name    string `json:"name" form:"name,strip,lower" validate:"required,maxlen=30,regexp=^[a-z][a-z0-9]{00x2C29}$"`
	Size    int64  `json:"size,omitempty"`
	Comment string `json:"comment,omitempty" form:"comment" validate:"omitempty,maxlen=250"`
}

type SchemaQuery struct {
	args.Pager
	args.Sorter
	Name string `json:"name" form:"name,strip"`
}

func (sq *SchemaQuery) Normalize(limits []int) {
	sq.Sorter.Normalize(
		"name",
		"comment",
		"size",
	)
	sq.Pager.Normalize(limits...)
}

type SchemaManager interface {
	ExistsSchema(s string) (bool, error)
	ListSchemas() ([]string, error)
	CreateSchema(name string, comment string) error
	CommentSchema(name string, comment string) error
	RenameSchema(old string, new string) error
	DeleteSchema(name string) error
	CountSchemas(sq *SchemaQuery) (total int, err error)
	FindSchemas(sq *SchemaQuery) (schemas []*SchemaInfo, err error)
}
