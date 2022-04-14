package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type Statement struct {
  Parent *Statement
  Location Span
  Abstract bool
	Target *Reference
  Plural bool
  Labels []*Datapoint
  Definitions []*Definition
  SubStatements []*Statement
}

func ParseStatement(text string) (*Statement, *ParserError) {
  return ReadStatement(NewCliffScanner(strings.NewReader(text)))
}

func ReadStatement(scanner *Scanner) (statement *Statement, err *ParserError){
  statement = new(Statement)
  statement.Location = *scanner.Position()

  statement.Target, err = ReadReference(scanner)
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
    return statement.readArrayDefinitions(scanner)
  } else if isExpressionToken(token.Token) {
    def, err := ReadDefinition(scanner)
    if err != nil {
      return err
    }
    statement.Definitions = append(statement.Definitions, def)
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
    return statement.readArrayDefinitions(scanner)
  } else if first.Token == WHEN {
    def := &Definition{}
    cnd, err := ReadExpression(scanner)
    if err != nil {
      return err
    }
    def.condition = cnd
    statement.Definitions = append(statement.Definitions, def)
  } else if isExpressionToken(first.Token) {
    def, err := ReadDefinition(scanner)
    if err != nil {
      return err
    }
    statement.Definitions = append(statement.Definitions, def)
  } else {
    return NewParserError(*scanner.Position(), fmt.Sprintf("Unexpected token %s '%s'", first.Token, first.Literal))
  }
  return nil
}

func (statement *Statement) readArrayDefinitions(scanner *Scanner) *ParserError {
  ft := scanner.Peek()

  logger, _ := zap.NewDevelopment()
  l := logger.Sugar()
  l.Debug(scanner.Position().String(), " | Reading alternative definitions")

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

  parentWhitespace := scanner.GetMinOffset()
  currentWhitespace := -1

  for ft = scanner.scanWhitespace(); ft.Length >= scanner.GetMinOffset(); ft = scanner.scanWhitespace() {
    if (currentWhitespace < 0) {
      currentWhitespace = ft.Length
      scanner.SetMinOffset(currentWhitespace)
    } else if (ft.Length != currentWhitespace) {
      return NewParserError(*scanner.Position(),
        fmt.Sprintf("Alternative definition offset mismatch, was expecting whitespace length %d characters, but got %d", currentWhitespace, ft.Length))
    }

    if def, err := ReadDefinition(scanner); err != nil {
      return err
    } else {
      l.Debug(scanner.Position().String(), "| Read Definition ->", def.String())
      statement.Definitions = append(statement.Definitions, def)
    }

    if eol := scanner.Peek(); eol.Token == EOF {
      break
    } else if eol.Token != EOL {
      return NewParserError(*scanner.Position(), "Expected EOL but got %s", eol.Token)
    }

    scanner.Scan()
  }

  scanner.SetMinOffset(parentWhitespace)

  return nil
}

func (statement *Statement) readCompoundStatements(scanner *Scanner) *ParserError {
  scanner.scanWhitespace()
  tok := scanner.Peek()
  if tok.Token != EOL {
    return NewParserError(*scanner.Position(), fmt.Sprintf("unexpected token in compound statement start: %s %s", tok.Token, tok.Literal))
  }

  offset := statement.Location.StartColumn
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
  ss.Parent = s
  if err != nil {
    return WrapParserError(err, "Failed to read sub-statement")
  }

  s.SubStatements = append(s.SubStatements, ss)
  return nil
}

func (statement *Statement) readPositionExpression(scanner *Scanner) *ParserError {
  first := scanner.Peek()
  if first.Token != AT {
    return NewParserError(*scanner.Position(), "trying to read position expression that does not start with AT")
  }
  scanner.Scan()
  scanner.scanWhitespace()
  ps := &Statement{Location: *scanner.Position(), Parent: statement}
  ps.Target = NewReference("position")

  statement.SubStatements = append(statement.SubStatements, ps)
  return nil
}

func (s *Statement) Dependencies() []*Datapoint {
  return nil
}

func (s *Statement) String() string {
  if result, err := json.Marshal(s); err != nil {
    return fmt.Sprintf("ERROR: %s", err)
  } else {
    return string(result)
  }
}
