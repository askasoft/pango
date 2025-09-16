package mysqlx

import "fmt"

func ResetAutoIncrementSQL(table string, starts ...int64) string {
	start := int64(1)
	if len(starts) > 0 {
		start = starts[0]
	}

	sql := fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = %d", table, start)

	return sql
}
