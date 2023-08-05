package ref

import (
	"reflect"
	"testing"
)

type testSetter struct {
	Int      int
	String   string
	Bool     bool
	StrSlice []string
}

func TestSetProperty(t *testing.T) {
	cs := []struct {
		o *testSetter
		k string
		v any
		w any
		g func(o *testSetter) any
	}{
		{&testSetter{}, "Int", "1", 1, func(o *testSetter) any { return o.Int }},
		{&testSetter{}, "String", true, "true", func(o *testSetter) any { return o.String }},
		{&testSetter{}, "StrSlice", nil, []string{}, func(o *testSetter) any { return o.StrSlice }},
	}

	for i, c := range cs {
		err := SetProperty(c.o, c.k, c.v)
		a := c.g(c.o)
		if err != nil || !reflect.DeepEqual(c.w, a) {
			t.Errorf("[%d] SetProperty(%s, %v) = %v, get %v, want %v", i, c.k, c.v, err, a, c.w)
		}
	}
}
