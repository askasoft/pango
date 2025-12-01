package gel

import (
	"github.com/askasoft/pango/cal"
)

type mathPositive struct {
	singleOp
}

func (mp mathPositive) Category() int {
	return opMath
}

func (mp mathPositive) String() string {
	return "+"
}

func (mp mathPositive) Priority() int {
	return 2
}

func (mp mathPositive) Calculate(ec elCtx) (any, error) {
	return mp.calcRight(ec)
}

type mathNegate struct {
	singleOp
}

func (mn mathNegate) Category() int {
	return opMath
}

func (mn mathNegate) String() string {
	return "-"
}

func (mn mathNegate) Priority() int {
	return 2
}

func (mn mathNegate) Calculate(ec elCtx) (any, error) {
	rval, err := mn.calcRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := mn.IsReturnNull(ec, rval); ok || err != nil {
		return nil, err
	}

	return cal.Negate(rval)
}

type mathAdd struct {
	doubleOp
}

func (ma mathAdd) Category() int {
	return opMath
}

func (ma mathAdd) String() string {
	return "+"
}

func (ma mathAdd) Priority() int {
	return 4
}

func (ma mathAdd) Calculate(ec elCtx) (any, error) {
	lval, rval, err := ma.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	ok, err := ma.IsReturnNull(ec, lval, rval)
	if err != nil {
		return nil, err
	}

	if ok {
		if lval == nil && rval == nil {
			return nil, nil
		}
		if lval == nil {
			return rval, nil
		}
		return lval, nil
	}

	return cal.Add(lval, rval)
}

type mathSub struct {
	doubleOp
}

func (ms mathSub) Category() int {
	return opMath
}

func (ms mathSub) String() string {
	return "-"
}

func (ms mathSub) Priority() int {
	return 4
}

func (ms mathSub) Calculate(ec elCtx) (any, error) {
	lval, rval, err := ms.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	ok, err := ms.IsReturnNull(ec, lval, rval)
	if err != nil {
		return ok, err
	}

	if ok {
		if lval == nil && rval == nil {
			return nil, nil
		}
		if rval == nil {
			return lval, nil
		}
		return cal.Negate(rval)
	}

	return cal.Sub(lval, rval)
}

type mathMul struct {
	doubleOp
}

func (mm mathMul) Category() int {
	return opMath
}

func (mm mathMul) String() string {
	return "*"
}

func (mm mathMul) Priority() int {
	return 3
}

func (mm mathMul) Calculate(ec elCtx) (any, error) {
	lval, rval, err := mm.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := mm.IsReturnNull(ec, lval, rval); ok || err != nil {
		return nil, err
	}

	return cal.Multiply(lval, rval)
}

type mathDiv struct {
	doubleOp
}

func (md mathDiv) Category() int {
	return opMath
}

func (md mathDiv) String() string {
	return "/"
}

func (md mathDiv) Priority() int {
	return 3
}

func (md mathDiv) Calculate(ec elCtx) (any, error) {
	lval, rval, err := md.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := md.IsReturnNull(ec, lval, rval); ok || err != nil {
		return nil, err
	}

	return cal.Divide(lval, rval)
}

type mathMod struct {
	doubleOp
}

func (mm mathMod) Category() int {
	return opMath
}

func (mm mathMod) String() string {
	return "%"
}

func (mm mathMod) Priority() int {
	return 3
}

func (mm mathMod) Calculate(ec elCtx) (any, error) {
	lval, rval, err := mm.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if ok, err := mm.IsReturnNull(ec, lval, rval); ok || err != nil {
		return nil, err
	}

	return cal.Mod(lval, rval)
}
