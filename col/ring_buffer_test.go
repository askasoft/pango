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

func TestRingBufferInterface(t *testing.T) {
	var l List = NewRingBuffer()
	if l == nil {
		t.Error("RingBuffer is not a List")
	}

	var q Queue = NewRingBuffer()
	if q == nil {
		t.Error("RingBuffer is not a Queue")
	}

	var dq Queue = NewRingBuffer()
	if dq == nil {
		t.Error("RingBuffer is not a Deque")
	}

	var s Sortable = NewRingBuffer()
	if s == nil {
		t.Error("RingBuffer is not a Sortable")
	}
}

func TestRingBufferNew(t *testing.T) {
	rb1 := NewRingBuffer()

	if av := rb1.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	rb2 := NewRingBuffer(1, "b")

	if av := rb2.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}

	if av := rb2.Get(0); av != 1 {
		t.Errorf("Got %v expected %v", av, 1)
	}

	if av := rb2.Get(1); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestRingBufferSimple(t *testing.T) {
	rb := NewRingBuffer()

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
			t.Error("remove", i, "had value", x)
		}
	}
}

func TestRingBufferWrapping(t *testing.T) {
	rb := NewRingBuffer()

	for i := 0; i < minArrayCap; i++ {
		rb.Push(i)
	}
	for i := 0; i < 3; i++ {
		rb.Poll()
		rb.Push(minArrayCap + i)
	}

	for i := 0; i < minArrayCap; i++ {
		v, _ := rb.Peek()
		if v.(int) != i+3 {
			t.Error("peek", i, "had value", v)
		}
		rb.Poll()
	}
}

func TestRingBufferLength(t *testing.T) {
	rb := NewRingBuffer()

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
	list := NewRingBuffer()
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
	list := NewRingBuffer()

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
	list := NewRingBuffer()

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
	l := NewRingBuffer()

	for i := 1; i <= 100; i++ {
		l.Add(i)
	}

	l.Delete(101)
	if l.Len() != 100 {
		t.Error("RingBuffer.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		l.Delete(i)
		if l.Len() != 100-i {
			t.Errorf("RingBuffer.Delete(%v) failed, l.Len() = %v, want %v", i, l.Len(), 100-i)
		}
	}

	if !l.IsEmpty() {
		t.Error("RingBuffer.IsEmpty() should return true")
	}
}

func TestRingBufferDeleteAll(t *testing.T) {
	l := NewRingBuffer()

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

func TestRingBufferRemove(t *testing.T) {
	list := NewRingBuffer()
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
	list := NewRingBuffer()
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
	rb := NewRingBuffer()

	for i := 0; i < 1000; i++ {
		rb.Push(i)
		for j := 0; j < rb.Len(); j++ {
			v := rb.Get(j).(int)
			if v != j {
				t.Errorf("[%d] index %d = %d, want %d", i, j, v, j)
			}
		}
	}
}

func TestRingBufferGetNegative(t *testing.T) {
	rb := NewRingBuffer()

	for i := 0; i < 1000; i++ {
		rb.Push(i)
		for j := 1; j <= rb.Len(); j++ {
			if rb.Get(-j).(int) != rb.Len()-j {
				t.Errorf("index %d doesn't contain %d", -j, rb.Len()-j)
			}
		}
	}
}

func TestRingBufferGetOutOfRangePanics(t *testing.T) {
	rb := NewRingBuffer()

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
	list := NewRingBuffer()
	list.Add("a")
	list.Add("b", "c")
	list.Swap(0, 1)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestRingBufferClear(t *testing.T) {
	list := NewRingBuffer()
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
	list := NewRingBuffer()
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

func TestRingBufferValues(t *testing.T) {
	list := NewRingBuffer()
	if av, ev := fmt.Sprintf("%v", list.Values()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}

	list.Add("a")
	list.Add("b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestRingBufferInsert(t *testing.T) {
	list := NewRingBuffer()
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
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abbc"; av != ev {
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
	list := NewRingBuffer()
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

func TestRingBufferIteratorPrevOnEmpty(t *testing.T) {
	list := NewRingBuffer()
	it := list.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestRingBufferIteratorNextOnEmpty(t *testing.T) {
	list := NewRingBuffer()
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestRingBufferIteratorPrev(t *testing.T) {
	list := NewRingBuffer()
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
	list := NewRingBuffer()
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
	list := NewRingBuffer()

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

func assertRingBufferIteratorRemove(t *testing.T, i int, it Iterator, w *RingBuffer) int {
	v := it.Value()

	it.Remove()

	w.Delete(v)

	it.SetValue(9999)

	rb := it.(*ringBufferIterator).rb
	if rb.Contains(v) {
		t.Fatalf("[%d] l.Contains(%v) = true", i, v)
	}

	if rb.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, rb.String(), w.String())
	}

	return v.(int)
}

func TestRingBufferIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewRingBuffer()
		w := NewRingBuffer()

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
	l := NewRingBuffer()
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

func TestRingBufferSort(t *testing.T) {
	for i := 1; i < 100; i++ {
		l := NewRingBuffer()

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

func checkRingBufferLen(t *testing.T, l *RingBuffer, len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkRingBuffer(t *testing.T, l *RingBuffer, evs []interface{}) {
	if !checkRingBufferLen(t, l, len(evs)) {
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

func TestRingBufferExtending(t *testing.T) {
	l1 := NewRingBuffer(1, 2, 3)
	l2 := NewRingBuffer()
	l2.Add(4)
	l2.Add(5)

	l3 := NewRingBuffer()
	l3.AddAll(l1)
	checkRingBuffer(t, l3, []interface{}{1, 2, 3})
	l3.AddAll(l2)
	checkRingBuffer(t, l3, []interface{}{1, 2, 3, 4, 5})

	l3 = NewRingBuffer()
	l3.PushHeadAll(l2)
	checkRingBuffer(t, l3, []interface{}{4, 5})
	l3.PushHeadAll(l1)
	checkRingBuffer(t, l3, []interface{}{1, 2, 3, 4, 5})

	checkRingBuffer(t, l1, []interface{}{1, 2, 3})
	checkRingBuffer(t, l2, []interface{}{4, 5})

	l3 = NewRingBuffer()
	l3.PushTailAll(l1)
	checkRingBuffer(t, l3, []interface{}{1, 2, 3})
	l3.PushTailAll(l3)
	checkRingBuffer(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = NewRingBuffer()
	l3.PushHeadAll(l1)
	checkRingBuffer(t, l3, []interface{}{1, 2, 3})
	l3.PushHeadAll(l3)
	checkRingBuffer(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = NewRingBuffer()
	l1.PushTailAll(l3)
	checkRingBuffer(t, l1, []interface{}{1, 2, 3})
	l1.PushHeadAll(l3)
	checkRingBuffer(t, l1, []interface{}{1, 2, 3})

	l1.Clear()
	l2.Clear()
	l3.Clear()
	l1.PushTail(1, 2, 3)
	checkRingBuffer(t, l1, []interface{}{1, 2, 3})
	l2.PushTail(4, 5)
	checkRingBuffer(t, l2, []interface{}{4, 5})
	l3.PushTailAll(l1)
	checkRingBuffer(t, l3, []interface{}{1, 2, 3})
	l3.PushTail(4, 5)
	checkRingBuffer(t, l3, []interface{}{1, 2, 3, 4, 5})
	l3.PushHead(4, 5)
	checkRingBuffer(t, l3, []interface{}{4, 5, 1, 2, 3, 4, 5})
}

func TestRingBufferContains2(t *testing.T) {
	l := NewRingBuffer(1, 11, 111, "1", "11", "111")

	n := (100+1)/101 + 110

	if !l.Contains(n) {
		t.Errorf("RingBuffer [%v] should contains %v", l, n)
	}

	n++
	if l.Contains(n) {
		t.Errorf("RingBuffer [%v] should not contains %v", l, n)
	}

	s := str.Repeat("1", 3)

	if !l.Contains(s) {
		t.Errorf("RingBuffer [%v] should contains %v", l, s)
	}

	s += "0"
	if l.Contains(s) {
		t.Errorf("RingBuffer [%v] should not contains %v", l, s)
	}
}

func TestRingBufferQueue(t *testing.T) {
	q := NewRingBuffer()

	if _, ok := q.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}

	for i := 0; i < 100; i++ {
		q.Push(i)
	}

	for i := 0; i < 100; i++ {
		v, _ := q.Peek()
		if v.(int) != i {
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
	q := NewRingBuffer()

	if _, ok := q.PeekHead(); ok {
		t.Error("should return false when peeking empty queue")
	}

	if _, ok := q.PeekTail(); ok {
		t.Error("should return false when peeking empty queue")
	}

	ea := []T{}
	for i := 0; i < 100; i++ {
		if i&1 == 0 {
			q.PushHead(i)
			a := make([]T, len(ea)+1)
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
		a *RingBuffer
	}{
		{`[]`, NewRingBuffer()},
		{`["a","b","c"]`, NewRingBuffer("a", "b", "c")},
	}

	for i, c := range cs {
		a := NewRingBuffer()
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
