package ref

import (
	"reflect"
	"testing"
)

type testInner struct {
	InInt int
	instr string
	inval string
}

func (ti *testInner) GetInstr() string {
	return ti.instr
}
func (ti *testInner) SetInstr(s string) {
	ti.instr = s
}
func (ti testInner) GetInval() string {
	return ti.inval
}
func (ti *testInner) SetInval(s string) {
	ti.inval = s
}

type testObject struct {
	testInner

	Int      int
	str      string
	val      string
	Bool     bool
	StrSlice []string
}

func (to *testObject) GetStr() string {
	return to.str
}
func (to *testObject) SetStr(s string) {
	to.str = s
}
func (to testObject) GetVal() string {
	return to.val
}
func (to *testObject) SetVal(s string) {
	to.val = s
}

func TestSetGetProperty(t *testing.T) {
	cs := []struct {
		o *testObject
		k string
		v any
		w any
		g func(o *testObject) any
	}{
		{&testObject{}, "InInt", "11", 11, func(o *testObject) any { return o.InInt }},
		{&testObject{}, "instr", "s", "s", func(o *testObject) any { return o.instr }},
		{&testObject{}, "inval", "v", "v", func(o *testObject) any { return o.inval }},
		{&testObject{}, "Int", "1", 1, func(o *testObject) any { return o.Int }},
		{&testObject{}, "str", true, "true", func(o *testObject) any { return o.str }},
		{&testObject{}, "val", true, "true", func(o *testObject) any { return o.val }},
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

		if c.k != "instr" && c.k != "str" {
			g, err := GetProperty(*c.o, c.k)
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
}
