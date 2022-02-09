package parser

import "fmt"

type Statement struct {
  parent *Statement
  location Span
  abstract bool
	target *Reference
  plural bool
  labels []*Datapoint
  definitions []*Definition
  subStatements []*Statement
}

func ReadStatement(scanner *Scanner) (statement *Statement, err *ParserError){
  statement = new(Statement)
  statement.location = *scanner.Position()

  statement.target, err = ReadReference(scanner)
  if err != nil {
    return nil, err
  }

  operator := scanner.Peek();
  if operator.Token == IS {
    return statement, statement.fillSingularStatement(scanner)
  } else if operator.Token == ARE {
    return statement, statement.fillPluralStatement(scanner)
  } else if operator.Token == EOL {
    scanner.Scan()
    return statement, statement.readCompoundStatements(scanner)
  }

  return nil, NewParserError(*scanner.Position(), fmt.Sprintf("Expected definition but got something else: %s %s", operator.Token, operator.Literal))
}

func (statement *Statement) fillPluralStatement(scanner *Scanner) *ParserError {
  token := scanner.Peek()
  if token.Token != ARE {
    return scanner.Error("tried to fill plural statement starting with wrong token: %s %s", token.Token, token.Literal)
  }
  scanner.Scan()
  scanner.scanWhitespace()

  token = scanner.Peek()
  if token.Token == AT {
    return statement.readPositionExpression(scanner)
  } else if token.Token == COLON {
    return statement.readAlternativeDefinitions(scanner)
  } else if isExpressionToken(token.Token) {
    def, err := ReadDefinition(scanner)
    if err != nil {
      return err
    }
    statement.definitions = append(statement.definitions, def)
  } else {
    return scanner.Error("Unexpected token %s '%s'", token.Token, token.Literal)
  }
  return nil
}

func (statement *Statement) fillSingularStatement(scanner *Scanner) *ParserError {
  start := scanner.Peek()
  if start.Token != IS {
    return NewParserError(*scanner.Position(), fmt.Sprintf("tried to fill singular statement starting with wrong token: %s %s", start.Token, start.Literal))
  }
  scanner.Scan()
  scanner.scanWhitespace()

  first := scanner.Peek()
  if first.Token == AT {
    return statement.readPositionExpression(scanner)
  } else if first.Token == COLON {
    return statement.readAlternativeDefinitions(scanner)
  } else if first.Token == WHEN {
    def := &Definition{}
    cnd, err := ReadExpression(scanner)
    if err != nil {
      return err
    }
    def.condition = cnd
    statement.definitions = append(statement.definitions, def)
  } else if isExpressionToken(first.Token) {
    def, err := ReadDefinition(scanner)
    if err != nil {
      return err
    }
    statement.definitions = append(statement.definitions, def)
  } else {
    return NewParserError(*scanner.Position(), fmt.Sprintf("Unexpected token %s '%s'", first.Token, first.Literal))
  }
  return nil
}

func (statement *Statement) readAlternativeDefinitions(scanner *Scanner) *ParserError {
  ft := scanner.Peek()
  if ft.Token != COLON {
    return NewParserError(*scanner.Position(), "tried to read alternative definitions starting with token that is not COLON")
  }
  scanner.Scan()

  scanner.scanWhitespace()
  ft = scanner.Peek()
  if ft.Token != EOL {
    return NewParserError(*scanner.Position(), fmt.Sprintf("expected new line but got: %s", ft.Token))
  }
  scanner.Scan()

  for ft = scanner.Peek(); ft == nil || ft.Token == MINUS; ft = scanner.Peek() {
    if ft == nil {
      return NewParserError(*scanner.Position(), "ft is null")
    }

    def, err := ReadDefinition(scanner)
    if err != nil {
      return err
    }
    statement.definitions = append(statement.definitions, def)
  }
  return nil
}

func (statement *Statement) readCompoundStatements(scanner *Scanner) *ParserError {
  scanner.scanWhitespace()
  tok := scanner.Peek()
  if tok.Token != EOL {
    return NewParserError(*scanner.Position(), fmt.Sprintf("unexpected token in compound statement start: %s %s", tok.Token, tok.Literal))
  }

  offset := statement.location.StartColumn
  for mine, e := scanner.scanOffset(offset); mine; mine, e = scanner.scanOffset(offset) {
    if e != nil {
      return ExtendParserError(*scanner.Position(), e)
    }

    nextToken := scanner.Peek()
    if nextToken.Token == WS {
      err := statement.readSubStatement(scanner)
      if err != nil {
        return err
      }
    }
    //todo
  }
  return nil
}

func (s *Statement) readSubStatement(scanner *Scanner) *ParserError {
  ot := scanner.Peek()
  if ot.Token != WS {
    return NewParserError(*scanner.Position(), "substatements must start with whitespace")
  }
  scanner.Scan()

  ss, err := ReadStatement(scanner)
  ss.parent = s
  if err != nil {
    return WrapParserError(err, "Failed to read sub-statement")
  }

  s.subStatements = append(s.subStatements, ss)
  return nil
}

func (statement *Statement) readPositionExpression(scanner *Scanner) *ParserError {
  first := scanner.Peek()
  if first.Token != AT {
    return NewParserError(*scanner.Position(), "trying to read position expression that does not start with AT")
  }
  scanner.Scan()
  scanner.scanWhitespace()
  ps := &Statement{location: *scanner.Position(), parent: statement}
  ps.target = NewReference("position")

  statement.subStatements = append(statement.subStatements, ps)
  return nil
}

func (s *Statement) Dependencies() []*Datapoint {
  return nil
}

func (s *Statement) Target() *Reference {
  return s.target
}

func (s *Statement) Definitions() []*Definition {
  return s.definitions
}
