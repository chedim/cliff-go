package parser;

import "reflect"

type Type reflect.Kind
type AValue interface {
  Value() interface{}
	Type() Type
  String() string
  Equals(o AValue) ABoolean
}

func (t Type) Type() Type {
  return Type(reflect.Interface)
}
