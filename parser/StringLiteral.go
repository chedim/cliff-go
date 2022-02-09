package parser

import (
	"fmt"
	"strconv"
)

type StringLiteral struct {
	location Span
	value string
}

var escapeCharacterMap = map[rune]rune{
	'\\': '\\',
  '\'': '\'',
  '"':  '"',
	'r':  '\r',
	'n':  '\n',
	't':  '\t',
}

func ReadString(scanner *Scanner, delim Token) (*StringLiteral, *ParserError) {
	first := scanner.Peek()

	if first.Token != delim {
    return nil, NewParserError(*scanner.Position(), fmt.Sprint("tried to read ", delim, " string at position that does not contain ", delim))
	}

	result := &StringLiteral{location: *scanner.Position()}
	scanner.Scan()
  scanner.PreserveCase(true)
	escaped := false
	for token := scanner.Peek(); token.Token != delim || escaped; token = scanner.Peek() {
		if token.Token == EOF {
			return nil, NewParserError(result.location, "Infinite string")
		} else if token.Token == SLASH && !escaped {
			escaped = true
      scanner.Scan()
			continue
		} else if escaped {
			escapedChar := rune(token.Literal[0])
      escaped = false
			if escapedChar == 'u' {
				if len(token.Literal) != 5 {
					return nil, NewParserError(*scanner.Position(), "Unsupported unicode string")
				} else {
          i, err := strconv.ParseInt(token.Literal[1:], 16, 32)
          if err != nil {
            return nil, ExtendParserError(*scanner.Position(), err)
          }
          token.Literal = string(rune(i))
				}
			} else {
        escvalue, exists := escapeCharacterMap[escapedChar]
        if !exists {
          return nil, NewParserError(*scanner.Position(), "Unknown escape code")
        }
        token.Literal = string(escvalue)
      }
    }

    result.value += token.Literal
		scanner.Scan()
	}

  tok := scanner.Peek()
  if tok.Token != delim {
    return nil, NewParserError(*scanner.Position(), fmt.Sprintf("Wrong string end delimeter: %s %s", tok.Token, tok.Literal))
  }
  scanner.Scan()
  scanner.PreserveCase(false)
	return result, nil
}

func (s *StringLiteral) Span() *Span {
  return &s.location
}

func (s *StringLiteral) Value() AValue {
  var result AValue = String(s.value)
  return result
}
