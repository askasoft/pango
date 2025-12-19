package vad

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

// extractTypeInternal gets the actual underlying type of field value.
// It will dive into pointers, customTypes and return you the
// underlying value and it's kind.
func (v *validate) extractTypeInternal(current reflect.Value, nullable bool) (reflect.Value, reflect.Kind, bool) {

BEGIN:
	switch current.Kind() {
	case reflect.Pointer:
		nullable = true

		if current.IsNil() {
			return current, reflect.Pointer, nullable
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Interface:
		nullable = true

		if current.IsNil() {
			return current, reflect.Interface, nullable
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Invalid:
		return current, reflect.Invalid, nullable

	default:
		if v.v.hasCustomFuncs {
			if fn, ok := v.v.customFuncs[current.Type()]; ok {
				current = reflect.ValueOf(fn(current))
				goto BEGIN
			}
		}

		return current, current.Kind(), nullable
	}
}

// getStructFieldOKInternal traverses a struct to retrieve a specific field denoted by the provided namespace and
// returns the field, field kind and whether is was successful in retrieving the field at all.
//
// NOTE: when not successful ok will be false, this can happen when a nested struct is nil and so the field
// could not be retrieved because it didn't exist.
func (v *validate) getStructFieldOKInternal(val reflect.Value, namespace string) (current reflect.Value, kind reflect.Kind, nullable bool, found bool) {

BEGIN:
	current, kind, nullable = v.ExtractType(val)
	if kind == reflect.Invalid {
		return
	}

	if namespace == "" {
		found = true
		return
	}

	switch kind {
	case reflect.Pointer, reflect.Interface:
		return

	case reflect.Struct:
		typ := current.Type()
		fld := namespace
		var ns string

		if !typ.ConvertibleTo(timeType) {
			idx := strings.Index(namespace, namespaceSeparator)

			if idx != -1 {
				fld = namespace[:idx]
				ns = namespace[idx+1:]
			} else {
				ns = ""
			}

			bracketIdx := strings.Index(fld, leftBracket)
			if bracketIdx != -1 {
				fld = fld[:bracketIdx]
				ns = namespace[bracketIdx:]
			}

			val = current.FieldByName(fld)
			namespace = ns
			goto BEGIN
		}

	case reflect.Array, reflect.Slice:
		idx := strings.Index(namespace, leftBracket)
		idx2 := strings.Index(namespace, rightBracket)

		arrIdx, _ := strconv.Atoi(namespace[idx+1 : idx2])
		if arrIdx >= current.Len() {
			return
		}

		startIdx := idx2 + 1
		if startIdx < len(namespace) {
			if namespace[startIdx:startIdx+1] == namespaceSeparator {
				startIdx++
			}
		}

		val = current.Index(arrIdx)
		namespace = namespace[startIdx:]
		goto BEGIN

	case reflect.Map:
		idx := strings.Index(namespace, leftBracket) + 1
		idx2 := strings.Index(namespace, rightBracket)

		endIdx := idx2
		if endIdx+1 < len(namespace) {
			if namespace[endIdx+1:endIdx+2] == namespaceSeparator {
				endIdx++
			}
		}

		key := namespace[idx:idx2]

		switch current.Type().Key().Kind() {
		case reflect.Int:
			i, _ := strconv.Atoi(key)
			val = current.MapIndex(reflect.ValueOf(i))
			namespace = namespace[endIdx+1:]

		case reflect.Int8:
			i, _ := strconv.ParseInt(key, 10, 8)
			val = current.MapIndex(reflect.ValueOf(int8(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Int16:
			i, _ := strconv.ParseInt(key, 10, 16)
			val = current.MapIndex(reflect.ValueOf(int16(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Int32:
			i, _ := strconv.ParseInt(key, 10, 32)
			val = current.MapIndex(reflect.ValueOf(int32(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Int64:
			i, _ := strconv.ParseInt(key, 10, 64)
			val = current.MapIndex(reflect.ValueOf(i))
			namespace = namespace[endIdx+1:]

		case reflect.Uint:
			i, _ := strconv.ParseUint(key, 10, 0)
			val = current.MapIndex(reflect.ValueOf(uint(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint8:
			i, _ := strconv.ParseUint(key, 10, 8)
			val = current.MapIndex(reflect.ValueOf(uint8(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint16:
			i, _ := strconv.ParseUint(key, 10, 16)
			val = current.MapIndex(reflect.ValueOf(uint16(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint32:
			i, _ := strconv.ParseUint(key, 10, 32)
			val = current.MapIndex(reflect.ValueOf(uint32(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint64:
			i, _ := strconv.ParseUint(key, 10, 64)
			val = current.MapIndex(reflect.ValueOf(i))
			namespace = namespace[endIdx+1:]

		case reflect.Float32:
			f, _ := strconv.ParseFloat(key, 32)
			val = current.MapIndex(reflect.ValueOf(float32(f)))
			namespace = namespace[endIdx+1:]

		case reflect.Float64:
			f, _ := strconv.ParseFloat(key, 64)
			val = current.MapIndex(reflect.ValueOf(f))
			namespace = namespace[endIdx+1:]

		case reflect.Bool:
			b, _ := strconv.ParseBool(key)
			val = current.MapIndex(reflect.ValueOf(b))
			namespace = namespace[endIdx+1:]

		// reflect.Type = string
		default:
			val = current.MapIndex(reflect.ValueOf(key))
			namespace = namespace[endIdx+1:]
		}

		goto BEGIN
	}

	// if got here there was more namespace, cannot go any deeper
	panic("Invalid field namespace")
}

var splitParamsRegex = regexp.MustCompile(`'[^']*'|\S+`)

// `'[^']*'|\S+`
func splits(param string) []string {
	vals := splitParamsRegex.FindAllString(param, -1)
	for i := 0; i < len(vals); i++ {
		vals[i] = strings.ReplaceAll(vals[i], "'", "")
	}
	return vals
}

func split2(tag, param string) (string, string) {
	ps := str.FieldsAny(param, "~～")
	if len(ps) != 2 {
		panic(fmt.Errorf("%s: invalid between expression %q", tag, param))
	}

	return str.Strip(ps[0]), str.Strip(ps[1])
}

// asInt returns the parameter as a int64
// or panics if it can't convert
func asInt(tag, param string) int64 {
	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(tag, err)
	return i
}

// asInt2 returns the parameter as two int64
// or panics if it can't convert
func asInt2(tag, param string) (int64, int64) {
	p1, p2 := split2(tag, param)
	i1 := asInt(tag, p1)
	i2 := asInt(tag, p2)
	return i1, i2
}

// asIntFromTimeDuration parses param as time.Duration and returns it as int64
// or panics on error.
func asIntFromTimeDuration(tag, param string) int64 {
	d, err := tmu.ParseDuration(param)
	if err != nil {
		// attempt parsing as an integer assuming nanosecond precision
		return asInt(tag, param)
	}
	return int64(d)
}

// asInt2FromTimeDuration parses param as time.Duration and returns it as two int64
// or panics on error.
func asInt2FromTimeDuration(tag, param string) (int64, int64) {
	p1, p2 := split2(tag, param)
	i1 := asIntFromTimeDuration(tag, p1)
	i2 := asIntFromTimeDuration(tag, p2)
	return i1, i2
}

// asIntFromType calls the proper function to parse param as int64,
// given a field's Type t.
func asIntFromType(tag string, t reflect.Type, param string) int64 {
	switch t {
	case timeDurationType:
		return asIntFromTimeDuration(tag, param)
	default:
		return asInt(tag, param)
	}
}

// asInt2FromType calls the proper function to parse param as int64,
// given a field's Type t.
func asInt2FromType(tag string, t reflect.Type, param string) (int64, int64) {
	switch t {
	case timeDurationType:
		return asInt2FromTimeDuration(tag, param)
	default:
		return asInt2(tag, param)
	}
}

// asUint returns the parameter as a uint64
// or panics if it can't convert
func asUint(tag, param string) uint64 {
	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(tag, err)
	return i
}

// asUint2 returns the parameter as two uint64
// or panics if it can't convert
func asUint2(tag, param string) (uint64, uint64) {
	p1, p2 := split2(tag, param)
	i1 := asUint(tag, p1)
	i2 := asUint(tag, p2)
	return i1, i2
}

// asFloat returns the parameter as a float64
// or panics if it can't convert
func asFloat(tag, param string) float64 {
	i, err := strconv.ParseFloat(param, 64)
	panicIf(tag, err)
	return i
}

// asFloat2 returns the parameter as two float64
// or panics if it can't convert
func asFloat2(tag, param string) (float64, float64) {
	p1, p2 := split2(tag, param)
	i1 := asFloat(tag, p1)
	i2 := asFloat(tag, p2)
	return i1, i2
}

// asBool returns the parameter as a bool
// or panics if it can't convert
func asBool(tag, param string) bool {
	i, err := strconv.ParseBool(param)
	panicIf(tag, err)
	return i
}

// asTime returns the parameter as a time
// or panics if it can't convert
func asTime(tag, param string) time.Time {
	if param == "" || param == "now" {
		return time.Now()
	}

	t, err := tmu.ParseInLocation(param, time.Local)
	panicIf(tag, err)
	return t
}

// asTime2 returns the parameter as two time
// or panics if it can't convert
func asTime2(tag, param string) (time.Time, time.Time) {
	p1, p2 := split2(tag, param)
	t1 := asTime(tag, p1)
	t2 := asTime(tag, p2)
	return t1, t2
}

func panicIf(tag string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", tag, err))
	}
}
