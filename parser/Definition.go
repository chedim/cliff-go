package parser

import (
	"fmt"
)

type Definition struct {
	value     AnExpression
	condition AnExpression
  related   AnExpression
}

func ReadDefinition(scanner *Scanner) (*Definition, *ParserError) {
	result := &Definition{}

  if tok := scanner.Peek(); !isExpressionToken(tok.Token) {
    if (tok.Token == WHEN) {
      LogDebug("Detected early when")
      result.readConditionSection(scanner)
    } else if (tok.Token == AFTER) {
      LogDebug("Detected early after")
      result.readConditionSection(scanner)
    } else {
      return nil, NewParserError(*scanner.Position(), fmt.Sprintf("tried to read definition starting from a token that is not an expression token: %s '%s'", tok.Token, tok.Literal))
    }

    scanner.scanWhitespace()

    if tok = scanner.Peek(); tok.Token == EOL {
      return result, nil
    } else if tok.Token == COLON {
      scanner.Scan()
    } else {
      return nil, scanner.Error("Expected colon but got: %s", tok.Literal)
    }
	}

  if exp, err := ReadExpression(scanner); err != nil {
    return nil, err
  } else {
    result.value = exp
  }

  for tok := scanner.Peek(); tok.Token == WHEN || tok.Token == AFTER; tok = scanner.Peek() {
    if tok.Token == WHEN {
      if (result.condition != nil) {
        return nil, scanner.Error("doublewhen: %s -- %+v ... %+v", tok.Token, result.condition.Span(), result.value.Value())
      }
      result.readConditionSection(scanner)
    } else if tok.Token == AFTER {
      if (result.related != nil) {
        return nil, scanner.Error("doubleafter")
      }
      result.readRelatedSection(scanner)
    }
  }

  return result, nil
}

func (d *Definition) readConditionSection(scanner *Scanner) *ParserError {
  LogDebug("Reading condition section")
  if tok := scanner.Peek(); tok.Token != WHEN {
    return scanner.Error("Tried to read condition starting with a token that is not WHEN: %s", tok.Token)
  }
  scanner.Scan()
  scanner.scanWhitespace()
  exp, err := ReadExpression(scanner)
  if err != nil {
    return err
  }
  d.condition = exp
  return nil
}

func (d *Definition) readRelatedSection(scanner *Scanner) *ParserError {
  if tok := scanner.Peek(); tok.Token != AFTER {
    return scanner.Error("Tried to read related definition section starting with a token that is not AFTER: %s", tok.Token)
  }
  scanner.Scan()
  scanner.scanWhitespace()
  if exp, err := ReadExpression(scanner); err == nil {
    d.related = exp
  } else {
    return err
  }
  return nil
}

func (d *Definition) String() string {
  s := d.value.String();
  if (d.condition != nil) {
    s += fmt.Sprintf(" when %s", d.condition.String())
  }
  if (d.related != nil) {
    s += fmt.Sprintf(" after %s ", d.related.String())
  }
  return s
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
