package ldt

import (
	"strings"
)

type checker func(r rune) bool

type detector interface {
	Count(string) bool
	Detect(string, Options) (Lang, float64)
	Chars() int
	Words() int
}

type scriptCounter struct {
	total int
	chars int
	words int
}

func (sc *scriptCounter) confidence() float64 {
	if sc.total > 0 {
		return float64(sc.chars) / float64(sc.total)
	}
	return 0
}

func (sc *scriptCounter) count(s string, check checker) bool {
	chars := 0
	for _, c := range s {
		if check(c) {
			chars++
		}
		sc.total++
	}

	if chars > 0 {
		sc.words++
		sc.chars += chars
	}

	return false
}

func (sc *scriptCounter) Chars() int {
	return sc.chars
}

func (sc *scriptCounter) Words() int {
	return sc.words
}

type Detectors []detector

func (ds Detectors) Len() int {
	return len(ds)
}

func (ds Detectors) Less(i, j int) bool {
	return lessDetectors(ds, i, j)
}

func lessDetectors(ds Detectors, i, j int) bool {
	a, b := ds[i], ds[j]
	if a.Words() == b.Words() {
		return a.Chars() < b.Chars()
	}
	return a.Words() < b.Words()
}

func (ds Detectors) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

func (ds Detectors) Best() detector {
	if len(ds) == 0 {
		return nil
	}

	m := 0
	for i := 1; i < len(ds); i++ {
		if ds.Less(m, i) {
			m = i
		}
	}

	return ds[m]
}

// detect returns only the detector of the given text.
func detect(text string) detector {
	detectors := Detectors{
		&latinDetector{},
		&cyrillicDetector{},
		&arabicDetector{},
		&devanagariDetector{},
		&ethiopicDetector{},
		&hebrewDetector{},
		&bengaliDetector{},
		&georgianDetector{},
		&greekDetector{},
		&kannadaDetector{},
		&tamilDetector{},
		&thaiDetector{},
		&gujaratiDetector{},
		&gurmukhiDetector{},
		&teluguDetector{},
		&malayalamDetector{},
		&oriyaDetector{},
		&myanmarDetector{},
		&khmerDetector{},
		&sinhalaDetector{},
		&cjkDetector{},
	}

	ss := strings.FieldsFunc(text, isStopChar)
	for _, s := range ss {
		for _, d := range detectors {
			if d.Count(s) {
				break
			}
		}
	}

	return detectors.Best()
}

type latinDetector struct {
	scriptCounter
}

func (ld *latinDetector) Count(s string) bool {
	return ld.count(s, isLatin)
}

func (ld *latinDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, latinLangs)
}

type cyrillicDetector struct {
	scriptCounter
}

func (cd *cyrillicDetector) Count(s string) bool {
	return cd.count(s, isCyrillic)
}

func (cd *cyrillicDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, cyrillicLangs)
}

type arabicDetector struct {
	scriptCounter
}

func (ad *arabicDetector) Count(s string) bool {
	return ad.count(s, isArabic)
}

func (ad *arabicDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, arabicLangs)
}

type devanagariDetector struct {
	scriptCounter
}

func (dd *devanagariDetector) Count(s string) bool {
	return dd.count(s, isDevanagari)
}

func (dd *devanagariDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, devanagariLangs)
}

type hebrewDetector struct {
	scriptCounter
}

func (hd *hebrewDetector) Count(s string) bool {
	return hd.count(s, isHebrew)
}

func (hd *hebrewDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, hebrewLangs)
}

type ethiopicDetector struct {
	scriptCounter
}

func (ed *ethiopicDetector) Count(s string) bool {
	return ed.count(s, isEthiopic)
}

func (ed *ethiopicDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, ethiopicLangs)
}

type bengaliDetector struct {
	scriptCounter
}

func (bd *bengaliDetector) Count(s string) bool {
	return bd.count(s, isBengali)
}

func (bd *bengaliDetector) Detect(text string, options Options) (Lang, float64) {
	return Ben, bd.confidence()
}

type georgianDetector struct {
	scriptCounter
}

func (gd *georgianDetector) Count(s string) bool {
	return gd.count(s, isGeorgian)
}

func (gd *georgianDetector) Detect(text string, options Options) (Lang, float64) {
	return Kat, gd.confidence()
}

type greekDetector struct {
	scriptCounter
}

func (gd *greekDetector) Count(s string) bool {
	return gd.count(s, isGreek)
}

func (gd *greekDetector) Detect(text string, options Options) (Lang, float64) {
	return Ell, gd.confidence()
}

type gujaratiDetector struct {
	scriptCounter
}

func (gd *gujaratiDetector) Count(s string) bool {
	return gd.count(s, isGujarati)
}

func (gd *gujaratiDetector) Detect(text string, options Options) (Lang, float64) {
	return Guj, gd.confidence()
}

type gurmukhiDetector struct {
	scriptCounter
}

func (gd *gurmukhiDetector) Count(s string) bool {
	return gd.count(s, isGurmukhi)
}

func (gd *gurmukhiDetector) Detect(text string, options Options) (Lang, float64) {
	return Pan, gd.confidence()
}

type kannadaDetector struct {
	scriptCounter
}

func (kd *kannadaDetector) Count(s string) bool {
	return kd.count(s, isKannada)
}

func (kd *kannadaDetector) Detect(text string, options Options) (Lang, float64) {
	return Kan, kd.confidence()
}

type khmerDetector struct {
	scriptCounter
}

func (kd *khmerDetector) Count(s string) bool {
	return kd.count(s, isKhmer)
}

func (kd *khmerDetector) Detect(text string, options Options) (Lang, float64) {
	return Khm, kd.confidence()
}

type malayalamDetector struct {
	scriptCounter
}

func (md *malayalamDetector) Count(s string) bool {
	return md.count(s, isMalayalam)
}

func (md *malayalamDetector) Detect(text string, options Options) (Lang, float64) {
	return Mal, md.confidence()
}

type myanmarDetector struct {
	scriptCounter
}

func (md *myanmarDetector) Count(s string) bool {
	return md.count(s, isMyanmar)
}

func (md *myanmarDetector) Detect(text string, options Options) (Lang, float64) {
	return Mya, md.confidence()
}

type oriyaDetector struct {
	scriptCounter
}

func (od *oriyaDetector) Count(s string) bool {
	return od.count(s, isOriya)
}

func (od *oriyaDetector) Detect(text string, options Options) (Lang, float64) {
	return Ori, od.confidence()
}

type sinhalaDetector struct {
	scriptCounter
}

func (sd *sinhalaDetector) Count(s string) bool {
	return sd.count(s, isSinhala)
}

func (sd *sinhalaDetector) Detect(text string, options Options) (Lang, float64) {
	return Sin, sd.confidence()
}

type tamilDetector struct {
	scriptCounter
}

func (td *tamilDetector) Count(s string) bool {
	return td.count(s, isTamil)
}

func (td *tamilDetector) Detect(text string, options Options) (Lang, float64) {
	return Tam, td.confidence()
}

type thaiDetector struct {
	scriptCounter
}

func (td *thaiDetector) Count(s string) bool {
	return td.count(s, isThai)
}

func (td *thaiDetector) Detect(text string, options Options) (Lang, float64) {
	return Tha, td.confidence()
}

type teluguDetector struct {
	scriptCounter
}

func (td *teluguDetector) Count(s string) bool {
	return td.count(s, isTelugu)
}

func (td *teluguDetector) Detect(text string, options Options) (Lang, float64) {
	return Tel, td.confidence()
}

type cjkDetector struct {
	scriptCounter
	c, j, k int
}

func (cjk *cjkDetector) Count(s string) bool {
	cjk.total++

	chars := 0
	for _, c := range s {
		switch {
		case isHan(c):
			cjk.c++
			chars++
		case isKana(c):
			cjk.j++
			chars++
		case isHangul(c):
			cjk.k++
			chars++
		}
	}

	if chars > 0 {
		cjk.words += chars
		cjk.chars += chars
		cjk.total += chars
	}

	return false
}

func (cjk *cjkDetector) Detect(text string, options Options) (Lang, float64) {
	var confidence float64
	if cjk.total > 0 {
		confidence = float64(cjk.words) / float64(cjk.total)
	}

	if cjk.j > 0 || cjk.k > 0 {
		if cjk.j > cjk.k {
			return Jpn, confidence
		}
		return Kor, confidence
	}
	return Zho, confidence
}
