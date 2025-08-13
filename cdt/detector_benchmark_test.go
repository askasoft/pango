package cdt

import (
	"testing"

	"github.com/askasoft/pango/iox/fsu"
)

var files = []string{
	"8859_1_en.html",
	"8859_1_da.html",
	"8859_1_de.html",
	"8859_1_es.html",
	"8859_1_fr.html",
	"8859_1_pt.html",
	"shift_jis.html",
	"gb18030.html",
	"euc_jp.html",
	"euc_kr.html",
	"big5.html",
	"utf8.html",
}

func benchmarkReadFile(b *testing.B, name string) []byte {
	fn := testFilename(name)
	bs, err := fsu.ReadFile(fn)
	if err != nil {
		b.Fatalf("Failed to read file %q: %v", fn, err)
	}
	return bs
}

func BenchmarkDetectBestConcurrent(b *testing.B) {
	textDetector := NewTextDetector()
	bss := [][]byte{}
	for _, f := range files {
		bss = append(bss, benchmarkReadFile(b, f))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, bs := range bss {
			textDetector.DetectBestConcurrent(bs)
		}
	}
}

func BenchmarkDetectBestSequential(b *testing.B) {
	textDetector := NewTextDetector()

	bss := [][]byte{}
	for _, f := range files {
		bss = append(bss, benchmarkReadFile(b, f))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, bs := range bss {
			textDetector.DetectBestSequential(bs)
		}
	}
}

func (d *Detector) DetectBestConcurrent(b []byte) (*Result, error) {
	input := newRecognizerInput(b, d.stripTag)
	outputChan := make(chan recognizerOutput)
	for _, r := range d.recognizers {
		go matchHelper(r, input, outputChan)
	}

	var output Result
	for i := 0; i < len(d.recognizers); i++ {
		o := <-outputChan
		if output.Confidence < o.Confidence {
			output = Result(o)
		}
	}

	if output.Confidence == 0 {
		return nil, ErrNotDetected
	}
	return &output, nil
}

// DetectBestSync returns the Result with highest Confidence.
func (d *Detector) DetectBestSequential(b []byte) (r *Result, err error) {
	input := newRecognizerInput(b, d.stripTag)

	var best Result
	for _, r := range d.recognizers {
		rout := r.Match(input)
		if rout.Confidence == 100 {
			best = Result(rout)
			return &best, nil
		}
		if best.Confidence < rout.Confidence {
			best = Result(rout)
		}
	}
	if best.Confidence == 0 {
		return nil, ErrNotDetected
	}
	return &best, nil
}
