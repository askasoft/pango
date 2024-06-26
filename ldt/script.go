package ldt

import (
	"strings"
)

type checker func(r rune) bool

type Detector interface {
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

type Detectors []Detector

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

func (ds Detectors) Best() Detector {
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

func AllDetectors() Detectors {
	return Detectors{
		&LatinDetector{},
		&CyrillicDetector{},
		&ArabicDetector{},
		&DevanagariDetector{},
		&EthiopicDetector{},
		&HebrewDetector{},
		&BengaliDetector{},
		&GeorgianDetector{},
		&GreekDetector{},
		&KannadaDetector{},
		&TamilDetector{},
		&ThaiDetector{},
		&GujaratiDetector{},
		&GurmukhiDetector{},
		&TeluguDetector{},
		&MalayalamDetector{},
		&OriyaDetector{},
		&MyanmarDetector{},
		&KhmerDetector{},
		&SinhalaDetector{},
		&CJKDetector{},
	}
}

// detect returns only the Detector of the given text.
func detect(text string, detectors []Detector) Detector {
	if len(detectors) == 0 {
		detectors = AllDetectors()
	}

	ss := strings.FieldsFunc(text, isStopChar)
	for _, s := range ss {
		for _, d := range detectors {
			if d.Count(s) {
				break
			}
		}
	}

	return Detectors(detectors).Best()
}

type LatinDetector struct {
	scriptCounter
}

func (ld *LatinDetector) Count(s string) bool {
	return ld.count(s, isLatin)
}

func (ld *LatinDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, latinLangs, ld.confidence())
}

type CyrillicDetector struct {
	scriptCounter
}

func (cd *CyrillicDetector) Count(s string) bool {
	return cd.count(s, isCyrillic)
}

func (cd *CyrillicDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, cyrillicLangs, cd.confidence())
}

type ArabicDetector struct {
	scriptCounter
}

func (ad *ArabicDetector) Count(s string) bool {
	return ad.count(s, isArabic)
}

func (ad *ArabicDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, arabicLangs, ad.confidence())
}

type DevanagariDetector struct {
	scriptCounter
}

func (dd *DevanagariDetector) Count(s string) bool {
	return dd.count(s, isDevanagari)
}

func (dd *DevanagariDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, devanagariLangs, dd.confidence())
}

type HebrewDetector struct {
	scriptCounter
}

func (hd *HebrewDetector) Count(s string) bool {
	return hd.count(s, isHebrew)
}

func (hd *HebrewDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, hebrewLangs, hd.confidence())
}

type EthiopicDetector struct {
	scriptCounter
}

func (ed *EthiopicDetector) Count(s string) bool {
	return ed.count(s, isEthiopic)
}

func (ed *EthiopicDetector) Detect(text string, options Options) (Lang, float64) {
	return detectLangInProfiles(text, options, ethiopicLangs, ed.confidence())
}

type BengaliDetector struct {
	scriptCounter
}

func (bd *BengaliDetector) Count(s string) bool {
	return bd.count(s, isBengali)
}

func (bd *BengaliDetector) Detect(text string, options Options) (Lang, float64) {
	return Ben, bd.confidence()
}

type GeorgianDetector struct {
	scriptCounter
}

func (gd *GeorgianDetector) Count(s string) bool {
	return gd.count(s, isGeorgian)
}

func (gd *GeorgianDetector) Detect(text string, options Options) (Lang, float64) {
	return Kat, gd.confidence()
}

type GreekDetector struct {
	scriptCounter
}

func (gd *GreekDetector) Count(s string) bool {
	return gd.count(s, isGreek)
}

func (gd *GreekDetector) Detect(text string, options Options) (Lang, float64) {
	return Ell, gd.confidence()
}

type GujaratiDetector struct {
	scriptCounter
}

func (gd *GujaratiDetector) Count(s string) bool {
	return gd.count(s, isGujarati)
}

func (gd *GujaratiDetector) Detect(text string, options Options) (Lang, float64) {
	return Guj, gd.confidence()
}

type GurmukhiDetector struct {
	scriptCounter
}

func (gd *GurmukhiDetector) Count(s string) bool {
	return gd.count(s, isGurmukhi)
}

func (gd *GurmukhiDetector) Detect(text string, options Options) (Lang, float64) {
	return Pan, gd.confidence()
}

type KannadaDetector struct {
	scriptCounter
}

func (kd *KannadaDetector) Count(s string) bool {
	return kd.count(s, isKannada)
}

func (kd *KannadaDetector) Detect(text string, options Options) (Lang, float64) {
	return Kan, kd.confidence()
}

type KhmerDetector struct {
	scriptCounter
}

func (kd *KhmerDetector) Count(s string) bool {
	return kd.count(s, isKhmer)
}

func (kd *KhmerDetector) Detect(text string, options Options) (Lang, float64) {
	return Khm, kd.confidence()
}

type MalayalamDetector struct {
	scriptCounter
}

func (md *MalayalamDetector) Count(s string) bool {
	return md.count(s, isMalayalam)
}

func (md *MalayalamDetector) Detect(text string, options Options) (Lang, float64) {
	return Mal, md.confidence()
}

type MyanmarDetector struct {
	scriptCounter
}

func (md *MyanmarDetector) Count(s string) bool {
	return md.count(s, isMyanmar)
}

func (md *MyanmarDetector) Detect(text string, options Options) (Lang, float64) {
	return Mya, md.confidence()
}

type OriyaDetector struct {
	scriptCounter
}

func (od *OriyaDetector) Count(s string) bool {
	return od.count(s, isOriya)
}

func (od *OriyaDetector) Detect(text string, options Options) (Lang, float64) {
	return Ori, od.confidence()
}

type SinhalaDetector struct {
	scriptCounter
}

func (sd *SinhalaDetector) Count(s string) bool {
	return sd.count(s, isSinhala)
}

func (sd *SinhalaDetector) Detect(text string, options Options) (Lang, float64) {
	return Sin, sd.confidence()
}

type TamilDetector struct {
	scriptCounter
}

func (td *TamilDetector) Count(s string) bool {
	return td.count(s, isTamil)
}

func (td *TamilDetector) Detect(text string, options Options) (Lang, float64) {
	return Tam, td.confidence()
}

type ThaiDetector struct {
	scriptCounter
}

func (td *ThaiDetector) Count(s string) bool {
	return td.count(s, isThai)
}

func (td *ThaiDetector) Detect(text string, options Options) (Lang, float64) {
	return Tha, td.confidence()
}

type TeluguDetector struct {
	scriptCounter
}

func (td *TeluguDetector) Count(s string) bool {
	return td.count(s, isTelugu)
}

func (td *TeluguDetector) Detect(text string, options Options) (Lang, float64) {
	return Tel, td.confidence()
}

type CJKDetector struct {
	scriptCounter
	c, j, k int
}

func (cjk *CJKDetector) Count(s string) bool {
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

func (cjk *CJKDetector) Detect(text string, options Options) (Lang, float64) {
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
