package parser

import "reflect"

type String string

func (s String) Type() Type {
	return Type(reflect.String)
}

