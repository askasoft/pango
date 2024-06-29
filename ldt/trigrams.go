package ldt

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type trigram struct {
	trigram string
	count   int
}

// convert punctuations and digits to space.
func toTrigramChar(ch rune) rune {
	if isStopChar(ch) {
		return ' '
	}
	return unicode.ToLower(ch)
}

func getTrigramsWithPositions(text string) map[string]int {
	counters := count(text)
	trigrams := make([]trigram, len(counters))

	i := 0
	for tg, count := range counters {
		trigrams[i] = trigram{tg, count}
		i++
	}

	// Sort in descending order by number of occurrences and trigrams
	sort.SliceStable(trigrams, func(i, j int) bool {
		a := trigrams[i]
		b := trigrams[j]
		if a.count == b.count {
			return strings.Compare(a.trigram, a.trigram) > 0
		}
		return a.count > b.count
	})

	// we're interested only in the first 600 (2 * MAX_TRIGRAM_DISTANCE)
	size := textTrigramSize
	if len(trigrams) < size {
		size = len(trigrams)
	}

	trigramsWithPositions := make(map[string]int, size)

	for i := 0; i < size; i++ {
		trigramsWithPositions[trigrams[i].trigram] = i
	}
	return trigramsWithPositions
}

func count(s string) map[string]int {
	// In order to improve performance, define the initial capacity for trigrams hash map
	z := len(s)
	if z == 0 {
		return nil
	}

	if z > 1024 {
		z = 1024
	}
	trigrams := make(map[string]int, z)

	var c1, c2, c3 rune
	c1 = ' '
	c2, z = utf8.DecodeRuneInString(s)
	c2 = toTrigramChar(c2)
	s = s[z:]

	next := true
	for next {
		if len(s) == 0 {
			c3 = ' '
			next = false
		} else {
			c3, z = utf8.DecodeRuneInString(s)
			c3 = toTrigramChar(c3)
			s = s[z:]
		}

		if !(c2 == ' ' && (c1 == ' ' || c3 == ' ')) {
			trigram := string([]rune{c1, c2, c3})
			trigrams[trigram]++
		}
		c1 = c2
		c2 = c3
	}

	return trigrams
}
