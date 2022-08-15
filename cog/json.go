//go:build go1.18
// +build go1.18

package cog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/pandafw/pango/ref"
)

func jsonAddColItem[T any](c Collection[T], v any) error {
	var tv T

	av, err := ref.Convert(v, reflect.TypeOf(tv))
	if err != nil {
		return err
	}

	c.Add(av.(T))
	return nil
}

func jsonUnmarshalCol[T any](data []byte, c Collection[T]) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	// must open with a delim token '['
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("expect JSON array open with '['")
	}

	for dec.More() {
		var v T
		err = dec.Decode(&v)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		err = jsonAddColItem(c, v)
		if err != nil {
			return err
		}
	}

	t, err = dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != ']' {
		return fmt.Errorf("expect JSON array close with ']'")
	}

	t, err = dec.Token()
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("expect end of JSON array but got more token: %T: %v or err: %w", t, t, err)
	}

	return nil
}

func jsonAddMapItem[K any, V any](am Map[K, V], k string, v any) error {
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

func jsonUnmarshalMap[K any, V any](data []byte, am Map[K, V]) error {
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
func jsonMarshalCol[T any](c Collection[T]) (res []byte, err error) {
	if c.IsEmpty() {
		return []byte("[]"), nil
	}

	res = append(res, '[')
	c.Each(func(v T) {
		var bs []byte
		bs, err = json.Marshal(v)
		if err != nil {
			return
		}
		res = append(res, bs...)
		res = append(res, ',')
	})
	res[len(res)-1] = ']'
	return
}

func jsonMarshalMap[K any, V any](m Map[K, V]) (res []byte, err error) {
	if m.IsEmpty() {
		return []byte("{}"), nil
	}

	res = append(res, '{')
	m.Each(func(k K, v V) {
		var bs []byte
		s := fmt.Sprintf("%v", k)
		res = append(res, fmt.Sprintf("%q:", s)...)
		bs, err = json.Marshal(v)
		if err != nil {
			return
		}
		res = append(res, bs...)
		res = append(res, ',')
	})
	res[len(res)-1] = '}'
	return
}
