package lang_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/olivoil/pkg/validate/internal/lang"
	"github.com/stretchr/testify/assert"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok lang.Token
		lit string
		pos int
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: lang.EOF},
		{s: `#`, tok: lang.ILLEGAL, lit: `#`},
		{s: `+`, tok: lang.ILLEGAL, lit: `+`},
		{s: `-`, tok: lang.ILLEGAL, lit: `-`},
		{s: `*`, tok: lang.ILLEGAL, lit: `*`},
		{s: `/`, tok: lang.BADREGEX, lit: ``},
		{s: `%`, tok: lang.ILLEGAL, lit: `%`},
		{s: ` `, tok: lang.WS, lit: " "},
		{s: "\t", tok: lang.WS, lit: "\t"},
		{s: "\n", tok: lang.WS, lit: "\n"},
		{s: "\r", tok: lang.WS, lit: "\n"},
		{s: "\r\n", tok: lang.WS, lit: "\n"},
		{s: "\rX", tok: lang.WS, lit: "\n"},
		{s: "\n\r", tok: lang.WS, lit: "\n\n"},
		{s: " \n\t \r\n\t", tok: lang.WS, lit: " \n\t \n\t"},
		{s: " foo", tok: lang.WS, lit: " "},

		// Logical operators
		{s: `AND`, tok: lang.AND},
		{s: `and`, tok: lang.AND},
		{s: `|`, tok: lang.OR},
		{s: `OR`, tok: lang.OR},
		{s: `or`, tok: lang.OR},
		{s: `!`, tok: lang.NOT},
		{s: `NOT`, tok: lang.NOT},
		{s: `not`, tok: lang.NOT},

		// Misc. tokens
		{s: `(`, tok: lang.LPAREN},
		{s: `)`, tok: lang.RPAREN},
		{s: `,`, tok: lang.COMMA},

		// Identifiers
		{s: `required`, tok: lang.IDENT, lit: `required`},
		{s: `required()`, tok: lang.IDENT, lit: `required`},
		{s: `foo`, tok: lang.IDENT, lit: `foo`},
		{s: `phone`, tok: lang.IDENT, lit: `phone`},
		{s: `range(1,2)`, tok: lang.IDENT, lit: `range`},

		{s: `true`, tok: lang.TRUE},
		{s: `false`, tok: lang.FALSE},

		{s: `'testing 123!'`, tok: lang.STRING, lit: `testing 123!`},
		{s: `'string'`, tok: lang.STRING, lit: `string`},
		{s: `'foo\nbar'`, tok: lang.STRING, lit: "foo\nbar"},

		{s: `1`, tok: lang.NUMBER, lit: "1"},
		{s: `2.3`, tok: lang.NUMBER, lit: "2.3"},

		{s: `EACH`, tok: lang.EACH},
		{s: `each(!zero)`, tok: lang.EACH},
	}

	for i, tc := range tests {
		t.Run(tc.s, func(t *testing.T) {
			s := lang.NewScanner(strings.NewReader(tc.s))
			tok, pos, lit := s.Scan()
			assert.Equal(t, tc.tok, tok, fmt.Sprintf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tc.s, tc.tok.String(), tok.String(), lit))
			assert.Equal(t, tc.pos, pos, fmt.Sprintf("%d. %q pos mismatch: exp=%#v got=%#v", i, tc.s, tc.pos, pos))
			assert.Equal(t, tc.lit, lit, fmt.Sprintf("%d. %q literal mismatch: exp=%q got=%q", i, tc.s, tc.lit, lit))
		})
	}
}

type multiScanResult struct {
	tok lang.Token
	pos int
	lit string
}

func (r multiScanResult) String() string {
	return fmt.Sprintf("{token: %s, pos: %d, lit: %s}", r.tok.String(), r.pos, r.lit)
}

// Ensure the scanner can scan a series of tokens correctly.
func TestScanner_Scan_Multi(t *testing.T) {
	tests := map[string][]multiScanResult{
		`required,!contains('example.com'),range(4,15.5),email|phone`: []multiScanResult{
			{tok: lang.IDENT, pos: 0, lit: `required`},
			{tok: lang.COMMA, pos: 8, lit: ``},
			{tok: lang.NOT, pos: 9, lit: ``},
			{tok: lang.IDENT, pos: 10, lit: `contains`},
			{tok: lang.LPAREN, pos: 18, lit: ``},
			{tok: lang.STRING, pos: 18, lit: `example.com`},
			{tok: lang.RPAREN, pos: 32, lit: ``},
			{tok: lang.COMMA, pos: 33, lit: ``},
			{tok: lang.IDENT, pos: 34, lit: `range`},
			{tok: lang.LPAREN, pos: 39, lit: ``},
			{tok: lang.NUMBER, pos: 40, lit: `4`},
			{tok: lang.COMMA, pos: 41, lit: ``},
			{tok: lang.NUMBER, pos: 42, lit: `15.5`},
			{tok: lang.RPAREN, pos: 46, lit: ``},
			{tok: lang.COMMA, pos: 47, lit: ``},
			{tok: lang.IDENT, pos: 48, lit: `email`},
			{tok: lang.OR, pos: 53, lit: ``},
			{tok: lang.IDENT, pos: 54, lit: `phone`},
			{tok: lang.EOF, pos: 60, lit: ``},
		},
		`required and not email`: []multiScanResult{
			{tok: lang.IDENT, pos: 0, lit: `required`},
			{tok: lang.WS, pos: 8, lit: ` `},
			{tok: lang.AND, pos: 9, lit: ``},
			{tok: lang.WS, pos: 12, lit: ` `},
			{tok: lang.NOT, pos: 13},
			{tok: lang.WS, pos: 16, lit: ` `},
			{tok: lang.IDENT, pos: 17, lit: `email`},
			{tok: lang.EOF, pos: 23, lit: ``},
		},
		`required,match(/^payments\./)`: []multiScanResult{
			{tok: lang.IDENT, pos: 0, lit: `required`},
			{tok: lang.COMMA, pos: 8, lit: ``},
			{tok: lang.IDENT, pos: 9, lit: `match`},
			{tok: lang.LPAREN, pos: 14, lit: ``},
			{tok: lang.REGEX, pos: 14, lit: `^payments\.`},
			{tok: lang.RPAREN, pos: 28, lit: ``},
			{tok: lang.EOF, pos: 29, lit: ``},
		},
		`required,each(required)`: []multiScanResult{
			{tok: lang.IDENT, pos: 0, lit: `required`},
			{tok: lang.COMMA, pos: 8, lit: ``},
			{tok: lang.EACH, pos: 9, lit: ``},
			{tok: lang.LPAREN, pos: 13, lit: ``},
			{tok: lang.IDENT, pos: 14, lit: `required`},
			{tok: lang.RPAREN, pos: 22, lit: ``},
			{tok: lang.EOF, pos: 23, lit: ``},
		},
	}

	for v, exp := range tests {
		t.Run(v, func(t *testing.T) {
			s := lang.NewScanner(strings.NewReader(v))

			// Continually scan until we reach the end.
			var act []multiScanResult
			for {
				tok, pos, lit := s.Scan()
				act = append(act, multiScanResult{tok, pos, lit})
				if tok == lang.EOF {
					break
				}
			}

			// Verify the token counts match.
			assert.Len(t, act, len(exp), "token count mismatch: exp=%d, got=%d", len(exp), len(act))

			// Verify each token matches.
			for i := range exp {
				assert.Equal(t, act[i], exp[i], "%d. token mismatch:\n\nexp=token: %s, pos: %d, lit: %s\n\ngot=token: %s, pos: %d, lit: %s", i, exp[i].tok, exp[i].pos, exp[i].lit, act[i].tok, act[i].pos, act[i].lit)
			}
		})
	}
}
