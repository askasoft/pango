package myx

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

func IsUniqueViolationError(err error) bool {
	if err != nil {
		var myErr *mysql.MySQLError
		if ok := errors.As(err, &myErr); ok {
			if myErr.Number == 1062 {
				return true
			}
		}
	}
	return false
}
