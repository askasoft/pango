package squ

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/str"
)

// Vector is a wrapper for []float64 to implement sql.Scanner and driver.Valuer.
type Vector struct {
	vec []float64
}

// NewVector creates a new Vector from a slice of float64.
func NewVector(vec []float64) Vector {
	return Vector{vec: vec}
}

// Slice returns the underlying slice of float64.
func (v Vector) Slice() []float64 {
	return v.vec
}

// String returns a string representation of the vector.
func (v Vector) String() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i := 0; i < len(v.vec); i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatFloat(float64(v.vec[i]), 'f', -1, 32))
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

	v.vec = make([]float64, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.vec[i] = n
	}
	return nil
}

// Scan implements the sql.Scanner interface.
func (v *Vector) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return v.Parse(bye.UnsafeString(src))
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
