package pgsqlx

import (
	"fmt"

	"github.com/askasoft/pango/asg"
)

func ResetSequenceSQL(table, column string, starts ...int64) string {
	start := max(asg.First(starts), 1)

	sql := fmt.Sprintf(
		"SELECT SETVAL('%s_%s_seq', GREATEST((SELECT MAX(%s)+1 FROM %s), %d), false)",
		table,
		column,
		column,
		table,
		start,
	)

	return sql
}
