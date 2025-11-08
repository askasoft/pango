package cal

import (
	"errors"
	"fmt"
	"testing"
)

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func add[T number](a, b T) T {
	return a + b
}

func sub[T number](a, b T) T {
	return a - b
}

func mul[T number](a, b T) T {
	return a * b
}

func div[T number](a, b T) T {
	return a / b
}

var addTests = []struct {
	a, b, w any
	e       error
}{
	{int(1), int(2), int(1) + int(2), nil},
	{int(1), int8(2), int(1) + int(2), nil},
	{int(1), int16(2), int(1) + int(2), nil},
	{int(1), int32(2), int(1) + int(2), nil},
	{int(1), int64(2), int(1) + int(2), nil},
	{int(1), uint(2), int(1) + int(2), nil},
	{int(1), uint8(1<<8 - 1), int(1) + int(1<<8-1), nil},
	{int(1), uint16(1<<16 - 1), int(1) + int(1<<16-1), nil},
	{int(1), uint32(1<<32 - 1), int(1) + int(1<<32-1), nil},
	{int(1), uint64(1<<64 - 1), int(1) + int(-1), nil},
	{int(1), float32(0.1), float32(1) + float32(0.1), nil},
	{int(1), float64(0.1), float64(1) + float64(0.1), nil},
	{int(1), "2", "12", nil},

	{int8(1), int(2), int(1) + int(2), nil},
	{int8(1), int8(2), int8(1) + int8(2), nil},
	{int8(1), int16(2), int16(1) + int16(2), nil},
	{int8(1), int32(2), int32(1) + int32(2), nil},
	{int8(1), int64(2), int64(1) + int64(2), nil},
	{int8(1), uint(2), uint(1) + uint(2), nil},
	{int8(1), uint8(1<<8 - 1), int8(1) + int8(-1), nil},
	{int8(1), uint16(1<<16 - 1), add(uint16(1), uint16(1<<16-1)), nil},
	{int8(1), uint32(1<<32 - 1), add(uint32(1), uint32(1<<32-1)), nil},
	{int8(1), uint64(1<<64 - 1), add(uint64(1), uint64(1<<64-1)), nil},
	{int8(1), float32(0.1), float32(1) + float32(0.1), nil},
	{int8(1), float64(0.1), float64(1) + float64(0.1), nil},
	{int8(1), "2", "12", nil},

	{int16(1), int(2), int(1) + int(2), nil},
	{int16(1), int8(2), int16(1) + int16(2), nil},
	{int16(1), int16(2), int16(1) + int16(2), nil},
	{int16(1), int32(2), int32(1) + int32(2), nil},
	{int16(1), int64(2), int64(1) + int64(2), nil},
	{int16(1), uint(2), uint(1) + uint(2), nil},
	{int16(1), uint8(1<<8 - 1), int16(1) + int16(1<<8-1), nil},
	{int16(1), uint16(1<<16 - 1), add(int16(1), int16(-1)), nil},
	{int16(1), uint32(1<<32 - 1), add(uint32(1), uint32(1<<32-1)), nil},
	{int16(1), uint64(1<<64 - 1), add(uint64(1), uint64(1<<64-1)), nil},
	{int16(1), float32(0.1), float32(1) + float32(0.1), nil},
	{int16(1), float64(0.1), float64(1) + float64(0.1), nil},
	{int16(1), "2", "12", nil},

	{int32(1), int(2), int(1) + int(2), nil},
	{int32(1), int8(2), int32(1) + int32(2), nil},
	{int32(1), int16(2), int32(1) + int32(2), nil},
	{int32(1), int32(2), int32(1) + int32(2), nil},
	{int32(1), int64(2), int64(1) + int64(2), nil},
	{int32(1), uint(2), uint(1) + uint(2), nil},
	{int32(1), uint8(1<<8 - 1), int32(1) + int32(1<<8-1), nil},
	{int32(1), uint16(1<<16 - 1), add(int32(1), int32(1<<16-1)), nil},
	{int32(1), uint32(1<<32 - 1), add(int32(1), int32(-1)), nil},
	{int32(1), uint64(1<<64 - 1), add(uint64(1), uint64(1<<64-1)), nil},
	{int32(1), float32(0.1), float32(1) + float32(0.1), nil},
	{int32(1), float64(0.1), float64(1) + float64(0.1), nil},
	{int32(1), "2", "12", nil},

	{int64(1), int(2), int64(1) + int64(2), nil},
	{int64(1), int8(2), int64(1) + int64(2), nil},
	{int64(1), int16(2), int64(1) + int64(2), nil},
	{int64(1), int32(2), int64(1) + int64(2), nil},
	{int64(1), int64(2), int64(1) + int64(2), nil},
	{int64(1), uint(2), int64(1) + int64(2), nil},
	{int64(1), uint8(1<<8 - 1), int64(1) + int64(1<<8-1), nil},
	{int64(1), uint16(1<<16 - 1), add(int64(1), int64(1<<16-1)), nil},
	{int64(1), uint32(1<<32 - 1), add(int64(1), int64(1<<32-1)), nil},
	{int64(1), uint64(1<<64 - 1), add(int64(1), int64(-1)), nil},
	{int64(1), float32(0.1), float32(1) + float32(0.1), nil},
	{int64(1), float64(0.1), float64(1) + float64(0.1), nil},
	{int64(1), "2", "12", nil},

	{uint(1), int(2), int(1) + int(2), nil},
	{uint(1), int8(2), uint(1) + uint(2), nil},
	{uint(1), int16(2), uint(1) + uint(2), nil},
	{uint(1), int32(2), uint(1) + uint(2), nil},
	{uint(1), int64(2), int64(1) + int64(2), nil},
	{uint(1), uint(2), uint(1) + uint(2), nil},
	{uint(1), uint8(1<<8 - 1), uint(1) + uint(1<<8-1), nil},
	{uint(1), uint16(1<<16 - 1), uint(1) + uint(1<<16-1), nil},
	{uint(1), uint32(1<<32 - 1), uint(1) + uint(1<<32-1), nil},
	{uint(1), uint64(1<<64 - 1), add(uint(1), uint(1<<64-1)), nil},
	{uint(1), float32(0.1), float32(1) + float32(0.1), nil},
	{uint(1), float64(0.1), float64(1) + float64(0.1), nil},
	{uint(1), "2", "12", nil},

	{uint8(1), int(2), int(1) + int(2), nil},
	{uint8(1), int8(2), int8(1) + int8(2), nil},
	{uint8(1), int16(2), int16(1) + int16(2), nil},
	{uint8(1), int32(2), int32(1) + int32(2), nil},
	{uint8(1), int64(2), int64(1) + int64(2), nil},
	{uint8(1), uint(2), uint(1) + uint(2), nil},
	{uint8(1), uint8(1<<8 - 1), add(uint8(1), uint8(1<<8-1)), nil},
	{uint8(1), uint16(1<<16 - 1), add(uint16(1), uint16(1<<16-1)), nil},
	{uint8(1), uint32(1<<32 - 1), add(uint32(1), uint32(1<<32-1)), nil},
	{uint8(1), uint64(1<<64 - 1), add(uint64(1), uint64(1<<64-1)), nil},
	{uint8(1), float32(0.1), float32(1) + float32(0.1), nil},
	{uint8(1), float64(0.1), float64(1) + float64(0.1), nil},
	{uint8(1), "2", "12", nil},

	{uint16(1), int(2), int(1) + int(2), nil},
	{uint16(1), int8(2), uint16(1) + uint16(2), nil},
	{uint16(1), int16(2), int16(1) + int16(2), nil},
	{uint16(1), int32(2), int32(1) + int32(2), nil},
	{uint16(1), int64(2), int64(1) + int64(2), nil},
	{uint16(1), uint(2), uint(1) + uint(2), nil},
	{uint16(1), uint8(1<<8 - 1), uint16(1) + uint16(1<<8-1), nil},
	{uint16(1), uint16(1<<16 - 1), add(uint16(1), uint16(1<<16-1)), nil},
	{uint16(1), uint32(1<<32 - 1), add(uint32(1), uint32(1<<32-1)), nil},
	{uint16(1), uint64(1<<64 - 1), add(uint64(1), uint64(1<<64-1)), nil},
	{uint16(1), float32(0.1), float32(1) + float32(0.1), nil},
	{uint16(1), float64(0.1), float64(1) + float64(0.1), nil},
	{uint16(1), "2", "12", nil},

	{uint32(1), int(2), int(1) + int(2), nil},
	{uint32(1), int8(2), uint32(1) + uint32(2), nil},
	{uint32(1), int16(2), uint32(1) + uint32(2), nil},
	{uint32(1), int32(2), int32(1) + int32(2), nil},
	{uint32(1), int64(2), int64(1) + int64(2), nil},
	{uint32(1), uint(2), uint(1) + uint(2), nil},
	{uint32(1), uint8(1<<8 - 1), uint32(1) + uint32(1<<8-1), nil},
	{uint32(1), uint16(1<<16 - 1), add(uint32(1), uint32(1<<16-1)), nil},
	{uint32(1), uint32(1<<32 - 1), add(uint32(1), uint32(1<<32-1)), nil},
	{uint32(1), uint64(1<<64 - 1), add(uint64(1), uint64(1<<64-1)), nil},
	{uint32(1), float32(0.1), float32(1) + float32(0.1), nil},
	{uint32(1), float64(0.1), float64(1) + float64(0.1), nil},
	{uint32(1), "2", "12", nil},

	{uint64(1), int(2), int(1) + int(2), nil},
	{uint64(1), int8(2), uint64(1) + uint64(2), nil},
	{uint64(1), int16(2), uint64(1) + uint64(2), nil},
	{uint64(1), int32(2), uint64(1) + uint64(2), nil},
	{uint64(1), int64(2), int64(1) + int64(2), nil},
	{uint64(1), uint(2), uint64(1) + uint64(2), nil},
	{uint64(1), uint8(1<<8 - 1), uint64(1) + uint64(1<<8-1), nil},
	{uint64(1), uint16(1<<16 - 1), add(uint64(1), uint64(1<<16-1)), nil},
	{uint64(1), uint32(1<<32 - 1), add(uint64(1), uint64(1<<32-1)), nil},
	{uint64(1), uint64(1<<64 - 1), add(uint64(1), uint64(1<<64-1)), nil},
	{uint64(1), float32(0.1), float32(1) + float32(0.1), nil},
	{uint64(1), float64(0.1), float64(1) + float64(0.1), nil},
	{uint64(1), "2", "12", nil},

	{float32(1), int(2), float32(1) + float32(2), nil},
	{float32(1), int8(2), float32(1) + float32(2), nil},
	{float32(1), int16(2), float32(1) + float32(2), nil},
	{float32(1), int32(2), float32(1) + float32(2), nil},
	{float32(1), int64(2), float32(1) + float32(2), nil},
	{float32(1), uint(2), float32(1) + float32(2), nil},
	{float32(1), uint8(1<<8 - 1), float32(1) + float32(1<<8-1), nil},
	{float32(1), uint16(1<<16 - 1), add(float32(1), float32(1<<16-1)), nil},
	{float32(1), uint32(1<<32 - 1), add(float32(1), float32(1<<32-1)), nil},
	{float32(1), uint64(1<<64 - 1), add(float32(1), float32(1<<64-1)), nil},
	{float32(1), float32(0.1), float32(1) + float32(0.1), nil},
	{float32(1), float64(0.1), float64(1) + float64(0.1), nil},
	{float32(1), "2", "12", nil},

	{float64(1), int(2), float64(1) + float64(2), nil},
	{float64(1), int8(2), float64(1) + float64(2), nil},
	{float64(1), int16(2), float64(1) + float64(2), nil},
	{float64(1), int32(2), float64(1) + float64(2), nil},
	{float64(1), int64(2), float64(1) + float64(2), nil},
	{float64(1), uint(2), float64(1) + float64(2), nil},
	{float64(1), uint8(1<<8 - 1), float64(1) + float64(1<<8-1), nil},
	{float64(1), uint16(1<<16 - 1), add(float64(1), float64(1<<16-1)), nil},
	{float64(1), uint32(1<<32 - 1), add(float64(1), float64(1<<32-1)), nil},
	{float64(1), uint64(1<<64 - 1), add(float64(1), float64(1<<64-1)), nil},
	{float64(1), float32(0.1), float64(1) + float64(float32(0.1)), nil},
	{float64(1), float64(0.1), float64(1) + float64(0.1), nil},
	{float64(1), "2", "12", nil},

	{"1", int(2), "12", nil},
	{"1", int8(2), "12", nil},
	{"1", int16(2), "12", nil},
	{"1", int32(2), "12", nil},
	{"1", int64(2), "12", nil},
	{"1", uint(2), "12", nil},
	{"1", uint8(2), "12", nil},
	{"1", uint16(2), "12", nil},
	{"1", uint32(2), "12", nil},
	{"1", uint64(2), "12", nil},
	{"1", float32(-1), "1-1", nil},
	{"1", float64(-1), "1-1", nil},
	{"1", "2", "12", nil},
}

func TestAdd(t *testing.T) {
	for i, c := range addTests {
		r, e := Adds(c.a, c.b)

		if c.e != nil {
			if fmt.Sprint(c.e) != fmt.Sprint(e) {
				t.Errorf("[%d] Add(%T(%v), %T(%v)) = (%v), want: (%v)", i, c.a, c.a, c.b, c.b, e, c.e)
			}
			continue
		}

		if c.w != r {
			t.Errorf("[%d] Add(%T(%v), %T(%v)) = %T(%v), want: %T(%v)", i, c.a, c.a, c.b, c.b, r, r, c.w, c.w)
		}
	}
}

var subTests = []struct {
	a, b, w any
	e       error
}{
	{int(1), int(2), int(1) - int(2), nil},
	{int(1), int8(2), int(1) - int(2), nil},
	{int(1), int16(2), int(1) - int(2), nil},
	{int(1), int32(2), int(1) - int(2), nil},
	{int(1), int64(2), int(1) - int(2), nil},
	{int(1), uint(2), int(1) - int(2), nil},
	{int(1), uint8(1<<8 - 1), int(1) - int(1<<8-1), nil},
	{int(1), uint16(1<<16 - 1), int(1) - int(1<<16-1), nil},
	{int(1), uint32(1<<32 - 1), int(1) - int(1<<32-1), nil},
	{int(1), uint64(1<<64 - 1), int(1) - int(-1), nil},
	{int(1), float32(0.1), float32(1) - float32(0.1), nil},
	{int(1), float64(0.1), float64(1) - float64(0.1), nil},
	{int(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{int8(1), int(2), int(1) - int(2), nil},
	{int8(1), int8(2), int8(1) - int8(2), nil},
	{int8(1), int16(2), int16(1) - int16(2), nil},
	{int8(1), int32(2), int32(1) - int32(2), nil},
	{int8(1), int64(2), int64(1) - int64(2), nil},
	{int8(1), uint(2), sub(uint(1), uint(2)), nil},
	{int8(1), uint8(1<<8 - 1), int8(1) - int8(-1), nil},
	{int8(1), uint16(1<<16 - 1), sub(uint16(1), uint16(1<<16-1)), nil},
	{int8(1), uint32(1<<32 - 1), sub(uint32(1), uint32(1<<32-1)), nil},
	{int8(1), uint64(1<<64 - 1), sub(uint64(1), uint64(1<<64-1)), nil},
	{int8(1), float32(0.1), float32(1) - float32(0.1), nil},
	{int8(1), float64(0.1), float64(1) - float64(0.1), nil},
	{int8(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{int16(1), int(2), int(1) - int(2), nil},
	{int16(1), int8(2), int16(1) - int16(2), nil},
	{int16(1), int16(2), int16(1) - int16(2), nil},
	{int16(1), int32(2), int32(1) - int32(2), nil},
	{int16(1), int64(2), int64(1) - int64(2), nil},
	{int16(1), uint(2), sub(uint(1), uint(2)), nil},
	{int16(1), uint8(1<<8 - 1), int16(1) - int16(1<<8-1), nil},
	{int16(1), uint16(1<<16 - 1), sub(int16(1), int16(-1)), nil},
	{int16(1), uint32(1<<32 - 1), sub(uint32(1), uint32(1<<32-1)), nil},
	{int16(1), uint64(1<<64 - 1), sub(uint64(1), uint64(1<<64-1)), nil},
	{int16(1), float32(0.1), float32(1) - float32(0.1), nil},
	{int16(1), float64(0.1), float64(1) - float64(0.1), nil},
	{int16(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{int32(1), int(2), int(1) - int(2), nil},
	{int32(1), int8(2), int32(1) - int32(2), nil},
	{int32(1), int16(2), int32(1) - int32(2), nil},
	{int32(1), int32(2), int32(1) - int32(2), nil},
	{int32(1), int64(2), int64(1) - int64(2), nil},
	{int32(1), uint(2), sub(uint(1), uint(2)), nil},
	{int32(1), uint8(1<<8 - 1), int32(1) - int32(1<<8-1), nil},
	{int32(1), uint16(1<<16 - 1), sub(int32(1), int32(1<<16-1)), nil},
	{int32(1), uint32(1<<32 - 1), sub(int32(1), int32(-1)), nil},
	{int32(1), uint64(1<<64 - 1), sub(uint64(1), uint64(1<<64-1)), nil},
	{int32(1), float32(0.1), float32(1) - float32(0.1), nil},
	{int32(1), float64(0.1), float64(1) - float64(0.1), nil},
	{int32(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{int64(1), int(2), int64(1) - int64(2), nil},
	{int64(1), int8(2), int64(1) - int64(2), nil},
	{int64(1), int16(2), int64(1) - int64(2), nil},
	{int64(1), int32(2), int64(1) - int64(2), nil},
	{int64(1), int64(2), int64(1) - int64(2), nil},
	{int64(1), uint(2), int64(1) - int64(2), nil},
	{int64(1), uint8(1<<8 - 1), int64(1) - int64(1<<8-1), nil},
	{int64(1), uint16(1<<16 - 1), sub(int64(1), int64(1<<16-1)), nil},
	{int64(1), uint32(1<<32 - 1), sub(int64(1), int64(1<<32-1)), nil},
	{int64(1), uint64(1<<64 - 1), sub(int64(1), int64(-1)), nil},
	{int64(1), float32(0.1), float32(1) - float32(0.1), nil},
	{int64(1), float64(0.1), float64(1) - float64(0.1), nil},
	{int64(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{uint(1), int(2), int(1) - int(2), nil},
	{uint(1), int8(2), sub(uint(1), uint(2)), nil},
	{uint(1), int16(2), sub(uint(1), uint(2)), nil},
	{uint(1), int32(2), sub(uint(1), uint(2)), nil},
	{uint(1), int64(2), int64(1) - int64(2), nil},
	{uint(1), uint(2), sub(uint(1), uint(2)), nil},
	{uint(1), uint8(1<<8 - 1), sub(uint(1), uint(1<<8-1)), nil},
	{uint(1), uint16(1<<16 - 1), sub(uint(1), uint(1<<16-1)), nil},
	{uint(1), uint32(1<<32 - 1), sub(uint(1), uint(1<<32-1)), nil},
	{uint(1), uint64(1<<64 - 1), sub(uint(1), uint(1<<64-1)), nil},
	{uint(1), float32(0.1), float32(1) - float32(0.1), nil},
	{uint(1), float64(0.1), float64(1) - float64(0.1), nil},
	{uint(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{uint8(1), int(2), int(1) - int(2), nil},
	{uint8(1), int8(2), int8(1) - int8(2), nil},
	{uint8(1), int16(2), int16(1) - int16(2), nil},
	{uint8(1), int32(2), int32(1) - int32(2), nil},
	{uint8(1), int64(2), int64(1) - int64(2), nil},
	{uint8(1), uint(2), sub(uint(1), uint(2)), nil},
	{uint8(1), uint8(1<<8 - 1), sub(uint8(1), uint8(1<<8-1)), nil},
	{uint8(1), uint16(1<<16 - 1), sub(uint16(1), uint16(1<<16-1)), nil},
	{uint8(1), uint32(1<<32 - 1), sub(uint32(1), uint32(1<<32-1)), nil},
	{uint8(1), uint64(1<<64 - 1), sub(uint64(1), uint64(1<<64-1)), nil},
	{uint8(1), float32(0.1), float32(1) - float32(0.1), nil},
	{uint8(1), float64(0.1), float64(1) - float64(0.1), nil},
	{uint8(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{uint16(1), int(2), int(1) - int(2), nil},
	{uint16(1), int8(2), sub(uint16(1), uint16(2)), nil},
	{uint16(1), int16(2), int16(1) - int16(2), nil},
	{uint16(1), int32(2), int32(1) - int32(2), nil},
	{uint16(1), int64(2), int64(1) - int64(2), nil},
	{uint16(1), uint(2), sub(uint(1), uint(2)), nil},
	{uint16(1), uint8(1<<8 - 1), sub(uint16(1), uint16(1<<8-1)), nil},
	{uint16(1), uint16(1<<16 - 1), sub(uint16(1), uint16(1<<16-1)), nil},
	{uint16(1), uint32(1<<32 - 1), sub(uint32(1), uint32(1<<32-1)), nil},
	{uint16(1), uint64(1<<64 - 1), sub(uint64(1), uint64(1<<64-1)), nil},
	{uint16(1), float32(0.1), float32(1) - float32(0.1), nil},
	{uint16(1), float64(0.1), float64(1) - float64(0.1), nil},
	{uint16(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{uint32(1), int(2), int(1) - int(2), nil},
	{uint32(1), int8(2), sub(uint32(1), uint32(2)), nil},
	{uint32(1), int16(2), sub(uint32(1), uint32(2)), nil},
	{uint32(1), int32(2), int32(1) - int32(2), nil},
	{uint32(1), int64(2), int64(1) - int64(2), nil},
	{uint32(1), uint(2), sub(uint(1), uint(2)), nil},
	{uint32(1), uint8(1<<8 - 1), sub(uint32(1), uint32(1<<8-1)), nil},
	{uint32(1), uint16(1<<16 - 1), sub(uint32(1), uint32(1<<16-1)), nil},
	{uint32(1), uint32(1<<32 - 1), sub(uint32(1), uint32(1<<32-1)), nil},
	{uint32(1), uint64(1<<64 - 1), sub(uint64(1), uint64(1<<64-1)), nil},
	{uint32(1), float32(0.1), float32(1) - float32(0.1), nil},
	{uint32(1), float64(0.1), float64(1) - float64(0.1), nil},
	{uint32(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{uint64(1), int(2), int(1) - int(2), nil},
	{uint64(1), int8(2), sub(uint64(1), uint64(2)), nil},
	{uint64(1), int16(2), sub(uint64(1), uint64(2)), nil},
	{uint64(1), int32(2), sub(uint64(1), uint64(2)), nil},
	{uint64(1), int64(2), int64(1) - int64(2), nil},
	{uint64(1), uint(2), sub(uint64(1), uint64(2)), nil},
	{uint64(1), uint8(1<<8 - 1), sub(uint64(1), uint64(1<<8-1)), nil},
	{uint64(1), uint16(1<<16 - 1), sub(uint64(1), uint64(1<<16-1)), nil},
	{uint64(1), uint32(1<<32 - 1), sub(uint64(1), uint64(1<<32-1)), nil},
	{uint64(1), uint64(1<<64 - 1), sub(uint64(1), uint64(1<<64-1)), nil},
	{uint64(1), float32(0.1), float32(1) - float32(0.1), nil},
	{uint64(1), float64(0.1), float64(1) - float64(0.1), nil},
	{uint64(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{float32(1), int(2), float32(1) - float32(2), nil},
	{float32(1), int8(2), float32(1) - float32(2), nil},
	{float32(1), int16(2), float32(1) - float32(2), nil},
	{float32(1), int32(2), float32(1) - float32(2), nil},
	{float32(1), int64(2), float32(1) - float32(2), nil},
	{float32(1), uint(2), float32(1) - float32(2), nil},
	{float32(1), uint8(1<<8 - 1), float32(1) - float32(1<<8-1), nil},
	{float32(1), uint16(1<<16 - 1), sub(float32(1), float32(1<<16-1)), nil},
	{float32(1), uint32(1<<32 - 1), sub(float32(1), float32(1<<32-1)), nil},
	{float32(1), uint64(1<<64 - 1), sub(float32(1), float32(1<<64-1)), nil},
	{float32(1), float32(0.1), float32(1) - float32(0.1), nil},
	{float32(1), float64(0.1), float64(1) - float64(0.1), nil},
	{float32(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{float64(1), int(2), float64(1) - float64(2), nil},
	{float64(1), int8(2), float64(1) - float64(2), nil},
	{float64(1), int16(2), float64(1) - float64(2), nil},
	{float64(1), int32(2), float64(1) - float64(2), nil},
	{float64(1), int64(2), float64(1) - float64(2), nil},
	{float64(1), uint(2), float64(1) - float64(2), nil},
	{float64(1), uint8(1<<8 - 1), float64(1) - float64(1<<8-1), nil},
	{float64(1), uint16(1<<16 - 1), sub(float64(1), float64(1<<16-1)), nil},
	{float64(1), uint32(1<<32 - 1), sub(float64(1), float64(1<<32-1)), nil},
	{float64(1), uint64(1<<64 - 1), sub(float64(1), float64(1<<64-1)), nil},
	{float64(1), float32(0.1), float64(1) - float64(float32(0.1)), nil},
	{float64(1), float64(0.1), float64(1) - float64(0.1), nil},
	{float64(1), "2", "12", errors.New("subtract: unsupported type for 'string'")},

	{"1", int(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", int8(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", int16(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", int32(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", int64(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", uint(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", uint8(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", uint16(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", uint32(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", uint64(2), "12", errors.New("subtract: unsupported type for 'string'")},
	{"1", float32(-1), "1-1", errors.New("subtract: unsupported type for 'string'")},
	{"1", float64(-1), "1-1", errors.New("subtract: unsupported type for 'string'")},
	{"1", "2", "12", errors.New("subtract: unsupported type for 'string'")},
}

func TestSub(t *testing.T) {
	for i, c := range subTests {
		r, e := Subs(c.a, c.b)

		if c.e != nil {
			if fmt.Sprint(c.e) != fmt.Sprint(e) {
				t.Errorf("[%d] Sub(%T(%v), %T(%v)) = (%v), want: (%v)", i, c.a, c.a, c.b, c.b, e, c.e)
			}
			continue
		}

		if c.w != r {
			t.Errorf("[%d] Sub(%T(%v), %T(%v)) = %T(%v), want: %T(%v)", i, c.a, c.a, c.b, c.b, r, r, c.w, c.w)
		}
	}
}

func TestNegate(t *testing.T) {
	tests := []struct {
		name string
		a    any
		w    any
		e    error
	}{
		// Signed integers
		{"int", int(10), int(-10), nil},
		{"int8", int8(10), int8(-10), nil},
		{"int16", int16(10), int16(-10), nil},
		{"int32", int32(10), int32(-10), nil},
		{"int64", int64(10), int64(-10), nil},

		// Unsigned integers
		{"uint", uint(10), int(-10), nil},
		{"uint8", uint8(10), int8(-10), nil},
		{"uint16", uint16(10), int16(-10), nil},
		{"uint32", uint32(10), int32(-10), nil},
		{"uint64", uint64(10), int64(-10), nil},

		// Floating points
		{"float32", float32(3.14), float32(-3.14), nil},
		{"float64", float64(2.718), float64(-2.718), nil},

		// Unknown type
		{"string", "hello", "hello", errors.New("negate: unsupported type for 'string'")},
		{"nil", nil, nil, nil},
	}

	for i, c := range tests {
		r, e := Negate(c.a)

		if c.e != nil {
			if fmt.Sprint(c.e) != fmt.Sprint(e) {
				t.Errorf("[%d:%s] Negate(%T(%v)) = (%v), want: (%v)", i, c.name, c.a, c.a, e, c.e)
			}
			continue
		}

		if c.w != r {
			t.Errorf("[%d:%s] Negate(%T(%v)) = %T(%v), want: %T(%v)", i, c.name, c.a, c.a, r, r, c.w, c.w)
		}
	}
}

var mulTests = []struct {
	a, b, w any
	e       error
}{
	{int(1), int(2), int(1) * int(2), nil},
	{int(1), int8(2), int(1) * int(2), nil},
	{int(1), int16(2), int(1) * int(2), nil},
	{int(1), int32(2), int(1) * int(2), nil},
	{int(1), int64(2), int(1) * int(2), nil},
	{int(1), uint(2), int(1) * int(2), nil},
	{int(1), uint8(1<<8 - 1), int(1) * int(1<<8-1), nil},
	{int(1), uint16(1<<16 - 1), int(1) * int(1<<16-1), nil},
	{int(1), uint32(1<<32 - 1), int(1) * int(1<<32-1), nil},
	{int(1), uint64(1<<64 - 1), int(1) * int(-1), nil},
	{int(1), float32(0.1), float32(1) * float32(0.1), nil},
	{int(1), float64(0.1), float64(1) * float64(0.1), nil},
	{int(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{int8(1), int(2), int(1) * int(2), nil},
	{int8(1), int8(2), int8(1) * int8(2), nil},
	{int8(1), int16(2), int16(1) * int16(2), nil},
	{int8(1), int32(2), int32(1) * int32(2), nil},
	{int8(1), int64(2), int64(1) * int64(2), nil},
	{int8(1), uint(2), mul(uint(1), uint(2)), nil},
	{int8(1), uint8(1<<8 - 1), int8(1) * int8(-1), nil},
	{int8(1), uint16(1<<16 - 1), mul(uint16(1), uint16(1<<16-1)), nil},
	{int8(1), uint32(1<<32 - 1), mul(uint32(1), uint32(1<<32-1)), nil},
	{int8(1), uint64(1<<64 - 1), mul(uint64(1), uint64(1<<64-1)), nil},
	{int8(1), float32(0.1), float32(1) * float32(0.1), nil},
	{int8(1), float64(0.1), float64(1) * float64(0.1), nil},
	{int8(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{int16(1), int(2), int(1) * int(2), nil},
	{int16(1), int8(2), int16(1) * int16(2), nil},
	{int16(1), int16(2), int16(1) * int16(2), nil},
	{int16(1), int32(2), int32(1) * int32(2), nil},
	{int16(1), int64(2), int64(1) * int64(2), nil},
	{int16(1), uint(2), mul(uint(1), uint(2)), nil},
	{int16(1), uint8(1<<8 - 1), int16(1) * int16(1<<8-1), nil},
	{int16(1), uint16(1<<16 - 1), mul(int16(1), int16(-1)), nil},
	{int16(1), uint32(1<<32 - 1), mul(uint32(1), uint32(1<<32-1)), nil},
	{int16(1), uint64(1<<64 - 1), mul(uint64(1), uint64(1<<64-1)), nil},
	{int16(1), float32(0.1), float32(1) * float32(0.1), nil},
	{int16(1), float64(0.1), float64(1) * float64(0.1), nil},
	{int16(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{int32(1), int(2), int(1) * int(2), nil},
	{int32(1), int8(2), int32(1) * int32(2), nil},
	{int32(1), int16(2), int32(1) * int32(2), nil},
	{int32(1), int32(2), int32(1) * int32(2), nil},
	{int32(1), int64(2), int64(1) * int64(2), nil},
	{int32(1), uint(2), mul(uint(1), uint(2)), nil},
	{int32(1), uint8(1<<8 - 1), int32(1) * int32(1<<8-1), nil},
	{int32(1), uint16(1<<16 - 1), mul(int32(1), int32(1<<16-1)), nil},
	{int32(1), uint32(1<<32 - 1), mul(int32(1), int32(-1)), nil},
	{int32(1), uint64(1<<64 - 1), mul(uint64(1), uint64(1<<64-1)), nil},
	{int32(1), float32(0.1), float32(1) * float32(0.1), nil},
	{int32(1), float64(0.1), float64(1) * float64(0.1), nil},
	{int32(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{int64(1), int(2), int64(1) * int64(2), nil},
	{int64(1), int8(2), int64(1) * int64(2), nil},
	{int64(1), int16(2), int64(1) * int64(2), nil},
	{int64(1), int32(2), int64(1) * int64(2), nil},
	{int64(1), int64(2), int64(1) * int64(2), nil},
	{int64(1), uint(2), int64(1) * int64(2), nil},
	{int64(1), uint8(1<<8 - 1), int64(1) * int64(1<<8-1), nil},
	{int64(1), uint16(1<<16 - 1), mul(int64(1), int64(1<<16-1)), nil},
	{int64(1), uint32(1<<32 - 1), mul(int64(1), int64(1<<32-1)), nil},
	{int64(1), uint64(1<<64 - 1), mul(int64(1), int64(-1)), nil},
	{int64(1), float32(0.1), float32(1) * float32(0.1), nil},
	{int64(1), float64(0.1), float64(1) * float64(0.1), nil},
	{int64(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{uint(1), int(2), int(1) * int(2), nil},
	{uint(1), int8(2), mul(uint(1), uint(2)), nil},
	{uint(1), int16(2), mul(uint(1), uint(2)), nil},
	{uint(1), int32(2), mul(uint(1), uint(2)), nil},
	{uint(1), int64(2), int64(1) * int64(2), nil},
	{uint(1), uint(2), mul(uint(1), uint(2)), nil},
	{uint(1), uint8(1<<8 - 1), mul(uint(1), uint(1<<8-1)), nil},
	{uint(1), uint16(1<<16 - 1), mul(uint(1), uint(1<<16-1)), nil},
	{uint(1), uint32(1<<32 - 1), mul(uint(1), uint(1<<32-1)), nil},
	{uint(1), uint64(1<<64 - 1), mul(uint(1), uint(1<<64-1)), nil},
	{uint(1), float32(0.1), float32(1) * float32(0.1), nil},
	{uint(1), float64(0.1), float64(1) * float64(0.1), nil},
	{uint(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{uint8(1), int(2), int(1) * int(2), nil},
	{uint8(1), int8(2), int8(1) * int8(2), nil},
	{uint8(1), int16(2), int16(1) * int16(2), nil},
	{uint8(1), int32(2), int32(1) * int32(2), nil},
	{uint8(1), int64(2), int64(1) * int64(2), nil},
	{uint8(1), uint(2), mul(uint(1), uint(2)), nil},
	{uint8(1), uint8(1<<8 - 1), mul(uint8(1), uint8(1<<8-1)), nil},
	{uint8(1), uint16(1<<16 - 1), mul(uint16(1), uint16(1<<16-1)), nil},
	{uint8(1), uint32(1<<32 - 1), mul(uint32(1), uint32(1<<32-1)), nil},
	{uint8(1), uint64(1<<64 - 1), mul(uint64(1), uint64(1<<64-1)), nil},
	{uint8(1), float32(0.1), float32(1) * float32(0.1), nil},
	{uint8(1), float64(0.1), float64(1) * float64(0.1), nil},
	{uint8(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{uint16(1), int(2), int(1) * int(2), nil},
	{uint16(1), int8(2), mul(uint16(1), uint16(2)), nil},
	{uint16(1), int16(2), int16(1) * int16(2), nil},
	{uint16(1), int32(2), int32(1) * int32(2), nil},
	{uint16(1), int64(2), int64(1) * int64(2), nil},
	{uint16(1), uint(2), mul(uint(1), uint(2)), nil},
	{uint16(1), uint8(1<<8 - 1), mul(uint16(1), uint16(1<<8-1)), nil},
	{uint16(1), uint16(1<<16 - 1), mul(uint16(1), uint16(1<<16-1)), nil},
	{uint16(1), uint32(1<<32 - 1), mul(uint32(1), uint32(1<<32-1)), nil},
	{uint16(1), uint64(1<<64 - 1), mul(uint64(1), uint64(1<<64-1)), nil},
	{uint16(1), float32(0.1), float32(1) * float32(0.1), nil},
	{uint16(1), float64(0.1), float64(1) * float64(0.1), nil},
	{uint16(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{uint32(1), int(2), int(1) * int(2), nil},
	{uint32(1), int8(2), mul(uint32(1), uint32(2)), nil},
	{uint32(1), int16(2), mul(uint32(1), uint32(2)), nil},
	{uint32(1), int32(2), int32(1) * int32(2), nil},
	{uint32(1), int64(2), int64(1) * int64(2), nil},
	{uint32(1), uint(2), mul(uint(1), uint(2)), nil},
	{uint32(1), uint8(1<<8 - 1), mul(uint32(1), uint32(1<<8-1)), nil},
	{uint32(1), uint16(1<<16 - 1), mul(uint32(1), uint32(1<<16-1)), nil},
	{uint32(1), uint32(1<<32 - 1), mul(uint32(1), uint32(1<<32-1)), nil},
	{uint32(1), uint64(1<<64 - 1), mul(uint64(1), uint64(1<<64-1)), nil},
	{uint32(1), float32(0.1), float32(1) * float32(0.1), nil},
	{uint32(1), float64(0.1), float64(1) * float64(0.1), nil},
	{uint32(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{uint64(1), int(2), int(1) * int(2), nil},
	{uint64(1), int8(2), mul(uint64(1), uint64(2)), nil},
	{uint64(1), int16(2), mul(uint64(1), uint64(2)), nil},
	{uint64(1), int32(2), mul(uint64(1), uint64(2)), nil},
	{uint64(1), int64(2), int64(1) * int64(2), nil},
	{uint64(1), uint(2), mul(uint64(1), uint64(2)), nil},
	{uint64(1), uint8(1<<8 - 1), mul(uint64(1), uint64(1<<8-1)), nil},
	{uint64(1), uint16(1<<16 - 1), mul(uint64(1), uint64(1<<16-1)), nil},
	{uint64(1), uint32(1<<32 - 1), mul(uint64(1), uint64(1<<32-1)), nil},
	{uint64(1), uint64(1<<64 - 1), mul(uint64(1), uint64(1<<64-1)), nil},
	{uint64(1), float32(0.1), float32(1) * float32(0.1), nil},
	{uint64(1), float64(0.1), float64(1) * float64(0.1), nil},
	{uint64(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{float32(1), int(2), float32(1) * float32(2), nil},
	{float32(1), int8(2), float32(1) * float32(2), nil},
	{float32(1), int16(2), float32(1) * float32(2), nil},
	{float32(1), int32(2), float32(1) * float32(2), nil},
	{float32(1), int64(2), float32(1) * float32(2), nil},
	{float32(1), uint(2), float32(1) * float32(2), nil},
	{float32(1), uint8(1<<8 - 1), float32(1) * float32(1<<8-1), nil},
	{float32(1), uint16(1<<16 - 1), mul(float32(1), float32(1<<16-1)), nil},
	{float32(1), uint32(1<<32 - 1), mul(float32(1), float32(1<<32-1)), nil},
	{float32(1), uint64(1<<64 - 1), mul(float32(1), float32(1<<64-1)), nil},
	{float32(1), float32(0.1), float32(1) * float32(0.1), nil},
	{float32(1), float64(0.1), float64(1) * float64(0.1), nil},
	{float32(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{float64(1), int(2), float64(1) * float64(2), nil},
	{float64(1), int8(2), float64(1) * float64(2), nil},
	{float64(1), int16(2), float64(1) * float64(2), nil},
	{float64(1), int32(2), float64(1) * float64(2), nil},
	{float64(1), int64(2), float64(1) * float64(2), nil},
	{float64(1), uint(2), float64(1) * float64(2), nil},
	{float64(1), uint8(1<<8 - 1), float64(1) * float64(1<<8-1), nil},
	{float64(1), uint16(1<<16 - 1), mul(float64(1), float64(1<<16-1)), nil},
	{float64(1), uint32(1<<32 - 1), mul(float64(1), float64(1<<32-1)), nil},
	{float64(1), uint64(1<<64 - 1), mul(float64(1), float64(1<<64-1)), nil},
	{float64(1), float32(0.1), float64(1) * float64(float32(0.1)), nil},
	{float64(1), float64(0.1), float64(1) * float64(0.1), nil},
	{float64(1), "2", "12", errors.New("multiply: unsupported type for 'string'")},

	{"1", int(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", int8(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", int16(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", int32(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", int64(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", uint(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", uint8(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", uint16(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", uint32(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", uint64(2), "12", errors.New("multiply: unsupported type for 'string'")},
	{"1", float32(-1), "1-1", errors.New("multiply: unsupported type for 'string'")},
	{"1", float64(-1), "1-1", errors.New("multiply: unsupported type for 'string'")},
	{"1", "2", "12", errors.New("multiply: unsupported type for 'string'")},
}

func TestMultiply(t *testing.T) {
	for i, c := range mulTests {
		r, e := Multiplys(c.a, c.b)

		if c.e != nil {
			if fmt.Sprint(c.e) != fmt.Sprint(e) {
				t.Errorf("[%d] Multiply(%T(%v), %T(%v)) = (%v), want: (%v)", i, c.a, c.a, c.b, c.b, e, c.e)
			}
			continue
		}

		if c.w != r {
			t.Errorf("[%d] Multiply(%T(%v), %T(%v)) = %T(%v), want: %T(%v)", i, c.a, c.a, c.b, c.b, r, r, c.w, c.w)
		}
	}
}

var divTests = []struct {
	a, b, w any
	e       error
}{
	{int(1), int(2), int(1) / int(2), nil},
	{int(1), int8(2), int(1) / int(2), nil},
	{int(1), int16(2), int(1) / int(2), nil},
	{int(1), int32(2), int(1) / int(2), nil},
	{int(1), int64(2), int(1) / int(2), nil},
	{int(1), uint(2), int(1) / int(2), nil},
	{int(1), uint8(1<<8 - 1), int(1) / int(1<<8-1), nil},
	{int(1), uint16(1<<16 - 1), int(1) / int(1<<16-1), nil},
	{int(1), uint32(1<<32 - 1), int(1) / int(1<<32-1), nil},
	{int(1), uint64(1<<64 - 1), int(1) / int(-1), nil},
	{int(1), float32(0.1), float32(1) / float32(0.1), nil},
	{int(1), float64(0.1), float64(1) / float64(0.1), nil},
	{int(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{int8(1), int(2), int(1) / int(2), nil},
	{int8(1), int8(2), int8(1) / int8(2), nil},
	{int8(1), int16(2), int16(1) / int16(2), nil},
	{int8(1), int32(2), int32(1) / int32(2), nil},
	{int8(1), int64(2), int64(1) / int64(2), nil},
	{int8(1), uint(2), div(uint(1), uint(2)), nil},
	{int8(1), uint8(1<<8 - 1), int8(1) / int8(-1), nil},
	{int8(1), uint16(1<<16 - 1), div(uint16(1), uint16(1<<16-1)), nil},
	{int8(1), uint32(1<<32 - 1), div(uint32(1), uint32(1<<32-1)), nil},
	{int8(1), uint64(1<<64 - 1), div(uint64(1), uint64(1<<64-1)), nil},
	{int8(1), float32(0.1), float32(1) / float32(0.1), nil},
	{int8(1), float64(0.1), float64(1) / float64(0.1), nil},
	{int8(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{int16(1), int(2), int(1) / int(2), nil},
	{int16(1), int8(2), int16(1) / int16(2), nil},
	{int16(1), int16(2), int16(1) / int16(2), nil},
	{int16(1), int32(2), int32(1) / int32(2), nil},
	{int16(1), int64(2), int64(1) / int64(2), nil},
	{int16(1), uint(2), div(uint(1), uint(2)), nil},
	{int16(1), uint8(1<<8 - 1), int16(1) / int16(1<<8-1), nil},
	{int16(1), uint16(1<<16 - 1), div(int16(1), int16(-1)), nil},
	{int16(1), uint32(1<<32 - 1), div(uint32(1), uint32(1<<32-1)), nil},
	{int16(1), uint64(1<<64 - 1), div(uint64(1), uint64(1<<64-1)), nil},
	{int16(1), float32(0.1), float32(1) / float32(0.1), nil},
	{int16(1), float64(0.1), float64(1) / float64(0.1), nil},
	{int16(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{int32(1), int(2), int(1) / int(2), nil},
	{int32(1), int8(2), int32(1) / int32(2), nil},
	{int32(1), int16(2), int32(1) / int32(2), nil},
	{int32(1), int32(2), int32(1) / int32(2), nil},
	{int32(1), int64(2), int64(1) / int64(2), nil},
	{int32(1), uint(2), div(uint(1), uint(2)), nil},
	{int32(1), uint8(1<<8 - 1), int32(1) / int32(1<<8-1), nil},
	{int32(1), uint16(1<<16 - 1), div(int32(1), int32(1<<16-1)), nil},
	{int32(1), uint32(1<<32 - 1), div(int32(1), int32(-1)), nil},
	{int32(1), uint64(1<<64 - 1), div(uint64(1), uint64(1<<64-1)), nil},
	{int32(1), float32(0.1), float32(1) / float32(0.1), nil},
	{int32(1), float64(0.1), float64(1) / float64(0.1), nil},
	{int32(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{int64(1), int(2), int64(1) / int64(2), nil},
	{int64(1), int8(2), int64(1) / int64(2), nil},
	{int64(1), int16(2), int64(1) / int64(2), nil},
	{int64(1), int32(2), int64(1) / int64(2), nil},
	{int64(1), int64(2), int64(1) / int64(2), nil},
	{int64(1), uint(2), int64(1) / int64(2), nil},
	{int64(1), uint8(1<<8 - 1), int64(1) / int64(1<<8-1), nil},
	{int64(1), uint16(1<<16 - 1), div(int64(1), int64(1<<16-1)), nil},
	{int64(1), uint32(1<<32 - 1), div(int64(1), int64(1<<32-1)), nil},
	{int64(1), uint64(1<<64 - 1), div(int64(1), int64(-1)), nil},
	{int64(1), float32(0.1), float32(1) / float32(0.1), nil},
	{int64(1), float64(0.1), float64(1) / float64(0.1), nil},
	{int64(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{uint(1), int(2), int(1) / int(2), nil},
	{uint(1), int8(2), div(uint(1), uint(2)), nil},
	{uint(1), int16(2), div(uint(1), uint(2)), nil},
	{uint(1), int32(2), div(uint(1), uint(2)), nil},
	{uint(1), int64(2), int64(1) / int64(2), nil},
	{uint(1), uint(2), div(uint(1), uint(2)), nil},
	{uint(1), uint8(1<<8 - 1), div(uint(1), uint(1<<8-1)), nil},
	{uint(1), uint16(1<<16 - 1), div(uint(1), uint(1<<16-1)), nil},
	{uint(1), uint32(1<<32 - 1), div(uint(1), uint(1<<32-1)), nil},
	{uint(1), uint64(1<<64 - 1), div(uint(1), uint(1<<64-1)), nil},
	{uint(1), float32(0.1), float32(1) / float32(0.1), nil},
	{uint(1), float64(0.1), float64(1) / float64(0.1), nil},
	{uint(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{uint8(1), int(2), int(1) / int(2), nil},
	{uint8(1), int8(2), int8(1) / int8(2), nil},
	{uint8(1), int16(2), int16(1) / int16(2), nil},
	{uint8(1), int32(2), int32(1) / int32(2), nil},
	{uint8(1), int64(2), int64(1) / int64(2), nil},
	{uint8(1), uint(2), div(uint(1), uint(2)), nil},
	{uint8(1), uint8(1<<8 - 1), div(uint8(1), uint8(1<<8-1)), nil},
	{uint8(1), uint16(1<<16 - 1), div(uint16(1), uint16(1<<16-1)), nil},
	{uint8(1), uint32(1<<32 - 1), div(uint32(1), uint32(1<<32-1)), nil},
	{uint8(1), uint64(1<<64 - 1), div(uint64(1), uint64(1<<64-1)), nil},
	{uint8(1), float32(0.1), float32(1) / float32(0.1), nil},
	{uint8(1), float64(0.1), float64(1) / float64(0.1), nil},
	{uint8(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{uint16(1), int(2), int(1) / int(2), nil},
	{uint16(1), int8(2), div(uint16(1), uint16(2)), nil},
	{uint16(1), int16(2), int16(1) / int16(2), nil},
	{uint16(1), int32(2), int32(1) / int32(2), nil},
	{uint16(1), int64(2), int64(1) / int64(2), nil},
	{uint16(1), uint(2), div(uint(1), uint(2)), nil},
	{uint16(1), uint8(1<<8 - 1), div(uint16(1), uint16(1<<8-1)), nil},
	{uint16(1), uint16(1<<16 - 1), div(uint16(1), uint16(1<<16-1)), nil},
	{uint16(1), uint32(1<<32 - 1), div(uint32(1), uint32(1<<32-1)), nil},
	{uint16(1), uint64(1<<64 - 1), div(uint64(1), uint64(1<<64-1)), nil},
	{uint16(1), float32(0.1), float32(1) / float32(0.1), nil},
	{uint16(1), float64(0.1), float64(1) / float64(0.1), nil},
	{uint16(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{uint32(1), int(2), int(1) / int(2), nil},
	{uint32(1), int8(2), div(uint32(1), uint32(2)), nil},
	{uint32(1), int16(2), div(uint32(1), uint32(2)), nil},
	{uint32(1), int32(2), int32(1) / int32(2), nil},
	{uint32(1), int64(2), int64(1) / int64(2), nil},
	{uint32(1), uint(2), div(uint(1), uint(2)), nil},
	{uint32(1), uint8(1<<8 - 1), div(uint32(1), uint32(1<<8-1)), nil},
	{uint32(1), uint16(1<<16 - 1), div(uint32(1), uint32(1<<16-1)), nil},
	{uint32(1), uint32(1<<32 - 1), div(uint32(1), uint32(1<<32-1)), nil},
	{uint32(1), uint64(1<<64 - 1), div(uint64(1), uint64(1<<64-1)), nil},
	{uint32(1), float32(0.1), float32(1) / float32(0.1), nil},
	{uint32(1), float64(0.1), float64(1) / float64(0.1), nil},
	{uint32(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{uint64(1), int(2), int(1) / int(2), nil},
	{uint64(1), int8(2), div(uint64(1), uint64(2)), nil},
	{uint64(1), int16(2), div(uint64(1), uint64(2)), nil},
	{uint64(1), int32(2), div(uint64(1), uint64(2)), nil},
	{uint64(1), int64(2), int64(1) / int64(2), nil},
	{uint64(1), uint(2), div(uint64(1), uint64(2)), nil},
	{uint64(1), uint8(1<<8 - 1), div(uint64(1), uint64(1<<8-1)), nil},
	{uint64(1), uint16(1<<16 - 1), div(uint64(1), uint64(1<<16-1)), nil},
	{uint64(1), uint32(1<<32 - 1), div(uint64(1), uint64(1<<32-1)), nil},
	{uint64(1), uint64(1<<64 - 1), div(uint64(1), uint64(1<<64-1)), nil},
	{uint64(1), float32(0.1), float32(1) / float32(0.1), nil},
	{uint64(1), float64(0.1), float64(1) / float64(0.1), nil},
	{uint64(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{float32(1), int(2), float32(1) / float32(2), nil},
	{float32(1), int8(2), float32(1) / float32(2), nil},
	{float32(1), int16(2), float32(1) / float32(2), nil},
	{float32(1), int32(2), float32(1) / float32(2), nil},
	{float32(1), int64(2), float32(1) / float32(2), nil},
	{float32(1), uint(2), float32(1) / float32(2), nil},
	{float32(1), uint8(1<<8 - 1), float32(1) / float32(1<<8-1), nil},
	{float32(1), uint16(1<<16 - 1), div(float32(1), float32(1<<16-1)), nil},
	{float32(1), uint32(1<<32 - 1), div(float32(1), float32(1<<32-1)), nil},
	{float32(1), uint64(1<<64 - 1), div(float32(1), float32(1<<64-1)), nil},
	{float32(1), float32(0.1), float32(1) / float32(0.1), nil},
	{float32(1), float64(0.1), float64(1) / float64(0.1), nil},
	{float32(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{float64(1), int(2), float64(1) / float64(2), nil},
	{float64(1), int8(2), float64(1) / float64(2), nil},
	{float64(1), int16(2), float64(1) / float64(2), nil},
	{float64(1), int32(2), float64(1) / float64(2), nil},
	{float64(1), int64(2), float64(1) / float64(2), nil},
	{float64(1), uint(2), float64(1) / float64(2), nil},
	{float64(1), uint8(1<<8 - 1), float64(1) / float64(1<<8-1), nil},
	{float64(1), uint16(1<<16 - 1), div(float64(1), float64(1<<16-1)), nil},
	{float64(1), uint32(1<<32 - 1), div(float64(1), float64(1<<32-1)), nil},
	{float64(1), uint64(1<<64 - 1), div(float64(1), float64(1<<64-1)), nil},
	{float64(1), float32(0.1), float64(1) / float64(float32(0.1)), nil},
	{float64(1), float64(0.1), float64(1) / float64(0.1), nil},
	{float64(1), "2", "12", errors.New("divide: unsupported type for 'string'")},

	{"1", int(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", int8(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", int16(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", int32(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", int64(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", uint(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", uint8(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", uint16(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", uint32(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", uint64(2), "12", errors.New("divide: unsupported type for 'string'")},
	{"1", float32(-1), "1-1", errors.New("divide: unsupported type for 'string'")},
	{"1", float64(-1), "1-1", errors.New("divide: unsupported type for 'string'")},
	{"1", "2", "12", errors.New("divide: unsupported type for 'string'")},
}

func TestDivide(t *testing.T) {
	for i, c := range divTests {
		r, e := Divides(c.a, c.b)

		if c.e != nil {
			if fmt.Sprint(c.e) != fmt.Sprint(e) {
				t.Errorf("[%d] Divide(%T(%v), %T(%v)) = (%v), want: (%v)", i, c.a, c.a, c.b, c.b, e, c.e)
			}
			continue
		}

		if c.w != r {
			t.Errorf("[%d] Divide(%T(%v), %T(%v)) = %T(%v), want: %T(%v)", i, c.a, c.a, c.b, c.b, r, r, c.w, c.w)
		}
	}
}

var modTests = []struct {
	a, b, w any
	e       error
}{
	// {int(9), int(2), int(9) % int(2), nil},
	// {int(9), int8(2), int(9) % int(2), nil},
	// {int(9), int16(2), int(9) % int(2), nil},
	// {int(9), int32(2), int(9) % int(2), nil},
	// {int(9), int64(2), int(9) % int(2), nil},
	// {int(9), uint(2), int(9) % int(2), nil},
	// {int(9), uint8(1<<8 - 1), int(9) % int(1<<8-1), nil},
	// {int(9), uint16(1<<16 - 1), int(9) % int(1<<16-1), nil},
	// {int(9), uint32(1<<32 - 1), int(9) % int(1<<32-1), nil},
	// {int(9), uint64(1<<64 - 1), int(9) % int(-1), nil},
	{int(9), float32(2), int64(9) % int64(2), nil},
	{int(9), float64(2), int64(9) % int64(2), nil},
	{int(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{int8(9), int(2), int(9) % int(2), nil},
	{int8(9), int8(2), int8(9) % int8(2), nil},
	{int8(9), int16(2), int16(9) % int16(2), nil},
	{int8(9), int32(2), int32(9) % int32(2), nil},
	{int8(9), int64(2), int64(9) % int64(2), nil},
	{int8(9), uint(2), uint(9) % uint(2), nil},
	{int8(9), uint8(1<<8 - 1), int8(9) % int8(-1), nil},
	{int8(9), uint16(1<<16 - 1), uint16(9) % uint16(1<<16-1), nil},
	{int8(9), uint32(1<<32 - 1), uint32(9) % uint32(1<<32-1), nil},
	{int8(9), uint64(1<<64 - 1), uint64(9) % uint64(1<<64-1), nil},
	{int8(9), float32(2), int64(9) % int64(2), nil},
	{int8(9), float64(2), int64(9) % int64(2), nil},
	{int8(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{int16(9), int(2), int(9) % int(2), nil},
	{int16(9), int8(2), int16(9) % int16(2), nil},
	{int16(9), int16(2), int16(9) % int16(2), nil},
	{int16(9), int32(2), int32(9) % int32(2), nil},
	{int16(9), int64(2), int64(9) % int64(2), nil},
	{int16(9), uint(2), uint(9) % uint(2), nil},
	{int16(9), uint8(1<<8 - 1), int16(9) % int16(1<<8-1), nil},
	{int16(9), uint16(1<<16 - 1), int16(9) % int16(-1), nil},
	{int16(9), uint32(1<<32 - 1), uint32(9) % uint32(1<<32-1), nil},
	{int16(9), uint64(1<<64 - 1), uint64(9) % uint64(1<<64-1), nil},
	{int16(9), float32(2), int64(9) % int64(2), nil},
	{int16(9), float64(2), int64(9) % int64(2), nil},
	{int16(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{int32(9), int(2), int(9) % int(2), nil},
	{int32(9), int8(2), int32(9) % int32(2), nil},
	{int32(9), int16(2), int32(9) % int32(2), nil},
	{int32(9), int32(2), int32(9) % int32(2), nil},
	{int32(9), int64(2), int64(9) % int64(2), nil},
	{int32(9), uint(2), uint(9) % uint(2), nil},
	{int32(9), uint8(1<<8 - 1), int32(9) % int32(1<<8-1), nil},
	{int32(9), uint16(1<<16 - 1), int32(9) % int32(1<<16-1), nil},
	{int32(9), uint32(1<<32 - 1), int32(9) % int32(-1), nil},
	{int32(9), uint64(1<<64 - 1), uint64(9) % uint64(1<<64-1), nil},
	{int32(9), float32(2), int64(9) % int64(2), nil},
	{int32(9), float64(2), int64(9) % int64(2), nil},
	{int32(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{int64(9), int(2), int64(9) % int64(2), nil},
	{int64(9), int8(2), int64(9) % int64(2), nil},
	{int64(9), int16(2), int64(9) % int64(2), nil},
	{int64(9), int32(2), int64(9) % int64(2), nil},
	{int64(9), int64(2), int64(9) % int64(2), nil},
	{int64(9), uint(2), int64(9) % int64(2), nil},
	{int64(9), uint8(1<<8 - 1), int64(9) % int64(1<<8-1), nil},
	{int64(9), uint16(1<<16 - 1), int64(9) % int64(1<<16-1), nil},
	{int64(9), uint32(1<<32 - 1), int64(9) % int64(1<<32-1), nil},
	{int64(9), uint64(1<<64 - 1), int64(9) % int64(-1), nil},
	{int64(9), float32(2), int64(9) % int64(2), nil},
	{int64(9), float64(2), int64(9) % int64(2), nil},
	{int64(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{uint(9), int(2), int(9) % int(2), nil},
	{uint(9), int8(2), uint(9) % uint(2), nil},
	{uint(9), int16(2), uint(9) % uint(2), nil},
	{uint(9), int32(2), uint(9) % uint(2), nil},
	{uint(9), int64(2), int64(9) % int64(2), nil},
	{uint(9), uint(2), uint(9) % uint(2), nil},
	{uint(9), uint8(1<<8 - 1), uint(9) % uint(1<<8-1), nil},
	{uint(9), uint16(1<<16 - 1), uint(9) % uint(1<<16-1), nil},
	{uint(9), uint32(1<<32 - 1), uint(9) % uint(1<<32-1), nil},
	{uint(9), uint64(1<<64 - 1), uint(9) % uint(1<<64-1), nil},
	{uint(9), float32(2), int64(9) % int64(2), nil},
	{uint(9), float64(2), int64(9) % int64(2), nil},
	{uint(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{uint8(9), int(2), int(9) % int(2), nil},
	{uint8(9), int8(2), int8(9) % int8(2), nil},
	{uint8(9), int16(2), int16(9) % int16(2), nil},
	{uint8(9), int32(2), int32(9) % int32(2), nil},
	{uint8(9), int64(2), int64(9) % int64(2), nil},
	{uint8(9), uint(2), uint(9) % uint(2), nil},
	{uint8(9), uint8(1<<8 - 1), uint8(9) % uint8(1<<8-1), nil},
	{uint8(9), uint16(1<<16 - 1), uint16(9) % uint16(1<<16-1), nil},
	{uint8(9), uint32(1<<32 - 1), uint32(9) % uint32(1<<32-1), nil},
	{uint8(9), uint64(1<<64 - 1), uint64(9) % uint64(1<<64-1), nil},
	{uint8(9), float32(2), int64(9) % int64(2), nil},
	{uint8(9), float64(2), int64(9) % int64(2), nil},
	{uint8(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{uint16(9), int(2), int(9) % int(2), nil},
	{uint16(9), int8(2), uint16(9) % uint16(2), nil},
	{uint16(9), int16(2), int16(9) % int16(2), nil},
	{uint16(9), int32(2), int32(9) % int32(2), nil},
	{uint16(9), int64(2), int64(9) % int64(2), nil},
	{uint16(9), uint(2), uint(9) % uint(2), nil},
	{uint16(9), uint8(1<<8 - 1), uint16(9) % uint16(1<<8-1), nil},
	{uint16(9), uint16(1<<16 - 1), uint16(9) % uint16(1<<16-1), nil},
	{uint16(9), uint32(1<<32 - 1), uint32(9) % uint32(1<<32-1), nil},
	{uint16(9), uint64(1<<64 - 1), uint64(9) % uint64(1<<64-1), nil},
	{uint16(9), float32(2), int64(9) % int64(2), nil},
	{uint16(9), float64(2), int64(9) % int64(2), nil},
	{uint16(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{uint32(9), int(2), int(9) % int(2), nil},
	{uint32(9), int8(2), uint32(9) % uint32(2), nil},
	{uint32(9), int16(2), uint32(9) % uint32(2), nil},
	{uint32(9), int32(2), int32(9) % int32(2), nil},
	{uint32(9), int64(2), int64(9) % int64(2), nil},
	{uint32(9), uint(2), uint(9) % uint(2), nil},
	{uint32(9), uint8(1<<8 - 1), uint32(9) % uint32(1<<8-1), nil},
	{uint32(9), uint16(1<<16 - 1), uint32(9) % uint32(1<<16-1), nil},
	{uint32(9), uint32(1<<32 - 1), uint32(9) % uint32(1<<32-1), nil},
	{uint32(9), uint64(1<<64 - 1), uint64(9) % uint64(1<<64-1), nil},
	{uint32(9), float32(2), int64(9) % int64(2), nil},
	{uint32(9), float64(2), int64(9) % int64(2), nil},
	{uint32(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{uint64(9), int(2), int(9) % int(2), nil},
	{uint64(9), int8(2), uint64(9) % uint64(2), nil},
	{uint64(9), int16(2), uint64(9) % uint64(2), nil},
	{uint64(9), int32(2), uint64(9) % uint64(2), nil},
	{uint64(9), int64(2), int64(9) % int64(2), nil},
	{uint64(9), uint(2), uint64(9) % uint64(2), nil},
	{uint64(9), uint8(1<<8 - 1), uint64(9) % uint64(1<<8-1), nil},
	{uint64(9), uint16(1<<16 - 1), uint64(9) % uint64(1<<16-1), nil},
	{uint64(9), uint32(1<<32 - 1), uint64(9) % uint64(1<<32-1), nil},
	{uint64(9), uint64(1<<64 - 1), uint64(9) % uint64(1<<64-1), nil},
	{uint64(9), float32(2), int64(9) % int64(2), nil},
	{uint64(9), float64(2), int64(9) % int64(2), nil},
	{uint64(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{float32(9), int(2), int64(9) % int64(2), nil},
	{float32(9), int8(2), int64(9) % int64(2), nil},
	{float32(9), int16(2), int64(9) % int64(2), nil},
	{float32(9), int32(2), int64(9) % int64(2), nil},
	{float32(9), int64(2), int64(9) % int64(2), nil},
	{float32(9), uint(2), int64(9) % int64(2), nil},
	{float32(9), uint8(1<<8 - 1), int64(9) % int64(1<<8-1), nil},
	{float32(9), uint16(1<<16 - 1), int64(9) % int64(1<<16-1), nil},
	{float32(9), uint32(1<<32 - 1), int64(9) % int64(1<<32-1), nil},
	{float32(9), uint64(2), int64(9) % int64(2), nil},
	{float32(9), float32(2), int64(9) % int64(2), nil},
	{float32(9), float64(2), int64(9) % int64(2), nil},
	{float32(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{float64(9), int(2), int64(9) % int64(2), nil},
	{float64(9), int8(2), int64(9) % int64(2), nil},
	{float64(9), int16(2), int64(9) % int64(2), nil},
	{float64(9), int32(2), int64(9) % int64(2), nil},
	{float64(9), int64(2), int64(9) % int64(2), nil},
	{float64(9), uint(2), int64(9) % int64(2), nil},
	{float64(9), uint8(1<<8 - 1), int64(9) % int64(1<<8-1), nil},
	{float64(9), uint16(1<<16 - 1), int64(9) % int64(1<<16-1), nil},
	{float64(9), uint32(1<<32 - 1), int64(9) % int64(1<<32-1), nil},
	{float64(9), uint64(2), int64(9) % int64(2), nil},
	{float64(9), float32(2), int64(9) % int64(2), nil},
	{float64(9), float64(2), int64(9) % int64(2), nil},
	{float64(9), "2", "12", errors.New("mod: unsupported type for 'string'")},

	{"1", int(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", int8(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", int16(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", int32(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", int64(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", uint(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", uint8(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", uint16(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", uint32(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", uint64(2), "12", errors.New("mod: unsupported type for 'string'")},
	{"1", float32(-1), "1-1", errors.New("mod: unsupported type for 'string'")},
	{"1", float64(-1), "1-1", errors.New("mod: unsupported type for 'string'")},
	{"1", "2", "12", errors.New("mod: unsupported type for 'string'")},
}

func TestMod(t *testing.T) {
	for i, c := range modTests {
		r, e := Mod(c.a, c.b)

		if c.e != nil {
			if fmt.Sprint(c.e) != fmt.Sprint(e) {
				t.Errorf("[%d] Mod(%T(%v), %T(%v)) = (%v), want: (%v)", i, c.a, c.a, c.b, c.b, e, c.e)
			}
			continue
		}

		if c.w != r {
			t.Errorf("[%d] Mod(%T(%v), %T(%v)) = %T(%v), want: %T(%v)", i, c.a, c.a, c.b, c.b, r, r, c.w, c.w)
		}
	}
}
