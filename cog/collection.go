package cog

import "fmt"

// Compare will make type assertion (see CompareString(a,b) for example),
// which will panic if a or b are not of the asserted type.
//
// Should return a int:
//
//	negative , if a < b
//	zero     , if a == b
//	positive , if a > b
type Compare[T any] func(a, b T) int

// Less will make type assertion (see LessString(a,b) for example),
// which will panic if a or b are not of the asserted type.
//
// Should return a bool:
//
//	true , if a < b
//	false, if a >= b
type Less[T any] func(a, b T) bool

type SortIF[T any] interface {
	// Len returns the length of the collection.
	Len() int

	// Get returns the value at the specified index in this list
	Get(index int) T

	// Swap swaps values of two items at the given index.
	Swap(i, j int)
}

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

	// Add adds item v to the collection
	Add(v T)

	// AddAll adds items of vs to the collection
	AddAll(vs ...T)

	// AddCol adds all items of another collection
	AddCol(ac Collection[T])

	// Remove remove all items with associated value v
	Remove(v T)

	// RemoveAll remove all items in the array vs
	RemoveAll(vs ...T)

	// RemoveCol remove all of this collection's elements that are also contained in the specified collection
	RemoveCol(ac Collection[T])

	// RemoveIter remove all items in the iterator it
	RemoveIter(it Iterator[T])

	// RemoveFunc remove all items that function f returns true
	RemoveFunc(f func(T) bool)

	// Contains Test to see if the collection contains item v
	Contains(v T) bool

	// ContainsAll Test to see if the collection contains all items of vs
	ContainsAll(vs ...T) bool

	// ContainsCol Test to see if the collection contains all items of another collection
	ContainsCol(ac Collection[T]) bool

	// ContainsIter Test to see if the collection contains all items of iterator 'it'
	ContainsIter(it Iterator[T]) bool

	// RetainAll Retains only the elements in this collection that are contained in the argument array vs.
	RetainAll(vs ...T)

	// RetainCol Retains only the elements in this collection that are contained in the specified collection.
	RetainCol(ac Collection[T])

	// RetainFunc Retains all items that function f returns true
	RetainFunc(f func(T) bool)

	// Values returns a slice contains all the items of the collection
	Values() []T

	Eachable[int, T]
}

// Set a set interface
type Set[T any] Collection[T]

// List a doubly linked list interface
type List[T any] interface {
	Collection[T]

	ReverseEachable[int, T]

	Iterable[T]

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
	InsertCol(index int, ac Collection[T])

	// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
	Index(v T) int

	// IndexFunc returns the index of the first true returned by function f in this list, or -1 if this list does not contain v.
	IndexFunc(f func(T) bool) int

	// DeleteAt remove the item at the specified position in this list
	DeleteAt(index int)

	// Swap swaps values of two items at the given index.
	Swap(i, j int)
}

// P key/value pair
type P[K any, V any] struct {
	Key K
	Val V
}

func (p *P[K, V]) String() string {
	return fmt.Sprintf("%v = %v", p.Key, p.Val)
}

// KV create key/value pair
func KV[K any, V any](k K, v V) P[K, V] {
	return P[K, V]{Key: k, Val: v}
}

// Map map interface
type Map[K any, V any] interface {
	Container

	// Get looks for the given key, and returns the value associated with it,
	// or nil if not found. The boolean it returns says whether the key is ok in the map.
	Get(key K) (V, bool)

	// MustGet looks for the given key, and returns the value associated with it.
	// Panic if not found.
	MustGet(key K) V

	// SafeGet looks for the given key, and returns the value associated with it.
	// If not found, return defaults[0] or zero V.
	SafeGet(key K, defaults ...V) V

	// Set sets the paired key-value items, and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	// Example: lm.Set("k1", "v1", "k2", "v2")
	Set(key K, value V) (ov V, ok bool)

	// SetIfAbsent sets the key-value item if the key does not exists in the map,
	// and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	SetIfAbsent(key K, value V) (ov V, ok bool)

	// SetEntries set items from key-value items array, override the existing items
	SetEntries(pairs ...P[K, V])

	// Copy copy items from another map am, override the existing items
	Copy(am Map[K, V])

	// Remove remove the item with key k,
	// and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	Remove(k K) (ov V, ok bool)

	// RemoveAll remove all items with key of ks.
	RemoveAll(ks ...K)

	// Contains Test to see if the map contains the given key
	Contains(k K) bool

	// ContainsAll Test to see if the map contains all keys of ks
	ContainsAll(ks ...K) bool

	// Keys returns the key slice
	Keys() []K

	// Values returns a slice contains all the items of the collection
	Values() []V

	// Entries returns the key-value pair slice
	Entries() []P[K, V]

	Eachable[K, V]
}

// IterableMap A iterable map interface
type IterableMap[K any, V any] interface {
	Map[K, V]

	ReverseEachable[K, V]

	Iterable2[K, V]
}

// Stack A stack interface
type Stack[T any] interface {
	// Pop Retrieves and removes the tail of this queue, or returns (nil, false) if this stack is empty.
	Pop() (T, bool)

	// Push add item v to the tail of stack
	Push(v T)
}

// Queue A queue interface
type Queue[T any] interface {
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
type Deque[T any] interface {
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
