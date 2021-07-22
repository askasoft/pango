package col

// Collection is base interface that all data structures implement.
type Collection interface {
	IsEmpty() bool
	Len() int
	Clear()
	Values() []interface{}
}
