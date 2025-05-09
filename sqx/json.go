package sqx

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var ErrNotBytes = errors.New("sqx: type assertion to []byte failed")

type JSONObject map[string]any

func (jo JSONObject) Value() (driver.Value, error) {
	return json.Marshal(jo)
}

func (jo *JSONObject) Scan(value any) error {
	return JSONScan(value, jo)
}

type JSONArray []any

func (ja JSONArray) Value() (driver.Value, error) {
	return json.Marshal(ja)
}

func (ja *JSONArray) Scan(value any) error {
	return JSONScan(value, ja)
}

type JSONStringArray []string

func (jsa JSONStringArray) Value() (driver.Value, error) {
	return json.Marshal(jsa)
}

func (jsa *JSONStringArray) Scan(value any) error {
	return JSONScan(value, jsa)
}

type JSONIntArray []int

func (jia JSONIntArray) Value() (driver.Value, error) {
	return json.Marshal(jia)
}

func (jia *JSONIntArray) Scan(value any) error {
	return JSONScan(value, jia)
}

type JSONInt64Array []int64

func (jia JSONInt64Array) Value() (driver.Value, error) {
	return json.Marshal(jia)
}

func (jia *JSONInt64Array) Scan(value any) error {
	return JSONScan(value, jia)
}

var jsonNull = []byte{'n', 'u', 'l', 'l'}

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
