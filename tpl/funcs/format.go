package funcs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"

	"github.com/askasoft/pango/num"
)

// JSON returns a json marshal string.
func JSON(a any) (template.JS, error) {
	bs, err := json.Marshal(a)
	return template.JS(bs), err //nolint: gosec
}

// Comma produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. Comma(834142) -> 834,142
// e.g. Comma(834142.1234, 3) -> 834,142.123
// e.g. Comma(834142, "_") -> 834_142
func Comma(n any, args ...any) (string, error) {
	c, d := ",", 0
	for _, a := range args {
		switch p := a.(type) {
		case string:
			c = p
		case int:
			d = p
		}
	}

	v := reflect.ValueOf(n)
	switch v.Kind() {
	case reflect.Int8, reflect.Uint8:
		return num.Itoa(int(v.Int())), nil
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return num.CommaInt(v.Int(), c), nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return num.CommaUint(v.Uint(), c), nil
	case reflect.Float32, reflect.Float64:
		if d > 0 {
			return num.CommaFloatWithDigits(v.Float(), d, c), nil
		}
		return num.CommaFloat(v.Float(), c), nil
	default:
		return "", fmt.Errorf("Comma: unknown type for '%v' (%T)", n, n)
	}
}

// HumanSize returns a human-readable approximation of a size
// with specified precision digit numbers (default: 2) (eg. "2.75 MB", "796 KB").
func HumanSize(n any, p ...int) (string, error) {
	v := reflect.ValueOf(n)

	switch v.Kind() {
	case reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return num.HumanSize(float64(v.Int()), p...), nil
	case reflect.Float32, reflect.Float64:
		return num.HumanSize(v.Float(), p...), nil
	default:
		return "", fmt.Errorf("HumanSize: unknown type for '%v' (%T)", n, n)
	}
}
