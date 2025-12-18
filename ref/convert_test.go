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

func TestCastTo(t *testing.T) {
	cs := []struct {
		w any
		s any
		b reflect.Type
	}{
		{"a", "a", nil},
		{[]byte{'s', '2', 'b'}, "s2b", TypeBytes},
		{"b2s", []byte{'b', '2', 's'}, TypeString},
		{"1.123", 1.123, TypeString},
		{true, "true", TypeBool},
		{int(-0776), "-0776", TypeInt},
		{int8(-8), "-8", TypeInt8},
		{int16(-16), "-16", TypeInt16},
		{int32(-32), "-32", TypeInt32},
		{int64(-64), "-64", TypeInt64},
		{uint(0776), "0776", TypeUint},
		{uint8(8), "8", TypeUint8},
		{uint16(16), "16", TypeUint16},
		{uint32(32), "32", TypeUint32},
		{uint64(64), "64", TypeUint64},
		{float32(0.32), "0.32", TypeFloat32},
		{float64(0.64), "0.64", TypeFloat64},
		{time.Second, "1s", TypeDuration},
		{time.Second, 1000000000, TypeDuration},
		{utcMilli(0).Add(time.Second), "1970-01-01T00:00:01Z", TypeTime},
		{utcMilli(0).Add(time.Second), 1000, TypeTime},
		{int(0), nil, TypeInt},
		{int32(0), "", TypeInt32},
		{[]string{"1", "2"}, []int{1, 2}, TypeStrings},
		{map[string]int{"1": 11, "2": 22}, map[int]string{1: "11", 2: "22"}, reflect.TypeFor[map[string]int]()},
	}
	for i, c := range cs {
		a, err := CastTo(c.s, c.b)
		if err != nil {
			t.Errorf("[%d] CastTo(%v, %q) Failed: %v, want %v", i, c.s, c.b, err, c.w)
		}
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] CastTo(%v, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}
