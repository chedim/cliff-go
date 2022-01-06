package parser

import "reflect"

type Integer int64

func (i Integer) Type() Type {
  return Type(reflect.Int64)
}

func (i Integer) Add(o AValue) AValue {
	return i + o.(Integer)
}

func (i Integer) Sub(o AValue) AValue {
  return i - o.(Integer)
}

func (i Integer) Mul(o AValue) AValue {
  return i * o.(Integer)
}

func (i Integer) Div(o AValue) AValue {
  return i / o.(Integer).Integer()
}

func (i Integer) Float() Float {
  return Float(i)
}

func (i Integer) Integer() Integer {
	return i
}

func (i Integer) CompareTo(o AValue) Integer {
  ov := o.(Integer).Integer()
  if i > ov {
    return 1
  } else if i < ov {
    return -1
  }
  return 0
}
