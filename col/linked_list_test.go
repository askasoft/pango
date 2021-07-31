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

	if av, ok := list2.Get(0); av != 1 || !ok {
		t.Errorf("Got %v expected %v", av, 1)
	}

	if av, ok := list2.Get(1); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}

	if av, ok := list2.Get(2); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
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
	if av, ok := list.Get(2); av != "c" || !ok {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestLinkedListRemove(t *testing.T) {
	list := NewLinkedList()
	list.Add("a")
	list.Add("b", "c")
	list.Remove(2)
	if av, ok := list.Get(2); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
	}
	list.Remove(1)
	list.Remove(0)
	list.Remove(0) // no effect
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedListGet(t *testing.T) {
	list := NewLinkedList()
	list.Add("a")
	list.Add("b", "c")
	if av, ok := list.Get(0); av != "a" || !ok {
		t.Errorf("Got %v expected %v", av, "a")
	}
	if av, ok := list.Get(1); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}
	if av, ok := list.Get(2); av != "c" || !ok {
		t.Errorf("Got %v expected %v", av, "c")
	}
	if av, ok := list.Get(3); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
	}
	list.Remove(0)
	if av, ok := list.Get(0); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestLinkedListSwap(t *testing.T) {
	list := NewLinkedList()
	list.Add("a")
	list.Add("b", "c")
	list.Swap(0, 1)
	if av, ok := list.Get(0); av != "b" || !ok {
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

func TestLinkedListValues(t *testing.T) {
	list := NewLinkedList()
	list.Add("a")
	list.Add("b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
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

func TestLinkedListInsert(t *testing.T) {
	list := NewLinkedList()
	list.Insert(0, "b", "c")
	list.Insert(0, "a")
	list.Insert(10, "x") // ignore
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
	list.Set(4, "d")  // ignore
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

func TestLinkedListIteratorNextOnEmpty(t *testing.T) {
	list := NewLinkedList()
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
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

func TestLinkedListIteratorPrevOnEmpty(t *testing.T) {
	list := NewLinkedList()
	it := list.Iterator()
	for it.Prev() {
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

	it.Reset()
	it.Prev()
	if value := it.Value(); value != "c" {
		t.Errorf("Got %v expected %v", value, "c")
	}
}

func TestLinkedListIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewLinkedList()

		for n := 0; n < i; n++ {
			l.Add(n)
		}

		it := l.Iterator()

		it.Remove()
		if l.Len() != i {
			t.Errorf("[%d] l.Len() == %v, want %v", i, l.Len(), i)
		}

		// remove middle
		x := rand.Intn(i-4) + 1
		for j := 0; j <= x; j++ {
			it.Next()
		}

		v := it.Value().(int)
		it.Remove()
		if l.Len() != i-1 {
			t.Errorf("[%d] l.Len() == %v, want %v", i, l.Len(), i-1)
		}
		if l.Contains(v) {
			t.Errorf("[%d] l.Contains(%v) = true", i, v)
		}

		it.Next()
		if v+1 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		it.Remove()
		if l.Contains(v + 1) {
			t.Errorf("[%d] l.Contains(%v) = true", i, v+1)
		}

		it.Prev()
		if v-1 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		it.Remove()
		if l.Contains(v - 1) {
			t.Errorf("[%d] l.Contains(%v) = true", i, v-1)
		}

		// remove first
		for it.Prev() {
		}
		it.Remove()
		if l.Contains(0) {
			t.Errorf("[%d] l.Contains(%v) = true", i, 0)
		}
		if it.Prev() {
			t.Errorf("[%d] l.Prev() = true", i)
		}

		// remove last
		for it.Next() {
		}
		it.Remove()
		if l.Contains(i - 1) {
			t.Errorf("[%d] l.Contains(%v) = true", i, i-1)
		}
		if it.Next() {
			t.Errorf("[%d] l.Next() = true", i)
		}

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				it.Remove()
			}
		} else {
			for it.Next() {
				it.Remove()
			}
		}
		if !l.IsEmpty() {
			t.Errorf("[%d] l.IsEmpty() = true", i)
		}
	}
}

func TestLinkedListLazyInit(t *testing.T) {
	{
		l := &LinkedList{}
		if l.Len() != 0 {
			t.Error("l.Len() != 0")
		}
		if !l.IsEmpty() {
			t.Error("l.IsEmpty() = true")
		}
		if len(l.Values()) != 0 {
			t.Error("len(l.Values()) != 0")
		}
		if l.Contains(1) {
			t.Error("l.Contains(1) = true, want false")
		}
		if l.Len() != 0 {
			t.Error("l.Len() != 0")
		}
		if l.Front() != nil {
			t.Error("l.Front() != nil")
		}
		if l.Back() != nil {
			t.Error("l.Back() != nil")
		}
		if n, i := l.Search(1); n != -1 || i != nil {
			t.Errorf("l.Search(1) == (%v, %v), want (%v, %v)", n, i, -1, nil)
		}
	}
	{
		l := &LinkedList{}
		l.PushFront(1)
		if l.Len() != 1 {
			t.Errorf("l.Len() = %v, want 1", l.Len())
		}
	}
	{
		l := &LinkedList{}
		l.PushBack(1)
		if l.Len() != 1 {
			t.Errorf("l.Len() = %v, want 1", l.Len())
		}
	}
	{
		l := &LinkedList{}
		l.Delete(1)
		if l.Len() != 0 {
			t.Error("l.Len() != 0")
		}
	}
	{
		l := &LinkedList{}
		l.Each(func(v interface{}) {})
		l.ReverseEach(func(v interface{}) {})
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

func checkLinkedListPointers(t *testing.T, l *LinkedList, lis []*LinkedListItem) {
	root := &l.root

	if !checkLinkedListLen(t, l, len(lis)) {
		return
	}

	// zero length lists must be the zero value or properly initialized (sentinel circle)
	if len(lis) == 0 {
		if l.root.next != nil && l.root.next != root || l.root.prev != nil && l.root.prev != root {
			t.Errorf("l.root.next = %p, l.root.prev = %p; both should both be nil or %p", l.root.next, l.root.prev, root)
		}
		return
	}
	// len(lis) > 0

	// check internal and external prev/next connections
	for i, e := range lis {
		prev := root
		var Prev *LinkedListItem
		if i > 0 {
			prev = lis[i-1]
			Prev = prev
		}
		if p := e.prev; p != prev {
			t.Errorf("elt[%d](%p).prev = %p, want %p", i, e, p, prev)
		}
		if p := e.Prev(); p != Prev {
			t.Errorf("elt[%d](%p).Prev() = %p, want %p", i, e, p, Prev)
		}

		next := root
		var Next *LinkedListItem
		if i < len(lis)-1 {
			next = lis[i+1]
			Next = next
		}
		if n := e.next; n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, e, n, next)
		}
		if n := e.Next(); n != Next {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, e, n, Next)
		}
	}
}

func TestLinkedListBasic(t *testing.T) {
	l := NewLinkedList()
	checkLinkedListPointers(t, l, []*LinkedListItem{})

	// Single item list
	e := l.PushFront("a")
	checkLinkedListPointers(t, l, []*LinkedListItem{e})
	l.MoveToFront(e)
	checkLinkedListPointers(t, l, []*LinkedListItem{e})
	l.MoveToBack(e)
	checkLinkedListPointers(t, l, []*LinkedListItem{e})
	e.Remove()
	checkLinkedListPointers(t, l, []*LinkedListItem{})

	// Bigger list
	e2 := l.PushFront(2)
	e1 := l.PushFront(1)
	e3 := l.PushBack(3)
	e4 := l.PushBack("banana")
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e2, e3, e4})

	e2.Remove()
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e3, e4})

	l.MoveToFront(e3) // move from middle
	checkLinkedListPointers(t, l, []*LinkedListItem{e3, e1, e4})

	l.MoveToFront(e1)
	l.MoveToBack(e3) // move from middle
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e4, e3})

	l.MoveToFront(e3) // move from back
	checkLinkedListPointers(t, l, []*LinkedListItem{e3, e1, e4})
	l.MoveToFront(e3) // should be no-op
	checkLinkedListPointers(t, l, []*LinkedListItem{e3, e1, e4})

	l.MoveToBack(e3) // move from front
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e4, e3})
	l.MoveToBack(e3) // should be no-op
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e4, e3})

	e2 = l.InsertBefore(e1, 2) // insert before front
	checkLinkedListPointers(t, l, []*LinkedListItem{e2, e1, e4, e3})
	e2.Remove()

	e2 = l.InsertBefore(e4, 2) // insert before middle
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e2, e4, e3})
	e2.Remove()

	e2 = l.InsertBefore(e3, 2) // insert before back
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e4, e2, e3})
	e2.Remove()

	e2 = l.InsertAfter(e1, 2) // insert after front
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e2, e4, e3})
	e2.Remove()

	e2 = l.InsertAfter(e4, 2) // insert after middle
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e4, e2, e3})
	e2.Remove()

	e2 = l.InsertAfter(e3, 2) // insert after back
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e4, e3, e2})
	e2.Remove()

	// Check standard iteration.
	sum := 0
	for e := l.Front(); e != nil; e = e.Next() {
		if i, ok := e.Value().(int); ok {
			sum += i
		}
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// Clear all items by iterating
	var next *LinkedListItem
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		e.Remove()
	}
	checkLinkedListPointers(t, l, []*LinkedListItem{})
}

func checkLinkedList(t *testing.T, l *LinkedList, evs []interface{}) {
	if !checkLinkedListLen(t, l, len(evs)) {
		return
	}

	for i, e := 0, l.Front(); e != nil; i, e = i+1, e.Next() {
		v := e.Value().(int)
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

func TestLinkedListRemove2(t *testing.T) {
	l := NewLinkedList()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e2})
	e := l.Front()
	e.Remove()
	checkLinkedListPointers(t, l, []*LinkedListItem{e2})
	e.Remove()
	checkLinkedListPointers(t, l, []*LinkedListItem{e2})
}

func TestLinkedListInsertBefore(t *testing.T) {
	l1 := NewLinkedList()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := NewLinkedList()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(e, 8)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestLinkedListRemove1(t *testing.T) {
	l := NewLinkedList(1, 2, 3)

	e := l.Front().Next()
	e.Remove()
	if e.Value() != 2 {
		t.Errorf("e.value = %d, want 2", e.Value())
	}
	if e.Next() != l.Back() {
		t.Errorf("e.Next() != l.Back()")
	}
	if e.Prev() != l.Front() {
		t.Errorf("e.Prev() != l.Front()")
	}
}

func TestLinkedListMove(t *testing.T) {
	l := NewLinkedList()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	e4 := l.PushBack(4)

	l.MoveAfter(e3, e3)
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e2, e3, e4})
	l.MoveBefore(e2, e2)
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e2, e3, e4})

	l.MoveAfter(e2, e3)
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e2, e3, e4})
	l.MoveBefore(e2, e2)
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e2, e3, e4})

	l.MoveBefore(e4, e2)
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e3, e2, e4})
	e2, e3 = e3, e2

	l.MoveBefore(e1, e4)
	checkLinkedListPointers(t, l, []*LinkedListItem{e4, e1, e2, e3})
	e1, e2, e3, e4 = e4, e1, e2, e3

	l.MoveAfter(e1, e4)
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e4, e2, e3})
	e2, e3, e4 = e4, e2, e3

	l.MoveAfter(e3, e2)
	checkLinkedListPointers(t, l, []*LinkedListItem{e1, e3, e2, e4})
	e2, e3 = e3, e2
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an item of l.
func TestLinkedListInsertBeforeUnknownMark(t *testing.T) {
	l := NewLinkedList(1, 2, 3)
	l.InsertBefore(new(LinkedListItem), 1)
	checkLinkedList(t, l, []interface{}{1, 2, 3})
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an item of l.
func TestLinkedListInsertAfterUnknownMark(t *testing.T) {
	l := NewLinkedList(1, 2, 3)
	l.InsertAfter(new(LinkedListItem), 1)
	checkLinkedList(t, l, []interface{}{1, 2, 3})
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an item of l.
func TestLinkedListMoveUnknownMark(t *testing.T) {
	l1 := NewLinkedList()
	e1 := l1.PushBack(1)

	l2 := NewLinkedList()
	e2 := l2.PushBack(2)

	l1.MoveAfter(e2, e1)
	checkLinkedList(t, l1, []interface{}{1})
	checkLinkedList(t, l2, []interface{}{2})

	l1.MoveBefore(e2, e1)
	checkLinkedList(t, l1, []interface{}{1})
	checkLinkedList(t, l2, []interface{}{2})
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

func TestLinkedListItem(t *testing.T) {
	l := NewLinkedList()

	for i := 0; i < 100; i++ {
		l.PushBack(i)
	}

	if l.Item(100) != nil {
		t.Error("LinkedList.Item(100) != nil")
	}
	if l.Item(-101) != nil {
		t.Error("LinkedList.Item(-101) != nil")
	}

	for n := 0; n < 100; n++ {
		var a interface{}

		w := n
		s := n
		i := l.Item(s)
		if i != nil {
			a = i.Value()
		}
		if i == nil || a != w {
			t.Errorf("LinkedList.Item(%v).Value = %v, want %v", s, a, w)
		}

		s = n - 100
		i = l.Item(s)
		if i != nil {
			a = i.Value()
		}
		if i == nil || a != w {
			t.Errorf("LinkedList.Item(%v).Value = %v, want %v", s, a, w)
		}
	}
}

func TestLinkedListSwapItem(t *testing.T) {
	l := NewLinkedList()

	for i := 0; i < 100; i++ {
		l.PushBack(i)
	}

	if l.SwapItem(l.Item(0), l.Front()) {
		t.Error("l.SwapItem(l.Item(0), l.Front()) = true, want false")
	}

	for i := 0; i < 50; i++ {
		ia := l.Item(i)
		ib := l.Item(-i - 1)
		a := l.SwapItem(ia, ib)
		if !a {
			t.Errorf("LinkedList.SwapItem(%v, %v) = %v, want %v", ia, ib, a, true)
		}
	}

	for i := 0; i < 100; i++ {
		e := 100 - 1 - i
		a := l.Item(i).Value()
		if a != e {
			t.Errorf("LinkedList.Item(%v).Value = %v, want %v", i, a, e)
		}
	}
}

func TestLinkedListSearch(t *testing.T) {
	l := NewLinkedList(1, 11)

	n111 := l.PushBack(111)

	l.PushBack("1", "11")
	s111 := l.PushBack("111")

	n := (100+1)/101 + 110

	sn, se := l.Search(n)
	if n111 != se || sn != 2 {
		t.Errorf("LinkedList [%v] should contains %v", l, n)
	}

	n++
	sn, se = l.Search(n)
	if se != nil || sn != -1 {
		t.Errorf("LinkedList [%v] should not contains %v", l, n)
	}

	s := str.Repeat("1", 3)

	sn, se = l.Search(s)
	if s111 != se || sn != 5 {
		t.Errorf("LinkedList [%v] should contains %v", l, s)
	}

	s += "0"
	sn, se = l.Search(s)
	if se != nil || sn != -1 {
		t.Errorf("LinkedList [%v] should not contains %v", l, s)
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

func TestLinkedListItemString(t *testing.T) {
	cs := []struct {
		e string
		s interface{}
	}{
		{"a", "a"},
		{"1", 1},
	}

	for _, c := range cs {
		i := &LinkedListItem{value: c.s}
		a := i.String()
		if a != c.e {
			t.Errorf("LinkedListItem(%v).String() = %q, want %q", c.s, a, c.e)
		}
	}
}
