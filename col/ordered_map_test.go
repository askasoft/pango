package col

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMapBasicFeatures(t *testing.T) {
	n := 100
	om := NewOrderedMap()

	// set(i, 2 * i)
	for i := 0; i < n; i++ {
		ov, ok := interface{}(nil), false

		assertLenEqual(t, om, i)
		if i%2 == 0 {
			ov, ok = om.Set(i, 2*i)
		} else {
			ov, ok = om.SetIfAbsent(i, 2*i)
		}
		assertLenEqual(t, om, i+1)

		assert.Nil(t, ov)
		assert.False(t, ok)

		ov, ok = om.SetIfAbsent(i, 3*i)
		assert.Equal(t, 2*i, ov)
		assert.True(t, ok)
	}

	// get what we just set
	for i := 0; i < n; i++ {
		value, ok := om.Get(i)

		assert.Equal(t, 2*i, value)
		assert.True(t, ok)
	}

	// get items of what we just set
	for i := 0; i < n; i++ {
		item := om.Item(i)

		assert.NotNil(t, item)
		assert.Equal(t, 2*i, item.Value)
	}

	// keys
	ks := make([]interface{}, n)
	for i := 0; i < n; i++ {
		ks[i] = i
	}
	assert.Equal(t, ks, om.Keys())

	// items
	mis := om.Items()
	assert.Equal(t, n, len(mis))
	for i := 0; i < n; i++ {
		assert.Equal(t, i, mis[i].Key())
		assert.Equal(t, i*2, mis[i].Value)
	}

	// values
	vs := make([]interface{}, n)
	for i := 0; i < n; i++ {
		vs[i] = i * 2
	}
	assert.Equal(t, vs, om.Values())

	// forward iteration
	i := 0
	for mi := om.Front(); mi != nil; mi = mi.Next() {
		assert.Equal(t, i, mi.key)
		assert.Equal(t, 2*i, mi.Value)
		i++
	}

	// backward iteration
	i = n - 1
	for mi := om.Back(); mi != nil; mi = mi.Prev() {
		assert.Equal(t, i, mi.key)
		assert.Equal(t, 2*i, mi.Value)
		i--
	}

	// forward iteration starting from known key
	i = 42
	for mi := om.Item(i); mi != nil; mi = mi.Next() {
		assert.Equal(t, i, mi.key)
		assert.Equal(t, 2*i, mi.Value)
		i++
	}

	// double values for items with even keys
	for j := 0; j < n/2; j++ {
		i = 2 * j
		old, ok := om.Set(i, 4*i)

		assert.Equal(t, 2*i, old)
		assert.True(t, ok)
	}

	// and delete itmes with odd keys
	for j := 0; j < n/2; j++ {
		i = 2*j + 1
		assertLenEqual(t, om, n-j)
		value, ok := om.Delete(i)
		assertLenEqual(t, om, n-j-1)

		assert.Equal(t, 2*i, value)
		assert.True(t, ok)

		// deleting again shouldn't change anything
		value, ok = om.Delete(i)
		assertLenEqual(t, om, n-j-1)
		assert.Nil(t, value)
		assert.False(t, ok)
	}

	// get the whole range
	for j := 0; j < n/2; j++ {
		i = 2 * j
		value, ok := om.Get(i)
		assert.Equal(t, 4*i, value)
		assert.True(t, ok)

		i = 2*j + 1
		value, ok = om.Get(i)
		assert.Nil(t, value)
		assert.False(t, ok)
	}

	// check iterations again
	i = 0
	for mi := om.Front(); mi != nil; mi = mi.Next() {
		assert.Equal(t, i, mi.key)
		assert.Equal(t, 4*i, mi.Value)
		i += 2
	}
	i = 2 * ((n - 1) / 2)
	for mi := om.Back(); mi != nil; mi = mi.Prev() {
		assert.Equal(t, i, mi.key)
		assert.Equal(t, 4*i, mi.Value)
		i -= 2
	}
}

func TestOrderedMapUpdatingDoesntChangePairsOrder(t *testing.T) {
	om := NewOrderedMap("foo", "bar", 12, 28, 78, 100, "bar", "baz")

	old, ok := om.Set(78, 102)
	assert.Equal(t, 100, old)
	assert.True(t, ok)

	assertOrderedPairsEqual(t, om,
		[]interface{}{"foo", 12, 78, "bar"},
		[]interface{}{"bar", 28, 102, "baz"})
}

func TestOrderedMapDeletingAndReinsertingChangesPairsOrder(t *testing.T) {
	om := NewOrderedMap()
	om.Set("foo", "bar")
	om.Set(12, 28)
	om.Set(78, 100)
	om.Set("bar", "baz")

	// delete a item
	old, ok := om.Delete(78)
	assert.Equal(t, 100, old)
	assert.True(t, ok)

	// re-insert the same item
	old, ok = om.Set(78, 100)
	assert.Nil(t, old)
	assert.False(t, ok)

	assertOrderedPairsEqual(t, om,
		[]interface{}{"foo", 12, "bar", 78},
		[]interface{}{"bar", 28, "baz", 100})
}

func TestOrderedMapEmptyMapOperations(t *testing.T) {
	om := NewOrderedMap()

	old, ok := om.Get("foo")
	assert.Nil(t, old)
	assert.False(t, ok)

	old, ok = om.Delete("bar")
	assert.Nil(t, old)
	assert.False(t, ok)

	assertLenEqual(t, om, 0)

	assert.Nil(t, om.Front())
	assert.Nil(t, om.Back())
}

type dummyTestStruct struct {
	value string
}

func TestOrderedMapPackUnpackStructs(t *testing.T) {
	om := NewOrderedMap()
	om.Set("foo", dummyTestStruct{"foo!"})
	om.Set("bar", dummyTestStruct{"bar!"})

	value, ok := om.Get("foo")
	assert.True(t, ok)
	if assert.NotNil(t, value) {
		assert.Equal(t, "foo!", value.(dummyTestStruct).value)
	}

	value, ok = om.Set("bar", dummyTestStruct{"baz!"})
	assert.True(t, ok)
	if assert.NotNil(t, value) {
		assert.Equal(t, "bar!", value.(dummyTestStruct).value)
	}

	value, ok = om.Get("bar")
	assert.True(t, ok)
	if assert.NotNil(t, value) {
		assert.Equal(t, "baz!", value.(dummyTestStruct).value)
	}
}

func TestOrderedMapShuffle(t *testing.T) {
	ranLen := 100

	for _, n := range []int{0, 10, 20, 100, 1000, 10000} {
		t.Run(fmt.Sprintf("shuffle test with %d items", n), func(t *testing.T) {
			om := NewOrderedMap()

			keys := make([]interface{}, n)
			values := make([]interface{}, n)

			for i := 0; i < n; i++ {
				// we prefix with the number to ensure that we don't get any duplicates
				keys[i] = fmt.Sprintf("%d_%s", i, randomHexString(t, ranLen))
				values[i] = randomHexString(t, ranLen)

				value, ok := om.Set(keys[i], values[i])
				assert.Nil(t, value)
				assert.False(t, ok)
			}

			assertOrderedPairsEqual(t, om, keys, values)
		})
	}
}

func TestOrderedMapTemplateRange(t *testing.T) {
	om := NewOrderedMap("z", "Z", "a", "A")
	tmpl, err := template.New("test").Parse("{{range $e := .om.Items}}[ {{$e.Key}} = {{$e.Value}} ]{{end}}")
	if err != nil {
		assert.Fail(t, err.Error())
	}

	cm := map[string]interface{}{
		"om": om,
	}
	sb := &strings.Builder{}
	err = tmpl.Execute(sb, cm)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, "[ z = Z ][ a = A ]", sb.String())
}

/* Test helpers */
func assertOrderedPairsEqual(t *testing.T, om *OrderedMap, expectedKeys, expectedValues []interface{}) {
	assertOrderedPairsEqualFromNewest(t, om, expectedKeys, expectedValues)
	assertOrderedPairsEqualFromOldest(t, om, expectedKeys, expectedValues)
}

func assertOrderedPairsEqualFromNewest(t *testing.T, om *OrderedMap, expectedKeys, expectedValues []interface{}) {
	if assert.Equal(t, len(expectedKeys), len(expectedValues)) && assert.Equal(t, len(expectedKeys), om.Len()) {
		i := om.Len() - 1
		for item := om.Back(); item != nil; item = item.Prev() {
			assert.Equal(t, expectedKeys[i], item.key)
			assert.Equal(t, expectedValues[i], item.Value)
			i--
		}
	}
}

func assertOrderedPairsEqualFromOldest(t *testing.T, om *OrderedMap, expectedKeys, expectedValues []interface{}) {
	if assert.Equal(t, len(expectedKeys), len(expectedValues)) && assert.Equal(t, len(expectedKeys), om.Len()) {
		i := om.Len() - 1
		for item := om.Back(); item != nil; item = item.Prev() {
			assert.Equal(t, expectedKeys[i], item.key)
			assert.Equal(t, expectedValues[i], item.Value)
			i--
		}
	}
}

func assertLenEqual(t *testing.T, om *OrderedMap, expectedLen int) {
	assert.Equal(t, expectedLen, om.Len())

	// also check the list length, for good measure
	assert.Equal(t, expectedLen, om.list.Len())
}

func randomHexString(t *testing.T, length int) string {
	b := length / 2
	randBytes := make([]byte, b)

	if n, err := rand.Read(randBytes); err != nil || n != b {
		if err == nil {
			err = fmt.Errorf("only got %v random bytes, expected %v", n, b)
		}
		t.Fatal(err)
	}

	return hex.EncodeToString(randBytes)
}

func TestOrderedMapString(t *testing.T) {
	assert.Equal(t, `{"1":1,"3":3,"2":2}`, fmt.Sprintf("%s", NewOrderedMap("1", 1, "3", 3, "2", 2)))
}

/*----------- JOSN Test -----------------*/
func TestOrderedMapMarshal(t *testing.T) {
	om := NewOrderedMap()
	om.Set("a", 34)
	om.Set("b", []int{3, 4, 5})
	b, err := json.Marshal(om)
	if err != nil {
		t.Fatalf("Marshal OrderedMap: %v", err)
	}
	// fmt.Printf("%q\n", b)
	const expected = "{\"a\":34,\"b\":[3,4,5]}"
	if !bytes.Equal(b, []byte(expected)) {
		t.Errorf("Marshal OrderedMap: %q not equal to expected %q", b, expected)
	}
}

func ExampleOrderedMap_UnmarshalJSON() {
	const jsonStream = `{
  "country"     : "United States",
  "countryCode" : "US",
  "region"      : "CA",
  "regionName"  : "California",
  "city"        : "Mountain View",
  "zip"         : "94043",
  "lat"         : 37.4192,
  "lon"         : -122.0574,
  "timezone"    : "America/Los_Angeles",
  "isp"         : "Google Cloud",
  "org"         : "Google Cloud",
  "as"          : "AS15169 Google Inc.",
  "mobile"      : true,
  "proxy"       : false,
  "query"       : "35.192.xx.xxx"
}`

	// compare with if using a regular generic map, the unmarshalled result
	//  is a map with unpredictable order of keys
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonStream), &m)
	if err != nil {
		fmt.Println("error:", err)
	}
	for key := range m {
		// fmt.Printf("%-12s: %v\n", key, m[key])
		_ = key
	}

	// use the OrderedMap to Unmarshal from JSON object
	om := NewOrderedMap()
	err = json.Unmarshal([]byte(jsonStream), om)
	if err != nil {
		fmt.Println("error:", err)
	}

	// loop over all key-value pairs,
	// it is ok to call Set append-modify new key-value pairs,
	// but not safe to call Delete during iteration.
	for me := om.Front(); me != nil; me = me.Next() {
		fmt.Printf("%-12s: %v\n", me.Key(), me.Value)
		if me.Key() == "city" {
			om.Set("mobile", false)
			om.Set("extra", 42)
		}
	}

	// Output:
	// country     : United States
	// countryCode : US
	// region      : CA
	// regionName  : California
	// city        : Mountain View
	// zip         : 94043
	// lat         : 37.4192
	// lon         : -122.0574
	// timezone    : America/Los_Angeles
	// isp         : Google Cloud
	// org         : Google Cloud
	// as          : AS15169 Google Inc.
	// mobile      : false
	// proxy       : false
	// query       : 35.192.xx.xxx
	// extra       : 42
}

func TestOrderedMapUnmarshalFromInvalid(t *testing.T) {
	om := NewOrderedMap()

	om.Set("m", math.NaN())
	b, err := json.Marshal(om)
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", b, err)
	}
	// fmt.Println(om, b, err)
	om.Delete("m")

	err = json.Unmarshal([]byte("[]"), om)
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error")
	}

	err = json.Unmarshal([]byte("["), om)
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte(nil))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte("{}3"))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte("{"))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte("{]"))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte(`{"a": 3, "b": [{`))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte(`{"a": 3, "b": [}`))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}
	// fmt.Println("error:", om, err)
}

func TestOrderedMapUnmarshal(t *testing.T) {
	var (
		data  = []byte(`{"as":"AS15169 Google Inc.","city":"Mountain View","country":"United States","countryCode":"US","isp":"Google Cloud","lat":37.4192,"lon":-122.0574,"org":"Google Cloud","query":"35.192.25.53","region":"CA","regionName":"California","status":"success","timezone":"America/Los_Angeles","zip":"94043"}`)
		pairs = []interface{}{
			"as", "AS15169 Google Inc.",
			"city", "Mountain View",
			"country", "United States",
			"countryCode", "US",
			"isp", "Google Cloud",
			"lat", 37.4192,
			"lon", -122.0574,
			"org", "Google Cloud",
			"query", "35.192.25.53",
			"region", "CA",
			"regionName", "California",
			"status", "success",
			"timezone", "America/Los_Angeles",
			"zip", "94043",
		}
		obj = NewOrderedMap(pairs...)
	)

	om := NewOrderedMap()
	err := json.Unmarshal(data, om)
	if err != nil {
		t.Fatalf("Unmarshal OrderedMap: %v", err)
	}

	// check by Has and GetValue
	for i := 0; i+1 < len(pairs); i += 2 {
		k := pairs[i]
		v := pairs[i+1]

		if !om.Has(k) {
			t.Fatalf("expect key %q exists in Unmarshaled OrderedMap", k)
		}
		value, ok := om.Get(k)
		if !ok || value != v {
			t.Fatalf("expect for key %q: the value %v should equal to %v, in Unmarshaled OrderedMap", k, value, v)
		}
	}

	b, err := json.MarshalIndent(om, "", "  ")
	if err != nil {
		t.Fatalf("Unmarshal OrderedMap: %v", err)
	}
	const expected = `{
  "as": "AS15169 Google Inc.",
  "city": "Mountain View",
  "country": "United States",
  "countryCode": "US",
  "isp": "Google Cloud",
  "lat": 37.4192,
  "lon": -122.0574,
  "org": "Google Cloud",
  "query": "35.192.25.53",
  "region": "CA",
  "regionName": "California",
  "status": "success",
  "timezone": "America/Los_Angeles",
  "zip": "94043"
}`
	if !bytes.Equal(b, []byte(expected)) {
		t.Fatalf("Unmarshal OrderedMap marshal indent from %#v not equal to expected: %q\n", om, expected)
	}

	if !reflect.DeepEqual(om, obj) {
		t.Fatalf("Unmarshal OrderedMap not deeply equal: %#v %#v", om, obj)
	}

	val, ok := om.Delete("org")
	if !ok {
		t.Fatalf("org should exist")
	}
	om.Set("org", val)
	b, err = json.MarshalIndent(om, "", "  ")
	// fmt.Println("after delete", om, string(b), err)
	if err != nil {
		t.Fatalf("Unmarshal OrderedMap: %v", err)
	}
	const expected2 = `{
  "as": "AS15169 Google Inc.",
  "city": "Mountain View",
  "country": "United States",
  "countryCode": "US",
  "isp": "Google Cloud",
  "lat": 37.4192,
  "lon": -122.0574,
  "query": "35.192.25.53",
  "region": "CA",
  "regionName": "California",
  "status": "success",
  "timezone": "America/Los_Angeles",
  "zip": "94043",
  "org": "Google Cloud"
}`
	if !bytes.Equal(b, []byte(expected2)) {
		t.Fatalf("Unmarshal OrderedMap marshal indent from %#v not equal to expected: %s\n", om, expected2)
	}
}

func TestOrderedMapUnmarshalNested(t *testing.T) {
	var (
		data = []byte(`{"a": true, "b": [3, 4, { "b": "3", "d": [] }]}`)
		obj  = NewOrderedMap(
			"a", true,
			"b", JSONArray{float64(3), float64(4), NewOrderedMap("b", "3", "d", JSONArray{})},
		)
	)

	om := NewOrderedMap()
	err := json.Unmarshal(data, om)
	if err != nil {
		t.Fatalf("Unmarshal OrderedMap: %v", err)
	}

	if !reflect.DeepEqual(om, obj) {
		t.Fatalf("Unmarshal OrderedMap not deeply equal: %#v expected %#v", om, obj)
	}
}

func ExampleNewOrderedMap() {
	// initialize from a list of key-value pairs
	om := NewOrderedMap(
		"country", "United States",
		"countryCode", "US",
		"region", "CA",
		"regionName", "California",
		"city", "Mountain View",
		"zip", "94043",
		"lat", 37.4192,
		"lon", -122.0574,
		"timezone", "America/Los_Angeles",
		"isp", "Google Cloud",
		"org", "Google Cloud",
		"as", "AS15169 Google Inc.",
		"mobile", true,
		"proxy", false,
		"query", "35.192.xx.xxx",
	)

	for me := om.Back(); me != nil; me = me.Prev() {
		fmt.Printf("%-12s: %v\n", me.Key(), me.Value)
	}

	// Output:
	// query       : 35.192.xx.xxx
	// proxy       : false
	// mobile      : true
	// as          : AS15169 Google Inc.
	// org         : Google Cloud
	// isp         : Google Cloud
	// timezone    : America/Los_Angeles
	// lon         : -122.0574
	// lat         : 37.4192
	// zip         : 94043
	// city        : Mountain View
	// regionName  : California
	// region      : CA
	// countryCode : US
	// country     : United States
}

var unmarshalTests = []struct {
	in  string
	out interface{}
	err interface{}
}{
	{in: "{}", out: NewOrderedMap()},
	{in: `{"a": 3}`, out: NewOrderedMap("a", float64(3))},
	{in: `{"a": 3, "b": true}`, out: NewOrderedMap("a", float64(3), "b", true)},
	{in: `{"a": 3, "b": true, "c": null}`, out: NewOrderedMap("a", float64(3), "b", true, "c", nil)},
	{in: `{"a": 3, "c": null, "d": []}`, out: NewOrderedMap("a", float64(3), "c", nil, "d", JSONArray{})},
	{in: `{"a": 3, "c": null, "d": [3,4,true]}`, out: NewOrderedMap(
		"a", float64(3), "c", nil, "d", JSONArray{
			float64(3), float64(4), true,
		})},
	{in: `{"a": 3, "c": null, "d": [3,4,true, { "inner": "abc" }]}`, out: NewOrderedMap(
		"a", float64(3), "c", nil, "d", JSONArray([]interface{}{
			float64(3), float64(4), true, NewOrderedMap("inner", "abc"),
		}))},
}

func TestOrderedMapUnmarshals(t *testing.T) {
	for i, tt := range unmarshalTests {
		in := []byte(tt.in)

		v := NewOrderedMap()
		dec := json.NewDecoder(bytes.NewReader(in))
		if err := dec.Decode(v); !reflect.DeepEqual(err, tt.err) {
			t.Errorf("#%d: %v, want %v", i, err, tt.err)
			continue
		} else if err != nil {
			continue
		}

		if !reflect.DeepEqual(v, tt.out) {
			act, _ := json.Marshal(v)
			exp, _ := json.Marshal(tt.out)
			t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v\nact: %s\nexp:%s", i, v, tt.out, string(act), string(exp))
			continue
		}

		// Check round trip also decodes correctly.
		if tt.err == nil {
			enc, err := json.Marshal(v)
			if err != nil {
				t.Errorf("#%d: error re-marshaling: %v", i, err)
				continue
			}

			vv := NewOrderedMap()
			dec = json.NewDecoder(bytes.NewReader(enc))
			if err := dec.Decode(vv); err != nil {
				t.Errorf("#%d: error re-unmarshaling %#q: %v", i, enc, err)
				continue
			}
			if !reflect.DeepEqual(v, vv) {
				t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, v, vv)
				t.Errorf("     In: %q", strings.Map(noSpace, string(in)))
				t.Errorf("Marshal: %q", strings.Map(noSpace, string(enc)))
				continue
			}
		}
	}
}

func noSpace(c rune) rune {
	if isSpace(byte(c)) { //only used for ascii
		return -1
	}
	return c
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
