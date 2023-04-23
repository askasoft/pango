package main

import (
	"fmt"

	"github.com/askasoft/pango/vad"
)

// MyStruct ..
type MyStruct struct {
	String string `validate:"is-awesome"`
}

// use a single instance of Validate, it caches struct info
var validate *vad.Validate

func main() {

	validate = vad.New()
	validate.RegisterValidation("is-awesome", ValidateMyVal)

	s := MyStruct{String: "awesome"}

	err := validate.Struct(s)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}

	s.String = "not awesome"
	err = validate.Struct(s)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}
}

// ValidateMyVal implements vad.Func
func ValidateMyVal(fl vad.FieldLevel) bool {
	return fl.Field().String() == "awesome"
}
