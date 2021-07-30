package col

// Eachable a value each interface for collection
type Eachable interface {
	// Each call f for each item in the collection
	Each(f func(value interface{}))
}

// ReverseEachable a value each interface for collection
type ReverseEachable interface {
	// ReverseEach call f for each item in the collection with reverse order
	ReverseEach(f func(value interface{}))
}

// Iterable a value iterable interface for collection
type Iterable interface {
	// Iterator returns a iterator for collection
	Iterator() Iterator
}

// Iterator is stateful iterator for collection.
type Iterator interface {
	// Prev moves the iterator to the previous element and returns true if there was a previous element in collection.
	// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
	// Modifies the state of the iterator.
	Prev() bool

	// Next moves the iterator to the next element and returns true if there was a next element in the collection.
	// If Next() returns true, then next element's value can be retrieved by Value().
	// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
	// Modifies the state of the iterator.
	Next() bool

	// Value returns the current element's value.
	Value() interface{}

	// SetValue set the value to the item
	SetValue(v interface{})

	// Remove remove the current element
	Remove()

	// Reset resets the iterator to its initial state (one-before-first/one-after-last)
	// Call Next()/Prev() to fetch the first/last element if any.
	Reset()
}

// Eachable2 a key/value each interface for collection
type Eachable2 interface {
	// Each call f for each key/value in the collection
	Each(f func(key interface{}, value interface{}))
}

// ReverseEachable2 a key/value reverse each interface for collection
type ReverseEachable2 interface {
	// ReverseEach call f for each key/value in the collection with reverse order
	ReverseEach(f func(key interface{}, value interface{}))
}

// Iterable2 a key/value iterable interface for collection
type Iterable2 interface {
	// Iterator returns a iterator for collection
	Iterator() Iterator2
}

// Iterator2 is stateful iterator for Key/Value paire.
type Iterator2 interface {
	Iterator

	// Key returns the current element's key.
	Key() interface{}
}
