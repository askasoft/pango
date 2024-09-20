package sqx

import (
	"github.com/askasoft/pango/mag"
	"github.com/askasoft/pango/str"
)

type Quote interface {
	Quote(s string) string
}

type Quoter []rune

var (
	QuoteDefault   = Quoter{'"', '"'}
	QuoteBackticks = Quoter{'`', '`'}
	QuoteBrackets  = Quoter{'[', ']'}
)

var quoters = map[string]Quoter{
	"mysql":     QuoteBackticks,
	"nrmysql":   QuoteBackticks,
	"sqlserver": QuoteBrackets,
	"azuresql":  QuoteBrackets,
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

// Quotes quote string 's' in 'ss' with quote marks [2]rune, return (m[0] + s + m[1])
func (quoter Quoter) Quotes(ss ...string) []string {
	for i, s := range ss {
		ss[i] = quoter.Quote(s)
	}
	return ss
}

// Quote quote string 's' with quotes [2]rune.
// Returns (quoter[0] + s + quoter[1]), if 's' does not contains any "!\"#$%&'()*+,-/:;<=>?@[\\]^`{|}~" characters.
func (quoter Quoter) Quote(s string) string {
	if str.ContainsAny(s, " !\"#$%&'()*+,-/:;<=>?@[\\]^`{|}~") {
		return s
	}

	qs := quoter
	if len(qs) == 0 {
		qs = QuoteDefault
	}

	ss := str.FieldsByte(s, '.')
	for i, s := range ss {
		ss[i] = string(qs[0]) + s + string(qs[1])
	}

	return str.Join(ss, ".")
}
