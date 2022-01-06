package parser;

type BinaryOperatorExpression struct {
  span Span
  arguments []AValue
  operator BinaryOperator
}

type Addable interface {
  AValue
  Add(other AValue) AValue;
}

type Subtractable interface {
  AValue
  Sub(other AValue) AValue
}

type Multipliable interface {
  AValue
  Mul(by AValue) AValue
}

type Dividable interface {
  AValue
  Div(by AValue) AValue
}

type Comparable interface {
  AValue
  CompareTo(o AValue) Integer
}

type BinaryOperator func(left AValue, right AValue) AValue

func NewBinaryOperator(span Span, left AValue, right AValue, operator BinaryOperator) AnExpression {
  return &BinaryOperatorExpression{
    span: span,
    arguments: []AValue{left.(Addable), right.(Addable)},
    operator: operator,
  }
}

func (o *BinaryOperatorExpression) Span() *Span {
  return &o.span
}

func (o *BinaryOperatorExpression) Value() AValue {
  return o.operator(o.arguments[0], o.arguments[1])
}

func (o *BinaryOperatorExpression) Arguments() *[]AValue {
  return &o.arguments
}
