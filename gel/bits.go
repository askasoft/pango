package gel

import "github.com/askasoft/pango/cal"

type bitNot struct {
	singleOp
}

func (bn *bitNot) Category() int {
	return opBits
}

func (bn *bitNot) Operator() string {
	return "~"
}

func (bn *bitNot) Priority() int {
	return 2
}

func (bn *bitNot) Calculate(ec elCtx) (any, error) {
	rval, err := bn.calcRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := bn.IsReturnNull(ec, rval); ok || err != nil {
		return nil, err
	}

	return cal.BitNot(rval)
}

type bitAnd struct {
	doubleOp
}

func (ba *bitAnd) Category() int {
	return opBits
}

func (ba *bitAnd) Operator() string {
	return "&"
}

func (ba *bitAnd) Priority() int {
	return 7
}

func (ba *bitAnd) Calculate(ec elCtx) (any, error) {
	lval, rval, err := ba.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := ba.IsReturnNull(ec, lval, rval); ok || err != nil {
		return nil, err
	}

	return cal.BitAnd(lval, rval)
}

type bitOr struct {
	doubleOp
}

func (bo *bitOr) Category() int {
	return opBits
}

func (bo *bitOr) Operator() string {
	return "|"
}

func (bo *bitOr) Priority() int {
	return 9
}

func (bo *bitOr) Calculate(ec elCtx) (any, error) {
	lval, rval, err := bo.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := bo.IsReturnNull(ec, lval, rval); ok || err != nil {
		return nil, err
	}

	return cal.BitOr(lval, rval)
}

type bitXor struct {
	doubleOp
}

func (bx *bitXor) Category() int {
	return opBits
}

func (bx *bitXor) Operator() string {
	return "^"
}

func (bx *bitXor) Priority() int {
	return 8
}

func (bx *bitXor) Calculate(ec elCtx) (any, error) {
	lval, rval, err := bx.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := bx.IsReturnNull(ec, lval, rval); ok || err != nil {
		return nil, err
	}

	return cal.BitXor(lval, rval)
}

type bitLeft struct {
	doubleOp
}

func (bl *bitLeft) Category() int {
	return opBits
}

func (bl *bitLeft) Operator() string {
	return "<<"
}

func (bl *bitLeft) Priority() int {
	return 5
}

func (bl *bitLeft) Calculate(ec elCtx) (any, error) {
	lval, rval, err := bl.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := bl.IsReturnNull(ec, lval, rval); ok || err != nil {
		return nil, err
	}

	return cal.BitLeft(lval, rval)
}

type bitRight struct {
	doubleOp
}

func (br *bitRight) Category() int {
	return opBits
}

func (br *bitRight) Operator() string {
	return ">>"
}

func (br *bitRight) Priority() int {
	return 5
}

func (br *bitRight) Calculate(ec elCtx) (any, error) {
	lval, rval, err := br.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := br.IsReturnNull(ec, lval, rval); ok || err != nil {
		return nil, err
	}

	return cal.BitRight(lval, rval)
}
