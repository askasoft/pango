package ref

import (
	"reflect"
	"time"
)

var (
	TypeString   = reflect.TypeOf("")
	TypeBool     = reflect.TypeOf(false)
	TypeByte     = reflect.TypeOf(byte(0))
	TypeRune     = reflect.TypeOf(rune(0))
	TypeInt      = reflect.TypeOf(int(0))
	TypeInt8     = reflect.TypeOf(int8(0))
	TypeInt16    = reflect.TypeOf(int16(0))
	TypeInt32    = reflect.TypeOf(int32(0))
	TypeInt64    = reflect.TypeOf(int64(0))
	TypeUint     = reflect.TypeOf(uint(0))
	TypeUint8    = reflect.TypeOf(uint8(0))
	TypeUint16   = reflect.TypeOf(uint16(0))
	TypeUint32   = reflect.TypeOf(uint32(0))
	TypeUint64   = reflect.TypeOf(uint64(0))
	TypeFloat32  = reflect.TypeOf(float32(0))
	TypeFloat64  = reflect.TypeOf(float64(0))
	TypeDuration = reflect.TypeOf(time.Duration(0))
	TypeTime     = reflect.TypeOf(time.Time{})

	TypeAnys      = reflect.TypeOf([]any{})
	TypeStrings   = reflect.TypeOf([]string{})
	TypeBools     = reflect.TypeOf([]bool{})
	TypeBytes     = reflect.TypeOf([]byte{})
	TypeRunes     = reflect.TypeOf([]rune{})
	TypeInts      = reflect.TypeOf([]int{})
	TypeInt8s     = reflect.TypeOf([]int8{})
	TypeInt16s    = reflect.TypeOf([]int16{})
	TypeInt32s    = reflect.TypeOf([]int32{})
	TypeInt64s    = reflect.TypeOf([]int64{})
	TypeUints     = reflect.TypeOf([]uint{})
	TypeUint8s    = reflect.TypeOf([]uint8{})
	TypeUint16s   = reflect.TypeOf([]uint16{})
	TypeUint32s   = reflect.TypeOf([]uint32{})
	TypeUint64s   = reflect.TypeOf([]uint64{})
	TypeFloat32s  = reflect.TypeOf([]float32{})
	TypeFloat64s  = reflect.TypeOf([]float64{})
	TypeDurations = reflect.TypeOf([]time.Duration{})
	TypeTimes     = reflect.TypeOf([]time.Time{})

	TypeAnyStrMap     = reflect.TypeOf(map[any]string{})
	TypeStrStrMap     = reflect.TypeOf(map[string]string{})
	TypeIntStrMap     = reflect.TypeOf(map[int]string{})
	TypeInt8StrMap    = reflect.TypeOf(map[int8]string{})
	TypeInt16StrMap   = reflect.TypeOf(map[int16]string{})
	TypeInt32StrMap   = reflect.TypeOf(map[int32]string{})
	TypeInt64StrMap   = reflect.TypeOf(map[int64]string{})
	TypeUintStrMap    = reflect.TypeOf(map[uint]string{})
	TypeUint8StrMap   = reflect.TypeOf(map[uint8]string{})
	TypeUint16StrMap  = reflect.TypeOf(map[uint16]string{})
	TypeUint32StrMap  = reflect.TypeOf(map[uint32]string{})
	TypeUint64StrMap  = reflect.TypeOf(map[uint64]string{})
	TypeFloat32StrMap = reflect.TypeOf(map[float32]string{})
	TypeFloat64StrMap = reflect.TypeOf(map[float64]string{})
)

// IsStringType return true if v is an string
func IsStringType(v any) bool {
	switch v.(type) {
	case string:
		return true
	default:
		return false
	}
}

// IsIntType return true if v is an integer
func IsIntType(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return true
	default:
		return false
	}
}

// IsUintType return true if v is an integer
func IsUintType(v any) bool {
	switch v.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

// IsFloatType return true if v is a float
func IsFloatType(v any) bool {
	switch v.(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}

// IsComplexType return true if v is a complex
func IsComplexType(v any) bool {
	switch v.(type) {
	case complex64, complex128:
		return true
	default:
		return false
	}
}

// IsArrayType return true if v is a array
func IsArrayType(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Array
}

// IsSliceType return true if v is a array
func IsSliceType(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Slice
}

// IsMapType return true if v is a map
func IsMapType(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Map
}

// IsPtrType return true if v is a ptr
func IsPtrType(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}
