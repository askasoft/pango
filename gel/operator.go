package gel

import (
	"errors"

	"github.com/askasoft/pango/cog"
)

type invoker interface {
	Invoke(elCtx, []any) (any, error)
}

const (
	opOpen = 1 + iota
	opClose
	opBits
	opLogic
	opMath
	opObject
)

type operator interface {
	Category() int

	Operator() string

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
