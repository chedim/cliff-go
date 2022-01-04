// Code generated by "stringer -type=Token"; DO NOT EDIT.

package parser

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[EOF-1]
	_ = x[EOL-2]
	_ = x[WS-3]
	_ = x[WORD-4]
	_ = x[NUMBER-5]
	_ = x[ASTERISK-6]
	_ = x[COMMA-7]
	_ = x[COLON-8]
	_ = x[SEMICOLON-9]
	_ = x[PLUS-10]
	_ = x[MINUS-11]
	_ = x[DOT-12]
	_ = x[QUOTE-13]
	_ = x[DQUOTE-14]
	_ = x[LBRA-15]
	_ = x[RBRA-16]
	_ = x[LPAREN-17]
	_ = x[RPAREN-18]
	_ = x[LCURL-19]
	_ = x[RCURL-20]
	_ = x[SLASH-21]
	_ = x[BSLASH-22]
	_ = x[OTHER-23]
	_ = x[A-24]
	_ = x[AN-25]
	_ = x[ARE-26]
	_ = x[AT-27]
	_ = x[IS-28]
	_ = x[WHEN-29]
	_ = x[AND-30]
	_ = x[OR-31]
	_ = x[OF-32]
	_ = x[THEN-33]
	_ = x[TRUE-34]
	_ = x[FALSE-35]
	_ = x[ND-36]
	_ = x[RD-37]
	_ = x[TH-38]
}

const _Token_name = "ILLEGALEOFEOLWSWORDNUMBERASTERISKCOMMACOLONSEMICOLONPLUSMINUSDOTQUOTEDQUOTELBRARBRALPARENRPARENLCURLRCURLSLASHBSLASHOTHERAANAREATISWHENANDOROFTHENTRUEFALSENDRDTH"

var _Token_index = [...]uint8{0, 7, 10, 13, 15, 19, 25, 33, 38, 43, 52, 56, 61, 64, 69, 75, 79, 83, 89, 95, 100, 105, 110, 116, 121, 122, 124, 127, 129, 131, 135, 138, 140, 142, 146, 150, 155, 157, 159, 161}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}