package tmu

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/num"
)

const (
	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
	Day         = time.Hour * 24
)

// Atod convert string to time.Duration.
// if not found or convert error, returns the default defs[0] value.
func Atod(s string, defs ...time.Duration) time.Duration {
	if d, err := ParseDuration(s); err == nil {
		return d
	}
	return asg.First(defs)
}

// HumanDuration returns a string representing the duration in the form "3d23h3m5s".
// Leading zero units are omitted. As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
// that the leading digit is non-zero. The zero duration formats as 0s.
func HumanDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	var sb strings.Builder

	if d < 0 {
		d = -d
		sb.WriteByte('-')
	}

	switch {
	case d < Second:
		// Special case: if duration is smaller than a second, use smaller units, like 1.2ms
		switch {
		case d < Microsecond:
			sb.WriteString(num.Ltoa(int64(d)))
			sb.WriteString("ns")
		case d < Millisecond:
			sb.WriteString(num.FtoaWithDigits(float64(d)/float64(Microsecond), 3))
			sb.WriteString("µs")
		default:
			sb.WriteString(num.FtoaWithDigits(float64(d)/float64(Millisecond), 3))
			sb.WriteString("ms")
		}
	case d < Minute:
		sb.WriteString(num.FtoaWithDigits(float64(d)/float64(Second), 3))
		sb.WriteString("s")
	default:
		for d > Second {
			switch {
			case d > Day:
				sb.WriteString(num.Ltoa(int64(d / Day)))
				sb.WriteByte('d')
				d = d % Day
			case d > Hour:
				sb.WriteString(num.Ltoa(int64(d / Hour)))
				sb.WriteByte('h')
				d = d % Hour
			case d > Minute:
				sb.WriteString(num.Ltoa(int64(d / Minute)))
				sb.WriteByte('m')
				d = d % Minute
			case d > Second:
				sb.WriteString(num.Ltoa(int64(d / Second)))
				sb.WriteByte('s')
				d = d % Second
			}
		}
	}

	return sb.String()
}

var errLeadingInt = errors.New("tmu: bad [0-9]*") // never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt[bytes []byte | string](s bytes) (x uint64, rem bytes, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, rem, errLeadingInt
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, rem, errLeadingInt
		}
	}
	return x, s[i:], nil
}

// leadingFraction consumes the leading [0-9]* from s.
// It is used only for fractions, so does not return an error on overflow,
// it just stops accumulating precision.
func leadingFraction(s string) (x uint64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > (1<<63-1)/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + uint64(c) - '0'
		if y > 1<<63 {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}

var unitMap = map[string]uint64{
	"ns": uint64(Nanosecond),
	"us": uint64(Microsecond),
	"µs": uint64(Microsecond), // U+00B5 = micro symbol
	"μs": uint64(Microsecond), // U+03BC = Greek letter mu
	"ms": uint64(Millisecond),
	"s":  uint64(Second),
	"m":  uint64(Minute),
	"h":  uint64(Hour),
	"d":  uint64(Day),
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h", "d".
func ParseDuration(s string) (time.Duration, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d uint64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, fmt.Errorf("tmu: invalid duration %q", orig)
	}
	for s != "" {
		var (
			v, f  uint64      // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if s[0] != '.' && (s[0] < '0' || s[0] > '9') {
			return 0, fmt.Errorf("tmu: invalid duration %q", orig)
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, fmt.Errorf("tmu: invalid duration %q", orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0, fmt.Errorf("tmu: invalid duration %q", orig)
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, fmt.Errorf("tmu: missing unit in duration %q", orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, fmt.Errorf("tmu: unknown unit %q in duration %q", u, orig)
		}
		if v > 1<<63/unit {
			// overflow
			return 0, fmt.Errorf("tmu: invalid duration %q", orig)
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += uint64(float64(f) * (float64(unit) / scale))
			if v > 1<<63 {
				// overflow
				return 0, fmt.Errorf("tmu: invalid duration %q", orig)
			}
		}
		d += v
		if d > 1<<63 {
			return 0, fmt.Errorf("tmu: invalid duration %q", orig)
		}
	}
	if neg {
		return -time.Duration(d), nil
	}
	if d > 1<<63-1 {
		return 0, fmt.Errorf("tmu: invalid duration %q", orig)
	}
	return time.Duration(d), nil
}
