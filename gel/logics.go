package gel

import (
	"errors"
	"regexp"

	"github.com/askasoft/pango/cal"
	"github.com/askasoft/pango/cas"
	"github.com/askasoft/pango/ref"
)

type logicNot struct {
	singleOp
}

func (ln *logicNot) Category() int {
	return opLogic
}

func (ln *logicNot) String() string {
	return "!"
}

func (ln *logicNot) Priority() int {
	return 2
}

func (ln *logicNot) Calculate(ec elCtx) (any, error) {
	rval, err := ln.calcRight(ec)
	if err != nil {
		return nil, err
	}

	return ref.IsZero(rval), nil
}

type logicAnd struct {
	doubleOp
}

func (la *logicAnd) Category() int {
	return opLogic
}

func (la *logicAnd) String() string {
	return "&&"
}

func (la *logicAnd) Priority() int {
	return 11
}

func (la *logicAnd) Calculate(ec elCtx) (any, error) {
	lval, rval, err := la.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	return cal.LogicAnd(lval, rval), nil
}

type logicOr struct {
	doubleOp
}

func (lo *logicOr) Category() int {
	return opLogic
}

func (lo *logicOr) String() string {
	return "||"
}

func (lo *logicOr) Priority() int {
	return 12
}

func (lo *logicOr) Calculate(ec elCtx) (any, error) {
	lval, rval, err := lo.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	return cal.LogicOr(lval, rval), nil
}

type logicNilable struct {
	singleOp
}

func (ln *logicNilable) Category() int {
	return opLogic
}

func (ln *logicNilable) String() string {
	return "@"
}

func (ln *logicNilable) Priority() int {
	return 2
}

func (ln *logicNilable) Calculate(ec elCtx) (any, error) {
	rval, _ := ln.calcRight(ec)
	return rval, nil
}

type logicOrable struct {
	doubleOp
}

func (lo *logicOrable) Category() int {
	return opLogic
}

func (lo *logicOrable) String() string {
	return "|||"
}

func (lo *logicOrable) Priority() int {
	return 12
}

func (lo *logicOrable) Calculate(ec elCtx) (any, error) {
	lval, rval, err := lo.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if !ref.IsZero(lval) {
		return lval, nil
	}

	return rval, nil
}

type logicEq struct {
	doubleOp
}

func (le *logicEq) Category() int {
	return opLogic
}

func (le *logicEq) String() string {
	return "=="
}

func (le *logicEq) Priority() int {
	return 6
}

func (le *logicEq) Calculate(ec elCtx) (any, error) {
	lval, rval, err := le.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	return cal.LogicEq(lval, rval)
}

type logicNeq struct {
	doubleOp
}

func (ln *logicNeq) Category() int {
	return opLogic
}

func (ln *logicNeq) String() string {
	return "!="
}

func (ln *logicNeq) Priority() int {
	return 6
}

func (ln *logicNeq) Calculate(ec elCtx) (any, error) {
	lval, rval, err := ln.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	return cal.LogicNeq(lval, rval)
}

type logicGt struct {
	doubleOp
}

func (lg *logicGt) Category() int {
	return opLogic
}

func (lg *logicGt) String() string {
	return ">"
}

func (lg *logicGt) Priority() int {
	return 6
}

func (lg *logicGt) Calculate(ec elCtx) (any, error) {
	lval, rval, err := lg.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	return cal.LogicGt(lval, rval)
}

type logicGte struct {
	doubleOp
}

func (lg *logicGte) Category() int {
	return opLogic
}

func (lg *logicGte) String() string {
	return ">="
}

func (lg *logicGte) Priority() int {
	return 6
}

func (lg *logicGte) Calculate(ec elCtx) (any, error) {
	lval, rval, err := lg.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	return cal.LogicGte(lval, rval)
}

type logicLt struct {
	doubleOp
}

func (lt *logicLt) Category() int {
	return opLogic
}

func (lt *logicLt) String() string {
	return "<"
}

func (lt *logicLt) Priority() int {
	return 6
}

func (lt *logicLt) Calculate(ec elCtx) (any, error) {
	lval, rval, err := lt.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	return cal.LogicLt(lval, rval)
}

type logicLte struct {
	doubleOp
}

func (lt *logicLte) Category() int {
	return opLogic
}

func (lt *logicLte) String() string {
	return "<="
}

func (lt *logicLte) Priority() int {
	return 6
}

func (lt *logicLte) Calculate(ec elCtx) (any, error) {
	lval, rval, err := lt.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	return cal.LogicLte(lval, rval)
}

type logicRem struct {
	doubleOp
}

func (lr *logicRem) Category() int {
	return opLogic
}

func (lr *logicRem) String() string {
	return "=="
}

func (lr *logicRem) Priority() int {
	return 6
}

func (lr *logicRem) Calculate(ec elCtx) (any, error) {
	lval, rval, err := lr.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	lstr, err := cas.ToString(lval)
	if err != nil {
		return nil, err
	}

	rstr, err := cas.ToString(rval)
	if err != nil {
		return nil, err
	}

	rex, err := regexp.Compile(rstr)
	if err != nil {
		return nil, err
	}

	return rex.MatchString(lstr), nil
}

type logicQuestion struct {
	doubleOp
}

func (lq *logicQuestion) Category() int {
	return opLogic
}

func (lq *logicQuestion) String() string {
	return "?"
}

func (lq *logicQuestion) Priority() int {
	return 13
}

func (lq *logicQuestion) Calculate(ec elCtx) (any, error) {
	lval, err := lq.calcLeft(ec)
	if err != nil {
		return nil, err
	}

	return !ref.IsZero(lval), nil
}

type logicQuestionSelect struct {
	doubleOp
}

func (lqs *logicQuestionSelect) Category() int {
	return opLogic
}

func (lqs *logicQuestionSelect) String() string {
	return ":"
}

func (lqs *logicQuestionSelect) Priority() int {
	return 13
}

func (lqs *logicQuestionSelect) Calculate(ec elCtx) (any, error) {
	if lq, ok := lqs.left.(*logicQuestion); ok {
		lv, err := lq.Calculate(ec)
		if err != nil {
			return nil, err
		}

		if lv.(bool) {
			return lq.calcRight(ec)
		}

		return lqs.calcRight(ec)
	}

	return nil, errors.New("gel: invalid ternary operator")
}
