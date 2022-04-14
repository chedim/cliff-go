package parser

import (
	"reflect"
)

type String string

func (s String) Type() Type {
	return Type(reflect.String)
}

func (s String) String() string {
  return string(s)
}

func (s String) Equals(o AValue) ABoolean {
  return NewBooleanValue(s.Value() == o.Value())
}

func (s String) Value() interface{} {
  return string(s)
}
