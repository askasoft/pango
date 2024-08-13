package wcu

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"mime"

	"github.com/askasoft/pango/cdt"
)

var boms = []struct {
	bom []byte
	enc string
}{
	{[]byte{0xfe, 0xff}, "utf-16be"},
	{[]byte{0xff, 0xfe}, "utf-16le"},
	{[]byte{0xef, 0xbb, 0xbf}, "utf-8"},
}

func CharsetFromContentType(contentType string) string {
	if _, params, err := mime.ParseMediaType(contentType); err == nil {
		if cs, ok := params["charset"]; ok {
			return cs
		}
	}
	return ""
}

// DetectCharsetBytes detect charset for the input byte array
func DetectCharsetBytes(data []byte, html ...bool) (string, error) {
	for _, b := range boms {
		if bytes.HasPrefix(data, b.bom) {
			return b.enc, nil
		}
	}

	cd := cdt.NewDetector(html...)
	cr, err := cd.DetectBest(data)
	if err != nil {
		return "", err
	}

	return cr.Charset, nil
}

// DetectCharsetReader detect charset for the input reader
// detect: the maximum bytes to read for detect. if detect < 0, read all `r` into buffer and return bytes.Reader.
func DetectCharsetReader(r io.Reader, detect int, html ...bool) (io.Reader, string, error) {
	if detect < 0 {
		data, err := io.ReadAll(r)
		if err != nil {
			return r, "", err
		}

		cs, err := DetectCharsetBytes(data, html...)
		return bytes.NewReader(data), cs, err
	}

	br := bufio.NewReaderSize(r, detect)

	data, err := br.Peek(detect)
	if err != nil && !errors.Is(err, io.EOF) {
		return br, "", err
	}

	cs, err := DetectCharsetBytes(data, html...)
	return br, cs, err
}
