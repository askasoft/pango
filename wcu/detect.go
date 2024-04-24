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

func DetectCharsetReader(r io.Reader, html ...bool) (io.Reader, string, error) {
	br := bufio.NewReaderSize(r, 1024)

	data, err := br.Peek(1024)
	if err != nil && !errors.Is(err, io.EOF) {
		return br, "", err
	}

	for _, b := range boms {
		if bytes.HasPrefix(data, b.bom) {
			return br, b.enc, nil
		}
	}

	cd := cdt.NewDetector(html...)
	cr, err := cd.DetectBest(data)
	if err != nil {
		return br, "", err
	}

	return br, cr.Charset, nil
}
