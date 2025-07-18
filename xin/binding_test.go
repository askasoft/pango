package xin

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/askasoft/pango/test/assert"
	"github.com/askasoft/pango/xin/binding"
	"github.com/askasoft/pango/xin/validate"
)

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}

type FooStruct struct {
	Foo string `json:"foo" form:"foo" xml:"foo" validate:"required"`
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
		Hoge *int `json:"hoge" validate:"required"`
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
		Hoge *int `json:"foo" validate:"required"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"boen": 0}`)
	err := binding.JSON.Bind(req, &obj)
	assert.NoError(t, err)
	err = vd.ValidateStruct(&obj)
	assert.Error(t, err)
}
