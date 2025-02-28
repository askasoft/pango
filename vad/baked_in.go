package vad

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/askasoft/pango/str"
)

// Func accepts a FieldLevel interface for all validation needs. The return
// value should be true when validation succeeds.
type Func func(fl FieldLevel) bool

// FuncCtx accepts a context.Context and FieldLevel interface for all
// validation needs. The return value should be true when validation succeeds.
type FuncCtx func(ctx context.Context, fl FieldLevel) bool

// wrapFunc wraps noramal Func makes it compatible with FuncCtx
func wrapFunc(fn Func) FuncCtx {
	if fn == nil {
		return nil // be sure not to wrap a bad function.
	}
	return func(ctx context.Context, fl FieldLevel) bool {
		return fn(fl)
	}
}

var (
	restrictedTags = map[string]struct{}{
		diveTag:           {},
		keysTag:           {},
		endKeysTag:        {},
		structOnlyTag:     {},
		omitempty:         {},
		skipValidationTag: {},
		utf8HexComma:      {},
		utf8Pipe:          {},
		noStructLevelTag:  {},
		requiredTag:       {},
		isdefault:         {},
	}

	// bakedInAliases is a default mapping of a single validation tag that
	// defines a common or complex set of validation(s) to simplify
	// adding validation to structs.
	bakedInAliases = map[string]string{
		"alpha":           "letter",
		"alphanum":        "letternum",
		"alphaunicode":    "utfletter",
		"alphanumunicode": "utfletternum",
		"bic":             "swiftCode",
		"iscolor":         "hexcolor|rgb|rgba|hsl|hsla",
	}

	// bakedInValidators is the default map of ValidationFunc
	// you can add, remove or even replace items to suite your needs,
	// or even disregard and use your own map if so desired.
	bakedInValidators = map[string]Func{
		"required":                      hasValue,
		"required_if":                   requiredIf,
		"required_unless":               requiredUnless,
		"required_with":                 requiredWith,
		"required_with_all":             requiredWithAll,
		"required_without":              requiredWithout,
		"required_without_all":          requiredWithoutAll,
		"excluded_with":                 excludedWith,
		"excluded_with_all":             excludedWithAll,
		"excluded_without":              excludedWithout,
		"excluded_without_all":          excludedWithoutAll,
		"isdefault":                     isDefault,
		"len":                           isLen,
		"minlen":                        isMinLen,
		"maxlen":                        isMaxLen,
		"btwlen":                        isBtwLen,
		"min":                           isGte,
		"max":                           isLte,
		"eq":                            isEq,
		"ne":                            isNe,
		"lt":                            isLt,
		"lte":                           isLte,
		"gt":                            isGt,
		"gte":                           isGte,
		"btw":                           isBtw,
		"eqfield":                       isEqField,
		"nefield":                       isNeField,
		"gtfield":                       isGtField,
		"ltfield":                       isLtField,
		"gtefield":                      isGteField,
		"ltefield":                      isLteField,
		"eqcsfield":                     isEqCrossStructField,
		"necsfield":                     isNeCrossStructField,
		"gtcsfield":                     isGtCrossStructField,
		"ltcsfield":                     isLtCrossStructField,
		"gtecsfield":                    isGteCrossStructField,
		"ltecsfield":                    isLteCrossStructField,
		"fieldcontains":                 fieldContains,
		"fieldexcludes":                 fieldExcludes,
		"letter":                        isLetter,
		"letternum":                     isLetterNumber,
		"utfletter":                     isUTFLetter,
		"utfletternum":                  isUTFLetterNumber,
		"boolean":                       isBoolean,
		"decimal":                       isDecimal,
		"numeric":                       isNumeric,
		"number":                        isNumber,
		"hexadecimal":                   isHexadecimal,
		"hexcolor":                      isHexColor,
		"rgb":                           isRGB,
		"rgba":                          isRGBA,
		"hsl":                           isHSL,
		"hsla":                          isHSLA,
		"e164":                          isE164,
		"email":                         isEmail,
		"uri":                           isURI,
		"url":                           isURL,
		"httpurl":                       isHttpURL,
		"httpsurl":                      isHttpsURL,
		"httpxurl":                      isHttpxURL,
		"base64":                        isBase64,
		"base64url":                     isBase64URL,
		"regexp":                        isRegexp,
		"contains":                      contains,
		"containsany":                   containsAny,
		"excludes":                      excludes,
		"excludesall":                   excludesAll,
		"startswith":                    startsWith,
		"endswith":                      endsWith,
		"startsnotwith":                 startsNotWith,
		"endsnotwith":                   endsNotWith,
		"isbn":                          isISBN,
		"isbn10":                        isISBN10,
		"isbn13":                        isISBN13,
		"btc_addr":                      isBitcoinAddress,
		"btc_addr_bech32":               isBitcoinBech32Address,
		"uuid":                          isUUID,
		"uuid3":                         isUUIDv3,
		"uuid4":                         isUUIDv4,
		"uuid5":                         isUUIDv5,
		"ulid":                          isULID,
		"ascii":                         isASCII,
		"printable":                     isUTFPrintable,
		"printascii":                    isASCIIPrintable,
		"multibyte":                     hasMultibyte,
		"datauri":                       isDataURI,
		"latitude":                      isLatitude,
		"longitude":                     isLongitude,
		"ssn":                           isSSN,
		"ipv4":                          isIPv4,
		"ipv6":                          isIPv6,
		"ip":                            isIP,
		"cidrv4":                        isCIDRv4,
		"cidrv6":                        isCIDRv6,
		"cidr":                          isCIDR,
		"mac":                           isMAC,
		"hostname":                      isHostnameRFC952,  // RFC 952
		"hostname_rfc1123":              isHostnameRFC1123, // RFC 1123
		"fqdn":                          isFQDN,
		"unique":                        isUnique,
		"oneof":                         isOneOf,
		"json":                          isJSON,
		"jsonobject":                    isJSONObject,
		"jsonarray":                     isJSONArray,
		"jwt":                           isJWT,
		"hostname_port":                 isHostnamePort,
		"lowercase":                     isLowercase,
		"uppercase":                     isUppercase,
		"datetime":                      isDatetime,
		"duration":                      isDuration,
		"timezone":                      isTimeZone,
		"postcode_iso3166_alpha2":       isPostcodeByIso3166Alpha2,
		"postcode_iso3166_alpha2_field": isPostcodeByIso3166Alpha2Field,
		"swiftcode":                     isSwiftCode,
		"semver":                        isSemverFormat,
		"dns_rfc1035_label":             isDnsRFC1035LabelFormat,
	}
)

func isOneOf(fl FieldLevel) bool {
	vs := splits(fl.Param())

	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	default:
		panic(fmt.Sprintf("oneof: bad field type %T", field.Interface()))
	}
	for _, o := range vs {
		if o == v {
			return true
		}
	}
	return false
}

// isUnique is the validation function for validating if each array|slice|map value is unique
func isUnique(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()
	v := reflect.ValueOf(struct{}{})

	switch field.Kind() {
	case reflect.Slice, reflect.Array:
		elem := field.Type().Elem()
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		if param == "" {
			m := reflect.MakeMap(reflect.MapOf(elem, v.Type()))

			for i := 0; i < field.Len(); i++ {
				m.SetMapIndex(reflect.Indirect(field.Index(i)), v)
			}
			return field.Len() == m.Len()
		}

		sf, ok := elem.FieldByName(param)
		if !ok {
			panic(fmt.Sprintf("unique: bad field name %s", param))
		}

		sfTyp := sf.Type
		if sfTyp.Kind() == reflect.Ptr {
			sfTyp = sfTyp.Elem()
		}

		m := reflect.MakeMap(reflect.MapOf(sfTyp, v.Type()))
		for i := 0; i < field.Len(); i++ {
			m.SetMapIndex(reflect.Indirect(reflect.Indirect(field.Index(i)).FieldByName(param)), v)
		}
		return field.Len() == m.Len()
	case reflect.Map:
		m := reflect.MakeMap(reflect.MapOf(field.Type().Elem(), v.Type()))

		for _, k := range field.MapKeys() {
			m.SetMapIndex(field.MapIndex(k), v)
		}
		return field.Len() == m.Len()
	default:
		panic(fmt.Sprintf("unique: bad field type %T", field.Interface()))
	}
}

func isStringField(name string, fl FieldLevel) {
	field := fl.Field()
	if field.Kind() != reflect.String {
		panic(fmt.Sprintf("%s: bad field type %T", name, field.Interface()))
	}
}

// isMAC is the validation function for validating if the field's value is a valid MAC address.
func isMAC(fl FieldLevel) bool {
	isStringField("mac", fl)

	return IsMAC(fl.Field().String())
}

// isCIDRv4 is the validation function for validating if the field's value is a valid v4 CIDR address.
func isCIDRv4(fl FieldLevel) bool {
	isStringField("cidrv4", fl)

	return IsCIDRv4(fl.Field().String())
}

// isCIDRv6 is the validation function for validating if the field's value is a valid v6 CIDR address.
func isCIDRv6(fl FieldLevel) bool {
	isStringField("cidrv6", fl)

	return IsCIDRv6(fl.Field().String())
}

// isCIDR is the validation function for validating if the field's value is a valid v4 or v6 CIDR address.
func isCIDR(fl FieldLevel) bool {
	isStringField("cidr", fl)

	return IsCIDR(fl.Field().String())
}

// isIPv4 is the validation function for validating if a value is a valid v4 IP address.
func isIPv4(fl FieldLevel) bool {
	isStringField("ipv4", fl)

	return IsIPv4(fl.Field().String())
}

// isIPv6 is the validation function for validating if the field's value is a valid v6 IP address.
func isIPv6(fl FieldLevel) bool {
	isStringField("ipv6", fl)

	return IsIPv6(fl.Field().String())
}

// isIP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
func isIP(fl FieldLevel) bool {
	isStringField("ip", fl)

	return IsIP(fl.Field().String())
}

// isSSN is the validation function for validating if the field's value is a valid SSN.
func isSSN(fl FieldLevel) bool {
	isStringField("ssn", fl)

	return IsSSN(fl.Field().String())
}

// isLongitude is the validation function for validating if the field's value is a valid longitude coordinate.
func isLongitude(fl FieldLevel) bool {
	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 64)
	default:
		panic(fmt.Sprintf("longitude: bad field type %T", field.Interface()))
	}

	return IsLongitude(v)
}

// isLatitude is the validation function for validating if the field's value is a valid latitude coordinate.
func isLatitude(fl FieldLevel) bool {
	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 64)
	default:
		panic(fmt.Sprintf("latitude: bad field type %T", field.Interface()))
	}

	return IsLatitude(v)
}

// isDataURI is the validation function for validating if the field's value is a valid data URI.
func isDataURI(fl FieldLevel) bool {
	isStringField("datauri", fl)

	return IsDataURI(fl.Field().String())
}

// hasMultibyte is the validation function for validating if the field's value has a multi byte character.
func hasMultibyte(fl FieldLevel) bool {
	isStringField("multibyte", fl)

	return str.HasMultibyte(fl.Field().String())
}

// isPrintable is the validation function for validating if the field's value is a valid printable character.
func isUTFPrintable(fl FieldLevel) bool {
	isStringField("printable", fl)

	return str.IsUTFPrintable(fl.Field().String())
}

// isASCIIPrintable is the validation function for validating if the field's value is a valid printable ASCII character.
func isASCIIPrintable(fl FieldLevel) bool {
	isStringField("printascii", fl)

	return str.IsASCIIPrintable(fl.Field().String())
}

// isASCII is the validation function for validating if the field's value is a valid ASCII character.
func isASCII(fl FieldLevel) bool {
	isStringField("ascii", fl)

	return str.IsASCII(fl.Field().String())
}

// isUUIDv5 is the validation function for validating if the field's value is a valid v5 UUID.
func isUUIDv5(fl FieldLevel) bool {
	isStringField("uuid5", fl)

	return IsUUIDv5(fl.Field().String())
}

// isUUIDv4 is the validation function for validating if the field's value is a valid v4 UUID.
func isUUIDv4(fl FieldLevel) bool {
	isStringField("uuid4", fl)

	return IsUUIDv4(fl.Field().String())
}

// isUUIDv3 is the validation function for validating if the field's value is a valid v3 UUID.
func isUUIDv3(fl FieldLevel) bool {
	isStringField("uuid3", fl)

	return IsUUIDv3(fl.Field().String())
}

// isUUID is the validation function for validating if the field's value is a valid UUID of any version.
func isUUID(fl FieldLevel) bool {
	isStringField("uuid", fl)

	return IsUUID(fl.Field().String())
}

// isULID is the validation function for validating if the field's value is a valid ULID.
func isULID(fl FieldLevel) bool {
	isStringField("ulid", fl)

	return uLIDRegex.MatchString(fl.Field().String())
}

// isISBN is the validation function for validating if the field's value is a valid v10 or v13 ISBN.
func isISBN(fl FieldLevel) bool {
	isStringField("isbn", fl)

	return isISBN10(fl) || isISBN13(fl)
}

// isISBN13 is the validation function for validating if the field's value is a valid v13 ISBN.
func isISBN13(fl FieldLevel) bool {
	isStringField("isbn13", fl)

	return IsISBN13(fl.Field().String())
}

// isISBN10 is the validation function for validating if the field's value is a valid v10 ISBN.
func isISBN10(fl FieldLevel) bool {
	isStringField("isbn10", fl)

	return IsISBN10(fl.Field().String())
}

// isBitcoinAddress is the validation function for validating if the field's value is a valid btc address
func isBitcoinAddress(fl FieldLevel) bool {
	isStringField("btc_addr", fl)

	return IsBitcoinAddress(fl.Field().String())
}

// isBitcoinBech32Address is the validation function for validating if the field's value is a valid bech32 btc address
func isBitcoinBech32Address(fl FieldLevel) bool {
	isStringField("btc_addr_bech32", fl)

	return IsBitcoinBech32Address(fl.Field().String())
}

// excludesAll is the validation function for validating that the field's value does not contain any of the characters specified within the param.
func excludesAll(fl FieldLevel) bool {
	isStringField("excludesall", fl)

	return !containsAny(fl)
}

// excludes is the validation function for validating that the field's value does not contain the text specified within the param.
func excludes(fl FieldLevel) bool {
	isStringField("excludes", fl)

	return !contains(fl)
}

// containsAny is the validation function for validating that the field's value contains any of the characters specified within the param.
func containsAny(fl FieldLevel) bool {
	isStringField("containsany", fl)

	return strings.ContainsAny(fl.Field().String(), fl.Param())
}

// contains is the validation function for validating that the field's value contains the text specified within the param.
func contains(fl FieldLevel) bool {
	isStringField("contains", fl)

	return strings.Contains(fl.Field().String(), fl.Param())
}

// isRegexp is the validation function for validating that the field's value match the regular expression specified within the param.
func isRegexp(fl FieldLevel) bool {
	isStringField("regexp", fl)

	re := regexp.MustCompile(fl.Param())
	return re.MatchString(fl.Field().String())
}

// startsWith is the validation function for validating that the field's value starts with the text specified within the param.
func startsWith(fl FieldLevel) bool {
	isStringField("startswith", fl)

	return strings.HasPrefix(fl.Field().String(), fl.Param())
}

// endsWith is the validation function for validating that the field's value ends with the text specified within the param.
func endsWith(fl FieldLevel) bool {
	isStringField("endswith", fl)

	return strings.HasSuffix(fl.Field().String(), fl.Param())
}

// startsNotWith is the validation function for validating that the field's value does not start with the text specified within the param.
func startsNotWith(fl FieldLevel) bool {
	isStringField("startsnotwith", fl)

	return !startsWith(fl)
}

// endsNotWith is the validation function for validating that the field's value does not end with the text specified within the param.
func endsNotWith(fl FieldLevel) bool {
	isStringField("endsnotwith", fl)

	return !endsWith(fl)
}

// fieldContains is the validation function for validating if the current field's value contains the field specified by the param's value.
func fieldContains(fl FieldLevel) bool {
	isStringField("fieldcontains", fl)

	currentField, _, ok := fl.GetStructFieldOK()
	if !ok {
		return false
	}

	return strings.Contains(fl.Field().String(), currentField.String())
}

// fieldExcludes is the validation function for validating if the current field's value excludes the field specified by the param's value.
func fieldExcludes(fl FieldLevel) bool {
	isStringField("fieldexcludes", fl)

	currentField, _, ok := fl.GetStructFieldOK()
	if !ok {
		return true
	}

	return !strings.Contains(fl.Field().String(), currentField.String())
}

// isNeField is the validation function for validating if the current field's value is not equal to the field specified by the param's value.
func isNeField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()

	if !ok || currentKind != kind {
		return true
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() != currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() != currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() != currentField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) != int64(currentField.Len())

	case reflect.Bool:
		return field.Bool() != currentField.Bool()

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return true
		}

		if fieldType == timeType {
			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)
			return !fieldTime.Equal(t)
		}
	}

	// default reflect.String:
	return field.String() != currentField.String()
}

// isNe is the validation function for validating that the field's value does not equal the provided param value.
func isNe(fl FieldLevel) bool {
	return !isEq(fl)
}

// isLteCrossStructField is the validation function for validating if the current field's value is less than or equal to the field, within a separate struct, specified by the param's value.
func isLteCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() <= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() <= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() <= topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) <= int64(topField.Len())

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)
			return fieldTime.Before(topTime) || fieldTime.Equal(topTime)
		}
	}

	// default reflect.String:
	return field.String() <= topField.String()
}

// isLtCrossStructField is the validation function for validating if the current field's value is less than the field, within a separate struct, specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func isLtCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() < topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() < topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() < topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) < int64(topField.Len())

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)
			return fieldTime.Before(topTime)
		}
	}

	// default reflect.String:
	return field.String() < topField.String()
}

// isGteCrossStructField is the validation function for validating if the current field's value is greater than or equal to the field, within a separate struct, specified by the param's value.
func isGteCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() >= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() >= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() >= topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) >= int64(topField.Len())

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)
			return fieldTime.After(topTime) || fieldTime.Equal(topTime)
		}
	}

	// default reflect.String:
	return field.String() >= topField.String()
}

// isGtCrossStructField is the validation function for validating if the current field's value is greater than the field, within a separate struct, specified by the param's value.
func isGtCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() > topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() > topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() > topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) > int64(topField.Len())

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)
			return fieldTime.After(topTime)
		}
	}

	// default reflect.String:
	return field.String() > topField.String()
}

// isNeCrossStructField is the validation function for validating that the current field's value is not equal to the field, within a separate struct, specified by the param's value.
func isNeCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return true
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() != field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return topField.Uint() != field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() != field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) != int64(field.Len())

	case reflect.Bool:
		return topField.Bool() != field.Bool()

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return true
		}

		if fieldType == timeType {
			t := field.Interface().(time.Time)
			fieldTime := topField.Interface().(time.Time)
			return !fieldTime.Equal(t)
		}
	}

	// default reflect.String:
	return topField.String() != field.String()
}

// isEqCrossStructField is the validation function for validating that the current field's value is equal to the field, within a separate struct, specified by the param's value.
func isEqCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() == field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return topField.Uint() == field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() == field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) == int64(field.Len())

	case reflect.Bool:
		return topField.Bool() == field.Bool()

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			t := field.Interface().(time.Time)
			fieldTime := topField.Interface().(time.Time)
			return fieldTime.Equal(t)
		}
	}

	// default reflect.String:
	return topField.String() == field.String()
}

// isEqField is the validation function for validating if the current field's value is equal to the field specified by the param's value.
func isEqField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() == currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() == currentField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) == int64(currentField.Len())

	case reflect.Bool:
		return field.Bool() == currentField.Bool()

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {
			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)
			return fieldTime.Equal(t)
		}
	}

	// default reflect.String:
	return field.String() == currentField.String()
}

// isEq is the validation function for validating if the current field's value is equal to the param's value.
func isEq(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		return field.String() == param

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)
		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		return field.Float() == p

	case reflect.Bool:
		p := asBool(param)
		return field.Bool() == p

	case reflect.Struct:
		if field.Type() == timeType {
			p := asTime(param)
			t := field.Interface().(time.Time)
			return t.Equal(p)
		}
	}

	panic(fmt.Sprintf("eq: bad field type %T", field.Interface()))
}

// isPostcodeByIso3166Alpha2 validates by value which is country code in iso 3166 alpha 2
// example: `postcode_iso3166_alpha2=US`
func isPostcodeByIso3166Alpha2(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	reg, found := postCodeRegexDict[param]
	if !found {
		return false
	}

	return reg.MatchString(field.String())
}

// isPostcodeByIso3166Alpha2 validates by field which represents for a value of country code in iso 3166 alpha 2
// example: `postcode_iso3166_alpha2_field=CountryCode`
func isPostcodeByIso3166Alpha2Field(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	if param == "" {
		return false
	}
	currentField, kind, _, found := fl.GetStructFieldOKAdvanced2(fl.Parent(), param)
	if !found {
		return false
	}

	if kind != reflect.String {
		panic(fmt.Sprintf("postcode_iso3166_alpha2_field: bad field type %T", currentField.Interface()))
	}

	reg, found := postCodeRegexDict[currentField.String()]
	if !found {
		return false
	}

	return reg.MatchString(field.String())
}

// isBase64 is the validation function for validating if the current field's value is a valid base 64.
func isBase64(fl FieldLevel) bool {
	isStringField("base64", fl)

	return IsBase64(fl.Field().String())
}

// isBase64URL is the validation function for validating if the current field's value is a valid base64 URL safe string.
func isBase64URL(fl FieldLevel) bool {
	isStringField("base64url", fl)

	return IsBase64URL(fl.Field().String())
}

// isURI is the validation function for validating if the current field's value is a valid URI.
func isURI(fl FieldLevel) bool {
	isStringField("uri", fl)

	return IsURI(fl.Field().String())
}

// isURL is the validation function for validating if the current field's value is a valid URL.
func isURL(fl FieldLevel) bool {
	isStringField("url", fl)

	return IsURL(fl.Field().String())
}

// isHttpURL is the validation function for validating if the current field's value is a valid http:// URL.
func isHttpURL(fl FieldLevel) bool {
	isStringField("httpurl", fl)

	return IsHttpURL(fl.Field().String())
}

// isHttpsURL is the validation function for validating if the current field's value is a valid https:// URL.
func isHttpsURL(fl FieldLevel) bool {
	isStringField("httpsurl", fl)

	return IsHttpsURL(fl.Field().String())
}

// isHttpxURL is the validation function for validating if the current field's value is a valid https?:// URL.
func isHttpxURL(fl FieldLevel) bool {
	isStringField("httpxurl", fl)

	return IsHttpxURL(fl.Field().String())
}

// isE164 is the validation function for validating if the current field's value is a valid e.164 formatted phone number.
func isE164(fl FieldLevel) bool {
	isStringField("e164", fl)

	return e164Regex.MatchString(fl.Field().String())
}

// isEmail is the validation function for validating if the current field's value is a valid email address.
func isEmail(fl FieldLevel) bool {
	isStringField("email", fl)

	return IsEmail(fl.Field().String())
}

// isHSLA is the validation function for validating if the current field's value is a valid HSLA color.
func isHSLA(fl FieldLevel) bool {
	isStringField("hsla", fl)

	return IsHSLAColor(fl.Field().String())
}

// isHSL is the validation function for validating if the current field's value is a valid HSL color.
func isHSL(fl FieldLevel) bool {
	isStringField("hsl", fl)

	return IsHSLColor(fl.Field().String())
}

// isRGBA is the validation function for validating if the current field's value is a valid RGBA color.
func isRGBA(fl FieldLevel) bool {
	isStringField("rgba", fl)

	return IsRGBAColor(fl.Field().String())
}

// isRGB is the validation function for validating if the current field's value is a valid RGB color.
func isRGB(fl FieldLevel) bool {
	isStringField("rgb", fl)

	return IsRGBColor(fl.Field().String())
}

// isHexColor is the validation function for validating if the current field's value is a valid HEX color.
func isHexColor(fl FieldLevel) bool {
	isStringField("hexcolor", fl)

	return IsHexColor(fl.Field().String())
}

// isHexadecimal is the validation function for validating if the current field's value is a valid hexadecimal.
func isHexadecimal(fl FieldLevel) bool {
	isStringField("hexadecimal", fl)

	return str.IsHexadecimal(fl.Field().String())
}

// isNumber is the validation function for validating if the current field's value is a valid number.
func isNumber(fl FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return true
	default:
		return str.IsNumber(fl.Field().String())
	}
}

// isNumeric is the validation function for validating if the current field's value is a valid numeric value.
func isNumeric(fl FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return true
	default:
		return str.IsNumeric(fl.Field().String())
	}
}

// isDecimal is the validation function for validating if the current field's value is a valid decimal value.
func isDecimal(fl FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return true
	default:
		return str.IsDecimal(fl.Field().String())
	}
}

// isLetter is the validation function for validating if the current field's value is a valid letter value.
func isLetter(fl FieldLevel) bool {
	isStringField("letter", fl)

	return str.IsLetter(fl.Field().String())
}

// isLetterNumber is the validation function for validating if the current field's value is a valid alphanumeric value.
func isLetterNumber(fl FieldLevel) bool {
	isStringField("letternumber", fl)

	return str.IsLetterNumber(fl.Field().String())
}

// isUTFLetter is the validation function for validating if the current field's value is a valid alpha unicode value.
func isUTFLetter(fl FieldLevel) bool {
	isStringField("utfletter", fl)

	return str.IsUTFLetter(fl.Field().String())
}

// isUTFLetterNumber is the validation function for validating if the current field's value is a valid alphanumeric unicode value.
func isUTFLetterNumber(fl FieldLevel) bool {
	isStringField("utfletternumber", fl)

	return str.IsUTFLetterNumber(fl.Field().String())
}

// isBoolean is the validation function for validating if the current field's value can be safely converted to a boolean.
func isBoolean(fl FieldLevel) bool {
	_, err := strconv.ParseBool(fl.Field().String())
	return err == nil
}

// isDefault is the opposite of required aka hasValue
func isDefault(fl FieldLevel) bool {
	return !hasValue(fl)
}

// isEmpty is the validation function for validating if the current field's value is not the default static value.
// check recursively if the field is a pointer.
func isEmpty(fl FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return field.IsNil()
	default:
		return !field.IsValid() || field.Interface() == reflect.Zero(field.Type()).Interface()
	}
}

// hasValue is the validation function for validating if the current field's value is not the default static value.
func hasValue(fl FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		if fl.(*validate).fldIsPointer && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

// requireCheckField is a func for check field kind
func requireCheckFieldKind(fl FieldLevel, param string, defaultNotFoundValue bool) bool {
	field := fl.Field()
	kind := field.Kind()
	var nullable, found bool
	if len(param) > 0 {
		field, kind, nullable, found = fl.GetStructFieldOKAdvanced2(fl.Parent(), param)
		if !found {
			return defaultNotFoundValue
		}
	}
	switch kind {
	case reflect.Invalid:
		return defaultNotFoundValue
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return field.IsNil()
	default:
		if nullable && field.Interface() != nil {
			return false
		}
		return field.IsValid() && field.Interface() == reflect.Zero(field.Type()).Interface()
	}
}

// requireCheckFieldValue is a func for check field value
func requireCheckFieldValue(fl FieldLevel, param string, value string, defaultNotFoundValue bool) bool {
	field, kind, _, found := fl.GetStructFieldOKAdvanced2(fl.Parent(), param)
	if !found {
		return defaultNotFoundValue
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == asInt(value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() == asUint(value)

	case reflect.Float32, reflect.Float64:
		return field.Float() == asFloat(value)

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) == asInt(value)

	case reflect.Bool:
		return field.Bool() == asBool(value)
	}

	// default reflect.String:
	return field.String() == value
}

// requiredIf is the validation function
// The field under validation must be present and not empty only if all the other specified fields are equal to the value following with the specified field.
func requiredIf(fl FieldLevel) bool {
	params := splits(fl.Param())
	if len(params)%2 != 0 {
		panic(fmt.Sprintf("Bad param number for required_if %s", fl.FieldName()))
	}
	for i := 0; i < len(params); i += 2 {
		if !requireCheckFieldValue(fl, params[i], params[i+1], false) {
			return true
		}
	}
	return hasValue(fl)
}

// requiredUnless is the validation function
// The field under validation must be present and not empty only unless all the other specified fields are equal to the value following with the specified field.
func requiredUnless(fl FieldLevel) bool {
	params := splits(fl.Param())
	if len(params)%2 != 0 {
		panic(fmt.Sprintf("Bad param number for required_unless %s", fl.FieldName()))
	}

	for i := 0; i < len(params); i += 2 {
		if requireCheckFieldValue(fl, params[i], params[i+1], false) {
			return true
		}
	}
	return hasValue(fl)
}

// excludedWith is the validation function
// The field under validation must not be present or is empty if any of the other specified fields are present.
func excludedWith(fl FieldLevel) bool {
	params := splits(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return !hasValue(fl)
		}
	}
	return true
}

// requiredWith is the validation function
// The field under validation must be present and not empty only if any of the other specified fields are present.
func requiredWith(fl FieldLevel) bool {
	params := splits(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return hasValue(fl)
		}
	}
	return true
}

// excludedWithAll is the validation function
// The field under validation must not be present or is empty if all of the other specified fields are present.
func excludedWithAll(fl FieldLevel) bool {
	params := splits(fl.Param())
	for _, param := range params {
		if requireCheckFieldKind(fl, param, true) {
			return true
		}
	}
	return !hasValue(fl)
}

// requiredWithAll is the validation function
// The field under validation must be present and not empty only if all of the other specified fields are present.
func requiredWithAll(fl FieldLevel) bool {
	params := splits(fl.Param())
	for _, param := range params {
		if requireCheckFieldKind(fl, param, true) {
			return true
		}
	}
	return hasValue(fl)
}

// excludedWithout is the validation function
// The field under validation must not be present or is empty when any of the other specified fields are not present.
func excludedWithout(fl FieldLevel) bool {
	if requireCheckFieldKind(fl, strings.TrimSpace(fl.Param()), true) {
		return !hasValue(fl)
	}
	return true
}

// requiredWithout is the validation function
// The field under validation must be present and not empty only when any of the other specified fields are not present.
func requiredWithout(fl FieldLevel) bool {
	if requireCheckFieldKind(fl, strings.TrimSpace(fl.Param()), true) {
		return hasValue(fl)
	}
	return true
}

// excludedWithoutAll is the validation function
// The field under validation must not be present or is empty when all of the other specified fields are not present.
func excludedWithoutAll(fl FieldLevel) bool {
	params := splits(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return true
		}
	}
	return !hasValue(fl)
}

// requiredWithoutAll is the validation function
// The field under validation must be present and not empty only when all of the other specified fields are not present.
func requiredWithoutAll(fl FieldLevel) bool {
	params := splits(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return true
		}
	}
	return hasValue(fl)
}

// isLen is the validation function for validating if the current field's length or rune count is equal to the param's value.
func isLen(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		p := asInt(param)
		return int64(utf8.RuneCountInString(field.String())) == p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		return int64(field.Len()) == p
	}

	panic(fmt.Sprintf("len: bad field type %T", field.Interface()))
}

// isMaxLen is the validation function for validating if the current field's length or rune count is less than or equal to the param's value.
func isMaxLen(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		p := asInt(param)
		return int64(utf8.RuneCountInString(field.String())) <= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		return int64(field.Len()) <= p
	}

	panic(fmt.Sprintf("maxlen: bad field type %T", field.Interface()))
}

// isMinLen is the validation function for validating if the current field's length or rune count is greater than or equal to the param's value.
func isMinLen(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		p := asInt(param)
		return int64(utf8.RuneCountInString(field.String())) >= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		return int64(field.Len()) >= p
	}

	panic(fmt.Sprintf("minlen: bad field type %T", field.Interface()))
}

// isBtwLen is the validation function for validating if the current field's length or rune count is between the param's value.
func isBtwLen(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		p1, p2 := asInt2(param)
		l := int64(utf8.RuneCountInString(field.String()))
		return l >= p1 && l <= p2

	case reflect.Slice, reflect.Map, reflect.Array:
		p1, p2 := asInt2(param)
		l := int64(field.Len())
		return l >= p1 && l <= p2
	}

	panic(fmt.Sprintf("btwlen: bad field type %T", field.Interface()))
}

// isGteField is the validation function for validating if the current field's value is greater than or equal to the field specified by the param's value.
func isGteField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() >= currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() >= currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() >= currentField.Float()

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {
			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)
			return fieldTime.After(t) || fieldTime.Equal(t)
		}
	}

	// default reflect.String
	return len(field.String()) >= len(currentField.String())
}

// isGtField is the validation function for validating if the current field's value is greater than the field specified by the param's value.
func isGtField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() > currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() > currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() > currentField.Float()

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {
			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)
			return fieldTime.After(t)
		}
	}

	// default reflect.String
	return len(field.String()) > len(currentField.String())
}

// isGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isGte(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		return field.String() >= param

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)
		return field.Int() >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		return field.Uint() >= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		return field.Float() >= p

	case reflect.Struct:
		if field.Type() == timeType {
			p := asTime(param)
			t := field.Interface().(time.Time)
			return t.After(p) || t.Equal(p)
		}
	}

	panic(fmt.Sprintf("gte: bad field type %T", field.Interface()))
}

// isGt is the validation function for validating if the current field's value is greater than the param's value.
func isGt(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		return field.String() > param

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)
		return field.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		return field.Uint() > p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		return field.Float() > p

	case reflect.Struct:
		if field.Type() == timeType {
			p := asTime(param)
			t := field.Interface().(time.Time)
			return t.After(p)
		}
	}

	panic(fmt.Sprintf("gt: bad field type %T", field.Interface()))
}

// isLteField is the validation function for validating if the current field's value is less than or equal to the field specified by the param's value.
func isLteField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() <= currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() <= currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() <= currentField.Float()

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {
			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)
			return fieldTime.Before(t) || fieldTime.Equal(t)
		}
	}

	// default reflect.String
	return len(field.String()) <= len(currentField.String())
}

// isLtField is the validation function for validating if the current field's value is less than the field specified by the param's value.
func isLtField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() < currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() < currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() < currentField.Float()

	case reflect.Struct:
		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {
			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)
			return fieldTime.Before(t)
		}
	}

	// default reflect.String
	return len(field.String()) < len(currentField.String())
}

// isLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func isLte(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		return field.String() <= param

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)
		return field.Int() <= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		return field.Uint() <= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		return field.Float() <= p

	case reflect.Struct:
		if field.Type() == timeType {
			p := asTime(param)
			t := field.Interface().(time.Time)
			return t.Before(p) || t.Equal(p)
		}
	}

	panic(fmt.Sprintf("lte: bad field type %T", field.Interface()))
}

// isLt is the validation function for validating if the current field's value is less than the param's value.
func isLt(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		return field.String() < param

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		return int64(field.Len()) < p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)
		return field.Int() < p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		return field.Uint() < p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		return field.Float() < p

	case reflect.Struct:
		if field.Type() == timeType {
			p := asTime(param)
			t := field.Interface().(time.Time)
			return t.Before(p)
		}
	}

	panic(fmt.Sprintf("lt: bad field type %T", field.Interface()))
}

// isBtw is the validation function for validating if the current field's value is between the param's value "min max".
func isBtw(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		p1, p2 := split2(param)
		return field.String() >= p1 && field.String() <= p2

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p1, p2 := asInt2FromType(field.Type(), param)
		i := field.Int()
		return i >= p1 && i <= p2

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p1, p2 := asUint2(param)
		u := field.Uint()
		return u >= p1 && u <= p2

	case reflect.Float32, reflect.Float64:
		p1, p2 := asFloat2(param)
		f := field.Float()
		return f >= p1 && f <= p2

	case reflect.Struct:
		if field.Type() == timeType {
			p1, p2 := asTime2(param)
			t := field.Interface().(time.Time)
			return !(p1.After(t) || p2.Before(t))
		}
	}

	panic(fmt.Sprintf("btw: bad field type %T", field.Interface()))
}

func isHostnameRFC952(fl FieldLevel) bool {
	isStringField("hostname", fl)

	return hostnameRegexRFC952.MatchString(fl.Field().String())
}

func isHostnameRFC1123(fl FieldLevel) bool {
	isStringField("hostname_rfc1123", fl)

	return hostnameRegexRFC1123.MatchString(fl.Field().String())
}

func isFQDN(fl FieldLevel) bool {
	isStringField("fqdn", fl)

	val := fl.Field().String()
	if val == "" {
		return false
	}

	return fqdnRegexRFC1123.MatchString(val)
}

// isJSON is the validation function for validating if the current field's value is a valid json string.
func isJSON(fl FieldLevel) bool {
	isStringField("json", fl)

	return json.Valid(str.UnsafeBytes(fl.Field().String()))
}

// isJSONObject is the validation function for validating if the current field's value is a valid json object string.
func isJSONObject(fl FieldLevel) bool {
	isStringField("jsonobject", fl)

	return IsJSONObject(fl.Field().String())
}

// isJSONArray is the validation function for validating if the current field's value is a valid json array string.
func isJSONArray(fl FieldLevel) bool {
	isStringField("jsonarray", fl)

	return IsJSONArray(fl.Field().String())
}

// isJWT is the validation function for validating if the current field's value is a valid JWT string.
func isJWT(fl FieldLevel) bool {
	isStringField("jwt", fl)

	return IsJWT(fl.Field().String())
}

// isHostnamePort validates a <dns>:<port> combination for fields typically used for socket address.
func isHostnamePort(fl FieldLevel) bool {
	isStringField("hostname_port", fl)

	val := fl.Field().String()
	host, port, err := net.SplitHostPort(val)
	if err != nil {
		return false
	}

	// Port must be a iny <= 65535.
	if portNum, err := strconv.ParseInt(port, 10, 32); err != nil || portNum > 65535 || portNum < 1 {
		return false
	}

	// If host is specified, it should match a DNS name
	if host != "" {
		return hostnameRegexRFC1123.MatchString(host)
	}
	return true
}

// isLowercase is the validation function for validating if the current field's value is a lowercase string.
func isLowercase(fl FieldLevel) bool {
	isStringField("lowercase", fl)

	return str.IsLowerCase(fl.Field().String())
}

// isUppercase is the validation function for validating if the current field's value is an uppercase string.
func isUppercase(fl FieldLevel) bool {
	isStringField("uppercase", fl)

	return str.IsUpperCase(fl.Field().String())
}

// isDatetime is the validation function for validating if the current field's value is a valid datetime string.
func isDatetime(fl FieldLevel) bool {
	isStringField("datetime", fl)

	field := fl.Field()
	param := fl.Param()

	_, err := time.Parse(param, field.String())
	return err == nil
}

// isDuration is the validation function for validating if the current field's value is a valid duration string.
func isDuration(fl FieldLevel) bool {
	isStringField("duration", fl)

	_, err := time.ParseDuration(fl.Field().String())
	return err == nil
}

// isTimeZone is the validation function for validating if the current field's value is a valid time zone string.
func isTimeZone(fl FieldLevel) bool {
	isStringField("timezone", fl)

	val := fl.Field().String()

	// empty value is converted to UTC by time.LoadLocation but disallow it as it is not a valid time zone name
	if val == "" {
		return false
	}

	// Local value is converted to the current system time zone by time.LoadLocation but disallow it as it is not a valid time zone name
	if strings.ToLower(val) == "local" {
		return false
	}

	_, err := time.LoadLocation(val)
	return err == nil
}

// isSwiftCode is the validation function for validating if the current field's value is a valid Business Identifier Code (SWIFT code), defined in ISO 9362
func isSwiftCode(fl FieldLevel) bool {
	isStringField("swiftcode", fl)

	return IsSwiftCode(fl.Field().String())
}

// isSemverFormat is the validation function for validating if the current field's value is a valid semver version, defined in Semantic Versioning 2.0.0
func isSemverFormat(fl FieldLevel) bool {
	isStringField("semver", fl)
	return semverRegex.MatchString(fl.Field().String())
}

// isDnsRFC1035LabelFormat is the validation function
// for validating if the current field's value is
// a valid dns RFC 1035 label, defined in RFC 1035.
func isDnsRFC1035LabelFormat(fl FieldLevel) bool {
	isStringField("dns_rfc1035_label", fl)

	return dnsRegexRFC1035Label.MatchString(fl.Field().String())
}
