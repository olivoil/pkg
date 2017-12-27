package lang

import "strings"

// Token represents a lexical token.
type Token int

const (
	// ILLEGAL and the following are special tokens.
	ILLEGAL Token = iota
	EOF
	WS

	literalBeg
	// IDENT and the following are literal tokens.
	IDENT      // validation name
	BOUNDPARAM // $param
	NUMBER     // 12.3
	INTEGER    // 12
	DURATION   // 12h
	STRING     // "abc"
	BADSTRING  // "abc
	TRUE       // true
	FALSE      // false
	REGEX      // Regular expressions
	BADESCAPE  // \q
	BADREGEX   // `.*
	literalEnd

	LPAREN // (
	RPAREN // )
	COMMA  // ,
	DOT    // .

	operatorBeg
	// OR and the following are Operators.
	OR  // |
	AND // ,
	NOT // !
	operatorEnd

	keywordBeg
	EACH // each
	keywordEnd
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	// Literals
	IDENT:      "IDENT",
	BOUNDPARAM: "BOUNDPARAM",
	NUMBER:     "NUMBER",
	INTEGER:    "INTEGER",
	DURATION:   "DURATION",
	STRING:     "STRING",
	BADSTRING:  "BADSTRING",
	TRUE:       "TRUE",
	FALSE:      "FALSE",
	REGEX:      "REGEX",
	BADESCAPE:  "BADESCAPE",
	BADREGEX:   "BADREGEX",

	LPAREN: "(",
	RPAREN: ")",
	COMMA:  ",",
	DOT:    ".",

	// Operators
	OR:  "OR",
	AND: "AND",
	NOT: "NOT",

	// Keywords
	EACH: "EACH",
}

var synonyms = map[Token][]string{
	OR:  []string{"OR", "|"},
	AND: []string{"AND", "&&", ","},
	NOT: []string{"NOT", "!"},
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for tok := keywordBeg + 1; tok < keywordEnd; tok++ {
		keywords[strings.ToLower(tokens[tok])] = tok
	}
	for tok, ss := range synonyms {
		for _, s := range ss {
			keywords[strings.ToLower(s)] = tok
		}
	}
	keywords["true"] = TRUE
	keywords["false"] = FALSE
}

// String returns the string representation of the token.
func (tok Token) String() string {
	if tok >= 0 && tok < Token(len(tokens)) {
		return tokens[tok]
	}
	return ""
}

// MarshalText implements encoding.TextMarshaler
func (tok Token) MarshalText() ([]byte, error) {
	return []byte(tok.String()), nil
}

// Precedence returns the operator precedence of the binary operator token.
func (tok Token) Precedence() int {
	switch tok {
	case NOT:
		return 1
	case OR:
		return 2
	case AND:
	case COMMA:
		return 3
	}
	return 0
}

// Lookup returns the token associated with a given string.
func Lookup(ident string) Token {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}

// isOperator returns true for operator tokens.
func (tok Token) isOperator() bool { return tok > operatorBeg && tok < operatorEnd }

// isKeyword returns true for keyword tokens.
func (tok Token) isKeyword() bool { return tok > keywordBeg && tok < keywordEnd }

// isLiteral returns true for literal tokens.
func (tok Token) isLiteral() bool { return tok > literalBeg && tok < literalEnd }

// tokstr returns a literal if provided, otherwise returns the token string.
func tokstr(tok Token, lit string) string {
	if lit != "" {
		return lit
	}
	return tok.String()
}
