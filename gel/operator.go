package gel

import (
	"errors"
	"fmt"

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
	String() string
	Category() int
	Operands() int
	Priority() int
	Wrap(op operator, operand cog.Queue[any]) error

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

func (sop *singleOp) Operands() int {
	return 1
}

func (sop *singleOp) Wrap(op operator, operand cog.Queue[any]) error {
	var ok bool

	sop.right, ok = operand.Poll()
	if !ok {
		return fmt.Errorf("gel: operator %q missing right operand", op)
	}

	return nil
}

func (sop *singleOp) calcRight(ec elCtx) (any, error) {
	return sop.calcItem(ec, sop.right)
}

func (sop *singleOp) IsReturnNull(ec elCtx, rval any) (bool, error) {
	if rval == nil {
		if ec.Strict {
			return false, fmt.Errorf("gel: operator %q right object is nil", sop)
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

func (dop *doubleOp) Operands() int {
	return 2
}

func (dop *doubleOp) Wrap(op operator, rpn cog.Queue[any]) error {
	var ok bool

	dop.right, ok = rpn.Poll()
	if !ok {
		return fmt.Errorf("gel: operator %q missing left and right operand", op)
	}

	dop.left, ok = rpn.Poll()
	if !ok {
		return fmt.Errorf("gel: operator %q missing left or right operand", op)
	}

	return nil
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
