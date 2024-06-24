package mbs

func Merge(ms ...map[rune]rune) map[rune]rune {
	sz := 0
	for _, m := range ms {
		sz += len(m)
	}

	mm := make(map[rune]rune, sz)
	for _, m := range ms {
		for k, v := range m {
			mm[k] = v
		}
	}
	return mm
}

func Reverse(m map[rune]rune) map[rune]rune {
	rm := make(map[rune]rune, len(m))
	for k, v := range m {
		rm[v] = k
	}
	return rm
}

func HasRune(s string, f func(rune) bool) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if f(c) {
			return true
		}
	}

	return false
}
