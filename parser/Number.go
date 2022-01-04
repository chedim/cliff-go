package parser

import (
	"fmt"
	"strconv"
)

type Number struct {
	span    *Span
	isFloat bool
	vFloat  float64
	vInt    int64
}

func ReadNumber(scanner *Scanner) (*Number, *ParserError) {
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
	result := &Number{
		span: vInt.Span,
	}

  tok, err := scanner.Peek()
	if err != nil {
		return nil, ExtendParserError(*scanner.Position(), err)
	}
	if tok.Token == DOT {
		result.isFloat = true
		result.span = result.span.Extend(tok.Span)
		scanner.Scan()
		tok, e = scanner.Peek()
		if e != nil {
			return nil, ExtendParserError(*scanner.Position(), e)
		}
		if tok.Token == NUMBER {
			vFraction := scanner.scanNumber()
      result.vFloat, e = strconv.ParseFloat(fmt.Sprintf("%s.%s", vInt.Literal, vFraction.Literal), 64)
      if e != nil {
        return nil, ExtendParserError(*scanner.Position(), e)
      }
			result.span = result.span.Extend(vFraction.Span)
		} else {
      result.vFloat, e = strconv.ParseFloat(vInt.Literal, 64)
      if e != nil {
        return nil, ExtendParserError(*scanner.Position(), e)
      }
    }
	} else {
    result.vInt, e = strconv.ParseInt(vInt.Literal, 10, 64)
    if e != nil {
      return nil, ExtendParserError(*scanner.Position(), e)
    }
  }
  return result, nil
}

func (n *Number) Span() *Span {
	return n.span
}

func (n *Number) Value() AValue {
  var ret AValue
  if n.isFloat {
    ret = n.vFloat
  } else {
    ret = n.vInt
  }
  return &ret
}
