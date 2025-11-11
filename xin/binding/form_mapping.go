package binding

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/lut"
	"github.com/askasoft/pango/mag"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
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

	if ptrVal.Kind() == reflect.Map {
		if pointed != nil {
			ptr = pointed
		}
		return setFormMap(ptr, form)
	}

	return mappingByPtr(ptr, formSource(form), tag)
}

// setter tries to set value on a walking by fields of a struct
type setter interface {
	TrySet(field reflect.Value, key string, opts options) (isSet bool, err *FieldBindError)
}

type formSource map[string][]string

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form formSource) TrySet(field reflect.Value, key string, opts options) (isSet bool, err *FieldBindError) {
	return setByForm(field, form, key, opts)
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
	name := str.SubstrBeforeByte(field.Tag.Get(tag), ',')

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

type options struct {
	tag      reflect.StructTag
	defaults string // default value
	split    bool   // split to slice
	valid    bool   // to valid utf-8
	strip    bool   // strip leading trailing whitespace, remove empty string
	ascii    bool   // convert full width runes to ascii rune
	lower    bool   // convert to lower case
	upper    bool   // convert to upper case
}

func tryToSetValue(prefix string, value reflect.Value, rsf reflect.StructField, setter setter, tag string) (isSet bool, err *FieldBindError) {
	key, next, _ := str.CutByte(rsf.Tag.Get(tag), ',')

	if key == "" { // default value is FieldName
		key = rsf.Name
	}

	if key == "" { // when field is "emptyField" variable
		return false, nil
	}

	var opts options

	opts.tag = rsf.Tag

	var opt string
	for next != "" {
		opt, next, _ = str.CutByte(next, ',')
		k, v, _ := str.CutByte(opt, '=')
		switch k {
		case "default":
			opts.defaults = v
		case "split":
			opts.split = true
		case "valid":
			opts.valid = true
		case "strip":
			opts.strip = true
		case "ascii":
			opts.ascii = true
		case "lower":
			opts.lower = true
		case "upper":
			opts.upper = true
		}
	}

	if prefix != "" {
		key = prefix + "." + key
	}

	return setter.TrySet(value, key, opts)
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

func trimFormValues(vs []string, opts options) []string {
	if opts.split {
		var ss []string
		for _, v := range vs {
			ss = append(ss, str.Fields(v)...)
		}
		vs = ss
	}
	if opts.valid {
		vs = str.ToValidUTF8s(vs, "")
	}
	if opts.strip {
		vs = str.Strips(vs)
	}
	if opts.ascii {
		vs = lut.ToASCIIs(vs)
	}
	if opts.lower {
		vs = str.ToLowers(vs)
	}
	if opts.upper {
		vs = str.ToUppers(vs)
	}
	return vs
}

func getFormValues(form map[string][]string, key string, opts options) ([]string, bool) {
	vs, ok := form[key]
	if !ok {
		akey := alterFormKey(key)
		if akey != "" {
			vs, ok = form[akey]
		}
	}

	if ok {
		vs = trimFormValues(vs, opts)
	}

	return vs, ok
}

func getFormValuesOrDefaults(form map[string][]string, key string, opts options) (vs []string, ok bool) {
	vs, ok = getFormValues(form, key, opts)
	if !ok {
		if opts.defaults == "" {
			return
		}
		ok, vs = true, trimFormValues([]string{opts.defaults}, opts)
	}
	return
}

func setByForm(field reflect.Value, form map[string][]string, key string, opts options) (isSet bool, be *FieldBindError) {
	var (
		vs  []string
		ok  bool
		err error
	)

	switch field.Kind() {
	case reflect.Map:
		vs, isSet, err = setMap(field, form, key, opts)
	case reflect.Slice:
		vs, ok = getFormValuesOrDefaults(form, key, opts)
		if !ok {
			return
		}
		isSet, err = true, setSlice(field, vs, opts)
	case reflect.Array:
		vs, ok = getFormValuesOrDefaults(form, key, opts)
		if !ok {
			return
		}
		if len(vs) > field.Len() {
			isSet, err = false, fmt.Errorf("form: %q is not valid value for %s", vs, field.Type().String())
		} else {
			isSet, err = true, setArray(field, vs, opts)
		}
	default:
		vs, ok = getFormValuesOrDefaults(form, key, opts)
		if !ok {
			return
		}

		val := asg.First(vs)

		isSet, err = true, setWithProperType(field, val, opts)
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

func setWithProperType(field reflect.Value, val string, opts options) error {
	switch field.Kind() {
	case reflect.Int:
		return setIntField(field, stripNumValue(val, opts), 0)
	case reflect.Int8:
		return setIntField(field, stripNumValue(val, opts), 8)
	case reflect.Int16:
		return setIntField(field, stripNumValue(val, opts), 16)
	case reflect.Int32:
		return setIntField(field, stripNumValue(val, opts), 32)
	case reflect.Int64:
		switch field.Interface().(type) {
		case time.Duration:
			return setTimeDuration(field, val)
		}
		return setIntField(field, stripNumValue(val, opts), 64)
	case reflect.Uint:
		return setUintField(field, stripNumValue(val, opts), 0)
	case reflect.Uint8:
		return setUintField(field, stripNumValue(val, opts), 8)
	case reflect.Uint16:
		return setUintField(field, stripNumValue(val, opts), 16)
	case reflect.Uint32:
		return setUintField(field, stripNumValue(val, opts), 32)
	case reflect.Uint64:
		return setUintField(field, stripNumValue(val, opts), 64)
	case reflect.Bool:
		return setBoolField(field, stripNumValue(val, opts))
	case reflect.Float32:
		return setFloatField(field, stripNumValue(val, opts), 32)
	case reflect.Float64:
		return setFloatField(field, stripNumValue(val, opts), 64)
	case reflect.String:
		return setStringField(field, val)
	case reflect.Struct:
		switch field.Interface().(type) {
		case time.Time:
			return setTimeField(field, val, opts)
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
		return setWithProperType(field.Elem(), val, opts)
	default:
		return fmt.Errorf("form: unknown type %v", field.Kind())
	}
}

func stripNumValue(val string, opts options) string {
	if opts.strip {
		return str.RemoveAny(val, ",_")
	}
	return val
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

func setTimeField(field reflect.Value, val string, opts options) error {
	if val == "" {
		field.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	tf := strings.ToLower(opts.tag.Get("time_format"))

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
	if isUTC, _ := strconv.ParseBool(opts.tag.Get("time_utc")); isUTC {
		loc = time.UTC
	}

	if locTag := opts.tag.Get("time_location"); locTag != "" {
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

func setMap(field reflect.Value, form map[string][]string, key string, opts options) ([]string, bool, error) {
	vs, ok := getFormValues(form, key, opts)
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

	mk, me := mt.Key(), mt.Elem()

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

		// convert key
		kv := reflect.ValueOf(k)
		if kv.Type() != mk {
			cv, err := ref.CastTo(key, mk)
			if err != nil {
				return vs, false, fmt.Errorf("form: invalid map key type - %w", err)
			}
			kv = reflect.ValueOf(cv)
		}

		ps = trimFormValues(ps, opts)
		vs = append(vs, ps...)

		var val any
		switch len(ps) {
		case 0:
			val = ""
		case 1:
			val = ps[0]
		default:
			val = ps
		}

		// convert value
		vv := reflect.ValueOf(val)
		if vv.Type() != me {
			cv, err := ref.CastTo(val, me)
			if err != nil {
				return vs, false, fmt.Errorf("form: invalid map value type - %w", err)
			}
			vv = reflect.ValueOf(cv)
		}

		// set map
		field.SetMapIndex(kv, vv)
		ok = true
	}

	return vs, ok, nil
}

func setSlice(field reflect.Value, vals []string, opts options) error {
	slice := reflect.MakeSlice(field.Type(), len(vals), len(vals))
	if err := setArray(slice, vals, opts); err != nil {
		return err
	}
	field.Set(slice)
	return nil
}

func setArray(field reflect.Value, vals []string, opts options) error {
	for i, s := range vals {
		if err := setWithProperType(field.Index(i), s, opts); err != nil {
			return err
		}
	}
	return nil
}

func setTimeDuration(field reflect.Value, val string) error {
	d, err := tmu.ParseDuration(val)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(d))
	return nil
}

func setFormMap(dict any, form map[string][]string) error {
	switch m := dict.(type) {
	case map[string][]string:
		mag.Copy(m, form)
		return nil
	case map[string]string:
		for k, vs := range form {
			m[k] = asg.First(vs) // pick first
		}
		return nil
	case map[string]any:
		for k, vs := range form {
			switch len(vs) {
			case 0:
				m[k] = ""
			case 1:
				m[k] = vs[0]
			default:
				m[k] = vs
			}
		}
		return nil
	default:
		me := reflect.TypeOf(dict).Elem()
		if me.Kind() == reflect.Slice {
			// element is slice
			for k, vs := range form {
				if _, err := ref.MapSet(dict, k, vs); err != nil {
					return err
				}
			}
			return nil
		}

		// element is not slice
		for k, vs := range form {
			if _, err := ref.MapSet(dict, k, asg.First(vs)); err != nil {
				return err
			}
		}
		return nil
	}
}
