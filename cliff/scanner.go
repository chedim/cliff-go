package cliff

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
)

var eof = rune(0)
var plc = pluralize.NewClient()

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
  offset int
  line int
  column int
  t *Tokenized
  lineBuffer *string
  peekedLine *string
  preserveCase bool
}

type Tokenized struct {
  *Span
  Token
  Literal string
  IsPlural bool
  Keyword bool
}

func NewCliffScanner(r io.Reader) *Scanner {
  return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
  res, _, e := s.r.ReadRune()
  if e != nil {
    return eof
  }
  s.offset++
  if res != '\n' {
    s.column++
  } else {
    s.line++
    s.column = 0
  }
  return res
}

func (s *Scanner) peek() rune {
  if s.lineBuffer != nil && s.column < len(*s.lineBuffer) {
    return []rune(*s.lineBuffer)[s.column]
  }

  ch, _, e := s.r.ReadRune()
  if e != nil {
    return eof
  }
  s.r.UnreadRune()
  return ch
}

func (s *Scanner) Peek() (*Tokenized, error) {
  if s.t != nil {
    return s.t, nil
  }
  s.t = s.Scan()
  return s.t, nil
}

var specialCharacters = map[rune]Token {
  '*' : ASTERISK,
  ',' : COMMA,
  ':' : COLON,
  ';' : SEMICOLON,
  '+' : PLUS,
  '-' : MINUS,
  '.' : DOT,
  '\'': QUOTE,
  '"' : DQUOTE,
  '[' : LBRA,
  ']' : RBRA,
  '(' : LPAREN,
  ')' : RPAREN,
  '{' : LCURL,
  '}' : RCURL,
  '\n': EOL,
  '\\': SLASH,
}

func (s *Scanner) Scan() (result *Tokenized) {
  if s.t != nil {
    result = s.t
    s.t = nil
    return
  }
  ch := s.peek()

  if ch == eof {
    return &Tokenized{Span: s.Position(), Token: EOF}
  } else if isWhitespace(ch) {
    return s.scanWhitespace()
  } else if isLetter(ch) {
    return s.scanWord()
  } else if isNumber(ch) {
    return s.scanNumber()
  } else {
    result = &Tokenized{
      Span: s.Position(),
      Literal: string(ch),
      Token: specialCharacters[ch],
    }
    result.Length = 1
    result.EndColumn++
    s.read()
  }
  return
}

func (s *Scanner) scanWhitespace() (result *Tokenized) {
  result = new(Tokenized)
  result.Token = WS
  result.Span = s.Position()

  for ch := s.peek(); isWhitespace(ch); ch = s.peek() {
    result.Length++
    result.EndColumn++

    if s.preserveCase {
      result.Literal += string(ch)
    } else {
      result.Literal += strings.ToLower(string(ch))
    }
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
  fmt.Printf("scanning words\n")
  for tok, e := s.Peek(); e == nil && (tok.Token == WS || tok.Token == WORD); tok, e = s.Peek() {
    if tok.Token == WORD {
      fmt.Printf("Scanwords: %s %s\n", tok.Token, tok.Literal)
      toks = append(toks, tok)
    }
    s.Scan()
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

    result.Length++
    result.EndColumn++ // can't have line breaks in words
    if s.preserveCase {
      result.Literal += string(ch)
    } else {
      result.Literal += strings.ToLower(string(ch))
    }

    s.read()
  }
  result.IsPlural = plc.IsPlural(result.Literal)
  return detectKeyword(result)
}

var Keywords = map[string]Token{
  "a": A,
  "is": IS,
  "are": ARE,
  "when": WHEN,
  "then": THEN,
  "and": AND,
  "or": OR,
  "of": OF,
  "nd": ND,
  "rd": RD,
  "th": TH,
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
    result.Literal += string(ch)
    result.Length++
    result.EndColumn++ // can't have line breaks in numbers 
    s.read()
  }
  return
}

func (s *Scanner) scanOffset(size int) (bool, *ParserError) {
  tok, e := s.Peek()
  if e != nil {
    return false, ExtendParserError(*s.Position(), e)
  }

  if tok.Token != WS || tok.Length < size {
    return false, nil
  }

  for i := 0; i < size; i++ {
    if tok.Literal[i] != ' ' {
      return false, NewParserError(*s.Position(), "Non-space character in offset")
    }
  }

  if tok.Length == size {
    s.Scan()
  } else {
    tok.StartColumn += size
    tok.Literal = tok.Literal[size:]
  }

  return true, nil
}

func (s *Scanner) PreserveCase(pc bool) (old bool) {
  old = s.preserveCase
  s.preserveCase = pc
  return
}
