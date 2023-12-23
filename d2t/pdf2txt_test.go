package d2t

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestExtractTextFromPdfFile(t *testing.T) {
	cs := []struct {
		f string
		w string
	}{
		{"hello.pdf", "                      Hello world!\n1.   EN:\n     How do you do?\n2.   JA:\n     今日はお元気ですか？\n3.   ZH:\n     今天         好吗？\n"},
		{"table.pdf", "                      Table\n\n\n\nThis is a cat. That is a dog. This is a fox.\n\n1-1            1-2               1-3\n\n\n\nThis    is   a This     is    an That is pear.\nbird.          apple.            2-3\n\n2-1            2-2\n"},
	}

	for i, c := range cs {
		bw := &bytes.Buffer{}
		err := ExtractTextFromPdfFile("testdata/"+c.f, bw)
		if err != nil {
			fmt.Printf("[%d] ExtractTextFromPdfFile(%s): %v\n", i, c.f, err)
		} else {
			a := bw.String()
			if c.w != a {
				t.Errorf("[%d] ExtractTextFromPdfFile(%s):\nACTUAL: %q\n  WANT: %q\n", i, c.f, a, c.w)
			}
		}
	}
}

func TestExtractTextFromPdfReader(t *testing.T) {
	cs := []struct {
		f string
		w string
	}{
		{"hello.pdf", "                      Hello world!\n1.   EN:\n     How do you do?\n2.   JA:\n     今日はお元気ですか？\n3.   ZH:\n     今天         好吗？\n"},
		{"table.pdf", "                      Table\n\n\n\nThis is a cat. That is a dog. This is a fox.\n\n1-1            1-2               1-3\n\n\n\nThis    is   a This     is    an That is pear.\nbird.          apple.            2-3\n\n2-1            2-2\n"},
	}

	for i, c := range cs {
		fr, err := os.Open("testdata/" + c.f)
		if err != nil {
			t.Errorf("[%d] TestExtractTextFromPdfReader(%s): %v\n", i, c.f, err)
			continue
		}
		defer fr.Close()

		bw := &bytes.Buffer{}
		err = ExtractTextFromPdfReader(fr, bw)
		if err != nil {
			fmt.Printf("[%d] ExtractTextFromPdfFile(%s): %v\n", i, c.f, err)
			continue
		}

		a := bw.String()
		if c.w != a {
			t.Errorf("[%d] ExtractTextFromPdfFile(%s):\nACTUAL: %q\n  WANT: %q\n", i, c.f, a, c.w)
		}
	}
}
