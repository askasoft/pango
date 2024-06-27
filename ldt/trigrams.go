package ldt

import (
	"sort"
	"strings"
	"unicode"
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
	return ch
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

func count(text string) map[string]int {
	var txt []rune
	for _, r := range text {
		txt = append(txt, unicode.ToLower(toTrigramChar(r)))
	}
	txt = append(txt, ' ')

	var r1, r2, r3 rune
	trigrams := map[string]int{}

	r1 = ' '
	r2 = txt[0]
	for i := 1; i < len(txt); i++ {
		r3 = txt[i]
		if !(r2 == ' ' && (r1 == ' ' || r3 == ' ')) {
			trigram := string([]rune{r1, r2, r3})
			trigrams[trigram]++
		}
		r1 = r2
		r2 = r3
	}

	return trigrams
}
