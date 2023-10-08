package fdk

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTimeString(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")

	tm := Time{time.Date(2020, 10, 1, 13, 14, 15, 0, jst)}

	fmt.Println(tm.String())
	fmt.Println(tm.Time)
}

func TestTimeParse(t *testing.T) {
	tml, _ := time.Parse(TimeFormat, "2020-01-02T03:04:05Z")
	fmt.Println(tml.String())

	tmu, _ := time.ParseInLocation(TimeFormat, "2020-01-02T03:04:05Z", time.UTC)
	fmt.Println(tmu.String())
}

func TestTimeSpentJSONUnmarshall(t *testing.T) {
	o := struct {
		S TimeSpent `json:"s"`
	}{}

	s := `{ "s": "10:20" }`

	err := json.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Fatal(err)
	}

	if o.S.String() != "10:20" {
		t.Errorf("want 10:20, but %s", o.S.String())
	}
}

func TestTimeSpentJSONMarshall(t *testing.T) {
	o := struct {
		S TimeSpent `json:"s,omitempty"`
	}{}

	bs, err := json.Marshal(&o)
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "{}" {
		t.Errorf("want {}, but %s", string(bs))
	}

	o.S = 361
	bs, err = json.Marshal(&o)
	if err != nil {
		t.Fatal(err)
	}

	w := `{"s":"06:01"}`
	if string(bs) != w {
		t.Errorf("want %s, but %s", w, string(bs))
	}
}
