package wcu

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path"
	"testing"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

func DetectEncodingReader(r io.Reader, contentType string) (io.Reader, encoding.Encoding, string, error) {
	br := bufio.NewReaderSize(r, 1024)

	data, err := br.Peek(1024)
	if err != nil && !errors.Is(err, io.EOF) {
		return br, nil, "", err
	}

	enc, name, _ := charset.DetermineEncoding(data, contentType)
	return br, enc, name, nil
}

func FailedTestDetectEncodingReader(t *testing.T) {
	for i, c := range cs {
		sf := path.Join("testdata", c+".txt")
		fr, err := os.Open(sf)
		if err != nil {
			t.Fatalf("[%d] Failed to read file %q: %v", i, sf, err)
		}
		defer fr.Close()

		_, _, a, err := DetectEncodingReader(fr, "")
		if err != nil {
			t.Fatalf("[%d] Failed to DetectEncodingReader(%q): %v", i, sf, err)
		}

		if a != c {
			t.Errorf("[%d] DetectEncodingReader(%q) = %v, want %v", i, c, a, c)
		}
	}
}
