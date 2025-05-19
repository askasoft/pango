package vad

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	defaultTagName        = "validate"
	utf8HexComma          = "0x2C"
	utf8Pipe              = "0x7C"
	tagSeparator          = ","
	orSeparator           = "|"
	tagKeySeparator       = "="
	structOnlyTag         = "structonly"
	noStructLevelTag      = "nostructlevel"
	omitempty             = "omitempty"
	isempty               = "isempty"
	requiredWithoutAllTag = "required_without_all"
	requiredWithoutTag    = "required_without"
	requiredWithTag       = "required_with"
	requiredWithAllTag    = "required_with_all"
	requiredIfTag         = "required_if"
	requiredUnlessTag     = "required_unless"
	excludedWithoutAllTag = "excluded_without_all"
	excludedWithoutTag    = "excluded_without"
	excludedWithTag       = "excluded_with"
	excludedWithAllTag    = "excluded_with_all"
	skipValidationTag     = "-"
	diveTag               = "dive"
	keysTag               = "keys"
	endKeysTag            = "endkeys"
	requiredTag           = "required"
	namespaceSeparator    = "."
	leftBracket           = "["
	rightBracket          = "]"
	restrictedTagChars    = ".[],|=+()`~!@#$%^&*\\\"/?<>{}"
	restrictedAliasErr    = "Alias '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
	restrictedTagErr      = "Tag '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
)

var (
	timeDurationType = reflect.TypeOf(time.Duration(0))
	timeType         = reflect.TypeOf(time.Time{})

	defaultCField   = &cField{namesEqual: true}
	utf8HexReplacer = strings.NewReplacer(utf8HexComma, ",", utf8Pipe, "|")
)

// FilterFunc is the type used to filter fields using
// StructFiltered(...) function.
// returning true results in the field being filtered/skipped from
// validation
type FilterFunc func(ns []byte) bool

// CustomTypeFunc allows for overriding or adding custom field type handler functions
// field = field value of the type to return a value to be validated
// example Valuer from sql drive see https://golang.org/src/database/sql/driver/types.go?s=1210:1293#L29
type CustomTypeFunc func(field reflect.Value) any

// TagNameFunc allows for adding of a custom tag name parser
type TagNameFunc func(field reflect.StructField) string

type internalValidationFuncWrapper struct {
	fn                 FuncEx
	runValidationOnNil bool
}

// Validate contains the validator settings and cache
type Validate struct {
	tagName          string
	pool             *sync.Pool
	tagNameFunc      TagNameFunc
	structLevelFuncs map[reflect.Type]StructLevelFunc
	customFuncs      map[reflect.Type]CustomTypeFunc
	aliases          map[string]string
	validations      map[string]internalValidationFuncWrapper
	tagCache         *tagCache
	structCache      *structCache
	hasCustomFuncs   bool
	hasTagNameFunc   bool
}

// New returns a new instance of 'validate' with sane defaults.
// Validate is designed to be thread-safe and used as a singleton instance.
// It caches information about your struct and validations,
// in essence only parsing your validation tags once per struct type.
// Using multiple instances neglects the benefit of caching.
func New() *Validate {
	v := &Validate{
		tagName:     defaultTagName,
		aliases:     make(map[string]string, len(bakedInAliases)),
		validations: make(map[string]internalValidationFuncWrapper, len(bakedInValidators)),
		tagCache:    &tagCache{tmap: map[string]*cTag{}},
		structCache: &structCache{smap: map[reflect.Type]*cStruct{}},
	}

	// must copy alias validators for separate validations to be used in each validator instance
	for k, val := range bakedInAliases {
		v.RegisterAlias(k, val)
	}

	// must copy validators for separate validations to be used in each instance
	for k, fn := range bakedInValidators {
		switch k {
		// these require that even if the value is nil that the validation should run, omitempty still overrides this behaviour
		case requiredIfTag, requiredUnlessTag, requiredWithTag, requiredWithAllTag, requiredWithoutTag, requiredWithoutAllTag,
			excludedWithTag, excludedWithAllTag, excludedWithoutTag, excludedWithoutAllTag:
			v.registerValidation(k, wrapFunc(k, fn), true, true)
		default:
			// no need to error check here, baked in will always be valid
			v.registerValidation(k, wrapFunc(k, fn), true, false)
		}
	}

	v.pool = &sync.Pool{
		New: func() any {
			return &validate{
				v:        v,
				ns:       make([]byte, 0, 64),
				actualNs: make([]byte, 0, 64),
				misc:     make([]byte, 32),
			}
		},
	}

	return v
}

// SetTagName allows for changing of the default tag name of 'validate'
func (v *Validate) SetTagName(name string) {
	v.tagName = name
}

// RegisterTagNameFunc registers a function to get alternate names for StructFields.
//
// eg. to use the names which have been specified for JSON representations of structs, rather than normal Go field names:
//
//	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
//	    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
//	    // skip if tag key says it should be ignored
//	    if name == "-" {
//	        return ""
//	    }
//	    return name
//	})
func (v *Validate) RegisterTagNameFunc(fn TagNameFunc) {
	v.tagNameFunc = fn
	v.hasTagNameFunc = true
}

// RegisterValidation adds a validation with the given tag
//
// NOTES:
// - if the key already exists, the previous validation function will be replaced.
// - this method is not thread-safe it is intended that these all be registered prior to any validation
func (v *Validate) RegisterValidation(tag string, fn Func, callValidationEvenIfNull ...bool) {
	v.RegisterValidationEx(tag, wrapFunc(tag, fn), callValidationEvenIfNull...)
}

// RegisterValidation adds a validation with the given tag
//
// NOTES:
// - if the key already exists, the previous validation function will be replaced.
// - this method is not thread-safe it is intended that these all be registered prior to any validation
func (v *Validate) RegisterValidationEx(tag string, fn FuncEx, callValidationEvenIfNull ...bool) {
	var nilCheckable bool
	if len(callValidationEvenIfNull) > 0 {
		nilCheckable = callValidationEvenIfNull[0]
	}
	v.registerValidation(tag, fn, false, nilCheckable)
}

func (v *Validate) registerValidation(tag string, fn FuncEx, bakedIn bool, nilCheckable bool) {
	if len(tag) == 0 {
		panic("tag cannot be empty")
	}

	if fn == nil {
		panic("function cannot be empty")
	}

	_, ok := restrictedTags[tag]
	if !bakedIn && (ok || strings.ContainsAny(tag, restrictedTagChars)) {
		panic(fmt.Sprintf(restrictedTagErr, tag))
	}
	v.validations[tag] = internalValidationFuncWrapper{fn: fn, runValidationOnNil: nilCheckable}
}

// RegisterAlias registers a mapping of a single validation tag that
// defines a common or complex set of validation(s) to simplify adding validation
// to structs.
//
// NOTE: this function is not thread-safe it is intended that these all be registered prior to any validation
func (v *Validate) RegisterAlias(alias, tags string) {
	_, ok := restrictedTags[alias]

	if ok || strings.ContainsAny(alias, restrictedTagChars) {
		panic(fmt.Sprintf(restrictedAliasErr, alias))
	}

	v.aliases[alias] = tags
}

// RegisterStructValidation registers a StructLevelFunc against a number of types.
//
// NOTE:
// - this method is not thread-safe it is intended that these all be registered prior to any validation
func (v *Validate) RegisterStructValidation(fn StructLevelFunc, types ...any) {
	if v.structLevelFuncs == nil {
		v.structLevelFuncs = make(map[reflect.Type]StructLevelFunc)
	}

	for _, t := range types {
		tv := reflect.ValueOf(t)
		if tv.Kind() == reflect.Ptr {
			t = reflect.Indirect(tv).Interface()
		}

		v.structLevelFuncs[reflect.TypeOf(t)] = fn
	}
}

// RegisterCustomTypeFunc registers a CustomTypeFunc against a number of types
//
// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
func (v *Validate) RegisterCustomTypeFunc(fn CustomTypeFunc, types ...any) {
	if v.customFuncs == nil {
		v.customFuncs = make(map[reflect.Type]CustomTypeFunc)
	}

	for _, t := range types {
		v.customFuncs[reflect.TypeOf(t)] = fn
	}

	v.hasCustomFuncs = true
}

// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (v *Validate) Struct(s any) (err error) {
	val := reflect.ValueOf(s)
	top := val

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType) {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	// good to validate
	vd := v.pool.Get().(*validate)
	vd.top = top
	vd.isPartial = false
	// vd.hasExcludes = false // only need to reset in StructPartial and StructExcept

	vd.validateStruct(top, val, val.Type(), vd.ns[0:0], vd.actualNs[0:0], nil)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}

	v.pool.Put(vd)
	return
}

// StructFiltered validates a structs exposed fields, that pass the FilterFunc check and automatically validates
// nested structs, unless otherwise specified.
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (v *Validate) StructFiltered(s any, fn FilterFunc) (err error) {
	val := reflect.ValueOf(s)
	top := val

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType) {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	// good to validate
	vd := v.pool.Get().(*validate)
	vd.top = top
	vd.isPartial = true
	vd.ffn = fn
	// vd.hasExcludes = false // only need to reset in StructPartial and StructExcept

	vd.validateStruct(top, val, val.Type(), vd.ns[0:0], vd.actualNs[0:0], nil)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}

	v.pool.Put(vd)
	return
}

// StructPartial validates the fields passed in only, ignoring all others.
// Fields may be provided in a namespaced fashion relative to the  struct provided
// eg. NestedStruct.Field or NestedArrayField[0].Struct.Name
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (v *Validate) StructPartial(s any, fields ...string) (err error) {
	val := reflect.ValueOf(s)
	top := val

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType) {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	// good to validate
	vd := v.pool.Get().(*validate)
	vd.top = top
	vd.isPartial = true
	vd.ffn = nil
	vd.hasExcludes = false
	vd.includeExclude = make(map[string]struct{})

	typ := val.Type()
	name := typ.Name()

	for _, k := range fields {
		flds := strings.Split(k, namespaceSeparator)
		if len(flds) > 0 {
			vd.misc = append(vd.misc[0:0], name...)
			// Don't append empty name for unnamed structs
			if len(vd.misc) != 0 {
				vd.misc = append(vd.misc, '.')
			}

			for _, s := range flds {
				idx := strings.Index(s, leftBracket)

				if idx != -1 {
					for idx != -1 {
						vd.misc = append(vd.misc, s[:idx]...)
						vd.includeExclude[string(vd.misc)] = struct{}{}

						idx2 := strings.Index(s, rightBracket)
						idx2++
						vd.misc = append(vd.misc, s[idx:idx2]...)
						vd.includeExclude[string(vd.misc)] = struct{}{}
						s = s[idx2:]
						idx = strings.Index(s, leftBracket)
					}
				} else {
					vd.misc = append(vd.misc, s...)
					vd.includeExclude[string(vd.misc)] = struct{}{}
				}

				vd.misc = append(vd.misc, '.')
			}
		}
	}

	vd.validateStruct(top, val, typ, vd.ns[0:0], vd.actualNs[0:0], nil)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}

	v.pool.Put(vd)
	return
}

// StructExcept validates all fields except the ones passed in.
// Fields may be provided in a namespaced fashion relative to the  struct provided
// i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (v *Validate) StructExcept(s any, fields ...string) (err error) {
	val := reflect.ValueOf(s)
	top := val

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType) {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	// good to validate
	vd := v.pool.Get().(*validate)
	vd.top = top
	vd.isPartial = true
	vd.ffn = nil
	vd.hasExcludes = true
	vd.includeExclude = make(map[string]struct{})

	typ := val.Type()
	name := typ.Name()

	for _, key := range fields {
		vd.misc = vd.misc[0:0]

		if len(name) > 0 {
			vd.misc = append(vd.misc, name...)
			vd.misc = append(vd.misc, '.')
		}

		vd.misc = append(vd.misc, key...)
		vd.includeExclude[string(vd.misc)] = struct{}{}
	}

	vd.validateStruct(top, val, typ, vd.ns[0:0], vd.actualNs[0:0], nil)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}

	v.pool.Put(vd)
	return
}

// Map validates map data form a map of tags
func (v *Validate) Map(data map[string]any, rules map[string]any) map[string]any {
	errs := make(map[string]any)
	for field, rule := range rules {
		if reflect.ValueOf(rule).Kind() == reflect.Map && reflect.ValueOf(data[field]).Kind() == reflect.Map {
			err := v.Map(data[field].(map[string]any), rule.(map[string]any))
			if len(err) > 0 {
				errs[field] = err
			}
		} else if reflect.ValueOf(rule).Kind() == reflect.Map {
			errs[field] = errors.New("The field: '" + field + "' is not a map to dive")
		} else {
			err := v.Field(field, data[field], rule.(string))
			if err != nil {
				errs[field] = err
			}
		}
	}
	return errs
}

// Field validates a single variable using tag style validation.
// eg.
// var i int
// validate.Field("a", i, "gt=1,lt=10")
//
// WARNING: a struct can be passed for validation eg. time.Time is a struct or
// if you have a custom type and have registered a custom type handler, so must
// allow it; however unforeseen validations will occur if trying to validate a
// struct that is meant to be passed to 'validate.Struct'
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
// validate Array, Slice and maps fields which may contain more than one error
func (v *Validate) Field(field string, value any, tag string) (err error) {
	if len(tag) == 0 || tag == skipValidationTag {
		return nil
	}

	ctag := v.fetchCacheTag(tag)
	val := reflect.ValueOf(value)
	vd := v.pool.Get().(*validate)
	vd.top = val
	vd.isPartial = false

	cf := defaultCField
	if len(field) > 0 {
		cf = &cField{name: field, altName: field, namesEqual: true}
	}

	vd.traverseField(val, val, vd.ns[0:0], vd.actualNs[0:0], cf, ctag)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}
	v.pool.Put(vd)
	return
}

// Var validates a single variable using tag style validation.
// eg.
// var i int
// validate.Var(i, "gt=1,lt=10")
//
// WARNING: a struct can be passed for validation eg. time.Time is a struct or
// if you have a custom type and have registered a custom type handler, so must
// allow it; however unforeseen validations will occur if trying to validate a
// struct that is meant to be passed to 'validate.Struct'
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
// validate Array, Slice and maps fields which may contain more than one error
func (v *Validate) Var(value any, tag string) error {
	return v.Field("", value, tag)
}

// VarWithValue validates a single variable, against another variable/field's value using tag style validation
// eg.
// s1 := "abcd"
// s2 := "abcd"
// validate.VarWithValue(s1, s2, "eqcsfield") // returns true
//
// WARNING: a struct can be passed for validation eg. time.Time is a struct or
// if you have a custom type and have registered a custom type handler, so must
// allow it; however unforeseen validations will occur if trying to validate a
// struct that is meant to be passed to 'validate.Struct'
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
// validate Array, Slice and maps fields which may contain more than one error
func (v *Validate) VarWithValue(field any, other any, tag string) (err error) {
	if len(tag) == 0 || tag == skipValidationTag {
		return nil
	}

	ctag := v.fetchCacheTag(tag)
	otherVal := reflect.ValueOf(other)
	vd := v.pool.Get().(*validate)
	vd.top = otherVal
	vd.isPartial = false
	vd.traverseField(otherVal, reflect.ValueOf(field), vd.ns[0:0], vd.actualNs[0:0], defaultCField, ctag)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}
	v.pool.Put(vd)
	return
}
