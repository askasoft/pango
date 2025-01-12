 Pango Collection
=====================================================================

A Go Collection/Container package.


### Container
All data structures implement the container interface with the following methods:

```go
// Container the base container interface
type Container interface {
	// Len returns the length of the container.
	Len() int

	// IsEmpty returns true if the container length == 0
	IsEmpty() bool

	// Clear clears the container
	Clear()
}
```

Containers are either ordered or unordered. All ordered containers provide [stateful iterators](#iterator).

| **Structure**                   | **Ordered** | **Iterator** | **Sortable** |
| :------------------------------ | :---------: | :----------: | :----------: |
| [ArrayList](#arraylist)         |      Y      |      Y       |      Y       |
| [LinkedList](#linkedlist)       |             |      Y       |      Y       |
| [HashSet](#hashset)             |             |              |              |
| [LinkedHashSet](#linkedhashset) |      Y      |      Y       |      Y       |
| [TreeSet](#treeset)             |      Y      |      Y       |              |
| [HashMap](#hashmap)             |             |              |              |
| [LinkedHashMap](#linkedhashmap) |      Y      |      Y       |              |
| [TreeMap](#treemap)             |      Y      |      Y       |              |


### Collection
Base interface for [List](#list) and [Set](#set).

Extends [Container](#container) interface.


```go
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

	// RetainAll Retains only the elements in this collection that are contained in the argument array vs.
	RetainAll(vs ...T)

	// RetainCol Retains only the elements in this collection that are contained in the specified collection.
	RetainCol(ac Collection[T])

	// Values returns a slice contains all the items of the collection
	Values() []T

	Eachable[T]
}
```

### List

A list is a data structure that stores values and may have repeated values.

Extends [Collection](#collection) interface.

```go
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
```

#### ArrayList

A [list](#list) backed by a dynamic array that grows implicitly.

Implements [List](#list), [Iterator](#iterator) interfaces.

```go
package main

import (
	"cmp"
	"github.com/askasoft/pango/cog/arraylist"
)

func main() {
	list := arraylist.NewArrayList[string]()
	list.Add("a")                            // ["a"]
	list.AddAll("c", "b")                    // ["a","c","b"]
	list.Sort(cmp.Less[string])              // ["a","b","c"]
	_ = list.Get(0)                          // "a"
	_ = list.Get(100)                        // panic
	_ = list.Contains("a")                   // true
	_ = list.ContainsAll("a", "b", "c")      // true
	_ = list.ContainsAll("a", "b", "c", "d") // false
	list.Swap(0, 1)                          // ["b","a",c"]
	list.DeleteAt(2)                         // ["b","a"]
	list.DeleteAt(1)                         // ["b"]
	list.DeleteAt(0)                         // []
	_ = list.IsEmpty()                       // true
	_ = list.Len()                           // 0
	list.Add("a")                            // ["a"]
	list.Clear()                             // []
	list.Insert(0, "b")                      // ["b"]
	list.Insert(0, "a")                      // ["a","b"]
}
```

#### LinkedList

A [list](#list) where each element points to the next and previous elements in the list.

Implements [List](#list), [Iterator](#iterator) interfaces.

```go
package main

import (
	"cmp"
	"github.com/askasoft/pango/cog/linkedlist"
)

func main() {
	list := linkedlist.NewLinkedList[string]()
	list.Add("a")                            // ["a"]
	list.AddAll("c", "b")                    // ["a","c","b"]
	list.Sort(cmp.Less[string])              // ["a","b","c"]
	_ = list.Get(0)                          // "a"
	_ = list.Get(100)                        // panic
	_ = list.Contains("a")                   // true
	_ = list.ContainsAll("a", "b", "c")      // true
	_ = list.ContainsAll("a", "b", "c", "d") // false
	list.Swap(0, 1)                          // ["b","a",c"]
	list.DeleteAt(2)                         // ["b","a"]
	list.DeleteAt(1)                         // ["b"]
	list.DeleteAt(0)                         // []
	_ = list.IsEmpty()                       // true
	_ = list.Len()                           // 0
	list.Add("a")                            // ["a"]
	list.Clear()                             // []
	list.Insert(0, "b")                      // ["b"]
	list.Insert(0, "a")                      // ["a","b"]
}
```


### Set

A set is a data structure that can store elements and has no repeated values. It is a computer implementation of the mathematical concept of a finite set. Unlike most other collection types, rather than retrieving a specific element from a set, one typically tests an element for membership in a set. This structure is often used to ensure that no duplicates are present in a container.

```go
// Set a set interface
type Set = Collection
```

#### HashSet

A [set](#set) backed by a hash table (actually a Go's map). It makes no guarantees as to the iteration order of the set.

Implements [Set](#set) interfaces.

```go
package main

import (
	"github.com/askasoft/pango/cog/hashset"
)

func main() {
	set := hashset.NewHashSet[int]()
	set.Add(1)                // 1
	set.AddAll(2, 2, 3, 4, 5) // 3, 1, 2, 4, 5 (random order, duplicates ignored)
	set.Remove(4)             // 5, 3, 2, 1 (random order)
	set.RemoveAll(2, 3)       // 1, 5 (random order)
	set.Contains(1)           // true
	set.ContainsAll(1, 5)     // true
	set.ContainsAll(1, 6)     // false
	_ = set.Values()          // []int{5,1} (random order)
	set.Clear()               // empty
	set.IsEmpty()             // true
	set.Len()                 // 0
}
```

#### LinkedHashSet

A [set](#set) that preserves insertion-order. Data structure is backed by a hash table to store values and doubly-linked list to store insertion ordering.

Implements [Set](#set), [Iterator](#iterator) interfaces.

```go
package main

import (
	"github.com/askasoft/pango/cog/linkedhashset"
)

func main() {
	set := linkedhashset.NewLinkedHashSet[int]()
	set.Add(5)                // 5
	set.AddAll(4, 4, 3, 2, 1) // 5, 4, 3, 2, 1 (in insertion-order, duplicates ignored)
	set.Add(4)                // 5, 4, 3, 2, 1 (duplicates ignored, insertion-order unchanged)
	set.Remove(4)             // 5, 3, 2, 1 (in insertion-order)
	set.RemoveAll(2, 3)       // 5, 1 (in insertion-order)
	set.Contains(1)           // true
	set.ContainsAll(1, 5)     // true
	set.ContainsAll(1, 6)     // false
	_ = set.Values()          // []int{5, 1} (in insertion-order)
	set.Clear()               // empty
	set.IsEmpty()             // true
	set.Len()                 // 0
}
```

#### TreeSet

A [set](#set) backed by a red-black-tree to keep the elements ordered with respect to the comparator.

Implements [Set](#set), [Iterator](#iterator) interfaces.

```go
package main

import (
	"cmp"
	"github.com/askasoft/pango/cog/treeset"
)

func main() {
	set := treeset.NewTreeSet(cmp.Compare[int])
	set.Add(1)                // 1
	set.AddAll(2, 2, 3, 4, 5) // 1, 2, 3, 4, 5 (in order, duplicates ignored)
	set.Remove(4)             // 1, 2, 3, 5 (in order)
	set.RemoveAll(2, 3)       // 1, 5 (in order)
	set.Contains(1)           // true
	set.ContainsAll(1, 5)     // true
	set.ContainsAll(1, 6)     // false
	_ = set.Values()          // []int{1,5} (in order)
	set.Clear()               // empty
	set.IsEmpty()             // true
	set.Len()                 // 0
}
```

### Map

A Map is a data structure that maps keys to values. A map cannot contain duplicate keys and each key can map to at most one value.

Extends [Container](#container) interface.

```go
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

	Eachable2[K, V]
}
```

#### HashMap

A [map](#map) based on hash tables. Keys are unordered.

Implements [Map](#map) interfaces.

```go
package main

import (
	"github.com/askasoft/pango/cog/hashmap"
)

func main() {
	m := hashmap.NewHashMap[int, string]()
	m.Set(1, "x")   // 1->x
	m.Set(2, "b")   // 2->b, 1->x (random order)
	m.Set(1, "a")   // 2->b, 1->a (random order)
	_, _ = m.Get(2) // b, true
	_, _ = m.Get(3) // nil, false
	_ = m.Values()  // []interface {}{"b", "a"} (random order)
	_ = m.Keys()    // []interface {}{1, 2} (random order)
	m.Remove(1)     // 2->b
	m.Clear()       // empty
	m.IsEmpty()     // true
	m.Len()         // 0
}
```

#### LinkedHashMap

A [map](#map) that preserves insertion-order. It is backed by a hash table to store values and doubly-linked list to store ordering.

Implements [Map](#map), [Iterator](#iterator) interfaces.

```go
package main

import (
	"github.com/askasoft/pango/cog/linkedhashmap"
)

func main() {
	m := linkedhashmap.NewLinkedHashMap[int, string]()
	m.Set(2, "b")   // 2->b
	m.Set(1, "x")   // 2->b, 1->x (insertion-order)
	m.Set(1, "a")   // 2->b, 1->a (insertion-order)
	_, _ = m.Get(2) // b, true
	_, _ = m.Get(3) // nil, false
	_ = m.Values()  // []interface {}{"b", "a"} (insertion-order)
	_ = m.Keys()    // []interface {}{2, 1} (insertion-order)
	m.Remove(1)     // 2->b
	m.Clear()       // empty
	m.IsEmpty()     // true
	m.Len()         // 0
}
```

#### TreeMap

A [map](#map) based on red-black tree. Keys are ordered with respect to the comparator.

Implements [Map](#map), [Iterator2](#iterator2) interfaces.

```go
package main

import (
	"cmp"
	"github.com/askasoft/pango/cog/treemap"
)

func main() {
	m := treemap.NewTreeMap[int, string](cmp.Compare[int])
	m.Set(1, "x")   // 1->x
	m.Set(2, "b")   // 1->x, 2->b (in order)
	m.Set(1, "a")   // 1->a, 2->b (in order)
	_, _ = m.Get(2) // b, true
	_, _ = m.Get(3) // nil, false
	_ = m.Values()  // []interface {}{"a", "b"} (in order)
	_ = m.Keys()    // []interface {}{1, 2} (in order)
	m.Remove(1)     // 2->b
	m.Clear()       // empty
	m.IsEmpty()     // true
	m.Len()         // 0

	// Other:
	m.Head()  // Returns the minimum key and its value from map.
	m.Tail()  // Returns the maximum key and its value from map.
}
```

### Iterator

All ordered containers have stateful iterators. Typically an iterator is obtained by _Iterator()_ function of an ordered container. Once obtained, iterator's _Next()_ function moves the iterator to the next element and returns true if there was a next element. If there was an element, then element's can be obtained by iterator's _Value()_ function.

Note: it is unsafe to use Iterator.Remove() element or Iterator.SetValue() while iterating, 
but the iterator's Prev()/Next() may be different after SetValue().

Typical usage:
```go
// forward
for it := list.Iterator(); it.Next(); {
	value := it.Value()
	...
}

// backward
for it := list.Iterator(); it.Prev(); {
	value := it.Value()
	...
}
```

### Iterator2

All ordered maps have stateful iterators. Typically an iterator is obtained by _Iterator()_ function of an ordered map. Once obtained, iterator's _Next()_ function moves the iterator to the next element and returns true if there was a next element. If there was an element, then element's can be obtained by iterator's _Key()_, _Value()_ function.

Note: it is unsafe to use Iterator.Remove() element or Iterator.SetValue() while iterating.

Typical usage:
```go
// forward
for it := list.Iterator(); it.Next(); {
	key, value := it.Key(), it.Value()
	...
}

// backward
for it := list.Iterator(); it.Prev(); {
	key, value := it.Key(), it.Value()
	...
}
```

 Compare
-----------------------------------------------------------------------

Various helper functions used by sortable collections.

### Comparator

Some data structures (e.g. TreeMap, TreeSet) require a comparator function to automatically keep their elements sorted upon insertion. This comparator is necessary during the initalization.

Comparator is defined as:

```go
// Should return a int:
//   negative : if a < b
//   zero     : if a == b
//   positive : if a > b
```

Comparator signature:

```go
type Compare[T any] func(a, b T) int
```

Writing custom comparators is easy:

```go
package main

import (
	"fmt"
	"github.com/askasoft/pango/cog/treeset"
)

type User struct {
	id   int
	name string
}

// Custom comparator (sort by IDs)
func byID(a, b any) int {

	// Type assertion, program will panic if this is not respected
	c1 := a.(User)
	c2 := b.(User)

	switch {
	case c1.id > c2.id:
		return 1
	case c1.id < c2.id:
		return -1
	default:
		return 0
	}
}

func main() {
	set := treeset.NewTreeSet(byID)

	set.Add(User{2, "Second"})
	set.Add(User{3, "Third"})
	set.Add(User{1, "First"})
	set.Add(User{4, "Fourth"})

	fmt.Println(set) // {1 First}, {2 Second}, {3 Third}, {4 Fourth}
}
```

### Less

Some data structures require a less compare function to sort it's elements (e.g. ArrayList.Sort()).

Less comparator is defined as:

```go
// Should return a bool:
//    true : if a < b
//    false: if a >= b
```

Less signature:

```go
type Less[T any] func(a, b T) bool
```
