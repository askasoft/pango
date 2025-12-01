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

func (o *openOp) Operands() int {
	return 0
}

func (o *openOp) Priority() int {
	return 100
}

func (o *openOp) Wrap(operator, cog.Queue[any]) error {
	return nil
}

type closeOp struct {
}

func (o *closeOp) Category() int {
	return opClose
}

func (o *closeOp) Operands() int {
	return 0
}

func (o *closeOp) Priority() int {
	return 100
}

func (o *closeOp) Wrap(operator, cog.Queue[any]) error {
	return nil
}

var lbraceop = &lBraceOp{}

type lBraceOp struct {
	openOp
}

func (lb *lBraceOp) String() string {
	return "{"
}

func (lb *lBraceOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '{'")
}

var rbraceop = &rBraceOp{}

type rBraceOp struct {
	closeOp
}

func (rb *rBraceOp) String() string {
	return "}"
}

func (rb *rBraceOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '}'")
}

var lparenthesisop = &lParenthesisOp{}

type lParenthesisOp struct {
	openOp
}

func (lp *lParenthesisOp) String() string {
	return "("
}

func (lp *lParenthesisOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '('")
}

var rparenthesisop = &rParenthesisOp{}

type rParenthesisOp struct {
	closeOp
}

func (rp *rParenthesisOp) String() string {
	return "("
}

func (rp *rParenthesisOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by ')'")
}

var lbracketop = &lBracketOp{}

type lBracketOp struct {
	openOp
}

func (lb *lBracketOp) String() string {
	return "("
}

func (lb *lBracketOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by '['")
}

var rbracketop = &rBracketOp{}

type rBracketOp struct {
	closeOp
}

func (rb *rBracketOp) String() string {
	return "("
}

func (rb *rBracketOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: Calculate() is unsupported by ']'")
}
