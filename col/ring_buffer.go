package col

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/askasoft/pango/bye"
)

// RingBuffer A fast Golang queue using a ring-buffer, based on the version suggested by Dariusz Górecki.
// Using this instead of other, simpler, queue implementations (slice+append or linked list) provides substantial memory and time benefits, and fewer GC pauses.
// The queue implemented here is as fast as it is in part because it is not thread-safe.
type RingBuffer struct {
	data            []T
	head, tail, len int
}

// NewRingBuffer constructs and returns a new RingBuffer.
// Example: NewRingBuffer(1, 2, 3)
func NewRingBuffer(vs ...T) *RingBuffer {
	size := doubleup(minArrayCap, len(vs))
	rb := &RingBuffer{
		data: make([]T, size),
	}

	rb.Pushs(vs...)
	return rb
}

// Cap returns the capcity of the buffer.
func (rb *RingBuffer) Cap() int {
	return len(rb.data)
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the number of elements currently stored in the buffer.
func (rb *RingBuffer) Len() int {
	return rb.len
}

// IsEmpty returns true if the container length == 0
func (rb *RingBuffer) IsEmpty() bool {
	return rb.len == 0
}

// Clear clears list al.
func (rb *RingBuffer) Clear() {
	rb.head, rb.tail, rb.len = 0, 0, 0
	rb.shrink()
}

// Add add item v.
func (rb *RingBuffer) Add(v T) {
	rb.Insert(rb.len, v)
}

// Adds adds all items of vs.
func (rb *RingBuffer) Adds(vs ...T) {
	rb.Inserts(rb.len, vs...)
}

// AddCol adds all items of another collection
func (rb *RingBuffer) AddCol(ac Collection) {
	rb.InsertCol(rb.len, ac)
}

// Remove remove all items with associated value v of vs
func (rb *RingBuffer) Remove(v T) {
	for i := rb.Index(v); i >= 0; i = rb.Index(v) {
		rb.DeleteAt(i)
	}
}

// Removes remove all items with associated value v of vs
func (rb *RingBuffer) Removes(vs ...T) {
	if rb.IsEmpty() {
		return
	}

	for _, v := range vs {
		rb.Remove(v)
	}
}

// RemoveFunc remove all items that function f returns true
func (rb *RingBuffer) RemoveFunc(f func(T) bool) {
	if rb.IsEmpty() {
		return
	}

	it := rb.Iterator()
	for it.Next() {
		if f(it.Value()) {
			it.Remove()
		}
	}
}

// RemoveCol remove all of this collection's elements that are also contained in the specified collection
func (rb *RingBuffer) RemoveCol(ac Collection) {
	if rb.IsEmpty() || ac.IsEmpty() {
		return
	}

	if rb == ac {
		rb.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			rb.Remove(it.Value())
		}
		return
	}

	rb.Removes(ac.Values()...)
}

// Contain Test to see if the list contains the value v
func (rb *RingBuffer) Contain(v T) bool {
	return rb.Index(v) >= 0
}

// Contains Test to see if the RingBuffer contains the value v
func (rb *RingBuffer) Contains(vs ...T) bool {
	if len(vs) == 0 {
		return true
	}

	if rb.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if rb.Index(v) < 0 {
			return false
		}
	}
	return true
}

// ContainCol Test to see if the collection contains all items of another collection
func (rb *RingBuffer) ContainCol(ac Collection) bool {
	if ac.IsEmpty() || rb == ac {
		return true
	}

	if rb.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if rb.Index(it.Value()) < 0 {
				return false
			}
		}
		return true
	}

	return rb.Contains(ac.Values()...)
}

// Retains Retains only the elements in this collection that are contained in the argument array vs.
func (rb *RingBuffer) Retains(vs ...T) {
	if rb.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		rb.Clear()
		return
	}

	it := rb.Iterator()
	for it.Next() {
		if !contains(vs, it.Value()) {
			it.Remove()
		}
	}
}

// RetainCol Retains only the elements in this collection that are contained in the specified collection.
func (rb *RingBuffer) RetainCol(ac Collection) {
	if rb.IsEmpty() || rb == ac {
		return
	}

	if ac.IsEmpty() {
		rb.Clear()
		return
	}

	it := rb.Iterator()
	for it.Next() {
		if !ac.Contain(it.Value()) {
			it.Remove()
		}
	}
}

// Values returns a slice contains all the items of the RingBuffer rb
func (rb *RingBuffer) Values() []T {
	if rb.len == 0 {
		return rb.data[:0]
	}

	if rb.head <= rb.tail {
		return rb.data[rb.head : rb.tail+1]
	}

	data := make([]T, len(rb.data))
	copy(data, rb.data[rb.head:])
	copy(data[rb.len-rb.tail-1:], rb.data[0:rb.tail+1])

	rb.head = 0
	rb.tail = rb.len - 1
	rb.data = data
	return rb.data[0:rb.len]
}

// Each call f for each item in the RingBuffer
func (rb *RingBuffer) Each(f func(T)) {
	if rb.head <= rb.tail {
		for i := rb.head; i <= rb.tail; i++ {
			f(rb.data[i])
		}
	} else {
		l := len(rb.data)
		for i := rb.head; i < l; i++ {
			f(rb.data[i])
		}
		for i := 0; i <= rb.tail; i++ {
			f(rb.data[i])
		}
	}
}

// ReverseEach call f for each item in the RingBuffer with reverse order
func (rb *RingBuffer) ReverseEach(f func(T)) {
	if rb.head <= rb.tail {
		for i := rb.tail; i >= rb.head; i-- {
			f(rb.data[i])
		}
	} else {
		l := len(rb.data)
		for i := rb.tail; i >= 0; i-- {
			f(rb.data[i])
		}
		for i := l - 1; i >= rb.head; i-- {
			f(rb.data[i])
		}
	}
}

// Iterator returns a iterator for the RingBuffer
func (rb *RingBuffer) Iterator() Iterator {
	return &ringBufferIterator{rb, -1, -1}
}

//-----------------------------------------------------------
// implements List interface

// Get returns the item at the specified position in this RingBuffer
// if i < -rb.Len() or i >= rb.Len(), panic
// if i < 0, returns rb.Get(rb.Len() + i)
func (rb *RingBuffer) Get(index int) T {
	index = rb.checkItemIndex(index)

	return rb.data[index]
}

// Set set the v at the specified index in this RingBuffer and returns the old value.
func (rb *RingBuffer) Set(index int, v T) (ov T) {
	index = rb.checkItemIndex(index)

	ov = rb.data[index]
	rb.data[index] = v
	return
}

// Insert insert value v at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than RingBuffer's size
// Note: position equal to RingBuffer's size is valid, i.e. append.
func (rb *RingBuffer) Insert(index int, v T) {
	idx := rb.checkSizeIndex(index)

	if rb.expand(1) {
		index = rb.checkSizeIndex(index)
	} else {
		index = idx
	}

	if rb.len == 0 {
		rb.data[0] = v
		rb.tail = 0
	} else if index == rb.tail+1 {
		l := len(rb.data)
		rb.tail++
		if rb.tail >= l {
			rb.tail -= l
		}
		rb.data[rb.tail] = v
	} else if index == rb.head {
		l := len(rb.data)
		rb.head--
		if rb.head < 0 {
			rb.head += l
		}
		rb.data[rb.head] = v
	} else if index > rb.head {
		head, tail := rb.head-1, rb.tail
		if head < 0 {
			tail -= head
			head = 0
		}
		if head != rb.head {
			copy(rb.data[head:rb.head], rb.data[rb.head:rb.head+rb.head-head])
		}
		if tail != rb.tail {
			x := tail - rb.tail
			for i, j := rb.tail, 0; j < x; i, j = i-1, j+1 {
				rb.data[i+x] = rb.data[i]
			}
		}
		rb.data[head] = v
		rb.head, rb.tail = head, tail
	} else {
		// 0 < index < tail < head
		for i, x := rb.tail, rb.tail-index+1; i >= index; i-- {
			rb.data[i+x] = rb.data[i]
		}
		rb.data[index] = v
	}

	rb.len++
}

// Inserts inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than RingBuffer's size
// Note: position equal to RingBuffer's size is valid, i.e. append.
func (rb *RingBuffer) Inserts(index int, vs ...T) {
	idx := rb.checkSizeIndex(index)

	n := len(vs)
	if n == 0 {
		return
	}

	if rb.expand(n) {
		index = rb.checkSizeIndex(index)
	} else {
		index = idx
	}

	if rb.len == 0 {
		copy(rb.data, vs)
		rb.tail = n - 1
	} else if index == rb.tail+1 {
		l := len(rb.data)
		rb.tail += n
		if rb.tail >= l {
			rb.tail -= l
			copy(rb.data[index:l], vs[:n-rb.tail-1])
			copy(rb.data, vs[n-rb.tail-1:])
		} else {
			copy(rb.data[index:], vs)
		}
	} else if index == rb.head {
		l := len(rb.data)
		rb.head -= n
		if rb.head < 0 {
			rb.head += l
			copy(rb.data, vs[n-index:])
			copy(rb.data[rb.head:], vs[:n-index])
		} else {
			copy(rb.data[rb.head:], vs)
		}
	} else if index > rb.head {
		head, tail := rb.head-n, rb.tail
		if head < 0 {
			tail -= head
			head = 0
		}
		if head != rb.head {
			copy(rb.data[head:rb.head], rb.data[rb.head:rb.head+rb.head-head])
		}
		if tail != rb.tail {
			x := tail - rb.tail
			for i, j := rb.tail, 0; j < x; i, j = i-1, j+1 {
				rb.data[i+x] = rb.data[i]
			}
		}
		copy(rb.data[head:head+n], vs)
		rb.head, rb.tail = head, tail
	} else {
		// 0 < index < tail < head
		for i, x := rb.tail, rb.tail-index+1; i >= index; i-- {
			rb.data[i+x] = rb.data[i]
		}
		copy(rb.data[index:], vs)
	}

	rb.len += n
}

// InsertCol inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than RingBuffer's size
// Note: position equal to RingBuffer's size is valid, i.e. append.
func (rb *RingBuffer) InsertCol(index int, ac Collection) {
	rb.Inserts(index, ac.Values()...)
}

// Index returns the index of the first occurrence of the specified v in this RingBuffer, or -1 if this RingBuffer does not contain v.
func (rb *RingBuffer) Index(v T) int {
	return rb.IndexFunc(func(d T) bool {
		return d == v
	})
}

// IndexFunc returns the index of the first true returned by function f in this list, or -1 if this list does not contain v.
func (rb *RingBuffer) IndexFunc(f func(T) bool) int {
	if rb.len == 0 {
		return -1
	}

	if rb.head <= rb.tail {
		for i := rb.head; i <= rb.tail; i++ {
			if f(rb.data[i]) {
				return i - rb.head
			}
		}
		return -1
	}

	for i := rb.head; i < rb.len; i++ {
		if f(rb.data[i]) {
			return i - rb.head
		}
	}
	for i := 0; i <= rb.tail; i++ {
		if f(rb.data[i]) {
			return i
		}
	}
	return -1
}

// DeleteAt remove the item at the specified position in this RingBuffer.
func (rb *RingBuffer) DeleteAt(index int) {
	index = rb.checkItemIndex(index)
	rb.removeAt(index)
	rb.shrink()
}

func (rb *RingBuffer) removeAt(index int) {
	rb.data[index] = nil
	rb.len--

	if rb.len == 0 {
		rb.head, rb.tail = 0, 0
	} else if index == rb.head {
		rb.head++
		if rb.head >= len(rb.data) {
			rb.head = 0
		}
	} else if index == rb.tail {
		rb.tail--
		if rb.tail < 0 {
			rb.tail = len(rb.data) - 1
		}
	} else if index > rb.head {
		if rb.head < rb.tail {
			copy(rb.data[index:rb.tail], rb.data[index+1:rb.tail+1])
			rb.data[rb.tail] = nil
			rb.tail--
		} else {
			copy(rb.data[index:], rb.data[index+1:])
			rb.data[len(rb.data)-1] = rb.data[0]
			if rb.tail > 0 {
				copy(rb.data[0:rb.tail], rb.data[1:rb.tail+1])
			}
			rb.data[rb.tail] = nil
			rb.tail--
			if rb.tail < 0 {
				rb.tail = len(rb.data) - 1
			}
		}
	} else {
		// 0 < index < tail < head
		copy(rb.data[index:], rb.data[index+1:rb.tail+1])
		rb.data[rb.tail] = nil
		rb.tail--
	}
}

// Swap swaps values of two items at the given index.
func (rb *RingBuffer) Swap(i, j int) {
	i = rb.checkItemIndex(i)
	j = rb.checkItemIndex(j)

	if i != j {
		rb.data[i], rb.data[j] = rb.data[j], rb.data[i]
	}
}

// Sort Sorts this RingBuffer according to the order induced by the specified Comparator.
func (rb *RingBuffer) Sort(less Less) {
	if rb.len < 2 {
		return
	}
	sort.Sort(&sorter{rb, less})
}

// Head get the first item of RingBuffer.
func (rb *RingBuffer) Head() (v T) {
	v, _ = rb.PeekHead()
	return
}

// Tail get the last item of RingBuffer.
func (rb *RingBuffer) Tail() (v T) {
	v, _ = rb.PeekTail()
	return
}

//--------------------------------------------------------------------
// implements Queue interface

// Peek get the first item of RingBuffer.
func (rb *RingBuffer) Peek() (v T, ok bool) {
	return rb.PeekHead()
}

// Poll get and remove the first item of RingBuffer.
func (rb *RingBuffer) Poll() (T, bool) {
	return rb.PollHead()
}

// Push insert item v at the tail of RingBuffer rb.
func (rb *RingBuffer) Push(v T) {
	rb.Insert(rb.len, v)
}

// Pushs inserts all items of vs at the tail of RingBuffer rb.
func (rb *RingBuffer) Pushs(vs ...T) {
	rb.Inserts(rb.len, vs...)
}

// MustPeek Retrieves, but does not remove, the head of this queue, panic if this queue is empty.
func (rb *RingBuffer) MustPeek() T {
	if v, ok := rb.Peek(); ok {
		return v
	}

	panic("RingBuffer: MustPeek() called on empty queue")
}

// MustPoll Retrieves and removes the head of this queue, panic if this queue is empty.
func (rb *RingBuffer) MustPoll() T {
	if v, ok := rb.Poll(); ok {
		return v
	}

	panic("RingBuffer: MustPoll() called on empty queue")
}

//--------------------------------------------------------------------
// implements Deque interface

// PeekHead get the first item of RingBuffer.
func (rb *RingBuffer) PeekHead() (v T, ok bool) {
	if rb.IsEmpty() {
		return
	}

	v, ok = rb.data[rb.head], true
	return
}

// PeekTail get the last item of RingBuffer.
func (rb *RingBuffer) PeekTail() (v T, ok bool) {
	if rb.IsEmpty() {
		return
	}

	v, ok = rb.data[rb.tail], true
	return
}

// PollHead get and remove the first item of RingBuffer.
func (rb *RingBuffer) PollHead() (v T, ok bool) {
	v, ok = rb.PeekHead()
	if ok {
		rb.removeAt(rb.head)
		rb.shrink()
	}

	return
}

// PollTail get and remove the last item of RingBuffer.
func (rb *RingBuffer) PollTail() (v T, ok bool) {
	v, ok = rb.PeekTail()
	if ok {
		rb.removeAt(rb.tail)
		rb.shrink()
	}
	return
}

// PushHead insert item v at the head of RingBuffer rb.
func (rb *RingBuffer) PushHead(v T) {
	rb.Insert(0, v)
}

// PushHeads inserts all items of vs at the head of RingBuffer rb.
func (rb *RingBuffer) PushHeads(vs ...T) {
	rb.Inserts(0, vs...)
}

// PushHeadCol inserts a copy of another collection at the head of RingBuffer rb.
// The rb and ac may be the same. They must not be nil.
func (rb *RingBuffer) PushHeadCol(ac Collection) {
	rb.InsertCol(0, ac)
}

// PushTail insert item v at the tail of RingBuffer rb.
func (rb *RingBuffer) PushTail(v T) {
	rb.Insert(rb.len, v)
}

// PushTails inserts all items of vs at the tail of RingBuffer rb.
func (rb *RingBuffer) PushTails(vs ...T) {
	rb.Inserts(rb.len, vs...)
}

// PushTailCol inserts a copy of another collection at the tail of RingBuffer rb.
// The rb and ac may be the same. They must not be nil.
func (rb *RingBuffer) PushTailCol(ac Collection) {
	rb.InsertCol(rb.len, ac)
}

//-----------------------------------------------------------

// String print RingBuffer to string
func (rb *RingBuffer) String() string {
	bs, _ := json.Marshal(rb)
	return bye.UnsafeString(bs)
}

//-----------------------------------------------------------

// expand expand the buffer to guarantee space for n more elements.
func (rb *RingBuffer) expand(x int) bool {
	c := len(rb.data)
	if rb.len+x <= c {
		return false
	}

	c = doubleup(c, c+x)
	rb.resize(c)
	return true
}

// resize down if data is less than 1/4 full.
func (rb *RingBuffer) shrink() {
	if len(rb.data) > minArrayCap && (rb.len<<2) == len(rb.data) {
		rb.resize(rb.len)
	}
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than 1/4 full
func (rb *RingBuffer) resize(n int) {
	data := make([]T, n)

	if rb.len > 0 {
		if rb.head <= rb.tail {
			copy(data, rb.data[rb.head:rb.tail+1])
		} else {
			n := copy(data, rb.data[rb.head:])
			copy(data[n:], rb.data[:rb.tail+1])
		}
	}

	rb.head = 0
	rb.tail = rb.len - 1
	rb.data = data
}

func (rb *RingBuffer) checkItemIndex(index int) int {
	if index >= rb.len || index < -rb.len {
		panic(fmt.Sprintf("RingBuffer out of bounds: index=%d, len=%d", index, rb.len))
	}

	if index < 0 {
		index += rb.len
	}

	index += rb.head
	sz := len(rb.data)
	if index >= sz {
		index -= sz
	}

	return index
}

func (rb *RingBuffer) checkSizeIndex(index int) int {
	if index > rb.len || index < -rb.len {
		panic(fmt.Sprintf("RingBuffer out of bounds: index=%d, len=%d", index, rb.len))
	}

	if index < 0 {
		index += rb.len
	}

	index += rb.head
	sz := len(rb.data)
	if index > sz {
		index -= sz
	}

	return index
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (rb *RingBuffer) addJSONArrayItem(v T) jsonArray {
	rb.Add(v)
	return rb
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(rb)
func (rb *RingBuffer) MarshalJSON() ([]byte, error) {
	return jsonMarshalArray(rb)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, rb)
func (rb *RingBuffer) UnmarshalJSON(data []byte) error {
	rb.Clear()
	return jsonUnmarshalArray(data, rb)
}
