package sch

import (
	"fmt"

	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

const (
	Daily   = 'd'
	Weekly  = 'w'
	Monthly = 'm'
)

// Periodic a periodic expression parser
// ┌───────────── unit (d: daily, w: weekly, m: monthly)
// │ ┌───────────── day (daily: 0, weekly: 1-7 monday to sunday, monthly: 1-31)
// │ │ ┌───────────── hour (0 - 23)
// │ │ │
// * * *
type Periodic struct {
	Unit rune // 'd': daily, 'w': weekly, 'm': monthly
	Day  int  // 0 for daily, 1-7 for weekly, 1-31 for monthly
	Hour int  // 0-23
}

// ParsePeriodic parses a periodic expression.
func ParsePeriodic(expr string) (p Periodic, err error) {
	err = p.Parse(expr)
	return
}

// Cron returns the cron expression for the periodic schedule.
func (p *Periodic) Cron() string {
	switch p.Unit {
	case Daily:
		return fmt.Sprintf("0 %d * * *", p.Hour)
	case Weekly:
		return fmt.Sprintf("0 %d * * %d", p.Hour, p.Day)
	case Monthly:
		return fmt.Sprintf("0 %d %d * *", p.Hour, p.Day)
	default:
		return ""
	}
}

func (p *Periodic) String() string {
	return fmt.Sprintf("%c %d %d", p.Unit, p.Day, p.Hour)
}

// Parse parses a periodic expression.
func (p *Periodic) Parse(expr string) (err error) {
	ss := str.Fields(str.ToLower(expr))

	if len(ss) != 3 {
		err = fmt.Errorf("periodic: expression must consist of 3 fields (found %d in %q)", len(ss), expr)
		return
	}

	p.Unit = rune(ss[0][0])
	p.Day = num.Atoi(ss[1])
	p.Hour = num.Atoi(ss[2])

	if p.Hour < 0 || p.Hour > 23 {
		err = fmt.Errorf("periodic: invalid hour %d (must be 0-23)", p.Hour)
		return
	}

	switch p.Unit {
	case Daily:
	case Weekly:
		if p.Day < 1 || p.Day > 7 {
			err = fmt.Errorf("periodic: invalid day %d (must be 1-7)", p.Day)
		}
	case Monthly:
		if p.Day < 1 || p.Day > 31 {
			err = fmt.Errorf("periodic: invalid day %d (must be 1-31)", p.Day)
		}
	default:
		err = fmt.Errorf("periodic: invalid unit %c (must be one of %c, %c, %c)", p.Unit, Daily, Weekly, Monthly)
	}
	return
}
