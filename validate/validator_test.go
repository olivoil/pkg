package validate_test

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/olivoil/pkg/validate"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidator_Struct_WithCustomValidation(t *testing.T) {
	validator := validate.New(
		validate.WithCustomValidation("odd", validate.SimpleValidationFunc(func(i interface{}) error {
			f, err := strconv.ParseFloat(fmt.Sprintf("%v", i), 64)
			if err != nil {
				return errors.Wrapf(validate.ErrIncompatibleFieldType, "odd requires a number, got %v", i)
			}

			if math.Mod(f, 2) == 0 {
				return fmt.Errorf("expected %v to be odd", i)
			}

			return nil
		})),
	)

	s := struct {
		Count int `validate:"odd"`
	}{
		Count: 5,
	}

	err := validator.Struct(s)
	assert.Nil(t, err)
	s.Count = 4
	err = validator.Struct(s)
	assert.NotNil(t, err)
}

func TestValidator_Struct_WithTagname(t *testing.T) {
	validator := validate.New(validate.WithTagname(`v`))

	s := struct {
		Name string `v:"required"`
	}{
		Name: "not_nil",
	}

	err := validator.Struct(s)
	assert.Nil(t, err)
	s.Name = ""
	err = validator.Struct(s)
	assert.NotNil(t, err)
}

func TestValidator_WithValidationRuleRequired_True(t *testing.T) {
	tests := testCases{
		"validation required on all fields": []testCase{
			{
				V: struct {
					Name string
				}{
					Name: "missing",
				},
				IsValid: false,
			},
			{
				V: struct {
					Name string `validate:"-"`
				}{
					Name: "explicitely missing",
				},
				IsValid: true,
			},
			{
				V: struct {
					Name string `validate:"required"`
				}{
					Name: "validated",
				},
				IsValid: true,
			},
		},
	}

	tests.Test(t, validate.New(validate.WithValidationRuleRequired(true)))
}

func TestValidator_WithFieldValidationRequired_False(t *testing.T) {
	tests := testCases{
		"validation optional by default": []testCase{
			{
				V: struct {
					Name string
				}{
					Name: "missing",
				},
				IsValid: true,
			},
			{
				V: struct {
					Name string `validate:""`
				}{
					Name: "missing",
				},
				IsValid: true,
			},
			{
				V: struct {
					Name string `validate:"-"`
				}{
					Name: "explicitely missing",
				},
				IsValid: true,
			},
		},
	}

	tests.Test(t, validate.New())
}
