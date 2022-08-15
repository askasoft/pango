package sch

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pandafw/pango/str"
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

type CronSequencer struct {
	Location    *time.Location
	expression  string
	seconds     [60]bool
	minutes     [60]bool
	hours       [24]bool
	daysOfWeek  [8]bool
	daysOfMonth [32]bool
	months      [13]bool
}

func (cs *CronSequencer) Parse(expression string) (err error) {
	cs.expression = expression
	if cs.Location == nil {
		cs.Location = time.Local
	}

	fields := str.Fields(expression)
	if len(fields) != 6 {
		err = fmt.Errorf("Cron expression must consist of 6 fields (found %d in \"%s\")", len(fields), expression)
		return
	}

	err = cs.setNumberHits(cs.seconds[:], fields[0], 0, 59)
	if err != nil {
		return
	}

	err = cs.setNumberHits(cs.minutes[:], fields[1], 0, 59)
	if err != nil {
		return
	}

	err = cs.setNumberHits(cs.hours[:], fields[2], 0, 23)
	if err != nil {
		return
	}

	err = cs.setDaysOfMonth(cs.daysOfMonth[:], fields[3])
	if err != nil {
		return
	}

	err = cs.setMonths(cs.months[:], fields[4])
	if err != nil {
		return
	}

	err = cs.setDays(cs.daysOfWeek[:], cs.replaceOrdinals(fields[5], weekdayAbbrs), 7)
	if err != nil {
		return
	}
	if cs.daysOfWeek[7] {
		// Sunday can be represented as 0 or 7
		cs.daysOfWeek[0] = true
		cs.daysOfWeek[7] = false
	}

	return
}

func (cs *CronSequencer) setDaysOfMonth(bits []bool, field string) (err error) {
	// Days of month start with 1 (in Cron and Golang)
	err = cs.setDays(bits, field, 31)
	if err != nil {
		return
	}

	// ... and remove it from the front
	bits[0] = false
	return
}

func (cs *CronSequencer) setDays(bits []bool, field string, max int) error {
	if str.ContainsByte(field, '?') {
		field = "*"
	}
	return cs.setNumberHits(bits, field, 0, max)
}

func (cs *CronSequencer) setMonths(bits []bool, value string) (err error) {
	value = cs.replaceOrdinals(value, monthAbbrs)

	// Months start with 1 in Cron and golang
	err = cs.setNumberHits(bits, value, 1, 12)
	if err != nil {
		return
	}

	// ... and remove it from the front
	bits[0] = false
	return
}

func (cs *CronSequencer) replaceOrdinals(value string, alias []string) string {
	value = str.ToUpper(value)
	for i, a := range alias {
		value = str.ReplaceAll(value, a, strconv.Itoa(i))
	}
	return value
}

func (cs *CronSequencer) setNumberHits(bits []bool, value string, min, max int) error {
	fields := str.FieldsRune(value, ',')
	for _, field := range fields {
		if !str.ContainsByte(field, '/') {
			// Not an incrementer so it must be a range (possibly empty)
			start, end, err := cs.getRange(field, min, max)
			if err != nil {
				return err
			}
			for i := start; i <= end; i++ {
				bits[i] = true
			}
		} else {
			split := str.FieldsRune(field, '/')
			if len(split) != 2 {
				return fmt.Errorf("Incrementer has more than two fields: '%s' in expression \"%s\"", field, cs.expression)
			}

			start, end, err := cs.getRange(split[0], min, max)
			if err != nil {
				return err
			}

			if !str.ContainsByte(split[0], '-') {
				end = max
			}

			delta, err := strconv.Atoi(split[1])
			if err != nil {
				return fmt.Errorf("Incrementer has invalid number: '%s' in expression \"%s\"", field, cs.expression)
			}
			for i := start; i <= end; i += delta {
				bits[i] = true
			}
		}
	}
	return nil
}

func (cs *CronSequencer) getRange(field string, min, max int) (start, end int, err error) {
	if str.ContainsByte(field, '*') {
		start = min
		end = max
		return
	}

	if !str.ContainsByte(field, '-') {
		start, err = strconv.Atoi(field)
		if err != nil {
			err = fmt.Errorf("Range has invalid number: '%s' in expression \"%s\"", field, cs.expression)
			return
		}
		end = start
	} else {
		split := str.FieldsRune(field, '-')
		if len(split) > 2 {
			err = fmt.Errorf("Range has more than two fields: '%s' in expression \"%s\"", field, cs.expression)
			return
		}

		start, err = strconv.Atoi(split[0])
		if err != nil {
			err = fmt.Errorf("Range has invalid number: '%s' in expression \"%s\"", field, cs.expression)
			return
		}
		end, err = strconv.Atoi(split[1])
		if err != nil {
			err = fmt.Errorf("Range has invalid number: '%s' in expression \"%s\"", field, cs.expression)
			return
		}
	}

	if start > max || end > max {
		err = fmt.Errorf("Range exceeds maximum (%d): '%s' in expression \"%s\"", max, field, cs.expression)
		return
	}

	if start < min || end < min {
		err = fmt.Errorf("Range less than minimum (%d): '%s' in expression \"%s\"", min, field, cs.expression)
		return
	}

	return
}

// * Get the next {@link Date} in the sequence matching the Cron pattern and
// * after the value provided. The return value will have a whole number of
// * seconds, and will be after the input value.
// * @param date a seed value
// * @return the next value matching the pattern
func (cs *CronSequencer) Next(date time.Time) time.Time {
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
	date = cs.doNext(dorg, dorg.Year())

	if date.Equal(dorg) {
		// We arrived at the original timestamp - round up to the next whole second and try again...
		date = date.Add(time.Second)
		date = cs.doNext(date, date.Year())
	}

	return date
}

func (cs *CronSequencer) doNext(date time.Time, dot int) time.Time {
	resets := []int{}

	second := date.Second()

	date, updateSecond := cs.findNext(date, cs.seconds[:], second, fieldSecond, time.Minute, resets)
	if second == updateSecond {
		resets = append(resets, fieldSecond)
	}

	minute := date.Minute()
	date, updateMinute := cs.findNext(date, cs.minutes[:], minute, fieldMinute, time.Hour, resets)
	if minute == updateMinute {
		resets = append(resets, fieldMinute)
	} else {
		date = cs.doNext(date, dot)
	}

	hour := date.Hour()
	date, updateHour := cs.findNext(date, cs.hours[:], hour, fieldHour, time.Hour*24, resets)
	if hour == updateHour {
		resets = append(resets, fieldHour)
	} else {
		date = cs.doNext(date, dot)
	}

	dayOfWeek := int(date.Weekday())
	dayOfMonth := date.Day()
	date, updateDayOfMonth := cs.findNextDay(date, dayOfMonth, dayOfWeek, resets)
	if dayOfMonth == updateDayOfMonth {
		resets = append(resets, fieldDay)
	} else {
		date = cs.doNext(date, dot)
	}

	month := int(date.Month())
	date, updateMonth := cs.findNextMonth(date, month, resets)
	if month != updateMonth {
		if date.Year()-dot > 4 {
			panic("Invalid cron expression \"" + cs.expression + "\" led to runaway search for next trigger")
		}
		date = cs.doNext(date, dot)
	}

	return date
}

func (cs *CronSequencer) findNextDay(date time.Time, dayOfMonth, dayOfWeek int, resets []int) (time.Time, int) {
	count := 0
	limit := 366

	for ; (!cs.daysOfMonth[dayOfMonth] || !cs.daysOfWeek[dayOfWeek]) && count < limit; count++ {
		date = date.AddDate(0, 0, 1)
		dayOfMonth = date.Day()
		dayOfWeek = int(date.Weekday())
		date = cs.reset(date, resets)
	}
	if count >= limit {
		panic("Overflow in day for expression \"" + cs.expression + "\"")
	}

	return date, dayOfMonth
}

func (cs *CronSequencer) findNextMonth(date time.Time, month int, resets []int) (time.Time, int) {
	nextValue := cs.nextSetBit(cs.months[:], month)

	// roll over if needed
	if nextValue == -1 {
		date = date.AddDate(1, 0, 0)
		date = cs.reset(date, []int{fieldMonth})
		nextValue = cs.nextSetBit(cs.months[:], 0)
	}
	if nextValue != month {
		date = cs.setField(date, fieldMonth, nextValue)
		date = cs.reset(date, resets)
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
func (cs *CronSequencer) findNext(date time.Time, bits []bool, value int, field int, nextDuration time.Duration, resets []int) (time.Time, int) {
	nextValue := cs.nextSetBit(bits, value)

	// roll over if needed
	if nextValue == -1 {
		date = date.Add(nextDuration)
		date = cs.reset(date, []int{field})
		nextValue = cs.nextSetBit(bits, 0)
	}
	if nextValue != value {
		date = cs.setField(date, field, nextValue)
		date = cs.reset(date, resets)
	}
	return date, nextValue
}

func (cs *CronSequencer) nextSetBit(bits []bool, start int) int {
	for i := start; i < len(bits); i++ {
		if bits[i] {
			return i
		}
	}
	return -1
}

func (cs *CronSequencer) setField(date time.Time, field int, value int) time.Time {
	vs := []int{date.Year(), int(date.Month()) - 1, date.Day() - 1, date.Hour(), date.Minute(), date.Second()}

	vs[field] = value

	return time.Date(vs[0], time.Month(vs[1]+1), vs[2]+1, vs[3], vs[4], vs[5], 0, cs.Location)
}

// Reset the calendar setting all the fields provided to zero.
func (cs *CronSequencer) reset(date time.Time, fields []int) time.Time {
	vs := []int{date.Year(), int(date.Month()) - 1, date.Day() - 1, date.Hour(), date.Minute(), date.Second()}

	for _, field := range fields {
		vs[field] = 0
	}

	return time.Date(vs[0], time.Month(vs[1]+1), vs[2]+1, vs[3], vs[4], vs[5], 0, cs.Location)
}
