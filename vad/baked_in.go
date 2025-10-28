package vad

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/sch"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/str/wildcard"
	"github.com/askasoft/pango/tmu"
)

// Func accepts a FieldLevel interface for all validation needs.
type Func func(fl FieldLevel) bool

// FuncEx accepts a FieldLevel interface for all validation needs.
type FuncEx func(fl FieldLevel) error

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
		isempty:           {},
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
		"required":                      isRequired,
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
		"isempty":                       isEmpty,
		"len":                           isLen,
		"minlen":                        isMinLen,
		"maxlen":                        isMaxLen,
		"btwlen":                        isBtwLen,
		"min":                           isMin,
		"max":                           isMax,
		"btw":                           isBtw,
		"eq":                            isEq,
		"ne":                            isNe,
		"lt":                            isLt,
		"lte":                           isLte,
		"gt":                            isGt,
		"gte":                           isGte,
		"eqfield":                       isEqField,
		"nefield":                       isNeField,
		"gtfield":                       isGtField,
		"ltfield":                       isLtField,
		"gtefield":                      isGteField,
		"ltefield":                      isLteField,
		"eqcsfield":                     isEqField,
		"necsfield":                     isNeField,
		"gtcsfield":                     isGtField,
		"ltcsfield":                     isLtField,
		"gtecsfield":                    isGteField,
		"ltecsfield":                    isLteField,
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
		"size":                          isSize,
		"minsize":                       isSizeGte,
		"maxsize":                       isSizeLte,
		"btwsize":                       isSizeBtw,
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
		"rematch":                       regexpMatch,
		"wcmatch":                       wildcardMatch,
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
		"cron":                          isCron,
		"periodic":                      isPeriodic,
		"timezone":                      isTimeZone,
		"postcode_iso3166_alpha2":       isPostcodeByIso3166Alpha2,
		"postcode_iso3166_alpha2_field": isPostcodeByIso3166Alpha2Field,
		"swiftcode":                     isSwiftCode,
		"semver":                        isSemver,
		"dns_rfc1035_label":             isDnsRFC1035Label,
	}
)

func wrapFunc(tag string, fn Func) FuncEx {
	if fn == nil {
		panic("function cannot be empty")
	}

	return func(fl FieldLevel) error {
		if fn(fl) {
			return nil
		}
		return errors.New(tag)
	}
}

func MustStringField(name string, fl FieldLevel) {
	field := fl.Field()
	if field.Kind() != reflect.String {
		panic(fmt.Errorf("%s: bad field type %T", name, field.Interface()))
	}
}

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
		panic(fmt.Errorf("oneof: bad field type %T", field.Interface()))
	}
	return asg.Contains(vs, v)
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
			for i := range field.Len() {
				m.SetMapIndex(reflect.Indirect(field.Index(i)), v)
			}
			return field.Len() == m.Len()
		}

		sf, ok := elem.FieldByName(param)
		if !ok {
			panic(fmt.Errorf("unique: bad field name %s", param))
		}

		sfTyp := sf.Type
		if sfTyp.Kind() == reflect.Ptr {
			sfTyp = sfTyp.Elem()
		}

		m := reflect.MakeMap(reflect.MapOf(sfTyp, v.Type()))
		for i := range field.Len() {
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
		panic(fmt.Errorf("unique: bad field type %T", field.Interface()))
	}
}

// isMAC is the validation function for validating if the field's value is a valid MAC address.
func isMAC(fl FieldLevel) bool {
	MustStringField("mac", fl)

	return IsMAC(fl.Field().String())
}

// isCIDRv4 is the validation function for validating if the field's value is a valid v4 CIDR address.
func isCIDRv4(fl FieldLevel) bool {
	MustStringField("cidrv4", fl)

	return IsCIDRv4(fl.Field().String())
}

// isCIDRv6 is the validation function for validating if the field's value is a valid v6 CIDR address.
func isCIDRv6(fl FieldLevel) bool {
	MustStringField("cidrv6", fl)

	return IsCIDRv6(fl.Field().String())
}

// isCIDR is the validation function for validating if the field's value is a valid v4 or v6 CIDR address.
func isCIDR(fl FieldLevel) bool {
	MustStringField("cidr", fl)

	return IsCIDR(fl.Field().String())
}

// isIPv4 is the validation function for validating if a value is a valid v4 IP address.
func isIPv4(fl FieldLevel) bool {
	MustStringField("ipv4", fl)

	return IsIPv4(fl.Field().String())
}

// isIPv6 is the validation function for validating if the field's value is a valid v6 IP address.
func isIPv6(fl FieldLevel) bool {
	MustStringField("ipv6", fl)

	return IsIPv6(fl.Field().String())
}

// isIP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
func isIP(fl FieldLevel) bool {
	MustStringField("ip", fl)

	return IsIP(fl.Field().String())
}

// isSSN is the validation function for validating if the field's value is a valid SSN.
func isSSN(fl FieldLevel) bool {
	MustStringField("ssn", fl)

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
		panic(fmt.Errorf("longitude: bad field type %T", field.Interface()))
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
		panic(fmt.Errorf("latitude: bad field type %T", field.Interface()))
	}

	return IsLatitude(v)
}

// isDataURI is the validation function for validating if the field's value is a valid data URI.
func isDataURI(fl FieldLevel) bool {
	MustStringField("datauri", fl)

	return IsDataURI(fl.Field().String())
}

// hasMultibyte is the validation function for validating if the field's value has a multi byte character.
func hasMultibyte(fl FieldLevel) bool {
	MustStringField("multibyte", fl)

	return str.HasMultibyte(fl.Field().String())
}

// isPrintable is the validation function for validating if the field's value is a valid printable character.
func isUTFPrintable(fl FieldLevel) bool {
	MustStringField("printable", fl)

	return str.IsUTFPrintable(fl.Field().String())
}

// isASCIIPrintable is the validation function for validating if the field's value is a valid printable ASCII character.
func isASCIIPrintable(fl FieldLevel) bool {
	MustStringField("printascii", fl)

	return str.IsASCIIPrintable(fl.Field().String())
}

// isASCII is the validation function for validating if the field's value is a valid ASCII character.
func isASCII(fl FieldLevel) bool {
	MustStringField("ascii", fl)

	return str.IsASCII(fl.Field().String())
}

// isUUIDv5 is the validation function for validating if the field's value is a valid v5 UUID.
func isUUIDv5(fl FieldLevel) bool {
	MustStringField("uuid5", fl)

	return IsUUIDv5(fl.Field().String())
}

// isUUIDv4 is the validation function for validating if the field's value is a valid v4 UUID.
func isUUIDv4(fl FieldLevel) bool {
	MustStringField("uuid4", fl)

	return IsUUIDv4(fl.Field().String())
}

// isUUIDv3 is the validation function for validating if the field's value is a valid v3 UUID.
func isUUIDv3(fl FieldLevel) bool {
	MustStringField("uuid3", fl)

	return IsUUIDv3(fl.Field().String())
}

// isUUID is the validation function for validating if the field's value is a valid UUID of any version.
func isUUID(fl FieldLevel) bool {
	MustStringField("uuid", fl)

	return IsUUID(fl.Field().String())
}

// isULID is the validation function for validating if the field's value is a valid ULID.
func isULID(fl FieldLevel) bool {
	MustStringField("ulid", fl)

	return IsULID(fl.Field().String())
}

// isISBN is the validation function for validating if the field's value is a valid v10 or v13 ISBN.
func isISBN(fl FieldLevel) bool {
	MustStringField("isbn", fl)

	return isISBN10(fl) || isISBN13(fl)
}

// isISBN13 is the validation function for validating if the field's value is a valid v13 ISBN.
func isISBN13(fl FieldLevel) bool {
	MustStringField("isbn13", fl)

	return IsISBN13(fl.Field().String())
}

// isISBN10 is the validation function for validating if the field's value is a valid v10 ISBN.
func isISBN10(fl FieldLevel) bool {
	MustStringField("isbn10", fl)

	return IsISBN10(fl.Field().String())
}

// isBitcoinAddress is the validation function for validating if the field's value is a valid btc address
func isBitcoinAddress(fl FieldLevel) bool {
	MustStringField("btc_addr", fl)

	return IsBitcoinAddress(fl.Field().String())
}

// isBitcoinBech32Address is the validation function for validating if the field's value is a valid bech32 btc address
func isBitcoinBech32Address(fl FieldLevel) bool {
	MustStringField("btc_addr_bech32", fl)

	return IsBitcoinBech32Address(fl.Field().String())
}

// excludesAll is the validation function for validating that the field's value does not contain any of the characters specified within the param.
func excludesAll(fl FieldLevel) bool {
	MustStringField("excludesall", fl)

	return !containsAny(fl)
}

// excludes is the validation function for validating that the field's value does not contain the text specified within the param.
func excludes(fl FieldLevel) bool {
	MustStringField("excludes", fl)

	return !contains(fl)
}

// containsAny is the validation function for validating that the field's value contains any of the characters specified within the param.
func containsAny(fl FieldLevel) bool {
	MustStringField("containsany", fl)

	return strings.ContainsAny(fl.Field().String(), fl.Param())
}

// contains is the validation function for validating that the field's value contains the text specified within the param.
func contains(fl FieldLevel) bool {
	MustStringField("contains", fl)

	return strings.Contains(fl.Field().String(), fl.Param())
}

// isRegexp is the validation function for validating that the field's value is a valid regular expression.
func isRegexp(fl FieldLevel) bool {
	MustStringField("regexp", fl)

	_, err := regexp.Compile(fl.Field().String())
	return err == nil
}

// regexpMatch is the validation function for validating that the field's value match the regular expression specified within the param.
func regexpMatch(fl FieldLevel) bool {
	MustStringField("rematch", fl)

	re := regexp.MustCompile(fl.Param())
	return re.MatchString(fl.Field().String())
}

// wildcardMatch is the validation function for validating that the field's value match the wildcard expression specified within the param.
func wildcardMatch(fl FieldLevel) bool {
	MustStringField("wcmatch", fl)

	return wildcard.Match(fl.Param(), fl.Field().String())
}

// startsWith is the validation function for validating that the field's value starts with the text specified within the param.
func startsWith(fl FieldLevel) bool {
	MustStringField("startswith", fl)

	return strings.HasPrefix(fl.Field().String(), fl.Param())
}

// endsWith is the validation function for validating that the field's value ends with the text specified within the param.
func endsWith(fl FieldLevel) bool {
	MustStringField("endswith", fl)

	return strings.HasSuffix(fl.Field().String(), fl.Param())
}

// startsNotWith is the validation function for validating that the field's value does not start with the text specified within the param.
func startsNotWith(fl FieldLevel) bool {
	MustStringField("startsnotwith", fl)

	return !startsWith(fl)
}

// endsNotWith is the validation function for validating that the field's value does not end with the text specified within the param.
func endsNotWith(fl FieldLevel) bool {
	MustStringField("endsnotwith", fl)

	return !endsWith(fl)
}

// fieldContains is the validation function for validating if the current field's value contains the field specified by the param's value.
func fieldContains(fl FieldLevel) bool {
	MustStringField("fieldcontains", fl)

	cfield, _, ok := fl.GetStructFieldOK()
	if !ok {
		return false
	}

	return strings.Contains(fl.Field().String(), cfield.String())
}

// fieldExcludes is the validation function for validating if the current field's value excludes the field specified by the param's value.
func fieldExcludes(fl FieldLevel) bool {
	MustStringField("fieldexcludes", fl)

	cfield, _, ok := fl.GetStructFieldOK()
	if !ok {
		return true
	}

	return !strings.Contains(fl.Field().String(), cfield.String())
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

	cfield, kind, _, found := fl.GetStructFieldOKAdvanced2(fl.Parent(), param)
	if !found {
		return false
	}

	if kind != reflect.String {
		panic(fmt.Errorf("postcode_iso3166_alpha2_field: bad field type %T", cfield.Interface()))
	}

	reg, found := postCodeRegexDict[cfield.String()]
	if !found {
		return false
	}

	return reg.MatchString(field.String())
}

// isBase64 is the validation function for validating if the current field's value is a valid base 64.
func isBase64(fl FieldLevel) bool {
	MustStringField("base64", fl)

	return IsBase64(fl.Field().String())
}

// isBase64URL is the validation function for validating if the current field's value is a valid base64 URL safe string.
func isBase64URL(fl FieldLevel) bool {
	MustStringField("base64url", fl)

	return IsBase64URL(fl.Field().String())
}

// isURI is the validation function for validating if the current field's value is a valid URI.
func isURI(fl FieldLevel) bool {
	MustStringField("uri", fl)

	return IsURI(fl.Field().String())
}

// isURL is the validation function for validating if the current field's value is a valid URL.
func isURL(fl FieldLevel) bool {
	MustStringField("url", fl)

	return IsURL(fl.Field().String())
}

// isHttpURL is the validation function for validating if the current field's value is a valid http:// URL.
func isHttpURL(fl FieldLevel) bool {
	MustStringField("httpurl", fl)

	return IsHttpURL(fl.Field().String())
}

// isHttpsURL is the validation function for validating if the current field's value is a valid https:// URL.
func isHttpsURL(fl FieldLevel) bool {
	MustStringField("httpsurl", fl)

	return IsHttpsURL(fl.Field().String())
}

// isHttpxURL is the validation function for validating if the current field's value is a valid https?:// URL.
func isHttpxURL(fl FieldLevel) bool {
	MustStringField("httpxurl", fl)

	return IsHttpxURL(fl.Field().String())
}

// isE164 is the validation function for validating if the current field's value is a valid e.164 formatted phone number.
func isE164(fl FieldLevel) bool {
	MustStringField("e164", fl)

	return IsE164(fl.Field().String())
}

// isEmail is the validation function for validating if the current field's value is a valid email address.
func isEmail(fl FieldLevel) bool {
	MustStringField("email", fl)

	return IsEmail(fl.Field().String())
}

// isHSLA is the validation function for validating if the current field's value is a valid HSLA color.
func isHSLA(fl FieldLevel) bool {
	MustStringField("hsla", fl)

	return IsHSLAColor(fl.Field().String())
}

// isHSL is the validation function for validating if the current field's value is a valid HSL color.
func isHSL(fl FieldLevel) bool {
	MustStringField("hsl", fl)

	return IsHSLColor(fl.Field().String())
}

// isRGBA is the validation function for validating if the current field's value is a valid RGBA color.
func isRGBA(fl FieldLevel) bool {
	MustStringField("rgba", fl)

	return IsRGBAColor(fl.Field().String())
}

// isRGB is the validation function for validating if the current field's value is a valid RGB color.
func isRGB(fl FieldLevel) bool {
	MustStringField("rgb", fl)

	return IsRGBColor(fl.Field().String())
}

// isHexColor is the validation function for validating if the current field's value is a valid HEX color.
func isHexColor(fl FieldLevel) bool {
	MustStringField("hexcolor", fl)

	return IsHexColor(fl.Field().String())
}

// isHexadecimal is the validation function for validating if the current field's value is a valid hexadecimal.
func isHexadecimal(fl FieldLevel) bool {
	MustStringField("hexadecimal", fl)

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

// isSize is the validation function for validating if the current field's value is a valid size value.
func isSize(fl FieldLevel) bool {
	MustStringField("size", fl)

	_, err := num.ParseSize(fl.Field().String())
	return err == nil
}

// isSizeGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isSizeGte(fl FieldLevel) bool {
	MustStringField("sizemin", fl)

	sz, err := num.ParseSize(fl.Field().String())
	if err != nil {
		return false
	}

	mz, err := num.ParseSize(fl.Param())
	if err != nil {
		panic(fmt.Errorf("sizemin: %w", err))
	}

	return sz >= mz
}

// isSizeLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func isSizeLte(fl FieldLevel) bool {
	MustStringField("sizemax", fl)

	sz, err := num.ParseSize(fl.Field().String())
	if err != nil {
		return false
	}

	mz, err := num.ParseSize(fl.Param())
	if err != nil {
		panic(fmt.Errorf("sizemax: %w", err))
	}

	return sz <= mz
}

// isSizeBtw is the validation function for validating if the current field's value is between the param's value "min~max".
func isSizeBtw(fl FieldLevel) bool {
	MustStringField("sizebtw", fl)

	sz, err := num.ParseSize(fl.Field().String())
	if err != nil {
		return false
	}

	p1, p2 := split2(fl.Param())

	z1, err := num.ParseSize(p1)
	if err != nil {
		panic(fmt.Errorf("sizebtw: %w", err))
	}

	z2, err := num.ParseSize(p2)
	if err != nil {
		panic(fmt.Errorf("sizebtw: %w", err))
	}

	return sz >= z1 && sz <= z2
}

// isLetter is the validation function for validating if the current field's value is a valid letter value.
func isLetter(fl FieldLevel) bool {
	MustStringField("letter", fl)

	return str.IsLetter(fl.Field().String())
}

// isLetterNumber is the validation function for validating if the current field's value is a valid alphanumeric value.
func isLetterNumber(fl FieldLevel) bool {
	MustStringField("letternumber", fl)

	return str.IsLetterNumber(fl.Field().String())
}

// isUTFLetter is the validation function for validating if the current field's value is a valid alpha unicode value.
func isUTFLetter(fl FieldLevel) bool {
	MustStringField("utfletter", fl)

	return str.IsUTFLetter(fl.Field().String())
}

// isUTFLetterNumber is the validation function for validating if the current field's value is a valid alphanumeric unicode value.
func isUTFLetterNumber(fl FieldLevel) bool {
	MustStringField("utfletternumber", fl)

	return str.IsUTFLetterNumber(fl.Field().String())
}

// isBoolean is the validation function for validating if the current field's value can be safely converted to a boolean.
func isBoolean(fl FieldLevel) bool {
	_, err := strconv.ParseBool(fl.Field().String())
	return err == nil
}

// isEmpty is the validation function for validating if the current field's value is not the default static value.
func isEmpty(fl FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return field.IsNil()
	default:
		if fl.(*validate).fldIsPointer {
			return field.Interface() == nil
		}
		return !field.IsValid() || field.IsZero()
	}
}

// isRequired is the validation function for validating if the current field's value is not the default static value.
func isRequired(fl FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		if fl.(*validate).fldIsPointer {
			return field.Interface() != nil
		}
		return field.IsValid() && !field.IsZero()
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
func requireCheckFieldValue(tag string, fl FieldLevel, param string, value string, defaultNotFoundValue bool) bool {
	field, kind, _, found := fl.GetStructFieldOKAdvanced2(fl.Parent(), param)
	if !found {
		return defaultNotFoundValue
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == asInt(tag, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() == asUint(tag, value)
	case reflect.Float32, reflect.Float64:
		return field.Float() == asFloat(tag, value)
	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) == asInt(tag, value)
	case reflect.Bool:
		return field.Bool() == asBool(tag, value)
	}

	// default reflect.String:
	return field.String() == value
}

// requiredIf is the validation function
// The field under validation must be present and not empty only if all the other specified fields are equal to the value following with the specified field.
func requiredIf(fl FieldLevel) bool {
	params := splits(fl.Param())
	if len(params)%2 != 0 {
		panic(fmt.Errorf("required_if: bad param number for  %s", fl.FieldName()))
	}

	for i := 0; i < len(params); i += 2 {
		if !requireCheckFieldValue("required_if", fl, params[i], params[i+1], false) {
			return true
		}
	}
	return isRequired(fl)
}

// requiredUnless is the validation function
// The field under validation must be present and not empty only unless all the other specified fields are equal to the value following with the specified field.
func requiredUnless(fl FieldLevel) bool {
	params := splits(fl.Param())
	if len(params)%2 != 0 {
		panic(fmt.Errorf("required_unless: bad param number for %s", fl.FieldName()))
	}

	for i := 0; i < len(params); i += 2 {
		if requireCheckFieldValue("required_unless", fl, params[i], params[i+1], false) {
			return true
		}
	}
	return isRequired(fl)
}

// requiredWith is the validation function
// The field under validation must be present and not empty only if any of the other specified fields are present.
func requiredWith(fl FieldLevel) bool {
	params := splits(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return isRequired(fl)
		}
	}
	return true
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
	return isRequired(fl)
}

// requiredWithout is the validation function
// The field under validation must be present and not empty only when any of the other specified fields are not present.
func requiredWithout(fl FieldLevel) bool {
	if requireCheckFieldKind(fl, strings.TrimSpace(fl.Param()), true) {
		return isRequired(fl)
	}
	return true
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
	return isRequired(fl)
}

// excludedWith is the validation function
// The field under validation must not be present or is empty if any of the other specified fields are present.
func excludedWith(fl FieldLevel) bool {
	params := splits(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return isEmpty(fl)
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
	return isEmpty(fl)
}

// excludedWithout is the validation function
// The field under validation must not be present or is empty when any of the other specified fields are not present.
func excludedWithout(fl FieldLevel) bool {
	if requireCheckFieldKind(fl, strings.TrimSpace(fl.Param()), true) {
		return isEmpty(fl)
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
	return isEmpty(fl)
}

func compareLength(tag string, fl FieldLevel) int {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		return utf8.RuneCountInString(field.String()) - int(asInt(tag, param))
	case reflect.Slice, reflect.Map, reflect.Array:
		return field.Len() - int(asInt(tag, param))
	}

	panic(fmt.Errorf("%s: bad field type %T", tag, field.Interface()))
}

// isLen is the validation function for validating if the current field's length or rune count is equal to the param's value.
func isLen(fl FieldLevel) bool {
	return compareLength("len", fl) == 0
}

// isMinLen is the validation function for validating if the current field's length or rune count is greater than or equal to the param's value.
func isMinLen(fl FieldLevel) bool {
	return compareLength("minlen", fl) >= 0
}

// isMaxLen is the validation function for validating if the current field's length or rune count is less than or equal to the param's value.
func isMaxLen(fl FieldLevel) bool {
	return compareLength("maxlen", fl) <= 0
}

// isBtwLen is the validation function for validating if the current field's length or rune count is between the param's value.
func isBtwLen(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		p1, p2 := asInt2("btwlen", param)
		l := int64(utf8.RuneCountInString(field.String()))
		return l >= p1 && l <= p2
	case reflect.Slice, reflect.Map, reflect.Array:
		p1, p2 := asInt2("btwlen", param)
		l := int64(field.Len())
		return l >= p1 && l <= p2
	}

	panic(fmt.Errorf("btwlen: bad field type %T", field.Interface()))
}

func compareValue(tag string, fl FieldLevel) int {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.Bool:
		return gog.If(field.Bool(), 1, 0) - gog.If(asBool("eq", param), 1, 0)
	case reflect.String:
		return str.Compare(field.String(), param)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cmp.Compare(field.Int(), asIntFromType(tag, field.Type(), param))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return cmp.Compare(field.Uint(), asUint(tag, param))
	case reflect.Float32, reflect.Float64:
		return cmp.Compare(field.Float(), asFloat(tag, param))
	case reflect.Struct:
		if field.Type().ConvertibleTo(timeType) {
			t := field.Convert(timeType).Interface().(time.Time)
			p := asTime(tag, param)
			return t.Compare(p)
		}
	}

	panic(fmt.Errorf("%s: bad field type %T", tag, field.Interface()))
}

// isNe is the validation function for validating that the field's value does not equal the provided param value.
func isNe(fl FieldLevel) bool {
	return compareValue("ne", fl) != 0
}

// isEq is the validation function for validating if the current field's value is equal to the param's value.
func isEq(fl FieldLevel) bool {
	return compareValue("eq", fl) == 0
}

// isGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isGte(fl FieldLevel) bool {
	return compareValue("gte", fl) >= 0
}

// isGt is the validation function for validating if the current field's value is greater than the param's value.
func isGt(fl FieldLevel) bool {
	return compareValue("gt", fl) > 0
}

// isLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func isLte(fl FieldLevel) bool {
	return compareValue("lte", fl) <= 0
}

// isLt is the validation function for validating if the current field's value is less than the param's value.
func isLt(fl FieldLevel) bool {
	return compareValue("lt", fl) < 0
}

// isMin is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isMin(fl FieldLevel) bool {
	return compareValue("min", fl) >= 0
}

// isMax is the validation function for validating if the current field's value is less than or equal to the param's value.
func isMax(fl FieldLevel) bool {
	return compareValue("max", fl) <= 0
}

// isBtw is the validation function for validating if the current field's value is between the param's value "min~max".
func isBtw(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		p1, p2 := split2(param)
		return field.String() >= p1 && field.String() <= p2
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p1, p2 := asInt2FromType("btw", field.Type(), param)
		i := field.Int()
		return i >= p1 && i <= p2
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p1, p2 := asUint2("btw", param)
		u := field.Uint()
		return u >= p1 && u <= p2
	case reflect.Float32, reflect.Float64:
		p1, p2 := asFloat2("btw", param)
		f := field.Float()
		return f >= p1 && f <= p2
	case reflect.Struct:
		if field.Type().ConvertibleTo(timeType) {
			p1, p2 := asTime2("btw", param)
			t := field.Convert(timeType).Interface().(time.Time)
			return !t.Before(p1) && !t.After(p2)
		}
	}

	panic(fmt.Errorf("btw: bad field type %T", field.Interface()))
}

// isNeField is the validation function for validating if the current field's value is not equal to the field specified by the param's value.
func isNeField(fl FieldLevel) bool {
	rv, ok := compareField(fl)
	return !ok || rv != 0
}

// isEqField is the validation function for validating if the current field's value is equal to the field specified by the param's value.
func isEqField(fl FieldLevel) bool {
	rv, ok := compareField(fl)
	return ok && rv == 0
}

func compareField(fl FieldLevel) (int, bool) {
	field := fl.Field()
	kind := field.Kind()

	cfield, cfkind, ok := fl.GetStructFieldOK()
	if !ok || cfkind != kind {
		return 0, false
	}

	switch kind {
	case reflect.Bool:
		return gog.If(field.Bool(), 1, 0) - gog.If(cfield.Bool(), 1, 0), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cmp.Compare(field.Int(), cfield.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return cmp.Compare(field.Uint(), cfield.Uint()), true
	case reflect.Float32, reflect.Float64:
		return cmp.Compare(field.Float(), cfield.Float()), true
	case reflect.Struct:
		fieldType := field.Type()
		if fieldType.ConvertibleTo(timeType) && cfield.Type().ConvertibleTo(timeType) {
			ft := field.Convert(timeType).Interface().(time.Time)
			ct := cfield.Convert(timeType).Interface().(time.Time)
			return ft.Compare(ct), true
		}

		// Not Same underlying type i.e. struct and time
		if fieldType != cfield.Type() {
			return 0, false
		}
	}

	// default reflect.String
	return str.Compare(field.String(), cfield.String()), true
}

// isGteField is the validation function for validating if the current field's value is greater than or equal to the field specified by the param's value.
func isGteField(fl FieldLevel) bool {
	rv, ok := compareField(fl)
	return ok && rv >= 0
}

// isGtField is the validation function for validating if the current field's value is greater than the field specified by the param's value.
func isGtField(fl FieldLevel) bool {
	rv, ok := compareField(fl)
	return ok && rv > 0
}

// isLteField is the validation function for validating if the current field's value is less than or equal to the field specified by the param's value.
func isLteField(fl FieldLevel) bool {
	rv, ok := compareField(fl)
	return ok && rv <= 0
}

// isLtField is the validation function for validating if the current field's value is less than the field specified by the param's value.
func isLtField(fl FieldLevel) bool {
	rv, ok := compareField(fl)
	return ok && rv < 0
}

func isHostnameRFC952(fl FieldLevel) bool {
	MustStringField("hostname", fl)

	return IsHostnameRFC952(fl.Field().String())
}

func isHostnameRFC1123(fl FieldLevel) bool {
	MustStringField("hostname_rfc1123", fl)

	return IsHostnameRFC1123(fl.Field().String())
}

func isFQDN(fl FieldLevel) bool {
	MustStringField("fqdn", fl)

	return IsFQDN(fl.Field().String())
}

// isJSON is the validation function for validating if the current field's value is a valid json string.
func isJSON(fl FieldLevel) bool {
	MustStringField("json", fl)

	return json.Valid(str.UnsafeBytes(fl.Field().String()))
}

// isJSONObject is the validation function for validating if the current field's value is a valid json object string.
func isJSONObject(fl FieldLevel) bool {
	MustStringField("jsonobject", fl)

	return IsJSONObject(fl.Field().String())
}

// isJSONArray is the validation function for validating if the current field's value is a valid json array string.
func isJSONArray(fl FieldLevel) bool {
	MustStringField("jsonarray", fl)

	return IsJSONArray(fl.Field().String())
}

// isJWT is the validation function for validating if the current field's value is a valid JWT string.
func isJWT(fl FieldLevel) bool {
	MustStringField("jwt", fl)

	return IsJWT(fl.Field().String())
}

// isHostnamePort validates a <dns>:<port> combination for fields typically used for socket address.
func isHostnamePort(fl FieldLevel) bool {
	MustStringField("hostname_port", fl)

	return IsHostnamePort(fl.Field().String())
}

// isLowercase is the validation function for validating if the current field's value is a lowercase string.
func isLowercase(fl FieldLevel) bool {
	MustStringField("lowercase", fl)

	return str.IsLowerCase(fl.Field().String())
}

// isUppercase is the validation function for validating if the current field's value is an uppercase string.
func isUppercase(fl FieldLevel) bool {
	MustStringField("uppercase", fl)

	return str.IsUpperCase(fl.Field().String())
}

// isDatetime is the validation function for validating if the current field's value is a valid datetime string.
func isDatetime(fl FieldLevel) bool {
	MustStringField("datetime", fl)

	field := fl.Field()
	param := fl.Param()

	_, err := time.Parse(param, field.String())
	return err == nil
}

// isDuration is the validation function for validating if the current field's value is a valid duration string.
func isDuration(fl FieldLevel) bool {
	MustStringField("duration", fl)

	_, err := tmu.ParseDuration(fl.Field().String())
	return err == nil
}

// isCron is the validation function for validating if the current field's value is a valid cron expression.
func isCron(fl FieldLevel) bool {
	MustStringField("cron", fl)

	_, err := sch.ParseCron(fl.Field().String())
	return err == nil
}

// isPeriodic is the validation function for validating if the current field's value is a valid periodic expression.
func isPeriodic(fl FieldLevel) bool {
	MustStringField("periodic", fl)

	_, err := sch.ParsePeriodic(fl.Field().String())
	return err == nil
}

// isTimeZone is the validation function for validating if the current field's value is a valid time zone string.
func isTimeZone(fl FieldLevel) bool {
	MustStringField("timezone", fl)

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
	MustStringField("swiftcode", fl)

	return IsSwiftCode(fl.Field().String())
}

// isSemver is the validation function for validating if the current field's value is a valid semver version, defined in Semantic Versioning 2.0.0
func isSemver(fl FieldLevel) bool {
	MustStringField("semver", fl)

	return IsSemver(fl.Field().String())
}

// isDnsRFC1035Label is the validation function
// for validating if the current field's value is
// a valid dns RFC 1035 label, defined in RFC 1035.
func isDnsRFC1035Label(fl FieldLevel) bool {
	MustStringField("dns_rfc1035_label", fl)

	return IsDnsRFC1035Label(fl.Field().String())
}
