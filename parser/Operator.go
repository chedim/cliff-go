package parser;

type Operator interface {
  Span() *Span
  Value() AValue
  Arguments() *[]AValue
}
