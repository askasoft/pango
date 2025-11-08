package num

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Signed interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type Unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// IfZero returns (a == 0 ? b : a)
func IfZero[T Number](a, b T) T {
	if a == 0 {
		return b
	}
	return a
}

// Abs returns the absolute value of x.
func Abs[T Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
