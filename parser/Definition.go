package parser

type Definition struct {
	value     AnExpression
	condition AnExpression
}

func ReadDefinition(scanner *Scanner) (*Definition, *ParserError) {
	result := &Definition{}

	tok, e := scanner.Peek()
	if e != nil {
		return nil, ExtendParserError(*scanner.Position(), e)
	}
	if tok.Token != IS && tok.Token != MINUS {
		return nil, NewParserError(*scanner.Position(), "tried to read definition starting from a token that is not IS or MINUS")
	}
	scanner.Scan()

	scanner.scanWhitespace()
	tok, e = scanner.Peek()
	if e != nil {
		return nil, ExtendParserError(*scanner.Position(), e)
	}
	if tok.Token != WHEN {
		exp, err := ReadExpression(scanner)
		if err != nil {
			return nil, err
		}
		result.value = exp
	}

	tok, e = scanner.Peek()
	if e != nil {
		return nil, ExtendParserError(*scanner.Position(), e)
	}

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
