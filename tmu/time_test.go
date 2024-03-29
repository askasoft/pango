package tmu

import (
	"encoding/json"
	"testing"
)

func TestParseHHMM(t *testing.T) {
	cs := []struct {
		s string
		w HHMM
	}{
		{"09:00", 540},
		{"08:00", 480},
		{"360", 360},
	}

	for i, c := range cs {
		a, err := ParseHHMM(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ParseHHMM(%q) = (%d, %v), want %d", i, c.s, a, err, c.w)
		}
	}
}

func TestHHMMUnmarshallJSON(t *testing.T) {
	cs := []struct {
		js string
		hm HHMM
	}{
		{`{ "s": "09:00" }`, 540},
		{`{ "s": "10:20" }`, 620},
		{`{ "s": "360" }`, 360},
	}

	o := struct {
		S HHMM `json:"s"`
	}{}

	for i, c := range cs {
		err := json.Unmarshal([]byte(c.js), &o)
		if err != nil || o.S != c.hm {
			t.Errorf("[%d] HHMMUnmarshallJSON(%s) = (%d, %v), want %d", i, c.js, o.S, err, c.hm)
		}
	}
}

func TestHHMMMarshallJSON(t *testing.T) {
	cs := []struct {
		hm HHMM
		js string
	}{
		{540, `{"s":"09:00"}`},
		{620, `{"s":"10:20"}`},
		{360, `{"s":"06:00"}`},
	}

	o := struct {
		S HHMM `json:"s,omitempty"`
	}{}

	for i, c := range cs {
		o.S = c.hm
		bs, err := json.Marshal(&o)
		js := string(bs)
		if err != nil || c.js != js {
			t.Errorf("[%d] HHMMMarshallJSON(%d) = (%s, %v), want %q", i, c.hm, js, err, c.js)
		}
	}
}

func TestParseMMSS(t *testing.T) {
	cs := []struct {
		s string
		w MMSS
	}{
		{"09:00", 540},
		{"08:00", 480},
		{"360", 360},
	}

	for i, c := range cs {
		a, err := ParseMMSS(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ParseMMSS(%q) = (%d, %v), want %d", i, c.s, a, err, c.w)
		}
	}
}

func TestMMSSUnmarshallJSON(t *testing.T) {
	cs := []struct {
		js string
		ms MMSS
	}{
		{`{ "s": "09:00" }`, 540},
		{`{ "s": "10:20" }`, 620},
		{`{ "s": "360" }`, 360},
	}

	o := struct {
		S MMSS `json:"s"`
	}{}

	for i, c := range cs {
		err := json.Unmarshal([]byte(c.js), &o)
		if err != nil || o.S != c.ms {
			t.Errorf("[%d] MMSSUnmarshallJSON(%s) = (%d, %v), want %d", i, c.js, o.S, err, c.ms)
		}
	}
}

func TestMMSSMarshallJSON(t *testing.T) {
	cs := []struct {
		ms MMSS
		js string
	}{
		{540, `{"s":"09:00"}`},
		{620, `{"s":"10:20"}`},
		{360, `{"s":"06:00"}`},
	}

	o := struct {
		S MMSS `json:"s,omitempty"`
	}{}

	for i, c := range cs {
		o.S = c.ms
		bs, err := json.Marshal(&o)
		js := string(bs)
		if err != nil || c.js != js {
			t.Errorf("[%d] MMSSMarshallJSON(%d) = (%s, %v), want %q", i, c.ms, js, err, c.js)
		}
	}
}
