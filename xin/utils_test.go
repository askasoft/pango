package xin

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/test/assert"
)

func BenchmarkParseAccept(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseAccept("text/html , application/xhtml+xml,application/xml;q=0.9,  */* ;q=0.8")
	}
}

type testStruct struct {
	T *testing.T
}

func (t *testStruct) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	assert.Equal(t.T, "POST", req.Method)
	assert.Equal(t.T, "/path", req.URL.Path)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "hello")
}

func TestWrap(t *testing.T) {
	router := New()
	router.POST("/path", WrapH(&testStruct{t}))
	router.GET("/path2", WrapF(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/path2", req.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "hola!")
	}))

	w := performRequest(router, "POST", "/path")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "hello", w.Body.String())

	w = performRequest(router, "GET", "/path2")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "hola!", w.Body.String())
}

func TestLastChar(t *testing.T) {
	assert.Equal(t, uint8('a'), lastChar("hola"))
	assert.Equal(t, uint8('s'), lastChar("adios"))
	assert.Panics(t, func() { lastChar("") })
}

func TestParseAccept(t *testing.T) {
	parts := parseAccept("text/html , application/xhtml+xml,application/xml;q=0.9,  */* ;q=0.8")
	assert.Len(t, parts, 4)
	assert.Equal(t, "text/html", parts[0])
	assert.Equal(t, "application/xhtml+xml", parts[1])
	assert.Equal(t, "application/xml", parts[2])
	assert.Equal(t, "*/*", parts[3])
}

func TestChooseData(t *testing.T) {
	A := "a"
	B := "b"
	assert.Equal(t, A, chooseData(A, B))
	assert.Equal(t, B, chooseData(nil, B))
	assert.Panics(t, func() { chooseData(nil, nil) })
}

func TestFilterFlags(t *testing.T) {
	result := filterFlags("text/html ")
	assert.Equal(t, "text/html", result)

	result = filterFlags("text/html;")
	assert.Equal(t, "text/html", result)
}

func TestJoinPaths(t *testing.T) {
	assert.Equal(t, "", joinPaths("", ""))
	assert.Equal(t, "/", joinPaths("", "/"))
	assert.Equal(t, "/a", joinPaths("/a", ""))
	assert.Equal(t, "/a/", joinPaths("/a/", ""))
	assert.Equal(t, "/a/", joinPaths("/a/", "/"))
	assert.Equal(t, "/a/", joinPaths("/a", "/"))
	assert.Equal(t, "/a/hola", joinPaths("/a", "/hola"))
	assert.Equal(t, "/a/hola", joinPaths("/a/", "/hola"))
	assert.Equal(t, "/a/hola/", joinPaths("/a/", "/hola/"))
	assert.Equal(t, "/a/hola/", joinPaths("/a/", "/hola//"))
}

func TestMarshalXMLforH(t *testing.T) {
	cs := []struct {
		h H
		w string
	}{
		{H{"x": "X"}, "<map><x>X</x></map>"},
		{H{"": "doc", "x": "X"}, "<doc><x>X</x></doc>"},
	}

	for i, c := range cs {
		bs, err := xml.Marshal(c.h)
		assert.NoError(t, err)

		a := str.UnsafeString(bs)
		assert.Equal(t, c.w, a, fmt.Sprintf("#%d", i))
	}
}
