package binding

import (
	"net/http"
	"net/textproto"
	"reflect"
)

type headerBinding struct{}

func (headerBinding) Name() string {
	return "header"
}

func (headerBinding) Bind(req *http.Request, obj any) error {
	return mapHeader(obj, req.Header)
}

func mapHeader(ptr any, h map[string][]string) error {
	return mappingByPtr(ptr, headerSource(h), "header")
}

type headerSource map[string][]string

func (hs headerSource) TrySet(field reflect.Value, key string, opts options) (bool, *FieldBindError) {
	return setByForm(field, hs, textproto.CanonicalMIMEHeaderKey(key), opts)
}
