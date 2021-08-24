 Pango
=====================================================================

![](https://github.com/pandafw/pango/raw/master/logo.png) [![Build Status](https://travis-ci.com/pandafw/pango.svg?branch=master)](https://travis-ci.com/pandafw/pango) [![codecov](https://codecov.io/gh/pandafw/pango/branch/master/graph/badge.svg)](https://codecov.io/gh/pandafw/pango) [![Apache 2](https://img.shields.io/badge/license-Apache%202-green)](https://www.apache.org/licenses/LICENSE-2.0.html) ![](https://github.com/pandafw/pango/raw/master/logo.png)



Pango is a GO development utility library.

| **Package**                     | **Description**                         |
| :------------------------------- | :-------------------------------------- |
| [cmp](#compare)                  | a Compare/Less helper functions package |
| [col](#col-collectioncontainer)  | a Collection/Container package          |
| [ini](#int)                      | a INI file read/write package           |
| iox                              | a IO utility package                    |
| [iox/fswatch](#ioxfswatch)       | recursive directory notifications built as a wrapper around fsnotify (golang) |
| [log](#log)                      | a Log like Log4j                        |
| [net/email](#netemail)           | a email sender package                  |
| [net/httpx](#nethttpx)           | a package for http                      |
| [net/slack](#netslack)           | a slack webhook sender package          |
| [str](#str)                      | a strings package wrapper and provides some useful utility functions      |
| [tpl](#tpl)                      | a text/template, html/template wrapper for recursively load template file |
| [x/ginx/gindump](#xginxgindump)  | a http request/response dumper middleware for gin                         |
| [x/ginx/ginfile](#xginxginfile)  | a static file handler with Cache-Control header support for gin           |
| [x/ginx/gingzip](#xginxginzip)   | a gzip encoding support middleware for gin                                |
| [x/ginx/ginlog](#xginxginlog)    | a access logger middleware for gin                                        |


### How to install?

	go get github.com/pandafw/pango



 col (Collection/Container)
-----------------------------------------------------------------------
All data structures implement the container interface with the following methods:

### Container

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
| [SortedList](#sortedlist)       |      Y      |      Y       |              |
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
```

### List

A list is a data structure that stores values and may have repeated values.

Extends [Collection](#collection) interface.

```go
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
```

#### ArrayList

A [list](#list) backed by a dynamic array that grows implicitly.

Implements [List](#list), [Iterator](#iterator) interfaces.

```go
package main

import (
	"github.com/pandafw/pango/col"
	"github.com/pandafw/pango/cmp"
)

func main() {
	list := col.NewArrayList()
	list.Add("a")                         // ["a"]
	list.Add("c", "b")                    // ["a","c","b"]
	list.Sort(cmp.LessString)             // ["a","b","c"]
	_ = list.Get(0)                       // "a"
	_ = list.Get(100)                     // panic
	_ = list.Contains("a", "b", "c")      // true
	_ = list.Contains("a", "b", "c", "d") // false
	list.Swap(0, 1)                       // ["b","a",c"]
	list.Remove(2)                        // ["b","a"]
	list.Remove(1)                        // ["b"]
	list.Remove(0)                        // []
	list.Remove(0)                        // [] (ignored)
	_ = list.IsEmpty()                    // true
	_ = list.Len()                        // 0
	list.Add("a")                         // ["a"]
	list.Clear()                          // []
	list.Insert(0, "b")                   // ["b"]
	list.Insert(0, "a")                   // ["a","b"]
}
```

#### LinkedList

A [list](#list) where each element points to the next and previous elements in the list.

Implements [List](#list), [Iterator](#iterator) interfaces.

```go
package main

import (
	"github.com/pandafw/pango/col"
	"github.com/pandafw/pango/cmp"
)

func main() {
	list := col.NewLinkedList()
	list.Add("a")                         // ["a"]
	list.Add("c", "b")                    // ["a","c","b"]
	list.Sort(cmp.LessString)             // ["a","b","c"]
	_ = list.Get(0)                       // "a"
	_ = list.Get(100)                     // panic
	_ = list.Contains("a", "b", "c")      // true
	_ = list.Contains("a", "b", "c", "d") // false
	list.Swap(0, 1)                       // ["b","a",c"]
	list.Remove(2)                        // ["b","a"]
	list.Remove(1)                        // ["b"]
	list.Remove(0)                        // []
	list.Remove(0)                        // [] (ignored)
	_ = list.IsEmpty()                    // true
	_ = list.Len()                        // 0
	list.Add("a")                         // ["a"]
	list.Clear()                          // []
	list.Insert(0, "b")                   // ["b"]
	list.Insert(0, "a")                   // ["a","b"]
}
```

#### SortedList

A [list](#list) where each element points to the next and previous elements in the list and automatically keep the elements ordered with respect to the comparator.

Implements [List](#list), [Iterator](#iterator) interfaces.

```go
package main

import (
	"github.com/pandafw/pango/col"
	"github.com/pandafw/pango/cmp"
)

func main() {
	list := NewSortedList(cmp.LessInt) // empty (keys are of type int)
	list.Add(1)                        // 1
	list.Add(2, 2, 3, 4, 5)            // 1, 2, 2, 3, 4, 5 (in order)
	_ = list.Get(0)                    // 1
	_ = list.Get(100)                  // panic
	list.Delete(4)                     // 1, 2, 2, 3, 5 (in order)
	list.Delete(2, 3)                  // 1, 5 (in order)
	list.Contains(1)                   // true
	list.Contains(1, 5)                // true
	list.Contains(1, 6)                // false
	_ = list.Values()                  // []int{1,5} (in order)
	list.Clear()                       // empty
	list.IsEmpty()                     // true
	list.Len()                         // 0
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
	"github.com/pandafw/pango/col"
)

func main() {
	set := col.NewHashSet()    // empty
	set.Add(1)             // 1
	set.Add(2, 2, 3, 4, 5) // 3, 1, 2, 4, 5 (random order, duplicates ignored)
	set.Delete(4)          // 5, 3, 2, 1 (random order)
	set.Delete(2, 3)       // 1, 5 (random order)
	set.Contains(1)        // true
	set.Contains(1, 5)     // true
	set.Contains(1, 6)     // false
	_ = set.Values()       // []int{5,1} (random order)
	set.Clear()            // empty
	set.IsEmpty()          // true
	set.Len()              // 0
}
```

#### LinkedHashSet

A [set](#set) that preserves insertion-order. Data structure is backed by a hash table to store values and doubly-linked list to store insertion ordering.

Implements [Set](#set), [Iterator](#iterator) interfaces.

```go
package main

import (
	"github.com/pandafw/pango/col"
	"github.com/pandafw/pango/cmp"
)

func main() {
	set := col.NewLinkedHashSet() // empty
	set.Add(5)                 // 5
	set.Add(4, 4, 3, 2, 1)     // 5, 4, 3, 2, 1 (in insertion-order, duplicates ignored)
	set.Add(4)                 // 5, 4, 3, 2, 1 (duplicates ignored, insertion-order unchanged)
	set.Delete(4)              // 5, 3, 2, 1 (in insertion-order)
	set.Delete(2, 3)           // 5, 1 (in insertion-order)
	set.Contains(1)            // true
	set.Contains(1, 5)         // true
	set.Contains(1, 6)         // false
	_ = set.Values()           // []int{5, 1} (in insertion-order)
	set.Clear()                // empty
	set.IsEmpty()              // true
	set.Len()                  // 0
}
```

#### TreeSet

A [set](#set) backed by a red-black-tree to keep the elements ordered with respect to the comparator.

Implements [Set](#set), [Iterator](#iterator) interfaces.

```go
package main

import (
	"github.com/pandafw/pango/col"
	"github.com/pandafw/pango/cmp"
)

func main() {
	set := col.NewTreeSet(cmp.CompareInt) // empty (keys are of type int)
	set.Add(1)                            // 1
	set.Add(2, 2, 3, 4, 5)                // 1, 2, 3, 4, 5 (in order, duplicates ignored)
	set.Delete(4)                         // 1, 2, 3, 5 (in order)
	set.Delete(2, 3)                      // 1, 5 (in order)
	set.Contains(1)                       // true
	set.Contains(1, 5)                    // true
	set.Contains(1, 6)                    // false
	_ = set.Values()                      // []int{1,5} (in order)
	set.Clear()                           // empty
	set.IsEmpty()                         // true
	set.Len()                             // 0
}
```

### Map

A Map is a data structure that maps keys to values. A map cannot contain duplicate keys and each key can map to at most one value.

Extends [Container](#container) interface.

```go
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
```

#### HashMap

A [map](#map) based on hash tables. Keys are unordered.

Implements [Map](#map) interfaces.

```go
package main

import (
	"github.com/pandafw/pango/col"
)

func main() {
	m := col.NewHashMap() // empty
	m.Set(1, "x")      // 1->x
	m.Set(2, "b")      // 2->b, 1->x (random order)
	m.Set(1, "a")      // 2->b, 1->a (random order)
	_, _ = m.Get(2)    // b, true
	_, _ = m.Get(3)    // nil, false
	_ = m.Values()     // []interface {}{"b", "a"} (random order)
	_ = m.Keys()       // []interface {}{1, 2} (random order)
	m.Delete(1)        // 2->b
	m.Clear()          // empty
	m.IsEmpty()        // true
	m.Len()            // 0
}
```

#### LinkedHashMap

A [map](#map) that preserves insertion-order. It is backed by a hash table to store values and doubly-linked list to store ordering.

Implements [Map](#map), [Iterator](#iterator) interfaces.

```go
package main

import (
	"github.com/pandafw/pango/col"
)

func main() {
	m := col.NewLinkedHashMap() // empty (keys are of type int)
	m.Set(2, "b")            // 2->b
	m.Set(1, "x")            // 2->b, 1->x (insertion-order)
	m.Set(1, "a")            // 2->b, 1->a (insertion-order)
	_, _ = m.Get(2)          // b, true
	_, _ = m.Get(3)          // nil, false
	_ = m.Values()           // []interface {}{"b", "a"} (insertion-order)
	_ = m.Keys()             // []interface {}{2, 1} (insertion-order)
	m.Delete(1)              // 2->b
	m.Clear()                // empty
	m.IsEmpty()              // true
	m.Len()                  // 0
}
```

#### TreeMap

A [map](#map) based on red-black tree. Keys are ordered with respect to the comparator.

Implements [Map](#map), [Iterator2](#iterator2) interfaces.

```go
package main

import (
	"github.com/pandafw/pango/col"
	"github.com/pandafw/pango/cmp"
)

func main() {
	m := col.NewTreeMap(cmp.CompareInt) // empty (keys are of type int)
	m.Set(1, "x")                   // 1->x
	m.Set(2, "b")                   // 1->x, 2->b (in order)
	m.Set(1, "a")                   // 1->a, 2->b (in order)
	_, _ = m.Get(2)                 // b, true
	_, _ = m.Get(3)                 // nil, false
	_ = m.Values()                  // []interface {}{"a", "b"} (in order)
	_ = m.Keys()                    // []interface {}{1, 2} (in order)
	m.Delete(1)                     // 2->b
	m.Clear()                       // empty
	m.IsEmpty()                     // true
	m.Len()                         // 0

	// Other:
	m.Front() // Returns the minimum key and its value from map.
	m.Back()  // Returns the maximum key and its value from map.
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

Various helper functions used by [Collection](colcontainercollection) package.

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
type Compare func(a, b interface{}) int
```

All common comparators for builtin types are included in the package:

```go
func CompareString(a, b interface{}) int
func CompareInt(a, b interface{}) int
func CompareInt8(a, b interface{}) int
func CompareInt16(a, b interface{}) int
func CompareInt32(a, b interface{}) int
func CompareInt64(a, b interface{}) int
func CompareUInt(a, b interface{}) int
func CompareUInt8(a, b interface{}) int
func CompareUInt16(a, b interface{}) int
func CompareUInt32(a, b interface{}) int
func CompareUInt64(a, b interface{}) int
func CompareFloat32(a, b interface{}) int
func CompareFloat64(a, b interface{}) int
func CompareByte(a, b interface{}) int
func CompareRune(a, b interface{}) int
```

Writing custom comparators is easy:

```go
package main

import (
	"fmt"
	"github.com/pandafw/pango/col"
)

type User struct {
	id   int
	name string
}

// Custom comparator (sort by IDs)
func byID(a, b interface{}) int {

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
	set := col.NewTreeSet(byID)

	set.Add(User{2, "Second"})
	set.Add(User{3, "Third"})
	set.Add(User{1, "First"})
	set.Add(User{4, "Fourth"})

	fmt.Println(set) // {1 First}, {2 Second}, {3 Third}, {4 Fourth}
}
```

### Less

Some data structures (e.g. SortedList) require a Less compare function to automatically keep their elements sorted upon insertion. 
This comparator is necessary during the initalization. 
Some data structures require a less compare function to sort it's elements (e.g. ArrayList.Sort()).

Less comparator is defined as:

```go
// Should return a bool:
//    true : if a < b
//    false: if a >= b
```

Comparator signature:

```go
type Less func(a, b interface{}) bool
```

All common comparators for builtin types are included in the package:

```go
func LessString(a, b interface{}) bool
func LessByte(a, b interface{}) bool
func LessRune(a, b interface{}) bool
func LessInt(a, b interface{}) bool
func LessInt8(a, b interface{}) bool
func LessInt16(a, b interface{}) bool
func LessInt32(a, b interface{}) bool
func LessInt64(a, b interface{}) bool
func LessUint(a, b interface{}) bool
func LessUint8(a, b interface{}) bool
func LessUint16(a, b interface{}) bool
func LessUint32(a, b interface{}) bool
func LessUint64(a, b interface{}) bool
func LessFloat32(a, b interface{}) bool
func LessFloat64(a, b interface{}) bool
```


 iox/fswatch
-----------------------------------------------------------------------
recursive directory notifications built as a wrapper around fsnotify (golang)

This is a wrapper around https://github.com/fsnotify/fsnotify instead of only monitoring a top level folder,
it allows you to monitor all folders underneath the folder you specify.

### Example:

(error handling omitted to improve readability)

```golang
	import "github.com/panafw/pango/iox/fswatch"

	// works exactly like fsnotify and implements the same API.
	watcher, err := fswatch.NewFileWatcher()

	// watch recursive and recieve events with callback function
	watcher.AddRecursive("watchdir", fswatch.OpALL, "", func(path string, op fswatch.Op) {
		fmt.Printf("%s %s\n", path, op)
	})

```

 log
-----------------------------------------------------------------------
log is a Go log manager. It can use many log writers. This package is inspired by https://github.com/pandafw/panda/tree/master/panda-core/src/main/java/panda/log .


### What writers are supported?

As of now this log support stream(console), file, slack, smtp, connection(tcp), webhook.


### How to use it?

First you must import it

```golang
import (
	"github.com/pandafw/pango/log"
)
```

Then init a Log (example with console writer)

```golang
	log := log.NewLog()
	log.SetWriter(&log.StreamWriter{Color:true})
```

Use it like this:

```golang
	log.Trace("trace")
	log.Debug("debug")
	log.Info("info")
	log.Warn("warning")
	log.Fatal("fatal")
```

#### File writer

Configure file writer like this:

```golang
	log := log.NewLog()
	log.SetWriter(&log.FileWriter{Path:"test.log"})
```

#### Conn writer

Configure like this:

```golang
	log := log.NewLog()
	log.SetWriter(&log.ConnWriter{Net:"tcp",Addr:":7020"})
	log.Info("info")
```

#### Slack writer

Configure like this:

```golang
	log := log.NewLog()
	log.SetWriter(&log.SlackWriter{
		Webhook: "https://hooks.slack.com/services/...",
		Channel: "alert",
		Username: "gotest",
	})
	log.Error("error")
```

#### SMTP writer

Configure like this:

```golang
	log := log.NewLog()
	log.SetWriter(&log.SMTPWriter{
		Host: "smtp.gmail.com",
		Port: 587,
		Username: "pangotest@gmail.com",
		Password: "xxxxxxxx",
		From: "xxxx@gmail.com",
		Tos: []string{"someone@gmail.com"},
	})
	log.Fatal("oh my god!")
```

#### Webhook writer

Configure like this:

```golang
	log := log.NewLog()
	log.SetWriter(&log.WebhookWriter{
		Webhook: "http://localhost:9200/pango/logs",
		ContentType: "application/json",
		Timeout: time.Second*5,
	})
	log.Fatal("fatal error!")
```

#### Multiple writer

Configure like this:

```golang
	log := log.NewLog()
	fw := &log.FileWriter{Path:"test.log"}
	ww := &log.WebhookWriter{
		Webhook: "http://localhost:9200/pango/logs",
		ContentType: "application/json",
		Timeout: time.Second*5,
	}
	log.SetWriter(log.NewMultiWriter(fw, ww))
	log.Fatal("fatal error!")
```


#### Async writer

Configure like this:

```golang
	log := log.NewLog()
	ww := &log.WebhookWriter{
		Webhook: "http://localhost:9200/pango/logs",
		ContentType: "application/json",
		Timeout: time.Second*5,
	}
	log.SetWriter(log.NewAsyncWriter(ww, 1000))
	log.Fatal("fatal error!")
```


### Configure from ini file

```golang
	log := log.NewLog()
	log.Config("log.ini")
```

#### log.ini

```ini
# log configuration #

### log async ###
async = 1000

### log format ###
#format=json:{"level":%l, "file":%S, "func":%F, "msg": %m}%n
format=text:%l %S %F() - %m%n%T

### log writer ###
writer = stdout, stderr, tcp, dailyfile, slack, smtp, webhook

### log level ###
[level]
* = info
sql = debug
http = trace

### stdout writer ###
[writer.stdout]
format = %l - %m%n%T
filter = name:out level:debug

### tcp writer ###
[writer.tcp]
addr = localhost:9999
timeout = 5s
format = %l - %m%n%T
filter = level:error

### file writer ###
[writer.dailyfile]
_ = file
path = /tmp/gotest/logs/test.log
dirPerm = 0777
maxDays = 7
format = %l %S:%L %F() - %m%n%T
filter = level:error

### slack writer ###
[writer.slack]
subject = %l - %m 
channel = develop
username = gotest
webhook = https://hooks.slack.com/services/...
timeout = 5s
format = %l - %m%n%T
filter = level:error

### smtp writer ###
[writer.smtp]
host = localhost
port = 25
username = -----
password = xxxxxxx
from = pango@google.com
to = to1@test.com, to2@test.com
cc = cc1@test.com, cc2@test.com
timeout = 5s
subject = %l - %m 
format = %l - %m%n%T
filter = level:error

### webhook writer ###
[writer.webhook]
webhook = http://localhost:9200/pango/logs
contentType = application/json
timeout = 5s
format = json:{"when":%t{2006-01-02T15:04:05.000Z07:00}, "level":%l, "file":%S, "line":%L, "func":%F, "msg": %m, "stack": %T}%n
filter = level:error
```
