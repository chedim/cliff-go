package parser

import "fmt"

type BinaryExpression struct {
	span     Span
	Left     AnExpression
	Right    AnExpression
	Operator BinaryOperator
}

func NewBinaryExpression(left AnExpression, right AnExpression, operator BinaryOperator) AnExpression {
	return &BinaryExpression{
		span:     *left.Span().Extend(right.Span()),
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (b *BinaryExpression) Span() *Span {
	return &b.span
}

func (b *BinaryExpression) Value() AValue {
	return b.Operator(b.Left.Value(), b.Right.Value())
}

func (b *BinaryExpression) String() string {
  return fmt.Sprintf("%s %v %s", b.Left.String(), b.Operator, b.Right.String())
}
