package col

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

	// Add adds items of vs
	Add(vs ...interface{})

	// AddAll adds all items of another collection
	AddAll(ac Collection)

	// Delete delete all items of vs
	Delete(vs ...interface{})

	// DeleteAll delete all of this collection's elements that are also contained in the specified collection
	DeleteAll(ac Collection)

	// Contains Test to see if the collection contains all items of vs
	Contains(vs ...interface{}) bool

	// ContainsAll Test to see if the collection contains all items of another collection
	ContainsAll(ac Collection) bool

	// Retain Retains only the elements in this collection that are contained in the argument array vs.
	Retain(vs ...interface{})

	// RetainAll Retains only the elements in this collection that are contained in the specified collection.
	RetainAll(ac Collection)

	// Values returns a slice contains all the items of the collection
	Values() []interface{}

	Eachable
}

// Set a set interface
type Set = Collection

// List a double linked list interface
type List interface {
	Collection

	ReverseEachable

	Iterable

	// Get returns the value at the specified index in this list
	Get(index int) interface{}

	// Set set the v at the specified index in this list and returns the old value.
	Set(index int, v interface{}) interface{}

	// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
	// Does not do anything if position is bigger than list's size
	// Note: position equal to list's size is valid, i.e. append.
	Insert(index int, vs ...interface{})

	// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
	// Does not do anything if position is bigger than list's size
	// Note: position equal to list's size is valid, i.e. append.
	InsertAll(index int, ac Collection)

	// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
	Index(v interface{}) int

	// Remove delete the item at the specified position in this list
	Remove(index int)

	// Swap swaps values of two items at the given index.
	Swap(i, j int)
}

// Map map interface
type Map interface {
	Container

	// Get looks for the given key, and returns the value associated with it,
	// or nil if not found. The boolean it returns says whether the key is ok in the map.
	Get(key interface{}) (interface{}, bool)

	// Set sets the paired key-value items, and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	// Example: lm.Set("k1", "v1", "k2", "v2")
	Set(kvs ...interface{}) (ov interface{}, ok bool)

	// SetAll set items from another map am, override the existing items
	SetAll(am Map)

	// SetIfAbsent sets the key-value item if the key does not exists in the map,
	// and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	// Example: lm.SetIfAbsent("k1", "v1", "k2", "v2")
	SetIfAbsent(kvs ...interface{}) (ov interface{}, ok bool)

	// Delete delete all items with key of ks,
	// and returns what `Get` would have returned
	// on that key prior to the call to `Set`.
	Delete(ks ...interface{}) (ov interface{}, ok bool)

	// Contains looks for the given key, and returns true if the key exists in the map.
	Contains(ks ...interface{}) bool

	// Keys returns the key slice
	Keys() []interface{}

	// Values returns a slice contains all the items of the collection
	Values() []interface{}

	Eachable2
}

// IterableMap a iterable map interface
type IterableMap interface {
	Map

	ReverseEachable2

	Iterable2
}
