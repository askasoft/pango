package gel

var (
	_ operator = &lBraceOp{}
	_ operator = &rBraceOp{}
	_ operator = &lParenthesisOp{}
	_ operator = &rParenthesisOp{}
	_ operator = &lBracketOp{}
	_ operator = &rBracketOp{}

	_ operator = &bitNot{}
	_ operator = &bitAnd{}
	_ operator = &bitOr{}
	_ operator = &bitXor{}
	_ operator = &bitLeft{}
	_ operator = &bitRight{}

	_ operator = &logicNot{}
	_ operator = &logicAnd{}
	_ operator = &logicOr{}
	_ operator = &logicNilable{}
	_ operator = &logicOrable{}
	_ operator = &logicEq{}
	_ operator = &logicNeq{}
	_ operator = &logicGt{}
	_ operator = &logicGte{}
	_ operator = &logicLt{}
	_ operator = &logicLte{}
	_ operator = &logicRem{}
	_ operator = &logicQuestion{}
	_ operator = &logicQuestionSelect{}

	_ operator = &mathPositive{}
	_ operator = &mathNegate{}
	_ operator = &mathAdd{}
	_ operator = &mathSub{}
	_ operator = &mathMul{}
	_ operator = &mathDiv{}
	_ operator = &mathMod{}

	_ operator = &accessOp{}
	_ operator = &arrayGetOp{}
	_ operator = &arrayEndOp{}
	_ operator = &arrayMakeOp{}
	_ operator = &commaOp{}
	_ operator = &funcInvokeOp{}
	_ operator = &funcEndOp{}
)
