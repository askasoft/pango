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
func JsonMarshalCol[T any](c cog.Collection[T]) ([]byte, error) {
	if c.IsEmpty() {
		return []byte("[]"), nil
	}

	var err error

	bb := &bytes.Buffer{}
	bb.WriteByte('[')

	je := json.NewEncoder(bb)
	c.Each(func(_ int, v T) bool {
		if err = je.Encode(v); err != nil {
			return false
		}

		// remove last '\n'
		bs := bb.Bytes()
		bs[len(bs)-1] = ','
		return true
	})

	bs := bb.Bytes()
	bs[len(bs)-1] = ']'
	return bs, err
}
