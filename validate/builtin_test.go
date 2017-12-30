package validate_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/olivoil/pkg/validate"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Title   string
	V       interface{}
	IsValid bool
}

type testCases map[string][]testCase

// Test all the test cases.
func (ts testCases) Test(t *testing.T, validator *validate.Validator) {
	t.Helper()

	for section, tcs := range ts {
		t.Run(section, func(t *testing.T) {
			for i, tc := range tcs {
				title := tc.Title
				if title == "" {
					title = fmt.Sprintf("%02d", i+1)
				}
				t.Run(title, func(t *testing.T) {
					err := validator.Struct(tc.V)
					if tc.IsValid {
						msg := ""
						if err != nil {
							msg = err.Error()
						}
						assert.NoError(t, err, msg)
					} else {
						assert.Error(t, err)
					}
				})
			}
		})
	}
}

func TestBuiltin_Required(t *testing.T) {
	tests := testCases{
		`string`: []testCase{
			{
				Title: "nil",
				V: struct {
					Name string `validate:"required"`
				}{},
				IsValid: false,
			},
			{
				Title: "not nil",
				V: struct {
					Name string `validate:"required"`
				}{
					Name: "ok",
				},
				IsValid: true,
			},
		},
		`number`: []testCase{
			{
				Title: "nil",
				V: struct {
					Count int `validate:"required"`
				}{},
				IsValid: false,
			},
			{
				Title: "not nil",
				V: struct {
					Count int `validate:"required"`
				}{
					Count: 2,
				},
				IsValid: true,
			},
		},
		`time`: []testCase{
			{
				Title: "nil",
				V: struct {
					Time time.Time `validate:"required"`
				}{},
				IsValid: false,
			},
			{
				Title: "not nil",
				V: struct {
					Time time.Time `validate:"required"`
				}{
					Time: time.Now(),
				},
				IsValid: true,
			},
		},
	}

	tests.Test(t, validate.New())
}

func TestBuiltin_Match(t *testing.T) {
	type V struct {
		Regexp *regexp.Regexp `validate:"required"`
		Name   string         `validate:"match($.Regexp)"`
	}

	tests := testCases{
		`regex`: []testCase{
			{
				Title: "empty value",
				V: V{
					Regexp: regexp.MustCompile(`^\d+$`),
				},
				IsValid: false,
			},
			{
				Title: "matches the regexp",
				V: V{
					Regexp: regexp.MustCompile(`^\d+$`),
					Name:   "123345",
				},
				IsValid: true,
			},
			{
				Title: "does not match the regexp",
				V: V{
					Regexp: regexp.MustCompile(`^\d+$`),
					Name:   "123a345",
				},
				IsValid: false,
			},
			{
				Title: "with missing required value",
				V: struct {
					Regexp *regexp.Regexp `validate:"required"`
					Name   string         `validate:"required,match($.Regexp)"`
				}{
					Regexp: regexp.MustCompile(`^\d+$`),
				},
				IsValid: false,
			},
		},
	}

	tests.Test(t, validate.New())
}

func TestBuiltin_Len(t *testing.T) {
	type Strings struct {
		Names []string `validate:"len(3),each(len(2))"`
	}

	tests := testCases{
		`strings`: []testCase{
			{
				Title: "works",
				V: Strings{
					Names: []string{"a1", "b2", "c3"},
				},
				IsValid: true,
			},
			{
				Title: "not enough items",
				V: Strings{
					Names: []string{"a1", "b2"},
				},
				IsValid: false,
			},
			{
				Title: "too many items",
				V: Strings{
					Names: []string{"a1", "b2", "c3", "d4"},
				},
				IsValid: false,
			},
			{
				Title: "items of wrong length",
				V: Strings{
					Names: []string{"a1", "b2", "c34"},
				},
				IsValid: false,
			},
		},
	}

	tests.Test(t, validate.New())
}

func TestBuiltin_Blacklist(t *testing.T) {
	type V struct {
		Name string `validate:"blacklist('tom','dan')"`
	}

	tests := testCases{
		`strings`: []testCase{
			{
				Title: "allows any other values",
				V: V{
					Name: "tam",
				},
				IsValid: true,
			},
			{
				Title: "denies blacklisted values",
				V: V{
					Name: "dan",
				},
				IsValid: false,
			},
		},
	}

	tests.Test(t, validate.New())
}

func TestBuiltin_Whitelist(t *testing.T) {
	type V struct {
		Name  string   `validate:"whitelist('tom','dan')"`
		Names []string `validate:"each(whitelist('tom','dan'))"`
		Count int      `validate:"whitelist(1,3,5,7)"`
	}

	tests := testCases{
		`strings`: []testCase{
			{
				Title: "works",
				V: V{
					Name:  "dan",
					Names: []string{"tom", "dan", "dan", "tom"},
					Count: 3,
				},
				IsValid: true,
			},
			{
				Title: "no name",
				V: V{
					Names: []string{"tom", "dan"},
					Count: 3,
				},
				IsValid: false,
			},
			{
				Title: "extra names",
				V: V{
					Name:  "tom",
					Names: []string{"tom", "dan", "mike"},
					Count: 3,
				},
				IsValid: false,
			},
			{
				Title: "even count",
				V: V{
					Name:  "tom",
					Names: []string{"tom", "dan"},
					Count: 2,
				},
				IsValid: false,
			},
		},
	}

	tests.Test(t, validate.New())
}

func TestBuiltin_Nil_Or(t *testing.T) {
	type V struct {
		Name  string   `validate:"nil or whitelist('tom','dan')"`
		Names []string `validate:"nil|each(whitelist('tom','dan'))"`
		Count int      `validate:"nil or whitelist(1,3,5,7)"`
	}

	tests := testCases{
		`strings`: []testCase{
			{
				Title: "works with correct values",
				V: V{
					Name:  "dan",
					Names: []string{"tom", "dan"},
					Count: 3,
				},
				IsValid: true,
			},
			{
				Title:   "works with zero values",
				V:       V{},
				IsValid: true,
			},
			{
				Title: "fails with incorrect values",
				V: V{
					Name:  "dom",
					Names: []string{"tom", "dan", "mike"},
					Count: 8,
				},
				IsValid: false,
			},
		},
	}

	tests.Test(t, validate.New())
}

func TestBuiltin_LessThan(t *testing.T) {
	type V struct {
		Float64 float64 `validate:"lt(5)"`
		Float32 float32 `validate:"lt(5)"`
		Int64   int64   `validate:"lt(5)"`
	}

	tests := testCases{
		`valid`: []testCase{
			{
				Title:   "allows nil",
				V:       V{},
				IsValid: true,
			},
			{
				Title: "allows formatted date",
				V: V{
					Timestamp: time.Now().Format(time.RFC3339),
				},
				IsValid: true,
			},
		},
		`invalid`: []testCase{
			{
				Title: "denies other date formats",
				V: V{
					Timestamp: time.Now().Format(time.Kitchen),
				},
				IsValid: false,
			},
			{
				Title: "denies other date formats",
				V: V{
					Timestamp: time.Now().Format(time.RFC822Z),
				},
				IsValid: false,
			},
		},
	}

	tests.Test(t, validate.New())
}

func TestBuiltin_RFC3339(t *testing.T) {
	type V struct {
		Timestamp string `validate:"nil or rfc3339"`
	}

	tests := testCases{
		`valid`: []testCase{
			{
				Title:   "allows nil",
				V:       V{},
				IsValid: true,
			},
			{
				Title: "allows formatted date",
				V: V{
					Timestamp: time.Now().Format(time.RFC3339),
				},
				IsValid: true,
			},
		},
		`invalid`: []testCase{
			{
				Title: "denies other date formats",
				V: V{
					Timestamp: time.Now().Format(time.Kitchen),
				},
				IsValid: false,
			},
			{
				Title: "denies other date formats",
				V: V{
					Timestamp: time.Now().Format(time.RFC822Z),
				},
				IsValid: false,
			},
		},
	}

	tests.Test(t, validate.New())
}
