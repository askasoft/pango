package xmlrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/str"
)

func EncodeMethodCall(w io.Writer, method string, args ...any) error {
	xe := NewEncoder(w)

	if err := xe.print(`<?xml version="1.0" encoding="UTF-8"?>`); err != nil {
		return err
	}
	if err := xe.printf("<methodCall><methodName>%s</methodName>", method); err != nil {
		return err
	}

	if len(args) > 0 {
		if err := xe.print("<params>"); err != nil {
			return err
		}

		for _, arg := range args {
			if err := xe.print("<param>"); err != nil {
				return err
			}

			if err := xe.Encode(arg); err != nil {
				return err
			}

			if err := xe.print("</param>"); err != nil {
				return err
			}
		}

		if err := xe.print("</params>"); err != nil {
			return err
		}
	}

	return xe.print("</methodCall>")
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func Encode(v any) ([]byte, error) {
	bb := &bytes.Buffer{}
	err := NewEncoder(bb).Encode(v)
	return bb.Bytes(), err
}

func (e *Encoder) Encode(v any) error {
	return e.encodeValue(reflect.ValueOf(v))
}

func (e *Encoder) print(a ...any) error {
	_, err := fmt.Fprint(e.w, a...)
	return err
}

func (e *Encoder) printf(f string, a ...any) error {
	_, err := fmt.Fprintf(e.w, f, a...)
	return err
}

func (e *Encoder) encodeValue(rv reflect.Value) error {
	if !rv.IsValid() {
		return e.print("<value><nil/></value>")
	}

	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return e.print("<value><nil/></value>")
		}
		rv = rv.Elem()
	}

	if err := e.print("<value>"); err != nil {
		return err
	}

	switch rv.Kind() {
	case reflect.Struct:
		switch v := rv.Interface().(type) {
		case time.Time:
			if err := e.printf("<dateTime.iso8601>%s</dateTime.iso8601>", v.Format(iso8601)); err != nil {
				return err
			}
		default:
			if err := e.encodeStruct(rv); err != nil {
				return err
			}
		}
	case reflect.Map:
		if err := e.encodeMap(rv); err != nil {
			return err
		}
	case reflect.Slice, reflect.Array:
		switch v := rv.Interface().(type) {
		case []byte:
			if err := e.encodeBase64(v); err != nil {
				return err
			}
		default:
			if err := e.encodeSlice(rv); err != nil {
				return err
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if err := e.printf("<int>%s</int>", strconv.FormatInt(rv.Int(), 10)); err != nil {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := e.printf("<i4>%s</i4>", strconv.FormatUint(rv.Uint(), 10)); err != nil {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if err := e.printf("<double>%s</double>", strconv.FormatFloat(rv.Float(), 'f', -1, rv.Type().Bits())); err != nil {
			return err
		}
	case reflect.Bool:
		if err := e.printf("<boolean>%s</boolean>", gog.If(rv.Bool(), "1", "0")); err != nil {
			return err
		}
	case reflect.String:
		if err := e.print("<string>"); err != nil {
			return err
		}
		if err := xml.EscapeText(e.w, str.UnsafeBytes(rv.String())); err != nil {
			return err
		}
		if err := e.print("</string>"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("xmlrpc encode error: unsupported type %v", rv.Type())
	}

	return e.print("</value>")
}

func (e *Encoder) encodeStruct(sv reflect.Value) error {
	if err := e.print("<struct>"); err != nil {
		return err
	}

	st := sv.Type()
	for i := range st.NumField() {
		fv := sv.Field(i)
		ft := st.Field(i)

		name := ft.Tag.Get("xmlrpc")
		// skip ignored fields.
		if name == "-" {
			continue
		}

		// if the tag has the omitempty property, skip it
		if strings.HasSuffix(name, ",omitempty") && fv.IsZero() {
			continue
		}

		name = strings.TrimSuffix(name, ",omitempty")
		if name == "" {
			name = ft.Name
		}

		if err := e.printf("<member><name>%s</name>", name); err != nil {
			return err
		}
		if err := e.encodeValue(fv); err != nil {
			return err
		}
		if err := e.print("</member>"); err != nil {
			return err
		}
	}

	return e.print("</struct>")
}

func (e *Encoder) encodeMap(rv reflect.Value) error {
	rt := rv.Type()
	if rt.Key().Kind() != reflect.String {
		return errors.New("xmlrpc encode error: only maps with string keys are supported")
	}

	if err := e.print("<struct>"); err != nil {
		return err
	}

	keys := rv.MapKeys()
	sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })

	for i := range rv.Len() {
		k := keys[i]
		v := rv.MapIndex(k)

		if err := e.printf("<member><name>%v</name>", k); err != nil {
			return err
		}

		if err := e.encodeValue(v); err != nil {
			return err
		}

		if err := e.print("</member>"); err != nil {
			return err
		}
	}

	return e.print("</struct>")
}

func (e *Encoder) encodeBase64(bs []byte) error {
	if err := e.print("<base64>"); err != nil {
		return err
	}

	be := base64.NewEncoder(base64.StdEncoding, e.w)
	if _, err := be.Write(bs); err != nil {
		return err
	}
	if err := be.Close(); err != nil {
		return err
	}

	return e.print("</base64>")
}

func (e *Encoder) encodeSlice(rv reflect.Value) error {
	if err := e.print("<array><data>"); err != nil {
		return err
	}

	for i := range rv.Len() {
		if err := e.encodeValue(rv.Index(i)); err != nil {
			return err
		}
	}

	return e.print("</data></array>")
}
