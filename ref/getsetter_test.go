package ref

import (
	"reflect"
	"testing"
)

type testObject struct {
	Int      int
	String   string
	Bool     bool
	StrSlice []string
}

func TestSetGetProperty(t *testing.T) {
	cs := []struct {
		o *testObject
		k string
		v any
		w any
		g func(o *testObject) any
	}{
		{&testObject{}, "Int", "1", 1, func(o *testObject) any { return o.Int }},
		{&testObject{}, "String", true, "true", func(o *testObject) any { return o.String }},
		{&testObject{}, "StrSlice", nil, []string{}, func(o *testObject) any { return o.StrSlice }},
	}

	for i, c := range cs {
		err := SetProperty(c.o, c.k, c.v)
		if err != nil {
			t.Errorf("[%d] SetProperty(%s, %v) = %v", i, c.k, c.v, err)
			continue
		}

		a := c.g(c.o)
		if !reflect.DeepEqual(c.w, a) {
			t.Errorf("[%d] SetProperty(%s, %v) get %v, want %v", i, c.k, c.v, a, c.w)
			continue
		}

		g, err := GetProperty(c.o, c.k)
		if err != nil {
			t.Errorf("[%d] GetProperty(%s) = %v", i, c.k, err)
			continue
		}
		if !reflect.DeepEqual(c.w, g) {
			t.Errorf("[%d] GetProperty(%s) = %v, want %v", i, c.k, a, c.w)
			continue
		}
	}
}
