package parser

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

func isExpressionToken(t Token) bool {
  return t == THE || t == WORD || t == NUMBER || t == QUOTE || t == MINUS || t == DQUOTE || t == TRUE || t == FALSE
}
type Scanner struct {
	r *bufio.Reader
  offset int
  line int
  column int
  t *Stack
  debugBuffer string
  preserveCase bool
  minOffset int
}

type Tokenized struct {
  *Span
  Token
  Literal string
  IsPlural bool
  Keyword bool
}

func NewCliffScanner(r io.Reader) *Scanner {
  return &Scanner{r: bufio.NewReader(r), t: NewStack()}
}

func (s *Scanner) read() rune {
  res, _, e := s.r.ReadRune()
  if e != nil {
    return eof
  }
  s.offset++
  if res != '\n' {
    s.debugBuffer = s.debugBuffer + string(res)
    s.column++
  } else {
    s.debugBuffer = ""
    s.line++
    s.column = 0
  }
  return res
}

func (s *Scanner) peek() rune {
  ch, _, e := s.r.ReadRune()
  if e != nil {
    return eof
  }
  s.r.UnreadRune()
  return ch
}

func (s *Scanner) Peek() *Tokenized {
  if (s.t.Len() == 0) {
    s.t.Push(s.Scan())
  }
  return s.t.Peek().(*Tokenized)
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
  if s.t.Len() > 0 {
    result = s.t.Pop().(*Tokenized)
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
  if (s.t.Len() > 0) {
    rc := s.t.Peek().(*Tokenized)
    if (rc.Token == WS) {
      result = s.t.Pop().(*Tokenized)
      return
    }
    result = &Tokenized{
      Token: WS,
      Span: s.Position(),
    }
    return
  }

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
    Debug: s.debugBuffer,
  }
}

func (s *Scanner) scanKeywords() (toks []*Tokenized) {
  toks = make([]*Tokenized, 0)
  for tok := s.Peek(); tok.Keyword; tok = s.Peek() {
    toks = append(toks, s.Scan())
  }
  return
}

func (s *Scanner) scanWords() (toks []*Tokenized) {
  toks = make([]*Tokenized, 0)
  for tok := s.Peek(); tok.Token == WS || tok.Token == WORD; tok = s.Peek() {
    if tok.Token == WORD {
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
  "the": THE,
  "is": IS,
  "are": ARE,
  "when": WHEN,
  "then": THEN,
  "and": AND,
  "true": TRUE,
  "false": FALSE,
  "or": OR,
  "of": OF,
  "nd": ND,
  "rd": RD,
  "th": TH,
  "=" : EQL,
  "first": FIRST,
  "last": LAST,
  "next": NEXT,
  "after": AFTER,
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
  tok := s.Peek()

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

func (s *Scanner) Error(msg string, args ...interface{}) (*ParserError) {
  return NewParserError(*s.Position(), fmt.Sprintf(msg, args...))
}

func (s *Scanner) GetMinOffset() int {
  return s.minOffset
}

func (s *Scanner) SetMinOffset(minOffset int) {
  s.minOffset = minOffset
}
