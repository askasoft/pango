package cas

import (
	"fmt"
	"strconv"
	"time"
)

func ToDuration(v any) (time.Duration, error) {
	if v == nil {
		return 0, nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return 0, nil
		}
		return time.ParseDuration(s)
	case int8:
		return time.Duration(s), nil
	case int16:
		return time.Duration(s), nil
	case int32:
		return time.Duration(s), nil
	case int64:
		return time.Duration(s), nil
	case int:
		return time.Duration(s), nil
	case uint8:
		return time.Duration(s), nil
	case uint16:
		return time.Duration(s), nil
	case uint32:
		return time.Duration(s), nil
	case uint64:
		return time.Duration(s), nil
	case uint:
		return time.Duration(s), nil
	case float32:
		return time.Duration(s), nil
	case float64:
		return time.Duration(s), nil
	}
	return 0, fmt.Errorf("cannot cast '%v' to time.Duration", v)
}

func utcMilli(msec int64) time.Time {
	return time.Unix(msec/1e3, (msec%1e3)*1e6).UTC()
}

func ToTime(v any) (time.Time, error) {
	if v == nil {
		return time.Time{}, nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return time.Time{}, nil
		}
		return time.Parse(time.RFC3339, s)
	case int8:
		return utcMilli(int64(s)), nil
	case int16:
		return utcMilli(int64(s)), nil
	case int32:
		return utcMilli(int64(s)), nil
	case int64:
		return utcMilli(int64(s)), nil
	case int:
		return utcMilli(int64(s)), nil
	case uint8:
		return utcMilli(int64(s)), nil
	case uint16:
		return utcMilli(int64(s)), nil
	case uint32:
		return utcMilli(int64(s)), nil
	case uint64:
		return utcMilli(int64(s)), nil
	case uint:
		return utcMilli(int64(s)), nil
	case float32:
		return utcMilli(int64(s)), nil
	case float64:
		return utcMilli(int64(s)), nil
	}
	return time.Time{}, fmt.Errorf("cannot cast '%v' to time.Time", v)
}

func ToInt(v any) (int, error) {
	if v == nil {
		return int(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return int(0), nil
		}
		i, err := strconv.ParseInt(s, 0, strconv.IntSize)
		return int(i), err
	case bool:
		if s {
			return int(1), nil
		}
		return int(0), nil
	case int8:
		return int(s), nil
	case int16:
		return int(s), nil
	case int32:
		return int(s), nil
	case int64:
		return int(s), nil
	case int:
		return int(s), nil
	case uint8:
		return int(s), nil
	case uint16:
		return int(s), nil
	case uint32:
		return int(s), nil
	case uint64:
		return int(s), nil
	case uint:
		return int(s), nil
	case float32:
		return int(s), nil
	case float64:
		return int(s), nil
	}
	return 0, fmt.Errorf("cannot cast '%v' to int", v)
}

func ToInt8(v any) (int8, error) {
	if v == nil {
		return 0, nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return 0, nil
		}
		i, err := strconv.ParseInt(s, 0, 8)
		return int8(i), err
	case bool:
		if s {
			return int8(1), nil
		}
		return 0, nil
	case int8:
		return int8(s), nil
	case int16:
		return int8(s), nil
	case int32:
		return int8(s), nil
	case int64:
		return int8(s), nil
	case int:
		return int8(s), nil
	case uint8:
		return int8(s), nil
	case uint16:
		return int8(s), nil
	case uint32:
		return int8(s), nil
	case uint64:
		return int8(s), nil
	case uint:
		return int8(s), nil
	case float32:
		return int8(s), nil
	case float64:
		return int8(s), nil
	}
	return 0, fmt.Errorf("cannot cast '%v' to int8", v)
}

func ToInt16(v any) (int16, error) {
	if v == nil {
		return int16(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return int16(0), nil
		}
		i, err := strconv.ParseInt(s, 0, 16)
		return int16(i), err
	case bool:
		if s {
			return int16(1), nil
		}
		return int16(0), nil
	case int8:
		return int16(s), nil
	case int16:
		return int16(s), nil
	case int32:
		return int16(s), nil
	case int64:
		return int16(s), nil
	case int:
		return int16(s), nil
	case uint8:
		return int16(s), nil
	case uint16:
		return int16(s), nil
	case uint32:
		return int16(s), nil
	case uint64:
		return int16(s), nil
	case uint:
		return int16(s), nil
	case float32:
		return int16(s), nil
	case float64:
		return int16(s), nil
	}
	return 0, fmt.Errorf("cannot cast '%v' to int16", v)
}

func ToInt32(v any) (int32, error) {
	if v == nil {
		return int32(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return int32(0), nil
		}
		i, err := strconv.ParseInt(s, 0, 32)
		return int32(i), err
	case bool:
		if s {
			return int32(1), nil
		}
		return int32(0), nil
	case int8:
		return int32(s), nil
	case int16:
		return int32(s), nil
	case int32:
		return int32(s), nil
	case int64:
		return int32(s), nil
	case int:
		return int32(s), nil
	case uint8:
		return int32(s), nil
	case uint16:
		return int32(s), nil
	case uint32:
		return int32(s), nil
	case uint64:
		return int32(s), nil
	case uint:
		return int32(s), nil
	case float32:
		return int32(s), nil
	case float64:
		return int32(s), nil
	}
	return 0, fmt.Errorf("cannot cast '%v' to int32", v)
}

func ToInt64(v any) (int64, error) {
	if v == nil {
		return int64(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return int64(0), nil
		}
		i, err := strconv.ParseInt(s, 0, 64)
		return int64(i), err
	case bool:
		if s {
			return int64(1), nil
		}
		return int64(0), nil
	case int8:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int32:
		return int64(s), nil
	case int64:
		return int64(s), nil
	case int:
		return int64(s), nil
	case uint8:
		return int64(s), nil
	case uint16:
		return int64(s), nil
	case uint32:
		return int64(s), nil
	case uint64:
		return int64(s), nil
	case uint:
		return int64(s), nil
	case float32:
		return int64(s), nil
	case float64:
		return int64(s), nil
	}
	return int64(0), fmt.Errorf("cannot cast '%v' to int64", v)
}

func ToUint(v any) (uint, error) {
	if v == nil {
		return uint(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return uint(0), nil
		}
		i, err := strconv.ParseUint(s, 0, strconv.IntSize)
		return uint(i), err
	case bool:
		if s {
			return uint(1), nil
		}
		return uint(0), nil
	case int8:
		return uint(s), nil
	case int16:
		return uint(s), nil
	case int32:
		return uint(s), nil
	case int64:
		return uint(s), nil
	case int:
		return uint(s), nil
	case uint8:
		return uint(s), nil
	case uint16:
		return uint(s), nil
	case uint32:
		return uint(s), nil
	case uint64:
		return uint(s), nil
	case uint:
		return uint(s), nil
	case float32:
		return uint(s), nil
	case float64:
		return uint(s), nil
	}
	return uint(0), fmt.Errorf("cannot cast '%v' to uint", v)
}

func ToUint8(v any) (uint8, error) {
	if v == nil {
		return uint8(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return uint8(0), nil
		}
		i, err := strconv.ParseUint(s, 0, 8)
		return uint8(i), err
	case bool:
		if s {
			return uint8(1), nil
		}
		return uint8(0), nil
	case int8:
		return uint8(s), nil
	case int16:
		return uint8(s), nil
	case int32:
		return uint8(s), nil
	case int64:
		return uint8(s), nil
	case int:
		return uint8(s), nil
	case uint8:
		return uint8(s), nil
	case uint16:
		return uint8(s), nil
	case uint32:
		return uint8(s), nil
	case uint64:
		return uint8(s), nil
	case uint:
		return uint8(s), nil
	case float32:
		return uint8(s), nil
	case float64:
		return uint8(s), nil
	}
	return uint8(0), fmt.Errorf("cannot cast '%v' to uint", v)
}

func ToUint16(v any) (uint16, error) {
	if v == nil {
		return uint16(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return uint16(0), nil
		}
		i, err := strconv.ParseUint(s, 0, 16)
		return uint16(i), err
	case bool:
		if s {
			return uint16(1), nil
		}
		return uint16(0), nil
	case int8:
		return uint16(s), nil
	case int16:
		return uint16(s), nil
	case int32:
		return uint16(s), nil
	case int64:
		return uint16(s), nil
	case int:
		return uint16(s), nil
	case uint8:
		return uint16(s), nil
	case uint16:
		return uint16(s), nil
	case uint32:
		return uint16(s), nil
	case uint64:
		return uint16(s), nil
	case uint:
		return uint16(s), nil
	case float32:
		return uint16(s), nil
	case float64:
		return uint16(s), nil
	}
	return uint16(0), fmt.Errorf("cannot cast '%v' to uint16", v)
}

func ToUint32(v any) (uint32, error) {
	if v == nil {
		return uint32(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return uint32(0), nil
		}
		i, err := strconv.ParseUint(s, 0, 32)
		return uint32(i), err
	case bool:
		if s {
			return uint32(1), nil
		}
		return uint32(0), nil
	case int8:
		return uint32(s), nil
	case int16:
		return uint32(s), nil
	case int32:
		return uint32(s), nil
	case int64:
		return uint32(s), nil
	case int:
		return uint32(s), nil
	case uint8:
		return uint32(s), nil
	case uint16:
		return uint32(s), nil
	case uint32:
		return uint32(s), nil
	case uint64:
		return uint32(s), nil
	case uint:
		return uint32(s), nil
	case float32:
		return uint32(s), nil
	case float64:
		return uint32(s), nil
	}
	return uint32(0), fmt.Errorf("cannot cast '%v' to uint32", v)
}

func ToUint64(v any) (uint64, error) {
	if v == nil {
		return uint64(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return uint64(0), nil
		}
		i, err := strconv.ParseUint(s, 0, 64)
		return uint64(i), err
	case bool:
		if s {
			return uint64(1), nil
		}
		return uint64(0), nil
	case int8:
		return uint64(s), nil
	case int16:
		return uint64(s), nil
	case int32:
		return uint64(s), nil
	case int64:
		return uint64(s), nil
	case int:
		return uint64(s), nil
	case uint8:
		return uint64(s), nil
	case uint16:
		return uint64(s), nil
	case uint32:
		return uint64(s), nil
	case uint64:
		return uint64(s), nil
	case uint:
		return uint64(s), nil
	case float32:
		return uint64(s), nil
	case float64:
		return uint64(s), nil
	}
	return uint64(0), fmt.Errorf("cannot cast '%v' to uint64", v)
}

func ToFloat32(v any) (float32, error) {
	if v == nil {
		return float32(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return float32(0), nil
		}
		i, err := strconv.ParseFloat(s, 32)
		return float32(i), err
	case bool:
		if s {
			return float32(1), nil
		}
		return float32(0), nil
	case int8:
		return float32(s), nil
	case int16:
		return float32(s), nil
	case int32:
		return float32(s), nil
	case int64:
		return float32(s), nil
	case int:
		return float32(s), nil
	case uint8:
		return float32(s), nil
	case uint16:
		return float32(s), nil
	case uint32:
		return float32(s), nil
	case uint64:
		return float32(s), nil
	case uint:
		return float32(s), nil
	case float32:
		return float32(s), nil
	case float64:
		return float32(s), nil
	}
	return float32(0), fmt.Errorf("cannot cast '%v' to float32", v)
}

func ToFloat64(v any) (float64, error) {
	if v == nil {
		return float64(0), nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return float64(0), nil
		}
		i, err := strconv.ParseFloat(s, 64)
		return float64(i), err
	case bool:
		if s {
			return float64(1), nil
		}
		return float64(0), nil
	case int8:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int64:
		return float64(s), nil
	case int:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint64:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case float32:
		return float64(s), nil
	case float64:
		return float64(s), nil
	}
	return float64(0), fmt.Errorf("cannot cast '%v' to float64", v)
}

func ToBool(v any) (bool, error) {
	if v == nil {
		return false, nil
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return false, nil
		}
		return strconv.ParseBool(s)
	case bool:
		return s, nil
	case int8:
		return s != 0, nil
	case int16:
		return s != 0, nil
	case int32:
		return s != 0, nil
	case int64:
		return s != 0, nil
	case int:
		return s != 0, nil
	case uint8:
		return s != 0, nil
	case uint16:
		return s != 0, nil
	case uint32:
		return s != 0, nil
	case uint64:
		return s != 0, nil
	case uint:
		return s != 0, nil
	case float32:
		return s != 0, nil
	case float64:
		return s != 0, nil
	}
	return false, fmt.Errorf("cannot cast '%v' to bool", v)
}

func ToString(v any) (string, error) {
	if v == nil {
		return "", nil
	}

	switch s := v.(type) {
	case string:
		return s, nil
	case bool:
		if s {
			return "true", nil
		}
		return "false", nil
	case time.Duration:
		return s.String(), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int32:
		return strconv.FormatInt(int64(s), 10), nil
	case int64:
		return strconv.FormatInt(int64(s), 10), nil
	case int:
		return strconv.FormatInt(int64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	}
	return "", fmt.Errorf("cannot cast '%v' to string", v)
}
