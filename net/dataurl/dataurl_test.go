package dataurl

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name      string
		mediaType string
		data      []byte
		expect    string
	}{
		{
			name:      "simple text",
			mediaType: "text/plain",
			data:      []byte("hello"),
			expect:    "data:text/plain;base64,aGVsbG8=",
		},
		{
			name:      "empty data",
			mediaType: "text/plain",
			data:      []byte{},
			expect:    "data:text/plain;base64,",
		},
		{
			name:      "binary data",
			mediaType: "application/octet-stream",
			data:      []byte{0x00, 0x01, 0x02},
			expect:    "data:application/octet-stream;base64,AAEC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.mediaType, tt.data)
			if got != tt.expect {
				t.Fatalf("expected %s, got %s", tt.expect, got)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectType  string
		expectData  []byte
		expectError bool
	}{
		{
			name:        "valid text",
			input:       "data:text/plain;base64,aGVsbG8=",
			expectType:  "text/plain",
			expectData:  []byte("hello"),
			expectError: false,
		},
		{
			name:        "valid binary",
			input:       "data:application/octet-stream;base64,AAEC",
			expectType:  "application/octet-stream",
			expectData:  []byte{0x00, 0x01, 0x02},
			expectError: false,
		},
		{
			name:        "invalid prefix",
			input:       "text/plain;base64,aGVsbG8=",
			expectError: true,
		},
		{
			name:        "invalid base64",
			input:       "data:text/plain;base64,@@@",
			expectError: true,
		},
		{
			name:        "missing comma",
			input:       "data:text/plain;base64aGVsbG8=",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt, data, err := Decode(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if mt != tt.expectType {
				t.Fatalf("expected media type %s, got %s", tt.expectType, mt)
			}

			if !bytes.Equal(data, tt.expectData) {
				t.Fatalf("expected data %v, got %v", tt.expectData, data)
			}
		})
	}
}
