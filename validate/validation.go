package validate

// Validations is a set of `Validation` by tag.
type Validations map[string]Validation

// Get looks up a validation by name
func (v Validations) Get(name string) (Validation, error) {
	validation, ok := v[name]
	if !ok {
		return nil, ErrValidationNotFound
	}

	return validation, nil
}

// Set registers a Validation by name.
func (v Validations) Set(name string, val Validation) {
	v[name] = val
}

// Validation represents a single validation function.
type Validation interface {
	Validate(i interface{}, args ...interface{}) error
}

type ValidationFunc = InterfaceValidationWithInterfaceArgsFunc
type SimpleValidationFunc = SimpleInterfaceValidationFunc