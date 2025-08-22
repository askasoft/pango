package sch

import (
	"fmt"
	"strconv"
	"time"

	"github.com/askasoft/pango/str"
)

var weekdayAbbrs = []string{"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"}
var monthAbbrs = []string{"FOO", "JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}

const (
	_ = iota
	fieldMonth
	fieldDay
	fieldHour
	fieldMinute
	fieldSecond
)

// Cron a cron expression parser and time calculator
// ┌───────────── second (0 - 59) (omittable)
// │ ┌───────────── minute (0 - 59)
// │ │ ┌───────────── hour (0 - 23)
// │ │ │ ┌───────────── day of the month (1 - 31)
// │ │ │ │ ┌───────────── month (1 - 12)
// │ │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday; 7 is also Sunday)
// │ │ │ │ │ │
// │ │ │ │ │ │
// │ │ │ │ │ │
// * * * * * *
// Comma ( , ): used to separate items of a list. For example, "MON,WED,FRI".
// Dash ( - ) : used to define ranges. For example, "1-10"
// Slash (/)  : combined with ranges to specify step values. For example, */5 in the minutes field indicates every 5 minutes It is shorthand for the more verbose POSIX form "5,10,15,20,25,30,35,40,45,50,55,00"
type Cron struct {
	location    *time.Location
	expression  string
	seconds     uint64
	minutes     uint64
	hours       uint64
	daysOfWeek  uint64
	daysOfMonth uint64
	months      uint64
}

// ParseCron parse the cron expression
func ParseCron(expression string, location ...*time.Location) (cron Cron, err error) {
	err = cron.Parse(expression, location...)
	return
}

func (cron *Cron) String() string {
	return cron.expression
}

// Parse parse the cron expression
func (cron *Cron) Parse(expression string, location ...*time.Location) error {
	cron.expression = expression
	if len(location) > 0 {
		cron.location = location[0]
	}
	if cron.location == nil {
		cron.location = time.Local
	}

	fields := str.Fields(expression)
	if z := len(fields); z < 5 || z > 6 {
		return fmt.Errorf("cron: expression must consist of 5-6 fields (found %d in %q)", z, expression)
	}

	i := 0
	if len(fields) == 6 {
		if err := cron.setNumberHits(&cron.seconds, fields[i], 0, 59); err != nil {
			return err
		}
		i++
	} else {
		if err := cron.setNumberHits(&cron.seconds, "0", 0, 59); err != nil {
			return err
		}
	}

	if err := cron.setNumberHits(&cron.minutes, fields[i], 0, 59); err != nil {
		return err
	}
	i++

	if err := cron.setNumberHits(&cron.hours, fields[i], 0, 23); err != nil {
		return err
	}
	i++

	if err := cron.setDaysOfMonth(&cron.daysOfMonth, fields[i]); err != nil {
		return err
	}
	i++

	if err := cron.setMonths(&cron.months, fields[i]); err != nil {
		return err
	}
	i++

	if err := cron.setDays(&cron.daysOfWeek, cron.replaceOrdinals(fields[i], weekdayAbbrs), 7); err != nil {
		return err
	}

	if cron.daysOfWeek&(1<<7) != 0 {
		// Sunday can be represented as 0 or 7
		cron.daysOfWeek &^= (1 << 7)
		cron.daysOfWeek |= 1
	}

	return nil
}

func (cron *Cron) setDaysOfMonth(bits *uint64, field string) error {
	// Days of month start with 1 (in Cron and Golang)
	if err := cron.setDays(bits, field, 31); err != nil {
		return err
	}

	// ... and remove it from the front
	*bits &^= 1
	return nil
}

func (cron *Cron) setDays(bits *uint64, field string, max int) error {
	if str.ContainsByte(field, '?') {
		field = "*"
	}
	return cron.setNumberHits(bits, field, 0, max)
}

func (cron *Cron) setMonths(bits *uint64, value string) error {
	value = cron.replaceOrdinals(value, monthAbbrs)

	// Months start with 1 in Cron and golang
	if err := cron.setNumberHits(bits, value, 1, 12); err != nil {
		return err
	}

	// ... and remove it from the front
	*bits &^= 1
	return nil
}

func (cron *Cron) replaceOrdinals(value string, alias []string) string {
	value = str.ToUpper(value)
	for i, a := range alias {
		value = str.ReplaceAll(value, a, strconv.Itoa(i))
	}
	return value
}

func (cron *Cron) setNumberHits(bits *uint64, value string, min, max int) error {
	*bits = 0

	fields := str.FieldsRune(value, ',')
	for _, field := range fields {
		if str.ContainsByte(field, '/') {
			parts := str.FieldsRune(field, '/')
			if len(parts) != 2 {
				return fmt.Errorf("cron: invalid format of field %q in expression %q", field, cron.expression)
			}

			start, end, err := cron.getRange(parts[0], min, max)
			if err != nil {
				return err
			}

			if !str.ContainsByte(parts[0], '-') {
				end = max
			}

			delta, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("cron: invalid number of field %q in expression %q", field, cron.expression)
			}
			for i := start; i <= end; i += delta {
				*bits |= 1 << i
			}
		} else {
			// Not an incrementer so it must be a range (possibly empty)
			start, end, err := cron.getRange(field, min, max)
			if err != nil {
				return err
			}

			mask := ((uint64(1) << (end - start + 1)) - 1) << start
			*bits |= mask
		}
	}
	return nil
}

func (cron *Cron) getRange(field string, min, max int) (start, end int, err error) {
	if str.ContainsByte(field, '*') {
		start = min
		end = max
		return
	}

	if str.ContainsByte(field, '-') {
		parts := str.FieldsRune(field, '-')
		if len(parts) > 2 {
			err = fmt.Errorf("cron: invalid range format of field %q in expression %q", field, cron.expression)
			return
		}

		start, err = strconv.Atoi(parts[0])
		if err != nil {
			err = fmt.Errorf("cron: invalid range number of field %q in expression %q", field, cron.expression)
			return
		}

		end, err = strconv.Atoi(parts[1])
		if err != nil {
			err = fmt.Errorf("cron: invalid range number of field %q in expression %q", field, cron.expression)
			return
		}
	} else {
		start, err = strconv.Atoi(field)
		if err != nil {
			err = fmt.Errorf("cron: invalid range number of field %q in expression %q", field, cron.expression)
			return
		}
		end = start
	}

	if start > end {
		start, end = end, start
	}

	if start > max || end > max {
		err = fmt.Errorf("cron: exceeded maximum range (%d) of field %q in expression %q", max, field, cron.expression)
		return
	}

	if start < min || end < min {
		err = fmt.Errorf("cron: exceeded minimum range (%d) of field %q in expression %q", min, field, cron.expression)
		return
	}

	return
}

// Next get the next time in the sequence matching the Cron pattern and
// after the value provided. The return value will have a whole number of
// seconds, and will be after the input value.
// @param date a seed value
// @return the next value matching the pattern
func (cron *Cron) Next(date time.Time) time.Time {
	/*
		The plan:

		1 Round up to the next whole second

		2 If seconds match move on, otherwise find the next match:
		2.1 If next match is in the next minute then roll forwards

		3 If minute matches move on, otherwise find the next match
		3.1 If next match is in the next hour then roll forwards
		3.2 Reset the seconds and go to 2

		4 If hour matches move on, otherwise find the next match
		4.1 If next match is in the next day then roll forwards,
		4.2 Reset the minutes and seconds and go to 2

		...
	*/

	// First, just reset the milliseconds and try to calculate from there...
	dorg := date.Truncate(time.Second)
	date = cron.doNext(dorg, dorg.Year())

	if date.Equal(dorg) {
		// We arrived at the original timestamp - round up to the next whole second and try again...
		date = date.Add(time.Second)
		date = cron.doNext(date, date.Year())
	}

	return date
}

func (cron *Cron) doNext(date time.Time, dot int) time.Time {
	resets := []int{}

	second := date.Second()

	date, updateSecond := cron.findNext(date, cron.seconds, second, fieldSecond, time.Minute, resets)
	if second == updateSecond {
		resets = append(resets, fieldSecond)
	}

	minute := date.Minute()
	date, updateMinute := cron.findNext(date, cron.minutes, minute, fieldMinute, time.Hour, resets)
	if minute == updateMinute {
		resets = append(resets, fieldMinute)
	} else {
		date = cron.doNext(date, dot)
	}

	hour := date.Hour()
	date, updateHour := cron.findNext(date, cron.hours, hour, fieldHour, time.Hour*24, resets)
	if hour == updateHour {
		resets = append(resets, fieldHour)
	} else {
		date = cron.doNext(date, dot)
	}

	dayOfWeek := int(date.Weekday())
	dayOfMonth := date.Day()
	date, updateDayOfMonth := cron.findNextDay(date, dayOfMonth, dayOfWeek, resets)
	if dayOfMonth == updateDayOfMonth {
		resets = append(resets, fieldDay)
	} else {
		date = cron.doNext(date, dot)
	}

	month := int(date.Month())
	date, updateMonth := cron.findNextMonth(date, month, resets)
	if month != updateMonth {
		if date.Year()-dot > 4 {
			panic("cron: invalid cron expression \"" + cron.expression + "\" led to runaway search for next trigger")
		}
		date = cron.doNext(date, dot)
	}

	return date
}

func (cron *Cron) findNextDay(date time.Time, dayOfMonth, dayOfWeek int, resets []int) (time.Time, int) {
	count := 0
	limit := 366

	for ; (cron.daysOfMonth&(1<<dayOfMonth) == 0 || cron.daysOfWeek&(1<<dayOfWeek) == 0) && count < limit; count++ {
		date = date.AddDate(0, 0, 1)
		dayOfMonth = date.Day()
		dayOfWeek = int(date.Weekday())
		date = cron.reset(date, resets)
	}
	if count >= limit {
		panic("cron: overflow in day for expression \"" + cron.expression + "\"")
	}

	return date, dayOfMonth
}

func (cron *Cron) findNextMonth(date time.Time, month int, resets []int) (time.Time, int) {
	nextValue := cron.nextSetBit(cron.months, month)

	// roll over if needed
	if nextValue == -1 {
		date = date.AddDate(1, 0, 0)
		date = cron.reset(date, []int{fieldMonth})
		nextValue = cron.nextSetBit(cron.months, 0)
	}
	if nextValue != month {
		date = cron.setField(date, fieldMonth, nextValue)
		date = cron.reset(date, resets)
	}
	return date, nextValue
}

/**
 * Search the bits provided for the next set bit after the value provided,
 * and reset the calendar.
 * @param bits a {@link BitSet} representing the allowed values of the field
 * @param value the current value of the field
 * @param calendar the calendar to increment as we move through the bits
 * @param field the field to increment in the calendar (@see
 * {@link Calendar} for the static constants defining valid fields)
 * @param resets the Calendar field ids that should be reset (i.e. the ones of lower significance than the field of interest)
 * @return the value of the calendar field that is next in the sequence
 */
func (cron *Cron) findNext(date time.Time, bits uint64, value int, field int, nextDuration time.Duration, resets []int) (time.Time, int) {
	nextValue := cron.nextSetBit(bits, value)

	// roll over if needed
	if nextValue == -1 {
		date = date.Add(nextDuration)
		date = cron.reset(date, []int{field})
		nextValue = cron.nextSetBit(bits, 0)
	}
	if nextValue != value {
		date = cron.setField(date, field, nextValue)
		date = cron.reset(date, resets)
	}
	return date, nextValue
}

func (cron *Cron) nextSetBit(bits uint64, start int) int {
	var mask uint64 = 1 << start
	for i := start; i < 64; i++ {
		if (bits & mask) != 0 {
			return i
		}
		mask <<= 1
	}
	return -1
}

func (cron *Cron) setField(date time.Time, field int, value int) time.Time {
	vs := []int{date.Year(), int(date.Month()) - 1, date.Day() - 1, date.Hour(), date.Minute(), date.Second()}

	vs[field] = value

	return time.Date(vs[0], time.Month(vs[1]+1), vs[2]+1, vs[3], vs[4], vs[5], 0, cron.location)
}

// reset the calendar setting all the fields provided to zero.
func (cron *Cron) reset(date time.Time, fields []int) time.Time {
	vs := []int{date.Year(), int(date.Month()) - 1, date.Day() - 1, date.Hour(), date.Minute(), date.Second()}

	for _, field := range fields {
		vs[field] = 0
	}

	return time.Date(vs[0], time.Month(vs[1]+1), vs[2]+1, vs[3], vs[4], vs[5], 0, cron.location)
}
