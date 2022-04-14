package parser

import (
	"fmt"
	"reflect"
)

type Float float64

func (f Float) Type() Type {
	return Type(reflect.Float64)
}

func (f Float) Add(o AValue) AValue {
  return f + o.(Float)
}

func (f Float) Sub(o AValue) AValue {
  return f - o.(Float)
}

func (f Float) Mul(o AValue) AValue {
  return f * o.(Float)
}

func (f Float) Div(o AValue) AValue {
  return f / o.(Float)
}

func (f Float) Float() Float {
  return f
}

func (f Float) CompareTo(o AValue) Integer {
  fo, ok := o.(Float)
  if !ok {
    io, ok := o.(Integer)
    if !ok {
      panic("unsupported comparable")
    }
    fo = Float(io)
  }

  if (f > fo) {
    return 1
  } else if (f < fo) {
    return -1
  }
  return 0
}

func (f Float) String() string {
  return fmt.Sprintf("%f", f)
}

func (f Float) Value() interface{} {
  return float64(f)
}

func (f Float) Equals(o AValue) ABoolean {
  return NewBooleanValue(f.Value() == o.Value())
}
