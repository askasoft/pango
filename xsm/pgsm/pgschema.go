package pgsm

import (
	"fmt"

	"github.com/askasoft/pango/sqx"
)

func SQLCreateSchema(name string) string {
	return "CREATE SCHEMA " + name
}

func SQLCommentSchema(name string, comment string) string {
	return fmt.Sprintf("COMMENT ON SCHEMA %s IS '%s'", name, sqx.EscapeString(comment))
}

func SQLRenameSchema(old string, new string) string {
	return fmt.Sprintf("ALTER SCHEMA %s RENAME TO %s", old, new)
}

func SQLDeleteSchema(name string) string {
	return fmt.Sprintf("DROP SCHEMA %s CASCADE", name)
}
