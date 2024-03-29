package vad

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// IsEqual returns whether val1 is equal to val2 taking into account Pointers, Interfaces and their underlying types
func assertIsEqual(val1, val2 any) bool {
	v1 := reflect.ValueOf(val1)
	v2 := reflect.ValueOf(val2)

	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}

	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}

	if !v1.IsValid() && !v2.IsValid() {
		return true
	}

	switch v1.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if v1.IsNil() {
			v1 = reflect.ValueOf(nil)
		}
	}

	switch v2.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if v2.IsNil() {
			v2 = reflect.ValueOf(nil)
		}
	}

	v1Underlying := reflect.Zero(reflect.TypeOf(v1)).Interface()
	v2Underlying := reflect.Zero(reflect.TypeOf(v2)).Interface()

	if v1 == v1Underlying {
		if v2 == v2Underlying {
			return reflect.DeepEqual(v1, v2)
		}
		return reflect.DeepEqual(v1, v2.Interface())
	}

	if v2 == v2Underlying {
		return reflect.DeepEqual(v1.Interface(), v2)
	}

	return reflect.DeepEqual(v1.Interface(), v2.Interface())
}

// NotMatchRegex validates that value matches the regex, either string or *regex
// and throws an error with line number
// func assertNotMatchRegex(t *testing.T, value string, regex any) {
// 	assertNotMatchRegexSkip(t, 2, value, regex)
// }

// NotMatchRegexSkip validates that value matches the regex, either string or *regex
// and throws an error with line number
// but the skip variable tells NotMatchRegexSkip how far back on the stack to report the error.
// This is a building block to creating your own more complex validation functions.
// func assertNotMatchRegexSkip(t *testing.T, skip int, value string, regex any) {
// 	if r, ok, err := assertRegexMatches(regex, value); ok || err != nil {
// 		_, file, line, _ := runtime.Caller(skip)

// 		if err != nil {
// 			fmt.Printf("%s:%d %v error compiling regex %v\n", filepath.Base(file), line, value, r.String())
// 		} else {
// 			fmt.Printf("%s:%d %v matches regex %v\n", filepath.Base(file), line, value, r.String())
// 		}

// 		t.FailNow()
// 	}
// }

// MatchRegex validates that value matches the regex, either string or *regex
// and throws an error with line number
// func assertMatchRegex(t *testing.T, value string, regex any) {
// 	assertMatchRegexSkip(t, 2, value, regex)
// }

// MatchRegexSkip validates that value matches the regex, either string or *regex
// and throws an error with line number
// but the skip variable tells MatchRegexSkip how far back on the stack to report the error.
// This is a building block to creating your own more complex validation functions.
// func assertMatchRegexSkip(t *testing.T, skip int, value string, regex any) {
// 	if r, ok, err := assertRegexMatches(regex, value); !ok {
// 		_, file, line, _ := runtime.Caller(skip)

// 		if err != nil {
// 			fmt.Printf("%s:%d %v error compiling regex %v\n", filepath.Base(file), line, value, r.String())
// 		} else {
// 			fmt.Printf("%s:%d %v does not match regex %v\n", filepath.Base(file), line, value, r.String())
// 		}

// 		t.FailNow()
// 	}
// }

// func assertRegexMatches(regex any, value string) (*regexp.Regexp, bool, error) {
// 	var err error

// 	r, ok := regex.(*regexp.Regexp)

// 	// must be a string
// 	if !ok {
// 		if r, err = regexp.Compile(regex.(string)); err != nil {
// 			return r, false, err
// 		}
// 	}

// 	return r, r.MatchString(value), err
// }

// Equal validates that val1 is equal to val2 and throws an error with line number
func assertEqual(t *testing.T, val1, val2 any) {
	assertEqualSkip(t, 2, val1, val2)
}

// EqualSkip validates that val1 is equal to val2 and throws an error with line number
// but the skip variable tells EqualSkip how far back on the stack to report the error.
// This is a building block to creating your own more complex validation functions.
func assertEqualSkip(t *testing.T, skip int, val1, val2 any) {
	if !assertIsEqual(val1, val2) {
		_, file, line, _ := runtime.Caller(skip)
		fmt.Printf("%s:%d %v does not equal %v\n", filepath.Base(file), line, val1, val2)
		t.FailNow()
	}
}

// NotEqual validates that val1 is not equal val2 and throws an error with line number
func assertNotEqual(t *testing.T, val1, val2 any) {
	assertNotEqualSkip(t, 2, val1, val2)
}

// NotEqualSkip validates that val1 is not equal to val2 and throws an error with line number
// but the skip variable tells NotEqualSkip how far back on the stack to report the error.
// This is a building block to creating your own more complex validation functions.
func assertNotEqualSkip(t *testing.T, skip int, val1, val2 any) {
	if assertIsEqual(val1, val2) {
		_, file, line, _ := runtime.Caller(skip)
		fmt.Printf("%s:%d %v should not be equal %v\n", filepath.Base(file), line, val1, val2)
		t.FailNow()
	}
}

// PanicMatches validates that the panic output of running fn matches the supplied string
func assertPanicMatches(t *testing.T, fn func(), matches string) {
	assertPanicMatchesSkip(t, 2, fn, matches)
}

// PanicMatchesSkip validates that the panic output of running fn matches the supplied string
// but the skip variable tells PanicMatchesSkip how far back on the stack to report the error.
// This is a building block to creating your own more complex validation functions.
func assertPanicMatchesSkip(t *testing.T, skip int, fn func(), matches string) {
	_, file, line, _ := runtime.Caller(skip)

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)

			if err != matches {
				fmt.Printf("%s:%d Panic...  expected [%s] received [%s]", filepath.Base(file), line, matches, err)
				t.FailNow()
			}
		} else {
			fmt.Printf("%s:%d Panic Expected, none found...  expected [%s]", filepath.Base(file), line, matches)
			t.FailNow()
		}
	}()

	fn()
}
