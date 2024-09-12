package sqx

import (
	"sync"

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

var quotes sync.Map

func init() {
	defaultQuotes := map[Quoter][]string{
		QuoteMYSQL: {"mysql", "nrmysql"},
		QuoteMSSQL: {"sqlserver", "azuresql"},
	}

	for typ, drivers := range defaultQuotes {
		for _, driver := range drivers {
			QuoteDriver(driver, typ)
		}
	}
}

// GetQuoteer returns the quoter for a given database given a drivername.
func GetQuoter(driverName string) Quoter {
	quoter, ok := quotes.Load(driverName)
	if !ok {
		return QuoteDefault
	}
	return quoter.(Quoter)
}

// QuoteDriver sets the Quoter for driverName to quoter.
func QuoteDriver(driverName string, quoter Quoter) {
	quotes.Store(driverName, quoter)
}
