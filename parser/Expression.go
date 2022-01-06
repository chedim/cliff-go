package parser

import "fmt"


type AnExpression interface {
	Span() *Span
	Value() AValue
}

type NumberLiteral struct {
	value int
}

type ExpressionReader func(scanner *Scanner, stack *Stack) *ParserError
type ExpressionJoiner func(l AnExpression, r AnExpression) AnExpression

func binaryOperator(op BinaryOperator) ExpressionReader {
  return func(scanner *Scanner, stack *Stack) *ParserError {
    scanner.Scan()
    if stack.Len() == 0 {
      return NewParserError(*scanner.Position(), "missing left operand")
    }

    left := stack.Pop().(AnExpression)
    right, err := ReadExpression(scanner)
    if err != nil {
      return err
    }

    stack.Push(NewBinaryExpression(left, right, op))
    return nil
  }
}

var valueHandlers map[Token]ExpressionReader

func init() {
  valueHandlers = map[Token]ExpressionReader{
    QUOTE:  readSingleQuotedString,
    DQUOTE: readDoubleQuotedString,
    WORD: readReference,
    NUMBER: readNumber,
    WS: skipToken,
    EOL: skipToken,
    TRUE: trueExpression,
    FALSE: falseExpression,
    PLUS: binaryOperator(add),
    MINUS: binaryOperator(sub),
    SLASH: binaryOperator(div),
    ASTERISK: binaryOperator(mul),
  }
}

func add(l AValue, r AValue) AValue {
  return l.(Addable).Add(r)
}

func sub(l AValue, r AValue) AValue {
  return l.(Subtractable).Sub(r)
}

func div(l AValue, r AValue) AValue {
  return l.(Dividable).Div(r)
}

func mul(l AValue, r AValue) AValue {
  return l.(Multipliable).Mul(r)
}

func trueExpression(scanner *Scanner, stack *Stack) *ParserError {
  stack.Push(NewConstExpression(scanner.Position(), Bool(true)))
  scanner.Scan()
  return nil
}

func falseExpression(scanner *Scanner, stack *Stack) *ParserError {
  stack.Push(NewConstExpression(scanner.Position(), Bool(false)))
  scanner.Scan()
  return nil
}

func skipToken(scanner *Scanner, stack *Stack) *ParserError {
  tok, e := scanner.Peek()
  if e != nil {
    return ExtendParserError(*scanner.Position(), e)
  }
  if (tok.Token == EOL) {
    fmt.Println("skipping EOL")
    scanner.Scan()
  } else if tok.Token == WS {
    fmt.Printf("skipping WS '%s' @%d\n", tok.Literal, scanner.offset)
    scanner.scanWhitespace()
  } else {
    return NewParserError(*scanner.Position(), fmt.Sprintf("Unable to skip token %s %s", tok.Token, tok.Literal))
  }

  return nil
}

func readNumber(scanner *Scanner, s *Stack) *ParserError {
  n, e := ReadNumber(scanner)
  if e != nil {
    return e
  }
  s.Push(n)
  return nil
}

func readReference(scanner *Scanner, s *Stack) *ParserError {
  r, e := ReadReference(scanner)
  if e != nil {
    return e
  }
  s.Push(r)
  return nil
}

func readSingleQuotedString(scanner *Scanner, s *Stack) *ParserError {
	r, e := ReadString(scanner, QUOTE)
  if e != nil {
    return e
  }
  s.Push(r)
	return nil
}

func readDoubleQuotedString(scanner *Scanner, s *Stack) *ParserError {
	r, e := ReadString(scanner, DQUOTE)
  if e != nil {
    return e
  }
  s.Push(r)
	return nil
}

func ReadExpression(scanner *Scanner) (AnExpression, *ParserError) {
  stk := NewStack()
  start := *scanner.Position()

  for tok, err := scanner.Peek(); tok.Token != EOF; tok, err = scanner.Peek() {
    if err != nil {
      return nil, ExtendParserError(*scanner.Position(), err)
    }
    handler, exists := valueHandlers[tok.Token]
    if !exists {
      break
    }
    e := handler(scanner, stk)
    if e != nil {
      return nil, e
    }
  }

  if stk.Len() != 1 {
    return nil, NewParserError(start, "disjoint expression")
  }

  return stk.Pop().(AnExpression), nil
}
