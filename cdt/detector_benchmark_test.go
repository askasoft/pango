package cdt

import (
	"bytes"
	"testing"
)

func BenchmarkDetectBest(b *testing.B) {
	textDetector := NewTextDetector()
	aaaa := bytes.Repeat([]byte("A"), 1024)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		textDetector.DetectBest(aaaa)
	}
}

func BenchmarkDetectBestSync(b *testing.B) {
	textDetector := NewTextDetector()
	aaaa := bytes.Repeat([]byte("A"), 1024)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		textDetector.DetectBestSync(aaaa)
	}
}

// DetectBestSync returns the Result with highest Confidence.
func (d *Detector) DetectBestSync(b []byte) (r *Result, err error) {
	input := newRecognizerInput(b, d.stripTag)

	var best Result
	for _, r := range d.recognizers {
		rout := r.Match(input)
		if best.Confidence < rout.Confidence {
			best = Result(rout)
		}
	}
	if best.Confidence == 0 {
		return nil, ErrNotDetected
	}
	return &best, nil
}
