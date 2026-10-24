package mysqlx

import (
	"fmt"

	"github.com/askasoft/pango/asg"
)

func ResetAutoIncrementSQL(table string, starts ...int64) string {
	start := max(asg.First(starts), 1)

	sql := fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = %d", table, start)

	return sql
}

func NextAutoIncrementSQL(schema, table string) string {
	return fmt.Sprintf("SELECT AUTO_INCREMENT FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s'", schema, table)
}
