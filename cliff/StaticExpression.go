package cliff;


type StaticExpression struct {
  span *Span
  value AValue
}

func NewStaticExpression(span *Span, value AValue) *StaticExpression {
  return &StaticExpression{span, value}
}

func (s *StaticExpression) Span() *Span {
  return s.span
}

func (s *StaticExpression) Value() AValue {
  return s.value
}
