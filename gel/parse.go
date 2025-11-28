package gel

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/askasoft/pango/cog/linkedlist"
)

var parses = []func(rune, *bufio.Reader) (any, error){parseOperator, parseString, parseIdentifier, parseNumber}

type parser struct {
	br    *bufio.Reader
	items linkedlist.LinkedList[any]
	funcs linkedlist.LinkedList[*funcInvokeOp]
}

func parse(expr string) (*linkedlist.LinkedList[any], error) {
	p := &parser{
		br: bufio.NewReader(strings.NewReader(expr)),
	}
	err := p.Parse()
	return &p.items, err
}

func (p *parser) Parse() error {
	for {
		if err := p.parseItem(); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

func (p *parser) parseItem() error {
	r, err := p.skipSpace()
	if err != nil {
		return err
	}

	for _, parse := range parses {
		obj, err := parse(r, p.br)
		if obj != nil {
			if er2 := p.addItem(obj); er2 != nil {
				return er2
			}
			return err
		}
		if err != nil {
			return err
		}
	}

	return errors.New("gel: illegal state")
}

func (p *parser) addItem(item any) (err error) {
	if mi, ok := p.funcs.PeekTail(); ok && mi != nil {
		if mi.params <= 0 {
			switch item.(type) {
			case *commaOp, *rParenthesisOp, *rBraceOp: // , ], }
			default:
				mi.params++
			}
		} else {
			if _, ok := item.(*commaOp); ok {
				mi.params++
			}
		}
	}

	prev := p.items.Tail()

	switch item.(type) {
	case *lParenthesisOp: // '('
		if _, ok := prev.(elObj); ok {
			fio := &funcInvokeOp{}
			p.funcs.Add(fio)
			p.items.Add(fio)
		} else {
			p.funcs.Add(nil)
		}
		p.items.Add(item)
	case *rParenthesisOp: // ')'
		p.items.Add(item)
		if m, ok := p.funcs.Pop(); ok && m != nil {
			p.items.Add(&funcEndOp{})
		}
	case *lBracketOp: // '['
		p.items.Add(&arrayGetOp{})
		p.items.Add(item)
	case *rBracketOp: // ']'
		p.items.Add(item)
		p.items.Add(&arrayEndOp{})
	case *lBraceOp: // '{'
		fio := &funcInvokeOp{}
		p.funcs.Add(fio)
		p.items.Add(&arrayMakeOp{})
		p.items.Add(fio)
		p.items.Add(lparenthesisop)
	case *rBraceOp: // '}'
		if m, ok := p.funcs.Pop(); !ok || m == nil {
			err = errors.New("gel: missing opening brace '{'")
		} else {
			p.items.Add(rparenthesisop)
			p.items.Add(&funcEndOp{})
		}
	case *bitXor: // '^'
		if p.isBitNotable(prev) {
			p.items.Add(&bitNot{})
		} else {
			p.items.Add(item)
		}
	case *mathSub: // '-'
		if p.isNegative(prev) {
			p.items.Add(&mathNegate{})
		} else {
			p.items.Add(item)
		}
	default:
		p.items.Add(item)
	}

	return
}

func (p *parser) isBitNotable(prev any) bool {
	if prev == nil {
		return true
	}

	if op, ok := prev.(operator); ok {
		switch op.Category() {
		case opOpen, opBits, opMath, opLogic:
			return true
		}
	}
	return false
}

func (p *parser) isNegative(prev any) bool {
	if prev == nil {
		return true
	}

	if op, ok := prev.(operator); ok {
		switch op.Category() {
		case opOpen, opBits, opMath, opLogic:
			return true
		}
	}
	return false
}

func (p *parser) skipSpace() (rune, error) {
	for {
		r, _, err := p.br.ReadRune()
		if err != nil {
			return r, err
		}

		if !unicode.IsSpace(r) {
			return r, nil
		}
	}
}

func isGoIdentifierStart(r rune) bool {
	return r == '$' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isGoIdentifierPart(r rune) bool {
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func parseOperator(r rune, br *bufio.Reader) (any, error) {
	switch r {
	case '+':
		return &mathAdd{}, nil
	case '-':
		return &mathSub{}, nil
	case '*':
		return &mathMul{}, nil
	case '/':
		return &mathDiv{}, nil
	case '%':
		return &mathMod{}, nil
	case '~':
		r2, _, err := br.ReadRune()
		if err != nil {
			return &bitNot{}, err
		}
		switch r2 {
		case '=':
			return &logicRem{}, nil
		default:
			return &bitNot{}, br.UnreadRune()
		}
	case '^':
		return &bitXor{}, nil
	case '?':
		return &logicQuestion{}, nil
	case ':':
		return &logicQuestionSelect{}, nil
	case ',':
		return &commaOp{}, nil
	case '(':
		return lparenthesisop, nil
	case ')':
		return rparenthesisop, nil
	case '[':
		return lbracketop, nil
	case ']':
		return rbracketop, nil
	case '{':
		return lbraceop, nil
	case '}':
		return rbraceop, nil
	case '>':
		r2, _, err := br.ReadRune()
		if err != nil {
			return &logicGt{}, err
		}
		switch r2 {
		case '=':
			return &logicGte{}, nil
		case '>':
			return &bitRight{}, nil
		default:
			return &logicGt{}, br.UnreadRune()
		}
	case '<':
		r2, _, err := br.ReadRune()
		if err != nil {
			return &logicLt{}, err
		}
		switch r2 {
		case '=':
			return &logicLte{}, nil
		case '<':
			return &bitLeft{}, nil
		default:
			return &logicLt{}, br.UnreadRune()
		}
	case '=':
		r2, _, err := br.ReadRune()
		if err != nil {
			return nil, errors.New("gel: incorrect expression, missing character after '='")
		}
		switch r2 {
		case '=':
			return &logicEq{}, nil
		default:
			return nil, errors.New("gel: incorrect expression, illegal character after '='")
		}
	case '!':
		r2, _, err := br.ReadRune()
		if err != nil {
			return nil, errors.New("gel: incorrect expression, missing character after '!'")
		}
		switch r2 {
		case '=':
			return &logicNeq{}, nil
		case '!':
			return &logicNilable{}, nil
		default:
			return &logicNot{}, br.UnreadRune()
		}
	case '|':
		r2, _, err := br.ReadRune()
		if err != nil {
			return nil, errors.New("gel: incorrect expression, missing character after '!'")
		}
		switch r2 {
		case '|':
			r3, _, err := br.ReadRune()
			if err != nil {
				return &logicOr{}, err
			}
			if r3 == '|' {
				return &logicOrable{}, nil
			}
			return &logicOr{}, br.UnreadRune()
		default:
			return &bitOr{}, br.UnreadRune()
		}
	case '&':
		r2, _, err := br.ReadRune()
		if err != nil {
			return &bitAnd{}, err
		}
		switch r2 {
		case '&':
			return &logicAnd{}, nil
		default:
			return &bitAnd{}, br.UnreadRune()
		}
	case '.':
		r2, _, err := br.ReadRune()
		if err != nil {
			return nil, errors.New("gel: incorrect expression, missing character after '.'")
		}

		err = br.UnreadRune()

		if !isGoIdentifierStart(r2) {
			return nil, err
		}

		return &accessOp{}, err
	default:
		return nil, nil
	}
}

func parseIdentifier(r rune, br *bufio.Reader) (any, error) {
	if !isGoIdentifierStart(r) {
		return nil, nil
	}

	var err error
	var sb strings.Builder

	for {
		sb.WriteRune(r)

		r, _, err = br.ReadRune()
		if err != nil {
			break
		}

		if !isGoIdentifierPart(r) {
			err = br.UnreadRune()
			break
		}
	}

	s := sb.String()
	if s == "nil" {
		return elObj{}, err
	}
	if s == "true" {
		return true, err
	}
	if s == "false" {
		return false, err
	}
	return elObj{s}, err
}

func parseNumber(r rune, br *bufio.Reader) (any, error) {
	var err error

	if r == '.' || (r >= '0' && r <= '9') {
		dot := r == '.'

		var sb strings.Builder
		for {
			sb.WriteRune(r)

			r, _, err = br.ReadRune()
			if err != nil {
				break
			}

			switch {
			case r == '_' || r >= '0' && r <= '9':
			case r == '.':
				if dot {
					return nil, errors.New("gel: multiple '.' in number literal")
				}
				dot = true
			case r == 'l' || r == 'L':
				n, er := strconv.ParseInt(sb.String(), 0, 64)
				if er != nil {
					return n, er
				}
				return n, err
			case r == 'f' || r == 'F':
				n, er := strconv.ParseFloat(sb.String(), 64)
				if er != nil {
					return n, er
				}
				return float32(n), err
			case r == 'd' || r == 'D':
				n, er := strconv.ParseFloat(sb.String(), 64)
				if er != nil {
					return n, er
				}
				return n, err
			default:
				err = br.UnreadRune()
				goto BREAK
			}
		}

	BREAK:
		if dot {
			n, er := strconv.ParseFloat(sb.String(), 64)
			if er != nil {
				return n, er
			}
			return n, err
		}

		n, er := strconv.ParseInt(sb.String(), 0, 64)
		if er != nil {
			return n, er
		}
		return int(n), err
	}

	return nil, nil
}

func parseString(r rune, br *bufio.Reader) (any, error) {
	var err error

	switch r {
	case '\'', '"':
		end := r

		var sb strings.Builder
		for {
			r, _, err = br.ReadRune()
			if err != nil || r == end {
				break
			}

			if r == '\\' {
				err = unescape(br, &sb)
				if err != nil {
					return nil, err
				}
			} else {
				sb.WriteRune(r)
			}
		}
		return sb.String(), err
	}

	return nil, nil
}

func unescape(br *bufio.Reader, sb *strings.Builder) error {
	r, _, err := br.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("gel: missing escape character")
		}
		return err
	}

	switch r {
	case 'b':
		sb.WriteRune(' ')
	case 'f':
		sb.WriteRune('\f')
	case 'n':
		sb.WriteRune('\n')
	case 'r':
		sb.WriteRune('\r')
	case 't':
		sb.WriteRune('\t')
	case 'v':
		sb.WriteRune('\v')
	case '\\':
		sb.WriteRune('\\')
	case '\'':
		sb.WriteRune('\'')
	case '"':
		sb.WriteRune('"')
	case 'x':
		r, err = hex2rune(br, 2)
		sb.WriteRune(r)
	case 'u':
		r, err = hex2rune(br, 4)
		sb.WriteRune(r)
	default:
		return errors.New("gel: unexpected character after \\")
	}
	return err
}

func hex2rune(br *bufio.Reader, size int) (r rune, err error) {
	var sb strings.Builder

	for range size {
		r, _, err = br.ReadRune()
		if err != nil {
			break
		}
		sb.WriteRune(r)
	}

	if err != nil && sb.Len() != size {
		return 0, errors.New("gel: invalid hex rune expression")
	}

	var n int64
	n, err = strconv.ParseInt(sb.String(), 16, 32)

	return rune(n), err
}
