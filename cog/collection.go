//go:build go1.18
// +build go1.18

package cog

// Sortable a value each interface for collection
type Sortable[T any] interface {
	// Sorts this container according to the order induced by the specified Comparator.
	Sort(less Less[T])
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
type Collection[T any] interface {
	Container

	// Add adds items of vs to the collection
	Add(vs ...T)

	// AddAll adds all items of another collection
	AddAll(ac Collection[T])

	// Delete delete all items of vs
	Delete(vs ...T)

	// DeleteIf delete all items that function f returns true
	DeleteIf(f func(T) bool)

	// DeleteAll delete all of this collection's elements that are also contained in the specified collection
	DeleteAll(ac Collection[T])

	// Contains Test to see if the collection contains all items of vs
	Contains(vs ...T) bool

	// ContainsAll Test to see if the collection contains all items of another collection
	ContainsAll(ac Collection[T]) bool

	// Retain Retains only the elements in this collection that are contained in the argument array vs.
	Retain(vs ...T)

	// RetainAll Retains only the elements in this collection that are contained in the specified collection.
	RetainAll(ac Collection[T])

	// Values returns a slice contains all the items of the collection
	Values() []T

	Eachable[T]
}

// Set a set interface
type Set[T any] Collection[T]

// List a doubly linked list interface
type List[T any] interface {
	Collection[T]

	ReverseEachable[T]

	Iterable[T]

	// Get returns the value at the specified index in this list. If the index is
	// invalid, the call will panic. This method accepts both positive and
	// negative index values. Index 0 refers to the first element, and
	// index -1 refers to the last.
	Get(index int) T

	// Set set the v at the specified index in this list and returns the old value.
	Set(index int, v T) T

	// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
	// Does not do anything if position is bigger than list's size
	// Note: position equal to list's size is valid, i.e. append.
	Insert(index int, vs ...T)

	// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
	// Does not do anything if position is bigger than list's size
	// Note: position equal to list's size is valid, i.e. append.
	InsertAll(index int, ac Collection[T])

	// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
	Index(v T) int

	// IndexIf returns the index of the first true returned by function f in this list, or -1 if this list does not contain v.
	IndexIf(f func(T) bool) int

	// Remove delete the item at the specified position in this list
	Remove(index int)

	// Swap swaps values of two items at the given index.
	Swap(i, j int)
}

// P key/value pair
type P[K any, V any] struct {
	Key   K
	Value V
}

// Map map interface
type Map[K any, V any] interface {
	Container

	// Get looks for the given key, and returns the value associated with it,
	// or nil if not found. The boolean it returns says whether the key is ok in the map.
	Get(key K) (V, bool)

	// Set sets the paired key-value items, and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	// Example: lm.Set("k1", "v1", "k2", "v2")
	Set(key K, value V) (ov V, ok bool)

	// SetPairs set items from key-value items array, override the existing items
	SetPairs(pairs ...P[K, V])

	// SetAll set items from another map am, override the existing items
	SetAll(am Map[K, V])

	// SetIfAbsent sets the key-value item if the key does not exists in the map,
	// and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	SetIfAbsent(key K, value V) (ov V, ok bool)

	// Delete delete all items with key of ks,
	// and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	Delete(ks ...K) (ov V, ok bool)

	// Contains looks for the given key, and returns true if the key exists in the map.
	Contains(ks ...K) bool

	// Keys returns the key slice
	Keys() []K

	// Values returns a slice contains all the items of the collection
	Values() []V

	Eachable2[K, V]
}

// IterableMap a iterable map interface
type IterableMap[K any, V any] interface {
	Map[K, V]

	ReverseEachable2[K, V]

	Iterable2[K, V]
}

// Queue A queue interface
type Queue[T any] interface {
	// Peek Retrieves, but does not remove, the head of this queue, or returns (nil, false) if this queue is empty.
	Peek() (T, bool)

	// Poll Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
	Poll() (T, bool)

	// Push adds items of vs to the tail of queue
	Push(vs ...T)
}

// Deque A linear collection that supports element insertion and removal at both ends.
type Deque[T any] interface {
	// PeekHead Retrieves, but does not remove, the head of this queue, or returns (nil, false) if this queue is empty.
	PeekHead() (T, bool)

	// PollHead Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
	PollHead() (T, bool)

	// PushHead adds items of vs to the head of queue
	PushHead(vs ...T)

	// PeekTail Retrieves, but does not remove, the tail of this queue, or returns (nil, false) if this queue is empty.
	PeekTail() (T, bool)

	// PollTail Retrieves and removes the tail of this queue, or returns (nil, false) if this queue is empty.
	PollTail() (T, bool)

	// PushTail adds items of vs to the tail of queue
	PushTail(vs ...T)
}
