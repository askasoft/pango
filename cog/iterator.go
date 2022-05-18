package cog

// Eachable a value each interface for collection
type Eachable[T any] interface {
	// Each call f for each item in the collection
	Each(f func(value T))
}

// ReverseEachable a value each interface for collection
type ReverseEachable[T any] interface {
	// ReverseEach call f for each item in the collection with reverse order
	ReverseEach(f func(value T))
}

// Iterable a value iterable interface for collection
type Iterable[T any] interface {
	// Iterator returns a iterator for collection
	Iterator() Iterator[T]
}

// Iterator is stateful iterator for collection.
type Iterator[T any] interface {
	// Prev moves the iterator to the previous item and returns true if there was a previous item in collection.
	// If Prev() returns true, then previous item's index and value can be retrieved by Index() and Value().
	// Modifies the state of the iterator.
	Prev() bool

	// Next moves the iterator to the next item and returns true if there was a next item in the collection.
	// If Next() returns true, then next item's value can be retrieved by Value().
	// If Next() was called for the first time, then it will point the iterator to the first item if it exists.
	// Modifies the state of the iterator.
	Next() bool

	// Value returns the current item's value.
	Value() T

	// SetValue set the current item's value.
	SetValue(v T)

	// Remove remove the current item
	Remove()

	// Reset resets the iterator to its initial state (one-before-first/one-after-last)
	// Call Next()/Prev() to fetch the first/last item if any.
	Reset()
}

// Eachable2 a key/value each interface for collection
type Eachable2[K any, V any] interface {
	// Each call f for each key/value in the collection
	Each(f func(key K, value V))
}

// ReverseEachable2 a key/value reverse each interface for collection
type ReverseEachable2[K any, V any] interface {
	// ReverseEach call f for each key/value in the collection with reverse order
	ReverseEach(f func(key K, value V))
}

// Iterable2 a key/value iterable interface for collection
type Iterable2[K any, V any] interface {
	// Iterator returns a iterator for collection
	Iterator() Iterator2[K, V]
}

// Iterator2 is stateful iterator for Key/Value paire.
type Iterator2[K any, V any] interface {
	Iterator[V]

	// Key returns the current item's key.
	Key() K
}
