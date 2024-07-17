//go:build go1.18
// +build go1.18

package num

type Signed interface {
	int | int16 | int32 | int64 | float32 | float64
}

type Number interface {
	byte | int | int16 | int32 | int64 | uint | uint16 | uint32 | uint64 | float32 | float64
}

// Abs returns the absolute value of x.
func Abs[T Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// IfZero returns (a == 0 ? b : a)
func IfZero[T Number](a, b T) T {
	if a == 0 {
		return b
	}
	return a
}
