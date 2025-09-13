package xmlrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

const (
	iso8601        = "20060102T15:04:05"
	iso8601Z       = "20060102T15:04:05Z07:00"
	iso8601Hyphen  = "2006-01-02T15:04:05"
	iso8601HyphenZ = "2006-01-02T15:04:05Z07:00"
)

var (
	timeLayouts = []string{iso8601, iso8601Z, iso8601Hyphen, iso8601HyphenZ}
)

type TypeMismatchError string

func (e TypeMismatchError) Error() string { return string(e) }

type Decoder struct {
	*xml.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{xml.NewDecoder(r)}
}

func Decode(r io.Reader, v any) error {
	return NewDecoder(r).Decode(v)
}

func DecodeBytes(xml []byte, v any) error {
	return NewDecoder(bytes.NewReader(xml)).Decode(v)
}

func DecodeString(xml string, v any) error {
	return NewDecoder(strings.NewReader(xml)).Decode(v)
}

func (d *Decoder) Decode(v any) error {
	if err := d.decode(v); err != nil {
		return err
	}

	// read until end of document
	if err := d.Skip(); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func (d *Decoder) decode(v any) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		if t, ok := tok.(xml.StartElement); ok {
			switch t.Name.Local {
			case "value":
				rv := reflect.ValueOf(v)
				if rv.Kind() != reflect.Ptr {
					return errors.New("xmlrpc decode error: cannot decode to non-pointer value")
				}
				return d.decodeValue(rv.Elem())
			case "fault":
				fe := &FaultError{}
				if err := d.decode(fe); err != nil {
					return err
				}
				return fe
			}
		}
	}
}

func (d *Decoder) decodeValue(rv reflect.Value) error {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	var tok xml.Token
	var err error
	var tag string

	for {
		if tok, err = d.Token(); err != nil {
			return err
		}

		if t, ok := tok.(xml.EndElement); ok {
			if t.Name.Local == "value" {
				return nil
			}
			return fmt.Errorf("xmlrpc decode error: invalid element <%s>", t.Name.Local)
		}

		if t, ok := tok.(xml.StartElement); ok {
			tag = t.Name.Local
			break
		}

		// Treat value data without type identifier as string
		if t, ok := tok.(xml.CharData); ok {
			if value := strings.TrimSpace(string(t)); value != "" {
				if err = d.checkType(rv, reflect.String); err != nil {
					return err
				}

				rv.SetString(value)
				return nil
			}
		}
	}

	switch tag {
	case "struct":
		ismap := false
		pmap := rv
		valType := rv.Type()

		if err = d.checkType(rv, reflect.Struct); err != nil {
			if d.checkType(rv, reflect.Map) == nil {
				if valType.Key().Kind() != reflect.String {
					return fmt.Errorf("xmlrpc decode error: invalid map key type %v", valType.Key().Kind())
				}
				ismap = true
			} else if d.checkType(rv, reflect.Interface) == nil && rv.IsNil() {
				var dummy map[string]any
				valType = reflect.TypeOf(dummy)
				pmap = reflect.New(valType).Elem()
				rv.Set(pmap)
				ismap = true
			} else {
				return err
			}
		}

		var fields map[string]reflect.Value

		if !ismap {
			fields = make(map[string]reflect.Value)

			for i := 0; i < valType.NumField(); i++ {
				field := valType.Field(i)
				fieldVal := rv.FieldByName(field.Name)

				if fieldVal.CanSet() {
					name := field.Tag.Get("xmlrpc")
					name = strings.TrimSuffix(name, ",omitempty")
					if name == "-" {
						continue
					}
					if name == "" {
						name = field.Name
					}
					fields[name] = fieldVal
				}
			}
		} else {
			// Create initial empty map
			pmap.Set(reflect.MakeMap(valType))
		}

		// Process struct members.
	StructLoop:
		for {
			if tok, err = d.Token(); err != nil {
				return err
			}

			switch t := tok.(type) {
			case xml.StartElement:
				if t.Name.Local != "member" {
					return fmt.Errorf("xmlrpc decode error: invalid element <%s>", t.Name.Local)
				}

				tagName, fieldName, err := d.readTag()
				if err != nil {
					return err
				}
				if tagName != "name" {
					return fmt.Errorf("xmlrpc decode error: invalid element <%s>", tagName)
				}

				var fv reflect.Value
				ok := true

				if !ismap {
					fv, ok = fields[string(fieldName)]
				} else {
					fv = reflect.New(valType.Elem())
				}

				if ok {
					for {
						if tok, err = d.Token(); err != nil {
							return err
						}
						if t, ok := tok.(xml.StartElement); ok && t.Name.Local == "value" {
							if err = d.decodeValue(fv); err != nil {
								return err
							}

							// </value>
							if err = d.Skip(); err != nil {
								return err
							}

							break
						}
					}
				}

				// </member>
				if err = d.Skip(); err != nil {
					return err
				}

				if ismap {
					pmap.SetMapIndex(reflect.ValueOf(string(fieldName)), reflect.Indirect(fv))
					rv.Set(pmap)
				}
			case xml.EndElement:
				break StructLoop
			}
		}
	case "array":
		slice := rv
		if d.checkType(rv, reflect.Interface) == nil && rv.IsNil() {
			slice = reflect.ValueOf([]any{})
		} else if err = d.checkType(rv, reflect.Slice); err != nil {
			return err
		}

	ArrayLoop:
		for {
			if tok, err = d.Token(); err != nil {
				return err
			}

			switch t := tok.(type) {
			case xml.StartElement:
				var index int
				if t.Name.Local != "data" {
					return fmt.Errorf("xmlrpc decode error: invalid element <%s>", t.Name.Local)
				}
			DataLoop:
				for {
					if tok, err = d.Token(); err != nil {
						return err
					}

					switch tt := tok.(type) {
					case xml.StartElement:
						if tt.Name.Local != "value" {
							return fmt.Errorf("xmlrpc decode error: invalid element <%s>", t.Name.Local)
						}

						if index < slice.Len() {
							v := slice.Index(index)
							if v.Kind() == reflect.Interface {
								v = v.Elem()
							}
							if v.Kind() != reflect.Ptr {
								return errors.New("xmlrpc decode error: cannot decode to non-pointer array element")
							}
							if err = d.decodeValue(v); err != nil {
								return err
							}
						} else {
							v := reflect.New(slice.Type().Elem())
							if err = d.decodeValue(v); err != nil {
								return err
							}
							slice = reflect.Append(slice, v.Elem())
						}

						// </value>
						if err = d.Skip(); err != nil {
							return err
						}
						index++
					case xml.EndElement:
						rv.Set(slice)
						break DataLoop
					}
				}
			case xml.EndElement:
				break ArrayLoop
			}
		}
	default:
		if tok, err = d.Token(); err != nil {
			return err
		}

		var data string

		switch t := tok.(type) {
		case xml.EndElement:
			return nil
		case xml.CharData:
			data = str.UnsafeString(t)
		default:
			return fmt.Errorf("xmlrpc decode error: invalid token %v", t)
		}

		switch tag {
		case "int", "i4", "i8":
			if d.checkType(rv, reflect.Interface) == nil && rv.IsNil() {
				i, err := strconv.ParseInt(data, 10, 64)
				if err != nil {
					return err
				}

				pi := reflect.New(reflect.TypeOf(i)).Elem()
				pi.SetInt(i)
				rv.Set(pi)
			} else if err = d.checkType(rv, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64); err != nil {
				return err
			} else {
				i, err := strconv.ParseInt(data, 10, rv.Type().Bits())
				if err != nil {
					return err
				}

				rv.SetInt(i)
			}
		case "string":
			if d.checkType(rv, reflect.Interface) == nil && rv.IsNil() {
				pstr := reflect.New(reflect.TypeOf(data)).Elem()
				pstr.SetString(str.Clone(data))
				rv.Set(pstr)
			} else if err = d.checkType(rv, reflect.String); err != nil {
				return err
			} else {
				rv.SetString(str.Clone(data))
			}
		case "base64":
			bs, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				return err
			}

			if d.checkType(rv, reflect.Interface) == nil && rv.IsNil() {
				pstr := reflect.New(reflect.TypeOf(bs)).Elem()
				pstr.Set(reflect.ValueOf(bs))
				rv.Set(pstr)
			} else if rv.Type() != ref.TypeBytes {
				return TypeMismatchError(fmt.Sprintf("error: type mismatch error - can't decode []byte to %v", rv.Kind()))
			} else {
				rv.Set(reflect.ValueOf(bs))
			}
		case "dateTime.iso8601":
			t, err := tmu.Parse(data, timeLayouts...)
			if err != nil {
				return err
			}

			if d.checkType(rv, reflect.Interface) == nil && rv.IsNil() {
				ptime := reflect.New(reflect.TypeOf(t)).Elem()
				ptime.Set(reflect.ValueOf(t))
				rv.Set(ptime)
			} else if _, ok := rv.Interface().(time.Time); !ok {
				return TypeMismatchError(fmt.Sprintf("error: type mismatch error - can't decode time to %v", rv.Kind()))
			} else {
				rv.Set(reflect.ValueOf(t))
			}
		case "boolean":
			v, err := strconv.ParseBool(data)
			if err != nil {
				return err
			}

			if d.checkType(rv, reflect.Interface) == nil && rv.IsNil() {
				pv := reflect.New(reflect.TypeOf(v)).Elem()
				pv.SetBool(v)
				rv.Set(pv)
			} else if err = d.checkType(rv, reflect.Bool); err != nil {
				return err
			} else {
				rv.SetBool(v)
			}
		case "double":
			if d.checkType(rv, reflect.Interface) == nil && rv.IsNil() {
				i, err := strconv.ParseFloat(data, 64)
				if err != nil {
					return err
				}

				pdouble := reflect.New(reflect.TypeOf(i)).Elem()
				pdouble.SetFloat(i)
				rv.Set(pdouble)
			} else if err = d.checkType(rv, reflect.Float32, reflect.Float64); err != nil {
				return err
			} else {
				i, err := strconv.ParseFloat(data, rv.Type().Bits())
				if err != nil {
					return err
				}

				rv.SetFloat(i)
			}
		default:
			return fmt.Errorf("xmlrpc decode error: unsupported type %s", tag)
		}

		// </type>
		if err = d.Skip(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) readTag() (string, []byte, error) {
	var tok xml.Token
	var err error

	var name string
	for {
		if tok, err = d.Token(); err != nil {
			return "", nil, err
		}

		if t, ok := tok.(xml.StartElement); ok {
			name = t.Name.Local
			break
		}
	}

	value, err := d.readCharData()
	if err != nil {
		return "", nil, err
	}

	return name, value, d.Skip()
}

func (d *Decoder) readCharData() ([]byte, error) {
	tok, err := d.Token()
	if err != nil {
		return nil, err
	}

	if t, ok := tok.(xml.CharData); ok {
		return []byte(t.Copy()), nil
	}
	return nil, fmt.Errorf("xmlrpc decode error: invalid token %v", tok)
}

func (d *Decoder) checkType(rv reflect.Value, kinds ...reflect.Kind) error {
	if len(kinds) == 0 {
		return nil
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	for _, kind := range kinds {
		if rv.Kind() == kind {
			return nil
		}
	}

	return TypeMismatchError(fmt.Sprintf("xmlrpc decode error: type mismatch - can't decode %v to %v", kinds[0], rv.Kind()))
}
