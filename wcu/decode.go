package wcu

import (
	"fmt"
	"io"

	xcharset "golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

// DecodeBytes Converts the bytes with specified charset to UTF-8.
func DecodeBytes(src []byte, chatset string) ([]byte, error) {
	enc, _ := xcharset.Lookup(chatset)
	if enc == nil {
		return src, fmt.Errorf("unsupported charset %q", chatset)
	}

	des, _, err := transform.Bytes(enc.NewDecoder(), src)
	return des, err
}

// DecodeBytes Detect the chatset of the bytes and convert it to UTF-8.
func DetectAndDecodeBytes(src []byte, html ...bool) ([]byte, error) {
	cs, err := DetectCharsetBytes(src, html...)
	if err != nil {
		return src, err
	}

	return DecodeBytes(src, cs)
}

// Transform Transform the (transformed reader, encoding name, errorr) with the specified charset.
func Transform(r io.Reader, charset string) (io.Reader, string, error) {
	enc, name := xcharset.Lookup(charset)
	if enc == nil {
		return r, "", fmt.Errorf("unsupported charset %q", charset)
	}

	return transform.NewReader(r, enc.NewDecoder()), name, nil
}

// DetectAndTransform Detect the chatset from the reader and transform the reader with the detected xcharset.
func DetectAndTransform(r io.Reader, detect int, html ...bool) (io.Reader, string, error) {
	dr, cs, err := DetectCharsetReader(r, detect, html...)
	if err != nil {
		return dr, cs, err
	}

	return Transform(dr, cs)
}
