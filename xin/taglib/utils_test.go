package taglib

import (
	"sort"
	"strings"
	"testing"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/arraylist"
	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/str"
)

type teststrstrmap map[string]string
type testintstrmap map[int]string
type teststrs []string
type testints []int

func TestAsListStr(t *testing.T) {
	cs := []struct {
		dict bool
		sort bool
		list any
	}{
		{true, false, linkedhashmap.NewLinkedHashMap(cog.KV("1", "a"), cog.KV("2", "b"))},
		{true, true, map[string]string{"1": "a", "2": "b"}},
		{true, true, teststrstrmap{"1": "a", "2": "b"}},
		{false, false, arraylist.NewArrayList("1", "2")},
		{false, false, []string{"1", "2"}},
		{false, false, teststrs{"1", "2"}},
	}

	cts := []struct {
		k any
		w bool
	}{
		{"1", true},
		{"-1", false},
		{1, true},
		{-1, false},
	}

	for i, c := range cs {
		l := AsList(c.list)

		for _, ct := range cts {
			if _, ok := l.Get(ct.k); ok != ct.w {
				t.Errorf("[%d] %T Get(%v) = %v, want %v", i, c.list, ct.k, ok, ct.w)
			}
		}

		var ks []string
		var vs []string
		l.Each(func(k any, v string) bool {
			ks = append(ks, toString(k))
			vs = append(vs, v)
			return true
		})

		wks := "1,2"
		if c.sort {
			sort.Strings(ks)
			sort.Strings(vs)
		}
		aks := strings.Join(ks, ",")
		if wks != aks {
			t.Errorf("[%d] %T Each(keys) = %v, want %v", i, c.list, aks, wks)
		}

		wvs := str.If(c.dict, "a,b", wks)
		avs := strings.Join(vs, ",")
		if wvs != avs {
			t.Errorf("[%d] %T Each(values) = %v, want %v", i, c.list, avs, wvs)
		}
	}
}

func TestAsListInt(t *testing.T) {
	cs := []struct {
		dict bool
		sort bool
		list any
	}{
		{true, false, linkedhashmap.NewLinkedHashMap(cog.KV(1, "a"), cog.KV(2, "b"))},
		{true, true, map[int]string{1: "a", 2: "b"}},
		{true, true, testintstrmap{1: "a", 2: "b"}},
		{false, false, arraylist.NewArrayList(1, 2)},
		{false, false, []int{1, 2}},
		{false, false, testints{1, 2}},
	}

	cts := []struct {
		k any
		w bool
	}{
		{"1", true},
		{"-1", false},
		{1, true},
		{-1, false},
	}

	for i, c := range cs {
		l := AsList(c.list)

		for _, ct := range cts {
			if _, ok := l.Get(ct.k); ok != ct.w {
				t.Errorf("[%d] %T Get(%v) = %v, want %v", i, c.list, ct.k, ok, ct.w)
			}
		}

		var ks []string
		var vs []string
		l.Each(func(k any, v string) bool {
			ks = append(ks, toString(k))
			vs = append(vs, v)
			return true
		})

		wks := "1,2"
		if c.sort {
			sort.Strings(ks)
			sort.Strings(vs)
		}
		aks := strings.Join(ks, ",")
		if wks != aks {
			t.Errorf("[%d] %T Each(keys) = %v, want %v", i, c.list, aks, wks)
		}

		wvs := str.If(c.dict, "a,b", wks)
		avs := strings.Join(vs, ",")
		if wvs != avs {
			t.Errorf("[%d] %T Each(values) = %v, want %v", i, c.list, avs, wvs)
		}
	}
}

func TestAsValuesStr(t *testing.T) {
	cs := []any{
		arraylist.NewArrayList("1", "2"),
		[]string{"1", "2"},
		teststrs{"1", "2"},
	}

	cts := []struct {
		k any
		w bool
	}{
		{"1", true},
		{"-1", false},
		{1, true},
		{-1, false},
	}

	for i, c := range cs {
		vs := AsValues(c)

		for _, ct := range cts {
			if ok := vs.Contains(ct.k); ok != ct.w {
				t.Errorf("[%d] %T Contains(%v) = %v, want %v", i, c, ct.k, ok, ct.w)
			}
		}

		var as []string
		vs.Each(func(a any) bool {
			as = append(as, toString(a))
			return true
		})

		wes := "1,2"
		aes := strings.Join(as, ",")
		if wes != aes {
			t.Errorf("[%d] %T Each() = %v, want %v", i, c, aes, wes)
		}
	}
}

func TestAsValuesInt(t *testing.T) {
	cs := []any{
		arraylist.NewArrayList(1, 2),
		[]int{1, 2},
		testints{1, 2},
	}

	cts := []struct {
		k any
		w bool
	}{
		{"1", true},
		{"-1", false},
		{1, true},
		{-1, false},
	}

	for i, c := range cs {
		vs := AsValues(c)

		for _, ct := range cts {
			if ok := vs.Contains(ct.k); ok != ct.w {
				t.Errorf("[%d] %T Contains(%v) = %v, want %v", i, c, ct.k, ok, ct.w)
			}
		}

		var as []string
		vs.Each(func(a any) bool {
			as = append(as, toString(a))
			return true
		})

		wes := "1,2"
		aes := strings.Join(as, ",")
		if wes != aes {
			t.Errorf("[%d] %T Each() = %v, want %v", i, c, aes, wes)
		}
	}
}

func TestAsValuesAny(t *testing.T) {
	cs := []any{
		arraylist.NewArrayList[any](1, "2"),
		[]any{1, "2"},
	}

	cts := []struct {
		k any
		w bool
	}{
		{"1", false},
		{"-1", false},
		{1, true},
		{-1, false},
	}

	for i, c := range cs {
		vs := AsValues(c)

		for _, ct := range cts {
			if ok := vs.Contains(ct.k); ok != ct.w {
				t.Errorf("[%d] %T Contains(%v) = %v, want %v", i, c, ct.k, ok, ct.w)
			}
		}

		var as []string
		vs.Each(func(a any) bool {
			as = append(as, toString(a))
			return true
		})

		wes := "1,2"
		aes := strings.Join(as, ",")
		if wes != aes {
			t.Errorf("[%d] %T Each() = %v, want %v", i, c, aes, wes)
		}
	}
}

func TestToValues(t *testing.T) {
	type C struct {
		k any
		w bool
	}

	cs := []struct {
		v  any
		cs []C
	}{
		{"1", []C{{"1", true}, {1, true}, {-1, false}}},
		{2, []C{{"2", true}, {2, true}, {-1, false}}},
	}

	for i, c := range cs {
		vs := ToValues(c.v)

		for _, ct := range c.cs {
			if ok := vs.Contains(ct.k); ok != ct.w {
				t.Errorf("[%d] %T Contains(%v) = %v, want %v", i, c, ct.k, ok, ct.w)
			}
		}
	}
}
