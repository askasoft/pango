package cas

import (
	"testing"
	"time"
)

func TestToDuration(t *testing.T) {
	cs := []struct {
		s any
		w time.Duration
	}{
		{"", 0},
		{"1s", time.Second},
		{int8(100), time.Duration(100)},
		{int16(200), time.Duration(200)},
		{int32(300), time.Duration(300)},
		{int64(400), time.Duration(400)},
		{int(50000), time.Duration(50000)},
		{uint8(100), time.Duration(100)},
		{uint16(20), time.Duration(20)},
		{uint32(30), time.Duration(30)},
		{uint64(40), time.Duration(40)},
		{uint(5000), time.Duration(5000)},
		{float32(10), time.Duration(10)},
		{float64(20), time.Duration(20)},
	}

	for i, c := range cs {
		a, err := ToDuration(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToDuration(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToTime(t *testing.T) {
	tm0102, err := time.ParseInLocation("2006-1-2 15:04:05", "2000-01-02 00:00:00", time.Local)
	if err != nil {
		t.Fatal(err)
	}
	tm_1_2, err := time.ParseInLocation("2006-1-2 15:04:05", "2000-1-2 00:00:00", time.Local)
	if err != nil {
		t.Fatal(err)
	}

	cs := []struct {
		s any
		w time.Time
	}{
		{"", time.Time{}},
		{"1970-01-01T00:00:01Z", time.Unix(1, 0).UTC()},
		{"2000-01-02 00:00:00", tm0102},
		{"2000-1-2 00:00:00", tm_1_2},
		{int8(100), time.Unix(0, 100000000).UTC()},
		{int16(200), time.Unix(0, 200000000).UTC()},
		{int32(300), time.Unix(0, 300000000).UTC()},
		{int64(400), time.Unix(0, 400000000).UTC()},
		{int(50000), time.Unix(50, 0).UTC()},
		{uint8(100), time.Unix(0, 100000000).UTC()},
		{uint16(20), time.Unix(0, 20000000).UTC()},
		{uint32(30), time.Unix(0, 30000000).UTC()},
		{uint64(40), time.Unix(0, 40000000).UTC()},
		{uint(5000), time.Unix(5, 0).UTC()},
		{float32(10000), time.Unix(10, 0).UTC()},
		{float64(20000), time.Unix(20, 0).UTC()},
	}

	for i, c := range cs {
		a, err := ToTime(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToTime(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToInt(t *testing.T) {
	cs := []struct {
		s any
		w int
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), int(100)},
		{int16(200), int(200)},
		{int32(300), int(300)},
		{int64(400), int(400)},
		{int(50000), int(50000)},
		{uint8(100), int(100)},
		{uint16(20), int(20)},
		{uint32(30), int(30)},
		{uint64(40), int(40)},
		{uint(5000), int(5000)},
		{float32(1), int(1)},
		{float64(2), int(2)},
	}

	for i, c := range cs {
		a, err := ToInt(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToInt(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToInt8(t *testing.T) {
	cs := []struct {
		s any
		w int8
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), int8(100)},
		{int16(20), int8(20)},
		{int32(30), int8(30)},
		{int64(40), int8(40)},
		{int(50), int8(50)},
		{uint8(100), int8(100)},
		{uint16(20), int8(20)},
		{uint32(30), int8(30)},
		{uint64(40), int8(40)},
		{uint(50), int8(50)},
		{float32(1), int8(1)},
		{float64(2), int8(2)},
	}

	for i, c := range cs {
		a, err := ToInt8(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToInt8(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToInt16(t *testing.T) {
	cs := []struct {
		s any
		w int16
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), int16(100)},
		{int16(20), int16(20)},
		{int32(30), int16(30)},
		{int64(40), int16(40)},
		{int(5000), int16(5000)},
		{uint8(100), int16(100)},
		{uint16(20), int16(20)},
		{uint32(30), int16(30)},
		{uint64(40), int16(40)},
		{uint(5000), int16(5000)},
		{float32(1), int16(1)},
		{float64(2), int16(2)},
	}

	for i, c := range cs {
		a, err := ToInt16(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToInt16(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToInt32(t *testing.T) {
	cs := []struct {
		s any
		w int32
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), int32(100)},
		{int16(20), int32(20)},
		{int32(30), int32(30)},
		{int64(40), int32(40)},
		{int(5000), int32(5000)},
		{uint8(100), int32(100)},
		{uint16(20), int32(20)},
		{uint32(30), int32(30)},
		{uint64(40), int32(40)},
		{uint(5000), int32(5000)},
		{float32(1), int32(1)},
		{float64(2), int32(2)},
	}

	for i, c := range cs {
		a, err := ToInt32(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToInt32(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToInt64(t *testing.T) {
	cs := []struct {
		s any
		w int64
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), int64(100)},
		{int16(20), int64(20)},
		{int32(30), int64(30)},
		{int64(40), int64(40)},
		{int(5000), int64(5000)},
		{uint8(100), int64(100)},
		{uint16(20), int64(20)},
		{uint32(30), int64(30)},
		{uint64(40), int64(40)},
		{uint(5000), int64(5000)},
		{float32(1), int64(1)},
		{float64(2), int64(2)},
	}

	for i, c := range cs {
		a, err := ToInt64(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToInt64(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToUint(t *testing.T) {
	cs := []struct {
		s any
		w uint
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), uint(100)},
		{int16(200), uint(200)},
		{int32(300), uint(300)},
		{int64(400), uint(400)},
		{int(50000), uint(50000)},
		{uint8(100), uint(100)},
		{uint16(20), uint(20)},
		{uint32(30), uint(30)},
		{uint64(40), uint(40)},
		{uint(5000), uint(5000)},
		{float32(1), uint(1)},
		{float64(2), uint(2)},
	}

	for i, c := range cs {
		a, err := ToUint(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToUint(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToUint8(t *testing.T) {
	cs := []struct {
		s any
		w uint8
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), uint8(100)},
		{int16(20), uint8(20)},
		{int32(30), uint8(30)},
		{int64(40), uint8(40)},
		{int(50), uint8(50)},
		{uint8(100), uint8(100)},
		{uint16(20), uint8(20)},
		{uint32(30), uint8(30)},
		{uint64(40), uint8(40)},
		{uint(50), uint8(50)},
		{float32(1), uint8(1)},
		{float64(2), uint8(2)},
	}

	for i, c := range cs {
		a, err := ToUint8(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToUint8(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToUint16(t *testing.T) {
	cs := []struct {
		s any
		w uint16
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), uint16(100)},
		{int16(20), uint16(20)},
		{int32(30), uint16(30)},
		{int64(40), uint16(40)},
		{int(5000), uint16(5000)},
		{uint8(100), uint16(100)},
		{uint16(20), uint16(20)},
		{uint32(30), uint16(30)},
		{uint64(40), uint16(40)},
		{uint(5000), uint16(5000)},
		{float32(1), uint16(1)},
		{float64(2), uint16(2)},
	}

	for i, c := range cs {
		a, err := ToUint16(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToUint16(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToUint32(t *testing.T) {
	cs := []struct {
		s any
		w uint32
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), uint32(100)},
		{int16(20), uint32(20)},
		{int32(30), uint32(30)},
		{int64(40), uint32(40)},
		{int(5000), uint32(5000)},
		{uint8(100), uint32(100)},
		{uint16(20), uint32(20)},
		{uint32(30), uint32(30)},
		{uint64(40), uint32(40)},
		{uint(5000), uint32(5000)},
		{float32(1), uint32(1)},
		{float64(2), uint32(2)},
	}

	for i, c := range cs {
		a, err := ToUint32(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToUint32(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToUint64(t *testing.T) {
	cs := []struct {
		s any
		w uint64
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), uint64(100)},
		{int16(20), uint64(20)},
		{int32(30), uint64(30)},
		{int64(40), uint64(40)},
		{int(5000), uint64(5000)},
		{uint8(100), uint64(100)},
		{uint16(20), uint64(20)},
		{uint32(30), uint64(30)},
		{uint64(40), uint64(40)},
		{uint(5000), uint64(5000)},
		{float32(1), uint64(1)},
		{float64(2), uint64(2)},
	}

	for i, c := range cs {
		a, err := ToUint64(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToUint64(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToFloat32(t *testing.T) {
	cs := []struct {
		s any
		w float32
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), float32(100)},
		{int16(20), float32(20)},
		{int32(30), float32(30)},
		{int64(40), float32(40)},
		{int(5000), float32(5000)},
		{uint8(100), float32(100)},
		{uint16(20), float32(20)},
		{uint32(30), float32(30)},
		{uint64(40), float32(40)},
		{uint(5000), float32(5000)},
		{float32(1), float32(1)},
		{float64(2), float32(2)},
	}

	for i, c := range cs {
		a, err := ToFloat32(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToFloat32(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToFloat64(t *testing.T) {
	cs := []struct {
		s any
		w float64
	}{
		{nil, 0},
		{"", 0},
		{"1", 1},
		{int8(100), float64(100)},
		{int16(20), float64(20)},
		{int32(30), float64(30)},
		{int64(40), float64(40)},
		{int(5000), float64(5000)},
		{uint8(100), float64(100)},
		{uint16(20), float64(20)},
		{uint32(30), float64(30)},
		{uint64(40), float64(40)},
		{uint(5000), float64(5000)},
		{float32(1), float64(1)},
		{float64(2), float64(2)},
	}

	for i, c := range cs {
		a, err := ToFloat64(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToFloat64(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToBool(t *testing.T) {
	cs := []struct {
		s any
		w bool
	}{
		{nil, false},
		{"", false},
		{"1", true},
		{int8(100), true},
		{int16(20), true},
		{int32(30), true},
		{int64(40), true},
		{int(5000), true},
		{uint8(100), true},
		{uint16(20), true},
		{uint32(30), true},
		{uint64(40), true},
		{uint(5000), true},
		{float32(1), true},
		{float64(2), true},
	}

	for i, c := range cs {
		a, err := ToBool(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToBool(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}

func TestToString(t *testing.T) {
	cs := []struct {
		s any
		w string
	}{
		{nil, ""},
		{"", ""},
		{"1", "1"},
		{[]byte{'b', 's'}, "bs"},
		{true, "true"},
		{time.Second, "1s"},
		{int8(100), "100"},
		{int16(20), "20"},
		{int32(30), "30"},
		{int64(40), "40"},
		{int(5000), "5000"},
		{uint8(100), "100"},
		{uint16(20), "20"},
		{uint32(30), "30"},
		{uint64(40), "40"},
		{uint(5000), "5000"},
		{float32(1), "1"},
		{float64(2), "2"},
	}

	for i, c := range cs {
		a, err := ToString(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ToString(%v) = (%v, %v), want: %v", i, c.s, a, err, c.w)
		}
	}
}
