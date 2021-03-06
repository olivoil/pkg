// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package validate

import "github.com/pkg/errors"

// SimpleBoolValidationFunc is a custom validation function that can be applied to a bool value with no arguments.
type SimpleBoolValidationFunc func(i bool) error

// Validate implements Validation.
func (v SimpleBoolValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(bool)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a bool", i)
	}

	return v(t1)
}

// SimpleByteValidationFunc is a custom validation function that can be applied to a byte value with no arguments.
type SimpleByteValidationFunc func(i byte) error

// Validate implements Validation.
func (v SimpleByteValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(byte)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a byte", i)
	}

	return v(t1)
}

// SimpleComplex128ValidationFunc is a custom validation function that can be applied to a complex128 value with no arguments.
type SimpleComplex128ValidationFunc func(i complex128) error

// Validate implements Validation.
func (v SimpleComplex128ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(complex128)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a complex128", i)
	}

	return v(t1)
}

// SimpleComplex64ValidationFunc is a custom validation function that can be applied to a complex64 value with no arguments.
type SimpleComplex64ValidationFunc func(i complex64) error

// Validate implements Validation.
func (v SimpleComplex64ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(complex64)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a complex64", i)
	}

	return v(t1)
}

// SimpleErrorValidationFunc is a custom validation function that can be applied to a error value with no arguments.
type SimpleErrorValidationFunc func(i error) error

// Validate implements Validation.
func (v SimpleErrorValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(error)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a error", i)
	}

	return v(t1)
}

// SimpleFloat32ValidationFunc is a custom validation function that can be applied to a float32 value with no arguments.
type SimpleFloat32ValidationFunc func(i float32) error

// Validate implements Validation.
func (v SimpleFloat32ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(float32)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a float32", i)
	}

	return v(t1)
}

// SimpleFloat64ValidationFunc is a custom validation function that can be applied to a float64 value with no arguments.
type SimpleFloat64ValidationFunc func(i float64) error

// Validate implements Validation.
func (v SimpleFloat64ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(float64)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a float64", i)
	}

	return v(t1)
}

// SimpleIntValidationFunc is a custom validation function that can be applied to a int value with no arguments.
type SimpleIntValidationFunc func(i int) error

// Validate implements Validation.
func (v SimpleIntValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(int)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a int", i)
	}

	return v(t1)
}

// SimpleInt16ValidationFunc is a custom validation function that can be applied to a int16 value with no arguments.
type SimpleInt16ValidationFunc func(i int16) error

// Validate implements Validation.
func (v SimpleInt16ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(int16)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a int16", i)
	}

	return v(t1)
}

// SimpleInt32ValidationFunc is a custom validation function that can be applied to a int32 value with no arguments.
type SimpleInt32ValidationFunc func(i int32) error

// Validate implements Validation.
func (v SimpleInt32ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(int32)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a int32", i)
	}

	return v(t1)
}

// SimpleInt64ValidationFunc is a custom validation function that can be applied to a int64 value with no arguments.
type SimpleInt64ValidationFunc func(i int64) error

// Validate implements Validation.
func (v SimpleInt64ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(int64)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a int64", i)
	}

	return v(t1)
}

// SimpleInt8ValidationFunc is a custom validation function that can be applied to a int8 value with no arguments.
type SimpleInt8ValidationFunc func(i int8) error

// Validate implements Validation.
func (v SimpleInt8ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(int8)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a int8", i)
	}

	return v(t1)
}

// SimpleRuneValidationFunc is a custom validation function that can be applied to a rune value with no arguments.
type SimpleRuneValidationFunc func(i rune) error

// Validate implements Validation.
func (v SimpleRuneValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(rune)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a rune", i)
	}

	return v(t1)
}

// SimpleStringValidationFunc is a custom validation function that can be applied to a string value with no arguments.
type SimpleStringValidationFunc func(i string) error

// Validate implements Validation.
func (v SimpleStringValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(string)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a string", i)
	}

	return v(t1)
}

// SimpleUintValidationFunc is a custom validation function that can be applied to a uint value with no arguments.
type SimpleUintValidationFunc func(i uint) error

// Validate implements Validation.
func (v SimpleUintValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(uint)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a uint", i)
	}

	return v(t1)
}

// SimpleUint16ValidationFunc is a custom validation function that can be applied to a uint16 value with no arguments.
type SimpleUint16ValidationFunc func(i uint16) error

// Validate implements Validation.
func (v SimpleUint16ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(uint16)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a uint16", i)
	}

	return v(t1)
}

// SimpleUint32ValidationFunc is a custom validation function that can be applied to a uint32 value with no arguments.
type SimpleUint32ValidationFunc func(i uint32) error

// Validate implements Validation.
func (v SimpleUint32ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(uint32)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a uint32", i)
	}

	return v(t1)
}

// SimpleUint64ValidationFunc is a custom validation function that can be applied to a uint64 value with no arguments.
type SimpleUint64ValidationFunc func(i uint64) error

// Validate implements Validation.
func (v SimpleUint64ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(uint64)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a uint64", i)
	}

	return v(t1)
}

// SimpleUint8ValidationFunc is a custom validation function that can be applied to a uint8 value with no arguments.
type SimpleUint8ValidationFunc func(i uint8) error

// Validate implements Validation.
func (v SimpleUint8ValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(uint8)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a uint8", i)
	}

	return v(t1)
}

// SimpleUintptrValidationFunc is a custom validation function that can be applied to a uintptr value with no arguments.
type SimpleUintptrValidationFunc func(i uintptr) error

// Validate implements Validation.
func (v SimpleUintptrValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(uintptr)
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a uintptr", i)
	}

	return v(t1)
}

// SimpleInterfaceValidationFunc is a custom validation function that can be applied to a interface{} value with no arguments.
type SimpleInterfaceValidationFunc func(i interface{}) error

// Validate implements Validation.
func (v SimpleInterfaceValidationFunc) Validate(i interface{}, args ...interface{}) error {
	t1, ok := i.(interface{})
	if !ok {
		return errors.Wrapf(ErrIncompatibleFieldType, "expected %v to be a interface{}", i)
	}

	return v(t1)
}
