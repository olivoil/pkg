package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

// builtin validations.
var builtins = Validations{
	"nil":       SimpleValidationFunc(Nil),
	"required":  SimpleValidationFunc(Required),
	"match":     ValidationFunc(Match),
	"len":       InterfaceValidationWithInt64ArgsFunc(Len),
	"whitelist": ValidationFunc(Whitelist),
	"blacklist": ValidationFunc(Blacklist),
	"lt":        ValidationFunc(LessThan),
	"lte":       ValidationFunc(LessThanOrEqual),
	"gt":        ValidationFunc(GreaterThan),
	"gte":       ValidationFunc(GreaterThanOrEqual),
	"rfc3339":   SimpleStringValidationFunc(RFC3339),
}

// LessThan
func LessThan(i interface{}, args ...interface{}) error {
	switch v := i.(type) {
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		numbers, err := toNumbers(args)
		if err != nil {
			return err
		}

		n, _ := toNumber(v)
		return lessThanfloat64(n, numbers...)
	case time.Duration:
		durations, err := toDurations(args)
		if err != nil {
			return err
		}

		return lessThantimeDuration(v, durations...)
	}

	return errors.Wrapf(ErrIncompatibleFieldType, "lt expects a numeric field type or a duration, got %v", i)
}

func LessThanOrEqual(i interface{}, args ...interface{}) error {
	switch v := i.(type) {
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		numbers, err := toNumbers(args)
		if err != nil {
			return err
		}

		n, _ := toNumber(v)
		return lessThanOrEqualTofloat64(n, numbers...)
	case time.Duration:
		durations, err := toDurations(args)
		if err != nil {
			return err
		}

		return lessThanOrEqualTotimeDuration(v, durations...)
	}

	return errors.Wrapf(ErrIncompatibleFieldType, "lte expects a numeric field type or a duration, got %v", i)
}

func GreaterThan(i interface{}, args ...interface{}) error {
	switch v := i.(type) {
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		numbers, err := toNumbers(args)
		if err != nil {
			return err
		}

		n, _ := toNumber(v)
		return greaterThanfloat64(n, numbers...)
	case time.Duration:
		durations, err := toDurations(args)
		if err != nil {
			return err
		}

		return greaterThantimeDuration(v, durations...)
	}

	return errors.Wrapf(ErrIncompatibleFieldType, "lte expects a numeric field type or a duration, got %v", i)
}

func GreaterThanOrEqual(i interface{}, args ...interface{}) error {
	switch v := i.(type) {
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		numbers, err := toNumbers(args)
		if err != nil {
			return err
		}

		n, _ := toNumber(v)
		return greaterThanOrEqualTofloat64(n, numbers...)
	case time.Duration:
		durations, err := toDurations(args)
		if err != nil {
			return err
		}

		return greaterThanOrEqualTotimeDuration(v, durations...)
	}

	return errors.Wrapf(ErrIncompatibleFieldType, "lte expects a numeric field type or a duration, got %v", i)
}

// Blacklist validates `i` is not one of `args`.
func Blacklist(i interface{}, args ...interface{}) error {
	val, err := normalize(i)
	if err != nil {
		val = i
	}

	for _, arg := range args {
		if val == arg {
			return fmt.Errorf("%v in (%v)", i, args)
		}
	}

	return nil
}

// Whitelist validates `i` is one of `args`.
func Whitelist(i interface{}, args ...interface{}) error {
	val, err := normalize(i)
	if err != nil {
		val = i
	}

	for _, arg := range args {
		if val == arg {
			return nil
		}
	}
	return fmt.Errorf("%v not in (%v)", i, args)
}

// Lengther enables `Len` to be used on custom types.
type Lengther interface {
	Len() int
}

// Len validates `i` has a length of `args[0]`.
func Len(i interface{}, args ...int64) error {
	// reflect the kind of `i`.
	v := reflect.ValueOf(i)
	k := v.Kind()
	if k == reflect.Ptr {
		k = v.Elem().Kind()
	}

	// compare len if i's kind is compatible.
	switch k {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		actual := v.Len()
		expected := int(args[0])
		if actual != expected {
			return fmt.Errorf("expected %v to have length %d, got %d", i, expected, actual)
		}
		return nil
	}

	if l, ok := i.(Lengther); ok {
		actual := l.Len()
		expected := int(args[0])
		if actual != expected {
			return fmt.Errorf("expected %v to have length %d, got %d", i, expected, actual)
		}
		return nil
	}

	return errors.Wrapf(ErrIncompatibleFieldType, "cannot guess length of %v", i)
}

// Match validates `i` is a string that matches the regular expression `args[0]`.
func Match(i interface{}, args ...interface{}) error {
	s, ok := i.(string)
	if !ok {
		return errors.Wrap(ErrIncompatibleFieldType, "match requires a string value")
	}

	if len(args) != 1 {
		return errors.Wrap(ErrInvalidParamType, "match expects a *regexp.Regexp argument")
	}

	r, ok := args[0].(*regexp.Regexp)
	if !ok {
		return errors.Wrap(ErrInvalidParamType, "match expects a *regexp.Regexp argument")
	}

	if !r.MatchString(s) {
		return fmt.Errorf("%v does not match %s", i, r.String())
	}

	return nil
}

// Niller enables custom types to define their zero values.
type Niller interface {
	IsNil() bool
}

// Zeroer enables custom types to define their zero values.
type Zeroer interface {
	IsZero() bool
}

// Required indicates if `i` is not a nil value.
func Required(i interface{}) error {
	if err := Nil(i); err != nil {
		return nil
	}

	return fmt.Errorf("expected %v not to be nil", i)
}

// Nil indicates if `i` is a nil value.
func Nil(i interface{}) error {
	v := reflect.ValueOf(i)
	err := fmt.Errorf("expected %v to be nil", i)

	switch v.Kind() {
	case reflect.String, reflect.Array:
		if v.Len() != 0 {
			return err
		}

		return nil
	case reflect.Map, reflect.Slice:
		if v.Len() != 0 || !v.IsNil() {
			return err
		}

		return nil
	case reflect.Bool:
		if v.Bool() {
			return err
		}

		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() != 0 {
			return err
		}

		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if v.Uint() != 0 {
			return err
		}

		return nil
	case reflect.Float32, reflect.Float64:
		if v.Float() != 0 {
			return err
		}

		return nil
	case reflect.Interface, reflect.Ptr:
		if !v.IsNil() {
			return err
		}

		return nil
	default:
		if n, ok := v.Interface().(Niller); ok {
			if !n.IsNil() {
				return err
			}

			return nil
		}
		if n, ok := v.Interface().(Zeroer); ok {
			if !n.IsZero() {
				return err
			}

			return nil
		}
	}

	if !reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()) {
		return err
	}

	return nil
}

// RFC3339 validates `i` is formatted with time.RFC3339.
func RFC3339(i string) error {
	_, err := time.Parse(time.RFC3339, i)
	return err
}
