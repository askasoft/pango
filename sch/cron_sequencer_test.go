package sch

import (
	"testing"
	"time"

	"github.com/askasoft/pango/str"
)

const testCronTimeFormat = "2006-01-02T15:04:05"

func testCronNext(t *testing.T, cron string, sdt string, ns []string) {
	cs := &CronSequencer{}
	err := cs.Parse(cron)
	if err != nil {
		t.Fatal(err)
	}

	d, err := time.Parse(testCronTimeFormat, sdt)
	if err != nil {
		t.Fatal(err)
	}
	for i, sw := range ns {
		d = cs.Next(d)
		sa := d.Format(testCronTimeFormat)
		if sw != sa {
			t.Errorf("[%d] Got %v, want %v", i, sa, sw)
		}
	}
}

func testCronSequencer(t *testing.T, cron string, sdt string, ns []string) {
	testCronNext(t, cron, sdt, ns)

	if str.StartsWith(cron, "0 ") {
		testCronNext(t, cron[2:], sdt, ns)
	}
}

func TestCronEvery1Min(t *testing.T) {
	testCronSequencer(t, "0 * * * * *", "2000-01-01T00:00:00", []string{
		"2000-01-01T00:01:00",
		"2000-01-01T00:02:00",
		"2000-01-01T00:03:00",
		"2000-01-01T00:04:00",
		"2000-01-01T00:05:00",
		"2000-01-01T00:06:00",
		"2000-01-01T00:07:00",
	})
}

func TestCronEvery10Min_(t *testing.T) {
	testCronSequencer(t, "0 */10 * * * *", "2000-01-01T01:01:01", []string{
		"2000-01-01T01:10:00",
		"2000-01-01T01:20:00",
		"2000-01-01T01:30:00",
		"2000-01-01T01:40:00",
		"2000-01-01T01:50:00",
		"2000-01-01T02:00:00",
		"2000-01-01T02:10:00",
	})
}

func TestCronEvery10Min0(t *testing.T) {
	testCronSequencer(t, "0 0/10 * * * *", "2000-01-01T01:01:01", []string{
		"2000-01-01T01:10:00",
		"2000-01-01T01:20:00",
		"2000-01-01T01:30:00",
		"2000-01-01T01:40:00",
		"2000-01-01T01:50:00",
		"2000-01-01T02:00:00",
		"2000-01-01T02:10:00",
	})
}

func TestCronEvery10Min2(t *testing.T) {
	testCronSequencer(t, "0 2/10 * * * *", "2000-01-01T01:01:01", []string{
		"2000-01-01T01:02:00",
		"2000-01-01T01:12:00",
		"2000-01-01T01:22:00",
		"2000-01-01T01:32:00",
		"2000-01-01T01:42:00",
		"2000-01-01T01:52:00",
		"2000-01-01T02:02:00",
		"2000-01-01T02:12:00",
	})
}

func TestCronEvery20Min11(t *testing.T) {
	testCronSequencer(t, "0 11/20 * * * *", "2000-01-01T01:01:01", []string{
		"2000-01-01T01:11:00",
		"2000-01-01T01:31:00",
		"2000-01-01T01:51:00",
		"2000-01-01T02:11:00",
		"2000-01-01T02:31:00",
	})
}

func TestCronHour(t *testing.T) {
	testCronSequencer(t, "0 0 * * * *", "2000-01-01T01:01:01", []string{
		"2000-01-01T02:00:00",
		"2000-01-01T03:00:00",
		"2000-01-01T04:00:00",
	})
}

func TestCronHourRange(t *testing.T) {
	testCronSequencer(t, "0 0 8-10 * * *", "2000-01-01T01:01:01", []string{
		"2000-01-01T08:00:00",
		"2000-01-01T09:00:00",
		"2000-01-01T10:00:00",
		"2000-01-02T08:00:00",
	})
}

func TestCronHourMinRange(t *testing.T) {
	testCronSequencer(t, "0 0/30 8-10 * * *", "2000-01-01T01:01:01", []string{
		"2000-01-01T08:00:00",
		"2000-01-01T08:30:00",
		"2000-01-01T09:00:00",
		"2000-01-01T09:30:00",
		"2000-01-01T10:00:00",
		"2000-01-01T10:30:00",
		"2000-01-02T08:00:00",
		"2000-01-02T08:30:00",
	})
}

func TestCronEveryDay(t *testing.T) {
	testCronSequencer(t, "0 0 9 * * *", "2000-01-30T01:01:01", []string{
		"2000-01-30T09:00:00",
		"2000-01-31T09:00:00",
		"2000-02-01T09:00:00",
		"2000-02-02T09:00:00",
	})
}

func TestCronEverySundayDay0(t *testing.T) {
	testCronSequencer(t, "0 0 2 * * 0", "2025-02-09T01:01:01", []string{
		"2025-02-09T02:00:00",
		"2025-02-16T02:00:00",
		"2025-02-23T02:00:00",
		"2025-03-02T02:00:00",
		"2025-03-09T02:00:00",
		"2025-03-16T02:00:00",
	})
}

func TestCronEverySundayDay7(t *testing.T) {
	testCronSequencer(t, "0 0 2 * * 7", "2025-02-09T01:01:01", []string{
		"2025-02-09T02:00:00",
		"2025-02-16T02:00:00",
		"2025-02-23T02:00:00",
		"2025-03-02T02:00:00",
		"2025-03-09T02:00:00",
		"2025-03-16T02:00:00",
	})
}

func TestCronEverySundayDayA(t *testing.T) {
	testCronSequencer(t, "0 0 2 * * SUN", "2025-02-09T01:01:01", []string{
		"2025-02-09T02:00:00",
		"2025-02-16T02:00:00",
		"2025-02-23T02:00:00",
		"2025-03-02T02:00:00",
		"2025-03-09T02:00:00",
		"2025-03-16T02:00:00",
	})
}

func TestCronMonthRange(t *testing.T) {
	testCronSequencer(t, "0 0 9-10 * * MON-TUE", "2000-01-01T01:01:01", []string{
		"2000-01-03T09:00:00",
		"2000-01-03T10:00:00",
		"2000-01-04T09:00:00",
		"2000-01-04T10:00:00",
		"2000-01-10T09:00:00",
		"2000-01-10T10:00:00",
	})
}

func TestCronMonthDay31(t *testing.T) {
	testCronSequencer(t, "0 0 2 31 * *", "2000-01-01T01:01:01", []string{
		"2000-01-31T02:00:00",
		"2000-03-31T02:00:00",
		"2000-05-31T02:00:00",
		"2000-07-31T02:00:00",
		"2000-08-31T02:00:00",
		"2000-10-31T02:00:00",
	})
}
