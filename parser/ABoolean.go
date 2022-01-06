package parser

import "reflect"

type ABoolean interface {
  AValue
  Orable
  Andable
  Verifiable
  Negatable
}

type Negatable interface {
  AValue
  Negate() Negatable
}

type Andable interface {
  AValue
  And(other Verifiable) Verifiable
}

type Orable interface {
  AValue
  Or(other Verifiable) Verifiable
}

type Verifiable interface {
  IsTrue() Bool
}

type Bool bool

func NewBooleanValue(v bool) ABoolean {
	return Bool(v)
}

func (b Bool) Type() Type {
	return Type(reflect.Bool)
}

func (b Bool) And(o Verifiable) Verifiable {
  return b && o.IsTrue()
}

func (b Bool) Or(o Verifiable) Verifiable {
  return b || o.IsTrue()
}

func (b Bool) IsTrue() Bool {
  return b
}

func (b Bool) Negate() Negatable {
  return !b
}
