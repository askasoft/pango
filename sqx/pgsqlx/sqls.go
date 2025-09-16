package pgsqlx

import (
	"fmt"
)

func ResetSequenceSQL(table, column string, starts ...int64) string {
	start := int64(1)
	if len(starts) > 0 {
		start = starts[0]
	}

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
