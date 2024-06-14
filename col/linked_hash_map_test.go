package col

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

func TestLinkedHashMapInterface(t *testing.T) {
	var _ Map = NewLinkedHashMap()
	var _ IterableMap = NewLinkedHashMap()
}

func TestLinkedHashMapBasicFeatures(t *testing.T) {
	lm := NewLinkedHashMap()

	n := 100
	// set(i, 2 * i)
	for i := 0; i < n; i++ {
		ov, ok := any(nil), false

		assertLenEqual("TestLinkedHashMapBasicFeatures", t, lm, i)
		if i%2 == 0 {
			ov, ok = lm.Set(i, 2*i)
		} else {
			ov, ok = lm.SetIfAbsent(i, 2*i)
		}
		assertLenEqual("TestLinkedHashMapBasicFeatures", t, lm, i+1)

		var w any
		if ov != w {
			t.Errorf("[%d] Set() val = %v, want %v", i, ov, w)
		}
		w = false
		if ok != w {
			t.Errorf("[%d] Set() ok = %v, want %v", i, ok, w)
		}

		ov, ok = lm.SetIfAbsent(i, 3*i)
		w = 2 * i
		if ov != w {
			t.Errorf("[%d] SetIfAbsent() val = %v, want %v", i, ov, w)
		}
		w = true
		if ok != w {
			t.Errorf("[%d] SetIfAbsent() ok = %v, want %v", i, ok, w)
		}
	}

	// get what we just set
	for i := 0; i < n; i++ {
		ov, ok := lm.Get(i)

		var w any
		w = 2 * i
		if ov != w {
			t.Errorf("[%d] get val = %v, want %v", i, ov, w)
		}
		w = true
		if ok != w {
			t.Errorf("[%d] get ok = %v, want %v", i, ok, w)
		}
	}

	// get items of what we just set
	for i := 0; i < n; i++ {
		w := 2 * i

		v, ok := lm.Get(i)

		if v != w || !ok {
			t.Errorf("lm[%d] = %v, want %v", i, v, w)
		}
	}

	// items
	mis := lm.Items()
	if n != len(mis) {
		t.Errorf("len(mis) = %v, want %v", len(mis), n)
	}
	for i := 0; i < n; i++ {
		if i != mis[i].Key() {
			t.Errorf("mis[%d].key = %v, want %v", i, mis[i].Key(), i)
		}
		if i*2 != mis[i].Value() {
			t.Errorf("mis[%d].value = %v, want %v", i, mis[i].Value(), i*2)
		}
	}

	// keys
	ks := make([]any, n)
	for i := 0; i < n; i++ {
		ks[i] = i
	}
	if !reflect.DeepEqual(ks, lm.Keys()) {
		t.Errorf("lm.Keys() = %v, want %v", lm.Keys(), ks)
	}

	// values
	vs := make([]any, n)
	for i := 0; i < n; i++ {
		vs[i] = i * 2
	}
	if !reflect.DeepEqual(vs, lm.Values()) {
		t.Errorf("lm.Values() = %v, want %v", lm.Values(), vs)
	}

	// entries
	ps := make([]P, n)
	for i := 0; i < n; i++ {
		ps[i] = P{i, i * 2}
	}
	if !reflect.DeepEqual(ps, lm.Entries()) {
		t.Errorf("lm.Entries() = %v, want %v", lm.Entries(), vs)
	}

	// forward iteration
	i := 0
	for it := lm.Iterator(); it.Next(); {
		if i != it.Key() {
			t.Errorf("[%d] it.Key() = %v, want %v", i, it.Key(), i)
		}
		if i*2 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), i*2)
		}
		i++
	}

	// backward iteration
	i = n - 1
	for it := lm.Iterator(); it.Prev(); {
		if i != it.Key() {
			t.Errorf("[%d] it.Key() = %v, want %v", i, it.Key(), i)
		}
		if i*2 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), i*2)
		}
		i--
	}

	// forward iteration starting from known key
	i = 42
	for it := lm.IteratorOf(i); it.Next(); {
		i++
		if i != it.Key() {
			t.Errorf("[%d] it.Key() = %v, want %v", i, it.Key(), i)
		}
		if i*2 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), i*2)
		}
	}

	// double values for items with even keys
	for j := 0; j < n/2; j++ {
		i = 2 * j
		ov, ok := lm.Set(i, 4*i)

		if 2*i != ov {
			t.Errorf("[%d] set val = %v, want %v", i, ov, 2*i)
		}
		if !ok {
			t.Errorf("[%d] set ok = false, want true", i)
		}
	}

	// and delete itmes with odd keys
	for j := 0; j < n/2; j++ {
		i = 2*j + 1
		assertLenEqual("TestLinkedHashMapBasicFeatures", t, lm, n-j)
		lm.Remove(i)
		assertLenEqual("TestLinkedHashMapBasicFeatures", t, lm, n-j-1)

		// deleting again shouldn't change anything
		lm.Removes(i)
		assertLenEqual("TestLinkedHashMapBasicFeatures", t, lm, n-j-1)
	}

	// get the whole range
	for j := 0; j < n/2; j++ {
		i = 2 * j
		ov, ok := lm.Get(i)
		if 4*i != ov {
			t.Errorf("[%d] get val = %v, want %v", i, ov, 4*i)
		}
		if !ok {
			t.Errorf("[%d] gel ok = %v, want %v", i, true, false)
		}

		i = 2*j + 1
		ov, ok = lm.Get(i)
		if nil != ov {
			t.Errorf("[%d] gel val = %v, want %v", i, ov, nil)
		}
		if ok {
			t.Errorf("[%d] gel ok = %v, want %v", i, ok, false)
		}
	}

	// check iterations again
	i = 0
	for it := lm.Iterator(); it.Next(); {
		if i != it.Key() {
			t.Errorf("[%d] it.Key() = %v, want %v", i, it.Key(), i)
		}
		if i*4 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), i*4)
		}
		i += 2
	}
	i = 2 * ((n - 1) / 2)
	for it := lm.Iterator(); it.Prev(); {
		if i != it.Key() {
			t.Errorf("[%d] it.Key() = %v, want %v", i, it.Key(), i)
		}
		if i*4 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), i*4)
		}
		i -= 2
	}
}

func TestLinkedHashMapUpdatingDoesntChangePairsOrder(t *testing.T) {
	lm := NewLinkedHashMap([]P{{"foo", "bar"}, {12, 28}, {78, 100}, {"bar", "baz"}}...)

	ov, ok := lm.Set(78, 102)
	if ov != 100 {
		t.Errorf("lm.Set(78, 102) = %v, want %v", ov, 100)
	}
	if !ok {
		t.Errorf("lm.Set(78, 102) = %v, want %v", ok, true)
	}

	assertOrderedPairsEqual(t, lm,
		[]any{"foo", 12, 78, "bar"},
		[]any{"bar", 28, 102, "baz"})
}

func TestLinkedHashMapDeletingAndReinsertingChangesPairsOrder(t *testing.T) {
	lm := NewLinkedHashMap()
	lm.Set("foo", "bar")
	lm.Set(12, 28)
	lm.Set(78, 100)
	lm.Set("bar", "baz")

	// delete a item
	lm.Remove(78)

	// re-insert the same item
	lm.Set(78, 100)

	assertOrderedPairsEqual(t, lm,
		[]any{"foo", 12, "bar", 78},
		[]any{"bar", 28, "baz", 100})
}

func TestLinkedHashMapEmptyMapOperations(t *testing.T) {
	lm := NewLinkedHashMap()

	ov, ok := lm.Get("foo")
	if ov != nil {
		t.Errorf("lm.Get(foo) = %v, want %v", ov, nil)
	}
	if ok {
		t.Errorf("lm.Get(foo) = %v, want %v", ok, false)
	}

	lm.Remove("bar")
	assertLenEqual("TestLinkedHashMapEmptyMapOperations", t, lm, 0)

	fn := lm.Head()
	if fn != nil {
		t.Errorf("lm.Head() = %v, want %v", fn, nil)
	}

	bn := lm.Tail()
	if bn != nil {
		t.Errorf("lm.Tail() = %v, want %v", bn, nil)
	}
}

type dummyTestStruct struct {
	value string
}

func TestLinkedHashMapPackUnpackStructs(t *testing.T) {
	lm := NewLinkedHashMap()
	lm.Set("foo", dummyTestStruct{"foo!"})
	lm.Set("bar", dummyTestStruct{"bar!"})

	ov, ok := lm.Get("foo")
	if !ok {
		t.Fatalf(`lm.Get("foo") = %v`, ok)
	}
	if "foo!" != ov.(dummyTestStruct).value {
		t.Fatalf(`lm.Get("foo") = %v, want %v`, ov, "foo!")
	}

	ov, ok = lm.Set("bar", dummyTestStruct{"baz!"})
	if !ok {
		t.Fatalf(`lm.Set("bar") = %v`, ok)
	}
	if "bar!" != ov.(dummyTestStruct).value {
		t.Fatalf(`lm.Set("bar") = %v, want %v`, ov, "bar!")
	}

	ov, ok = lm.Get("bar")
	if !ok {
		t.Fatalf(`lm.Get("bar") = %v`, ok)
	}
	if "baz!" != ov.(dummyTestStruct).value {
		t.Fatalf(`lm.Get("bar") = %v, want %v`, ov, "baz!")
	}
}

func TestLinkedHashMapShuffle(t *testing.T) {
	ranLen := 100

	for _, n := range []int{0, 10, 20, 100, 1000, 10000} {
		t.Run(fmt.Sprintf("shuffle test with %d items", n), func(t *testing.T) {
			lm := NewLinkedHashMap()

			keys := make([]any, n)
			values := make([]any, n)

			for i := 0; i < n; i++ {
				// we prefix with the number to ensure that we don't get any duplicates
				keys[i] = fmt.Sprintf("%d_%s", i, randomHexString(t, ranLen))
				values[i] = randomHexString(t, ranLen)

				ov, ok := lm.Set(keys[i], values[i])
				if ok {
					t.Fatalf(`[%d] lm.Set(%v) = %v`, i, keys[i], ok)
				}
				if ov != nil {
					t.Fatalf(`[%d] lm.Set(%v) = %v`, i, keys[i], ov)
				}
			}

			assertOrderedPairsEqual(t, lm, keys, values)
		})
	}
}

func TestLinkedHashMapTemplateRange(t *testing.T) {
	lm := NewLinkedHashMap([]P{{"z", "Z"}, {"a", "A"}}...)
	tmpl, err := template.New("test").Parse("{{range $e := .lm.Items}}[ {{$e.Key}} = {{$e.Value}} ]{{end}}")
	if err != nil {
		t.Fatal(err.Error())
	}

	cm := map[string]any{
		"lm": lm,
	}
	sb := &strings.Builder{}
	err = tmpl.Execute(sb, cm)
	if err != nil {
		t.Fatal(err.Error())
	}

	a := sb.String()
	w := "[ z = Z ][ a = A ]"
	if w != a {
		t.Errorf("tmpl.Execute() = %q, want %q", a, w)
	}
}

/* Test helpers */
func assertOrderedPairsEqual(t *testing.T, lm *LinkedHashMap, eks, evs []any) {
	assertOrderedPairsEqualFromNewest(t, lm, eks, evs)
	assertOrderedPairsEqualFromOldest(t, lm, eks, evs)
}

func assertOrderedPairsEqualFromNewest(t *testing.T, lm *LinkedHashMap, eks, evs []any) {
	if len(eks) != len(evs) {
		t.Errorf("len(keys) %v != len(vals) %v", len(eks), len(evs))
		return
	}

	if len(eks) != lm.Len() {
		t.Errorf("len(keys) %v != lm.Len %v", len(eks), lm.Len())
		return
	}

	i := lm.Len() - 1
	for it := lm.Iterator(); it.Prev(); {
		if eks[i] != it.Key() {
			t.Errorf("[%d] key = %v, want %v", i, it.Key(), eks[i])
		}

		if evs[i] != it.Value() {
			t.Errorf("[%d] val = %v, want %v", i, it.Value(), evs[i])
		}
		i--
	}
}

func assertOrderedPairsEqualFromOldest(t *testing.T, lm *LinkedHashMap, eks, evs []any) {
	if len(eks) != len(evs) {
		t.Errorf("len(keys) %v != len(vals) %v", len(eks), len(evs))
		return
	}

	if len(eks) != lm.Len() {
		t.Errorf("len(keys) %v != lm.Len %v", len(eks), lm.Len())
		return
	}

	i := 0
	for it := lm.Iterator(); it.Next(); {
		if eks[i] != it.Key() {
			t.Errorf("[%d] key = %v, want %v", i, it.Key(), eks[i])
		}

		if evs[i] != it.Value() {
			t.Errorf("[%d] val = %v, want %v", i, it.Value(), evs[i])
		}
		i++
	}
}

func assertLenEqual(n string, t *testing.T, lm *LinkedHashMap, w int) {
	if lm.Len() != w {
		t.Fatalf("%s: lm.Len() != %v", n, w)
	}
}

func randomHexString(t *testing.T, length int) string {
	b := length / 2
	randBytes := make([]byte, b)

	if n, err := rand.Read(randBytes); err != nil || n != b {
		if err == nil {
			err = fmt.Errorf("only got %v random bytes, expected %v", n, b)
		}
		t.Fatal(err)
	}

	return hex.EncodeToString(randBytes)
}

func TestLinkedHashMapString(t *testing.T) {
	w := `{"1":1,"3":3,"2":2}`
	a := NewLinkedHashMap([]P{{"1", 1}, {"3", 3}, {"2", 2}}...).String()
	if w != a {
		t.Errorf("TestLinkedHashMapString = %v, want %v", a, w)
	}
}

//--------------------------------------------

func TestLinkedHashMapPut(t *testing.T) {
	m := NewLinkedHashMap()
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
	if av, ev := m.Keys(), []any{5, 6, 7, 3, 4, 1, 2}; !testLinkedHashMapSameValues(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := m.Values(), []any{"e", "f", "g", "c", "d", "a", "b"}; !testLinkedHashMapSameValues(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := m.Entries(), []P{{5, "e"}, {6, "f"}, {7, "g"}, {3, "c"}, {4, "d"}, {1, "a"}, {2, "b"}}; !testLinkedHashMapSameEntries(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}

	// key,ev,expectedFound
	tests1 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, nil, false},
	}

	for _, test := range tests1 {
		// retrievals
		av, actualFound := m.Get(test[0])
		if av != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", av, test[1])
		}
	}
}

func TestLinkedHashMapRemove(t *testing.T) {
	m := NewLinkedHashMap()
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

	if av, ev := m.Keys(), []any{3, 4, 1, 2}; !testLinkedHashMapSameValues(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := m.Values(), []any{"c", "d", "a", "b"}; !testLinkedHashMapSameValues(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av := m.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 4)
	}

	tests2 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, nil, false},
		{6, nil, false},
		{7, nil, false},
		{8, nil, false},
	}

	for _, test := range tests2 {
		av, actualFound := m.Get(test[0])
		if av != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", av, test[1])
		}
	}

	m.Remove(1)
	m.Removes(4, 2, 3, 2)
	m.Remove(2)

	if av, ev := fmt.Sprintf("%s", m.Keys()), "[]"; av != ev {
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

func testLinkedHashMapSameValues(a []any, b []any) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func testLinkedHashMapSameEntries(a []P, b []P) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Key != b[i].Key && a[i].Value != b[i].Value {
			return false
		}
	}

	return true
}

func TestLinkedHashMapEach(t *testing.T) {
	m := NewLinkedHashMap()
	m.Set("c", 1)
	m.Set("a", 2)
	m.Set("b", 3)
	count := 0
	m.Each(func(key any, value any) bool {
		count++
		if av, ev := count, value; av != ev {
			t.Errorf("Got %v expected %v", av, ev)
		}
		switch value {
		case 1:
			if av, ev := key, "c"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 2:
			if av, ev := key, "a"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 3:
			if av, ev := key, "b"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			t.Errorf("Too many")
		}
		return true
	})
}

func TestLinkedHashMapIteratorNextOnEmpty(t *testing.T) {
	m := NewLinkedHashMap()
	it := m.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func TestLinkedHashMapIteratorPrevOnEmpty(t *testing.T) {
	m := NewLinkedHashMap()
	it := m.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func TestLinkedHashMapIteratorNext(t *testing.T) {
	m := NewLinkedHashMap()
	m.Set("c", 1)
	m.Set("a", 2)
	m.Set("b", 3)

	it := m.Iterator()
	count := 0
	for it.Next() {
		count++
		key := it.Key()
		value := it.Value()
		switch key {
		case "c":
			if av, ev := value, 1; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case "a":
			if av, ev := value, 2; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case "b":
			if av, ev := value, 3; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			t.Errorf("Too many")
		}
		if av, ev := value, count; av != ev {
			t.Errorf("Got %v expected %v", av, ev)
		}
	}
	if av, ev := count, 3; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashMapIteratorPrev(t *testing.T) {
	m := NewLinkedHashMap()
	m.Set("c", 1)
	m.Set("a", 2)
	m.Set("b", 3)

	it := m.Iterator()
	countDown := m.Len()
	for it.Prev() {
		key := it.Key()
		value := it.Value()
		switch key {
		case "c":
			if av, ev := value, 1; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case "a":
			if av, ev := value, 2; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case "b":
			if av, ev := value, 3; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			t.Errorf("Too many")
		}
		if av, ev := value, countDown; av != ev {
			t.Errorf("Got %v expected %v", av, ev)
		}
		countDown--
	}
	if av, ev := countDown, 0; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashMapIteratorReset(t *testing.T) {
	m := NewLinkedHashMap()
	it := m.Iterator()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	for it.Next() {
	}
	it.Reset()
	it.Next()
	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}

	it.Reset()
	it.Prev()
	if key, value := it.Key(), it.Value(); key != 2 || value != "b" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 2, "b")
	}
}

func assertLinkedHashMapIteratorRemove(t *testing.T, i int, it Iterator2, w *LinkedHashMap) int {
	it.Remove()

	k := it.Key()
	w.Remove(k)

	it.SetValue(9999)

	m := it.(*linkedHashMapIterator).lmap
	if m.Contains(k) {
		t.Fatalf("[%d] w.Contains(%v) = true", i, k)
	}

	if m.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, k, m.String(), w.String())
	}

	return k.(int)
}

func TestLinkedHashMapIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		m := NewLinkedHashMap()
		w := NewLinkedHashMap()

		for n := 0; n < i; n++ {
			m.Set(n, -n)
			w.Set(n, -n)
		}

		it := m.Iterator()

		// remove nothing
		it.Remove()
		w.Remove(it.Key())
		it.SetValue(9999)
		if m.Len() != i {
			t.Fatalf("[%d] m.Len() == %v, want %v", i, m.Len(), i)
		}

		// remove middle
		for j := 0; j <= m.Len()/2; j++ {
			it.Next()
		}

		v := assertLinkedHashMapIteratorRemove(t, i, it, w)

		it.Next()
		if v+1 != it.Key() {
			t.Fatalf("[%d] it.Key() = %v, want %v", i, it.Key(), v+1)
		}
		assertLinkedHashMapIteratorRemove(t, i, it, w)

		it.Prev()
		if v-1 != it.Key() {
			t.Fatalf("[%d] it.Key() = %v, want %v", i, it.Key(), v-1)
		}
		assertLinkedHashMapIteratorRemove(t, i, it, w)

		// remove first
		for it.Prev() {
		}
		assertLinkedHashMapIteratorRemove(t, i, it, w)

		// remove last
		for it.Next() {
		}
		assertLinkedHashMapIteratorRemove(t, i, it, w)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertLinkedHashMapIteratorRemove(t, i, it, w)
			}
		} else {
			for it.Next() {
				assertLinkedHashMapIteratorRemove(t, i, it, w)
			}
		}
		if !m.IsEmpty() {
			t.Fatalf("[%d] m.IsEmpty() = true", i)
		}
	}
}

func TestLinkedHashMapIteratorSetValue(t *testing.T) {
	m := NewLinkedHashMap()
	for i := 1; i <= 100; i++ {
		m.Set(i, i)
	}

	// forward
	for it := m.Iterator(); it.Next(); {
		it.SetValue(it.Value().(int) + 100)
	}
	for i := 1; i <= m.Len(); i++ {
		v, _ := m.Get(i)
		w := i + 100
		if v != w {
			t.Fatalf("Hash[%d] = %v, want %v", i, v, w)
		}
	}

	// backward
	for it := m.Iterator(); it.Prev(); {
		it.SetValue(it.Value().(int) + 100)
	}
	for i := 1; i <= m.Len(); i++ {
		v, _ := m.Get(i)
		w := i + 200
		if v != w {
			t.Fatalf("Hash[%d] = %v, want %v", i, v, w)
		}
	}
}

/*----------- JOSN Test -----------------*/
func TestLinkedHashMapMarshal(t *testing.T) {
	lm := NewLinkedHashMap()
	lm.Set("a", 34)
	lm.Set("b", []int{3, 4, 5})
	b, err := json.Marshal(lm)
	if err != nil {
		t.Fatalf("Marshal LinkedHashMap: %v", err)
	}
	// fmt.Printf("%q\n", b)
	const expected = "{\"a\":34,\"b\":[3,4,5]}"
	if !bytes.Equal(b, []byte(expected)) {
		t.Errorf("Marshal LinkedHashMap: %q not equal to expected %q", b, expected)
	}
}

func TestLinkedHashMapUnmarshalFromInvalid(t *testing.T) {
	lm := NewLinkedHashMap()

	lm.Set("m", math.NaN())
	b, err := json.Marshal(lm)
	if err == nil {
		t.Fatal("Unmarshal LinkedHashMap: expecting error:", b, err)
	}
	// fmt.Println(lm, b, err)
	lm.Remove("m")

	err = json.Unmarshal([]byte("[]"), lm)
	if err == nil {
		t.Fatal("Unmarshal LinkedHashMap: expecting error")
	}

	err = json.Unmarshal([]byte("["), lm)
	if err == nil {
		t.Fatal("Unmarshal LinkedHashMap: expecting error:", lm)
	}

	err = lm.UnmarshalJSON([]byte(nil))
	if err == nil {
		t.Fatal("Unmarshal LinkedHashMap: expecting error:", lm)
	}

	err = lm.UnmarshalJSON([]byte("{}3"))
	if err == nil {
		t.Fatal("Unmarshal LinkedHashMap: expecting error:", lm)
	}

	err = lm.UnmarshalJSON([]byte("{"))
	if err == nil {
		t.Fatal("Unmarshal LinkedHashMap: expecting error:", lm)
	}

	err = lm.UnmarshalJSON([]byte("{]"))
	if err == nil {
		t.Fatal("Unmarshal LinkedHashMap: expecting error:", lm)
	}

	err = lm.UnmarshalJSON([]byte(`{"a": 3, "b": [{`))
	if err == nil {
		t.Fatal("Unmarshal LinkedHashMap: expecting error:", lm)
	}

	err = lm.UnmarshalJSON([]byte(`{"a": 3, "b": [}`))
	if err == nil {
		t.Fatal("Unmarshal LinkedHashMap: expecting error:", lm)
	}
	// fmt.Println("error:", lm, err)
}

func TestLinkedHashMapUnmarshal(t *testing.T) {
	var (
		data  = []byte(`{"as":"AS15169 Google Inc.","city":"Mountain View","country":"United States","countryCode":"US","isp":"Google Cloud","lat":37.4192,"lon":-122.0574,"org":"Google Cloud","query":"35.192.25.53","region":"CA","regionName":"California","status":"success","timezone":"America/Los_Angeles","zip":"94043"}`)
		pairs = []P{
			{"as", "AS15169 Google Inc."},
			{"city", "Mountain View"},
			{"country", "United States"},
			{"countryCode", "US"},
			{"isp", "Google Cloud"},
			{"lat", 37.4192},
			{"lon", -122.0574},
			{"org", "Google Cloud"},
			{"query", "35.192.25.53"},
			{"region", "CA"},
			{"regionName", "California"},
			{"status", "success"},
			{"timezone", "America/Los_Angeles"},
			{"zip", "94043"},
		}
		obj = NewLinkedHashMap(pairs...)
	)

	lm := NewLinkedHashMap()
	err := json.Unmarshal(data, lm)
	if err != nil {
		t.Fatalf("Unmarshal LinkedHashMap: %v", err)
	}

	// check by Has and GetValue
	for _, p := range pairs {
		k := p.Key
		v := p.Value

		if !lm.Contains(k) {
			t.Fatalf("expect key %q exists in Unmarshaled LinkedHashMap", k)
		}
		value, ok := lm.Get(k)
		if !ok || value != v {
			t.Fatalf("expect for key %q: the value %v should equal to %v, in Unmarshaled LinkedHashMap", k, value, v)
		}
	}

	b, err := json.MarshalIndent(lm, "", "  ")
	if err != nil {
		t.Fatalf("Unmarshal LinkedHashMap: %v", err)
	}
	const expected = `{
  "as": "AS15169 Google Inc.",
  "city": "Mountain View",
  "country": "United States",
  "countryCode": "US",
  "isp": "Google Cloud",
  "lat": 37.4192,
  "lon": -122.0574,
  "org": "Google Cloud",
  "query": "35.192.25.53",
  "region": "CA",
  "regionName": "California",
  "status": "success",
  "timezone": "America/Los_Angeles",
  "zip": "94043"
}`
	if !bytes.Equal(b, []byte(expected)) {
		t.Fatalf("Unmarshal LinkedHashMap marshal indent from %#v not equal to expected: %q\n", lm, expected)
	}

	if !reflect.DeepEqual(lm, obj) {
		t.Fatalf("Unmarshal LinkedHashMap not deeply equal: %#v %#v", lm, obj)
	}

	val, ok := lm.Get("org")
	if !ok {
		t.Fatalf("org should exist")
	}
	lm.Remove("org")
	lm.Set("org", val)

	b, err = json.MarshalIndent(lm, "", "  ")
	// fmt.Println("after delete", lm, string(b), err)
	if err != nil {
		t.Fatalf("Unmarshal LinkedHashMap: %v", err)
	}
	const expected2 = `{
  "as": "AS15169 Google Inc.",
  "city": "Mountain View",
  "country": "United States",
  "countryCode": "US",
  "isp": "Google Cloud",
  "lat": 37.4192,
  "lon": -122.0574,
  "query": "35.192.25.53",
  "region": "CA",
  "regionName": "California",
  "status": "success",
  "timezone": "America/Los_Angeles",
  "zip": "94043",
  "org": "Google Cloud"
}`
	if !bytes.Equal(b, []byte(expected2)) {
		t.Fatalf("Unmarshal LinkedHashMap marshal indent from %#v not equal to expected: %s\n", lm, expected2)
	}
}

func TestLinkedHashMapUnmarshalNested(t *testing.T) {
	var (
		data = []byte(`{"a": true, "b": [3, 4, { "b": "3", "d": [] }]}`)
		obj  = NewLinkedHashMap(
			P{"a", true},
			P{"b", JSONArray{float64(3), float64(4), JSONObject{"b": "3", "d": JSONArray{}}}},
		)
	)

	lm := NewLinkedHashMap()
	err := json.Unmarshal(data, lm)
	if err != nil {
		t.Fatalf("Unmarshal LinkedHashMap: %v", err)
	}

	if !reflect.DeepEqual(lm, obj) {
		t.Fatalf("Unmarshal LinkedHashMap not deeply equal: %#v expected %#v", lm, obj)
	}
}

func TestLinkedHashMapUnmarshals(t *testing.T) {
	var unmarshalTests = []struct {
		in  string
		out any
		err any
	}{
		{in: "{}", out: NewLinkedHashMap()},
		{in: `{"a": 3}`, out: NewLinkedHashMap(P{"a", float64(3)})},
		{in: `{"a": 3, "b": true}`, out: NewLinkedHashMap([]P{{"a", float64(3)}, {"b", true}}...)},
		{in: `{"a": 3, "b": true, "c": null}`, out: NewLinkedHashMap([]P{{"a", float64(3)}, {"b", true}, {"c", nil}}...)},
		{in: `{"a": 3, "c": null, "d": []}`, out: NewLinkedHashMap([]P{{"a", float64(3)}, {"c", nil}, {"d", JSONArray{}}}...)},
		{in: `{"a": 3, "c": null, "d": [3,4,true]}`, out: NewLinkedHashMap([]P{
			{"a", float64(3)}, {"c", nil}, {"d", JSONArray{
				float64(3), float64(4), true,
			}}}...)},
		{in: `{"a": 3, "c": null, "d": [3,4,true, { "inner": "abc" }]}`, out: NewLinkedHashMap([]P{
			{"a", float64(3)}, {"c", nil}, {"d", JSONArray([]any{
				float64(3), float64(4), true, JSONObject{"inner": "abc"},
			})}}...)},
	}

	for i, tt := range unmarshalTests {
		in := []byte(tt.in)

		v := NewLinkedHashMap()
		dec := json.NewDecoder(bytes.NewReader(in))
		if err := dec.Decode(v); !reflect.DeepEqual(err, tt.err) {
			t.Errorf("#%d: %v, want %v", i, err, tt.err)
			continue
		} else if err != nil {
			continue
		}

		if !reflect.DeepEqual(v, tt.out) {
			act, _ := json.Marshal(v)
			exp, _ := json.Marshal(tt.out)
			t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v\nact: %s\nexp: %s", i, v, tt.out, string(act), string(exp))
			continue
		}

		// Check round trip also decodes correctly.
		if tt.err == nil {
			enc, err := json.Marshal(v)
			if err != nil {
				t.Errorf("#%d: error re-marshaling: %v", i, err)
				continue
			}

			vv := NewLinkedHashMap()
			dec = json.NewDecoder(bytes.NewReader(enc))
			if err := dec.Decode(vv); err != nil {
				t.Errorf("#%d: error re-unmarshaling %#q: %v", i, enc, err)
				continue
			}
			if !reflect.DeepEqual(v, vv) {
				t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, v, vv)
				t.Errorf("     In: %q", in)
				t.Errorf("Marshal: %q", enc)
				continue
			}
		}
	}
}
