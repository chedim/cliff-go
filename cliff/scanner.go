package cliff

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var eof = rune(0)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isNumber(ch rune) bool {
  return ch >= '0' && ch <= '9'
}

type Scanner struct {
	r *bufio.Reader
  offset int32
  line int32
  column int32
  t *Tokenized
}

type Tokenized struct {
  *Span
  Token
  Literal string
  Keyword bool
}

func Text(t []*Tokenized) string {
  r := ""
  for i := 0; i < len(t); i++ {
    r += t[i].Literal
    if (i < len(t)) {
      r += " "
    }
  }
  return r
}

func NewCliffScanner(r io.Reader) *Scanner {
  return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
  ch, _, err := s.r.ReadRune()
  if err != nil {
    return eof
  }
  s.offset++
  if ch == '\n' {
    s.line++
    s.column = 0
  } else {
    s.column++
  }
  return ch
}

func (s *Scanner) peek() rune {
  ch := s.read()
  s.r.UnreadRune()
  return ch
}

func (s *Scanner) Peek() (*Tokenized, error) {
  if s.t != nil {
    return nil, errors.New("Already peeked")
  }
  s.t = s.Scan()
  return s.t, nil
}

func (s *Scanner) Scan() (result *Tokenized) {
  if s.t != nil {
    result = s.t
    s.t = nil
    return
  }
  ch := s.peek()

  if isWhitespace(ch) {
    return s.scanWhitespace()
  } else if isLetter(ch) {
    return s.scanWord()
  } else if isNumber(ch) {
    return s.scanNumber()
  } else {
    result = &Tokenized{
      Span: s.Position(),
      Literal: string(ch),
    }

    result.Length = 1
    result.EndColumn++

    switch ch {
      case '*': result.Token = ASTERISK
      case ',': result.Token = COMMA
      case ':': result.Token = COLON
      case ';': result.Token = SEMICOLON
      case '+': result.Token = PLUS
      case '-': result.Token = MINUS
      case '.': result.Token = DOT
    }
  }
  return
}

func (s *Scanner) scanWhitespace() (result *Tokenized) {
  result = new(Tokenized)
  result.Token = WS
  result.Span = s.Position()

  for ch := s.peek(); isWhitespace(ch); ch = s.peek() {
    result.Length++
    if (ch == '\n') {
      result.EndLine++
      result.EndColumn = 0
    } else {
      result.EndColumn++
    }

    result.Literal += strings.ToLower(ch)
    s.read()
  }
  return
}

func (s *Scanner) Position() *Span {
  return &Span{
    Start: s.offset,
    Length: 0,
    StartLine: s.line,
    StartColumn: s.column,
    EndLine: s.line,
    EndColumn: s.column,
  }
}

func (s *Scanner) scanKeywords() (toks []*Tokenized) {
  toks = make([]*Tokenized, 0)
  for tok, e := s.Peek(); e == nil && tok.Keyword; tok, e = s.Peek() {
    toks = append(toks, s.Scan())
  }
  return
}


func (s *Scanner) scanWords() (toks []*Tokenized) {
  toks = make([]*Tokenized, 0)
  for tok, e := s.Peek(); e == nil && tok.Token == WORD; tok, e = s.Peek() {
    toks = append(toks, s.Scan())
  }
  return
}

func (s *Scanner) scanWord() (result *Tokenized) {
  result = &Tokenized{Token: WORD, Span: s.Position()}

  for ch := s.peek(); isLetter(ch) || isNumber(ch); ch = s.peek() {
    if (result.Literal == "" && isNumber(ch)) {
      result.Token = ILLEGAL
      return result
    }

    result.Span.Length++
    result.Span.EndColumn++ // can't have line breaks in words
    result.Literal += strings.ToLower(ch)

    s.read()
  }
  return detectKeyword(result)
}

var Keywords = map[string]Token{
  "is": IS,
  "are": ARE,
  "when": WHEN,
  "then": THEN,
  "and": AND,
  "or": OR,
  "of": OF,
}

func detectKeyword(in *Tokenized) *Tokenized {
  lclit := strings.ToLower(in.Literal)
  if tok, k := Keywords[lclit]; k {
    in.Token = tok
    in.Keyword = true
  }
  return in
}

func (s *Scanner) scanNumber() (result *Tokenized) {
  result = &Tokenized{
    Token: NUMBER,
    Span: s.Position(),
  }

  for ch := s.peek(); isNumber(ch); ch = s.peek() {
    result.Literal += strings.ToLower(ch)
    result.Length++
    result.EndColumn++ // can't have line breaks in numbers 
  }
  return
}
