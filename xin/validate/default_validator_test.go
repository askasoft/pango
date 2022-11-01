package validate

import (
	"errors"
	"testing"
)

func TestSliceValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  SliceValidationError
		want string
	}{
		{"has nil elements", SliceValidationError{errors.New("test error"), nil}, "[0]: test error"},
		{"has zero elements", SliceValidationError{}, ""},
		{"has one element", SliceValidationError{errors.New("test one error")}, "[0]: test one error"},
		{"has two elements",
			SliceValidationError{
				errors.New("first error"),
				errors.New("second error"),
			},
			"[0]: first error\n[1]: second error",
		},
		{"has many elements",
			SliceValidationError{
				errors.New("first error"),
				errors.New("second error"),
				nil,
				nil,
				nil,
				errors.New("last error"),
			},
			"[0]: first error\n[1]: second error\n[5]: last error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("SliceValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultValidator(t *testing.T) {
	v := NewStructValidator()

	type exampleStruct struct {
		A string `binding:"max=8"`
		B int    `binding:"gt=0"`
	}
	tests := []struct {
		name    string
		v       StructValidator
		obj     any
		wantErr bool
	}{
		{"validate nil obj", v, nil, false},
		{"validate int obj", v, 3, false},
		{"validate struct failed-1", v, exampleStruct{A: "123456789", B: 1}, true},
		{"validate struct failed-2", v, exampleStruct{A: "12345678", B: 0}, true},
		{"validate struct passed", v, exampleStruct{A: "12345678", B: 1}, false},
		{"validate *struct failed-1", v, &exampleStruct{A: "123456789", B: 1}, true},
		{"validate *struct failed-2", v, &exampleStruct{A: "12345678", B: 0}, true},
		{"validate *struct passed", v, &exampleStruct{A: "12345678", B: 1}, false},
		{"validate []struct failed-1", v, []exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate []struct failed-2", v, []exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate []struct passed", v, []exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate []*struct failed-1", v, []*exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate []*struct failed-2", v, []*exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate []*struct passed", v, []*exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate *[]struct failed-1", v, &[]exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate *[]struct failed-2", v, &[]exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate *[]struct passed", v, &[]exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate *[]*struct failed-1", v, &[]*exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate *[]*struct failed-2", v, &[]*exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate *[]*struct passed", v, &[]*exampleStruct{{A: "12345678", B: 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.ValidateStruct(tt.obj); (err != nil) != tt.wantErr {
				t.Errorf("defaultValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
