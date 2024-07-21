 Compare
-----------------------------------------------------------------------

Various helper functions used by [Collection](../cog/) package.

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

All common comparators for builtin types are included in the package:

```go
func CompareString(a, b string) int
func CompareInt(a, b int) int
func CompareInt8(a, b int8) int
func CompareInt16(a, b int16) int
func CompareInt32(a, b int32) int
func CompareInt64(a, b int64) int
func CompareUInt(a, b uint) int
func CompareUInt8(a, b uint8) int
func CompareUInt16(a, b uint16) int
func CompareUInt32(a, b uint32) int
func CompareUInt64(a, b uint64) int
func CompareFloat32(a, b float32) int
func CompareFloat64(a, b float64) int
func CompareByte(a, b byte) int
func CompareRune(a, b rune) int
```

### Less

Some data structures require a less compare function to sort it's elements (e.g. ArrayList.Sort()).

Less comparator is defined as:

```go
// Should return a bool:
//    true : if a < b
//    false: if a >= b
```

Comparator signature:

```go
type Less[T any] func(a, b T) bool
```

All common comparators for builtin types are included in the package:

```go
func LessString(a, b string) bool
func LessByte(a, b byte) bool
func LessRune(a, b rune) bool
func LessInt(a, b int) bool
func LessInt8(a, b int8) bool
func LessInt16(a, b int16) bool
func LessInt32(a, b int32) bool
func LessInt64(a, b int64) bool
func LessUint(a, b uint) bool
func LessUint8(a, b uint8) bool
func LessUint16(a, b uint16) bool
func LessUint32(a, b uint32) bool
func LessUint64(a, b uint64) bool
func LessFloat32(a, b float32) bool
func LessFloat64(a, b float64) bool
```


