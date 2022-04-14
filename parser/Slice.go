package parser

import (
	"fmt"
	"reflect"
)

type ASlice []AValue

type SliceExpression struct {
	location  Span
	datapoint Reference
	isHead    bool
  isNext    bool
	capacity  NumberExpression
}

func ReadSliceExpression(scanner *Scanner) (*SliceExpression, *ParserError) {
  result := SliceExpression{
    location: *scanner.Position(),
  }
  if tok := scanner.Peek(); tok.Token != THE {
    return nil, NewParserError(*scanner.Position(), fmt.Sprintf("tried to read slice at %s", tok.Token))
  }
  scanner.Scan()
  scanner.scanWhitespace()

  switch first := scanner.Peek(); first.Token {
    case FIRST:
      result.isHead = true
      scanner.Scan()
      scanner.scanWhitespace()
    case NEXT:
      result.isNext = true
      fallthrough
    case LAST:
      scanner.Scan()
      scanner.scanWhitespace()
      if capacity, e := ReadNumber(scanner); e == nil {
        result.capacity = *capacity
      } else {
        return nil, e
      }
    default:
      break;
  }


  if datapoint, e := ReadReference(scanner); e != nil {
    return nil, e
  } else {
    result.datapoint = *datapoint
  }

  result.location.Extend(scanner.Position())

  return &result, nil
}

func (se *SliceExpression) Span() *Span {
	return &se.location
}

func (se *SliceExpression) Type() Type {
	return Type(reflect.Array)
}

func (se *SliceExpression) Value() AValue {
	return se.Slice(se.capacity.Value().Value().(int))
}

func (se *SliceExpression) Slice(max int) (result ASlice) {
  return se.datapoint.Target().Slice(max)
}

func (se *SliceExpression) Target() (*Reference) {
  return &se.datapoint
}

func (se *SliceExpression) String() string {
  if se.isHead {
    return fmt.Sprintf("%s[0:%s]", se.Target().String(), se.capacity.String())
  }
  return fmt.Sprintf("%s[-%s:]", se.Target().String(), se.capacity.String())
}

func (se *SliceExpression) Equals(o AValue) ABoolean {
  return NewBooleanValue(se.Value() == o.Value())
}

func (se ASlice) Value() interface{} {
  return []AValue(se)
}

func (se ASlice) Equals(o AValue) ABoolean {
  return NewBooleanValue(se.Value() == o.Value())
}

func (se ASlice) String() string {
  return fmt.Sprint(([]AValue(se)))
}

func (se ASlice) Type() Type {
  return Type(reflect.Array)
}
