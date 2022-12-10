package ref

import (
	"reflect"
	"testing"
	"time"
)

func TestStringCompare(t *testing.T) {
	cs := []struct {
		e bool
		a bool
		m string
	}{
		{true, any("") != nil, `any("") != nil`},
		{true, any("") == "", `any("") == ""`},
	}

	for _, c := range cs {
		if c.a != c.e {
			t.Error(c.m)
		}
	}
}

func TestConvert(t *testing.T) {
	cs := []struct {
		e any
		s any
		b reflect.Type
	}{
		{"1.123", 1.123, reflect.TypeOf("")},
		{int32(0777), "0777", reflect.TypeOf(int32(0))},
		{time.Second, "1s", reflect.TypeOf(time.Duration(0))},
		{time.Second, 1000000000, reflect.TypeOf(time.Duration(0))},
		{utcMilli(0).Add(time.Second), "1970-01-01T00:00:01Z", reflect.TypeOf(time.Time{})},
		{utcMilli(0).Add(time.Second), 1000, reflect.TypeOf(time.Time{})},
		{int(0), nil, reflect.TypeOf(int(0))},
		{int32(0), "", reflect.TypeOf(int32(0))},
	}
	for i, c := range cs {
		a, err := Convert(c.s, c.b)
		if err != nil {
			t.Errorf("[%d] Convert(%q, %q) Failed: %v, want %v", i, c.s, c.b, err, c.e)
		}
		if a != c.e {
			t.Errorf("[%d] Convert(%q, %q) = %v, want %v", i, c.s, c.b, a, c.e)
		}
	}
}
