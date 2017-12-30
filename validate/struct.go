package validate

import (
	"reflect"

	"github.com/olivoil/pkg/validate/internal/lang"
)

// Struct validates all exported fields in a struct `i` against the rules in field tags.
func Struct(i interface{}) error {
	return New().Struct(i)
}

// Struct validates all exported fields in a struct `i` against the rules in field tags.
func (v *Validator) Struct(s interface{}) error {
	if s == nil {
		return nil
	}

	root := reflect.ValueOf(s)
	if root.Kind() == reflect.Interface || root.Kind() == reflect.Ptr {
		root = root.Elem()
	}
	if root.Kind() != reflect.Struct {
		return ErrInvalidParamType
	}

	var errs Errors
	for i := 0; i < root.NumField(); i++ {
		value := root.Field(i)
		structField := root.Type().Field(i)

		// filter out private struct fields
		if structField.PkgPath != "" {
			continue // Private field
		}

		rule := structField.Tag.Get(v.tagname)
		name := structFieldName(structField)

		// check if rule is missing unintentionally.
		if rule == "" && v.validationRuleRequired {
			return ErrMissingValidationRule
		}

		// validate inner struct.
		if value.Kind() == reflect.Struct && rule != "-" {
			if err := v.Struct(value.Interface()); err != nil {
				switch v := err.(type) {
				case Errors:
					errs = append(errs, v...)
				case Error:
					errs = append(errs, v)
				default:
					return err
				}
			}
		}

		// skip.
		if rule == "" || rule == "-" {
			continue
		}

		// parse rule into AST
		expr, err := lang.Parse(rule)
		if err != nil {
			return err
		}

		if err := v.validate(name, expr, value, root); err != nil {
			switch v := err.(type) {
			case Errors:
				errs = append(errs, v...)
			case Error:
				errs = append(errs, v)
			default:
				return err
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
