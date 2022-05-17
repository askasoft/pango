package col

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

	// Add adds items of vs to the collection
	Add(vs ...T)

	// AddAll adds all items of another collection
	AddAll(ac Collection)

	// Delete delete all items of vs
	Delete(vs ...T)

	// DeleteAll delete all of this collection's elements that are also contained in the specified collection
	DeleteAll(ac Collection)

	// Contains Test to see if the collection contains all items of vs
	Contains(vs ...T) bool

	// ContainsAll Test to see if the collection contains all items of another collection
	ContainsAll(ac Collection) bool

	// Retain Retains only the elements in this collection that are contained in the argument array vs.
	Retain(vs ...T)

	// RetainAll Retains only the elements in this collection that are contained in the specified collection.
	RetainAll(ac Collection)

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

	// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
	// Does not do anything if position is bigger than list's size
	// Note: position equal to list's size is valid, i.e. append.
	Insert(index int, vs ...T)

	// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
	// Does not do anything if position is bigger than list's size
	// Note: position equal to list's size is valid, i.e. append.
	InsertAll(index int, ac Collection)

	// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
	Index(v T) int

	// Remove delete the item at the specified position in this list
	Remove(index int)

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

// Map map interface
type Map interface {
	Container

	// Get looks for the given key, and returns the value associated with it,
	// or nil if not found. The boolean it returns says whether the key is ok in the map.
	Get(key K) (K, bool)

	// Set sets the paired key-value items, and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	// Example: lm.Set("k1", "v1", "k2", "v2")
	Set(key K, value V) (ov V, ok bool)

	// SetPairs set items from key-value items array, override the existing items
	SetPairs(pairs ...P)

	// SetAll set items from another map am, override the existing items
	SetAll(am Map)

	// SetIfAbsent sets the key-value item if the key does not exists in the map,
	// and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	SetIfAbsent(key K, value V) (ov K, ok bool)

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

	// Push adds items of vs to the tail of queue
	Push(vs ...T)
}

// Deque A linear collection that supports element insertion and removal at both ends.
type Deque interface {
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
