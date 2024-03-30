package sqx

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/askasoft/pango/str"
)

// StringArray is sa wrapper for []string to implement sql.Scanner and driver.Valuer.
type StringArray []string

// Slice returns the []string slice.
func (sa StringArray) Slice() []string {
	return sa
}

// String returns sa string representation of the array.
func (sa StringArray) String() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i := 0; i < len(sa); i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.Quote(sa[i]))
	}
	sb.WriteString("]")

	return sb.String()
}

// Parse parses sa string representation of sa array.
func (sa *StringArray) Parse(s string) error {
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

	*sa = v
	return nil
}

// Scan implements the sql.Scanner interface.
func (sa *StringArray) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return sa.Parse(str.UnsafeString(src))
	case string:
		return sa.Parse(src)
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

// Value implements the driver.Valuer interface.
func (sa StringArray) Value() (driver.Value, error) {
	return sa.String(), nil
}
