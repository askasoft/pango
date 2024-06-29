package ldt

import (
	"sort"
)

// Detect detects the language info of the given text.
func Detect(text string) Info {
	return DetectWithOptions(text, Options{})
}

// DetectLang detects only the language of the given text.
func DetectLang(text string) Lang {
	return Detect(text).Lang
}

// DetectWithOptions detects the language info of the given text with the provided options.
func DetectWithOptions(text string, options Options) Info {
	if len(text) > 0 {
		detector := detect(text, options.Detectors)
		if detector != nil {
			lang, confidence := detector.Detect(text, options)
			return Info{
				Lang:       lang,
				Confidence: confidence,
			}
		}
	}

	return Info{}
}

// DetectLangWithOptions detects only the language of the given text with the provided options.
func DetectLangWithOptions(text string, options Options) Lang {
	return DetectWithOptions(text, options).Lang
}

type langDistance struct {
	lang Lang
	dist int
}

func detectLangInProfiles(text string, options Options, langProfileList langProfileList, confidence float64) (Lang, float64) {
	trigrams := getTrigramsWithPositions(text)

	langDistances := make([]langDistance, 0, len(langProfileList))

	for lang, langTrigrams := range langProfileList {
		if options.exclude(lang) { // skip excluded languages.
			continue
		}

		if !options.include(lang) { // skip non-included languages.
			continue
		}

		dist := calculateDistance(langTrigrams, trigrams)
		langDistances = append(langDistances, langDistance{lang, dist})
	}

	switch len(langDistances) {
	case 0:
		return Unknown, 0
	case 1:
		return langDistances[0].lang, confidence
	default:
		return calculateConfidence(langDistances, trigrams)
	}
}

func calculateConfidence(langDistances []langDistance, trigrams map[string]int) (Lang, float64) {
	sort.SliceStable(langDistances, func(i, j int) bool {
		return langDistances[i].dist < langDistances[j].dist
	})

	maxDist := len(trigrams) * maxTrigramDistance
	if maxDist > maxTotalDistance {
		maxDist = maxTotalDistance
	}

	langDist1 := langDistances[0]
	langDist2 := langDistances[1]
	score1 := maxDist - langDist1.dist
	score2 := maxDist - langDist2.dist

	if score1 == 0 {
		// If score1 is 0, score2 is 0 as well, because array is sorted.
		// Therefore there is no language to return.
		return Unknown, 0
	}

	if score2 == 0 {
		// If score2 is 0, return first language, to prevent division by zero in the rate formula.
		// In this case confidence is calculated by another formula.
		// At this point there are two options:
		// * Text contains random characters that accidentally match trigrams of one of the languages
		// * Text really matches one of the languages.
		//
		// Number 500.0 is based on experiments and common sense expectations.
		confidence := float64(score1) / float64(500)
		if confidence > 1.0 {
			confidence = 1.0
		}
		return langDist1.lang, confidence
	}

	rate := float64(score1-score2) / float64(score2)

	// Hyperbola function. Everything that is above the function has confidence = 1.0
	// If rate is below, confidence is calculated proportionally.
	// Constants are used based on experiments.
	confidentRate := float64(12.0/float64(len(trigrams))) + 0.05

	var confidence float64
	if rate > confidentRate {
		confidence = 1.0
	} else {
		confidence = rate / confidentRate
	}

	return langDist1.lang, confidence
}

func calculateDistance(langTrigrams []string, textTrigrams map[string]int) int {
	var dist, totalDist int
	for i, trigram := range langTrigrams {
		if n, ok := textTrigrams[trigram]; ok {
			dist = abs(n - i)
		} else {
			dist = maxTrigramDistance
		}
		totalDist += dist
	}

	delta := maxTrigramDistance - len(textTrigrams)
	if delta > 0 {
		totalDist -= delta * maxTrigramDistance
	}

	if totalDist > maxTotalDistance {
		totalDist = maxTotalDistance
	}

	return totalDist
}
