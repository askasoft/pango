package ldt

// Lang represents a language following ISO 639-3 standard.
type Lang int

// Aka ...
const (
	Unknown Lang = iota
	Afr
	Aka
	Amh
	Arb
	Azj
	Bel
	Ben
	Bho
	Bul
	Ceb
	Ces
	Dan
	Deu
	Ell
	Eng
	Epo
	Est
	Fin
	Fra
	Guj
	Hat
	Hau
	Heb
	Hin
	Hrv
	Hun
	Ibo
	Ilo
	Ind
	Ita
	Jav
	Jpn
	Kan
	Kat
	Khm
	Kin
	Kor
	Kur
	Lav
	Lit
	Mai
	Mal
	Mar
	Mkd
	Mlg
	Mya
	Nep
	Nld
	Nno
	Nob
	Nya
	Ori
	Orm
	Pan
	Pes
	Pol
	Por
	Ron
	Run
	Rus
	Sin
	Skr
	Slv
	Sna
	Som
	Spa
	Srp
	Swe
	Tam
	Tel
	Tgl
	Tha
	Tir
	Tuk
	Tur
	Uig
	Ukr
	Urd
	Uzb
	Vie
	Ydd
	Yor
	Zho
	Zul
)

// CodeToLang gets enum by ISO 639-3 code as a string.
func CodeToLang(code string) Lang {
	lang := map[string]Lang{
		"afr": Afr,
		"aka": Aka,
		"amh": Amh,
		"arb": Arb,
		"azj": Azj,
		"bel": Bel,
		"ben": Ben,
		"bho": Bho,
		"bul": Bul,
		"ceb": Ceb,
		"ces": Ces,
		"dan": Dan,
		"deu": Deu,
		"ell": Ell,
		"eng": Eng,
		"epo": Epo,
		"est": Est,
		"fin": Fin,
		"fra": Fra,
		"guj": Guj,
		"hat": Hat,
		"hau": Hau,
		"heb": Heb,
		"hin": Hin,
		"hrv": Hrv,
		"hun": Hun,
		"ibo": Ibo,
		"ilo": Ilo,
		"ind": Ind,
		"ita": Ita,
		"jav": Jav,
		"jpn": Jpn,
		"kan": Kan,
		"kat": Kat,
		"khm": Khm,
		"kin": Kin,
		"kor": Kor,
		"kur": Kur,
		"lav": Lav,
		"lit": Lit,
		"mai": Mai,
		"mal": Mal,
		"mar": Mar,
		"mkd": Mkd,
		"mlg": Mlg,
		"mya": Mya,
		"nep": Nep,
		"nld": Nld,
		"nno": Nno,
		"nob": Nob,
		"nya": Nya,
		"ori": Ori,
		"orm": Orm,
		"pan": Pan,
		"pes": Pes,
		"pol": Pol,
		"por": Por,
		"ron": Ron,
		"run": Run,
		"rus": Rus,
		"sin": Sin,
		"skr": Skr,
		"slv": Slv,
		"sna": Sna,
		"som": Som,
		"spa": Spa,
		"srp": Srp,
		"swe": Swe,
		"tam": Tam,
		"tel": Tel,
		"tgl": Tgl,
		"tha": Tha,
		"tir": Tir,
		"tuk": Tuk,
		"tur": Tur,
		"uig": Uig,
		"ukr": Ukr,
		"urd": Urd,
		"uzb": Uzb,
		"vie": Vie,
		"ydd": Ydd,
		"yor": Yor,
		"zho": Zho,
		"zul": Zul,
	}

	if val, ok := lang[code]; ok {
		return val
	}

	return Unknown
}

// Iso6391 returns ISO 639-1 code of Lang as a string.
func (lang Lang) Iso6391() string {
	return LangToStringShort(lang)
}

// LangToStringShort converts enum into ISO 639-1 code as a string.
// Return empty string when there is no ISO 639-1 code.
// Deprecated: LangToStringShort is deprected and exists for historical
// compatibility. Please use `Lang.Iso6391()` instead.
func LangToStringShort(lang Lang) string {
	langMap := map[Lang]string{
		Afr: "af",
		Aka: "ak",
		Amh: "am",
		Arb: "ar",
		Azj: "az", // Azerbaijani iso 639-3 is aze, iso 639-1 az
		Bel: "be",
		Ben: "bn",
		Bho: "bh",
		Bul: "bg",
		Ceb: "", // No iso 639-1 code
		Ces: "cs",
		Dan: "da",
		Deu: "de",
		Ell: "el",
		Eng: "en",
		Epo: "eo",
		Est: "et",
		Fin: "fi",
		Fra: "fr",
		Guj: "gu",
		Hat: "ht",
		Hau: "ha",
		Heb: "he",
		Hin: "hi",
		Hrv: "hr",
		Hun: "hu",
		Ibo: "ig",
		Ilo: "", // No iso639-1
		Ind: "id",
		Ita: "it",
		Jav: "jv",
		Jpn: "ja",
		Kan: "kn",
		Kat: "ka",
		Khm: "km",
		Kin: "rw",
		Kor: "ko",
		Kur: "ku",
		Lav: "lv",
		Lit: "lt",
		Mai: "", // No iso639-1
		Mal: "ml",
		Mar: "mr",
		Mkd: "mk",
		Mlg: "mg",
		Mya: "my",
		Nep: "ne",
		Nld: "nl",
		Nno: "nn",
		Nob: "nb",
		Nya: "ny",
		Ori: "or",
		Orm: "om",
		Pan: "pa",
		Pes: "", // No iso639-1
		Pol: "pl",
		Por: "pt",
		Ron: "ro",
		Run: "rn",
		Rus: "ru",
		Sin: "si",
		Skr: "", // No iso639-1
		Slv: "sl",
		Sna: "sn",
		Som: "so",
		Spa: "es",
		Srp: "sr",
		Swe: "sv",
		Tam: "ta",
		Tel: "te",
		Tgl: "tl",
		Tha: "th",
		Tir: "ti",
		Tuk: "tk",
		Tur: "tr",
		Uig: "ug",
		Ukr: "uk",
		Urd: "ur",
		Uzb: "uz",
		Vie: "vi",
		Ydd: "", // No iso639-1
		Yor: "yo",
		Zho: "zh", // No iso 639-1, but http://www.loc.gov/standards/iso639-2/faq.html#24
		Zul: "zu",
	}

	if val, ok := langMap[lang]; ok {
		return val
	}

	return ""
}

// Iso6393 returns ISO 639-3 code of Lang as a string.
func (lang Lang) Iso6393() string {
	return LangToString(lang)
}

// LangToString converts enum into ISO 639-3 code as a string.
// Deprecated: LangToString is deprected and exists for historical
// compatibility. Please use `Lang.Iso6393()` instead.
func LangToString(lang Lang) string {
	langMap := map[Lang]string{
		Afr: "afr",
		Aka: "aka",
		Amh: "amh",
		Arb: "arb",
		Azj: "azj",
		Bel: "bel",
		Ben: "ben",
		Bho: "bho",
		Bul: "bul",
		Ceb: "ceb",
		Ces: "ces",
		Dan: "dan",
		Deu: "deu",
		Ell: "ell",
		Eng: "eng",
		Epo: "epo",
		Est: "est",
		Fin: "fin",
		Fra: "fra",
		Guj: "guj",
		Hat: "hat",
		Hau: "hau",
		Heb: "heb",
		Hin: "hin",
		Hrv: "hrv",
		Hun: "hun",
		Ibo: "ibo",
		Ilo: "ilo",
		Ind: "ind",
		Ita: "ita",
		Jav: "jav",
		Jpn: "jpn",
		Kan: "kan",
		Kat: "kat",
		Khm: "khm",
		Kin: "kin",
		Kor: "kor",
		Kur: "kur",
		Lav: "lav",
		Lit: "lit",
		Mai: "mai",
		Mal: "mal",
		Mar: "mar",
		Mkd: "mkd",
		Mlg: "mlg",
		Mya: "mya",
		Nep: "nep",
		Nld: "nld",
		Nno: "nno",
		Nob: "nob",
		Nya: "nya",
		Ori: "ori",
		Orm: "orm",
		Pan: "pan",
		Pes: "pes",
		Pol: "pol",
		Por: "por",
		Ron: "ron",
		Run: "run",
		Rus: "rus",
		Sin: "sin",
		Skr: "skr",
		Slv: "slv",
		Sna: "sna",
		Som: "som",
		Spa: "spa",
		Srp: "srp",
		Swe: "swe",
		Tam: "tam",
		Tel: "tel",
		Tgl: "tgl",
		Tha: "tha",
		Tir: "tir",
		Tuk: "tuk",
		Tur: "tur",
		Uig: "uig",
		Ukr: "ukr",
		Urd: "urd",
		Uzb: "uzb",
		Vie: "vie",
		Ydd: "ydd",
		Yor: "yor",
		Zho: "zho",
		Zul: "zul",
	}

	if val, ok := langMap[lang]; ok {
		return val
	}

	return ""
}

// String returns the human-readable name of the language as a string.
func (lang Lang) String() string {
	if val, ok := Langs[lang]; ok {
		return val
	}

	return ""
}

// Langs represents a map of Lang to language name.
var Langs = map[Lang]string{
	Afr: "Afrikaans",
	Aka: "Akan",
	Amh: "Amharic",
	Arb: "Arabic",
	Azj: "Azerbaijani",
	Bel: "Belarusian",
	Ben: "Bengali",
	Bho: "Bhojpuri",
	Bul: "Bulgarian",
	Ceb: "Cebuano",
	Ces: "Czech",
	Dan: "Danish",
	Deu: "German",
	Ell: "Greek",
	Eng: "English",
	Epo: "Esperanto",
	Est: "Estonian",
	Fin: "Finnish",
	Fra: "French",
	Guj: "Gujarati",
	Hat: "Haitian Creole",
	Hau: "Hausa",
	Heb: "Hebrew",
	Hin: "Hindi",
	Hrv: "Croatian",
	Hun: "Hungarian",
	Ibo: "Igbo",
	Ilo: "Ilocano",
	Ind: "Indonesian",
	Ita: "Italian",
	Jav: "Javanese",
	Jpn: "Japanese",
	Kan: "Kannada",
	Kat: "Georgian",
	Khm: "Khmer",
	Kin: "Kinyarwanda",
	Kor: "Korean",
	Kur: "Kurdish",
	Lav: "Latvian",
	Lit: "Lithuanian",
	Mai: "Maithili",
	Mal: "Malayalam",
	Mar: "Marathi",
	Mkd: "Macedonian",
	Mlg: "Malagasy",
	Mya: "Burmese",
	Nep: "Nepali",
	Nld: "Dutch",
	Nno: "Nynorsk",
	Nob: "Bokmal",
	Nya: "Chewa",
	Ori: "Oriya",
	Orm: "Oromo",
	Pan: "Punjabi",
	Pes: "Persian",
	Pol: "Polish",
	Por: "Portuguese",
	Ron: "Romanian",
	Run: "Rundi",
	Rus: "Russian",
	Sin: "Sinhalese",
	Skr: "Saraiki",
	Slv: "Slovene",
	Sna: "Shona",
	Som: "Somali",
	Spa: "Spanish",
	Srp: "Serbian",
	Swe: "Swedish",
	Tam: "Tamil",
	Tel: "Telugu",
	Tgl: "Tagalog",
	Tha: "Thai",
	Tir: "Tigrinya",
	Tuk: "Turkmen",
	Tur: "Turkish",
	Uig: "Uyghur",
	Ukr: "Ukrainian",
	Urd: "Urdu",
	Uzb: "Uzbek",
	Vie: "Vietnamese",
	Ydd: "Yiddish",
	Yor: "Yoruba",
	Zho: "Chinese",
	Zul: "Zulu",
}
