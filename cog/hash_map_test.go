package cog

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestHashMapInterface(t *testing.T) {
	var m Map[int, int] = NewHashMap[int, int]()
	if m == nil {
		t.Error("HashMap is not a Map")
	}
}

func TestHashMapAsHashMap(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}
	hm := AsHashMap(m)

	if hm.Len() != len(m) {
		t.Fatal("hm.Len() != len(m)")
	}
	hm.Clear()

	if hm.Len() != 0 {
		t.Fatal("hm.Len() != 0")
	}
	if hm.Len() != len(m) {
		t.Fatal("hm.Len() != len(m)")
	}
}

func TestHashMapSet(t *testing.T) {
	m := NewHashMap[int, string]()
	m.Set(5, "e")
	m.Set(6, "f")
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(4, "d")
	m.Set(1, "x")
	m.Set(2, "b")
	m.Set(1, "a") //overwrite

	if av := m.Len(); av != 7 {
		t.Errorf("Got %v expected %v", av, 7)
	}
	if av, ev := m.Keys(), []int{1, 2, 3, 4, 5, 6, 7}; !testHashMapSame(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := m.Values(), []string{"a", "b", "c", "d", "e", "f", "g"}; !testHashMapSame(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}

	// key,ev,expectedFound
	tests1 := []struct {
		k int
		v string
		f bool
	}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests1 {
		// retrievals
		av, actualFound := m.Get(test.k)
		if av != test.v || actualFound != test.f {
			t.Errorf("Got %v expected %v", av, test.v)
		}
	}
}

func TestHashMapRemove(t *testing.T) {
	m := NewHashMap[int, string]()
	m.Set(5, "e")
	m.Set(6, "f")
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(4, "d")
	m.Set(1, "x")
	m.Set(2, "b")
	m.Set(1, "a") //overwrite

	m.Delete(5)
	m.Delete(6)
	m.Delete(7)
	m.Delete(8)
	m.Delete(5)

	if av, ev := m.Keys(), []int{1, 2, 3, 4}; !testHashMapSame(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}

	if av, ev := m.Values(), []string{"a", "b", "c", "d"}; !testHashMapSame(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av := m.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 4)
	}

	tests2 := []struct {
		k int
		v string
		f bool
	}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "", false},
		{6, "", false},
		{7, "", false},
		{8, "", false},
	}

	for _, test := range tests2 {
		av, actualFound := m.Get(test.k)
		if av != test.v || actualFound != test.f {
			t.Errorf("Got %v expected %v", av, test.v)
		}
	}

	m.Delete(1)
	m.Delete(4)
	m.Delete(2)
	m.Delete(3)
	m.Delete(2)
	m.Delete(2)

	if av, ev := fmt.Sprintf("%v", m.Keys()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%s", m.Values()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av := m.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
	if av := m.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
}

func TestHashMapJSON(t *testing.T) {
	cs := []struct {
		s string
		w *HashMap[string, int]
	}{
		{`{}`, NewHashMap[string, int]()},
		{`{"a":1,"b":2,"c":3}`, NewHashMap([]P[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}...)},
	}

	for i, c := range cs {
		a := NewHashMap[string, int]()
		err := json.Unmarshal(([]byte)(c.s), &a)
		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%q) = %v", i, c.s, err)
		}
		if !testHashMapSame(a.Values(), c.w.Values()) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %v", i, c.s, a.Values(), c.w.Values())
		}
		if !testHashMapSame(a.Keys(), c.w.Keys()) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %v", i, c.s, a.Keys(), c.w.Keys())
		}

		bs, err := json.Marshal(a)
		if err != nil {
			t.Errorf("[%d] json.Marshal(%v) = %v", i, a, err)
		}

		m := make(map[string]int)
		err = json.Unmarshal(bs, &m)
		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%q) = %v", i, bs, err)
		}

		ks := make([]string, 0)
		vs := make([]int, 0)
		for k, v := range m {
			ks = append(ks, k)
			vs = append(vs, v)
		}

		if !testHashMapSame(ks, a.Keys()) || !testHashMapSame(vs, a.Values()) {
			t.Errorf("[%d] json.Marshal(%q) = %q, want %q", i, a.String(), string(bs), c.s)
		}
	}
}

func testHashMapSame[T any](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, av := range a {
		found := false
		for _, bv := range b {
			if reflect.DeepEqual(av, bv) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
