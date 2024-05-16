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
	TypeInt32    = reflect.TypeOf(int32(0))
	TypeInt64    = reflect.TypeOf(int64(0))
	TypeFloat32  = reflect.TypeOf(float32(0))
	TypeFloat64  = reflect.TypeOf(float64(0))
	TypeDuration = reflect.TypeOf(time.Duration(0))
	TypeTime     = reflect.TypeOf(time.Time{})
)

// IsIntType return true if v is an integer
func IsIntType(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
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
	return v != nil && reflect.TypeOf(v).Kind() == reflect.Array
}

// IsMapType return true if v is a map
func IsMapType(v any) bool {
	return v != nil && reflect.TypeOf(v).Kind() == reflect.Map
}

// IsPtrType return true if v is a ptr
func IsPtrType(v any) bool {
	return v != nil && reflect.TypeOf(v).Kind() == reflect.Ptr
}
