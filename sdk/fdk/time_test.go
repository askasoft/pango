package fdk

import (
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
