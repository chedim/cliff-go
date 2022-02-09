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
	vInt := scanner.Peek()
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

  tok := scanner.Peek()
	if tok.Token == DOT {
		result.span = result.span.Extend(tok.Span)
		scanner.Scan()
		tok = scanner.Peek()
    var val float64
		if tok.Token == NUMBER {
			vFraction := scanner.scanNumber()
      parsed, e := strconv.ParseFloat(fmt.Sprintf("%s.%s", vInt.Literal, vFraction.Literal), 64)
      if e != nil {
        return nil, ExtendParserError(*scanner.Position(), e)
      }
      val = parsed
			result.span = result.span.Extend(vFraction.Span)
		} else {
      parsed, e := strconv.ParseFloat(vInt.Literal, 64)
      if e != nil {
        return nil, ExtendParserError(*scanner.Position(), e)
      }
      val = parsed
    }
    result.value = Float(val)
	} else {
    val, e := strconv.ParseInt(vInt.Literal, 10, 64)
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
