package parser

type Reference struct {
	span     *Span
	names    []string
	abstract bool
  plural   bool
}

func NewReference(names ...string) *Reference {
	return &Reference{span: nil, names: names}
}

func ReadReference(scanner *Scanner) (*Reference, *ParserError) {
	tok := scanner.Peek()
	result := &Reference{}
	if tok.Token == A || tok.Token == AN {
		scanner.Scan()
		result.abstract = true
	} else if tok.Token != WORD {
		return nil, NewParserError(*scanner.Position(), "tried to read reference from non-WORD token")
	}

	toks := scanner.scanWords()
	if len(toks) == 0 {
		return nil, NewParserError(*scanner.Position(), "zero words scanned for a reference")
	}

	lt := toks[len(toks)-1]
	span := toks[0].Span.Extend(lt.Span)
	result.span = span
	result.names = NormalizedTextArray(toks)
  result.plural = anyPlural(toks)

	return result, nil
}

func anyPlural(toks []*Tokenized) bool {
  for _, tok := range toks {
    if tok.IsPlural {
      return true
    }
  }
  return false
}

func (r *Reference) Span() *Span {
	return r.span
}

func (r *Reference) Target() *Datapoint {
  return DatapointByName(r.names)
}

func (r *Reference) Value() AValue {
  return r.Target().Value()
}

func (r *Reference) Type() Type {
  return r.Target().Type()
}
