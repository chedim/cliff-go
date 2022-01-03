package cliff

import (
	"strings"
	"testing"
)


func TestIsWhitespace(t *testing.T) {
  if !isWhitespace(' ') {
    t.Error("Space is not whitespase")
  }

  if !isWhitespace('\t') {
    t.Error("Tab is not whitespace")
  }

  if isWhitespace('\n') {
    t.Error("Newline is whitespace")
  }
}

func TestIsLetter(t *testing.T) {
  for ch := rune(0); ch <= rune(255); ch++ {
    expected := (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
    actual := isLetter(ch)

    if expected && !actual {
      t.Errorf("Rune %c is not a letter", ch)
    } else if !expected && actual {
      t.Errorf("Rune %c is a letter", ch)
    }
  }
}

func TestIsNumber(t *testing.T) {
  for ch := rune(0); ch <= rune(255); ch++ {
    expected := (ch >= '0' && ch <= '9')
    actual := isNumber(ch)

    if expected && !actual {
      t.Errorf("Rune %c is not a number", ch)
    } else if !expected && actual {
      t.Errorf("Rune %c is a number", ch)
    }
  }
}

func TestHelloWorld(t *testing.T) {
  source := strings.NewReader("heLlo, world\n2nd line")

  scanner := NewCliffScanner(source)

  token, e := scanner.Peek()
  if e != nil {
    t.Errorf("Failed to peek %e", e)
  }
  if token.Token != WORD {
    t.Error("token is not a word")
  }

  if token.Literal != "hello" {
    t.Errorf("token is not hello: %s", token.Literal)
  }

  if token.Span.StartLine != 0 || token.Span.StartColumn != 0 {
    t.Errorf("wrong token start: %d:%d", token.Span.StartLine, token.Span.StartColumn)
  }

  if token.Length != 5 {
    t.Errorf("wrong token length: %d", token.Length)
  }

  if token.EndColumn != 5 {
    t.Errorf("wrong token end column: %d", token.EndColumn)
  }

  scanned := scanner.Scan()
  if (scanned.Literal != "hello") {
    t.Errorf("scanned literal is not hello: %s", scanned.Literal)
  }

  token, e = scanner.Peek()
  if e != nil {
    t.Errorf("failed to peek: %s", e)
  }

  if token.Token != COMMA {
    t.Errorf("token is not COMMA: %s", token.Token)
  }

  if token.Length != 1 {
    t.Errorf("token length is not 1: %d", token.Length)
  }

  if token.StartLine != 0 && token.StartColumn != 5 {
    t.Errorf("token start is not at 0,5: %d,%d", token.StartLine, token.StartColumn)
  }

  if token.EndLine != 0 && token.EndColumn != 6 {
    t.Errorf("token end is not at 0, 6: %d,%d", token.EndLine, token.EndColumn)
  }

  scanned = scanner.Scan()
  if scanned.Literal != "," {
    t.Errorf("scanned literal is not comma: %s", scanned.Literal)
  }

  token, e = scanner.Peek()
  if token.Literal != " " {
    t.Errorf("expected to peek a space but got: '%s'", token.Literal)
  }
  if token.Token != WS {
    t.Errorf("token is not whitespace")
  }
  if token.StartLine != 0 && token.StartColumn != 7 {
    t.Errorf("token start location is not 0,7: %d,%d", token.StartLine, token.StartColumn)
  }
  if token.EndLine != 0 && token.EndColumn != 8 {
    t.Errorf("token end location is not 0,8: %d, %d", token.EndLine, token.EndColumn)
  }

  scanned = scanner.Scan()
  if scanned.Literal != " " {
    t.Errorf("unexpected scan literal '%s'", scanned.Literal)
  }

  token, e = scanner.Peek()
  if e != nil {
    t.Errorf("failed to peek, %s", e)
  }

  if token.Token != WORD {
    t.Errorf("token is not WORD but %s", token.Token)
  }
  if token.Literal != "world" {
    t.Errorf("token literal is not 'world' but '%s'", token.Literal)
  }
  if token.EndColumn != 12 {
    t.Errorf("last line token end column is not 12: %d", token.EndColumn)
  }

  scanned = scanner.Scan()
  if scanned.Literal != token.Literal {
    t.Errorf("scanned literal is not 'world': %s", scanned.Literal)
  }

  token, e = scanner.Peek()
  if e != nil {
    t.Errorf("failed to peek: %s", e)
  }

  if token.Token != EOL {
    t.Errorf("token is not EOL: %s", token.Token)
  }
  if token.StartLine != 0 && token.StartColumn != 13 {
    t.Errorf("token start is invalid: %d,%d", token.StartLine, token.StartColumn)
  }
  if token.EndLine != 0 && token.EndColumn != 14 {
    t.Errorf("token end is invalid: %d, %d", token.EndLine, token.EndColumn)
  }

  scanned = scanner.Scan()
  if scanned.Literal != "\n" {
    t.Errorf("scanned literal is invalid: %s", scanned.Literal)
  }

  token, e = scanner.Peek()
  if e != nil {
    t.Errorf("failed to peek: %s", e)
  }
  if token.Token != NUMBER {
    t.Errorf("invalid token: %s", token.Token)
  }
  if token.Literal != "2" {
    t.Errorf("invalid literal: %s", token.Literal)
  }
  if token.StartLine != 1 && token.StartColumn != 0 {
    t.Errorf("invalid start position: %d,%d", token.StartLine, token.StartColumn)
  }
  if token.EndLine != 1 && token.EndColumn != 0 {
    t.Errorf("invalid end position: %d,%d", token.EndLine, token.EndColumn)
  }

  scanned = scanner.Scan()
  if scanned.Literal != "2" {
    t.Errorf("invalid scanned literal: %s", scanned.Literal)
  }

  token, e = scanner.Peek()
  if e != nil {
    t.Errorf("failed to peek: %s", e)
  }
  if token.Token != ND {
    t.Errorf("invalid token: %s", token.Token)
  }

  scanned = scanner.Scan()
  if scanned.Literal != "nd" {
    t.Errorf("invalid literal: %s", scanned.Literal)
  }
  if !scanned.Keyword {
    t.Errorf("ND should be a keyword")
  }

  scanned = scanner.Scan()
  if scanned.Token != WS {
    t.Errorf("invalid token: %s", scanned.Token)
  }

  scanned = scanner.Scan()
  if scanned.Token != WORD {
    t.Errorf("invalid token: %s", scanned.Token)
  }
  if scanned.Literal != "line" {
    t.Errorf("invalid literal: %s", scanned.Literal)
  }

  scanned = scanner.Scan()
  if scanned.Token != EOF {
    t.Errorf("token is not EOF: '%s' (%s)", scanned.Literal, scanned.Token)
  }
}
