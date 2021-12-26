package col

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/pandafw/pango/ars"
	"github.com/pandafw/pango/cmp"
)

func TestTreeSetInterface(t *testing.T) {
	var s Set = NewTreeSet(cmp.CompareInt)
	if s == nil {
		t.Error("TreeSet is not a Set")
	}
}

func TestTreeSetNew(t *testing.T) {
	tset := NewTreeSet(cmp.CompareString)
	if av := tset.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	tset = NewTreeSet(cmp.CompareString, "1", "b")
	if av := tset.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	if av := tset.Front(); av != "1" {
		t.Errorf("Got %v expected %v", av, 1)
	}
	if av := tset.Back(); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestTreeSetAdd(t *testing.T) {
	tset := NewTreeSet(cmp.CompareString)
	tset.Add("a")
	tset.Add("b", "c")
	if av := tset.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := tset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av := tset.Back(); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestTreeSetClear(t *testing.T) {
	tset := NewTreeSet(cmp.CompareString)
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
	tset := NewTreeSet(cmp.CompareString)
	tset.Add("a", "a")
	tset.Add("b", "c", "b", "c")
	if av := tset.Contains("a"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := tset.Contains("a", "b", "c"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := tset.Contains("a", "b", "c", "d"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	tset.Clear()
	if av := tset.Contains("a"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := tset.Contains("a", "b", "c"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestTreeSetValues(t *testing.T) {
	tset := NewTreeSet(cmp.CompareString)
	tset.Add("a", "a")
	tset.Add("b", "c", "b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", tset.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestTreeSetEach(t *testing.T) {
	tset := NewTreeSet(cmp.CompareString)
	tset.Add("a", "b", "c")
	tset.Add("a", "b", "c")
	index := 0
	tset.Each(func(value interface{}) {
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
	tset := NewTreeSet(cmp.CompareString)
	it := tset.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty tset")
	}
}

func TestTreeSetIteratorNext(t *testing.T) {
	tset := NewTreeSet(cmp.CompareString)
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
	tset := NewTreeSet(cmp.CompareString)
	it := tset.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty tset")
	}
}

func TestTreeSetIteratorPrev(t *testing.T) {
	tset := NewTreeSet(cmp.CompareString)
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
	tset := NewTreeSet(cmp.CompareString)
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
		tset := NewTreeSet(cmp.CompareInt)
		wset := NewTreeSet(cmp.CompareInt)

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
	tset := NewTreeSet(cmp.CompareInt)
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
		tset := NewTreeSet(cmp.CompareInt)

		a := make([]interface{}, 0, 20)
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

func checkTreeSetLen(t *testing.T, tset *TreeSet, len int) bool {
	if n := tset.Len(); n != len {
		t.Errorf("tset.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkTreeSet(t *testing.T, tset *TreeSet, evs []interface{}) {
	if !checkTreeSetLen(t, tset, len(evs)) {
		return
	}

	for i, it := 0, tset.Iterator(); it.Next(); i++ {
		v := it.Value().(int)
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}

	avs := tset.Values()
	for i, v := range avs {
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}
}

func TestTreeSetDelete(t *testing.T) {
	tset := NewTreeSet(cmp.CompareInt)

	for i := 1; i <= 100; i++ {
		tset.Add(i)
	}

	tset.Delete(101)
	if tset.Len() != 100 {
		t.Error("TreeSet.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		tset.Delete(i)
		if tset.Len() != 100-i {
			t.Errorf("TreeSet.Delete(%v) failed, tset.Len() = %v, want %v", i, tset.Len(), 100-i)
		}
	}

	if !tset.IsEmpty() {
		t.Error("TreeSet.IsEmpty() should return true")
	}
}

func TestTreeSetDeleteAll(t *testing.T) {
	tset := NewTreeSet(cmp.CompareInt)

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

func TestTreeSetString(t *testing.T) {
	e := "[1,2,3]"
	a := fmt.Sprintf("%s", NewTreeSet(cmp.CompareInt, 1, 3, 2))
	if a != e {
		t.Errorf(`fmt.Sprintf("%%s", NewTreeSet(1, 3, 2)) = %v, want %v`, a, e)
	}
}

func TestTreeSetMarshalJSON(t *testing.T) {
	cs := []struct {
		tset *TreeSet
		json string
	}{
		{NewTreeSet(cmp.CompareString, "0", "1"), `["0","1"]`},
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
		{`["1","0"]`, NewTreeSet(cmp.CompareString, "0", "1")},
	}

	for i, c := range cs {
		a := NewTreeSet(cmp.CompareString)
		err := json.Unmarshal([]byte(c.json), a)

		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%v) error: %v", i, c.json, err)
		}

		if a.String() != c.tset.String() {
			t.Errorf("[%d] json.Unmarshal(%v) = %v, want %v", i, c.json, a, c.tset)
		}
	}
}