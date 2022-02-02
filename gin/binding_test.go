package gin

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/pandafw/pango/gin/binding"
	"github.com/pandafw/pango/gin/validate"
	"github.com/stretchr/testify/assert"
)

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required,max=32"`
}

var vd = validate.NewStructValidator()

func TestValidationFails(t *testing.T) {
	var obj FooStruct
	req := requestWithBody("POST", "/", `{"bar": "foo"}`)
	err := binding.JSON.Bind(req, &obj)
	assert.NoError(t, err)
	err = vd.ValidateStruct(&obj)
	assert.Error(t, err)
}

func TestRequiredSucceeds(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"hoge" binding:"required"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"hoge": 0}`)
	err := binding.JSON.Bind(req, &obj)
	assert.NoError(t, err)
	err = vd.ValidateStruct(&obj)
	assert.NoError(t, err)
}

func TestRequiredFails(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"foo" binding:"required"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"boen": 0}`)
	err := binding.JSON.Bind(req, &obj)
	assert.NoError(t, err)
	err = vd.ValidateStruct(&obj)
	assert.Error(t, err)
}
