package wcu

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/saintfish/chardet"
	xcharset "golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

func NewDetector(html ...bool) (cd *chardet.Detector) {
	if len(html) > 0 && html[0] {
		cd = chardet.NewHtmlDetector()
	} else {
		cd = chardet.NewTextDetector()
	}
	return
}

func DetectCharsetBytes(data []byte, html ...bool) (string, error) {
	cd := NewDetector(html...)
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

	cd := NewDetector(html...)
	cr, err := cd.DetectBest(data)
	if err != nil {
		return br, "", err
	}

	return br, cr.Charset, nil
}

func DetectCharsetFile(filename string) (string, error) {
	fr, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	html := IsHTMLFile(filename)
	_, cs, err := DetectCharsetReader(fr, html)
	return cs, err
}

// DecodeBytes Converts the bytes with specified charset to UTF-8.
func DecodeBytes(src []byte, chatset string) ([]byte, error) {
	enc, _ := xcharset.Lookup(chatset)
	if enc == nil {
		return nil, fmt.Errorf("unsupported charset %q", chatset)
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

// Transform Transform the reader with the specified xcharset.
func Transform(r io.Reader, charset string) (io.Reader, error) {
	enc, _ := xcharset.Lookup(charset)
	if enc == nil {
		return r, fmt.Errorf("unsupported charset %q", charset)
	}

	return transform.NewReader(r, enc.NewDecoder()), nil
}

// DetectAndTransform Detect the chatset from the reader and transform the reader with the detected xcharset.
func DetectAndTransform(r io.Reader, html ...bool) (io.Reader, error) {
	dr, cs, err := DetectCharsetReader(r, html...)
	if err != nil {
		return dr, err
	}

	return Transform(dr, cs)
}
