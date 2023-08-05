package tags

type Iterator interface {
	// Next moves the iterator to the next item and returns true if there was a next item in the collection.
	// If Next() returns true, then next item's value can be retrieved by Value().
	// If Next() was called for the first time, then it will point the iterator to the first item if it exists.
	// Modifies the state of the iterator.
	Next() bool

	// Key returns the current item's key.
	Key() string

	// Value returns the current item's value.
	Value() string

	// Reset resets the iterator to its initial state (one-before-first/one-after-last)
	Reset()
}
