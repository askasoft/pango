//go:build go1.18
// +build go1.18

package jsoncol

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

func jsonAddColItem[T any](c cog.Collection[T], v any) error {
	var tv T

	av, err := ref.Convert(v, reflect.TypeOf(tv))
	if err != nil {
		return err
	}

	c.Add(av.(T))
	return nil
}

func JsonUnmarshalCol[T any](data []byte, c cog.Collection[T]) error {
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

// ---------------------------------------------------------------------
func JsonMarshalCol[T any](c cog.Collection[T]) (res []byte, err error) {
	if c.IsEmpty() {
		return []byte("[]"), nil
	}

	res = append(res, '[')
	c.Each(func(_ int, v T) bool {
		var bs []byte
		bs, err = json.Marshal(v)
		if err != nil {
			return false
		}
		res = append(res, bs...)
		res = append(res, ',')
		return true
	})
	res[len(res)-1] = ']'

	return
}
