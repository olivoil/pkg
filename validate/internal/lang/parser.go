package lang

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Parser represents an validation field tag parser.
type Parser struct {
	s *bufScanner
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: newBufScanner(r)}
}

// Parse parses an expression string and returns its AST.
func Parse(s string) (Expr, error) {
	return NewParser(strings.NewReader(s)).Parse(false)
}

// MustParse parses an expression string and returns its AST. Panic on error.
func MustParse(s string) Expr {
	expr, err := Parse(s)
	if err != nil {
		panic(err.Error())
	}
	return expr
}

// Parse parses an expression.
func (p *Parser) Parse(call bool) (Expr, error) {
	var err error

	// Dummy root node.
	root := &BinaryExpr{}

	// Parse a non-binary expression type to start.
	// This variable will always be the root of the expression tree.
	root.RHS, err = p.parseUnaryExpr()
	if err != nil {
		return nil, err
	}

	// Loop over operations and unary exprs and build a tree based on precendence.
	for {
		// If the next token is NOT an operator then return the expression.
		op, _, _ := p.ScanIgnoreWhitespace()

		// AND can be expressed as a comma
		if !call && op == COMMA {
			op = AND
		}

		if !op.isOperator() {
			p.Unscan()
			return root.RHS, nil
		}

		// Otherwise parse the next expression.
		var rhs Expr
		if rhs, err = p.parseUnaryExpr(); err != nil {
			return nil, err
		}

		// Find the right spot in the tree to add the new expression by
		// descending the RHS of the expression tree until we reach the last
		// BinaryExpr or a BinaryExpr whose RHS has an operator with
		// precedence >= the operator being added.
		for node := root; ; {
			r, ok := node.RHS.(*BinaryExpr)
			if !ok || r.Op.Precedence() >= op.Precedence() {
				// Add the new expression here and break.
				node.RHS = &BinaryExpr{LHS: node.RHS, RHS: rhs, Op: op}
				break
			}
			node = r
		}
	}
}

// parseUnaryExpr parses an non-binary expression.
func (p *Parser) parseUnaryExpr() (Expr, error) {
	// If the first token is a LPAREN then parse it as its own grouped expression.
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok == LPAREN {
		expr, err := p.Parse(false)
		if err != nil {
			return nil, err
		}

		// Expect an RPAREN at the end.
		if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != RPAREN {
			return nil, newParseError(tokstr(tok, lit), []string{")"}, pos)
		}

		return &ParenExpr{Expr: expr}, nil
	}
	p.Unscan()

	// If the first token is a NOT then parse it as a negative expression.
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok == NOT {
		expr, err := p.parseUnaryExpr()
		if err != nil {
			return nil, err
		}

		return &NegativeExpr{Expr: expr}, nil
	}
	p.Unscan()

	// If the first token is EACH then parse it as its own grouped expression.
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok == EACH {
		tok2, pos2, lit2 := p.ScanIgnoreWhitespace()
		if tok2 == LPAREN {
			expr, err := p.Parse(false)
			if err != nil {
				return nil, err
			}

			// Expect an RPAREN at the end.
			if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != RPAREN {
				return nil, newParseError(tokstr(tok, lit), []string{")"}, pos)
			}

			return &EachExpr{Expr: expr}, nil
		}

		return nil, newParseError(tokstr(tok2, lit2), []string{"("}, pos2)
	}
	p.Unscan()

	// Read next token.
	tok, pos, lit := p.ScanIgnoreWhitespace()
	switch tok {
	case IDENT:
		// return p.parseCall(lit)
		// If the next immediate token is a left parentheses, parse as function call.
		// Otherwise parse as a variable reference.
		if tok0, _, _ := p.Scan(); tok0 == LPAREN {
			return p.parseCall(lit)
		}

		p.Unscan() // Unscan the last token (wasn't an LPAREN)
		// p.Unscan() // Unscan the IDENT token

		// Parse it as a VarRef.
		// return p.ParseVarRef()
		// Parse it as a string.
		// return &Call{Val: lit}, nil
		return &Call{Name: lit}, nil
	case STRING:
		return &StringLiteral{Val: lit}, nil
	case NUMBER:
		v, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			return nil, &ParseError{Message: "unable to parse number", Pos: pos}
		}
		return &NumberLiteral{Val: v}, nil
	// case INTEGER:
	// 	v, err := strconv.ParseInt(lit, 10, 64)
	// 	if err != nil {
	// 		// The literal may be too large to fit into an int64. If it is, use an unsigned integer.
	// 		// The check for negative numbers is handled somewhere else so this should always be a positive number.
	// 		if v, err := strconv.ParseUint(lit, 10, 64); err == nil {
	// 			return &UnsignedLiteral{Val: v}, nil
	// 		}
	// 		return nil, &ParseError{Message: "unable to parse integer", Pos: pos}
	// 	}
	// 	return &IntegerLiteral{Val: v}, nil
	case TRUE, FALSE:
		return &BooleanLiteral{Val: (tok == TRUE)}, nil
	case REGEX:
		re, err := regexp.Compile(lit)
		if err != nil {
			return nil, &ParseError{Message: err.Error(), Pos: pos}
		}
		return &RegexLiteral{Val: re}, nil
	default:
		return nil, newParseError(tokstr(tok, lit), []string{"identifier", "string", "number", "bool"}, pos)
	}
}

// Scan returns the next token from the underlying scanner.
func (p *Parser) Scan() (tok Token, pos int, lit string) { return p.s.Scan() }

// ScanIgnoreWhitespace scans the next non-whitespace and non-comment token.
func (p *Parser) ScanIgnoreWhitespace() (tok Token, pos int, lit string) {
	for {
		tok, pos, lit = p.Scan()
		if tok == WS {
			continue
		}
		return
	}
}

// consumeWhitespace scans the next token if it's whitespace.
func (p *Parser) consumeWhitespace() {
	if tok, _, _ := p.Scan(); tok != WS {
		p.Unscan()
	}
}

// Unscan pushes the previously read token back onto the buffer.
func (p *Parser) Unscan() { p.s.Unscan() }

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Message  string
	Found    string
	Expected []string
	Pos      int
}

// newParseError returns a new instance of ParseError.
func newParseError(found string, expected []string, pos int) *ParseError {
	return &ParseError{Found: found, Expected: expected, Pos: pos}
}

// Error returns the string representation of the error.
func (e *ParseError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s at char %d", e.Message, e.Pos+1)
	}
	return fmt.Sprintf("found %s, expected %s at char %d", e.Found, strings.Join(e.Expected, ", "), e.Pos+1)
}

// parseCall parses a function call.
// This function assumes the function name and LPAREN have been consumed.
func (p *Parser) parseCall(name string) (*Call, error) {
	name = strings.ToLower(name)

	// Parse first function argument if one exists.
	var args []Expr
	re, err := p.parseRegex()
	if err != nil {
		return nil, err
	} else if re != nil {
		args = append(args, re)
	} else {
		// If there's a right paren then just return immediately.
		if tok, _, _ := p.Scan(); tok == RPAREN {
			return &Call{Name: name}, nil
		}
		p.Unscan()

		arg, err := p.Parse(true)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	// Parse additional function arguments if there is a comma.
	for {
		// If there's not a comma, stop parsing arguments.
		if tok, _, _ := p.ScanIgnoreWhitespace(); tok != COMMA {
			p.Unscan()
			break
		}

		re, err := p.parseRegex()
		if err != nil {
			return nil, err
		} else if re != nil {
			args = append(args, re)
			continue
		}

		// Parse an expression argument.
		arg, err := p.Parse(true)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	// There should be a right parentheses at the end.
	if tok, pos, lit := p.Scan(); tok != RPAREN {
		return nil, newParseError(tokstr(tok, lit), []string{")"}, pos)
	}

	return &Call{Name: name, Args: args}, nil
}

// parseRegex parses a regular expression.
func (p *Parser) parseRegex() (*RegexLiteral, error) {
	nextRune := p.peekRune()
	if isWhitespace(nextRune) {
		p.consumeWhitespace()
	}

	// If the next character is not a '/', then return nils.
	nextRune = p.peekRune()
	if nextRune != '/' {
		return nil, nil
	}

	tok, pos, lit := p.s.ScanRegex()

	if tok == BADESCAPE {
		msg := fmt.Sprintf("bad escape: %s", lit)
		return nil, &ParseError{Message: msg, Pos: pos}
	} else if tok == BADREGEX {
		msg := fmt.Sprintf("bad regex: %s", lit)
		return nil, &ParseError{Message: msg, Pos: pos}
	} else if tok != REGEX {
		return nil, newParseError(tokstr(tok, lit), []string{"regex"}, pos)
	}

	re, err := regexp.Compile(lit)
	if err != nil {
		return nil, &ParseError{Message: err.Error(), Pos: pos}
	}

	return &RegexLiteral{Val: re}, nil
}

// peekRune returns the next rune that would be read by the scanner.
func (p *Parser) peekRune() rune {
	r, _, _ := p.s.s.r.ReadRune()
	if r != eof {
		_ = p.s.s.r.UnreadRune()
	}

	return r
}
