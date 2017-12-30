package validate

import (
	"reflect"

	"github.com/olivoil/pkg/validate/internal/lang"
)

// Value validates a single value `i` against a rule `r`.
func Value(i interface{}, r string) error {
	return New().Value(i, r)
}

// Value validates a single value `i` against a rule `r`.
func (v *Validator) Value(i interface{}, rule string) error {
	// check if rule is missing unintentionally.
	if rule == "" && v.validationRuleRequired {
		return ErrMissingValidationRule
	}

	// skip.
	if rule == "" || rule == "-" {
		return nil
	}

	// parse rule into AST
	expr, err := lang.Parse(rule)
	if err != nil {
		return err
	}

	var errs Errors
	if err := v.validate("", expr, reflect.ValueOf(i), reflect.ValueOf(struct{}{})); err != nil {
		switch v := err.(type) {
		case Errors:
			errs = append(errs, v...)
		case Error:
			errs = append(errs, v)
		default:
			return err
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
