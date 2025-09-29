package pgsqlx

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/askasoft/pango/str"
)

// Vector is a wrapper for []float64 to implement sql.Scanner and driver.Valuer.
type Vector []float64

// Slice returns the []float64 slice.
func (v Vector) Slice() []float64 {
	return v
}

// String returns a string representation of the vector.
func (v Vector) String() string {
	if n := len(v); n > 0 {
		// There will be at least two brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '['

		b = strconv.AppendFloat(b, v[0], 'f', -1, 64)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, v[i], 'f', -1, 64)
		}

		return str.UnsafeString(append(b, ']'))
	}

	return "[]"
}

// Parse parses a string representation of a vector.
func (v *Vector) Parse(s string) error {
	if !str.StartsWithByte(s, '[') || !str.EndsWithByte(s, ']') {
		return nil
	}

	ss := strings.Split(s[1:len(s)-1], ",")

	a := make([]float64, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		a[i] = n
	}

	*v = a
	return nil
}

// Scan implements the sql.Scanner interface.
func (v *Vector) Scan(src any) (err error) {
	switch src := src.(type) {
	case []byte:
		return v.Parse(str.UnsafeString(src))
	case string:
		return v.Parse(src)
	case nil:
		*v = nil
		return nil
	default:
		return fmt.Errorf("pgsqlx: cannot convert %T to Vector", src)
	}
}

// Value implements the driver.Valuer interface.
func (v Vector) Value() (driver.Value, error) {
	if v == nil {
		return nil, nil
	}
	return v.String(), nil
}
