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
	"github.com/pandafw/pango/str"
)

func TestLinkedHashSetInterface(t *testing.T) {
	var s Set = NewLinkedHashSet()
	if s == nil {
		t.Error("LinkedHashSet is not a Set")
	}
}

func TestLinkedHashSetNew(t *testing.T) {
	lset1 := NewLinkedHashSet()
	if av := lset1.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	lset2 := NewLinkedHashSet(1, "b")
	if av := lset2.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	if av := lset2.Get(0); av != 1 {
		t.Errorf("Got %v expected %v", av, 1)
	}
	if av := lset2.Get(1); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestLinkedHashSetAdd(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a")
	lset.Add("b", "c")
	if av := lset.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av := lset.Get(2); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestLinkedHashSetRemove(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a")
	lset.Add("b", "c")
	lset.Remove(2)
	lset.Remove(1)
	lset.Remove(0)
	if av := lset.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedHashSetRemovePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedHashSet("a")
	list.Remove(1)
}

func TestLinkedHashSetGet(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "a")
	lset.Add("b", "c", "b", "c")
	if av := lset.Get(0); av != "a" {
		t.Errorf("Got %v expected %v", av, "a")
	}
	if av := lset.Get(1); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
	if av := lset.Get(2); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
	lset.Remove(0)
	if av := lset.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestLinkedHashSetGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedHashSet("a")
	list.Get(1)
}

func TestLinkedHashSetSwap(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "a")
	lset.Add("b", "c", "b", "c")
	lset.Swap(0, 1)
	if av := lset.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestLinkedHashSetClear(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("e", "f", "g", "a", "b", "c", "d")
	lset.Add("e", "f", "g", "a", "b", "c", "d")
	lset.Clear()
	if av := lset.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedHashSetContains(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "a")
	lset.Add("b", "c", "b", "c")
	if av := lset.Contains("a"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Contains("a", "b", "c"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Contains("a", "b", "c", "d"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	lset.Clear()
	if av := lset.Contains("a"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := lset.Contains("a", "b", "c"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestLinkedHashSetValues(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "a")
	lset.Add("b", "c", "b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", lset.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashSetInsert(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Insert(0, "b", "c", "b", "c")
	lset.Insert(0, "a", "a")
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	lset.Insert(3, "d") // append
	if av := lset.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 4)
	}
	if av, ev := fmt.Sprintf("%s%s%s%s", lset.Values()...), "abcd"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashSetInsertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedHashSet("a")
	list.Insert(2, "b")
}

func TestLinkedHashSetSet(t *testing.T) {
	lset := NewLinkedHashSet("0", "1")
	lset.Set(0, "a")
	lset.Set(1, "b")
	if av := lset.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	lset.Add("")
	lset.Set(2, "c") // last
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	lset.Set(1, "bb") // update
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av, ev := fmt.Sprintf("%s%s%s", lset.Values()...), "abbc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	lset.Set(2, "cc") // last to first traversal
	lset.Set(0, "aa") // first to last traversal
	if av, ev := fmt.Sprintf("%s%s%s", lset.Values()...), "aabbcc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashSetSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedHashSet("a")
	list.Set(1, "b")
}

func TestLinkedHashSetEach(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "b", "c")
	lset.Add("a", "b", "c")
	index := 0
	lset.Each(func(value interface{}) {
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

func TestLinkedHashSetIteratorNextOnEmpty(t *testing.T) {
	lset := NewLinkedHashSet()
	it := lset.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty lset")
	}
}

func TestLinkedHashSetIteratorNext(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "b", "c")
	lset.Add("a", "b", "c")
	it := lset.Iterator()
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

func TestLinkedHashSetIteratorPrevOnEmpty(t *testing.T) {
	lset := NewLinkedHashSet()
	it := lset.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty lset")
	}
}

func TestLinkedHashSetIteratorPrev(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "b", "c")
	lset.Add("a", "b", "c")
	it := lset.Iterator()
	count := 0
	index := lset.Len()
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

func TestLinkedHashSetIteratorReset(t *testing.T) {
	lset := NewLinkedHashSet()
	it := lset.Iterator()
	lset.Add("a", "b", "c")
	lset.Add("a", "b", "c")
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

func assertLinkedHashSetIteratorRemove(t *testing.T, i int, it Iterator, w *LinkedHashSet) int {
	it.Remove()

	v := it.Value()
	w.Delete(v)

	it.SetValue(9999)

	lset := it.(*linkedHashSetIterator).lset
	if lset.Contains(v) {
		t.Fatalf("[%d] lset.Contains(%v) = true", i, v)
	}

	if lset.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, lset.String(), w.String())
	}

	return v.(int)
}

func TestLinkedHashSetIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		lset := NewLinkedHashSet()
		wset := NewLinkedHashSet()

		for n := 0; n < i; n++ {
			lset.Add(n)
			wset.Add(n)
		}

		it := lset.Iterator()

		it.Remove()
		if lset.Len() != i {
			t.Fatalf("[%d] lset.Len() == %v, want %v", i, lset.Len(), i)
		}

		// remove middle
		for j := 0; j <= lset.Len()/2; j++ {
			it.Next()
		}

		v := assertLinkedHashSetIteratorRemove(t, i, it, wset)

		it.Next()
		if v+1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		assertLinkedHashSetIteratorRemove(t, i, it, wset)

		it.Prev()
		if v-1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		assertLinkedHashSetIteratorRemove(t, i, it, wset)

		// remove first
		for it.Prev() {
		}
		assertLinkedHashSetIteratorRemove(t, i, it, wset)

		// remove last
		for it.Next() {
		}
		assertLinkedHashSetIteratorRemove(t, i, it, wset)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertLinkedHashSetIteratorRemove(t, i, it, wset)
			}
		} else {
			for it.Next() {
				assertLinkedHashSetIteratorRemove(t, i, it, wset)
			}
		}
		if !lset.IsEmpty() {
			t.Fatalf("[%d] lset.IsEmpty() = true", i)
		}
	}
}

func TestLinkedHashSetIteratorSetValue(t *testing.T) {
	ls := NewLinkedHashSet()
	for i := 1; i <= 100; i++ {
		ls.Add(i)
	}

	// forward (1->2, 3->4, ... )
	for it := ls.Iterator(); it.Next(); {
		v := it.Value().(int) + 1
		it.SetValue(v)
	}
	for i := 1; i <= ls.Len(); i++ {
		v := ls.Get(i - 1).(int)
		w := i * 2
		if v != w {
			t.Fatalf("Set[%d] = %v, want %v", i-1, v, w)
		}
	}

	// backward (100 -> 98, 96 -> 94)
	for it := ls.Iterator(); it.Prev(); {
		v := it.Value().(int) - 2
		it.SetValue(v)
	}
	for i := 1; i <= ls.Len(); i++ {
		v := ls.Get(i - 1).(int)
		w := i*4 - 2
		if v != w {
			t.Fatalf("Set[%d] = %v, want %v", i-1, v, w)
		}
	}
}

func TestLinkedHashSetSort(t *testing.T) {
	for i := 1; i < 20; i++ {
		lset := NewLinkedHashSet()

		a := make([]interface{}, 0, 20)
		for n := i; n < 20; n++ {
			v := rand.Intn(1000)
			if !ars.Contains(a, v) {
				a = append(a)
			}
		}

		for j := len(a) - 1; j >= 0; j-- {
			lset.Add(a[j])
		}

		lset.Sort(cmp.LessInt)
		sort.Sort(inta(a))

		if !reflect.DeepEqual(a, lset.Values()) {
			t.Fatalf("%v != %v", a, lset.Values())
		}
	}
}

func checkLinkedHashSetLen(t *testing.T, lset *LinkedHashSet, len int) bool {
	if n := lset.Len(); n != len {
		t.Errorf("lset.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkLinkedHashSet(t *testing.T, lset *LinkedHashSet, evs []interface{}) {
	if !checkLinkedHashSetLen(t, lset, len(evs)) {
		return
	}

	for i, it := 0, lset.Iterator(); it.Next(); i++ {
		v := it.Value().(int)
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}

	avs := lset.Values()
	for i, v := range avs {
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}
}

func TestLinkedHashSetExtending(t *testing.T) {
	l1 := NewLinkedHashSet(1, 2, 3)
	l2 := NewLinkedHashSet()
	l2.PushBack(4)
	l2.PushBack(5)

	l3 := NewLinkedHashSet()
	l3.PushBackAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})
	l3.PushBackAll(l2)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3, 4, 5})

	l3 = NewLinkedHashSet()
	l3.PushFrontAll(l2)
	checkLinkedHashSet(t, l3, []interface{}{4, 5})
	l3.PushFrontAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3, 4, 5})

	checkLinkedHashSet(t, l1, []interface{}{1, 2, 3})
	checkLinkedHashSet(t, l2, []interface{}{4, 5})

	l3 = NewLinkedHashSet()
	l3.PushBackAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})
	l3.PushBackAll(l3)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})

	l3 = NewLinkedHashSet()
	l3.PushFrontAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})
	l3.PushFrontAll(l3)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})

	l3 = NewLinkedHashSet()
	l1.PushBackAll(l3)
	checkLinkedHashSet(t, l1, []interface{}{1, 2, 3})
	l1.PushFrontAll(l3)
	checkLinkedHashSet(t, l1, []interface{}{1, 2, 3})

	l1.Clear()
	l2.Clear()
	l3.Clear()
	l1.PushBack(1, 2, 3)
	checkLinkedHashSet(t, l1, []interface{}{1, 2, 3})
	l2.PushBack(4, 5)
	checkLinkedHashSet(t, l2, []interface{}{4, 5})
	l3.PushBackAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})
	l3.PushBack(4, 5)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3, 4, 5})
	l3.PushFront(4, 5)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3, 4, 5})
}

func TestLinkedHashSetContains2(t *testing.T) {
	lset := NewLinkedHashSet(1, 11, 111, "1", "11", "111")

	n := (100+1)/101 + 110

	if !lset.Contains(n) {
		t.Errorf("LinkedHashSet [%v] should contains %v", lset, n)
	}

	n++
	if lset.Contains(n) {
		t.Errorf("LinkedHashSet [%v] should not contains %v", lset, n)
	}

	s := str.Repeat("1", 3)

	if !lset.Contains(s) {
		t.Errorf("LinkedHashSet [%v] should contains %v", lset, s)
	}

	s += "0"
	if lset.Contains(s) {
		t.Errorf("LinkedHashSet [%v] should not contains %v", lset, s)
	}
}

func TestLinkedHashSetDelete(t *testing.T) {
	lset := NewLinkedHashSet()

	for i := 1; i <= 100; i++ {
		lset.PushBack(i)
	}

	lset.Delete(101)
	if lset.Len() != 100 {
		t.Error("LinkedHashSet.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		lset.Delete(i)
		if lset.Len() != 100-i {
			t.Errorf("LinkedHashSet.Delete(%v) failed, lset.Len() = %v, want %v", i, lset.Len(), 100-i)
		}
	}

	if !lset.IsEmpty() {
		t.Error("LinkedHashSet.IsEmpty() should return true")
	}
}

func TestLinkedHashSetDeleteAll(t *testing.T) {
	lset := NewLinkedHashSet()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			lset.PushBack(i)
		}
	}

	n := lset.Len()
	lset.Delete(100)
	if lset.Len() != n {
		t.Errorf("LinkedHashSet.Delete(100).Len() = %v, want %v", lset.Len(), n)
	}

	for i := 0; i < 100; i++ {
		n = lset.Len()
		lset.Delete(i)
		z := i % 10
		if z > 1 {
			z = 1
		}
		a := n - lset.Len()
		if a != z {
			t.Errorf("LinkedHashSet.Delete(%v) = %v, want %v", i, a, z)
		}
	}

	if !lset.IsEmpty() {
		t.Error("LinkedHashSet.IsEmpty() should return true")
	}
}

func TestLinkedHashSetString(t *testing.T) {
	e := "[1,3,2]"
	a := fmt.Sprintf("%s", NewLinkedHashSet(1, 3, 2))
	if a != e {
		t.Errorf(`fmt.Sprintf("%%s", NewLinkedHashSet(1, 3, 2)) = %v, want %v`, a, e)
	}
}

func TestLinkedHashSetMarshalJSON(t *testing.T) {
	cs := []struct {
		lset *LinkedHashSet
		json string
	}{
		{NewLinkedHashSet(0, 1, "0", "1", 0.0, 1.0, true, false), `[0,1,"0","1",0,1,true,false]`},
		{NewLinkedHashSet(0, "1", 2.0, [2]int{1, 2}), `[0,"1",2,[1,2]]`},
	}

	for i, c := range cs {
		bs, err := json.Marshal(c.lset)
		if err != nil {
			t.Errorf("[%d] json.Marshal(%v) error: %v", i, c.lset, err)
		}

		a := string(bs)
		if a != c.json {
			t.Errorf("[%d] json.Marshal(%v) = %q, want %q", i, c.lset, a, c.lset)
		}
	}
}

func TestLinkedHashSetUnmarshalJSON(t *testing.T) {
	type Case struct {
		json string
		lset *LinkedHashSet
	}

	cs := []Case{
		{`["0","1",0,1,true,false]`, NewLinkedHashSet("0", "1", 0.0, 1.0, true, false)},
		{`["1",2]`, NewLinkedHashSet("1", 2.0)},
	}

	for i, c := range cs {
		a := NewLinkedHashSet()
		err := json.Unmarshal([]byte(c.json), a)

		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%v) error: %v", i, c.json, err)
		}

		if !reflect.DeepEqual(a, c.lset) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %q", i, c.json, a, c.lset)
		}
	}
}
