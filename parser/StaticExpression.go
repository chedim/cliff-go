package parser;


type ConstExpression struct {
  span *Span
  value AValue
}

func NewConstExpression(span *Span, value AValue) *ConstExpression {
  return &ConstExpression{span, value}
}

func (s *ConstExpression) Span() *Span {
  return s.span
}

func (s *ConstExpression) Value() AValue {
  return s.value
}

func (s *ConstExpression) String() string {
  return s.value.String()
}
