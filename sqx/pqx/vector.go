package pqx

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
	var sb strings.Builder

	sb.WriteString("[")
	for i := 0; i < len(v); i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatFloat(float64(v[i]), 'f', -1, 32))
	}
	sb.WriteString("]")

	return sb.String()
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
func (v *Vector) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return v.Parse(str.UnsafeString(src))
	case string:
		return v.Parse(src)
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

// Value implements the driver.Valuer interface.
func (v Vector) Value() (driver.Value, error) {
	return v.String(), nil
}
