package parser

import "reflect"

type Collection []AValue

func (c Collection) Type() Type {
  return Type(reflect.Array)
}
