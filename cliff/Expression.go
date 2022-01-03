package cliff

import "fmt"

type AValue interface {
}

type AnExpression interface {
	Span() *Span
	Value() *AValue
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

var valueHandlers = map[Token]func(scanner *Scanner) (AnExpression, *ParserError){
	QUOTE:  readSingleQuotedString,
	DQUOTE: readDoubleQuotedString,
	//  WORD: ReadReference,
	//  NUMBER: ReadNumber,
}

func readSingleQuotedString(scanner *Scanner) (AnExpression, *ParserError) {
	r, e := ReadString(scanner, QUOTE)
	t := AnExpression(r)
	return t, e
}

func readDoubleQuotedString(scanner *Scanner) (AnExpression, *ParserError) {
	r, e := ReadString(scanner, DQUOTE)
	t := AnExpression(r)
	return t, e
}

func ReadExpression(scanner *Scanner) (AnExpression, *ParserError) {
	tok, err := scanner.Peek()
	if err != nil {
		return nil, ExtendParserError(*scanner.Position(), err)
	}

	if handler, exists := valueHandlers[tok.Token]; exists {
		return handler(scanner)
	}

	return nil, NewParserError(*scanner.Position(), fmt.Sprint("unexpected expression literal: ", tok.Token, " ", tok.Literal))
}
