package jsonmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/ref"
)

func jsonAddMapItem[K any, V any](am cog.Map[K, V], k string, v any) error {
	var tk K
	var tv V

	ak, err := ref.Convert(k, reflect.TypeOf(tk))
	if err != nil {
		return err
	}

	av, err := ref.Convert(v, reflect.TypeOf(tv))
	if err != nil {
		return err
	}

	am.Set(ak.(K), av.(V))
	return nil
}

func JsonUnmarshalMap[K any, V any](data []byte, am cog.Map[K, V]) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	// must open with a delim token '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return err
		}

		k, ok := t.(string)
		if !ok {
			return fmt.Errorf("expecting JSON key should be always a string: %T: %v", t, t)
		}

		var v V
		err = dec.Decode(&v)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		err = jsonAddMapItem(am, k, v)
		if err != nil {
			return err
		}
	}

	t, err = dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expect JSON object close with '}'")
	}

	t, err = dec.Token()
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("expect end of JSON object but got more token: %T: %v or err: %w", t, t, err)
	}

	return nil
}

// ---------------------------------------------------------------------
func JsonMarshalMap[K any, V any](m cog.Map[K, V]) ([]byte, error) {
	if m.IsEmpty() {
		return []byte("{}"), nil
	}

	var err error

	bb := &bytes.Buffer{}
	bb.WriteByte('{')

	je := json.NewEncoder(bb)
	m.Each(func(k K, v V) bool {
		if _, err = fmt.Fprintf(bb, "%q:", fmt.Sprint(k)); err != nil {
			return false
		}
		if err = je.Encode(v); err != nil {
			return false
		}

		// remove last '\n'
		bs := bb.Bytes()
		bs[len(bs)-1] = ','
		return true
	})

	bs := bb.Bytes()
	bs[len(bs)-1] = '}'
	return bs, err
}
