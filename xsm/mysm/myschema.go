package mysm

import (
	"fmt"

	"github.com/askasoft/pango/sqx"
)

var SysDBs = []string{"information_schema", "mysql", "performance_schema", "sys"}

func SQLCreateSchema(name string) string {
	return "CREATE SCHEMA " + name
}

func SQLCommentSchema(name string, comment string) string {
	return fmt.Sprintf("ALTER SCHEMA %s COMMENT = '%s'", name, sqx.EscapeString(comment))
}

func SQLDeleteSchema(name string) string {
	return fmt.Sprintf("DROP SCHEMA " + name)
}
