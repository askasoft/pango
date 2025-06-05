package cal

import (
	"testing"
)

func TestBitNot(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expected  any
		expectErr bool
	}{
		{"int", int(0b1010), int(^0b1010), false},
		{"int8", int8(0b1010), int8(^0b1010), false},
		{"uint8", uint8(0b1010), ^uint8(0b1010), false},
		{"string", "invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BitNot(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for input %v, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %v: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("BitNot(%v) = %v; want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestBitAnd(t *testing.T) {
	tests := []struct {
		name      string
		a, b      any
		expected  any
		expectErr bool
	}{
		{"int", int(0b1100), int(0b1010), int(0b1000), false},
		{"int8", int8(0b1100), int8(0b1010), int8(0b1000), false},
		{"uint8", uint8(0b1100), uint8(0b1010), uint8(0b1000), false},
		{"mixed", int(0b1100), int8(0b1010), int(0b1000), false},
		{"invalid", int(0b1100), "invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BitAnd(tt.a, tt.b)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for inputs %v and %v, got nil", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for inputs %v and %v: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("BitAnd(%v, %v) = %v; want %v", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func TestBitOr(t *testing.T) {
	tests := []struct {
		name      string
		a, b      any
		expected  any
		expectErr bool
	}{
		{"int", int(0b1100), int(0b1010), int(0b1110), false},
		{"int8", int8(0b1100), int8(0b1010), int8(0b1110), false},
		{"uint8", uint8(0b1100), uint8(0b1010), uint8(0b1110), false},
		{"mixed", int(0b1100), int8(0b1010), int(0b1110), false},
		{"invalid", int(0b1100), "invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BitOr(tt.a, tt.b)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for inputs %v and %v, got nil", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for inputs %v and %v: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("BitOr(%v, %v) = %v; want %v", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func TestBitXor(t *testing.T) {
	tests := []struct {
		name      string
		a, b      any
		expected  any
		expectErr bool
	}{
		{"int", int(0b1100), int(0b1010), int(0b0110), false},
		{"int8", int8(0b1100), int8(0b1010), int8(0b0110), false},
		{"uint8", uint8(0b1100), uint8(0b1010), uint8(0b0110), false},
		{"mixed", int(0b1100), int8(0b1010), int(0b0110), false},
		{"invalid", int(0b1100), "invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BitXor(tt.a, tt.b)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for inputs %v and %v, got nil", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for inputs %v and %v: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("BitXor(%v, %v) = %v; want %v", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func TestBitLeft(t *testing.T) {
	tests := []struct {
		name      string
		a, b      any
		expected  any
		expectErr bool
	}{
		{"int", int(0b0001), int(2), int(0b0100), false},
		{"int8", int8(0b0001), int8(2), int8(0b0100), false},
		{"uint8", uint8(0b0001), uint8(2), uint8(0b0100), false},
		{"mixed", int(0b0001), int8(2), int(0b0100), false},
		{"invalid", int(0b0001), "invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BitLeft(tt.a, tt.b)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for inputs %v and %v, got nil", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for inputs %v and %v: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("BitLeft(%v, %v) = %v; want %v", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func TestBitRight(t *testing.T) {
	tests := []struct {
		name      string
		a, b      any
		expected  any
		expectErr bool
	}{
		{"int", int(0b0100), int(2), int(0b0001), false},
		{"int8", int8(0b0100), int8(2), int8(0b0001), false},
		{"uint8", uint8(0b0100), uint8(2), uint8(0b0001), false},
		{"mixed", int(0b0100), int8(2), int(0b0001), false},
		{"invalid", int(0b0100), "invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BitRight(tt.a, tt.b)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for inputs %v and %v, got nil", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for inputs %v and %v: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("BitRight(%v, %v) = %v; want %v", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}
