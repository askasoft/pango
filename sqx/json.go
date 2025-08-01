package sqx

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/askasoft/pango/ref"
)

var ErrNotBytes = errors.New("sqx: type assertion to []byte failed")

type JSONObject map[string]any

func (jo JSONObject) Value() (driver.Value, error) {
	if jo == nil {
		return nil, nil
	}
	return json.Marshal(jo)
}

func (jo *JSONObject) Scan(value any) error {
	return JSONScan(value, jo)
}

type JSONStringObject map[string]string

func (jso JSONStringObject) Value() (driver.Value, error) {
	if jso == nil {
		return nil, nil
	}
	return json.Marshal(jso)
}

func (jso *JSONStringObject) Scan(value any) error {
	return JSONScan(value, jso)
}

type JSONArray []any

func (ja JSONArray) Value() (driver.Value, error) {
	if ja == nil {
		return nil, nil
	}
	return json.Marshal(ja)
}

func (ja *JSONArray) Scan(value any) error {
	return JSONScan(value, ja)
}

type JSONStringArray []string

func (jsa JSONStringArray) Value() (driver.Value, error) {
	if jsa == nil {
		return nil, nil
	}
	return json.Marshal(jsa)
}

func (jsa *JSONStringArray) Scan(value any) error {
	return JSONScan(value, jsa)
}

type JSONIntArray []int

func (jia JSONIntArray) Value() (driver.Value, error) {
	if jia == nil {
		return nil, nil
	}
	return json.Marshal(jia)
}

func (jia *JSONIntArray) Scan(value any) error {
	return JSONScan(value, jia)
}

type JSONInt64Array []int64

func (jia JSONInt64Array) Value() (driver.Value, error) {
	if jia == nil {
		return nil, nil
	}
	return json.Marshal(jia)
}

func (jia *JSONInt64Array) Scan(value any) error {
	return JSONScan(value, jia)
}

var jsonNull = []byte{'n', 'u', 'l', 'l'}

func JSONValue(value any) (driver.Value, error) {
	if ref.IsNil(value) {
		return nil, nil
	}
	return json.Marshal(value)
}

func JSONScan(value, dest any) error {
	if value == nil {
		return json.Unmarshal(jsonNull, dest)
	}

	b, ok := value.([]byte)
	if !ok {
		return ErrNotBytes
	}
	return json.Unmarshal(b, dest)
}
