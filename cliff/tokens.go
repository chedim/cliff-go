package cliff;

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

	// KEYWORDS
  ARE
	IS
	WHEN
  AND
  OR
  OF
  THEN
)

