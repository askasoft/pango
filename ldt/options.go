package ldt

// Options represents options that can be set when detecting a language or/and script such
// blacklisting languages to skip checking.
type Options struct {
	Detectors []Detector
	Includes  []Lang
	Excludes  []Lang
}

func (o Options) contains(langs []Lang, lang Lang) bool {
	for _, v := range langs {
		if lang == v {
			return true
		}
	}
	return false
}

func (o Options) include(lang Lang) bool {
	return len(o.Includes) == 0 || o.contains(o.Includes, lang)
}

func (o Options) exclude(lang Lang) bool {
	return len(o.Excludes) > 0 && o.contains(o.Excludes, lang)
}
