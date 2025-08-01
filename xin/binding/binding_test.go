package binding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/askasoft/pango/test/assert"
	"github.com/askasoft/pango/test/require"
)

type appkey struct {
	Appkey string `json:"appkey" form:"appkey"`
}

type QueryTest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	appkey
}

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" validate:"required,max=32"`
}

type FooBarStruct struct {
	FooStruct
	Bar string `msgpack:"bar" json:"bar" form:"bar" xml:"bar" validate:"required"`
}

type FooBarFileStruct struct {
	FooBarStruct
	File *multipart.FileHeader `form:"file" validate:"required"`
}

type FooBarFileFailStruct struct {
	FooBarStruct
	File *multipart.FileHeader `invalid_name:"file" validate:"required"`
	// for unexport test
	data *multipart.FileHeader `form:"data" validate:"required"`
}

type FooDefaultBarStruct struct {
	FooStruct
	Bar string `msgpack:"bar" json:"bar" form:"bar,default=hello" xml:"bar" validate:"required"`
}

type FooStructUseNumber struct {
	Foo any `json:"foo" validate:"required"`
}

type FooStructDisallowUnknownFields struct {
	Foo any `json:"foo" validate:"required"`
}

type FooBarStructForTimeType struct {
	TimeFoo   time.Time `form:"time_foo" time_format:"2006-01-02" time_utc:"1" time_location:"Asia/Chongqing"`
	TimeBar   time.Time `form:"time_bar" time_format:"2006-01-02" time_utc:"1"`
	NanoTime  time.Time `form:"nanoTime" time_format:"unixNano"`
	MicroTime time.Time `form:"microTime" time_format:"unixMicro"`
	MilliTime time.Time `form:"milliTime" time_format:"unixMilli"`
	UnixTime  time.Time `form:"unixTime" time_format:"unix"`
}

type FooStructForTimeTypeNotUnixFormat struct {
	NanoTime  time.Time `form:"nanoTime" time_format:"unixNano"`
	MicroTime time.Time `form:"microTime" time_format:"unixMicro"`
	MilliTime time.Time `form:"milliTime" time_format:"unixMilli"`
	UnixTime  time.Time `form:"unixTime" time_format:"unix"`
}

type FooStructForTimeTypeNotFormat struct {
	TimeFoo time.Time `form:"time_foo"`
}

type FooStructForTimeTypeFailFormat struct {
	TimeFoo time.Time `form:"time_foo" time_format:"2017-11-15"`
}

type FooStructForTimeTypeFailLocation struct {
	TimeFoo time.Time `form:"time_foo" time_format:"2006-01-02" time_location:"/asia/chongqing"`
}

type FooStructForMapType struct {
	MapFoo map[string]any `form:"map_foo"`
}

type FooStructForIgnoreFormTag struct {
	Foo *string `form:"-"`
}

type InvalidNameType struct {
	TestName string `invalid_name:"test_name"`
}

type InvalidNameMapType struct {
	MapFoo map[string]any `form:"map_foo"`
}

type FooStructForSliceType struct {
	SliceFoo []int `form:"slice_foo"`
}

type StructFooIdx struct {
	Idx int `form:"idx"`
}

type FooStructForStructType struct {
	StructFooIdx
}

type StructFooName struct {
	Name string `form:"name"`
}

type FooStructForStructPointerType struct {
	*StructFooName
}

type FooStructForSliceMapType struct {
	// Unknown type: not support map
	SliceMapFoo []map[string]any `form:"slice_map_foo"`
}

type FooStructForBoolType struct {
	BoolFoo bool `form:"bool_foo"`
}

type FooStructForStringPtrType struct {
	PtrStr *string `form:"ptr_str"`
	PtrNum *int    `form:"ptr_num" validate:"required"`
}

type FooStructForMapPtrType struct {
	PtrBar *map[string]any `form:"ptr_bar"`
}

func TestBindingDefault(t *testing.T) {
	assert.Equal(t, Form, Default("GET", ""))
	assert.Equal(t, Form, Default("GET", MIMEJSON))

	assert.Equal(t, JSON, Default("POST", MIMEJSON))
	assert.Equal(t, JSON, Default("PUT", MIMEJSON))

	assert.Equal(t, XML, Default("POST", MIMEXML))
	assert.Equal(t, XML, Default("PUT", MIMEXML2))

	assert.Equal(t, Form, Default("POST", MIMEPOSTForm))
	assert.Equal(t, Form, Default("PUT", MIMEPOSTForm))

	assert.Equal(t, FormMultipart, Default("POST", MIMEMultipartPOSTForm))
	assert.Equal(t, FormMultipart, Default("PUT", MIMEMultipartPOSTForm))
}

func TestBindingJSON(t *testing.T) {
	testBodyBinding(t,
		JSON, "json",
		"/", "/",
		`{"foo": "bar"}`, `{"foo": 1}`)
}

func TestBindingJSONSlice(t *testing.T) {
	EnableDecoderDisallowUnknownFields = true
	defer func() {
		EnableDecoderDisallowUnknownFields = false
	}()

	testBodyBindingSlice(t, JSON, "json", "/", "/", `[]`, ``)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"foo": 123}]`)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"bar": 123}]`)
}

func TestBindingJSONUseNumber(t *testing.T) {
	testBodyBindingUseNumber(t,
		JSON, "json",
		"/",
		`{"foo": 123}`)
}

func TestBindingJSONUseNumber2(t *testing.T) {
	testBodyBindingUseNumber2(t,
		JSON, "json",
		"/",
		`{"foo": 123}`)
}

func TestBindingJSONDisallowUnknownFields(t *testing.T) {
	testBodyBindingDisallowUnknownFields(t, JSON,
		"/", "/",
		`{"foo": "bar"}`, `{"foo": "bar", "what": "this"}`)
}

func TestBindingJSONStringMap(t *testing.T) {
	testBodyBindingStringMap(t, JSON,
		"/", "/",
		`{"foo": "bar", "hello": "world"}`, `{"num": 2}`)
}

func TestBindingForm(t *testing.T) {
	testFormBinding(t, "POST",
		"/", "/",
		"foo=bar&bar=foo", "bar2=foo")
}

func TestBindingForm2(t *testing.T) {
	testFormBinding(t, "GET",
		"/?foo=bar&bar=foo", "/?bar2=foo",
		"", "")
}

func TestBindingFormEmbeddedStruct(t *testing.T) {
	testFormBindingEmbeddedStruct(t, "POST",
		"/", "/",
		"page=1&size=2&appkey=test-appkey", "bar2=foo")
}

func TestBindingFormEmbeddedStruct2(t *testing.T) {
	testFormBindingEmbeddedStruct(t, "GET",
		"/?page=1&size=2&appkey=test-appkey", "/?bar2=foo",
		"", "")
}

func TestBindingFormDefaultValue(t *testing.T) {
	testFormBindingDefaultValue(t, "POST",
		"/", "/",
		"foo=bar", "bar2=foo")
}

func TestBindingFormDefaultValue2(t *testing.T) {
	testFormBindingDefaultValue(t, "GET",
		"/?foo=bar", "/?bar2=foo",
		"", "")
}

func TestBindingFormForTimePOST(t *testing.T) {
	testFormBindingForTime(t, "POST", "/", "time_foo=2017-11-15&time_bar=&nanoTime=1562400033123456789&microTime=1562400033123456&milliTime=1562400033123&unixTime=1562400033")
	testFormBindingForTimeNotUnixFormat(t, "POST", "/", "time_foo=2017-11-15&nanoTime=bad")
	testFormBindingForTimeNotUnixFormat(t, "POST", "/", "time_foo=2017-11-15&microTime=bad")
	testFormBindingForTimeNotUnixFormat(t, "POST", "/", "time_foo=2017-11-15&milliTime=bad")
	testFormBindingForTimeNotUnixFormat(t, "POST", "/", "time_foo=2017-11-15&unixTime=bad")
	testFormBindingForTimeNotFormat(t, "POST", "/", "time_foo=2017/11/15")
	testFormBindingForTimeFailFormat(t, "POST", "/", "time_foo=2017-11-15")
	testFormBindingForTimeFailLocation(t, "POST", "/", "time_foo=2017-11-15")
}

func TestBindingFormForTimeGET(t *testing.T) {
	testFormBindingForTime(t, "GET", "/?time_foo=2017-11-15&time_bar=&nanoTime=1562400033123456789&microTime=1562400033123456&milliTime=1562400033123&unixTime=1562400033", "")
	testFormBindingForTimeNotUnixFormat(t, "GET", "/?time_foo=2017-11-15&nanoTime=bad", "")
	testFormBindingForTimeNotUnixFormat(t, "GET", "/?time_foo=2017-11-15&microTime=bad", "")
	testFormBindingForTimeNotUnixFormat(t, "GET", "/?time_foo=2017-11-15&milliTime=bad", "")
	testFormBindingForTimeNotUnixFormat(t, "GET", "/?time_foo=2017-11-15&unixTime=bad", "")
	testFormBindingForTimeNotFormat(t, "GET", "/?time_foo=2017/11/15", "")
	testFormBindingForTimeFailFormat(t, "GET", "/?time_foo=2017-11-15", "")
	testFormBindingForTimeFailLocation(t, "GET", "/?time_foo=2017-11-15", "")
}

func TestFormBindingIgnoreField(t *testing.T) {
	testFormBindingIgnoreField(t, "POST",
		"/", "/",
		"-=bar", "")
}

func TestBindingFormInvalidName(t *testing.T) {
	testFormBindingInvalidName(t, "POST",
		"/", "/",
		"test_name=bar", "bar2=foo")
}

func TestBindingFormInvalidName2(t *testing.T) {
	testFormBindingInvalidName2(t, "POST",
		"/", "/",
		"map_foo=bar", "bar2=foo")
}

func TestBindingFormForType(t *testing.T) {
	testFormBindingForType(t, "POST",
		"/", "/",
		"map_foo={\"bar\":123}", "map_foo=1", "Map")

	testFormBindingForType(t, "POST",
		"/", "/",
		"slice_foo=1&slice_foo=2", "bar2=1&bar2=2", "Slice")

	testFormBindingForType(t, "GET",
		"/?slice_foo=1&slice_foo=2", "/?bar2=1&bar2=2",
		"", "", "Slice")

	testFormBindingForType(t, "POST",
		"/", "/",
		"slice_map_foo=1&slice_map_foo=2", "bar2=1&bar2=2", "SliceMap")

	testFormBindingForType(t, "GET",
		"/?slice_map_foo=1&slice_map_foo=2", "/?bar2=1&bar2=2",
		"", "", "SliceMap")

	testFormBindingForType(t, "POST",
		"/", "/",
		"ptr_str=test", "ptr_num=test", "Ptr")

	testFormBindingForType(t, "GET",
		"/?ptr_str=test", "/?ptr_num=test",
		"", "", "Ptr")

	testFormBindingForType(t, "POST",
		"/", "/",
		"idx=123", "id1=1", "Struct")

	testFormBindingForType(t, "GET",
		"/?idx=123", "/?id1=1",
		"", "", "Struct")

	testFormBindingForType(t, "POST",
		"/", "/",
		"name=thinkerou", "name1=ou", "StructPointer")

	testFormBindingForType(t, "GET",
		"/?name=thinkerou", "/?name1=ou",
		"", "", "StructPointer")
}

func TestBindingFormStringMap(t *testing.T) {
	testBodyBindingStringMap(t, Form,
		"/", "",
		`foo=bar&hello=world`, "")
	// Should pick the last value
	testBodyBindingStringMap(t, Form,
		"/", "",
		`foo=something&foo=bar&hello=world`, "")
}

func TestBindingFormStringSliceMap(t *testing.T) {
	obj := make(map[string][]string)
	req := requestWithBody("POST", "/", "foo=something&foo=bar&hello=world")
	req.Header.Add("Content-Type", MIMEPOSTForm)
	err := Form.Bind(req, &obj)
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	target := map[string][]string{
		"foo":   {"something", "bar"},
		"hello": {"world"},
	}
	assert.True(t, reflect.DeepEqual(obj, target))

	objInvalid := make(map[string][]int)
	req = requestWithBody("POST", "/", "foo=something&foo=bar&hello=world")
	req.Header.Add("Content-Type", MIMEPOSTForm)
	err = Form.Bind(req, &objInvalid)
	assert.Error(t, err)
}

func TestBindingQuery(t *testing.T) {
	testQueryBinding(t, "POST",
		"/?foo=bar&bar=foo", "/",
		"foo=unused", "bar2=foo")
}

func TestBindingQuery2(t *testing.T) {
	testQueryBinding(t, "GET",
		"/?foo=bar&bar=foo", "/?bar2=foo",
		"foo=unused", "")
}

func TestBindingQueryFail(t *testing.T) {
	testQueryBindingFail(t, "POST",
		"/?map_foo=", "/",
		"map_foo=unused", "bar2=foo")
}

func TestBindingQueryFail2(t *testing.T) {
	testQueryBindingFail(t, "GET",
		"/?map_foo=", "/?bar2=foo",
		"map_foo=unused", "")
}

func TestBindingQueryBoolFail(t *testing.T) {
	testQueryBindingBoolFail(t, "GET",
		"/?bool_foo=fasl", "/?bar2=foo",
		"bool_foo=unused", "")
}

func TestBindingQueryStringMap(t *testing.T) {
	b := Query

	obj := make(map[string]string)
	req := requestWithBody("GET", "/?foo=bar&hello=world", "")
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	assert.Equal(t, "bar", obj["foo"])
	assert.Equal(t, "world", obj["hello"])

	obj = make(map[string]string)
	req = requestWithBody("GET", "/?foo=bar&foo=2&hello=world", "") // should pick last
	err = b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	assert.Equal(t, "2", obj["foo"])
	assert.Equal(t, "world", obj["hello"])
}

func TestBindingXML(t *testing.T) {
	testBodyBinding(t,
		XML, "xml",
		"/", "/",
		"<map><foo>bar</foo></map>", "<map><bar>foo</bar></map>")
}

func TestBindingXMLFail(t *testing.T) {
	testBodyBindingFail(t,
		XML, "xml",
		"/", "/",
		"<map><foo>bar<foo></map>", "<map><bar>foo</bar></map>")
}

func createFormPostRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar&bar=foo"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createDefaultFormPostRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestForMap(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo={\"bar\":123}"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestForMap2(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo.a=123&map_foo.b=234"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestForMap3(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo[a]=123&map_foo[b]=234"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestForMapFail(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo=hello"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormFilesMultipartRequest(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("foo", "bar"))
	assert.NoError(t, mw.WriteField("bar", "foo"))

	f, err := os.Open("form.go")
	assert.NoError(t, err)
	defer f.Close()
	fw, err1 := mw.CreateFormFile("file", "form.go")
	assert.NoError(t, err1)
	_, err = io.Copy(fw, f)
	assert.NoError(t, err)

	req, err2 := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	assert.NoError(t, err2)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)

	return req
}

func createFormFilesMultipartRequestFail(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("foo", "bar"))
	assert.NoError(t, mw.WriteField("bar", "foo"))

	f, err := os.Open("form.go")
	assert.NoError(t, err)
	defer f.Close()
	fw, err1 := mw.CreateFormFile("file_foo", "form_foo.go")
	assert.NoError(t, err1)
	_, err = io.Copy(fw, f)
	assert.NoError(t, err)

	req, err2 := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	assert.NoError(t, err2)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)

	return req
}

func createFormMultipartRequest(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("foo", "bar"))
	assert.NoError(t, mw.WriteField("bar", "foo"))
	req, err := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func createFormMultipartRequestForMap(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("map_foo", "{\"bar\":123, \"name\":\"thinkerou\", \"pai\": 3.14}"))
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func createFormMultipartRequestForMapFail(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("map_foo", "3.14"))
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func TestBindingFormPost(t *testing.T) {
	req := createFormPostRequest(t)
	var obj FooBarStruct
	assert.NoError(t, FormPost.Bind(req, &obj))

	assert.Equal(t, "form-urlencoded", FormPost.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingDefaultValueFormPost(t *testing.T) {
	req := createDefaultFormPostRequest(t)
	var obj FooDefaultBarStruct
	assert.NoError(t, FormPost.Bind(req, &obj))

	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "hello", obj.Bar)
}

func TestBindingFormPostForMap(t *testing.T) {
	req := createFormPostRequestForMap(t)
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	assert.NoError(t, err)
	assert.NotNil(t, obj.MapFoo)
	assert.Equal(t, float64(123), obj.MapFoo["bar"])
}

func TestBindingFormPostForMap2(t *testing.T) {
	req := createFormPostRequestForMap2(t)
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	assert.NoError(t, err)
	assert.NotNil(t, obj.MapFoo)
	assert.Equal(t, "123", obj.MapFoo["a"])
	assert.Equal(t, "234", obj.MapFoo["b"])
}

func TestBindingFormPostForMap3(t *testing.T) {
	req := createFormPostRequestForMap3(t)
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	assert.NoError(t, err)
	assert.NotNil(t, obj.MapFoo)
	assert.Equal(t, "123", obj.MapFoo["a"])
	assert.Equal(t, "234", obj.MapFoo["b"])
}

func TestBindingFormPostForMapFail(t *testing.T) {
	req := createFormPostRequestForMapFail(t)
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingFormFilesMultipart(t *testing.T) {
	req := createFormFilesMultipartRequest(t)
	var obj FooBarFileStruct
	err := FormMultipart.Bind(req, &obj)
	assert.NoError(t, err)

	// file from os
	f, _ := os.Open("form.go")
	defer f.Close()
	fileActual, _ := io.ReadAll(f)

	// file from multipart
	mf, _ := obj.File.Open()
	defer mf.Close()
	fileExpect, _ := io.ReadAll(mf)

	assert.Equal(t, FormMultipart.Name(), "multipart/form-data")
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
	assert.Equal(t, fileExpect, fileActual)
}

func TestBindingFormFilesMultipartFail(t *testing.T) {
	req := createFormFilesMultipartRequestFail(t)
	var obj FooBarFileFailStruct
	err := FormMultipart.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Nil(t, obj.File)
	assert.Nil(t, obj.data)
}

func TestBindingFormMultipart(t *testing.T) {
	req := createFormMultipartRequest(t)
	var obj FooBarStruct
	assert.NoError(t, FormMultipart.Bind(req, &obj))

	assert.Equal(t, "multipart/form-data", FormMultipart.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingFormMultipartForMap(t *testing.T) {
	req := createFormMultipartRequestForMap(t)
	var obj FooStructForMapType
	err := FormMultipart.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, float64(123), obj.MapFoo["bar"].(float64))
	assert.Equal(t, "thinkerou", obj.MapFoo["name"].(string))
	assert.Equal(t, float64(3.14), obj.MapFoo["pai"].(float64))
}

func TestBindingFormMultipartForMapFail(t *testing.T) {
	req := createFormMultipartRequestForMapFail(t)
	var obj FooStructForMapType
	err := FormMultipart.Bind(req, &obj)
	assert.Error(t, err)
}

func TestHeaderBinding(t *testing.T) {
	h := Header
	assert.Equal(t, "header", h.Name())

	type tHeader struct {
		Limit int `header:"limit"`
	}

	var theader tHeader
	req := requestWithBody("GET", "/", "")
	req.Header.Add("limit", "1000")
	assert.NoError(t, h.Bind(req, &theader))
	assert.Equal(t, 1000, theader.Limit)

	req = requestWithBody("GET", "/", "")
	req.Header.Add("fail", `{fail:fail}`)

	type failStruct struct {
		Fail map[string]any `header:"fail"`
	}

	err := h.Bind(req, &failStruct{})
	assert.Error(t, err)
}

func TestURIBinding(t *testing.T) {
	b := URI
	assert.Equal(t, "uri", b.Name())

	type Tag struct {
		Name string `uri:"name"`
	}
	var tag Tag
	m := make(map[string][]string)
	m["name"] = []string{"thinkerou"}
	assert.NoError(t, b.BindURI(m, &tag))
	assert.Equal(t, "thinkerou", tag.Name)

	type NotSupportStruct struct {
		Name map[string]any `uri:"name"`
	}
	var not NotSupportStruct
	assert.Error(t, b.BindURI(m, &not))
	assert.Equal(t, map[string]any(nil), not.Name)
}

func TestURIInnerBinding(t *testing.T) {
	type S struct {
		Age int `uri:"age"`
	}

	type Tag struct {
		Name string `uri:"name"`
		S
	}

	expectedName := "mike"
	expectedAge := 25

	m := map[string][]string{
		"name": {expectedName},
		"age":  {strconv.Itoa(expectedAge)},
	}

	var tag Tag
	assert.NoError(t, URI.BindURI(m, &tag))
	assert.Equal(t, tag.Name, expectedName)
	assert.Equal(t, tag.S.Age, expectedAge)
}

func testFormBindingEmbeddedStruct(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := QueryTest{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, 1, obj.Page)
	assert.Equal(t, 2, obj.Size)
	assert.Equal(t, "test-appkey", obj.Appkey)
}

func testFormBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)

	obj = FooBarStruct{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingDefaultValue(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooDefaultBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "hello", obj.Bar)

	obj = FooDefaultBarStruct{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormBindingFail(t *testing.T) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormBindingMultipartFail(t *testing.T) {
	obj := FooBarStruct{}
	req, err := http.NewRequest("POST", "/", strings.NewReader("foo=bar"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+";boundary=testboundary")
	_, err = req.MultipartReader()
	assert.NoError(t, err)
	err = Form.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormPostBindingFail(t *testing.T) {
	b := FormPost
	assert.Equal(t, "form-urlencoded", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormMultipartBindingFail(t *testing.T) {
	b := FormMultipart
	assert.Equal(t, "multipart/form-data", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTime(t *testing.T, method, path, body string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStructForTimeType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)

	assert.NoError(t, err)
	assert.Equal(t, int64(1510675200), obj.TimeFoo.Unix())
	assert.Equal(t, "Asia/Chongqing", obj.TimeFoo.Location().String())
	assert.Equal(t, int64(-62135596800), obj.TimeBar.Unix())
	assert.Equal(t, "UTC", obj.TimeBar.Location().String())
	assert.Equal(t, int64(1562400033123456789), obj.NanoTime.UnixNano())
	assert.Equal(t, int64(1562400033123456), obj.MicroTime.UnixMicro())
	assert.Equal(t, int64(1562400033123), obj.MilliTime.UnixMilli())
	assert.Equal(t, int64(1562400033), obj.UnixTime.Unix())
}

func testFormBindingForTimeNotUnixFormat(t *testing.T, method, path, body string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForTimeTypeNotUnixFormat{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTimeNotFormat(t *testing.T, method, path, body string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForTimeTypeNotFormat{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTimeFailFormat(t *testing.T, method, path, body string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForTimeTypeFailFormat{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTimeFailLocation(t *testing.T, method, path, body string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForTimeTypeFailLocation{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingIgnoreField(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForIgnoreFormTag{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)

	assert.Nil(t, obj.Foo)
}

func testFormBindingInvalidName(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := InvalidNameType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "", obj.TestName)

	obj = InvalidNameType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingInvalidName2(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := InvalidNameMapType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)

	obj = InvalidNameMapType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForType(t *testing.T, method, path, badPath, body, badBody string, typ string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	switch typ {
	case "Slice":
		obj := FooStructForSliceType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2}, obj.SliceFoo)

		obj = FooStructForSliceType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Struct":
		obj := FooStructForStructType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, StructFooIdx{Idx: 123}, obj.StructFooIdx)
	case "StructPointer":
		obj := FooStructForStructPointerType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.NotNil(t, obj.StructFooName)
		assert.Equal(t, StructFooName{Name: "thinkerou"}, *obj.StructFooName)
	case "Map":
		obj := FooStructForMapType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, float64(123), obj.MapFoo["bar"].(float64))
	case "SliceMap":
		obj := FooStructForSliceMapType{}
		err := b.Bind(req, &obj)
		assert.Error(t, err)
	case "Ptr":
		obj := FooStructForStringPtrType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Nil(t, obj.PtrNum)
		assert.Equal(t, "test", *obj.PtrStr)

		obj = FooStructForStringPtrType{}
		req = requestWithBody(method, badPath, badBody)
		if method == "POST" {
			req.Header.Add("Content-Type", MIMEPOSTForm)
		}
		err = b.Bind(req, &obj)
		assert.Error(t, err)
	}
}

func testQueryBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, "query", b.Name())

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func testQueryBindingFail(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, "query", b.Name())

	obj := FooStructForMapType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testQueryBindingBoolFail(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, "query", b.Name())

	obj := FooStructForBoolType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBindingSlice(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	var obj1 []FooStruct
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj1)
	assert.NoError(t, err)

	var obj2 []FooStruct
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj2)
	assert.Error(t, err)
}

func testBodyBindingStringMap(t *testing.T, b Binding, path, badPath, body, badBody string) {
	obj := make(map[string]string)
	req := requestWithBody("POST", path, body)
	if b.Name() == "form" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	assert.Equal(t, "bar", obj["foo"])
	assert.Equal(t, "world", obj["hello"])

	if badPath != "" && badBody != "" {
		obj = make(map[string]string)
		req = requestWithBody("POST", badPath, badBody)
		err = b.Bind(req, &obj)
		assert.Error(t, err)
	}

	objInt := make(map[string]int)
	req = requestWithBody("POST", path, body)
	if b.Name() == "form" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err = b.Bind(req, &objInt)
	assert.Error(t, err)
}

func testBodyBindingUseNumber(t *testing.T, b Binding, name, path, body string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = true
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// we hope it is int64(123)
	v, e := obj.Foo.(json.Number).Int64()
	assert.NoError(t, e)
	assert.Equal(t, int64(123), v)
}

func testBodyBindingUseNumber2(t *testing.T, b Binding, name, path, body string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = false
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// it will return float64(123) if not use EnableDecoderUseNumber
	// maybe it is not hoped
	assert.Equal(t, float64(123), obj.Foo)
}

func testBodyBindingDisallowUnknownFields(t *testing.T, b Binding, path, badPath, body, badBody string) {
	EnableDecoderDisallowUnknownFields = true
	defer func() {
		EnableDecoderDisallowUnknownFields = false
	}()

	obj := FooStructDisallowUnknownFields{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStructDisallowUnknownFields{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "what")
}

func testBodyBindingFail(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
	assert.Equal(t, "", obj.Foo)

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

type failRead struct{}

func (f *failRead) Read(b []byte) (n int, err error) {
	return 0, errors.New("my fail")
}

func (f *failRead) Close() error {
	return nil
}

func TestPlainBinding(t *testing.T) {
	p := Plain
	assert.Equal(t, "plain", p.Name())

	var s string
	req := requestWithBody(http.MethodPost, "/", "test string")
	assert.NoError(t, p.Bind(req, &s))
	assert.Equal(t, "test string", s)

	var bs []byte
	req = requestWithBody(http.MethodPost, "/", "test []byte")
	assert.NoError(t, p.Bind(req, &bs))
	assert.Equal(t, bs, []byte("test []byte"))

	var i int
	req = requestWithBody(http.MethodPost, "/", "test fail")
	require.Error(t, p.Bind(req, &i))

	req = requestWithBody(http.MethodPost, "/", "")
	req.Body = &failRead{}
	require.Error(t, p.Bind(req, &s))

	req = requestWithBody(http.MethodPost, "/", "")
	assert.NoError(t, p.Bind(req, nil))

	var ptr *string
	req = requestWithBody(http.MethodPost, "/", "")
	assert.NoError(t, p.Bind(req, ptr))
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
