package ref

// IsIntType return true if v is an integer
func IsIntType(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64:
	case uint, uint8, uint16, uint32, uint64:
		return true
	}

	return false
}

// IsFloatType return true if v is a float
func IsFloatType(v any) bool {
	switch v.(type) {
	case float32, float64:
		return true
	}

	return false
}

// IsComplexType return true if v is a complex
func IsComplexType(v any) bool {
	switch v.(type) {
	case complex64, complex128:
		return true
	}

	return false
}
