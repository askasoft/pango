package ldt

import (
	"sort"
)

// Detect language and script of the given text.
func Detect(text string) Info {
	return DetectWithOptions(text, Options{})
}

// DetectLang detects only the language by a given text.
func DetectLang(text string) Lang {
	return Detect(text).Lang
}

// DetectLangWithOptions detects only the language of the given text with the provided options.
func DetectLangWithOptions(text string, options Options) Lang {
	return DetectWithOptions(text, options).Lang
}

// DetectWithOptions detects the language and script of the given text with the provided options.
func DetectWithOptions(text string, options Options) Info {
	detector := detect(text)
	if detector != nil {
		lang, confidence := detector.Detect(text, options)
		return Info{
			Lang:       lang,
			Confidence: confidence,
		}
	}

	return Info{}
}

type langDistance struct {
	lang Lang
	dist int
}

func detectLangInProfiles(text string, options Options, langProfileList langProfileList) (Lang, float64) {
	trigrams := getTrigramsWithPositions(text)

	langDistances := []langDistance{}

	for lang, langTrigrams := range langProfileList {
		if len(options.Whitelist) != 0 {
			//Skip non-whitelisted languages.
			if _, ok := options.Whitelist[lang]; !ok {
				continue
			}
		} else if len(options.Blacklist) != 0 {
			//skip blacklisted languages.
			if _, ok := options.Blacklist[lang]; ok {
				continue
			}
		}

		dist := calculateDistance(langTrigrams, trigrams)
		langDistances = append(langDistances, langDistance{lang, dist})
	}

	switch len(langDistances) {
	case 0:
		return UNKNOWN, 0
	case 1:
		return langDistances[0].lang, 1
	default:
		return calculateConfidence(langDistances, trigrams)
	}
}

func calculateConfidence(langDistances []langDistance, trigrams map[string]int) (Lang, float64) {
	sort.SliceStable(langDistances, func(i, j int) bool { return langDistances[i].dist < langDistances[j].dist })
	langDist1 := langDistances[0]
	langDist2 := langDistances[1]
	score1 := maxTotalDistance - langDist1.dist
	score2 := maxTotalDistance - langDist2.dist

	var confidence float64
	if score1 == 0 {
		// If score1 is 0, score2 is 0 as well, because array is sorted.
		// Therefore there is no language to return.
		return UNKNOWN, 0
	} else if score2 == 0 {
		// If score2 is 0, return first language, to prevent division by zero in the rate formula.
		// In this case confidence is calculated by another formula.
		// At this point there are two options:
		// * Text contains random characters that accidentally match trigrams of one of the languages
		// * Text really matches one of the languages.
		//
		// Number 500.0 is based on experiments and common sense expectations.
		confidence = float64((score1) / 500.0)
		if confidence > 1.0 {
			confidence = 1.0
		}
		return langDist1.lang, confidence
	}

	rate := float64((score1 - score2)) / float64(score2)

	// Hyperbola function. Everything that is above the function has confidence = 1.0
	// If rate is below, confidence is calculated proportionally.
	// Numbers 12.0 and 0.05 are obtained experimentally, so the function represents common sense.

	confidentRate := float64(12.0/float64(len(trigrams))) + 0.05
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

	return totalDist
}
