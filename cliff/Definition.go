package cliff;

type Definition struct {
  AnExpression
  Condition AnExpression
}

func ReadDefinition(scanner *Scanner) (*Definition, *ParserError) {
  result := &Definition{}

  tok, e := scanner.Peek();
  if e != nil {
    return nil, ExtendParserError(*scanner.Position(), e)
  }
  if tok.Token != IS && tok.Token != MINUS {
    return nil, NewParserError(*scanner.Position(), "tried to read definition starting from a token that is not IS or MINUS")
  }
  scanner.Scan()

  tok,e = scanner.Peek()
  if tok.Token != WHEN {
    exp, err := ReadExpression(scanner)
    if err != nil {
      return nil, err
    }
    result.AnExpression = exp
  } else {
    scanner.Scan()
    exp, err := ReadExpression(scanner)
    if err != nil {
      return nil, err
    }
    result.Condition = exp
  }

  return result, nil
}
