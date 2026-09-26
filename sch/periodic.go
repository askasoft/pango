package sch

import (
	"fmt"

	"github.com/askasoft/pango/str"
)

const (
	Daily   = 'd'
	Weekly  = 'w'
	Monthly = 'm'
)

// Periodic a periodic expression parser
// ┌───────────── unit (d: daily, w: weekly, m: monthly)
// │ ┌───────────── day (daily: 0, weekly: 1-7 monday to sunday, monthly: 1-31,32 is last day of month)
// │ │ ┌───────────── hour (0 - 23)
// │ │ │ ┌───────────── minute (0 - 59) (omittable)
// │ │ │ │ ┌───────────── second (0 - 59) (omittable)
// * * * * *
// Comma ( , ): used to separate items of a list. For example, "MON,WED,FRI".
// Dash ( - ) : used to define ranges. For example, "1-10"
type Periodic struct {
	expression string // original expression
}

// ParsePeriodic parses a periodic expression.
func ParsePeriodic(expr string) (p Periodic, err error) {
	err = p.Parse(expr)
	return
}

// Cron returns the cron expression for the periodic schedule.
func (p *Periodic) Cron() string {
	ss := str.Fields(str.ToLower(p.expression))

	if z := len(ss); z < 3 || z > 5 {
		return p.expression
	}

	minute := "0"
	if len(ss) > 3 {
		minute = ss[3]
	}

	second := ""
	if len(ss) > 4 {
		second = ss[4] + " "
	}

	switch ss[0][0] {
	case Daily:
		return fmt.Sprintf("%s%s %s * * *", second, minute, ss[2])
	case Weekly:
		return fmt.Sprintf("%s%s %s * * %s", second, minute, ss[2], ss[1])
	case Monthly:
		return fmt.Sprintf("%s%s %s %s * *", second, minute, ss[2], ss[1])
	default:
		return p.expression
	}
}

func (p *Periodic) String() string {
	return p.expression
}

// Parse parses a periodic expression.
func (p *Periodic) Parse(expr string) (err error) {
	p.expression = expr

	ss := str.Fields(str.ToLower(expr))

	if z := len(ss); z < 3 || z > 5 {
		err = fmt.Errorf("periodic: expression must consist of 3-5 fields (found %d in %q)", z, expr)
		return
	}

	unit := ss[0][0]
	days := ss[1]

	switch unit {
	case Daily:
	case Weekly:
		err = p.checkDays("weekdays", days, 1, 7)
	case Monthly:
		err = p.checkDays("days", days, 1, 32)
	default:
		err = fmt.Errorf("periodic: invalid unit %c (must be one of %c, %c, %c) in expression %q", unit, Daily, Weekly, Monthly, expr)
	}

	if err == nil {
		err = p.checkHours(ss[2])
	}
	if err == nil && len(ss) > 3 {
		err = p.checkMinutes(ss[3])
	}
	if err == nil && len(ss) > 4 {
		err = p.checkSeconds(ss[4])
	}

	return
}

func (p *Periodic) checkDays(name, value string, min, max int) error {
	if _, err := getNumberHits(name, value, min, max); err != nil {
		return fmt.Errorf("periodic: %w in expression %q", err, p.expression)
	}
	return nil
}

func (p *Periodic) checkHours(value string) error {
	if _, err := getNumberHits("hours", value, 0, 23); err != nil {
		return fmt.Errorf("periodic: %w in expression %q", err, p.expression)
	}
	return nil
}

func (p *Periodic) checkMinutes(value string) error {
	if _, err := getNumberHits("minutes", value, 0, 59); err != nil {
		return fmt.Errorf("periodic: %w in expression %q", err, p.expression)
	}
	return nil
}

func (p *Periodic) checkSeconds(value string) error {
	if _, err := getNumberHits("seconds", value, 0, 59); err != nil {
		return fmt.Errorf("periodic: %w in expression %q", err, p.expression)
	}
	return nil
}
