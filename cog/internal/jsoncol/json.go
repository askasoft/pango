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

	var v T

	tv := reflect.TypeOf(v)
	for dec.More() {
		if err := dec.Decode(&v); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		cv, err := ref.CastTo(v, tv)
		if err != nil {
			return err
		}

		c.Add(cv.(T))
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
