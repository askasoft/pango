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

func TestLinkedListInterface(t *testing.T) {
	var l List = NewLinkedList()
	if l == nil {
		t.Error("LinkedList is not a List")
	}
}

func TestLinkedListNew(t *testing.T) {
	list1 := NewLinkedList()

	if av := list1.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	list2 := NewLinkedList(1, "b")

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

func TestLinkedListAdd(t *testing.T) {
	list := NewLinkedList()
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

func TestLinkedListIndex(t *testing.T) {
	list := NewLinkedList()

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

func TestLinkedListDelete(t *testing.T) {
	l := NewLinkedList()

	for i := 1; i <= 100; i++ {
		l.PushBack(i)
	}

	l.Delete(101)
	if l.Len() != 100 {
		t.Error("LinkedList.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		l.Delete(i)
		if l.Len() != 100-i {
			t.Errorf("LinkedList.Delete(%v) failed, l.Len() = %v, want %v", i, l.Len(), 100-i)
		}
	}

	if !l.IsEmpty() {
		t.Error("LinkedList.IsEmpty() should return true")
	}
}

func TestLinkedListDeleteAll(t *testing.T) {
	l := NewLinkedList()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			l.PushBack(i)
		}
	}

	n := l.Len()
	l.Delete(100)
	if l.Len() != n {
		t.Errorf("LinkedList.Delete(100).Len() = %v, want %v", l.Len(), n)
	}
	for i := 0; i < 100; i++ {
		n = l.Len()
		z := i % 10
		l.Delete(i)
		a := n - l.Len()
		if a != z {
			t.Errorf("LinkedList.Delete(%v) = %v, want %v", i, a, z)
		}
	}

	if !l.IsEmpty() {
		t.Error("LinkedList.IsEmpty() should return true")
	}
}

func TestLinkedListRemove(t *testing.T) {
	list := NewLinkedList()
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

func TestLinkedListRemovePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedList("a")
	list.Remove(1)
}

func TestLinkedListGet(t *testing.T) {
	list := NewLinkedList()
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

func TestLinkedListGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedList("a")
	list.Get(1)
}

func TestLinkedListSwap(t *testing.T) {
	list := NewLinkedList()
	list.Add("a")
	list.Add("b", "c")
	list.Swap(0, 1)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestLinkedListClear(t *testing.T) {
	list := NewLinkedList()
	list.Add("e", "f", "g", "a", "b", "c", "d")
	list.Clear()
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedListContains(t *testing.T) {
	list := NewLinkedList()
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

func TestLinkedListContains2(t *testing.T) {
	l := NewLinkedList(1, 11, 111, "1", "11", "111")

	n := (100+1)/101 + 110

	if !l.Contains(n) {
		t.Errorf("LinkedList [%v] should contains %v", l, n)
	}

	n++
	if l.Contains(n) {
		t.Errorf("LinkedList [%v] should not contains %v", l, n)
	}

	s := str.Repeat("1", 3)

	if !l.Contains(s) {
		t.Errorf("LinkedList [%v] should contains %v", l, s)
	}

	s += "0"
	if l.Contains(s) {
		t.Errorf("LinkedList [%v] should not contains %v", l, s)
	}
}

func TestLinkedListValues(t *testing.T) {
	list := NewLinkedList()
	list.Add("a")
	list.Add("b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedListInsert(t *testing.T) {
	list := NewLinkedList()
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

func TestLinkedListInsertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedList("a")
	list.Insert(2, "b")
}

func TestLinkedListSet(t *testing.T) {
	list := NewLinkedList("", "")
	list.Set(0, "a")
	list.Set(1, "b")
	if av := list.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	list.Add("")
	list.Set(2, "c") // last
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
	list.Set(2, "cc") // last to first traversal
	list.Set(0, "aa") // first to last traversal
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "aabbcc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedListSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedList("a")
	list.Set(1, "b")
}

func TestLinkedListEach(t *testing.T) {
	list := NewLinkedList()
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

func TestLinkedListIteratorPrevOnEmpty(t *testing.T) {
	list := NewLinkedList()
	it := list.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestLinkedListIteratorNextOnEmpty(t *testing.T) {
	list := NewLinkedList()
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestLinkedListIteratorPrev(t *testing.T) {
	list := NewLinkedList()
	list.Add("a", "b", "c")
	it := list.Iterator()
	count := 0
	index := list.Len()
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

func TestLinkedListIteratorNext(t *testing.T) {
	list := NewLinkedList()
	list.Add("a", "b", "c")
	it := list.Iterator()
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

func TestLinkedListIteratorReset(t *testing.T) {
	list := NewLinkedList()

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

func assertLinkedListIteratorRemove(t *testing.T, i int, it Iterator, w *LinkedList) int {
	it.Remove()

	v := it.Value()
	w.Delete(v)

	it.SetValue(9999)

	l := it.(*linkedListIterator).list
	if l.Contains(v) {
		t.Fatalf("[%d] l.Contains(%v) = true", i, v)
	}

	if l.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, l.String(), w.String())
	}

	return v.(int)
}

func TestLinkedListIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewLinkedList()
		w := NewLinkedList()

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

		v := assertLinkedListIteratorRemove(t, i, it, w)

		it.Next()
		if v+1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		assertLinkedListIteratorRemove(t, i, it, w)

		it.Prev()
		if v-1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		assertLinkedListIteratorRemove(t, i, it, w)

		// remove first
		for it.Prev() {
		}
		assertLinkedListIteratorRemove(t, i, it, w)

		// remove last
		for it.Next() {
		}
		assertLinkedListIteratorRemove(t, i, it, w)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertLinkedListIteratorRemove(t, i, it, w)
			}
		} else {
			for it.Next() {
				assertLinkedListIteratorRemove(t, i, it, w)
			}
		}
		if !l.IsEmpty() {
			t.Fatalf("[%d] l.IsEmpty() = true", i)
		}
	}
}

func TestLinkedListIteratorSetValue(t *testing.T) {
	l := NewLinkedList()
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

func TestLinkedListSort(t *testing.T) {
	for i := 1; i < 100; i++ {
		l := NewLinkedList()

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

func checkLinkedListLen(t *testing.T, l *LinkedList, len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkLinkedList(t *testing.T, l *LinkedList, evs []interface{}) {
	if !checkLinkedListLen(t, l, len(evs)) {
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

func TestLinkedListExtending(t *testing.T) {
	l1 := NewLinkedList(1, 2, 3)
	l2 := NewLinkedList()
	l2.PushBack(4)
	l2.PushBack(5)

	l3 := NewLinkedList()
	l3.PushBackAll(l1)
	checkLinkedList(t, l3, []interface{}{1, 2, 3})
	l3.PushBackAll(l2)
	checkLinkedList(t, l3, []interface{}{1, 2, 3, 4, 5})

	l3 = NewLinkedList()
	l3.PushFrontAll(l2)
	checkLinkedList(t, l3, []interface{}{4, 5})
	l3.PushFrontAll(l1)
	checkLinkedList(t, l3, []interface{}{1, 2, 3, 4, 5})

	checkLinkedList(t, l1, []interface{}{1, 2, 3})
	checkLinkedList(t, l2, []interface{}{4, 5})

	l3 = NewLinkedList()
	l3.PushBackAll(l1)
	checkLinkedList(t, l3, []interface{}{1, 2, 3})
	l3.PushBackAll(l3)
	checkLinkedList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = NewLinkedList()
	l3.PushFrontAll(l1)
	checkLinkedList(t, l3, []interface{}{1, 2, 3})
	l3.PushFrontAll(l3)
	checkLinkedList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = NewLinkedList()
	l1.PushBackAll(l3)
	checkLinkedList(t, l1, []interface{}{1, 2, 3})
	l1.PushFrontAll(l3)
	checkLinkedList(t, l1, []interface{}{1, 2, 3})

	l1.Clear()
	l2.Clear()
	l3.Clear()
	l1.PushBack(1, 2, 3)
	checkLinkedList(t, l1, []interface{}{1, 2, 3})
	l2.PushBack(4, 5)
	checkLinkedList(t, l2, []interface{}{4, 5})
	l3.PushBackAll(l1)
	checkLinkedList(t, l3, []interface{}{1, 2, 3})
	l3.PushBack(4, 5)
	checkLinkedList(t, l3, []interface{}{1, 2, 3, 4, 5})
	l3.PushFront(4, 5)
	checkLinkedList(t, l3, []interface{}{4, 5, 1, 2, 3, 4, 5})
}

func TestLinkedListString(t *testing.T) {
	e := "[1,3,2]"
	a := fmt.Sprintf("%s", NewLinkedList(1, 3, 2))
	if a != e {
		t.Errorf(`fmt.Sprintf("%%s", NewLinkedList(1, 3, 2)) = %v, want %v`, a, e)
	}
}

func TestLinkedListMarshalJSON(t *testing.T) {
	cs := []struct {
		list *LinkedList
		json string
	}{
		{NewLinkedList(0, 1, "0", "1", 0.0, 1.0, true, false), `[0,1,"0","1",0,1,true,false]`},
		{NewLinkedList(0, "1", 2.0, []int{1, 2}, map[int]int{1: 10, 2: 20}), `[0,"1",2,[1,2],{"1":10,"2":20}]`},
	}

	for i, c := range cs {
		bs, err := json.Marshal(c.list)
		if err != nil {
			t.Errorf("[%d] json.Marshal(%v) error: %v", i, c.list, err)
		}

		a := string(bs)
		if a != c.json {
			t.Errorf("[%d] json.Marshal(%v) = %q, want %q", i, c.list, a, c.list)
		}
	}
}

func TestLinkedListUnmarshalJSON(t *testing.T) {
	type Case struct {
		json string
		list *LinkedList
	}

	cs := []Case{
		{`["0","1",0,1,true,false]`, NewLinkedList("0", "1", 0.0, 1.0, true, false)},
		{`["1",2,[1,2],{"1":10,"2":20}]`, NewLinkedList("1", 2.0, NewLinkedList(1.0, 2.0), JSONObject{"1": 10.0, "2": 20.0})},
	}

	for i, c := range cs {
		a := NewLinkedList()
		err := json.Unmarshal([]byte(c.json), a)

		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%v) error: %v", i, c.json, err)
		}

		if !reflect.DeepEqual(a, c.list) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %q", i, c.json, a, c.list)
		}
	}
}
