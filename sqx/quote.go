package sqx

import (
	"github.com/askasoft/pango/mag"
	"github.com/askasoft/pango/str"
)

type Quote interface {
	Quote(s string) string
}

type QuoteMarks []rune

var (
	QuoteMarksDefault = QuoteMarks{'"', '"'}
	QuoteMarksMSSQL   = QuoteMarks{'[', ']'}
	QuoteMarksMYSQL   = QuoteMarks{'`', '`'}
)

type Quoter int

const (
	QuoteDefault Quoter = iota
	QuoteMYSQL
	QuoteMSSQL
)

func (quoter Quoter) Marks() QuoteMarks {
	switch quoter {
	case QuoteMYSQL:
		return QuoteMarksMYSQL
	case QuoteMSSQL:
		return QuoteMarksMSSQL
	default:
		return QuoteMarksDefault
	}
}

// Quotes quote string 's' in 'ss' with quote marks [2]rune, return (m[0] + s + m[1])
func (quoter Quoter) Quotes(ss ...string) []string {
	for i, s := range ss {
		ss[i] = quoter.Quote(s)
	}
	return ss
}

// Quote quote string 's' with quote marks [2]rune.
// Returns (m[0] + s + m[1]), if 's' does not contains any "!\"#$%&'()*+,-/:;<=>?@[\\]^`{|}~" characters.
func (quoter Quoter) Quote(s string) string {
	if str.ContainsAny(s, " !\"#$%&'()*+,-/:;<=>?@[\\]^`{|}~") {
		return s
	}

	qms := quoter.Marks()

	ss := str.FieldsByte(s, '.')
	for i, s := range ss {
		ss[i] = string(qms[0]) + s + string(qms[1])
	}

	return str.Join(ss, ".")
}

var quoters map[string]Quoter

func init() {
	defaultQuoters := map[Quoter][]string{
		QuoteMYSQL: {"mysql", "nrmysql"},
		QuoteMSSQL: {"sqlserver", "azuresql"},
	}

	quoters = make(map[string]Quoter)
	for quoter, drivers := range defaultQuoters {
		for _, driver := range drivers {
			quoters[driver] = quoter
		}
	}
}

// GetQuoteer returns the quoter for a given database given a drivername.
func GetQuoter(driverName string) Quoter {
	quoter, ok := quoters[driverName]
	if !ok {
		return QuoteDefault
	}
	return quoter
}

// QuoteDriver sets the Quoter for driverName to quoter.
func QuoteDriver(driverName string, quoter Quoter) {
	nqs := make(map[string]Quoter)
	mag.Copy(nqs, quoters)

	nqs[driverName] = quoter
	quoters = nqs
}
