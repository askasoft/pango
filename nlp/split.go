package nlp

import (
	"unicode"
	"unicode/utf8"
)

// SplitMode は改行コードの扱いを指定する動作モード
type SplitMode int

const (
	SplitForceByEOL  SplitMode = iota // 改行コードで強制文切り
	SplitIgnoreEOL                    // 改行無視（スペース扱い）
	SplitOnBlankLine                  // 空白行を挟む場合のみ文切り
)

// Sentence は切り出された文（または空白の塊）とその開始位置を保持する構造体です。
type Sentence struct {
	Text  string // 文のテキスト、または文と文の間の空白・改行
	Start int    // 原文（string）におけるバイトインデックスの開始位置
}

type SentenceSplitter struct {
	Mode               SplitMode
	MaxRunesInBracket  int // 閉じ括弧忘れに対するセーフティネット
	MaxRunesInSentence int // 0以下の数値を指定した場合、長さによる強制文切りは行われません。
}

// 欧米語の主要な省略語リスト
var abbreviations = map[string]bool{
	"Mr": true, "Ms": true, "Mrs": true, "Dr": true, "Prof": true,
	"Inc": true, "Ltd": true, "Co": true, "Corp": true,
	"e.g": true, "i.e": true, "vs": true,
}

// 開き括弧から閉じ括弧へのマッピング（左右の区別が明確なもののみ管理）
var openToClose = map[rune]rune{
	'「': '」', '『': '』', '（': '）', '(': ')', '[': ']', '{': '}',
	'〈': '〉', '《': '》', '〔': '〕', '【': '】', '“': '”', '‘': '’',
}

// 文末記号の直後に連続して、同じ文の一部として巻き込むべき記号（他の文末記号 or 閉じ括弧）
func isSkippableAfterTerminal(r rune) bool {
	switch r {
	// 文末記号のバリエーション
	case '。', '！', '？', '!', '?', '.':
		return true
	// 閉じ括弧のバリエーション
	case '」', '』', ')', ']', '}', '〉', '》', '｣', '）', '"', '’', '”':
		return true
	default:
		return false
	}
}

// SplitSentences 日本語と英語（欧米語）が混在したテキストを、言語ごとの文脈や記号の特性、特殊な表現を考慮しながら適切な「文（Sentence）」に分割する。
// ## 基本設計思想
//   - **完全な復元性（100% Round-trip）**
//     原文の文字（スペース、タブ、改行を含む）を1ビットも破棄・改変しません。文と文の間に存在する空白文字も位置情報を保持したまま独立した要素として回収するため、分割された全要素のテキストを単純結合（ `res[0].Text + res[1].Text + ...` ）するだけで、**原文と100%完全一致する文字列を再構築可能**です。
//   - **位置情報（Byte Index）の提供**
//     分割された各要素には、UTF-8エンコードにおける原文からの開始バイトインデックス（ `Start` ）が付与されます。これにより、下流のタスクでのハイライト表示やアノテーションへのマッピングが容易になります。
//
// ## 文末判定ロジック
// ### 基本文末記号
// スタックが空（括弧の外側）である場合において、以下の文字を基本文末として認識します。
//
//	全角/半角記号: 。、！、？、!、?
//	ピリオド: . （ただし、後述の「例外判定」を通過したもののみ）
//
// ### 文末記号・閉じ括弧の連続巻き込み（コンボ処理）
// 文末記号が検出された際、その直後に連続する「他の文末記号」や「閉じ括弧」は、すべて現在の文の一部として末尾に巻き込んで切り出します。
//
//	例: Really?! → ? で文切りせず、後ろの ! まで含めて Really?! を1文とする。
//	例: お祝いです！）」 → ！」 までを1つの文の末尾として扱う。
//
// ### 改行コードの動作（NewlineModeによる制御）
// 引数で指定されたモードにより、改行コード（ \n、\r、\r\n ）の扱いが以下のように切り替わります。
// モード名	識別子	振る舞いの詳細
// 強制	SplitForceByEOL	改行コードを検出した時点で、文末記号の有無に関わらず強制的に文を区切る。
// 無視	SplitIgnoreEOL	改行コードを単なる空白文字（スペースと同等）として扱い、強制文切りを行わない。
// 連続なら改行	SplitOnBlankLine	改行の直後に「文字（非空白文字）」が続く場合は結合。途中に空白行（改行やスペースのみの行）を挟む場合のみ文を区切る（Windowsの \r\n も1つの改行として適切に跨ぎます）。
//
// ## 誤分割防止（例外処理・エスケープ）
// ### 括弧（クォーテーション）のネスト管理
// 括弧の内部にある文末記号（例: 「〜。〜」 や "(~.~)"）による誤分割を防ぐため、スタックによる状態管理を行います。スタックが空（括弧の外側）の場合のみ文切りが実行されます。
//
//	左右が明確な括弧（ペア管理）:
//	「」、『』、（）、()、[]、{}、〈〉、《》、〔〕、【】、“”、‘’
//
//	半角ダブルクォーテーション（ " ）のトグル処理:
//	左右の区別がない半角 " は、スタックの末尾を確認し、すでに " があれば「閉じる」、なければ「開く」というトグル制御で内部判定を行います（これにより引用符内のピリオド等での誤分割を防ぎます）。
//
//	アポストロフィ（ ' ）の除外:
//	英語の縮約形（ I'm や don't など）での誤判定を防ぐため、半角のシングルクォーテーション（アポストロフィ）はスタック管理の対象外（常に無視）とします。
//
//	セーフティネット（閉じ忘れ対策）:
//	閉じ括弧の書き忘れによるスタック崩壊（以降すべての文切りがスキップされる現象）を防ぐため、1つの文の文字数が 400文字 を超えた場合は強制的にスタックをクリアします。
//
// ### ピリオド（ . ）に特化した例外判定
// 英語圏のテキストにおける以下のケースでは、ピリオドが検出されても文末とはみなされません。
//
//	小数点: 前後が半角数字で挟まれている場合（例: 3.14 ）。
//	主要な省略語: 直前の単語が定義された省略語リスト（ Mr, Ms, Dr, Inc, Ltd, Co, Corp, e.g, i.e, vs など）に一致する場合。
//	URL / ドメイン: ピリオド（および後ろに続く閉じ括弧・文末記号）の直後に、スペースや改行を挟まらずすぐ次の文字が始まっている場合（例: google.com ）。
//
// ## 特殊なケースにおける挙動
//
//	先頭/末尾の空白:
//	テキストの最先頭、および最後尾にある空白や改行の塊も、それぞれ独立した空白要素として安全に切り出され、位置情報が保持されます。
//
//	文直後のスペース:
//	文が確定した直後に連続するスペースやタブは、文そのものには含めず、次の要素（空白要素）として分離します。これにより、文のテキストから余分なトリム処理を行う必要がなくなります。
//
// ## `maxLen` に 0 以下の数値を指定した場合、長さによる強制文切りは行われません。
func (ss *SentenceSplitter) SplitSentences(text string) []Sentence {
	if text == "" {
		return nil
	}

	sentences := make([]Sentence, 0, 16)
	length := len(text)
	start := 0
	i := 0

	// 先頭にある空白の塊を独立した要素として回収
	for i < length {
		r, size := utf8.DecodeRuneInString(text[i:])
		if unicode.IsSpace(r) {
			i += size
		} else {
			break
		}
	}
	if i > start {
		sentences = append(sentences, Sentence{Text: text[start:i], Start: start})
		start = i
	}

	// 括弧のネスト状態を管理するスタック
	bracketStack := make([]rune, 0, 8)
	sentenceRunes := 0

	for i < length {
		r, size := utf8.DecodeRuneInString(text[i:])
		nextIdx := i + size
		sentenceRunes++

		// セーフティネット
		if ss.MaxRunesInBracket > 0 && sentenceRunes > ss.MaxRunesInBracket {
			bracketStack = bracketStack[:0]
		} else {
			// --- 括弧のネスト状態の更新 ---
			if expectedClose, isOpening := openToClose[r]; isOpening {
				bracketStack = append(bracketStack, expectedClose)
			} else if r == '"' {
				// 半角ダブルクォートは同一記号のため、スタックの末尾を確認して開閉をトグル処理
				if len(bracketStack) > 0 && bracketStack[len(bracketStack)-1] == r {
					bracketStack = bracketStack[:len(bracketStack)-1]
				} else {
					bracketStack = append(bracketStack, r)
				}
			} else if len(bracketStack) > 0 && r == bracketStack[len(bracketStack)-1] {
				bracketStack = bracketStack[:len(bracketStack)-1]
			}
		}

		isTerminal := false
		isBreakByWhitespace := false // 改行等の空白による区切りか
		isForceSplitByLen := false   // 最大長さ超過による強制区切りか

		// --- 最大長さの判定ロジック ---
		if ss.MaxRunesInSentence > 0 && sentenceRunes >= ss.MaxRunesInSentence {
			isTerminal = true
			isForceSplitByLen = true
		}

		// 文末記号の判定（括弧の外側にいる場合、かつ最大長さによる強制文切りでない場合のみ）
		if !isTerminal && len(bracketStack) == 0 {
			if r == '。' || r == '！' || r == '？' || r == '!' || r == '?' {
				isTerminal = true
			} else if r == '\n' || r == '\r' {
				// --- 改行コードの動作モード判定 ---
				if ss.Mode == SplitForceByEOL {
					isTerminal = true
					isBreakByWhitespace = true
				} else if ss.Mode == SplitOnBlankLine {
					// 先読みして、次の非空白文字が出る前に「別の改行」があるか（＝空行か）をチェック
					hasBlankLine := false
					idx := nextIdx

					// Windowsの改行コード (\r\n) を1つの改行として跨ぐ処理
					if r == '\r' && idx < length {
						nr, nsize := utf8.DecodeRuneInString(text[idx:])
						if nr == '\n' {
							idx += nsize
						}
					}

					for idx < length {
						nr, nsize := utf8.DecodeRuneInString(text[idx:])
						if nr == '\n' || nr == '\r' {
							hasBlankLine = true
							break
						}
						if !unicode.IsSpace(nr) {
							break // 通常の文字が出たら空行ではない
						}
						idx += nsize
					}

					if hasBlankLine {
						isTerminal = true
						isBreakByWhitespace = true
					}
				}
				// SplitIgnoreEOL の場合は何もせず、通常の文字（スペース同等）として処理を流す
			} else if r == '.' {
				isTerminal = true

				// 小数点対策 (例: 3.14)
				if i > 0 {
					prevR, _ := utf8.DecodeLastRuneInString(text[:i])
					if unicode.IsDigit(prevR) && nextIdx < length {
						nextR, _ := utf8.DecodeRuneInString(text[nextIdx:])
						if unicode.IsDigit(nextR) {
							isTerminal = false
						}
					}
				}

				// 略語対策 (例: Mr. / e.g.)
				if isTerminal && i > start {
					j := i - 1
					for j >= start {
						cb := text[j]
						if (cb >= 'a' && cb <= 'z') || (cb >= 'A' && cb <= 'Z') || cb == '.' || cb == '-' {
							j--
						} else {
							break
						}
					}
					word := text[j+1 : i]
					if abbreviations[word] {
						isTerminal = false
					}
				}

				// URL/ドメイン対策 (例: google.com)
				if isTerminal {
					checkIdx := nextIdx
					for checkIdx < length {
						nr, nsize := utf8.DecodeRuneInString(text[checkIdx:])
						// ここも連続する記号をスキップするように変更
						if isSkippableAfterTerminal(nr) {
							checkIdx += nsize
						} else {
							break
						}
					}
					if checkIdx < length {
						afterR, _ := utf8.DecodeRuneInString(text[checkIdx:])
						if !unicode.IsSpace(afterR) {
							isTerminal = false
						}
					}
				}
			}
		}

		// 文切り確定処理
		if isTerminal {
			endIdx := nextIdx
			if isForceSplitByLen {
				// 最大長さによる強制区切りの場合は、現在の文字（nextIdx）までをそのまま文とする
				if start < endIdx {
					sentences = append(sentences, Sentence{Text: text[start:endIdx], Start: start})
				}
			} else if !isBreakByWhitespace {
				// 文末記号の直後に連続する「他の文末記号」や「閉じ括弧」をすべて同じ文に巻き込む
				for endIdx < length {
					nr, nsize := utf8.DecodeRuneInString(text[endIdx:])
					if isSkippableAfterTerminal(nr) {
						endIdx += nsize
					} else {
						break
					}
				}
				if start < endIdx {
					sentences = append(sentences, Sentence{Text: text[start:endIdx], Start: start})
				}
			} else {
				// 改行による区切りの場合、改行文字自体は「空白の塊」に含めるため手前で切る
				if start < i {
					sentences = append(sentences, Sentence{Text: text[start:i], Start: start})
				}
				endIdx = i // 空白スキャンを改行文字の位置から開始
			}

			// 文直後の空白（スペース・タブ・改行）を塊として1要素にする
			spaceEnd := endIdx
			for spaceEnd < length {
				nr, nsize := utf8.DecodeRuneInString(text[spaceEnd:])
				if unicode.IsSpace(nr) {
					spaceEnd += nsize
				} else {
					break
				}
			}

			if endIdx < spaceEnd {
				sentences = append(sentences, Sentence{Text: text[endIdx:spaceEnd], Start: endIdx})
			}

			start = spaceEnd
			i = spaceEnd
			bracketStack = bracketStack[:0]
			sentenceRunes = 0
			continue
		}

		i += size
	}

	// 最後に残った文字列の回収
	if start < length {
		// 文末記号がないまま終了した場合でも、末尾の空白があればきれいに分離する
		endTextIdx := length
		for endTextIdx > start {
			r, size := utf8.DecodeLastRuneInString(text[start:endTextIdx])
			if unicode.IsSpace(r) {
				endTextIdx -= size
			} else {
				break
			}
		}
		if endTextIdx > start {
			sentences = append(sentences, Sentence{Text: text[start:endTextIdx], Start: start})
		}
		if endTextIdx < length {
			sentences = append(sentences, Sentence{Text: text[endTextIdx:], Start: endTextIdx})
		}
	}

	return sentences
}
