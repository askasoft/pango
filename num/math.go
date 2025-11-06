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

// AbsInt returns the absolute value of x.
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// AbsInt16 returns the absolute value of x.
func AbsInt16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

// AbsInt32 returns the absolute value of x.
func AbsInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// AbsInt64 returns the absolute value of x.
func AbsInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// AbsFloat32 returns the absolute value of x.
func AbsFloat32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// AbsFloat64 returns the absolute value of x.
func AbsFloat64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
