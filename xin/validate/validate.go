package validate

import (
	"github.com/askasoft/pango/vad"
)

// StructValidator is the minimal interface which needs to be implemented in
// order for it to be used as the validator engine for ensuring the correctness
// of the request.
type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is a slice|array, the validation should be performed travel on every element.
	// If the received type is not a struct or slice|array, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(any) error

	// SetTagName allows for changing of the default tag name of 'validate'
	SetTagName(name string)

	// Engine returns the underlying validator engine which powers the
	// StructValidator implementation.
	Engine() any
}

// NewStructValidator is the default validator which implements the StructValidator interface.
func NewStructValidator() StructValidator {
	v := &defaultValidator{
		engine: vad.New(),
	}
	return v
}
