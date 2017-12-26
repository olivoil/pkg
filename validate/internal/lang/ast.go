package lang

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Quote String replacer.
	qsReplacer = strings.NewReplacer("\n", `\n`, `\`, `\\`, `'`, `\'`)

	// Quote Ident replacer.
	qiReplacer = strings.NewReplacer("\n", `\n`, `\`, `\\`, `"`, `\"`)
)

// Expr is a expression.
type Expr interface {
	expr()
	String() string
}

func (*BinaryExpr) expr()     {}
func (*ParenExpr) expr()      {}
func (*NegativeExpr) expr()   {}
func (*EachExpr) expr()       {}
func (*Call) expr()           {}
func (*StringLiteral) expr()  {}
func (*BooleanLiteral) expr() {}
func (*NumberLiteral) expr()  {}
func (*RegexLiteral) expr()   {}

// BinaryExpr represents an operation between two expressions.
type BinaryExpr struct {
	Op  Token // AND/OR
	LHS Expr
	RHS Expr
}

// String returns a string representation of the binary expression.
func (e *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.LHS.String(), e.Op.String(), e.RHS.String())
}

// ParenExpr represents a parenthesized expression.
type ParenExpr struct {
	Expr Expr
}

// String returns a string representation of the parenthesized expression.
func (e *ParenExpr) String() string { return fmt.Sprintf("(%s)", e.Expr.String()) }

// NegativeExpr represents the negation of a expression.
type NegativeExpr struct {
	Expr Expr
}

// String returns a string representation of the parenthesized expression.
func (e *NegativeExpr) String() string { return fmt.Sprintf("(NOT %s)", e.Expr.String()) }

// EachExpr represents an expression that applies to each item in a slice.
type EachExpr struct {
	Expr Expr
}

// String returns a string representation of the parenthesized expression.
func (e *EachExpr) String() string { return fmt.Sprintf("EACH(%s)", e.Expr.String()) }

// Call represents a function call.
type Call struct {
	Name string
	Args []Expr
}

// String returns a string representation of the call.
func (c *Call) String() string {
	// Join arguments.
	var str []string
	for _, arg := range c.Args {
		str = append(str, arg.String())
	}

	// Write function name and args.
	return fmt.Sprintf("%s(%s)", c.Name, strings.Join(str, ", "))
}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Val string
}

// String returns a string representation of the literal.
func (l *StringLiteral) String() string { return QuoteString(l.Val) }

// QuoteString returns a quoted string.
func QuoteString(s string) string {
	return `'` + qsReplacer.Replace(s) + `'`
}

// BooleanLiteral represents a boolean literal.
type BooleanLiteral struct {
	Val bool
}

// String returns a string representation of the literal.
func (l *BooleanLiteral) String() string {
	if l.Val {
		return "true"
	}
	return "false"
}

// NumberLiteral represents a numeric literal.
type NumberLiteral struct {
	Val float64
}

// String returns a string representation of the literal.
func (l *NumberLiteral) String() string { return strconv.FormatFloat(l.Val, 'f', 3, 64) }

// RegexLiteral represents a regular expression.
type RegexLiteral struct {
	Val *regexp.Regexp
}

// String returns a string representation of the literal.
func (r *RegexLiteral) String() string {
	if r.Val != nil {
		return fmt.Sprintf("/%s/", strings.Replace(r.Val.String(), `/`, `\/`, -1))
	}
	return ""
}
