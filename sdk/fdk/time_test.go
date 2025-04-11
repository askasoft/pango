package fdk

import (
	"fmt"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	tml, _ := time.Parse(TimeFormat, "2020-01-02T03:04:05Z")
	fmt.Println(tml.String())

	tmu, _ := time.ParseInLocation(TimeFormat, "2020-01-02T03:04:05Z", time.UTC)
	fmt.Println(tmu.String())

	tm3, _ := time.ParseInLocation(TimeFormat, "2020-01-02T03:04:05+08:00", time.UTC)
	fmt.Println(tm3.String())
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
