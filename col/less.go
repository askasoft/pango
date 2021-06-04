package col

// LessString string less function
func LessString(a interface{}, b interface{}) bool {
	return a.(string) < b.(string)
}

// LessInt int less function
func LessInt(a interface{}, b interface{}) bool {
	return a.(int) < b.(int)
}

// LessInt32 int32 less function
func LessInt32(a interface{}, b interface{}) bool {
	return a.(int32) < b.(int32)
}

// LessInt64 int64 less function
func LessInt64(a interface{}, b interface{}) bool {
	return a.(int64) < b.(int64)
}

// LessFloat32 float32 less function
func LessFloat32(a interface{}, b interface{}) bool {
	return a.(float32) < b.(float32)
}

// LessFloat64 float64 less function
func LessFloat64(a interface{}, b interface{}) bool {
	return a.(float64) < b.(float64)
}
