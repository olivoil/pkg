package validate

import (
	"errors"
	"fmt"
	"strings"
)

// non-validation errors.
var (
	ErrUnknownExpression         = errors.New("unknown expression")
	ErrMissingValidationRule     = errors.New("missing required validation rule")
	ErrInvalidParamType          = errors.New("invalid param type")
	ErrIncompatibleFieldType     = errors.New("incompatible field type")
	ErrUnknownValidationFunction = errors.New("unknown validation function")
	ErrValidationNotFound        = errors.New("validation not found")
)

// Errors holds one or several validation errors.
type Errors []Error

func (e Errors) Error() string {
	msg := []string{}

	for _, err := range e {
		msg = append(msg, err.Error())
	}

	return strings.Join(msg, "; ")
}

// Error represents a single validation error.
type Error struct {
	// Field indicates the name of the field that failed to validate.
	Field string
	// Validation indicates the tag of the validation that failed.
	Validation string
	// Err is the error from the validation function.
	Err error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s failed the '%s' validation: %s", e.Field, e.Validation, e.Err.Error())
}
