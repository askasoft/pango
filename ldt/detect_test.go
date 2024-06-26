package ldt

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		s string
		w Info
	}{
		{"愛している I Love You.", Info{Jpn, 0.5}},
		{"Además de todo lo anteriormente dicho, también encontramos...", Info{Spa, 0.6}},
		{"बहुत बहुत (धन्यवाद / शुक्रिया)!", Info{Hin, 0.7}},
		{"अनुच्छेद १, सबहि लोकानि आजादे जम्मेला आओर ओखिनियो के बराबर सम्मान आओर अघ्कार प्राप्त हवे। ओखिनियो के पास समझ-बूझ आओर अंत,करण के आवाज होखता आओर हुनको के दोसरा के साथ भाईचारे के बेवहार करे के होखला", Info{Bho, 1}},
		{"ኢትዮጵያ አፍሪቃ ውስጥ ናት", Info{Amh, 1}},
		{"لغتي العربية ليست كما يجب", Info{Arb, 1}},
		{"我爱你", Info{Zho, 1}},
		{"আমি তোমাকে ভালোবাস ", Info{Ben, 1}},
		{"울란바토르", Info{Kor, 1}},
		{"ყველა ადამიანი იბადება თავისუფალი და თანასწორი თავისი ღირსებითა და უფლებებით", Info{Kat, 1}},
		{"Όλοι οι άνθρωποι γεννιούνται ελεύθεροι και ίσοι στην αξιοπρέπεια και τα δικαιώματα.", Info{Ell, 1}},
		{"ಎಲ್ಲಾ ಮಾನವರ ಉಚಿತ ಮತ್ತು ಘನತೆ ಮತ್ತು ಹಕ್ಕುಗಳಲ್ಲಿ ಸಮಾನ ಹುಟ್ಟಿದ.", Info{Kan, 1}},
		{"நீங்கள் ஆங்கிலம் பேசுவீர்களா?", Info{Tam, 1}},
		{"มนุษย์ทุกคนเกิดมามีอิสระและเสมอภาคกันในศักดิ์ศรีและสิทธิ", Info{Tha, 1}},
		{"નાણાં મારા લોહીમાં છે", Info{Guj, 1}},
		{" ਗੁਰੂ ਗ੍ਰੰਥ ਸਾਹਿਬ ਜੀ", Info{Pan, 1}},
		{"నన్ను ఒంటరిగా వదిలేయ్", Info{Tel, 1}},
		{"എന്താണ് നിങ്ങളുടെ പേര് ?", Info{Mal, 1}},
		{"ମୁ ତୁମକୁ ଭଲ ପାଏ |", Info{Ori, 1}},
		{"အားလုံးလူသားတွေအခမဲ့နှင့်ဂုဏ်သိက္ခာနှင့်လူ့အခွင့်အရေးအတွက်တန်းတူဖွားမြင်ကြသည်။", Info{Mya, 1}},
		{"වෙලාව කියද?", Info{Sin, 1}},
		{"ពួកម៉ាកខ្ញុំពីរនាក់នេះ", Info{Khm, 1}},
		{"其疾如風、其徐如林、侵掠如火、不動如山、難知如陰、動如雷震。", Info{Zho, 1}},
		{"知彼知己、百戰不殆。不知彼而知己、一勝一負。不知彼不知己、毎戰必殆。", Info{Zho, 1}},
		{"支那の上海の或町です。", Info{Jpn, 1}},
		{"或日の暮方の事である。", Info{Jpn, 1}},
		{"今日は", Info{Jpn, 1}},
		{"コンニチハ", Info{Jpn, 1}},
		{"ﾀﾅｶ ﾀﾛｳ", Info{Jpn, 1}},
		{"どうもありがとう", Info{Jpn, 1}},
	}

	for i, c := range tests {
		a := Detect(c.s)

		if c.w.Lang != a.Lang || a.Confidence-c.w.Confidence >= 0.1 {
			t.Errorf("#%d Detect(%q) = %v %v, WANT %v %v", i, c.s, LangToString(a.Lang), a.Confidence, LangToString(c.w.Lang), c.w.Confidence)
		}
	}
}

func TestDetectLang(t *testing.T) {
	tests := map[string]Lang{
		"Та нічого, все нормально. А в тебе як?": Ukr,
		"Vouloir, c'est pouvoir":                                Fra,
		"Where there is a will there is a way":                  Eng,
		"Mi ŝategas la japanan kaj studas ĝin kelkajn jarojn 😊": Epo,
		"Te echo de menos":                                      Spa,
		"Buona notte e sogni d'oro!":                            Ita,
	}

	for text, want := range tests {
		got := DetectLang(text)
		if got != want {
			t.Fatalf("%s want %v got %v", text, LangToString(want), LangToString(got))
		}
	}
}

// Test detect with empty options and supported language and script
func TestDetectWithOptionsEmptySupportedLang(t *testing.T) {
	want := Info{Epo, 0.5}
	got := DetectWithOptions("La viro amas hundojn. Hundo estas la plej bona amiko de viro", Options{})
	if got.Lang != want.Lang || got.Confidence-want.Confidence >= 0.1 {
		t.Fatalf("want %v %v got %v %v", want.Lang, want.Confidence, got.Lang, got.Confidence)
	}
}

// Test detect with empty options and nonsupported script(Balinese)
func TestDetectWithOptionsEmptyNonSupportedLang(t *testing.T) {
	want := Info{UNKNOWN, 0}
	got := DetectWithOptions("ᬅᬓ᭄ᬱᬭᬯ᭄ᬬᬜ᭄ᬚᬦ", Options{})
	if got.Lang != want.Lang || got.Confidence-want.Confidence >= 0.1 {
		t.Fatalf("want %v %v got %v %v", want.Lang, want.Confidence, got.Lang, got.Confidence)
	}
}

func TestDetectWithOptionsWithBlacklist(t *testing.T) {
	text := "האקדמיה ללשון העברית"

	//All languages with Hebrew text blacklisted ... returns correct script but invalid language
	options1 := Options{
		Blacklist: map[Lang]bool{
			Heb: true,
			Ydd: true,
		},
	}

	want := Info{UNKNOWN, 0}
	got := DetectWithOptions(text, options1)
	if got.Lang != want.Lang || got.Confidence-want.Confidence >= 0.1 {
		t.Fatalf("Want %v %v got %v %v", LangToString(want.Lang), want.Confidence, LangToString(got.Lang), got.Confidence)
	}

	text = "Tu me manques"
	want = Info{Ilo, 0}
	options3 := Options{
		Blacklist: map[Lang]bool{
			Kur: true,
		},
	}
	got = DetectWithOptions(text, options3)
	if got.Lang != want.Lang || got.Confidence-want.Confidence >= 0.1 {
		t.Fatalf("Want %v %v got %v %v", LangToString(want.Lang), want.Confidence, LangToString(got.Lang), got.Confidence)
	}
}

func TestWithOptionsWithWhitelist(t *testing.T) {
	text := "Mi ne scias!"
	want := Info{Epo, 1}
	options2 := Options{
		Whitelist: map[Lang]bool{
			Epo: true,
			Ukr: true,
		},
	}
	got := DetectWithOptions(text, options2)
	if got.Lang != want.Lang || want.Confidence != got.Confidence {
		t.Fatalf("Want %v %v got %v %v", LangToString(want.Lang), want.Confidence, LangToString(got.Lang), got.Confidence)
	}
}

func TestDetectLangWithOptions(t *testing.T) {
	text := "All evil come from a single cause ... man's inability to sit still in a room"
	want := Eng
	//without blacklist
	got := DetectLangWithOptions(text, Options{})
	if want != got {
		t.Fatalf("want %s got %s", LangToString(want), LangToString(got))
	}

	//with blacklist
	options := Options{
		Blacklist: map[Lang]bool{
			Jav: true,
			Tgl: true,
			Nld: true,
			Uzb: true,
			Swe: true,
			Nob: true,
			Ceb: true,
			Ilo: true,
		},
	}

	got = DetectLangWithOptions(text, options)
	if want != got {
		t.Fatalf("want %s got %s", LangToString(want), LangToString(got))
	}
}

func TestWithMultipleExamples(t *testing.T) {
	examplesFile, err := os.Open("testdata/examples.json")
	if err != nil {
		t.Fatal("Error opening testdata/examples.json")
	}

	defer examplesFile.Close()

	byteValue, err := io.ReadAll(examplesFile)
	if err != nil {
		t.Fatal("Error reading testdata/examples.json")
	}

	var examples map[string]string
	err = json.Unmarshal(byteValue, &examples)
	if err != nil {
		t.Fatal("Error Unmarshalling json")
	}

	for lang, text := range examples {
		want := CodeToLang(lang)
		info := Detect(text)
		if info.Lang != want && !info.IsReliable() {
			t.Fatalf("want %v, got %v", Langs[want], Langs[info.Lang])
		}
	}
}
