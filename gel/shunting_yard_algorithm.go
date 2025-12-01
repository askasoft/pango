package gel

import (
	"errors"

	"github.com/askasoft/pango/cog/linkedlist"
)

// Shunting yard算法是一个用于将中缀表达式转换为后缀表达式的经典算法，由艾兹格·迪杰斯特拉引入，因其操作类似于火车编组场而得名。
// @see https://en.wikipedia.org/wiki/Shunting_yard_algorithm
type shuntingYard struct {
	ops linkedlist.LinkedList[operator]
	rpn linkedlist.LinkedList[any]
}

func (sy *shuntingYard) addOperator(op operator) error {
	switch op.(type) {
	case *lParenthesisOp, *lBracketOp, *lBraceOp: // '(' or '[' or '{'
		sy.ops.PushHead(op)
		return nil
	case *rParenthesisOp: // ')'
		for {
			if op, ok := sy.ops.PollHead(); ok {
				if _, ok := op.(*lParenthesisOp); ok {
					return nil
				}
				sy.rpn.Add(op)
			} else {
				return errors.New("gel: missing open bracket '('")
			}
		}
	case *rBracketOp: // ']'
		for {
			if op, ok := sy.ops.PollHead(); ok {
				if _, ok := op.(*lBracketOp); ok {
					return nil
				}
				sy.rpn.Add(op)
			} else {
				return errors.New("gel: missing open bracket '['")
			}
		}
	case *rBraceOp: // '}'
		for {
			if op, ok := sy.ops.PollHead(); ok {
				if _, ok := op.(*lBraceOp); ok {
					return nil
				}
				sy.rpn.Add(op)
			} else {
				return errors.New("gel: missing open brace '{'")
			}
		}
	}

	// 空,直接添加进操作符队列
	if sy.ops.IsEmpty() {
		sy.ops.PushHead(op)
		return nil
	}

	oh, ok := sy.ops.PeekHead()
	if ok {
		// 符号队列top元素优先级大于当前,则直接添加到
		if oh.Priority() > op.Priority() {
			sy.ops.PushHead(op)
			return nil
		}

		// for !!a
		if oh.Priority() == op.Priority() && oh.Operands() == op.Operands() && oh.Operands() == 1 {
			sy.ops.PushHead(op)
			return nil
		}

		// 一般情况,即优先级小于栈顶,那么直接弹出来,添加到逆波兰表达式中
		for ; ok && oh.Priority() <= op.Priority(); oh, ok = sy.ops.PeekHead() {
			// 三元表达式嵌套的特殊处理
			if _, ok := oh.(*logicQuestion); ok {
				if _, ok := op.(*logicQuestion); ok {
					break
				}
				if _, ok := op.(*logicQuestionSelect); ok {
					sy.ops.PollHead()
					sy.rpn.Add(oh)
					break
				}
			}

			sy.ops.PollHead()
			sy.rpn.Add(oh)
		}
	}

	sy.ops.PushHead(op)
	return nil
}

func (sy *shuntingYard) parseToRPN(expr string) error {
	sy.rpn.Clear()
	sy.ops.Clear()

	items, err := parse(expr)
	if err != nil {
		return err
	}

	for it := items.Iterator(); it.Next(); {
		item := it.Value()
		if op, ok := item.(operator); ok {
			if err := sy.addOperator(op); err != nil {
				return err
			}
			continue
		}
		sy.rpn.Add(item)
	}

	for {
		if op, ok := sy.ops.PollHead(); ok {
			sy.rpn.Add(op)
		} else {
			break
		}
	}

	return nil
}

// ParseToRPN Parse expression to Reverse Polish notation (RPN).
// 逆波兰表示法(逆波兰记法)是一种是由波兰数学家扬·武卡谢维奇1920年引入的数学表达式方式。
// 在逆波兰记法中，所有操作符置于操作数的后面，因此也被称为后缀表示法。
// @see https://en.wikipedia.org/wiki/Reverse_Polish_notation
func (sy *shuntingYard) ParseToRPN(expr string) (any, error) {
	if err := sy.parseToRPN(expr); err != nil {
		return nil, err
	}

	operand := linkedlist.NewLinkedList[any]()
	for it := sy.rpn.Iterator(); it.Next(); {
		obj := it.Value()
		if op, ok := obj.(operator); ok {
			if err := op.Wrap(op, operand); err != nil {
				return nil, err
			}
		}
		operand.PushHead(obj)
	}

	root, _ := operand.PeekHead()
	return root, nil
}
