package col

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type jsonArray interface {
	addJSONArrayItem(v T) jsonArray
}

type jsonObject interface {
	addJSONObjectItem(k string, v T)
}

// JSONArray json array type
type JSONArray []T

func (ja JSONArray) addJSONArrayItem(v T) jsonArray {
	return append(ja, v)
}

func newJSONArray() jsonArray {
	return JSONArray([]T{})
}

// JSONObject json object type
type JSONObject map[string]T

func (jo JSONObject) addJSONObjectItem(k string, v T) {
	jo[k] = v
}

func newJSONObject() jsonObject {
	return JSONObject(make(map[string]T))
}

func jsonUnmarshalArray(data []byte, ja jsonArray) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	// must open with a delim token '['
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("expect JSON array open with '['")
	}

	_, err = jsonParseArray(dec, ja)
	if err != nil {
		return err
	}

	t, err = dec.Token()
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("expect end of JSON array but got more token: %T: %v or err: %w", t, t, err)
	}

	return nil
}

func jsonUnmarshalObject(data []byte, jo jsonObject) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	// must open with a delim token '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	_, err = jsonParseObject(dec, jo)
	if err != nil {
		return err
	}

	t, err = dec.Token()
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("expect end of JSON object but got more token: %T: %v or err: %w", t, t, err)
	}

	return nil
}

func jsonParseObject(dec *json.Decoder, jo jsonObject) (jsonObject, error) {
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return nil, err
		}

		k, ok := t.(string)
		if !ok {
			return nil, fmt.Errorf("expecting JSON key should be always a string: %T: %v", t, t)
		}

		t, err = dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		var v T
		v, err = jsonHandleDelim(dec, t)
		if err != nil {
			return nil, err
		}

		jo.addJSONObjectItem(k, v)
	}

	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return nil, fmt.Errorf("expect JSON object close with '}'")
	}

	return jo, nil
}

func jsonParseArray(dec *json.Decoder, ja jsonArray) (jsonArray, error) {
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return nil, err
		}

		v, err := jsonHandleDelim(dec, t)
		if err != nil {
			return nil, err
		}
		ja = ja.addJSONArrayItem(v)
	}

	t, err := dec.Token()
	if err != nil {
		return nil, err
	}

	if delim, ok := t.(json.Delim); !ok || delim != ']' {
		return nil, fmt.Errorf("expect JSON array close with ']'")
	}

	return ja, nil
}

func jsonHandleDelim(dec *json.Decoder, t json.Token) (T, error) {
	if delim, ok := t.(json.Delim); ok {
		switch delim {
		case '{':
			jo, err := jsonParseObject(dec, newJSONObject())
			if err != nil {
				return nil, err
			}
			return jo, nil
		case '[':
			ja, err := jsonParseArray(dec, newJSONArray())
			if err != nil {
				return nil, err
			}
			return ja, nil
		default:
			return nil, fmt.Errorf("Unexpected delimiter: %q", delim)
		}
	}
	return t, nil
}

// ---------------------------------------------------------------------
func jsonMarshalArray(col Collection) (res []byte, err error) {
	if col.IsEmpty() {
		return []byte("[]"), nil
	}

	res = append(res, '[')
	col.Each(func(_ int, v T) bool {
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

func jsonMarshalObject(m Map) (res []byte, err error) {
	if m.IsEmpty() {
		return []byte("{}"), nil
	}

	res = append(res, '{')
	m.Each(func(k K, v V) bool {
		var bs []byte
		k = fmt.Sprintf("%v", k)
		res = append(res, fmt.Sprintf("%q:", k)...)
		bs, err = json.Marshal(v)
		if err != nil {
			return false
		}
		res = append(res, bs...)
		res = append(res, ',')
		return true
	})
	res[len(res)-1] = '}'
	return
}
