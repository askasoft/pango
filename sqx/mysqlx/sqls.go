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
