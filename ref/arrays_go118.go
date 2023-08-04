//go:build go1.18
// +build go1.18

package ref

// ArrayOf returns a []T{args[0], args[1], ...}
func ArrayOf[T any](args ...T) []T {
	return args
}
