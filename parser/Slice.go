package parser

import (
	"fmt"
	"reflect"
)

type ASlice interface {
	Slice(max int) (result []AValue)
}

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
	return se
}

func (se *SliceExpression) Slice(max int) (result []AValue) {
  return se.datapoint.Target().Slice(max)
}

func (se *SliceExpression) Target() (*Reference) {
  return &se.datapoint
}
