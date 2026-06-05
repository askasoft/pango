package sqx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
)

type JSONObject[T any] map[string]T

func (jo JSONObject[T]) String() string {
	return jsonx.Prettify(jo)
}

func (jo JSONObject[T]) Value() (driver.Value, error) {
	if jo == nil {
		return nil, nil
	}
	return json.Marshal(jo)
}

func (jo *JSONObject[T]) Scan(value any) error {
	return JSONScan(value, jo)
}

type (
	JSONAnyObject    = JSONObject[any]
	JSONStringObject = JSONObject[string]
)

type JSONArray[T any] []T

func (ja JSONArray[T]) Value() (driver.Value, error) {
	if ja == nil {
		return nil, nil
	}
	return json.Marshal(ja)
}

func (ja *JSONArray[T]) Scan(value any) error {
	return JSONScan(value, ja)
}

func (ja JSONArray[T]) String() string {
	return jsonx.Prettify(ja)
}

type (
	JSONAnyArray    = JSONArray[any]
	JSONStringArray = JSONArray[string]
	JSONIntArray    = JSONArray[int]
	JSONInt64Array  = JSONArray[int64]
)

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

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, dest)
	case string:
		return json.Unmarshal(str.UnsafeBytes(v), dest)
	default:
		return fmt.Errorf("sqx: cannot scan %T to JSON", value)
	}
}
