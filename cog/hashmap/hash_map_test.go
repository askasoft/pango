package hashmap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/askasoft/pango/cog"
)

func TestHashMapInterface(t *testing.T) {
	var _ cog.Map[int, int] = NewHashMap[int, int]()
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
	if av, ev := m.Keys(), []int{1, 2, 3, 4, 5, 6, 7}; !testHashMapSameValues(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := m.Values(), []string{"a", "b", "c", "d", "e", "f", "g"}; !testHashMapSameValues(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := m.Entries(), []cog.P[int, string]{
		cog.KV(1, "a"),
		cog.KV(2, "b"),
		cog.KV(3, "c"),
		cog.KV(4, "d"),
		cog.KV(5, "e"),
		cog.KV(6, "f"),
		cog.KV(7, "g")}; !testHashMapSameEntries(av, ev) {
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

	m.Remove(5)
	m.Removes(6, 7, 8)
	m.Remove(5)

	if av, ev := m.Keys(), []int{1, 2, 3, 4}; !testHashMapSameValues(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}

	if av, ev := m.Values(), []string{"a", "b", "c", "d"}; !testHashMapSameValues(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := m.Entries(), []cog.P[int, string]{
		cog.KV(1, "a"),
		cog.KV(2, "b"),
		cog.KV(3, "c"),
		cog.KV(4, "d")}; !testHashMapSameEntries(av, ev) {
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

	m.Remove(1)
	m.Removes(4, 2, 3, 2)
	m.Remove(2)

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
		{`{"a":1,"b":2,"c":3}`, NewHashMap([]cog.P[string, int]{cog.KV("a", 1), cog.KV("b", 2), cog.KV("c", 3)}...)},
	}

	for i, c := range cs {
		a := NewHashMap[string, int]()
		err := json.Unmarshal(([]byte)(c.s), &a)
		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%q) = %v", i, c.s, err)
		}
		if !testHashMapSameValues(a.Values(), c.w.Values()) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %v", i, c.s, a.Values(), c.w.Values())
		}
		if !testHashMapSameValues(a.Keys(), c.w.Keys()) {
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

		if !testHashMapSameValues(ks, a.Keys()) || !testHashMapSameValues(vs, a.Values()) {
			t.Errorf("[%d] json.Marshal(%q) = %q, want %q", i, a.String(), string(bs), c.s)
		}
	}
}

func testHashMapSameValues[T any](a []T, b []T) bool {
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

func testHashMapSameEntries[K any, V any](a []cog.P[K, V], b []cog.P[K, V]) bool {
	if len(a) != len(b) {
		return false
	}
	for _, ap := range a {
		found := false
		for _, bp := range b {
			if reflect.DeepEqual(ap.Key, ap.Key) && reflect.DeepEqual(ap.Value, bp.Value) {
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
