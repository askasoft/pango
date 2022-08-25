package validate

import (
	"bytes"
	"testing"
	"time"

	"github.com/pandafw/pango/vad"
	"github.com/stretchr/testify/assert"
)

type testInterface interface {
	String() string
}

type substructNoValidation struct {
	IString string
	IInt    int
}

type mapNoValidationSub map[string]substructNoValidation

type structNoValidationValues struct {
	substructNoValidation

	Boolean bool

	Uinteger   uint
	Integer    int
	Integer8   int8
	Integer16  int16
	Integer32  int32
	Integer64  int64
	Uinteger8  uint8
	Uinteger16 uint16
	Uinteger32 uint32
	Uinteger64 uint64

	Float32 float32
	Float64 float64

	String string

	Date time.Time

	Struct        substructNoValidation
	InlinedStruct struct {
		String  []string
		Integer int
	}

	IntSlice           []int
	IntPointerSlice    []*int
	StructPointerSlice []*substructNoValidation
	StructSlice        []substructNoValidation
	InterfaceSlice     []testInterface

	UniversalInterface any
	CustomInterface    testInterface

	FloatMap  map[string]float32
	StructMap mapNoValidationSub
}

func createNoValidationValues() structNoValidationValues {
	integer := 1
	s := structNoValidationValues{
		Boolean:            true,
		Uinteger:           1 << 29,
		Integer:            -10000,
		Integer8:           120,
		Integer16:          -20000,
		Integer32:          1 << 29,
		Integer64:          1 << 61,
		Uinteger8:          250,
		Uinteger16:         50000,
		Uinteger32:         1 << 31,
		Uinteger64:         1 << 62,
		Float32:            123.456,
		Float64:            123.456789,
		String:             "text",
		Date:               time.Time{},
		CustomInterface:    &bytes.Buffer{},
		Struct:             substructNoValidation{},
		IntSlice:           []int{-3, -2, 1, 0, 1, 2, 3},
		IntPointerSlice:    []*int{&integer},
		StructSlice:        []substructNoValidation{},
		UniversalInterface: 1.2,
		FloatMap: map[string]float32{
			"foo": 1.23,
			"bar": 232.323,
		},
		StructMap: mapNoValidationSub{
			"foo": substructNoValidation{},
			"bar": substructNoValidation{},
		},
		// StructPointerSlice []noValidationSub
		// InterfaceSlice     []testInterface
	}
	s.InlinedStruct.Integer = 1000
	s.InlinedStruct.String = []string{"first", "second"}
	s.IString = "substring"
	s.IInt = 987654
	return s
}

func TestValidateNoValidationValues(t *testing.T) {
	v := NewStructValidator()
	origin := createNoValidationValues()
	test := createNoValidationValues()
	empty := structNoValidationValues{}

	assert.Nil(t, v.ValidateStruct(test))
	assert.Nil(t, v.ValidateStruct(&test))
	assert.Nil(t, v.ValidateStruct(empty))
	assert.Nil(t, v.ValidateStruct(&empty))

	assert.Equal(t, origin, test)
}

type structNoValidationPointer struct {
	substructNoValidation

	Boolean bool

	Uinteger   *uint
	Integer    *int
	Integer8   *int8
	Integer16  *int16
	Integer32  *int32
	Integer64  *int64
	Uinteger8  *uint8
	Uinteger16 *uint16
	Uinteger32 *uint32
	Uinteger64 *uint64

	Float32 *float32
	Float64 *float64

	String *string

	Date *time.Time

	Struct *substructNoValidation

	IntSlice           *[]int
	IntPointerSlice    *[]*int
	StructPointerSlice *[]*substructNoValidation
	StructSlice        *[]substructNoValidation
	InterfaceSlice     *[]testInterface

	FloatMap  *map[string]float32
	StructMap *mapNoValidationSub
}

func TestValidateNoValidationPointers(t *testing.T) {
	v := NewStructValidator()

	//origin := createNoValidation_values()
	//test := createNoValidation_values()
	empty := structNoValidationPointer{}

	//assert.Nil(t, v.ValidateStruct(test))
	//assert.Nil(t, v.ValidateStruct(&test))
	assert.Nil(t, v.ValidateStruct(empty))
	assert.Nil(t, v.ValidateStruct(&empty))

	//assert.Equal(t, origin, test)
}

type Object map[string]any

func TestValidatePrimitives(t *testing.T) {
	v := NewStructValidator()

	obj := Object{"foo": "bar", "bar": 1}
	assert.NoError(t, v.ValidateStruct(obj))
	assert.NoError(t, v.ValidateStruct(&obj))
	assert.Equal(t, Object{"foo": "bar", "bar": 1}, obj)

	obj2 := []Object{{"foo": "bar", "bar": 1}, {"foo": "bar", "bar": 1}}
	assert.NoError(t, v.ValidateStruct(obj2))
	assert.NoError(t, v.ValidateStruct(&obj2))

	nu := 10
	assert.NoError(t, v.ValidateStruct(nu))
	assert.NoError(t, v.ValidateStruct(&nu))
	assert.Equal(t, 10, nu)

	str := "value"
	assert.NoError(t, v.ValidateStruct(str))
	assert.NoError(t, v.ValidateStruct(&str))
	assert.Equal(t, "value", str)
}

// structCustomValidation is a helper struct we use to check that
// custom validation can be registered on it.
// The `notone` binding directive is for custom validation and registered later.
type structCustomValidation struct {
	Integer int `binding:"notone"`
}

func notOne(f1 vad.FieldLevel) bool {
	if val, ok := f1.Field().Interface().(int); ok {
		return val != 1
	}
	return false
}

func TestValidatorEngine(t *testing.T) {
	v := NewStructValidator()

	// This validates that the function `notOne` matches
	// the expected function signature by `defaultValidator`
	// and by extension the validator library.
	engine, ok := v.Engine().(*vad.Validate)
	assert.True(t, ok)

	err := engine.RegisterValidation("notone", notOne)
	// Check that we can register custom validation without error
	assert.Nil(t, err)

	// Create an instance which will fail validation
	withOne := structCustomValidation{Integer: 1}
	errs := v.ValidateStruct(withOne)

	// Check that we got back non-nil errs
	assert.NotNil(t, errs)
	// Check that the error matches expectation
	assert.Error(t, errs, "", "", "notone")
}
