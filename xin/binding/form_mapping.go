package binding

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/str"
)

var (
	errUnknownType = errors.New("unknown type")

	// ErrConvertMapStringSlice can not covert to map[string][]string
	ErrConvertMapStringSlice = errors.New("can not convert to map slices of strings")

	// ErrConvertToMapString can not convert to map[string]string
	ErrConvertToMapString = errors.New("can not convert to map of strings")
)

func mapURI(ptr any, m map[string][]string) error {
	return mapFormByTag(ptr, m, "uri")
}

func mapForm(ptr any, form map[string][]string) error {
	return mapFormByTag(ptr, form, "form")
}

func MapFormWithTag(ptr any, form map[string][]string, tag string) error {
	return mapFormByTag(ptr, form, tag)
}

var emptyField = reflect.StructField{}

func mapFormByTag(ptr any, form map[string][]string, tag string) error {
	// Check if ptr is a map
	ptrVal := reflect.ValueOf(ptr)
	var pointed any
	if ptrVal.Kind() == reflect.Ptr {
		ptrVal = ptrVal.Elem()
		pointed = ptrVal.Interface()
	}
	if ptrVal.Kind() == reflect.Map &&
		ptrVal.Type().Key().Kind() == reflect.String {
		if pointed != nil {
			ptr = pointed
		}
		return setFormMap(ptr, form)
	}

	return mappingByPtr(ptr, formSource(form), tag)
}

// setter tries to set value on a walking by fields of a struct
type setter interface {
	TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (isSet bool, err *FieldBindError)
}

type formSource map[string][]string

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form formSource) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (isSet bool, err *FieldBindError) {
	return setByForm(value, field, form, key, opt)
}

func mappingByPtr(ptr any, setter setter, tag string) error {
	bes := &FieldBindErrors{}
	mapping("", reflect.ValueOf(ptr), emptyField, setter, tag, bes)
	if bes.IsEmpty() {
		return nil
	}
	return bes
}

func mapping(prefix string, value reflect.Value, field reflect.StructField, setter setter, tag string, bes *FieldBindErrors) (isSet bool) {
	if field.Tag.Get(tag) == "-" { // just ignoring this field
		return
	}

	vKind := value.Kind()

	if vKind == reflect.Ptr {
		isNew := false
		vPtr := value
		if value.IsNil() {
			isNew = true
			vPtr = reflect.New(value.Type().Elem())
		}
		isSet = mapping(prefix, vPtr.Elem(), field, setter, tag, bes)
		if isNew && isSet {
			value.Set(vPtr)
		}
		return
	}

	if vKind != reflect.Struct || !field.Anonymous {
		ok, err := tryToSetValue(prefix, value, field, setter, tag)
		if err != nil {
			bes.AddError(err)
			return false
		}
		if ok {
			return true
		}
	}

	if vKind == reflect.Struct {
		prefix = getStructFieldPrefix(prefix, field, tag)

		tValue := value.Type()

		for i := 0; i < value.NumField(); i++ {
			sf := tValue.Field(i)
			if sf.PkgPath != "" && !sf.Anonymous { // unexported
				continue
			}

			ok := mapping(prefix, value.Field(i), sf, setter, tag, bes)
			isSet = isSet || ok
		}
		return
	}
	return
}

func getStructFieldPrefix(prefix string, field reflect.StructField, tag string) string {
	name := field.Tag.Get(tag)
	name, _ = head(name, ",")

	if name == "" && !field.Anonymous {
		name = field.Name
	}

	if name != "" {
		if prefix != "" {
			prefix += "."
		}
		prefix += name
	}
	return prefix
}

type setOptions struct {
	stripValue   bool
	defaultValue string
}

func tryToSetValue(prefix string, value reflect.Value, field reflect.StructField, setter setter, tag string) (isSet bool, err *FieldBindError) {
	var key, opts string

	key = field.Tag.Get(tag)
	key, opts = head(key, ",")

	if key == "" { // default value is FieldName
		key = field.Name
	}

	if key == "" { // when field is "emptyField" variable
		return false, nil
	}

	var setOpts setOptions
	var opt string
	for len(opts) > 0 {
		opt, opts = head(opts, ",")
		k, v := head(opt, "=")
		switch k {
		case "default":
			setOpts.defaultValue = v
		case "strip":
			setOpts.stripValue = true
		}
	}

	if prefix != "" {
		key = prefix + "." + key
	}

	return setter.TrySet(value, field, key, setOpts)
}

func alterFormKey(key string) string {
	dot := strings.IndexByte(key, '.')
	if dot < 0 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString(key[:dot])
	sb.WriteByte('[')
	key = key[dot+1:]
	for key != "" {
		dot = strings.IndexByte(key, '.')
		if dot < 0 {
			sb.WriteString(key)
			break
		}
		sb.WriteString(key[:dot])
		sb.WriteString("][")
		key = key[dot+1:]
	}
	sb.WriteByte(']')
	return sb.String()
}

func setByForm(value reflect.Value, field reflect.StructField, form map[string][]string, key string, opt setOptions) (isSet bool, be *FieldBindError) {
	vs, ok := form[key]
	if !ok {
		akey := alterFormKey(key)
		if akey != "" {
			vs, ok = form[akey]
		}
	}

	if ok && opt.stripValue {
		vs = str.RemoveEmpties(str.Strips(vs))
		ok = (len(vs) > 0)
	}

	if !ok && opt.defaultValue == "" {
		return
	}

	var err error
	switch value.Kind() {
	case reflect.Slice:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		isSet, err = true, setSlice(vs, value, field)
	case reflect.Array:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		if len(vs) != value.Len() {
			isSet, err = false, fmt.Errorf("%q is not valid value for %s", vs, value.Type().String())
		} else {
			isSet, err = true, setArray(vs, value, field)
		}
	default:
		var val string
		if !ok {
			val = opt.defaultValue
		}

		if len(vs) > 0 {
			val = vs[0]
		}
		isSet, err = true, setWithProperType(val, value, field)
	}

	if err != nil {
		be = &FieldBindError{
			Name:   key,
			Values: vs,
			Cause:  err,
		}
	}
	return
}

func setWithProperType(val string, value reflect.Value, field reflect.StructField) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, 0, value)
	case reflect.Int8:
		return setIntField(val, 8, value)
	case reflect.Int16:
		return setIntField(val, 16, value)
	case reflect.Int32:
		return setIntField(val, 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val, value)
		}
		return setIntField(val, 64, value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	case reflect.String:
		return setStringField(val, field, value)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(val, field, value)
		}
		return json.Unmarshal(str.UnsafeBytes(val), value.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal(str.UnsafeBytes(val), value.Addr().Interface())
	default:
		return errUnknownType
	}
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setStringField(val string, structField reflect.StructField, field reflect.Value) error {
	field.SetString(val)
	return nil
}

var timeFormats = []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02", "15:04:05"}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	timeFormat := structField.Tag.Get("time_format")

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	if timeFormat == "" {
		timeFormat = timeFormats[0]
		for _, tf := range timeFormats {
			if len(tf) == len(val) {
				timeFormat = tf
				break
			}
		}
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func setArray(vals []string, value reflect.Value, field reflect.StructField) error {
	for i, s := range vals {
		err := setWithProperType(s, value.Index(i), field)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(vals []string, value reflect.Value, field reflect.StructField) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setArray(vals, slice, field)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setTimeDuration(val string, value reflect.Value) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

func setFormMap(ptr any, form map[string][]string) error {
	el := reflect.TypeOf(ptr).Elem()

	if el.Kind() == reflect.Slice {
		ptrMap, ok := ptr.(map[string][]string)
		if !ok {
			return ErrConvertMapStringSlice
		}
		for k, v := range form {
			ptrMap[k] = v
		}

		return nil
	}

	ptrMap, ok := ptr.(map[string]string)
	if !ok {
		return ErrConvertToMapString
	}
	for k, v := range form {
		ptrMap[k] = v[len(v)-1] // pick last
	}

	return nil
}
