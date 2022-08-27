//go:build go1.18
// +build go1.18

package cog

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestRingBufferInterface(t *testing.T) {
	var l List[int] = NewRingBuffer[int]()
	if l == nil {
		t.Error("RingBuffer is not a List")
	}

	var q Queue[int] = NewRingBuffer[int]()
	if q == nil {
		t.Error("RingBuffer is not a Queue")
	}

	var dq Deque[int] = NewRingBuffer[int]()
	if dq == nil {
		t.Error("RingBuffer is not a Deque")
	}

	var s Sortable[int] = NewRingBuffer[int]()
	if s == nil {
		t.Error("RingBuffer is not a Sortable")
	}
}

func TestRingBufferNew(t *testing.T) {
	rb1 := NewRingBuffer[int]()

	if av := rb1.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	rb2 := NewRingBuffer(1, 2)

	if av := rb2.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}

	if av := rb2.Get(0); av != 1 {
		t.Errorf("Got %v expected %v", av, 1)
	}

	if av := rb2.Get(1); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
}

func TestRingBufferSimple(t *testing.T) {
	rb := NewRingBuffer[int]()

	for i := 0; i < minArrayCap; i++ {
		rb.Push(i)
	}

	for i := 0; i < minArrayCap; i++ {
		v, _ := rb.Peek()
		if v != i {
			t.Error("peek", i, "had value", v)
		}

		x, _ := rb.Poll()
		if x != i {
			t.Error("poll", i, "had value", x)
		}
	}
}

func TestRingBufferSimple2(t *testing.T) {
	rb := &RingBuffer[int]{}

	for i := 0; i < minArrayCap; i++ {
		rb.Push(i)
	}

	for i := 0; i < minArrayCap; i++ {
		v, _ := rb.Peek()
		if v != i {
			t.Error("peek", i, "had value", v)
		}

		x, _ := rb.Poll()
		if x != i {
			t.Error("poll", i, "had value", x)
		}
	}
}

func TestRingBufferWrapping(t *testing.T) {
	rb := NewRingBuffer[int]()

	for i := 0; i < minArrayCap; i++ {
		rb.Push(i)
	}
	for i := 0; i < 3; i++ {
		rb.Poll()
		rb.Push(minArrayCap + i)
	}

	for i := 0; i < minArrayCap; i++ {
		v, _ := rb.Peek()
		if v != i+3 {
			t.Error("peek", i, "had value", v)
		}
		rb.Poll()
	}
}

func TestRingBufferLength(t *testing.T) {
	rb := NewRingBuffer[int]()

	if rb.Len() != 0 {
		t.Error("empty queue length not 0")
	}

	for i := 0; i < 1000; i++ {
		rb.Push(i)
		if rb.Len() != i+1 {
			t.Error("adding: queue with", i, "elements has length", rb.Len())
		}
	}
	for i := 0; i < 1000; i++ {
		rb.Poll()
		if rb.Len() != 1000-i-1 {
			t.Error("removing: queue with", 1000-i-i, "elements has length", rb.Len())
		}
	}
}

func TestRingBufferAdd(t *testing.T) {
	list := NewRingBuffer[string]()
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

func calcBufferCap(n int) int {
	c := minArrayCap
	for c < n {
		c <<= 1
	}
	return c
}

func TestRingBufferGrow(t *testing.T) {
	list := NewRingBuffer[int]()

	for i := 0; i < 1000; i++ {
		list.Add(i)
		if l := list.Len(); l != i+1 {
			t.Errorf("list.Len() = %v, want %v", l, i+1)
		}

		wc := calcBufferCap(list.Len())
		ac := list.Cap()
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

func TestRingBufferIndex(t *testing.T) {
	list := NewRingBuffer[string]()

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

func TestRingBufferDelete(t *testing.T) {
	l := NewRingBuffer[int]()

	for i := 1; i <= 100; i++ {
		l.Add(i)
	}

	l.Delete(101)
	if l.Len() != 100 {
		t.Error("RingBuffer.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		l.Delete(i, i)
		if l.Len() != 100-i {
			t.Errorf("RingBuffer.Delete(%v) failed, l.Len() = %v, want %v", i, l.Len(), 100-i)
		}
	}

	if !l.IsEmpty() {
		t.Error("RingBuffer.IsEmpty() should return true")
	}
}

func TestRingBufferDelete2(t *testing.T) {
	l := NewRingBuffer[int]()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			l.Add(i)
		}
	}

	n := l.Len()
	l.Delete(100)
	if l.Len() != n {
		t.Errorf("RingBuffer.Delete(100).Len() = %v, want %v", l.Len(), n)
	}
	for i := 0; i < 100; i++ {
		n = l.Len()
		z := i % 10
		l.Delete(i)
		a := n - l.Len()
		if a != z {
			t.Errorf("RingBuffer.Delete(%v) = %v, want %v", i, a, z)
		}
	}

	if !l.IsEmpty() {
		t.Error("RingBuffer.IsEmpty() should return true")
	}
}

func TestRingBufferDeleteAll(t *testing.T) {
	l := NewRingBuffer[int]()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			l.Add(i)
		}
	}

	n := l.Len()
	l.Delete(100)
	if l.Len() != n {
		t.Errorf("RingBuffer.Delete(100).Len() = %v, want %v", l.Len(), n)
	}
	for i := 0; i < 100; i++ {
		n = l.Len()
		z := i % 10
		l.DeleteAll(NewRingBuffer(i, i))
		a := n - l.Len()
		if a != z {
			t.Errorf("RingBuffer.Delete(%v) = %v, want %v", i, a, z)
		}
	}

	if !l.IsEmpty() {
		t.Error("RingBuffer.IsEmpty() should return true")
	}
}

func TestRingBufferRemove(t *testing.T) {
	list := NewRingBuffer[string]()
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

func TestRingBufferRemovePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewRingBuffer("a")
	list.Remove(1)
}

func TestRingBufferGet(t *testing.T) {
	list := NewRingBuffer[string]()
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

func TestRingBufferGet2(t *testing.T) {
	rb := NewRingBuffer[int]()

	for i := 0; i < 1000; i++ {
		rb.Push(i)
		for j := 0; j < rb.Len(); j++ {
			v := rb.Get(j)
			if v != j {
				t.Errorf("[%d] index %d = %d, want %d", i, j, v, j)
			}
		}
	}
}

func TestRingBufferGetNegative(t *testing.T) {
	rb := NewRingBuffer[int]()

	for i := 0; i < 1000; i++ {
		rb.Push(i)
		for j := 1; j <= rb.Len(); j++ {
			if rb.Get(-j) != rb.Len()-j {
				t.Errorf("index %d doesn't contain %d", -j, rb.Len()-j)
			}
		}
	}
}

func TestRingBufferGetOutOfRangePanics(t *testing.T) {
	rb := NewRingBuffer[int]()

	rb.Push(1, 2, 3)

	assertPanics(t, "should panic when negative index", func() {
		rb.Get(-4)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		rb.Get(4)
	})
}

func TestRingBufferGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewRingBuffer("a")
	list.Get(1)
}

func TestRingBufferSwap(t *testing.T) {
	list := NewRingBuffer[string]()
	list.Add("a")
	list.Add("b", "c")
	list.Swap(0, 1)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestRingBufferClear(t *testing.T) {
	list := NewRingBuffer[string]()
	list.Add("e", "f", "g", "a", "b", "c", "d")
	list.Clear()
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestRingBufferContains(t *testing.T) {
	list := NewRingBuffer[int]()

	a := []int{}
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
	if av := list.Contains(0); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.Contains(0, 1, 2); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestRingBufferRetain(t *testing.T) {
	for n := 0; n < 100; n++ {
		a := []int{}
		list := NewRingBuffer[int]()
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
			a = []int{}
			list.Retain()
			vs := list.Values()
			if len(vs) > 0 {
				t.Fatalf("%d Retain() = %v, want %v", n, vs, a)
			}
		}

		list.Clear()
		a = []int{}
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
			a = []int{}
			list.RetainAll(AsArrayList(a))
			vs := list.Values()
			if len(vs) > 0 {
				t.Fatalf("%d Retain() = %v, want %v", n, vs, a)
			}
		}
	}
}

func TestRingBufferValues(t *testing.T) {
	list := NewRingBuffer[string]()
	if av, ev := fmt.Sprintf("%v", list.Values()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}

	list.Add("a")
	list.Add("b", "c")
	if av, ev := fmt.Sprintf("%v", list.Values()), "[a b c]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestRingBufferInsert(t *testing.T) {
	list := NewRingBuffer[string]()
	list.Insert(0, "b", "c")
	list.Insert(0, "a")
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	list.Insert(3, "d") // append
	if av := list.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 4)
	}
	if av, ev := fmt.Sprintf("%v", list.Values()), "[a b c d]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestRingBufferInsertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewRingBuffer("a")
	list.Insert(2, "b")
}

func TestRingBufferSet(t *testing.T) {
	list := NewRingBuffer("", "")
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
	if av, ev := fmt.Sprintf("%v", list.Values()), "[a bb c]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestRingBufferSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewRingBuffer("a")
	list.Set(1, "b")
}

func TestRingBufferEach(t *testing.T) {
	list := NewRingBuffer[string]()
	list.Add("a", "b", "c")
	index := 0
	list.Each(func(value string) {
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

func TestRingBufferIteratorPrevOnEmpty(t *testing.T) {
	list := NewRingBuffer[int]()
	it := list.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestRingBufferIteratorNextOnEmpty(t *testing.T) {
	list := NewRingBuffer[int]()
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestRingBufferIteratorPrev(t *testing.T) {
	list := NewRingBuffer[string]()
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

func TestRingBufferIteratorNext(t *testing.T) {
	list := NewRingBuffer[string]()
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

func TestRingBufferIteratorReset(t *testing.T) {
	list := NewRingBuffer[string]()

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

func assertRingBufferIteratorRemove(t *testing.T, i int, it Iterator[int], w *RingBuffer[int]) int {
	v := it.Value()

	it.Remove()

	w.Delete(v)

	it.SetValue(9999)

	rb := it.(*ringBufferIterator[int]).rb
	if rb.Contains(v) {
		t.Fatalf("[%d] l.Contains(%v) = true", i, v)
	}

	if rb.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, rb.String(), w.String())
	}

	return v
}

func TestRingBufferIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewRingBuffer[int]()
		w := NewRingBuffer[int]()

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

		v := assertRingBufferIteratorRemove(t, i, it, w)

		it.Next()
		if v+1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		assertRingBufferIteratorRemove(t, i, it, w)

		it.Prev()
		if v-1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		assertRingBufferIteratorRemove(t, i, it, w)

		// remove first
		for it.Prev() {
		}
		assertRingBufferIteratorRemove(t, i, it, w)

		// remove last
		for it.Next() {
		}
		assertRingBufferIteratorRemove(t, i, it, w)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertRingBufferIteratorRemove(t, i, it, w)
			}
		} else {
			for it.Next() {
				assertRingBufferIteratorRemove(t, i, it, w)
			}
		}
		if !l.IsEmpty() {
			t.Fatalf("[%d] l.IsEmpty() = true", i)
		}
	}
}

func TestRingBufferIteratorSetValue(t *testing.T) {
	l := NewRingBuffer[int]()
	for i := 1; i <= 100; i++ {
		l.Add(i)
	}

	// forward
	for it := l.Iterator(); it.Next(); {
		it.SetValue(it.Value() + 100)
	}
	for i := 1; i <= l.Len(); i++ {
		v := l.Get(i - 1)
		w := i + 100
		if v != w {
			t.Fatalf("List[%d] = %v, want %v", i-1, v, w)
		}
	}

	// backward
	for it := l.Iterator(); it.Prev(); {
		it.SetValue(it.Value() + 100)
	}
	for i := 1; i <= l.Len(); i++ {
		v := l.Get(i - 1)
		w := i + 200
		if v != w {
			t.Fatalf("List[%d] = %v, want %v", i-1, v, w)
		}
	}
}

func TestRingBufferSort(t *testing.T) {
	for i := 1; i < 100; i++ {
		l := NewRingBuffer[int]()

		a := make([]int, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, rand.Intn(20))
		}

		for j := len(a) - 1; j >= 0; j-- {
			l.Add(a[j])
		}

		l.Sort(LessInt)
		sort.Sort(inta(a))

		if !reflect.DeepEqual(a, l.Values()) {
			t.Errorf("%v != %v", a, l.Values())
		}
	}
}

func checkRingBufferLen[T any](t *testing.T, l *RingBuffer[T], len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkRingBuffer(t *testing.T, l *RingBuffer[int], evs []int) {
	if !checkRingBufferLen(t, l, len(evs)) {
		return
	}

	for i, it := 0, l.Iterator(); it.Next(); i++ {
		v := it.Value()
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

func TestRingBufferExtending(t *testing.T) {
	l1 := NewRingBuffer(1, 2, 3)
	l2 := NewRingBuffer[int]()
	l2.Add(4)
	l2.Add(5)

	l3 := NewRingBuffer[int]()
	l3.AddAll(l1)
	checkRingBuffer(t, l3, []int{1, 2, 3})
	l3.AddAll(l2)
	checkRingBuffer(t, l3, []int{1, 2, 3, 4, 5})

	l3 = NewRingBuffer[int]()
	l3.PushHeadAll(l2)
	checkRingBuffer(t, l3, []int{4, 5})
	l3.PushHeadAll(l1)
	checkRingBuffer(t, l3, []int{1, 2, 3, 4, 5})

	checkRingBuffer(t, l1, []int{1, 2, 3})
	checkRingBuffer(t, l2, []int{4, 5})

	l3 = NewRingBuffer[int]()
	l3.PushTailAll(l1)
	checkRingBuffer(t, l3, []int{1, 2, 3})
	l3.PushTailAll(l3)
	checkRingBuffer(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = NewRingBuffer[int]()
	l3.PushHeadAll(l1)
	checkRingBuffer(t, l3, []int{1, 2, 3})
	l3.PushHeadAll(l3)
	checkRingBuffer(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = NewRingBuffer[int]()
	l1.PushTailAll(l3)
	checkRingBuffer(t, l1, []int{1, 2, 3})
	l1.PushHeadAll(l3)
	checkRingBuffer(t, l1, []int{1, 2, 3})

	l1.Clear()
	l2.Clear()
	l3.Clear()
	l1.PushTail(1, 2, 3)
	checkRingBuffer(t, l1, []int{1, 2, 3})
	l2.PushTail(4, 5)
	checkRingBuffer(t, l2, []int{4, 5})
	l3.PushTailAll(l1)
	checkRingBuffer(t, l3, []int{1, 2, 3})
	l3.PushTail(4, 5)
	checkRingBuffer(t, l3, []int{1, 2, 3, 4, 5})
	l3.PushHead(4, 5)
	checkRingBuffer(t, l3, []int{4, 5, 1, 2, 3, 4, 5})
}

func TestRingBufferQueue(t *testing.T) {
	q := NewRingBuffer[int]()

	if _, ok := q.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}

	for i := 0; i < 100; i++ {
		q.Push(i)
	}

	for i := 0; i < 100; i++ {
		v, _ := q.Peek()
		if v != i {
			t.Errorf("Peek(%d) = %v, want %v", i, v, i)
		}

		x, _ := q.Poll()
		if x != i {
			t.Errorf("Poll(%d) = %v, want %v", i, x, i)
		}
	}

	if _, ok := q.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}
}

func TestRingBufferDeque(t *testing.T) {
	q := NewRingBuffer[int]()

	if _, ok := q.PeekHead(); ok {
		t.Error("should return false when peeking empty queue")
	}

	if _, ok := q.PeekTail(); ok {
		t.Error("should return false when peeking empty queue")
	}

	ea := []int{}
	for i := 0; i < 100; i++ {
		if i&1 == 0 {
			q.PushHead(i)
			a := make([]int, len(ea)+1)
			a[0] = i
			copy(a[1:], ea)
			ea = a
		} else {
			q.PushTail(i)
			ea = append(ea, i)
		}

		vs := q.Values()
		if len(vs) != i+1 {
			t.Fatalf("(%d) = %v, want %v", i, len(vs), i+1)
		}
		if !reflect.DeepEqual(vs, ea) {
			t.Fatalf("(%d) = %v, want %v", i, vs, ea)
		}
	}

	for i := 0; i < 100; i++ {
		if i&1 == 0 {
			w := 100 - i - 2
			v, _ := q.PeekHead()
			if v != w {
				t.Fatalf("PeekHead(%d) = %v, want %v\n%v", i, v, w, q.Values())
			}

			x, _ := q.PollHead()
			if x != w {
				t.Fatalf("PeekHead(%d) = %v, want %v\n%v", i, x, w, q.Values())
			}
		} else {
			w := 100 - i
			v, _ := q.PeekTail()
			if v != w {
				t.Fatalf("PeekTail(%d) = %v, want %v\n%v", i, v, w, q.Values())
			}

			x, _ := q.PollTail()
			if x != w {
				t.Fatalf("PoolTail(%d) = %v, want %v\n%v", i, x, w, q.Values())
			}
		}
	}

	if _, ok := q.PollHead(); ok {
		t.Error("should return false when removing empty queue")
	}

	if _, ok := q.PollTail(); ok {
		t.Error("should return false when removing empty queue")
	}
}

func TestRingBufferJSON(t *testing.T) {
	cs := []struct {
		s string
		a *RingBuffer[string]
	}{
		{`[]`, NewRingBuffer[string]()},
		{`["a","b","c"]`, NewRingBuffer("a", "b", "c")},
	}

	for i, c := range cs {
		a := NewRingBuffer[string]()
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
