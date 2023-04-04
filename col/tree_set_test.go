package col

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/pandafw/pango/ars"
)

func TestTreeSetInterface(t *testing.T) {
	var s Set = NewTreeSet(CompareInt)
	if s == nil {
		t.Error("TreeSet is not a Set")
	}
}

func TestTreeSetNew(t *testing.T) {
	tset := NewTreeSet(CompareString)
	if av := tset.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	tset = NewTreeSet(CompareString, "1", "b")
	if av := tset.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	if av := tset.Head(); av != "1" {
		t.Errorf("Got %v expected %v", av, 1)
	}
	if av := tset.Tail(); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestTreeSetDebug(t *testing.T) {
	tset := NewTreeSet(CompareString)
	ev := "(empty)"
	av := tset.debug()
	if av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestTreeSetAdd(t *testing.T) {
	tset := NewTreeSet(CompareString)
	tset.Add("a")
	tset.Add("b", "c")
	if av := tset.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := tset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av := tset.Tail(); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestTreeSetClear(t *testing.T) {
	tset := NewTreeSet(CompareString)
	tset.Add("e", "f", "g", "a", "b", "c", "d")
	tset.Add("e", "f", "g", "a", "b", "c", "d")
	tset.Clear()
	if av := tset.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := tset.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestTreeSetContains(t *testing.T) {
	list := NewTreeSet(CompareInt)

	a := []T{}
	for i := 0; i < 100; i++ {
		a = append(a, i)
		list.Add(i)
	}
	a = append(a, 1000)

	for i := 0; i < 100; i++ {
		if !list.Contains(i) {
			t.Errorf("%d Contains() should return true", i)
		}
		if !list.Contains(a[0 : i+1]...) {
			t.Errorf("%d Contains(...) should return true", i)
		}
		if list.Contains(a...) {
			t.Errorf("%d Contains(...) should return false", i)
		}
		if !list.ContainsAll(AsArrayList(a[0 : i+1])) {
			t.Errorf("%d ContainsAll(...) should return true", i)
		}
		if list.ContainsAll(AsArrayList(a)) {
			t.Errorf("%d ContainsAll(...) should return false", i)
		}
	}

	list.Clear()
	if av := list.Contains("a"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.Contains("a", "b", "c"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestTreeSetRetain(t *testing.T) {
	for n := 0; n < 100; n++ {
		a := []T{}
		list := NewTreeSet(CompareInt)
		for i := 0; i < n; i++ {
			if i&1 == 0 {
				a = append(a, i)
			}
			list.Add(i)

			list.Retain(a...)
			vs := list.Values()
			if !reflect.DeepEqual(vs, a) {
				t.Fatalf("%d Retain() = %v, want %v", i, vs, a)
			}
		}

		{
			a = []T{}
			list.Retain()
			vs := list.Values()
			if len(vs) > 0 {
				t.Fatalf("%d Retain() = %v, want %v", n, vs, a)
			}
		}

		a = []T{}
		list.Clear()
		for i := 0; i < n; i++ {
			if i&1 == 0 {
				a = append(a, i)
			}
			list.Add(i)

			list.RetainAll(AsArrayList(a))
			vs := list.Values()
			if !reflect.DeepEqual(vs, a) {
				t.Fatalf("%d RetainAll() = %v, want %v", i, vs, a)
			}
		}

		{
			a = []T{}
			list.RetainAll(AsArrayList(a))
			vs := list.Values()
			if len(vs) > 0 {
				t.Fatalf("%d Retain() = %v, want %v", n, vs, a)
			}
		}
	}
}

func TestTreeSetValues(t *testing.T) {
	tset := NewTreeSet(CompareString)
	tset.Add("a", "a")
	tset.Add("b", "c", "b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", tset.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestTreeSetEach(t *testing.T) {
	tset := NewTreeSet(CompareString)
	tset.Add("a", "b", "c")
	tset.Add("a", "b", "c")
	index := 0
	tset.Each(func(value any) {
		switch index {
		case 0:
			if av, ev := value, "a"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 1:
			if av, ev := value, "b"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 2:
			if av, ev := value, "c"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			t.Errorf("Too many")
		}
		index++
	})
}

func TestTreeSetIteratorNextOnEmpty(t *testing.T) {
	tset := NewTreeSet(CompareString)
	it := tset.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty tset")
	}
}

func TestTreeSetIteratorNext(t *testing.T) {
	tset := NewTreeSet(CompareString)
	tset.Add("a", "b", "c")
	tset.Add("a", "b", "c")
	it := tset.Iterator()
	count := 0
	index := 0
	for it.Next() {
		count++
		value := it.Value()
		switch index {
		case 0:
			if av, ev := value, "a"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 1:
			if av, ev := value, "b"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 2:
			if av, ev := value, "c"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			t.Errorf("Too many")
		}
		index++
	}
	if av, ev := count, 3; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestTreeSetIteratorPrevOnEmpty(t *testing.T) {
	tset := NewTreeSet(CompareString)
	it := tset.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty tset")
	}
}

func TestTreeSetIteratorPrev(t *testing.T) {
	tset := NewTreeSet(CompareString)
	tset.Add("a", "b", "c")
	tset.Add("a", "b", "c")
	it := tset.Iterator()
	count := 0
	index := tset.Len()
	for it.Prev() {
		count++
		index--
		value := it.Value()
		switch index {
		case 0:
			if av, ev := value, "a"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 1:
			if av, ev := value, "b"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 2:
			if av, ev := value, "c"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			t.Errorf("Too many")
		}
	}
	if av, ev := count, 3; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestTreeSetIteratorReset(t *testing.T) {
	tset := NewTreeSet(CompareString)
	it := tset.Iterator()
	tset.Add("a", "b", "c")
	tset.Add("a", "b", "c")
	for it.Next() {
	}
	it.Reset()
	it.Next()
	if value := it.Value(); value != "a" {
		t.Errorf("Got %v expected %v", value, "a")
	}

	it.Reset()
	it.Prev()
	if value := it.Value(); value != "c" {
		t.Errorf("Got %v expected %v", value, "c")
	}
}

func assertTreeSetIteratorRemove(t *testing.T, i int, it Iterator, w *TreeSet) int {
	it.Remove()

	v := it.Value()
	w.Delete(v)

	it.SetValue(9999)

	tset := it.(*treeSetIterator).tree
	if tset.Contains(v) {
		t.Fatalf("[%d] tset.Contains(%v) = true", i, v)
	}

	if tset.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, tset.String(), w.String())
	}

	return v.(int)
}

func TestTreeSetIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		tset := NewTreeSet(CompareInt)
		wset := NewTreeSet(CompareInt)

		for n := 0; n < i; n++ {
			tset.Add(n)
			wset.Add(n)
		}

		it := tset.Iterator()

		it.Remove()
		if tset.Len() != i {
			t.Fatalf("[%d] tset.Len() == %v, want %v", i, tset.Len(), i)
		}

		// remove middle
		for j := 0; j <= tset.Len()/2; j++ {
			it.Next()
		}

		v := assertTreeSetIteratorRemove(t, i, it, wset)

		it.Next()
		if v+1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		assertTreeSetIteratorRemove(t, i, it, wset)

		it.Prev()
		if v-1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		assertTreeSetIteratorRemove(t, i, it, wset)

		// remove first
		for it.Prev() {
		}
		assertTreeSetIteratorRemove(t, i, it, wset)

		// remove last
		for it.Next() {
		}
		assertTreeSetIteratorRemove(t, i, it, wset)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertTreeSetIteratorRemove(t, i, it, wset)
			}
		} else {
			for it.Next() {
				assertTreeSetIteratorRemove(t, i, it, wset)
			}
		}
		if !tset.IsEmpty() {
			t.Fatalf("[%d] tset.IsEmpty() = true", i)
		}
	}
}

func TestTreeSetIteratorSetValue(t *testing.T) {
	tset := NewTreeSet(CompareInt)
	for i := 1; i <= 100; i++ {
		tset.Add(i)
	}

	// forward
	for it := tset.Iterator(); it.Next() && it.Value().(int) <= 100; it.Reset() {
		v := it.Value().(int) + 100
		it.SetValue(v)
	}
	// fmt.Println(tset)
	for i, it := 1, tset.Iterator(); it.Next(); i++ {
		v := it.Value()
		w := i + 100
		if v != w {
			t.Fatalf("Set[%d] = %v, want %v", i, v, w)
		}
	}

	// backward
	for it := tset.Iterator(); it.Prev() && it.Value().(int) <= 200; {
		v := it.Value().(int) + 100
		it.SetValue(v)
	}
	for i, it := 1, tset.Iterator(); it.Next(); i++ {
		v := it.Value()
		w := i + 200
		if v != w {
			t.Fatalf("Set[%d] = %v, want %v", i-1, v, w)
		}
	}
}

func TestTreeSetSort(t *testing.T) {
	for i := 1; i < 20; i++ {
		tset := NewTreeSet(CompareInt)

		a := make([]any, 0, 20)
		for n := i; n < 20; n++ {
			v := rand.Intn(1000)
			if !ars.Contains(a, v) {
				a = append(a)
			}
		}

		for j := len(a) - 1; j >= 0; j-- {
			tset.Add(a[j])
		}

		sort.Sort(inta(a))

		if !reflect.DeepEqual(a, tset.Values()) {
			t.Fatalf("%v != %v", a, tset.Values())
		}
	}
}

func TestTreeSetDelete(t *testing.T) {
	tset := NewTreeSet(CompareInt)

	for i := 1; i <= 100; i++ {
		tset.Add(i)
	}

	tset.DeleteIf(func(d any) bool {
		return d == 101
	})
	if tset.Len() != 100 {
		t.Error("TreeSet.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		tset.Delete(i, i)
		if tset.Len() != 100-i {
			t.Errorf("TreeSet.Delete(%v) failed, tset.Len() = %v, want %v", i, tset.Len(), 100-i)
		}
	}

	if !tset.IsEmpty() {
		t.Error("TreeSet.IsEmpty() should return true")
	}
}

func TestTreeSetDelete2(t *testing.T) {
	tset := NewTreeSet(CompareInt)

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			tset.Add(i)
		}
	}

	n := tset.Len()
	tset.Delete(100)
	if tset.Len() != n {
		t.Errorf("TreeSet.Delete(100).Len() = %v, want %v", tset.Len(), n)
	}

	for i := 0; i < 100; i++ {
		n = tset.Len()
		tset.Delete(i)
		z := i % 10
		if z > 1 {
			z = 1
		}
		a := n - tset.Len()
		if a != z {
			t.Errorf("TreeSet.Delete(%v) = %v, want %v", i, a, z)
		}
	}

	if !tset.IsEmpty() {
		t.Error("TreeSet.IsEmpty() should return true")
	}
}

func TestTreeSetDeleteAll(t *testing.T) {
	tset := NewTreeSet(CompareInt)

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			tset.Add(i)
		}
	}

	n := tset.Len()
	tset.Delete(100)
	if tset.Len() != n {
		t.Errorf("TreeSet.Delete(100).Len() = %v, want %v", tset.Len(), n)
	}

	for i := 0; i < 100; i++ {
		n = tset.Len()
		tset.DeleteAll(NewArrayList(i, i))
		z := i % 10
		if z > 1 {
			z = 1
		}
		a := n - tset.Len()
		if a != z {
			t.Errorf("TreeSet.Delete(%v) = %v, want %v", i, a, z)
		}
	}

	if !tset.IsEmpty() {
		t.Error("TreeSet.IsEmpty() should return true")
	}
}

func TestTreeSetString(t *testing.T) {
	e := "[1,2,3]"
	a := NewTreeSet(CompareInt, 1, 3, 2).String()
	if a != e {
		t.Errorf(`fmt.Sprintf("%%s", NewTreeSet(1, 3, 2)) = %v, want %v`, a, e)
	}
}

func TestTreeSetMarshalJSON(t *testing.T) {
	cs := []struct {
		tset *TreeSet
		json string
	}{
		{NewTreeSet(CompareString, "0", "1"), `["0","1"]`},
	}

	for i, c := range cs {
		bs, err := json.Marshal(c.tset)
		if err != nil {
			t.Errorf("[%d] json.Marshal(%v) error: %v", i, c.tset, err)
		}

		a := string(bs)
		if a != c.json {
			t.Errorf("[%d] json.Marshal(%v) = %q, want %q", i, c.tset, a, c.tset)
		}
	}
}

func TestTreeSetUnmarshalJSON(t *testing.T) {
	type Case struct {
		json string
		tset *TreeSet
	}

	cs := []Case{
		{`["1","0"]`, NewTreeSet(CompareString, "0", "1")},
	}

	for i, c := range cs {
		a := NewTreeSet(CompareString)
		err := json.Unmarshal([]byte(c.json), a)

		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%v) error: %v", i, c.json, err)
		}

		if a.String() != c.tset.String() {
			t.Errorf("[%d] json.Unmarshal(%v) = %v, want %v", i, c.json, a, c.tset)
		}
	}
}
