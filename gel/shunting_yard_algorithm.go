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

	// 符号队列top元素优先级大于当前,则直接添加到
	if oh, ok := sy.ops.PeekHead(); ok && oh.Priority() > op.Priority() {
		sy.ops.PushHead(op)
		return nil
	}

	// 一般情况,即优先级小于栈顶,那么直接弹出来,添加到逆波兰表达式中
	for {
		if oh, ok := sy.ops.PeekHead(); ok && oh.Priority() <= op.Priority() {
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
			continue
		}
		break
	}

	sy.ops.PushHead(op)
	return nil
}

// 转换成 逆波兰表示法（Reverse Polish notation，RPN，或逆波兰记法）
func (sy *shuntingYard) ParseToRPN(expr string) error {
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
