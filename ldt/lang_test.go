package ldt

import "testing"

func TestCodeToLang(t *testing.T) {
	tests := map[string]Lang{
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
		"cmn": Cmn,
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
		"zul": Zul,
		"xxx": -1,
	}

	for codeStr, lang := range tests {
		got := CodeToLang(codeStr)
		if got != lang {
			t.Fatalf("%s: want %v got %v", codeStr, lang, got)
		}
	}
}

func TestLangToString(t *testing.T) {
	tests := map[Lang]string{
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
		Cmn: "cmn",
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
		Zul: "zul",
		-1:  "",
	}

	for lang, codeStr := range tests {
		got := lang.Iso6393()
		if got != codeStr {
			t.Fatalf("%d: want %s got %s", lang, codeStr, got)
		}
	}
}

func TestLangToStringShort(t *testing.T) {
	tests := map[Lang]string{
		Afr: "af",
		Aka: "ak",
		Amh: "am",
		Arb: "ar",
		Azj: "az",
		Bel: "be",
		Ben: "bn",
		Bho: "bh",
		Bul: "bg",
		Ceb: "",
		Ces: "cs",
		Cmn: "zh",
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
		Ilo: "",
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
		Mai: "",
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
		Pes: "",
		Pol: "pl",
		Por: "pt",
		Ron: "ro",
		Run: "rn",
		Rus: "ru",
		Sin: "si",
		Skr: "",
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
		Ydd: "",
		Yor: "yo",
		Zul: "zu",
		-1:  "",
	}

	for lang, codeStr := range tests {
		got := lang.Iso6391()
		if got != codeStr {
			t.Fatalf("%d: want %s got %s", lang, codeStr, got)
		}
	}
}
