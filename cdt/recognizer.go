package cdt

type recognizer interface {
	Match(*recognizerInput) recognizerOutput
}

type recognizerOutput = Result

type recognizerInput struct {
	raw        []byte
	stripTag   bool
	byteStats  []int
	hasC1Bytes bool
}

type iEach interface {
	Each(func(b byte))
}

type rawInput []byte

func (ri rawInput) Each(f func(b byte)) {
	for _, b := range ri {
		f(b)
	}
}

type stripInput []byte

func (si stripInput) Each(f func(b byte)) {
	var inMarkup bool

	for _, c := range si {
		switch c {
		case '<':
			inMarkup = true
		case '>':
			inMarkup = false
		default:
			if !inMarkup {
				f(c)
			}
		}
	}
}

func newRecognizerInput(raw []byte, stripTag bool) *recognizerInput {
	if stripTag {
		stripTag = mayStripInput(raw)
	}

	byteStats := computeByteStats(raw, stripTag)
	hasC1Bytes := computeHasC1Bytes(byteStats)

	return &recognizerInput{
		raw:        raw,
		stripTag:   stripTag,
		byteStats:  byteStats,
		hasC1Bytes: hasC1Bytes,
	}
}

func mayStripInput(raw []byte) bool {
	var badTags, openTags, chars int
	var inMarkup bool

	for _, c := range raw {
		switch c {
		case '<':
			if inMarkup {
				badTags++
			}
			inMarkup = true
			openTags++
		case '>':
			inMarkup = false
		default:
			if !inMarkup {
				chars++
			}
		}
	}

	// If it looks like this input wasn't marked up, or if it looks like it's
	// essentially nothing but markup abandon the markup stripping.
	// Detection will have to work on the unstripped input.
	if openTags < 5 || openTags/5 < badTags || (chars < 100 && len(raw) > 600) {
		return false
	}
	return true
}

func computeByteStats(input []byte, stripTag bool) []int {
	r := make([]int, 256)

	if stripTag {
		stripInput(input).Each(func(b byte) { r[b]++ })
	} else {
		rawInput(input).Each(func(b byte) { r[b]++ })
	}
	return r
}

func computeHasC1Bytes(byteStats []int) bool {
	for _, count := range byteStats[0x80 : 0x9F+1] {
		if count > 0 {
			return true
		}
	}
	return false
}
