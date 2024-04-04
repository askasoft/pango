package sqx

import (
	"sync"
)

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

// Quote quote string 's' with quote marks [2]rune, return (m[0] + s + m[1])
func (quoter Quoter) Quote(s string) string {
	qms := quoter.Marks()
	return string(qms[0]) + s + string(qms[1])
}

var quotes sync.Map

func init() {
	defaultQuotes := map[Quoter][]string{
		QuoteMYSQL: {"mysql", "nrmysql"},
		QuoteMSSQL: {"sqlserver"},
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
