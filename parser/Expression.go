package parser

import (
	"fmt"

	"go.uber.org/zap"
)


type AnExpression interface {
	Span() *Span
	Value() AValue
  String() string
}

type NumberLiteral struct {
	value int
}

type ExpressionReader func(scanner *Scanner, stack *Stack) (bool, *ParserError)

func binaryOperator(op BinaryOperator) ExpressionReader {
  return func(scanner *Scanner, stack *Stack) (bool, *ParserError) {
    scanner.Scan()
    if stack.Len() == 0 {
      return false, NewParserError(*scanner.Position(), "missing left operand")
    }

    left := stack.Pop().(AnExpression)
    right, err := ReadExpression(scanner)
    if err != nil {
      return false, err
    }

    stack.Push(NewBinaryExpression(left, right, op))
    return true, nil
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
    THE: readSliceExpression,
    IS: binaryOperator(equals),
  }
}

func equals(l AValue, r AValue) AValue {
  return l.Equals(r)
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

func trueExpression(scanner *Scanner, stack *Stack) (bool, *ParserError) {
  if (stack.Len() != 0) {
    return false, nil
  }
  stack.Push(NewConstExpression(scanner.Position(), Bool(true)))
  scanner.Scan()
  return true, nil
}

func falseExpression(scanner *Scanner, stack *Stack) (bool, *ParserError) {
  if (stack.Len() != 0) {
    return false, nil
  }
  stack.Push(NewConstExpression(scanner.Position(), Bool(false)))
  scanner.Scan()
  return true, nil
}

func skipToken(scanner *Scanner, stack *Stack) (bool, *ParserError) {
  tok := scanner.Peek()
  if (tok.Token == EOL) {
    if (stack.Len() > 0) {
      return false, nil
    }
    scanner.Scan()
  } else if tok.Token == WS {
    scanner.scanWhitespace()
  } else {
    return false, NewParserError(*scanner.Position(), fmt.Sprintf("Unable to skip token %s %s", tok.Token, tok.Literal))
  }

  return true, nil
}

func readNumber(scanner *Scanner, s *Stack) (bool, *ParserError) {
  if (s.Len() != 0) {
    return false, nil
  }
  n, e := ReadNumber(scanner)
  if e != nil {
    return false, e
  }
  s.Push(n)
  return true, nil
}

func readReference(scanner *Scanner, s *Stack) (bool, *ParserError) {
  if (s.Len() != 0) {
    return false, nil
  }
  r, e := ReadReference(scanner)
  if e != nil {
    return false, e
  }
  s.Push(r)
  return true, nil
}

func readSingleQuotedString(scanner *Scanner, s *Stack) (bool, *ParserError) {
  if (s.Len() != 0) {
    // doesn't support left operands
    return false, nil
  }
	r, e := ReadString(scanner, QUOTE)
  if e != nil {
    return false, e
  }
  s.Push(r)
	return true, nil
}

func readDoubleQuotedString(scanner *Scanner, s *Stack) (bool, *ParserError) {
  if (s.Len() != 0) {
    // end of expression
    return false, nil
  }
	r, e := ReadString(scanner, DQUOTE)
  if e != nil {
    return false, e
  }
  s.Push(r)
	return true, nil
}

func readSliceExpression(scanner *Scanner, s *Stack) (bool, *ParserError) {
  if (s.Len() != 0) {
    return false, nil
  }
  r, e := ReadSliceExpression(scanner)
  if e != nil {
    return false, e
  }
  s.Push(r)
  return true, nil
}

func ReadExpression(scanner *Scanner) (AnExpression, *ParserError) {
  stk := NewStack()
  start := *scanner.Position()
  logger, _ := zap.NewDevelopment()
  l := logger.Sugar()
  l.Debug(scanner.Position().String(), " | Reading expression")

  for tok := scanner.Peek(); tok.Token != EOF; tok = scanner.Peek() {
    l.Debugf("expr token: %+v", tok)
    handler, exists := valueHandlers[tok.Token]
    if !exists {
      l.Debug("^^ no handler")
      break
    }
    l.Debug("^^ delegating to handler")
    read, e := handler(scanner, stk)
    if e != nil {
      return nil, e
    }

    if !read {
      l.Debug("Handler signalled stop of expression")
      break;
    }
  }

  if stk.Len() != 1 {
    return nil, NewParserError(start, "disjoint expression, stack length: %d; stack values:\n%s", stk.Len(), stk.String())
  }

  return stk.Pop().(AnExpression), nil
}
