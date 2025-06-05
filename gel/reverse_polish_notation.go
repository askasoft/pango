package gel

import (
	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/linkedlist"
)

// 逆波兰表示法（Reverse Polish notation，RPN，或逆波兰记法），
// 是一种是由波兰数学家扬·武卡谢维奇1920年引入的数学表达式方式，在逆波兰记法中，所有操作符置于操作数的后面，因此也被称为后缀表示法。
// @see https://en.wikipedia.org/wiki/Reverse_Polish_notation
type reversePolishNotation struct {
	ops *linkedlist.LinkedList[any]
}

func newReversePolishNotation(items cog.List[any]) (rpn reversePolishNotation) {
	rpn.ops = operatorTree(items)
	return
}

func (rpn *reversePolishNotation) Calculate(ec elCtx) (any, error) {
	if obj, ok := rpn.ops.PeekHead(); ok {
		if op, ok := obj.(operator); ok {
			return op.Calculate(ec)
		}
		if eo, ok := obj.(elObj); ok {
			return eo.Get(ec)
		}
		return obj, nil
	}
	return nil, nil
}

// 转换成操作树
func operatorTree(items cog.List[any]) *linkedlist.LinkedList[any] {
	operand := linkedlist.NewLinkedList[any]()
	for it := items.Iterator(); it.Next(); {
		obj := it.Value()
		if op, ok := obj.(operator); ok {
			op.Wrap(operand)
		}
		operand.PushHead(obj)
	}
	return operand
}
