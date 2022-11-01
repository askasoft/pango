package xin

import (
	"github.com/pandafw/pango/xin/binding"
)

// EnableJSONDecoderUseNumber sets true for binding.EnableDecoderUseNumber to
// call the UseNumber method on the JSON Decoder instance.
func EnableJSONDecoderUseNumber() {
	binding.EnableDecoderUseNumber = true
}

// EnableJSONDecoderDisallowUnknownFields sets true for binding.EnableDecoderDisallowUnknownFields to
// call the DisallowUnknownFields method on the JSON Decoder instance.
func EnableJSONDecoderDisallowUnknownFields() {
	binding.EnableDecoderDisallowUnknownFields = true
}
