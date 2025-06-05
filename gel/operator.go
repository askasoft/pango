package gel

import (
	"errors"

	"github.com/askasoft/pango/cog"
)

type invoker interface {
	Invoke(elCtx, []any) (any, error)
}

type operator interface {
	Priority() int

	Wrap(operand cog.Queue[any])

	Calculate(elCtx) (any, error)
}

type op struct {
}

func (op *op) calcItem(ec elCtx, obj any) (any, error) {
	if obj == nil {
		return nil, nil
	}

	if op, ok := obj.(operator); ok {
		return op.Calculate(ec)
	}

	if eo, ok := obj.(elObj); ok {
		return eo.Get(ec)
	}

	return obj, nil
}

type singleOp struct {
	op
	right any
}

func (sop *singleOp) Wrap(operand cog.Queue[any]) {
	sop.right, _ = operand.Poll()
}

func (sop *singleOp) calcRight(ec elCtx) (any, error) {
	return sop.calcItem(ec, sop.right)
}

func (sop *singleOp) IsReturnNull(ec elCtx, rval any) (bool, error) {
	if rval == nil {
		if ec.Strict {
			return false, errors.New("gel: right object is nil")
		}
		return true, nil
	}

	return false, nil
}

type doubleOp struct {
	op
	right any
	left  any
}

func (dop *doubleOp) Wrap(rpn cog.Queue[any]) {
	dop.right, _ = rpn.Poll()
	dop.left, _ = rpn.Poll()
}

func (dop *doubleOp) calcLeft(ec elCtx) (any, error) {
	return dop.calcItem(ec, dop.left)
}

func (dop *doubleOp) calcRight(ec elCtx) (any, error) {
	return dop.calcItem(ec, dop.right)
}

func (dop *doubleOp) calcLeftRight(ec elCtx) (lval, rval any, err error) {
	lval, err = dop.calcLeft(ec)
	if err != nil {
		return
	}
	rval, err = dop.calcRight(ec)
	return
}

func (dop *doubleOp) evalLeft(ec elCtx) (any, error) {
	if dop.left == nil {
		return ec.Object, nil
	}

	if op, ok := dop.left.(operator); ok {
		return op.Calculate(ec)
	}
	if eo, ok := dop.left.(elObj); ok {
		return eo.Get(ec)
	}
	return dop.left, nil
}

func (dop *doubleOp) IsReturnNull(ec elCtx, lval, rval any) (bool, error) {
	if lval == nil {
		if ec.Strict {
			return false, errors.New("gel: left object is nil")
		}
		return true, nil
	}

	if rval == nil {
		if ec.Strict {
			return false, errors.New("gel: right object is nil")
		}
		return true, nil
	}

	return false, nil
}

var lbraceop = &lBraceOp{}

type lBraceOp struct {
}

func (lb *lBraceOp) Operator() string {
	return "{"
}

func (lb *lBraceOp) Priority() int {
	return 100
}

func (lb *lBraceOp) Wrap(cog.Queue[any]) {
}

func (lb *lBraceOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '{'")
}

var rbraceop = &rBraceOp{}

type rBraceOp struct {
}

func (rb *rBraceOp) Operator() string {
	return "}"
}

func (rb *rBraceOp) Priority() int {
	return 100
}

func (rb *rBraceOp) Wrap(cog.Queue[any]) {
}

func (rb *rBraceOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '}'")
}

var lparenthesisop = &lParenthesisOp{}

type lParenthesisOp struct {
}

func (lp *lParenthesisOp) Operator() string {
	return "("
}

func (lp *lParenthesisOp) Priority() int {
	return 100
}

func (lp *lParenthesisOp) Wrap(cog.Queue[any]) {
}

func (lp *lParenthesisOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '('")
}

var rparenthesisop = &rParenthesisOp{}

type rParenthesisOp struct {
}

func (rp *rParenthesisOp) Operator() string {
	return "("
}

func (rp *rParenthesisOp) Priority() int {
	return 100
}

func (rp *rParenthesisOp) Wrap(cog.Queue[any]) {
}

func (rp *rParenthesisOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by ')'")
}

var lbracketop = &lBracketOp{}

type lBracketOp struct {
}

func (lb *lBracketOp) Operator() string {
	return "("
}

func (lb *lBracketOp) Priority() int {
	return 100
}

func (lb *lBracketOp) Wrap(cog.Queue[any]) {
}

func (lb *lBracketOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '['")
}

var rbracketop = &rBracketOp{}

type rBracketOp struct {
}

func (rb *rBracketOp) Operator() string {
	return "("
}

func (rb *rBracketOp) Priority() int {
	return 100
}

func (rb *rBracketOp) Wrap(cog.Queue[any]) {
}

func (rb *rBracketOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by ']'")
}
