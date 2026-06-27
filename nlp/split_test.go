package nlp

import (
	"reflect"
	"testing"
)

func TestSplitSentencesForceByEOL(t *testing.T) {
	cs := []struct {
		s string
		w []Sentence
	}{
		// 空文字・スペースのみのケース
		{"", nil},
		{"   ", []Sentence{{"   ", 0}}},

		// 通常の単一文章（句点・感嘆符・疑問符）
		{"こんにちは。", []Sentence{
			{`こんにちは。`, 0},
		}},
		{"Hello world!", []Sentence{
			{`Hello world!`, 0},
		}},

		// 複数の文章（連続するケース）
		{"最初の文章です。二番目の文章です。", []Sentence{
			{`最初の文章です。`, 0},
			{`二番目の文章です。`, 24}, // 「最初の文章です。」は24バイト
		}},
		{"Hello! How are you?", []Sentence{
			{`Hello!`, 0},
			{` `, 6},
			{`How are you?`, 7}, // 半角スペースを挟んで7バイト目から開始
		}},

		// 改行が含まれるケース
		{"一行目。\n二行目。", []Sentence{
			{`一行目。`, 0},
			{"\n", 12},
			{`二行目。`, 13}, // 「一行目。\n」のバイト数
		}},

		// 末尾に句切れ文字がないケース
		{"句点のない文章です", []Sentence{
			{`句点のない文章です`, 0},
		}},

		// 連続する句切れマーク（ModeForceの挙動検証用）
		{"本当に！？そう思う。", []Sentence{
			{`本当に！？`, 0},
			{`そう思う。`, 15},
		}},
		{"Go v1.26 release.", []Sentence{
			{`Go v1.26 release.`, 0}, // 小数点と句点の区別
		}},

		{"What...?", []Sentence{{`What...?`, 0}}},
		{"Really?!", []Sentence{{`Really?!`, 0}}},
		{"Hello\nI'm tom.", []Sentence{{`Hello`, 0}, {"\n", 5}, {`I'm tom.`, 6}}},
		{`I'm going next week. Get ready.`, []Sentence{{`I'm going next week.`, 0}, {` `, 20}, {`Get ready.`, 21}}},
		{`He said, "I'm going next week. Get ready. "`, []Sentence{{`He said, "I'm going next week. Get ready. "`, 0}}},
		{`彼は「来週行くよ。準備してね」と言った。`, []Sentence{{`彼は「来週行くよ。準備してね」と言った。`, 0}}},
		{`彼は"来週行くよ。準備してね"と言った。`, []Sentence{{`彼は"来週行くよ。準備してね"と言った。`, 0}}},
	}

	ss := SentenceSplitter{Mode: SplitForceByEOL}
	for i, c := range cs {
		a := ss.SplitSentences(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("#%d: SplitSentencesForceByEOL(%q) = %v\nWANT: %v", i, c.s, a, c.w)
		}
	}
}

func TestSplitSentencesIgnoreEOL(t *testing.T) {
	cs := []struct {
		s string
		w []Sentence
	}{
		// 空文字・スペースのみのケース
		{"", nil},
		{"   ", []Sentence{{"   ", 0}}},

		{"こんにちは。元気ですか？バグはありませんか？", []Sentence{
			{`こんにちは。`, 0},
			{`元気ですか？`, 18},
			{`バグはありませんか？`, 36},
		}},

		// 改行が含まれていても分割されないこと
		{"一行目の文章、\n二行目の文章です。", []Sentence{
			{"一行目の文章、\n二行目の文章です。", 0},
		}},
	}

	ss := SentenceSplitter{Mode: SplitIgnoreEOL}
	for i, c := range cs {
		a := ss.SplitSentences(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("#%d: SplitSentencesIgnoreEOL(%q) = %v\nWANT: %v", i, c.s, a, c.w)
		}
	}
}

func TestSplitSentencesOnBlankLine(t *testing.T) {
	cs := []struct {
		s string
		w []Sentence
	}{
		// 単一の改行では分割されないこと
		{"第一句。第二句、\n同じ段落内の新しい行です。", []Sentence{
			{`第一句。`, 0},
			{"第二句、\n同じ段落内の新しい行です。", 12},
		}},

		// 空行（\n\n）によって綺麗に2つのブロックに分割されること
		// 「セクション1。」＝15バイト ＋ 「\n\n」＝2バイト ＝ 計17バイト
		{"セクション1。\n\nセクション2。", []Sentence{
			{`セクション1。`, 0},
			{"\n\n", 19},
			{`セクション2。`, 21},
		}},

		// 連続する空行や、スペースが含まれる空行のケース（ロバスト性の検証）
		// 「最初」＝6バイト ＋ 「\n\n\n」＝3バイト ＝ 計9バイト
		{"最初\n\n\n最後", []Sentence{
			{`最初`, 0},
			{"\n\n\n", 6},
			{`最後`, 9},
		}},

		// 文末に空行があるケース
		{"末尾に空行あり\n\n", []Sentence{
			{`末尾に空行あり`, 0},
			{"\n\n", 21},
		}},
	}

	ss := SentenceSplitter{Mode: SplitOnBlankLine}
	for i, c := range cs {
		a := ss.SplitSentences(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("#%d: SplitSentencesOnBlankLine(%q) = %v\nWANT: %v", i, c.s, a, c.w)
		}
	}
}
