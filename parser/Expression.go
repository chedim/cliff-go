package parser

import "fmt"

type AValue interface {
}

type AnExpression interface {
	Span() *Span
	Value() AValue
}

type NumberLiteral struct {
	value int
}

type ReferenceExpression struct {
	value []*Tokenized
}

type ConditionalExpression struct {
	value     AnExpression
	condition AnExpression
}

type BinaryExpression struct {
	Operator
	Left  AnExpression
	Right AnExpression
}

var valueHandlers = map[Token]func(scanner *Scanner, stack *Stack) *ParserError {
	QUOTE:  readSingleQuotedString,
	DQUOTE: readDoubleQuotedString,
	WORD: readReference,
	NUMBER: readNumber,
  WS: skipToken,
  EOL: skipToken,
  TRUE: trueExpression,
  FALSE: falseExpression,
}

func trueExpression(scanner *Scanner, stack *Stack) *ParserError {
  stack.Push(NewConstExpression(scanner.Position(), true))
  scanner.Scan()
  return nil
}

func falseExpression(scanner *Scanner, stack *Stack) *ParserError {
  stack.Push(NewConstExpression(scanner.Position(), false))
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
