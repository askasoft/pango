package pgsqlx

import (
	"fmt"

	"github.com/askasoft/pango/asg"
)

func ResetSequenceSQL(table, column string, starts ...int64) string {
	start := max(asg.First(starts), 1)

	return fmt.Sprintf(
		"SELECT SETVAL(pg_get_serial_sequence('%s', '%s'), GREATEST((SELECT MAX(%s)+1 FROM %s), %d), false)",
		table,
		column,
		column,
		table,
		start,
	)
}

func NextSequenceSQL(table, column string) string {
	return fmt.Sprintf("SELECT NEXTVAL(pg_get_serial_sequence('%s', '%s'))", table, column)
}
