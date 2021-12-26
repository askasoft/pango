package col

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type jsonArray interface {
	addJSONArrayItem(v T) jsonArray
}

type jsonObject interface {
	addJSONObjectItem(k string, v T) jsonObject
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

func (jo JSONObject) addJSONObjectItem(k string, v T) jsonObject {
	jo[k] = v
	return jo
}

func newJSONObject() jsonObject {
	return JSONObject(make(map[string]T))
}

type jsonUnmarshaler struct {
	newArray  func() jsonArray
	newObject func() jsonObject
}

func (jd *jsonUnmarshaler) unmarshalJSONArray(data []byte, ja jsonArray) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	// must open with a delim token '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("expect JSON array open with '['")
	}

	ja, err = jd.parseJSONArray(dec, ja)
	if err != nil {
		return err
	}

	t, err = dec.Token()
	if err != io.EOF {
		return fmt.Errorf("expect end of JSON array but got more token: %T: %v or err: %v", t, t, err)
	}

	return nil
}

func (jd *jsonUnmarshaler) unmarshalJSONObject(data []byte, jo jsonObject) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	// must open with a delim token '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	jo, err = jd.parseJSONObject(dec, jo)
	if err != nil {
		return err
	}

	t, err = dec.Token()
	if err != io.EOF {
		return fmt.Errorf("expect end of JSON object but got more token: %T: %v or err: %v", t, t, err)
	}

	return nil
}

func (jd *jsonUnmarshaler) parseJSONObject(dec *json.Decoder, jo jsonObject) (jsonObject, error) {
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
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		var v T
		v, err = jd.handleDelim(dec, t)
		if err != nil {
			return nil, err
		}

		jo = jo.addJSONObjectItem(k, v)
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

func (jd *jsonUnmarshaler) parseJSONArray(dec *json.Decoder, ja jsonArray) (jsonArray, error) {
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return nil, err
		}

		v, err := jd.handleDelim(dec, t)
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

func (jd *jsonUnmarshaler) handleDelim(dec *json.Decoder, t json.Token) (T, error) {
	if delim, ok := t.(json.Delim); ok {
		switch delim {
		case '{':
			jo, err := jd.parseJSONObject(dec, jd.newObject())
			if err != nil {
				return nil, err
			}
			return jo, nil
		case '[':
			ja, err := jd.parseJSONArray(dec, jd.newArray())
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

//---------------------------------------------------------------------
func jsonMarshalIter(it Iterator) (res []byte, err error) {
	res = append(res, '[')
	for it.Next() {
		var b []byte
		b, err = json.Marshal(it.Value())
		if err != nil {
			return
		}
		res = append(res, b...)
		res = append(res, ',')
	}
	res[len(res)-1] = ']'
	return
}

func jsonMarshalSet(set Set) (res []byte, err error) {
	if set.IsEmpty() {
		return []byte("[]"), nil
	}

	if ic, ok := set.(Iterable); ok {
		return jsonMarshalIter(ic.Iterator())
	}

	var bs []byte
	res = append(res, '[')
	set.Each(func(v T) {
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

func jsonMarshalList(list List) (res []byte, err error) {
	if list.IsEmpty() {
		return []byte("[]"), nil
	}

	return jsonMarshalIter(list.Iterator())
}

func jsonMarshalHashMap(hmap map[K]V) (res []byte, err error) {
	if len(hmap) == 0 {
		return []byte("{}"), nil
	}

	var bs []byte
	res = append(res, '{')
	for k, v := range hmap {
		_, ok := k.(string)
		if !ok {
			err = fmt.Errorf("expecting JSON key should be always a string: %T: %v", k, k)
			return
		}

		res = append(res, fmt.Sprintf("%q:", k)...)
		bs, err = json.Marshal(v)
		if err != nil {
			return
		}
		res = append(res, bs...)
		res = append(res, ',')
	}
	res[len(res)-1] = '}'
	return
}

func jsonMarshalMap(im IterableMap) (res []byte, err error) {
	if im.IsEmpty() {
		return []byte("{}"), nil
	}

	var bs []byte

	res = append(res, '{')
	it := im.Iterator()
	for it.Next() {
		k, ok := it.Key().(string)
		if !ok {
			err = fmt.Errorf("expecting JSON key should be always a string: %T: %v", it.Key(), it.Key())
			return
		}

		res = append(res, fmt.Sprintf("%q:", k)...)
		bs, err = json.Marshal(it.Value())
		if err != nil {
			return
		}
		res = append(res, bs...)
		res = append(res, ',')
	}
	res[len(res)-1] = '}'
	return
}
