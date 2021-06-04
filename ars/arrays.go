package ars

import "bytes"

// IndexByte returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexByte(a []byte, c byte) int {
	return bytes.IndexByte(a, c)
}

// IndexInt returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexInt(a []int, c int) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexInt32 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexInt32(a []int32, c int32) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexInt64 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexInt64(a []int64, c int64) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexUint returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexUint(a []uint, c uint) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexUint32 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexUint32(a []uint32, c uint32) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexUint64 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexUint64(a []uint64, c uint64) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexRune returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexRune(a []rune, c rune) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexFloat32 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexFloat32(a []float32, c float32) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexFloat64 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexFloat64(a []float64, c float64) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}
