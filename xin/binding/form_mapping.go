package binding

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/lut"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

var (
	ErrUnknownType = errors.New("unknown type")

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

	if ptrVal.Kind() == reflect.Map && ptrVal.Type().Key().Kind() == reflect.String {
		if pointed != nil {
			ptr = pointed
		}
		return setFormMap(ptr, form)
	}

	return mappingByPtr(ptr, formSource(form), tag)
}

// setter tries to set value on a walking by fields of a struct
type setter interface {
	TrySet(rsf reflect.StructField, field reflect.Value, key string, opt *setOptions) (isSet bool, err *FieldBindError)
}

type formSource map[string][]string

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form formSource) TrySet(rsf reflect.StructField, field reflect.Value, key string, opt *setOptions) (isSet bool, err *FieldBindError) {
	return setByForm(rsf, field, form, key, opt)
}

func mappingByPtr(ptr any, setter setter, tag string) error {
	bes := FieldBindErrors{}
	mapping("", reflect.ValueOf(ptr), emptyField, setter, tag, &bes)
	if len(bes) == 0 {
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
			*bes = append(*bes, err)
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
	valid    bool   // to valid utf-8
	strip    bool   // strip leading trailing whitespace, remove empty string
	ascii    bool   // convert full width runes to ascii rune
	lower    bool   // convert to lower case
	upper    bool   // convert to upper case
	defaults string // default value
}

func tryToSetValue(prefix string, value reflect.Value, rsf reflect.StructField, setter setter, tag string) (isSet bool, err *FieldBindError) {
	var key, opts string

	key = rsf.Tag.Get(tag)
	key, opts = head(key, ",")

	if key == "" { // default value is FieldName
		key = rsf.Name
	}

	if key == "" { // when field is "emptyField" variable
		return false, nil
	}

	setOpts := &setOptions{}
	var opt string
	for len(opts) > 0 {
		opt, opts = head(opts, ",")
		k, v := head(opt, "=")
		switch k {
		case "default":
			setOpts.defaults = v
		case "valid":
			setOpts.valid = true
		case "strip":
			setOpts.strip = true
		case "ascii":
			setOpts.ascii = true
		case "lower":
			setOpts.lower = true
		case "upper":
			setOpts.upper = true
		}
	}

	if prefix != "" {
		key = prefix + "." + key
	}

	return setter.TrySet(rsf, value, key, setOpts)
}

func alterFormKey(key string) string {
	dot := strings.IndexByte(key, '.')
	if dot < 0 {
		return ""
	}

	var sb strings.Builder

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

func trimFormValues(vs []string, opt *setOptions) []string {
	if opt.valid {
		vs = str.ToValidUTF8s(vs, "")
	}
	if opt.strip {
		vs = str.RemoveEmpties(str.Strips(vs))
	}
	if opt.ascii {
		vs = lut.ToASCIIs(vs)
	}
	if opt.lower {
		vs = str.ToLowers(vs)
	}
	if opt.upper {
		vs = str.ToUppers(vs)
	}
	return vs
}

func getFormValues(form map[string][]string, key string, opt *setOptions) ([]string, bool) {
	vs, ok := form[key]
	if !ok {
		akey := alterFormKey(key)
		if akey != "" {
			vs, ok = form[akey]
		}
	}

	if ok {
		vs = trimFormValues(vs, opt)
	}

	return vs, ok
}

func setByForm(rsf reflect.StructField, field reflect.Value, form map[string][]string, key string, opt *setOptions) (isSet bool, be *FieldBindError) {
	var (
		vs  []string
		ok  bool
		err error
	)

	switch field.Kind() {
	case reflect.Map:
		vs, isSet, err = setMap(field, form, key, opt)
	case reflect.Slice:
		vs, ok = getFormValues(form, key, opt)
		if !ok {
			if opt.defaults == "" {
				return
			}
			vs = []string{opt.defaults}
		}
		isSet, err = true, setSlice(rsf, field, vs)
	case reflect.Array:
		vs, ok = getFormValues(form, key, opt)
		if !ok {
			if opt.defaults == "" {
				return
			}
			vs = []string{opt.defaults}
		}
		if len(vs) != field.Len() {
			isSet, err = false, fmt.Errorf("%q is not valid value for %s", vs, field.Type().String())
		} else {
			isSet, err = true, setArray(rsf, field, vs)
		}
	default:
		vs, ok = getFormValues(form, key, opt)
		if !ok && opt.defaults == "" {
			return
		}

		var val string
		if len(vs) > 0 {
			val = vs[0]
		}
		if val == "" {
			val = opt.defaults
		}
		isSet, err = true, setWithProperType(rsf, field, val)
	}

	if err != nil {
		be = &FieldBindError{
			Err:    err,
			Field:  key,
			Values: vs,
		}
	}
	return
}

func setWithProperType(rsf reflect.StructField, field reflect.Value, val string) error {
	switch field.Kind() {
	case reflect.Int:
		return setIntField(field, val, 0)
	case reflect.Int8:
		return setIntField(field, val, 8)
	case reflect.Int16:
		return setIntField(field, val, 16)
	case reflect.Int32:
		return setIntField(field, val, 32)
	case reflect.Int64:
		switch field.Interface().(type) {
		case time.Duration:
			return setTimeDuration(field, val)
		}
		return setIntField(field, val, 64)
	case reflect.Uint:
		return setUintField(field, val, 0)
	case reflect.Uint8:
		return setUintField(field, val, 8)
	case reflect.Uint16:
		return setUintField(field, val, 16)
	case reflect.Uint32:
		return setUintField(field, val, 32)
	case reflect.Uint64:
		return setUintField(field, val, 64)
	case reflect.Bool:
		return setBoolField(field, val)
	case reflect.Float32:
		return setFloatField(field, val, 32)
	case reflect.Float64:
		return setFloatField(field, val, 64)
	case reflect.String:
		return setStringField(field, val)
	case reflect.Struct:
		switch field.Interface().(type) {
		case time.Time:
			return setTimeField(rsf, field, val)
		case multipart.FileHeader:
			return nil
		}
		return json.Unmarshal(str.UnsafeBytes(val), field.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal(str.UnsafeBytes(val), field.Addr().Interface())
	case reflect.Ptr:
		if !field.Elem().IsValid() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setWithProperType(rsf, field.Elem(), val)
	default:
		return ErrUnknownType
	}
}

func setIntField(field reflect.Value, val string, bitSize int) error {
	if val == "" {
		field.SetInt(0)
		return nil
	}

	n, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(n)
	}
	return err
}

func setUintField(field reflect.Value, val string, bitSize int) error {
	if val == "" {
		field.SetUint(0)
		return nil
	}

	n, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(n)
	}
	return err
}

func setBoolField(field reflect.Value, val string) error {
	if val == "" {
		field.SetBool(false)
		return nil
	}

	b, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(b)
	}
	return err
}

func setFloatField(field reflect.Value, val string, bitSize int) error {
	if val == "" {
		field.SetFloat(0)
		return nil
	}

	f, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(f)
	}
	return err
}

func setStringField(field reflect.Value, val string) error {
	field.SetString(val)
	return nil
}

func setTimeField(rsf reflect.StructField, field reflect.Value, val string) error {
	if val == "" {
		field.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	tf := strings.ToLower(rsf.Tag.Get("time_format"))

	switch tf {
	case "unix", "unixmilli", "unixmicro", "unixnano":
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		var t time.Time
		switch tf {
		case "unixnano":
			t = time.Unix(n/int64(time.Second), n%int64(time.Second))
		case "unixmicro":
			t = time.UnixMicro(n)
		case "unixmilli":
			t = time.UnixMilli(n)
		default:
			t = time.Unix(n, 0)
		}

		field.Set(reflect.ValueOf(t))
		return nil
	}

	loc := time.Local
	if isUTC, _ := strconv.ParseBool(rsf.Tag.Get("time_utc")); isUTC {
		loc = time.UTC
	}

	if locTag := rsf.Tag.Get("time_location"); locTag != "" {
		tl, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		loc = tl
	}

	if tf != "" {
		t, err := time.ParseInLocation(tf, val, loc)
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(t))
		return nil
	}

	t, err := tmu.ParseInLocation(val, loc)
	if err != nil {
		return err
	}

	field.Set(reflect.ValueOf(t))
	return nil
}

func setMap(field reflect.Value, form map[string][]string, key string, opt *setOptions) ([]string, bool, error) {
	vs, ok := getFormValues(form, key, opt)
	if ok {
		for _, val := range vs {
			if err := json.Unmarshal(str.UnsafeBytes(val), field.Addr().Interface()); err != nil {
				return vs, false, err
			}
		}
	}

	mt := field.Type()
	if field.IsNil() {
		field.Set(reflect.MakeMap(mt))
	}

	px1, px2 := key+"[", key+"."
	for k, ps := range form {
		if str.EndsWithByte(k, ']') {
			if !str.StartsWith(k, px1) {
				continue
			}
			k = k[len(px1) : len(k)-1]
		} else {
			if !str.StartsWith(k, px2) {
				continue
			}
			k = k[len(px2):]
		}

		ps = trimFormValues(ps, opt)
		if opt.strip && len(ps) == 0 {
			continue
		}

		vs = append(vs, ps...)

		kv := reflect.ValueOf(k)

		var val any
		switch len(ps) {
		case 0:
			val = ""
		case 1:
			val = ps[0]
		default:
			val = ps
		}

		vv := reflect.ValueOf(val)
		vt := reflect.TypeOf(val)
		if vt.Kind() != mt.Elem().Kind() {
			cv, err := ref.CastTo(val, mt.Elem())
			if err != nil {
				return vs, false, fmt.Errorf("map: invalid value type - %w", err)
			}
			vv = reflect.ValueOf(cv)
		}
		field.SetMapIndex(kv, vv)
		ok = true
	}

	return vs, ok, nil
}

func setSlice(rsf reflect.StructField, field reflect.Value, vals []string) error {
	slice := reflect.MakeSlice(field.Type(), len(vals), len(vals))
	if err := setArray(rsf, slice, vals); err != nil {
		return err
	}
	field.Set(slice)
	return nil
}

func setArray(rsf reflect.StructField, field reflect.Value, vals []string) error {
	for i, s := range vals {
		if err := setWithProperType(rsf, field.Index(i), s); err != nil {
			return err
		}
	}
	return nil
}

func setTimeDuration(field reflect.Value, val string) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(d))
	return nil
}

func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

func setFormMap(dict any, form map[string][]string) error {
	el := reflect.TypeOf(dict).Elem()

	if el.Kind() == reflect.Slice {
		m, ok := dict.(map[string][]string)
		if !ok {
			return ErrConvertMapStringSlice
		}

		for k, v := range form {
			m[k] = v
		}
		return nil
	}

	if m, ok := dict.(map[string]string); ok {
		for k, v := range form {
			if len(v) > 0 {
				m[k] = v[len(v)-1] // pick last
			}
		}
	} else {
		for k, v := range form {
			if len(v) > 0 {
				if _, err := ref.MapSet(dict, k, v[len(v)-1]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
