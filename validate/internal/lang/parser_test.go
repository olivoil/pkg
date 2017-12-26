package lang_test

import (
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/olivoil/pkg/validate/internal/lang"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	var tests = []struct {
		skip bool
		s    string
		expr lang.Expr
		err  string
	}{
		{
			s: "required",
			expr: &lang.Call{
				Name: "required",
			},
		},
		{
			s: "required,numeric",
			expr: &lang.BinaryExpr{
				Op: lang.AND,
				LHS: &lang.Call{
					Name: "required",
				},
				RHS: &lang.Call{
					Name: "numeric",
				},
			},
		},
		{
			s: "required,!numeric",
			expr: &lang.BinaryExpr{
				Op: lang.AND,
				LHS: &lang.Call{
					Name: "required",
				},
				RHS: &lang.NegativeExpr{
					Expr: &lang.Call{
						Name: "numeric",
					},
				},
			},
		},
		{
			s: "len(3)",
			expr: &lang.Call{
				Name: "len",
				Args: []lang.Expr{
					&lang.NumberLiteral{Val: 3},
				},
			},
		},
		{
			s: "!len(4)",
			expr: &lang.NegativeExpr{
				Expr: &lang.Call{
					Name: "len",
					Args: []lang.Expr{
						&lang.NumberLiteral{Val: 4},
					},
				},
			},
		},
		{
			s: "required,each(required,uuidv4)",
			expr: &lang.BinaryExpr{
				Op: lang.AND,
				LHS: &lang.Call{
					Name: "required",
				},
				RHS: &lang.EachExpr{
					Expr: &lang.BinaryExpr{
						Op: lang.AND,
						LHS: &lang.Call{
							Name: "required",
						},
						RHS: &lang.Call{
							Name: "uuidv4",
						},
					},
				},
			},
		},
		{
			s: `required,!numeric,range(4,15),email|phone,match(/(?i)^(.+@example\.com)|(\+1\d{10})$/)`,
			expr: &lang.BinaryExpr{
				Op: lang.AND,
				LHS: &lang.BinaryExpr{
					Op: lang.AND,
					LHS: &lang.BinaryExpr{
						Op: lang.AND,
						LHS: &lang.BinaryExpr{
							Op: lang.AND,
							LHS: &lang.Call{
								Name: "required",
							},
							RHS: &lang.NegativeExpr{
								Expr: &lang.Call{
									Name: "numeric",
								},
							},
						},
						RHS: &lang.Call{
							Name: "range",
							Args: []lang.Expr{
								&lang.NumberLiteral{Val: 4},
								&lang.NumberLiteral{Val: 15},
							},
						},
					},
					RHS: &lang.BinaryExpr{
						Op: lang.OR,
						LHS: &lang.Call{
							Name: "email",
						},
						RHS: &lang.Call{
							Name: "phone",
						},
					},
				},
				RHS: &lang.Call{
					Name: "match",
					Args: []lang.Expr{
						&lang.RegexLiteral{
							Val: regexp.MustCompile(`(?i)^(.+@example\.com)|(\+1\d{10})$`),
						},
					},
				},
			},
		},
		{
			s: `required() AND NOT numeric() AND range(4.0, 15.0) AND email() OR phone() AND match(/(?i)^(.+@example\.com)|(\+1\d{10})$/)`,
			expr: &lang.BinaryExpr{
				Op: lang.AND,
				LHS: &lang.BinaryExpr{
					Op: lang.AND,
					LHS: &lang.BinaryExpr{
						Op: lang.AND,
						LHS: &lang.BinaryExpr{
							Op: lang.AND,
							LHS: &lang.Call{
								Name: "required",
							},
							RHS: &lang.NegativeExpr{
								Expr: &lang.Call{
									Name: "numeric",
								},
							},
						},
						RHS: &lang.Call{
							Name: "range",
							Args: []lang.Expr{
								&lang.NumberLiteral{Val: 4},
								&lang.NumberLiteral{Val: 15},
							},
						},
					},
					RHS: &lang.BinaryExpr{
						Op: lang.OR,
						LHS: &lang.Call{
							Name: "email",
						},
						RHS: &lang.Call{
							Name: "phone",
						},
					},
				},
				RHS: &lang.Call{
					Name: "match",
					Args: []lang.Expr{
						&lang.RegexLiteral{
							Val: regexp.MustCompile(`(?i)^(.+@example\.com)|(\+1\d{10})$`),
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		if tc.skip {
			continue
		}

		t.Run(tc.s, func(t *testing.T) {
			p := lang.NewParser(strings.NewReader(tc.s))
			expr, err := p.Parse(false)

			assert.Equal(t, tc.err, errstring(err))
			if tc.err == "" {
				if expr != nil {
					log.Printf("actual: %s\n", expr.String())
				}
				assert.Equal(t, tc.expr, expr)
			}
		})
	}
}

// errstring converts an error to its string representation.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
