package ref

import (
	"reflect"
	"runtime"
)

func NameOfFunc(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func IsZero(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	return !rv.IsValid() || rv.IsZero()
}
