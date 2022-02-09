package parser

import (
	"fmt"
)

type Definition struct {
	value     AnExpression
	condition AnExpression
}


func ReadDefinition(scanner *Scanner) (*Definition, *ParserError) {
	result := &Definition{}

  tok := scanner.Peek()
	if !isExpressionToken(tok.Token) {
		return nil, NewParserError(*scanner.Position(), fmt.Sprintf("tried to read definition starting from a token that is not an expression token: %s '%s'", tok.Token, tok.Literal))
	}

  exp, err := ReadExpression(scanner)
  if err != nil {
    return nil, err
  }
  result.value = exp

	tok = scanner.Peek()

	if tok.Token == WHEN {
		scanner.Scan()
		scanner.scanWhitespace()
		exp, err := ReadExpression(scanner)
		if err != nil {
			return nil, err
		}
		result.condition = exp
	}

	return result, nil
}

func (d *Definition) Span() *Span {
	return d.value.Span()
}

func (d *Definition) Value() AValue {
	return d.value.Value()
}

func (d *Definition) Definition() AnExpression {
	return d.value
}

func (d *Definition) Condition() AnExpression {
	return d.condition
}
