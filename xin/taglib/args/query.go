package args

import (
	"github.com/askasoft/pango/str"
)

type Query struct {
	Key string `form:"k"`
	Sorter
	Pager
}

func (q *Query) Normalize(columns []string, limits []int) {
	q.Key = str.Strip(q.Key)
	q.Sorter.Normalize(columns...)
	q.Pager.Normalize(limits...)
}
