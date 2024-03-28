package sqx

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/askasoft/pango/str"
)

// Array is a wrapper for []string to implement sql.Scanner and driver.Valuer.
type Array []string

// NewArray creates a new Array from a slice of string.
func NewArray(a []string) Array {
	return Array(a)
}

// Slice returns the underlying slice of string.
func (a Array) Slice() []string {
	return a
}

// String returns a string representation of the array.
func (a Array) String() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i := 0; i < len(a); i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.Quote(a[i]))
	}
	sb.WriteString("]")

	return sb.String()
}

// Parse parses a string representation of a array.
func (a *Array) Parse(s string) error {
	if !str.StartsWithByte(s, '[') || !str.EndsWithByte(s, ']') {
		return nil
	}

	ss := strings.Split(s[1:len(s)-1], ",")

	v := make([]string, len(ss))
	for i, s := range ss {
		n, err := strconv.Unquote(s)
		if err != nil {
			return err
		}
		v[i] = n
	}

	*a = v
	return nil
}

// Scan implements the sql.Scanner interface.
func (a *Array) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return a.Parse(str.UnsafeString(src))
	case string:
		return a.Parse(src)
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

// Value implements the driver.Valuer interface.
func (a Array) Value() (driver.Value, error) {
	return a.String(), nil
}
