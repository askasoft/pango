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

func utcMilli(msec int64) time.Time {
	return time.Unix(msec/1e3, (msec%1e3)*1e6).UTC()
}

func TestConvert(t *testing.T) {
	cs := []struct {
		w any
		s any
		b reflect.Type
	}{
		{"1.123", 1.123, TypeString},
		{true, "true", TypeBool},
		{int(0776), "0776", TypeInt},
		{int32(0777), "0777", TypeInt32},
		{int64(779), "779", TypeInt64},
		{time.Second, "1s", TypeDuration},
		{time.Second, 1000000000, TypeDuration},
		{utcMilli(0).Add(time.Second), "1970-01-01T00:00:01Z", TypeTime},
		{utcMilli(0).Add(time.Second), 1000, TypeTime},
		{int(0), nil, TypeInt},
		{int32(0), "", TypeInt32},
	}
	for i, c := range cs {
		a, err := Convert(c.s, c.b)
		if err != nil {
			t.Errorf("[%d] Convert(%v, %q) Failed: %v, want %v", i, c.s, c.b, err, c.w)
		}
		if a != c.w {
			t.Errorf("[%d] Convert(%v, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}
