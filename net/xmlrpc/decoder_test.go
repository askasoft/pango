package xmlrpc

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type book struct {
	Title  string
	Amount int
}

type bookUnexported struct {
	title  string
	amount int
}

var decodeTests = []struct {
	val any
	ptr any
	xml string
}{
	// nil
	{"", new(*string), "<value><nil/></value>"},

	// int, i4, i8
	{0, new(*int), "<value><int></int></value>"},
	{100, new(*int), "<value><int>100</int></value>"},
	{389451, new(*int), "<value><i4>389451</i4></value>"},
	{int64(45659074), new(*int64), "<value><i8>45659074</i8></value>"},

	// string
	{"Once upon a time", new(*string), "<value><string>Once upon a time</string></value>"},
	{"Mike & Mick <London, UK>", new(*string), "<value><string>Mike &amp; Mick &lt;London, UK&gt;</string></value>"},
	{"Once upon a time", new(*string), "<value>Once upon a time</value>"},

	// base64
	{[]byte("0123456789"), new(*[]byte), "<value><base64>MDEyMzQ1Njc4OQ==</base64></value>"},

	// boolean
	{true, new(*bool), "<value><boolean>1</boolean></value>"},
	{false, new(*bool), "<value><boolean>0</boolean></value>"},

	// double
	{12.134, new(*float32), "<value><double>12.134</double></value>"},
	{-12.134, new(*float32), "<value><double>-12.134</double></value>"},

	// datetime.iso8601
	{_time("2013-12-09T21:00:12Z"), new(*time.Time), "<value><dateTime.iso8601>20131209T21:00:12</dateTime.iso8601></value>"},
	{_time("2013-12-09T21:00:12Z"), new(*time.Time), "<value><dateTime.iso8601>20131209T21:00:12Z</dateTime.iso8601></value>"},
	{_time("2013-12-09T21:00:12-01:00"), new(*time.Time), "<value><dateTime.iso8601>20131209T21:00:12-01:00</dateTime.iso8601></value>"},
	{_time("2013-12-09T21:00:12+01:00"), new(*time.Time), "<value><dateTime.iso8601>20131209T21:00:12+01:00</dateTime.iso8601></value>"},
	{_time("2013-12-09T21:00:12Z"), new(*time.Time), "<value><dateTime.iso8601>2013-12-09T21:00:12</dateTime.iso8601></value>"},
	{_time("2013-12-09T21:00:12Z"), new(*time.Time), "<value><dateTime.iso8601>2013-12-09T21:00:12Z</dateTime.iso8601></value>"},
	{_time("2013-12-09T21:00:12-01:00"), new(*time.Time), "<value><dateTime.iso8601>2013-12-09T21:00:12-01:00</dateTime.iso8601></value>"},
	{_time("2013-12-09T21:00:12+01:00"), new(*time.Time), "<value><dateTime.iso8601>2013-12-09T21:00:12+01:00</dateTime.iso8601></value>"},

	// array
	{[]int{1, 5, 7}, new(*[]int), "<value><array><data><value><int>1</int></value><value><int>5</int></value><value><int>7</int></value></data></array></value>"},
	{[]any{"A", "5"}, new(any), "<value><array><data><value><string>A</string></value><value><string>5</string></value></data></array></value>"},
	{[]any{"A", int64(5)}, new(any), "<value><array><data><value><string>A</string></value><value><int>5</int></value></data></array></value>"},

	// struct
	{book{"War and Piece", 20}, new(*book), "<value><struct><member><name>Title</name><value><string>War and Piece</string></value></member><member><name>Amount</name><value><int>20</int></value></member></struct></value>"},
	{bookUnexported{title: "", amount: 0}, new(*bookUnexported), "<value><struct><member><name>title</name><value><string>War and Piece</string></value></member><member><name>amount</name><value><int>20</int></value></member></struct></value>"},
	{map[string]any{"Name": "John Smith"}, new(any), "<value><struct><member><name>Name</name><value><string>John Smith</string></value></member></struct></value>"},
	{map[string]any{}, new(any), "<value><struct></struct></value>"},
}

func _time(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(fmt.Sprintf("time parsing error: %v", err))
	}
	return t
}

func TestDecode(t *testing.T) {
	for i, tt := range decodeTests {
		v := reflect.New(reflect.TypeOf(tt.val))
		if err := DecodeString(tt.xml, v.Interface()); err != nil {
			t.Fatalf("#%d unmarshal error: %v", i, err)
		}

		v = v.Elem()

		if v.Kind() == reflect.Slice {
			vv := reflect.ValueOf(tt.val)
			if vv.Len() != v.Len() {
				t.Fatalf("#%d unmarshal error:\nexpected: %v\n     got: %v", i, tt.val, v.Interface())
			}
			for i := 0; i < v.Len(); i++ {
				if v.Index(i).Interface() != vv.Index(i).Interface() {
					t.Fatalf("#%d unmarshal error:\nexpected: %v\n     got: %v", i, tt.val, v.Interface())
				}
			}
		} else {
			a1 := v.Interface()
			a2 := any(tt.val)

			if !reflect.DeepEqual(a1, a2) {
				t.Fatalf("#%d unmarshal error:\nexpected: %v\n     got: %v", i, tt.val, v.Interface())
			}
		}
	}
}

func TestDecodeToPtr(t *testing.T) {
	for _, tt := range decodeTests {
		if err := DecodeString(tt.xml, tt.ptr); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
	}
}

func TestDecodeToNil(t *testing.T) {
	encoded := "<value><int>100</int></value>"
	err := DecodeString(encoded, nil)
	fmt.Println(err)
}

func TestTypeMismatchError(t *testing.T) {
	var s string

	encoded := "<value><int>100</int></value>"
	var err error

	if err = DecodeString(encoded, &s); err == nil {
		t.Fatal("unmarshal error: expected error, but didn't get it")
	}

	if _, ok := err.(TypeMismatchError); !ok {
		t.Fatal("unmarshal error: expected type mistmatch error, but didn't get it")
	}
}

func TestDecodeEmptyValueTag(t *testing.T) {
	var v int

	if err := DecodeString("<value/>", &v); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
}

const structEmptyXML = `
<value>
  <struct>
  </struct>
</value>
`

func TestDecodeEmptyStruct(t *testing.T) {
	var v any
	if err := DecodeString(structEmptyXML, &v); err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatalf("got nil map")
	}
}

const arrayValueXML = `
<value>
  <array>
    <data>
      <value><int>234</int></value>
      <value><boolean>1</boolean></value>
      <value><string>Hello World</string></value>
      <value><string>Extra Value</string></value>
    </data>
  </array>
</value>
`

func TestDecodeExistingArray(t *testing.T) {
	var (
		v1 int
		v2 bool
		v3 string

		v = []any{&v1, &v2, &v3}
	)
	if err := DecodeString(arrayValueXML, &v); err != nil {
		t.Fatal(err)
	}

	// check pre-existing values
	if want := 234; v1 != want {
		t.Fatalf("want %d, got %d", want, v1)
	}
	if want := true; v2 != want {
		t.Fatalf("want %t, got %t", want, v2)
	}
	if want := "Hello World"; v3 != want {
		t.Fatalf("want %s, got %s", want, v3)
	}
	// check the appended result
	if n := len(v); n != 4 {
		t.Fatalf("missing appended result")
	}
	if got, ok := v[3].(string); !ok || got != "Extra Value" {
		t.Fatalf("got %s, want %s", got, "Extra Value")
	}
}
