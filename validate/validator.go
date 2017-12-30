package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/olivoil/pkg/validate/internal/lang"
	"github.com/pkg/errors"
)

// Validator allows customization of the validation behavior.
type Validator struct {
	validations            Validations
	validationRuleRequired bool
	tagname                string
	err                    error
}

func (v *Validator) registerValidation(name string, f Validation) {
	v.validations.Set(name, f)
}

func (v *Validator) setValidationRuleRequired(required bool) {
	v.validationRuleRequired = required
}

func (v *Validator) setTagname(name string) {
	if name != "" {
		v.tagname = name
	}
}

// Option customizes a validator's behavior.
type Option func(*Validator)

// WithValidationRuleRequired causes validations to fail when struct fields
// do not include validations or are not explicitly marked as exempt (using `validate:"-"` or `validate:"email,omitempty"`).
func WithValidationRuleRequired(required bool) Option {
	return func(v *Validator) {
		v.setValidationRuleRequired(required)
	}
}

// WithCustomValidation adds a custom validation with tag `name`.
func WithCustomValidation(name string, f Validation) Option {
	return func(v *Validator) {
		v.registerValidation(name, f)
	}
}

// WithTagname changes the tag name used to set each struct field validation.
// By default, the tagname is `validate`
func WithTagname(name string) Option {
	return func(v *Validator) {
		v.setTagname(name)
	}
}

// New returns a validator with the specified options.
func New(options ...Option) *Validator {
	validator := &Validator{
		validations: builtins,
		tagname:     "validate",
	}

	for _, o := range options {
		o(validator)
	}

	return validator
}

// validate a value against a expression, optionally within a bounded context `s`.
// @param expr parsed AST expression for the validation rule to validate
// @param val value to validate
// @param s bounded context (typically the struct being validated)
func (v *Validator) validate(name string, expr lang.Expr, val, s reflect.Value) error {
	switch exp := expr.(type) {
	case *lang.BinaryExpr:
		fmt.Printf("checking %s\n", exp.String())
		fmt.Printf("LHS: %s\n", exp.LHS.String())
		fmt.Printf("Op: %s\n", exp.Op.String())
		fmt.Printf("RHS: %s\n\n", exp.RHS.String())

		err := v.validate(name, exp.LHS, val, s)
		if err != nil && exp.Op == lang.AND {
			return err
		}
		if err == nil && exp.Op == lang.OR {
			return nil
		}

		return v.validate(name, exp.RHS, val, s)
	case *lang.ParenExpr:
		return v.validate(name, exp.Expr, val, s)
	case *lang.NegativeExpr:
		if err := v.validate(name, exp.Expr, val, s); err == nil {
			return Error{Field: name, Validation: exp.String()}
		}

		return nil
	case *lang.EachExpr:
		switch reflect.TypeOf(val.Interface()).Kind() {
		case reflect.Slice, reflect.Array, reflect.String:
			for i := 0; i < val.Len(); i++ {
				if err := v.validate(name, exp.Expr, val.Index(i), s); err != nil {
					return err
				}
			}
			return nil
		default:
			fmt.Printf("each(%v)", val.Interface())
			return errors.Wrap(ErrIncompatibleFieldType, "each() requires the value to be an array, slice, or string")
		}
	case *lang.Call:
		// look up validation function
		f, ok := v.validations[exp.Name]
		if !ok {
			return errors.Wrapf(ErrUnknownValidationFunction, "unknown validation: %s", exp.Name)
		}

		// extract parameters
		params := []interface{}{}
		for _, arg := range exp.Args {
			switch a := arg.(type) {
			case *lang.BoundParam:
				p, err := getValueFromStruct(a.Path, s)
				if err != nil {
					return err
				}
				params = append(params, p)
			default:
				if l, ok := arg.(lang.Literal); ok {
					params = append(params, l.Interface())
				}
			}
		}

		// call validation
		if err := f.Validate(val.Interface(), params...); err != nil {
			return Error{Field: name, Validation: exp.String(), Err: err}
		}

		return nil
	}

	return errors.Wrap(ErrUnknownExpression, expr.String())
}

// getValueFromStruct resolves a value from a struct using a path separated by dots (i.e. `Account.User.Name.First`)
func getValueFromStruct(keyWithDots string, v reflect.Value) (interface{}, error) {
	keySlice := strings.Split(keyWithDots, ".")

	for _, key := range keySlice {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// we only accept structs
		if v.Kind() != reflect.Struct {
			return nil, ErrIncompatibleFieldType
		}

		v = v.FieldByName(key)
	}

	return v.Interface(), nil
}

func structFieldName(structField reflect.StructField) string {
	name := jsonName(structField.Tag.Get(`json`))

	if name == "" {
		name = structField.Name
	}

	return name
}

// jsonName returns the name of a `json:"name"` tag.
func jsonName(tag string) string {
	if tag == "" {
		return ""
	}

	split := strings.SplitN(tag, ",", 2)
	name := split[0]

	if name == "-" {
		return ""
	}

	return name
}
