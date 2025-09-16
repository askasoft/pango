package pgsqlx

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/askasoft/pango/str"
)

// BoolArray represents a one-dimensional array of the PostgreSQL boolean type.
type BoolArray []bool

func (a BoolArray) Slice() []bool {
	return a
}

// Scan implements the sql.Scanner interface.
func (a *BoolArray) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes(str.UnsafeBytes(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pqx: cannot convert %T to BoolArray", src)
}

func (a *BoolArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, "BoolArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(BoolArray, len(elems))
		for i, v := range elems {
			if len(v) != 1 {
				return fmt.Errorf("pqx: could not parse boolean array index %d: invalid boolean %q", i, v)
			}
			switch v[0] {
			case 't':
				b[i] = true
			case 'f':
				b[i] = false
			default:
				return fmt.Errorf("pqx: could not parse boolean array index %d: invalid boolean %q", i, v)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a BoolArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be exactly two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1+2*n)

		for i := 0; i < n; i++ {
			b[2*i] = ','
			if a[i] {
				b[1+2*i] = 't'
			} else {
				b[1+2*i] = 'f'
			}
		}

		b[0] = '{'
		b[2*n] = '}'

		return str.UnsafeString(b), nil
	}

	return "{}", nil
}

// Float64Array represents a one-dimensional array of the PostgreSQL double
// precision type.
type Float64Array []float64

func (a Float64Array) Slice() []float64 {
	return a
}

// Scan implements the sql.Scanner interface.
func (a *Float64Array) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes(str.UnsafeBytes(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pqx: cannot convert %T to Float64Array", src)
}

func (a *Float64Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, "Float64Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Float64Array, len(elems))
		for i, v := range elems {
			if b[i], err = strconv.ParseFloat(str.UnsafeString(v), 64); err != nil {
				return fmt.Errorf("pqx: parsing array element index %d: %w", i, err)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Float64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendFloat(b, a[0], 'f', -1, 64)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, a[i], 'f', -1, 64)
		}

		return str.UnsafeString(append(b, '}')), nil
	}

	return "{}", nil
}

// Float32Array represents a one-dimensional array of the PostgreSQL double
// precision type.
type Float32Array []float32

func (a Float32Array) Slice() []float32 {
	return a
}

// Scan implements the sql.Scanner interface.
func (a *Float32Array) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes(str.UnsafeBytes(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pqx: cannot convert %T to Float32Array", src)
}

func (a *Float32Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, "Float32Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Float32Array, len(elems))
		for i, v := range elems {
			var x float64
			if x, err = strconv.ParseFloat(str.UnsafeString(v), 32); err != nil {
				return fmt.Errorf("pqx: parsing array element index %d: %w", i, err)
			}
			b[i] = float32(x)
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Float32Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendFloat(b, float64(a[0]), 'f', -1, 32)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, float64(a[i]), 'f', -1, 32)
		}

		return str.UnsafeString(append(b, '}')), nil
	}

	return "{}", nil
}

// Int64Array represents a one-dimensional array of the PostgreSQL integer types.
type Int64Array []int64

func (a Int64Array) Slice() []int64 {
	return a
}

// Scan implements the sql.Scanner interface.
func (a *Int64Array) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes(str.UnsafeBytes(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pqx: cannot convert %T to Int64Array", src)
}

func (a *Int64Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, "Int64Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Int64Array, len(elems))
		for i, v := range elems {
			if b[i], err = strconv.ParseInt(str.UnsafeString(v), 10, 64); err != nil {
				return fmt.Errorf("pqx: parsing array element index %d: %w", i, err)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Int64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendInt(b, a[0], 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, a[i], 10)
		}

		return str.UnsafeString(append(b, '}')), nil
	}

	return "{}", nil
}

// Int32Array represents a one-dimensional array of the PostgreSQL integer types.
type Int32Array []int32

func (a Int32Array) Slice() []int32 {
	return a
}

// Scan implements the sql.Scanner interface.
func (a *Int32Array) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes(str.UnsafeBytes(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pqx: cannot convert %T to Int32Array", src)
}

func (a *Int32Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, "Int32Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Int32Array, len(elems))
		for i, v := range elems {
			x, err := strconv.ParseInt(str.UnsafeString(v), 10, 32)
			if err != nil {
				return fmt.Errorf("pqx: parsing array element index %d: %w", i, err)
			}
			b[i] = int32(x)
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Int32Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendInt(b, int64(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, int64(a[i]), 10)
		}

		return str.UnsafeString(append(b, '}')), nil
	}

	return "{}", nil
}

// IntArray represents a one-dimensional array of the PostgreSQL integer types.
type IntArray []int

func (a IntArray) Slice() []int {
	return a
}

// Scan implements the sql.Scanner interface.
func (a *IntArray) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes(str.UnsafeBytes(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pqx: cannot convert %T to Int32Array", src)
}

func (a *IntArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, "IntArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(IntArray, len(elems))
		for i, v := range elems {
			x, err := strconv.ParseInt(str.UnsafeString(v), 10, strconv.IntSize)
			if err != nil {
				return fmt.Errorf("pqx: parsing array element index %d: %w", i, err)
			}
			b[i] = int(x)
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendInt(b, int64(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, int64(a[i]), 10)
		}

		return str.UnsafeString(append(b, '}')), nil
	}

	return "{}", nil
}

// StringArray represents a one-dimensional array of the PostgreSQL character types.
type StringArray []string

func (a StringArray) Slice() []string {
	return a
}

// Scan implements the sql.Scanner interface.
func (a *StringArray) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes(str.UnsafeBytes(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pqx: cannot convert %T to StringArray", src)
}

func (a *StringArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, "StringArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(StringArray, len(elems))
		for i, v := range elems {
			if b[i] = str.UnsafeString(v); v == nil {
				return fmt.Errorf("pqx: parsing array element index %d: cannot convert nil to string", i)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+3*n)
		b[0] = '{'

		b = appendArrayQuotedBytes(b, str.UnsafeBytes(a[0]))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendArrayQuotedBytes(b, str.UnsafeBytes(a[i]))
		}

		return str.UnsafeString(append(b, '}')), nil
	}

	return "{}", nil
}

func appendArrayQuotedBytes(b, v []byte) []byte {
	b = append(b, '"')
	for {
		i := bytes.IndexAny(v, `"\`)
		if i < 0 {
			b = append(b, v...)
			break
		}
		if i > 0 {
			b = append(b, v[:i]...)
		}
		b = append(b, '\\', v[i])
		v = v[i+1:]
	}
	return append(b, '"')
}

// parseArray extracts the dimensions and elements of an array represented in
// text format. Only representations emitted by the backend are supported.
// Notably, whitespace around brackets and delimiters is significant, and NULL
// is case-sensitive.
//
// See http://www.postgresql.org/docs/current/static/arrays.html#ARRAYS-IO
var null = []byte("NULL")

func parseArray(src []byte) (dims []int, elems [][]byte, err error) {
	var depth, i int

	if len(src) < 1 || src[0] != '{' {
		return nil, nil, fmt.Errorf("pqx: unable to parse array; expected %q at offset %d", '{', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '{':
			depth++
			i++
		case '}':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '{':
			if depth == len(dims) {
				break Element
			}
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem = []byte{}
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if src[i] == ',' || src[i] == '}' {
					elem := src[start:i]
					if len(elem) == 0 {
						return nil, nil, fmt.Errorf("pqx: unable to parse array; unexpected %q at offset %d", src[i], i)
					}
					if bytes.Equal(elem, null) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if src[i] == ',' && depth > 0 {
			dims[depth-1]++
			i++
			goto Element
		} else if src[i] == '}' && depth > 0 {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pqx: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == '}' && depth > 0 {
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pqx: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("pqx: unable to parse array; expected %q at offset %d", '}', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("pqx: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}

func scanLinearArray(src []byte, typ string) (elems [][]byte, err error) {
	dims, elems, err := parseArray(src)
	if err != nil {
		return nil, err
	}
	if len(dims) > 1 {
		return nil, fmt.Errorf("pqx: cannot convert ARRAY%s to %s", strings.ReplaceAll(fmt.Sprint(dims), " ", "]["), typ)
	}
	return elems, err
}
