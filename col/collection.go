package col

import "fmt"

// T element type
type T = any

// Sortable a value each interface for collection
type Sortable interface {
	// Sorts this container according to the order induced by the specified Comparator.
	Sort(less Less)
}

// Container the base container interface
type Container interface {
	// Len returns the length of the container.
	Len() int

	// IsEmpty returns true if the container length == 0
	IsEmpty() bool

	// Clear clears the container
	Clear()
}

// Collection the base collection interface
type Collection interface {
	Container

	// Add adds item v to the collection
	Add(v T)

	// Adds adds items of vs to the collection
	Adds(vs ...T)

	// AddCol adds all items of another collection
	AddCol(ac Collection)

	// Remove remove all items with associated value v
	Remove(v T)

	// Removes remove all items in the array vs
	Removes(vs ...T)

	// RemoveCol remove all of this collection's elements that are also contained in the specified collection
	RemoveCol(ac Collection)

	// RemoveIter remove all items in the iterator it
	RemoveIter(it Iterator)

	// RemoveFunc remove all items that function f returns true
	RemoveFunc(f func(T) bool)

	// Contain Test to see if the collection contains item v
	Contain(v T) bool

	// Contains Test to see if the collection contains all items of vs
	Contains(vs ...T) bool

	// ContainCol Test to see if the collection contains all items of another collection
	ContainCol(ac Collection) bool

	// ContainIter Test to see if the collection contains all items of iterator 'it'
	ContainIter(it Iterator) bool

	// Retains Retains only the elements in this collection that are contained in the argument array vs.
	Retains(vs ...T)

	// RetainCol Retains only the elements in this collection that are contained in the specified collection.
	RetainCol(ac Collection)

	// RetainFunc Retains all items that function f returns true
	RetainFunc(f func(T) bool)

	// Values returns a slice contains all the items of the collection
	Values() []T

	Eachable
}

// Set a set interface
type Set = Collection

// List a doubly linked list interface
type List interface {
	Collection

	ReverseEachable

	Iterable

	// Get returns the value at the specified index in this list. If the index is
	// invalid, the call will panic. This method accepts both positive and
	// negative index values. Index 0 refers to the first element, and
	// index -1 refers to the last.
	Get(index int) T

	// Set set the v at the specified index in this list and returns the old value.
	Set(index int, v T) T

	// Insert insert item v at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
	// Does not do anything if position is bigger than list's size
	// Note: position equal to list's size is valid, i.e. append.
	Insert(index int, v T)

	// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
	// Does not do anything if position is bigger than list's size
	// Note: position equal to list's size is valid, i.e. append.
	Inserts(index int, vs ...T)

	// InsertCol inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
	// Does not do anything if position is bigger than list's size
	// Note: position equal to list's size is valid, i.e. append.
	InsertCol(index int, ac Collection)

	// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
	Index(v T) int

	// IndexFunc returns the index of the first true returned by function f in this list, or -1 if this list does not contain v.
	IndexFunc(f func(T) bool) int

	// DeleteAt remove the item at the specified position in this list
	DeleteAt(index int)

	// Swap swaps values of two items at the given index.
	Swap(i, j int)
}

// K KEY type
type K = any

// V VALUE type
type V = any

// P key/value pair
type P struct {
	Key   K
	Value V
}

func (p *P) String() string {
	return fmt.Sprintf("%v = %v", p.Key, p.Value)
}

// KV create key/value pair
func KV(k K, v V) P {
	return P{Key: k, Value: v}
}

// Map map interface
type Map interface {
	Container

	// Get looks for the given key, and returns the value associated with it,
	// or nil if not found. The boolean it returns says whether the key is ok in the map.
	Get(key K) (V, bool)

	// MustGet looks for the given key, and returns the value associated with it.
	// If not found, return defaults[0] or panic if defaults is not supplied.
	MustGet(key K, defaults ...V) V

	// Set sets the paired key-value items, and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	// Example: lm.Set("k1", "v1", "k2", "v2")
	Set(key K, value V) (ov V, ok bool)

	// SetIfAbsent sets the key-value item if the key does not exists in the map,
	// and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	SetIfAbsent(key K, value V) (ov V, ok bool)

	// SetEntries set items from key-value items array, override the existing items
	SetEntries(pairs ...P)

	// Copy copy items from another map am, override the existing items
	Copy(am Map)

	// Remove remove the item with key k,
	// and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	Remove(k K) (ov V, ok bool)

	// Removes remove all items with key of ks.
	Removes(ks ...K)

	// Contain looks for the given key, and returns true if the key exists in the map.
	Contain(k K) bool

	// Contains looks for the given key, and returns true if the key exists in the map.
	Contains(ks ...K) bool

	// Keys returns the key slice
	Keys() []K

	// Values returns a slice contains all the items of the collection
	Values() []V

	// Entries returns the key-value pair slice
	Entries() []P

	Eachable2
}

// IterableMap a iterable map interface
type IterableMap interface {
	Map

	ReverseEachable2

	Iterable2
}

// Queue A queue interface
type Queue interface {
	// Peek Retrieves, but does not remove, the head of this queue, or returns (nil, false) if this queue is empty.
	Peek() (T, bool)

	// Poll Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
	Poll() (T, bool)

	// Push add item v to the tail of queue
	Push(v T)

	// Pushs adds items of vs to the tail of queue
	Pushs(vs ...T)
}

// Deque A linear collection that supports element insertion and removal at both ends.
type Deque interface {
	// PeekHead Retrieves, but does not remove, the head of this queue, or returns (nil, false) if this queue is empty.
	PeekHead() (T, bool)

	// PollHead Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
	PollHead() (T, bool)

	// PushHead add item v to the head of queue
	PushHead(v T)

	// PushHeads adds items of vs to the head of queue
	PushHeads(vs ...T)

	// PeekTail Retrieves, but does not remove, the tail of this queue, or returns (nil, false) if this queue is empty.
	PeekTail() (T, bool)

	// PollTail Retrieves and removes the tail of this queue, or returns (nil, false) if this queue is empty.
	PollTail() (T, bool)

	// PushTail add item v to the tail of queue
	PushTail(v T)

	// PushTails adds items of vs to the tail of queue
	PushTails(vs ...T)
}
