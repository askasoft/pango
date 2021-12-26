package col

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/pandafw/pango/cmp"
	"github.com/pandafw/pango/str"
)

func TestArrayListInterface(t *testing.T) {
	var l List = NewArrayList()
	if l == nil {
		t.Error("ArrayList is not a List")
	}

	var s Sortable = NewArrayList()
	if s == nil {
		t.Error("ArrayList is not a Sortable")
	}
}

func TestArrayListNew(t *testing.T) {
	list1 := NewArrayList()

	if av := list1.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	list2 := NewArrayList(1, "b")

	if av := list2.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}

	if av := list2.Get(0); av != 1 {
		t.Errorf("Got %v expected %v", av, 1)
	}

	if av := list2.Get(1); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestArrayListAdd(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	if av := list.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av := list.Get(2); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestArrayListGrow(t *testing.T) {
	list := NewArrayList()

	for i := 0; i < 1000; i++ {
		list.Add(i)
		if l := list.Len(); l != i+1 {
			t.Errorf("list.Len() = %v, want %v", l, i+1)
		}

		wc := list.roundup(i + 1)
		ac := cap(list.data)
		if wc != ac {
			t.Errorf("list.Cap(%d) = %v, want %v", i+1, ac, wc)
		}

		for n := 0; n <= i; n++ {
			if v := list.Get(n); v != n {
				t.Errorf("list.Get(%d) = %v, want %v", n, v, n)
			}
		}
	}
}

func TestArrayListIndex(t *testing.T) {
	list := NewArrayList()

	expectedIndex := -1
	if index := list.Index("a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	list.Add("a")
	list.Add("b", "c")

	expectedIndex = 0
	if index := list.Index("a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 1
	if index := list.Index("b"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 2
	if index := list.Index("c"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}
}

func TestArrayListDelete(t *testing.T) {
	l := NewArrayList()

	for i := 1; i <= 100; i++ {
		l.PushBack(i)
	}

	l.Delete(101)
	if l.Len() != 100 {
		t.Error("ArrayList.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		l.Delete(i)
		if l.Len() != 100-i {
			t.Errorf("ArrayList.Delete(%v) failed, l.Len() = %v, want %v", i, l.Len(), 100-i)
		}
	}

	if !l.IsEmpty() {
		t.Error("ArrayList.IsEmpty() should return true")
	}
}

func TestArrayListDeleteAll(t *testing.T) {
	l := NewArrayList()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			l.PushBack(i)
		}
	}

	n := l.Len()
	l.Delete(100)
	if l.Len() != n {
		t.Errorf("ArrayList.Delete(100).Len() = %v, want %v", l.Len(), n)
	}
	for i := 0; i < 100; i++ {
		n = l.Len()
		z := i % 10
		l.Delete(i)
		a := n - l.Len()
		if a != z {
			t.Errorf("ArrayList.Delete(%v) = %v, want %v", i, a, z)
		}
	}

	if !l.IsEmpty() {
		t.Error("ArrayList.IsEmpty() should return true")
	}
}

func TestArrayListRemove(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	list.Remove(2)
	list.Remove(1)
	list.Remove(0)
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestArrayListRemovePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewArrayList("a")
	list.Remove(1)
}

func TestArrayListGet(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	if av := list.Get(0); av != "a" {
		t.Errorf("Got %v expected %v", av, "a")
	}
	if av := list.Get(1); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
	if av := list.Get(2); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
	list.Remove(0)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestArrayListGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewArrayList("a")
	list.Get(1)
}

func TestArrayListSwap(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	list.Swap(0, 1)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestArrayListClear(t *testing.T) {
	list := NewArrayList()
	list.Add("e", "f", "g", "a", "b", "c", "d")
	list.Clear()
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestArrayListContains(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	if av := list.Contains("a"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Contains("a", "b", "c"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Contains("a", "b", "c", "d"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	list.Clear()
	if av := list.Contains("a"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.Contains("a", "b", "c"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestArrayListValues(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListInsert(t *testing.T) {
	list := NewArrayList()
	list.Insert(0, "b", "c")
	list.Insert(0, "a")
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	list.Insert(3, "d") // append
	if av := list.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 4)
	}
	if av, ev := fmt.Sprintf("%s%s%s%s", list.Values()...), "abcd"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListInsertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewArrayList("a")
	list.Insert(2, "b")
}

func TestArrayListSet(t *testing.T) {
	list := NewArrayList("", "")
	list.Set(0, "a")
	list.Set(1, "b")
	if av := list.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	list.Add("")
	list.Set(2, "c")
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	list.Set(1, "bb") // update
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abbc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewArrayList("a")
	list.Set(1, "b")
}

func TestArrayListEach(t *testing.T) {
	list := NewArrayList()
	list.Add("a", "b", "c")
	index := 0
	list.Each(func(value interface{}) {
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

func TestArrayListIteratorPrevOnEmpty(t *testing.T) {
	list := NewArrayList()
	it := list.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestArrayListIteratorNextOnEmpty(t *testing.T) {
	list := NewArrayList()
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestArrayListIteratorPrev(t *testing.T) {
	list := NewArrayList()
	list.Add("a", "b", "c")
	it := list.Iterator()
	count := 0
	index := list.Len()
	for it.Prev() {
		index--
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
	}
	if av, ev := count, 3; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListIteratorNext(t *testing.T) {
	list := NewArrayList()
	list.Add("a", "b", "c")
	it := list.Iterator()
	count := 0
	for index := 0; it.Next(); index++ {
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
	}
	if av, ev := count, 3; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListIteratorReset(t *testing.T) {
	list := NewArrayList()

	it := list.Iterator()
	list.Add("a", "b", "c")

	for it.Next() {
	}
	it.Reset()
	it.Next()
	if value := it.Value(); value != "a" {
		t.Errorf("Got %v expected %v", value, "a")
	}

	for it.Prev() {
	}
	it.Reset()
	it.Prev()
	if value := it.Value(); value != "c" {
		t.Errorf("Got %v expected %v", value, "c")
	}
}

func assertArrayListIteratorRemove(t *testing.T, i int, it Iterator, w *ArrayList) int {
	v := it.Value()

	it.Remove()

	w.Delete(v)

	it.SetValue(9999)

	l := it.(*arrayListIterator).list
	if l.Contains(v) {
		t.Fatalf("[%d] l.Contains(%v) = true", i, v)
	}

	if l.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, l.String(), w.String())
	}

	return v.(int)
}

func TestArrayListIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewArrayList()
		w := NewArrayList()

		for n := 0; n < i; n++ {
			l.Add(n)
			w.Add(n)
		}

		it := l.Iterator()

		it.Remove()
		it.SetValue(9999)
		if l.Len() != i {
			t.Fatalf("[%d] l.Len() == %v, want %v", i, l.Len(), i)
		}

		// remove middle
		for j := 0; j <= l.Len()/2; j++ {
			it.Next()
		}

		v := assertArrayListIteratorRemove(t, i, it, w)

		it.Next()
		if v+1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		assertArrayListIteratorRemove(t, i, it, w)

		it.Prev()
		if v-1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		assertArrayListIteratorRemove(t, i, it, w)

		// remove first
		for it.Prev() {
		}
		assertArrayListIteratorRemove(t, i, it, w)

		// remove last
		for it.Next() {
		}
		assertArrayListIteratorRemove(t, i, it, w)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertArrayListIteratorRemove(t, i, it, w)
			}
		} else {
			for it.Next() {
				assertArrayListIteratorRemove(t, i, it, w)
			}
		}
		if !l.IsEmpty() {
			t.Fatalf("[%d] l.IsEmpty() = true", i)
		}
	}
}

func TestArrayListIteratorSetValue(t *testing.T) {
	l := NewArrayList()
	for i := 1; i <= 100; i++ {
		l.Add(i)
	}

	// forward
	for it := l.Iterator(); it.Next(); {
		it.SetValue(it.Value().(int) + 100)
	}
	for i := 1; i <= l.Len(); i++ {
		v := l.Get(i - 1).(int)
		w := i + 100
		if v != w {
			t.Fatalf("List[%d] = %v, want %v", i-1, v, w)
		}
	}

	// backward
	for it := l.Iterator(); it.Prev(); {
		it.SetValue(it.Value().(int) + 100)
	}
	for i := 1; i <= l.Len(); i++ {
		v := l.Get(i - 1).(int)
		w := i + 200
		if v != w {
			t.Fatalf("List[%d] = %v, want %v", i-1, v, w)
		}
	}
}

func TestArrayListSort(t *testing.T) {
	for i := 1; i < 100; i++ {
		l := NewArrayList()

		a := make([]interface{}, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, rand.Intn(20))
		}

		for j := len(a) - 1; j >= 0; j-- {
			l.Add(a[j])
		}

		l.Sort(cmp.LessInt)
		sort.Sort(inta(a))

		if !reflect.DeepEqual(a, l.Values()) {
			t.Errorf("%v != %v", a, l.Values())
		}
	}
}

func checkArrayListLen(t *testing.T, l *ArrayList, len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkArrayList(t *testing.T, l *ArrayList, evs []interface{}) {
	if !checkArrayListLen(t, l, len(evs)) {
		return
	}

	for i, it := 0, l.Iterator(); it.Next(); i++ {
		v := it.Value().(int)
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}

	avs := l.Values()
	for i, v := range avs {
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}
}

func TestArrayListExtending(t *testing.T) {
	l1 := NewArrayList(1, 2, 3)
	l2 := NewArrayList()
	l2.PushBack(4)
	l2.PushBack(5)

	l3 := NewArrayList()
	l3.PushBackAll(l1)
	checkArrayList(t, l3, []interface{}{1, 2, 3})
	l3.PushBackAll(l2)
	checkArrayList(t, l3, []interface{}{1, 2, 3, 4, 5})

	l3 = NewArrayList()
	l3.PushFrontAll(l2)
	checkArrayList(t, l3, []interface{}{4, 5})
	l3.PushFrontAll(l1)
	checkArrayList(t, l3, []interface{}{1, 2, 3, 4, 5})

	checkArrayList(t, l1, []interface{}{1, 2, 3})
	checkArrayList(t, l2, []interface{}{4, 5})

	l3 = NewArrayList()
	l3.PushBackAll(l1)
	checkArrayList(t, l3, []interface{}{1, 2, 3})
	l3.PushBackAll(l3)
	checkArrayList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = NewArrayList()
	l3.PushFrontAll(l1)
	checkArrayList(t, l3, []interface{}{1, 2, 3})
	l3.PushFrontAll(l3)
	checkArrayList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = NewArrayList()
	l1.PushBackAll(l3)
	checkArrayList(t, l1, []interface{}{1, 2, 3})
	l1.PushFrontAll(l3)
	checkArrayList(t, l1, []interface{}{1, 2, 3})

	l1.Clear()
	l2.Clear()
	l3.Clear()
	l1.PushBack(1, 2, 3)
	checkArrayList(t, l1, []interface{}{1, 2, 3})
	l2.PushBack(4, 5)
	checkArrayList(t, l2, []interface{}{4, 5})
	l3.PushBackAll(l1)
	checkArrayList(t, l3, []interface{}{1, 2, 3})
	l3.PushBack(4, 5)
	checkArrayList(t, l3, []interface{}{1, 2, 3, 4, 5})
	l3.PushFront(4, 5)
	checkArrayList(t, l3, []interface{}{4, 5, 1, 2, 3, 4, 5})
}

func TestArrayListContains2(t *testing.T) {
	l := NewArrayList(1, 11, 111, "1", "11", "111")

	n := (100+1)/101 + 110

	if !l.Contains(n) {
		t.Errorf("ArrayList [%v] should contains %v", l, n)
	}

	n++
	if l.Contains(n) {
		t.Errorf("ArrayList [%v] should not contains %v", l, n)
	}

	s := str.Repeat("1", 3)

	if !l.Contains(s) {
		t.Errorf("ArrayList [%v] should contains %v", l, s)
	}

	s += "0"
	if l.Contains(s) {
		t.Errorf("ArrayList [%v] should not contains %v", l, s)
	}
}
func TestArrayListJSON(t *testing.T) {
	cs := []struct {
		s string
		a *ArrayList
	}{
		{`[]`, NewArrayList()},
		{`["a","b","c"]`, NewArrayList("a", "b", "c")},
	}

	for i, c := range cs {
		a := NewArrayList()
		json.Unmarshal(([]byte)(c.s), &a)
		if !reflect.DeepEqual(a.Values(), c.a.Values()) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %v", i, c.s, a.Values(), c.a.Values())
		}

		bs, _ := json.Marshal(a)
		if string(bs) != c.s {
			t.Errorf("[%d] json.Marshal(%v) = %q, want %q", i, a.Values(), string(bs), c.s)
		}
	}
}
