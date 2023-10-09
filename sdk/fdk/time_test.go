package fdk

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	tml, _ := time.Parse(TimeFormat, "2020-01-02T03:04:05Z")
	fmt.Println(tml.String())

	tmu, _ := time.ParseInLocation(TimeFormat, "2020-01-02T03:04:05Z", time.UTC)
	fmt.Println(tmu.String())
}

func TestParseTimeSpent(t *testing.T) {
	cs := []struct {
		s string
		w TimeSpent
	}{
		{"09:00", 540},
		{"08:00", 480},
		{"360", 360},
	}

	for i, c := range cs {
		a, err := ParseTimeSpent(c.s)
		if err != nil || a != c.w {
			t.Errorf("[%d] ParseTimeSpent(%q) = (%d, %v), want %d", i, c.s, a, err, c.w)
		}
	}
}

func TestTimeSpentUnmarshallJSON(t *testing.T) {
	cs := []struct {
		js string
		ts TimeSpent
	}{
		{`{ "s": "09:00" }`, 540},
		{`{ "s": "10:20" }`, 620},
		{`{ "s": "360" }`, 360},
	}

	o := struct {
		S TimeSpent `json:"s"`
	}{}

	for i, c := range cs {
		err := json.Unmarshal([]byte(c.js), &o)
		if err != nil || o.S != c.ts {
			t.Errorf("[%d] TimeSpentUnmarshallJSON(%s) = (%d, %v), want %d", i, c.js, o.S, err, c.ts)
		}
	}
}

func TestTimeSpentMarshallJSON(t *testing.T) {
	cs := []struct {
		ts TimeSpent
		js string
	}{
		{540, `{"s":"09:00"}`},
		{620, `{"s":"10:20"}`},
		{360, `{"s":"06:00"}`},
	}

	o := struct {
		S TimeSpent `json:"s,omitempty"`
	}{}

	for i, c := range cs {
		o.S = c.ts
		bs, err := json.Marshal(&o)
		js := string(bs)
		if err != nil || c.js != js {
			t.Errorf("[%d] TimeSpentMarshallJSON(%d) = (%s, %v), want %q", i, c.ts, js, err, c.js)
		}
	}
}
