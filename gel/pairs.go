package gel

import (
	"errors"

	"github.com/askasoft/pango/cog"
)

type openOp struct {
}

func (o *openOp) Category() int {
	return opOpen
}

func (o *openOp) Priority() int {
	return 100
}

func (o *openOp) Wrap(cog.Queue[any]) {
}

type closeOp struct {
}

func (o *closeOp) Category() int {
	return opClose
}

func (o *closeOp) Priority() int {
	return 100
}

func (o *closeOp) Wrap(cog.Queue[any]) {
}

var lbraceop = &lBraceOp{}

type lBraceOp struct {
	openOp
}

func (lb *lBraceOp) Operator() string {
	return "{"
}

func (lb *lBraceOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '{'")
}

var rbraceop = &rBraceOp{}

type rBraceOp struct {
	closeOp
}

func (rb *rBraceOp) Operator() string {
	return "}"
}

func (rb *rBraceOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '}'")
}

var lparenthesisop = &lParenthesisOp{}

type lParenthesisOp struct {
	openOp
}

func (lp *lParenthesisOp) Operator() string {
	return "("
}

func (lp *lParenthesisOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '('")
}

var rparenthesisop = &rParenthesisOp{}

type rParenthesisOp struct {
	closeOp
}

func (rp *rParenthesisOp) Operator() string {
	return "("
}

func (rp *rParenthesisOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by ')'")
}

var lbracketop = &lBracketOp{}

type lBracketOp struct {
	openOp
}

func (lb *lBracketOp) Operator() string {
	return "("
}

func (lb *lBracketOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '['")
}

var rbracketop = &rBracketOp{}

type rBracketOp struct {
	closeOp
}

func (rb *rBracketOp) Operator() string {
	return "("
}

func (rb *rBracketOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by ']'")
}
