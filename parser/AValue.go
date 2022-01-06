package parser;

import "reflect"

type Type reflect.Kind
type AValue interface {
	Type() Type
}

func (t Type) Type() Type {
  return Type(reflect.Interface)
}
