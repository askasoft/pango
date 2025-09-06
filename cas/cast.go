package cas

import (
	"fmt"
	"strconv"
	"time"

	"github.com/askasoft/pango/tmu"
)

func ToDuration(v any) (d time.Duration, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case time.Duration:
		return o, nil
	case *time.Duration:
		return *o, nil
	case string:
		if o == "" {
			return
		}
		d, err = tmu.ParseDuration(o)
	case int8:
		d = time.Duration(o)
	case int16:
		d = time.Duration(o)
	case int32:
		d = time.Duration(o)
	case int64:
		d = time.Duration(o)
	case int:
		d = time.Duration(o)
	case uint8:
		d = time.Duration(o)
	case uint16:
		d = time.Duration(o)
	case uint32:
		d = time.Duration(o)
	case uint64:
		d = time.Duration(o)
	case uint:
		d = time.Duration(o)
	case float32:
		d = time.Duration(o)
	case float64:
		d = time.Duration(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to time.Duration", v)
	}
	return
}

func utcMilli(msec int64) time.Time {
	return time.Unix(msec/1e3, (msec%1e3)*1e6).UTC()
}

func ToTime(v any) (t time.Time, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case time.Time:
		return o, nil
	case *time.Time:
		return *o, nil
	case string:
		if o == "" {
			return
		}
		t, err = tmu.ParseInLocation(o, time.Local)
	case int8:
		t = utcMilli(int64(o))
	case int16:
		t = utcMilli(int64(o))
	case int32:
		t = utcMilli(int64(o))
	case int64:
		t = utcMilli(int64(o))
	case int:
		t = utcMilli(int64(o))
	case uint8:
		t = utcMilli(int64(o))
	case uint16:
		t = utcMilli(int64(o))
	case uint32:
		t = utcMilli(int64(o))
	case uint64:
		t = utcMilli(int64(o))
	case uint:
		t = utcMilli(int64(o))
	case float32:
		t = utcMilli(int64(o))
	case float64:
		t = utcMilli(int64(o))
	default:
		err = fmt.Errorf("cannot cast '%T' to time.Time", v)
	}
	return
}

func ToInt(v any) (n int, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseInt(o, 0, strconv.IntSize)
		n, err = int(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = int(o)
	case int16:
		n = int(o)
	case int32:
		n = int(o)
	case int64:
		n = int(o)
	case int:
		n = int(o)
	case uint8:
		n = int(o)
	case uint16:
		n = int(o)
	case uint32:
		n = int(o)
	case uint64:
		n = int(o)
	case uint:
		n = int(o)
	case float32:
		n = int(o)
	case float64:
		n = int(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to int", v)
	}
	return
}

func ToInt8(v any) (n int8, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseInt(o, 0, 8)
		n, err = int8(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = int8(o)
	case int16:
		n = int8(o)
	case int32:
		n = int8(o)
	case int64:
		n = int8(o)
	case int:
		n = int8(o)
	case uint8:
		n = int8(o)
	case uint16:
		n = int8(o)
	case uint32:
		n = int8(o)
	case uint64:
		n = int8(o)
	case uint:
		n = int8(o)
	case float32:
		n = int8(o)
	case float64:
		n = int8(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to int8", v)
	}
	return
}

func ToInt16(v any) (n int16, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseInt(o, 0, 16)
		n, err = int16(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = int16(o)
	case int16:
		n = int16(o)
	case int32:
		n = int16(o)
	case int64:
		n = int16(o)
	case int:
		n = int16(o)
	case uint8:
		n = int16(o)
	case uint16:
		n = int16(o)
	case uint32:
		n = int16(o)
	case uint64:
		n = int16(o)
	case uint:
		n = int16(o)
	case float32:
		n = int16(o)
	case float64:
		n = int16(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to int16", v)
	}
	return
}

func ToInt32(v any) (n int32, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseInt(o, 0, 32)
		n, err = int32(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = int32(o)
	case int16:
		n = int32(o)
	case int32:
		n = int32(o)
	case int64:
		n = int32(o)
	case int:
		n = int32(o)
	case uint8:
		n = int32(o)
	case uint16:
		n = int32(o)
	case uint32:
		n = int32(o)
	case uint64:
		n = int32(o)
	case uint:
		n = int32(o)
	case float32:
		n = int32(o)
	case float64:
		n = int32(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to int32", v)
	}
	return
}

func ToInt64(v any) (n int64, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		n, err = strconv.ParseInt(o, 0, 64)
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = int64(o)
	case int16:
		n = int64(o)
	case int32:
		n = int64(o)
	case int64:
		n = int64(o)
	case int:
		n = int64(o)
	case uint8:
		n = int64(o)
	case uint16:
		n = int64(o)
	case uint32:
		n = int64(o)
	case uint64:
		n = int64(o)
	case uint:
		n = int64(o)
	case float32:
		n = int64(o)
	case float64:
		n = int64(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to int64", v)
	}
	return
}

func ToUint(v any) (n uint, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseUint(o, 0, strconv.IntSize)
		n, err = uint(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = uint(o)
	case int16:
		n = uint(o)
	case int32:
		n = uint(o)
	case int64:
		n = uint(o)
	case int:
		n = uint(o)
	case uint8:
		n = uint(o)
	case uint16:
		n = uint(o)
	case uint32:
		n = uint(o)
	case uint64:
		n = uint(o)
	case uint:
		n = uint(o)
	case float32:
		n = uint(o)
	case float64:
		n = uint(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to uint", v)
	}
	return
}

func ToUint8(v any) (n uint8, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseUint(o, 0, 8)
		n, err = uint8(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = uint8(o)
	case int16:
		n = uint8(o)
	case int32:
		n = uint8(o)
	case int64:
		n = uint8(o)
	case int:
		n = uint8(o)
	case uint8:
		n = uint8(o)
	case uint16:
		n = uint8(o)
	case uint32:
		n = uint8(o)
	case uint64:
		n = uint8(o)
	case uint:
		n = uint8(o)
	case float32:
		n = uint8(o)
	case float64:
		n = uint8(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to uint", v)
	}
	return
}

func ToUint16(v any) (n uint16, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseUint(o, 0, 16)
		n, err = uint16(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = uint16(o)
	case int16:
		n = uint16(o)
	case int32:
		n = uint16(o)
	case int64:
		n = uint16(o)
	case int:
		n = uint16(o)
	case uint8:
		n = uint16(o)
	case uint16:
		n = uint16(o)
	case uint32:
		n = uint16(o)
	case uint64:
		n = uint16(o)
	case uint:
		n = uint16(o)
	case float32:
		n = uint16(o)
	case float64:
		n = uint16(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to uint16", v)
	}
	return
}

func ToUint32(v any) (n uint32, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseUint(o, 0, 32)
		n, err = uint32(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = uint32(o)
	case int16:
		n = uint32(o)
	case int32:
		n = uint32(o)
	case int64:
		n = uint32(o)
	case int:
		n = uint32(o)
	case uint8:
		n = uint32(o)
	case uint16:
		n = uint32(o)
	case uint32:
		n = uint32(o)
	case uint64:
		n = uint32(o)
	case uint:
		n = uint32(o)
	case float32:
		n = uint32(o)
	case float64:
		n = uint32(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to uint32", v)
	}
	return
}

func ToUint64(v any) (n uint64, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseUint(o, 0, 64)
		n, err = uint64(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = uint64(o)
	case int16:
		n = uint64(o)
	case int32:
		n = uint64(o)
	case int64:
		n = uint64(o)
	case int:
		n = uint64(o)
	case uint8:
		n = uint64(o)
	case uint16:
		n = uint64(o)
	case uint32:
		n = uint64(o)
	case uint64:
		n = uint64(o)
	case uint:
		n = uint64(o)
	case float32:
		n = uint64(o)
	case float64:
		n = uint64(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to uint64", v)
	}
	return
}

func ToFloat32(v any) (n float32, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		i, e := strconv.ParseFloat(o, 32)
		n, err = float32(i), e
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = float32(o)
	case int16:
		n = float32(o)
	case int32:
		n = float32(o)
	case int64:
		n = float32(o)
	case int:
		n = float32(o)
	case uint8:
		n = float32(o)
	case uint16:
		n = float32(o)
	case uint32:
		n = float32(o)
	case uint64:
		n = float32(o)
	case uint:
		n = float32(o)
	case float32:
		n = float32(o)
	case float64:
		n = float32(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to float32", v)
	}
	return
}

func ToFloat64(v any) (n float64, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		n, err = strconv.ParseFloat(o, 64)
	case bool:
		if o {
			n = 1
		}
	case int8:
		n = float64(o)
	case int16:
		n = float64(o)
	case int32:
		n = float64(o)
	case int64:
		n = float64(o)
	case int:
		n = float64(o)
	case uint8:
		n = float64(o)
	case uint16:
		n = float64(o)
	case uint32:
		n = float64(o)
	case uint64:
		n = float64(o)
	case uint:
		n = float64(o)
	case float32:
		n = float64(o)
	case float64:
		n = float64(o)
	default:
		err = fmt.Errorf("cannot cast '%T' to float64", v)
	}
	return
}

func ToBool(v any) (b bool, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		if o == "" {
			return
		}
		b, err = strconv.ParseBool(o)
	case bool:
		b = o
	case int8:
		b = o != 0
	case int16:
		b = o != 0
	case int32:
		b = o != 0
	case int64:
		b = o != 0
	case int:
		b = o != 0
	case uint8:
		b = o != 0
	case uint16:
		b = o != 0
	case uint32:
		b = o != 0
	case uint64:
		b = o != 0
	case uint:
		b = o != 0
	case float32:
		b = o != 0
	case float64:
		b = o != 0
	default:
		err = fmt.Errorf("cannot cast '%T' to bool", v)
	}
	return
}

type stringer interface {
	String() string
}

func ToString(v any) (s string, err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case string:
		s = o
	case []byte:
		s = string(o)
	case bool:
		if o {
			s = "true"
		} else {
			s = "false"
		}
	case int8:
		s = strconv.FormatInt(int64(o), 10)
	case int16:
		s = strconv.FormatInt(int64(o), 10)
	case int32:
		s = strconv.FormatInt(int64(o), 10)
	case int64:
		s = strconv.FormatInt(int64(o), 10)
	case int:
		s = strconv.FormatInt(int64(o), 10)
	case uint8:
		s = strconv.FormatUint(uint64(o), 10)
	case uint16:
		s = strconv.FormatUint(uint64(o), 10)
	case uint32:
		s = strconv.FormatUint(uint64(o), 10)
	case uint64:
		s = strconv.FormatUint(uint64(o), 10)
	case uint:
		s = strconv.FormatUint(uint64(o), 10)
	case float32:
		s = strconv.FormatFloat(float64(o), 'f', -1, 32)
	case float64:
		s = strconv.FormatFloat(o, 'f', -1, 64)
	default:
		if sr, ok := v.(stringer); ok {
			s = sr.String()
		} else {
			err = fmt.Errorf("cannot cast '%T' to string", v)
		}
	}
	return
}
