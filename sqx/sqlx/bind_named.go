package sqlx

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
)

// bindStruct binds a named parameter query with fields from a struct argument.
// The rules for binding field names to parameter names follow the same
// conventions as for StructScan, including obeying the `db` struct tags.
func bindStruct(binder Binder, query string, arg any, m *ref.Mapper) (string, []any, error) {
	bound, names, err := compileNamedQuery(binder, query)
	if err != nil {
		return "", []any{}, err
	}

	arglist, err := bindAnyArgs(names, arg, m)
	if err != nil {
		return "", []any{}, err
	}

	return bound, arglist, nil
}

// bindArray binds a named parameter query with fields from an array or slice of
// structs argument.
func bindArray(binder Binder, query string, arg any, m *ref.Mapper) (string, []any, error) {
	// do the initial binding with QUESTION;  if binder is not question,
	// we can rebind it at the end.
	bound, names, err := compileNamedQuery(BindQuestion, query)
	if err != nil {
		return "", nil, err
	}

	arrayValue := reflect.ValueOf(arg)
	arrayLen := arrayValue.Len()
	if arrayLen == 0 {
		return "", nil, fmt.Errorf("length of array is 0: %#v", arg)
	}

	var arglist = make([]any, 0, len(names)*arrayLen)
	for i := 0; i < arrayLen; i++ {
		elemArglist, err := bindAnyArgs(names, arrayValue.Index(i).Interface(), m)
		if err != nil {
			return "", nil, err
		}
		arglist = append(arglist, elemArglist...)
	}
	if arrayLen > 1 {
		bound = fixBound(bound, arrayLen)
	}

	// adjust binding type if we weren't on question
	if binder != BindQuestion {
		bound = binder.Rebind(bound)
	}
	return bound, arglist, nil
}

// bindMap binds a named parameter query with a map of arguments.
func bindMap(binder Binder, query string, args map[string]any) (string, []any, error) {
	bound, names, err := compileNamedQuery(binder, query)
	if err != nil {
		return "", []any{}, err
	}

	arglist, err := bindMapArgs(names, args)
	return bound, arglist, err
}

// -- Compilation of Named Queries

// Allow digits and letters in bind params;  additionally runes are
// checked against underscores, meaning that bind params can have be
// alphanumeric with underscores.  Mind the difference between unicode
// digits and numbers, where '5' is a digit but '五' is not.
var allowedBindRunes = []*unicode.RangeTable{unicode.Letter, unicode.Digit}

// compile a NamedQuery into an unbound query (using the '?' bindvar) and a list of names.
func compileNamedQuery(binder Binder, qs string) (query string, names []string, err error) {
	names = make([]string, 0, str.CountByte(qs, ':'))
	rebound := &strings.Builder{}
	rebound.Grow(len(qs))

	inName := false
	vars := 0

	_, n := utf8.DecodeLastRuneInString(qs)
	last := len(qs) - n

	name := &strings.Builder{}
	name.Grow(16)

	for i, b := range qs {
		// a ':' while we're in a name is an error
		if b == ':' {
			if inName {
				// if this is the second ':' in a '::' escape sequence, append a ':'
				if i > 0 && qs[i-1] == ':' {
					rebound.WriteByte(':')
					inName = false
					continue
				}

				err = errors.New("unexpected `:` while reading named param at " + strconv.Itoa(i))
				return query, names, err
			}

			name.Reset()
			inName = true
			continue
		}

		if inName {
			if i > 0 && b == '=' && name.Len() == 0 {
				rebound.WriteString(":=")
				inName = false
				continue
			}

			// if we're in a name, and this is an allowed character, continue
			if (unicode.IsOneOf(allowedBindRunes, b) || b == '_' || b == '.') && i != last {
				// append the byte to the name if we are in a name and not on the last byte
				name.WriteRune(b)
				continue
			}

			// if we're in a name and it's not an allowed character, the name is done
			inName = false
			// if this is the final byte of the string and it is part of the name, then
			// make sure to add it to the name
			if i == last && unicode.IsOneOf(allowedBindRunes, b) {
				name.WriteRune(b)
			}

			// add the string representation to the names list
			names = append(names, name.String())
			// add a proper bindvar for the binder
			switch binder {
			// oracle only supports named type bind vars even for positional
			case BindColon:
				rebound.WriteByte(':')
				rebound.WriteString(name.String())
			case BindDollar:
				vars++
				rebound.WriteByte('$')
				rebound.WriteString(strconv.Itoa(vars))
			case BindAt:
				vars++
				rebound.WriteString("@p")
				rebound.WriteString(strconv.Itoa(vars))
			case BindQuestion, BindUnknown:
				rebound.WriteByte('?')
			}

			// add this byte to string unless it was not part of the name
			if i != last {
				rebound.WriteRune(b)
			} else if !unicode.IsOneOf(allowedBindRunes, b) {
				rebound.WriteRune(b)
			}
			continue
		}

		// this is a normal byte and should just go onto the rebound query
		rebound.WriteRune(b)
	}

	return rebound.String(), names, err
}

func bindNamedMapper(binder Binder, query string, arg any, m *ref.Mapper) (string, []any, error) {
	t := reflect.TypeOf(arg)
	k := t.Kind()
	switch {
	case k == reflect.Map && t.Key().Kind() == reflect.String:
		m, ok := convertMapStringInterface(arg)
		if !ok {
			return "", nil, fmt.Errorf("sqlx.bindNamedMapper: unsupported map type: %T", arg)
		}
		return bindMap(binder, query, m)
	case k == reflect.Array || k == reflect.Slice:
		return bindArray(binder, query, arg, m)
	default:
		return bindStruct(binder, query, arg, m)
	}
}

// convertMapStringInterface attempts to convert v to map[string]any.
// Unlike v.(map[string]any), this function works on named types that
// are convertible to map[string]any as well.
func convertMapStringInterface(v any) (map[string]any, bool) {
	var m map[string]any
	mtype := reflect.TypeOf(m)
	t := reflect.TypeOf(v)
	if !t.ConvertibleTo(mtype) {
		return nil, false
	}
	return reflect.ValueOf(v).Convert(mtype).Interface().(map[string]any), true

}

func bindAnyArgs(names []string, arg any, m *ref.Mapper) ([]any, error) {
	if maparg, ok := convertMapStringInterface(arg); ok {
		return bindMapArgs(names, maparg)
	}
	return bindArgs(names, arg, m)
}

// private interface to generate a list of interfaces from a given struct
// type, given a list of names to pull out of the struct.  Used by public
// BindStruct interface.
func bindArgs(names []string, arg any, m *ref.Mapper) ([]any, error) {
	arglist := make([]any, 0, len(names))

	// grab the indirected value of arg
	v := ref.DerefValue(reflect.ValueOf(arg))

	err := m.TraversalsByNameFunc(v.Type(), names, func(i int, t []int) error {
		if len(t) == 0 {
			return fmt.Errorf("could not find name %s in %#v", names[i], arg)
		}

		val := ref.FieldByIndexesReadOnly(v, t)
		arglist = append(arglist, val.Interface())

		return nil
	})

	return arglist, err
}

// like bindArgs, but for maps.
func bindMapArgs(names []string, arg map[string]any) ([]any, error) {
	arglist := make([]any, 0, len(names))

	for _, name := range names {
		val, ok := arg[name]
		if !ok {
			return arglist, fmt.Errorf("could not find name %s in %#v", name, arg)
		}
		arglist = append(arglist, val)
	}
	return arglist, nil
}

var valuesReg = regexp.MustCompile(`\)\s*(?i)VALUES\s*\(`)

func findMatchingClosingBracketIndex(s string) int {
	count := 0
	for i, ch := range s {
		if ch == '(' {
			count++
		}
		if ch == ')' {
			count--
			if count == 0 {
				return i
			}
		}
	}
	return 0
}

func fixBound(bound string, loop int) string {
	loc := valuesReg.FindStringIndex(bound)
	// defensive guard when "VALUES (...)" not found
	if len(loc) < 2 {
		return bound
	}

	openingBracketIndex := loc[1] - 1
	index := findMatchingClosingBracketIndex(bound[openingBracketIndex:])
	// defensive guard. must have closing bracket
	if index == 0 {
		return bound
	}
	closingBracketIndex := openingBracketIndex + index + 1

	var buffer bytes.Buffer

	buffer.WriteString(bound[0:closingBracketIndex])
	for i := 1; i < loop; i++ {
		buffer.WriteString(",")
		buffer.WriteString(bound[openingBracketIndex:closingBracketIndex])
	}
	buffer.WriteString(bound[closingBracketIndex:])
	return buffer.String()
}
