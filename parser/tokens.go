package parser;

//go:generate stringer -type=Token
type Token int

const (
	// specials
	ILLEGAL Token = iota
	EOF
	EOL
	WS

	// literals
	WORD
  NUMBER

	// misc
	ASTERISK  // *
	COMMA     // ,
	COLON     // :
	SEMICOLON // ;
	PLUS      // +
	MINUS     // -
  DOT       // .
  QUOTE     // '
  DQUOTE    // "
  LBRA      // [
  RBRA      // ]
  LPAREN    // (
  RPAREN    // )
  LCURL     // {
  RCURL     // }
  SLASH     // /
  BSLASH    // \
  EQL       // =
  OTHER

	// KEYWORDS
  A
  AN
  THE
  ARE
  AT
	IS
	WHEN
  AND
  OR
  OF
  THEN
  TRUE
  FALSE
  ND
  RD
  TH
  FIRST
  LAST
  NEXT
  AFTER
)
