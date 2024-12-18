package tmu

import (
	"fmt"
	"strings"
	"time"

	"github.com/askasoft/pango/num"
)

const (
	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond //nolint: revive
	Millisecond = time.Millisecond //nolint: revive
	Second      = time.Second      //nolint: revive
	Minute      = time.Minute      //nolint: revive
	Hour        = time.Hour        //nolint: revive
	Day         = time.Hour * 24
)

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

	if d < Second {
		// Special case: if duration is smaller than a second, use smaller units, like 1.2ms
		switch {
		case d < Microsecond:
			fmt.Fprintf(&sb, "%dns", d)
		case d < Millisecond:
			sb.WriteString(num.FtoaWithDigits(float64(d)/float64(Microsecond), 3))
			sb.WriteString("µs")
		default:
			sb.WriteString(num.FtoaWithDigits(float64(d)/float64(Millisecond), 3))
			sb.WriteString("ms")
		}
	} else if d < Minute {
		sb.WriteString(num.FtoaWithDigits(float64(d)/float64(Second), 3))
		sb.WriteString("s")
	} else {
		for d > Second {
			switch {
			case d > Day:
				fmt.Fprintf(&sb, "%dd", d/Day)
				d = d % Day
			case d > Hour:
				fmt.Fprintf(&sb, "%dh", d/Hour)
				d = d % Hour
			case d > Minute:
				fmt.Fprintf(&sb, "%dm", d/Minute)
				d = d % Minute
			case d > Second:
				fmt.Fprintf(&sb, "%ds", d/Second)
				d = d % Second
			}
		}
	}

	return sb.String()
}
