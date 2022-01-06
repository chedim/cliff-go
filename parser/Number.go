package parser

import (
	"fmt"
	"reflect"
	"strconv"
)

type ANumber interface{
  AValue
  Addable
  Subtractable
  Dividable
  Multipliable
  Comparable
}

type NumberExpression struct {
	span    *Span
	value  ANumber
}

func ReadNumber(scanner *Scanner) (AnExpression, *ParserError) {
	vInt, e := scanner.Peek()
	if e != nil {
		return nil, ExtendParserError(*scanner.Position(), e)
	}
	if vInt.Token != NUMBER {
    return nil, NewParserError(
      *scanner.Position(),
      fmt.Sprintf("tried to read a number starting from not a NUMBER token but %s %s", vInt.Token, vInt.Literal),
    )
	}

  scanner.Scan()
	result := &NumberExpression{
		span: vInt.Span,
	}

  tok, err := scanner.Peek()
	if err != nil {
		return nil, ExtendParserError(*scanner.Position(), err)
	}
	if tok.Token == DOT {
		result.span = result.span.Extend(tok.Span)
		scanner.Scan()
		tok, e = scanner.Peek()
		if e != nil {
			return nil, ExtendParserError(*scanner.Position(), e)
		}
    var val float64
		if tok.Token == NUMBER {
			vFraction := scanner.scanNumber()
      val, e = strconv.ParseFloat(fmt.Sprintf("%s.%s", vInt.Literal, vFraction.Literal), 64)
      if e != nil {
        return nil, ExtendParserError(*scanner.Position(), e)
      }
			result.span = result.span.Extend(vFraction.Span)
		} else {
      val, e = strconv.ParseFloat(vInt.Literal, 64)
      if e != nil {
        return nil, ExtendParserError(*scanner.Position(), e)
      }
    }
    result.value = Float(val)
	} else {
    var val int64
    val, e = strconv.ParseInt(vInt.Literal, 10, 64)
    if e != nil {
      return nil, ExtendParserError(*scanner.Position(), e)
    }
    result.value = Integer(val)
  }
  return result, nil
}

func (n *NumberExpression) Type() Type {
  if n.value == nil {
    return Type(reflect.Invalid)
  }
  return n.value.Type()
}

func (n *NumberExpression) Span() *Span {
	return n.span
}

func (n *NumberExpression) Value() AValue {
  return n.value
}
