package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

var (
	ErrStructRequired    = errors.New("only struct is supported")
	ErrConstraintInvalid = errors.New("invalid constraint")
)

type Validator interface {
	Validate(value reflect.Value, field reflect.StructField, constraints Constraints) ValidationErrors
}

type Constraint struct {
	Name  string
	Value string
}

type Constraints []Constraint

type ValidationError struct {
	Field string
	Err   error

	Constraint Constraint
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	strBuilder := strings.Builder{}

	for _, err := range v {
		strBuilder.WriteString(
			fmt.Sprintf("field %q: %s (%s)", err.Field, err.Err, err.Constraint.Value) + "\n",
		)
	}

	return strBuilder.String()
}

func NewValidators() []Validator {
	return []Validator{
		NewIntValidator(),
		NewStrValidator(),
		NewSliceValidator(),
	}
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Struct {
		return ErrStructRequired
	}

	valErrors := make(ValidationErrors, 0)
	validators := NewValidators()

	valueType := value.Type()

	for i := 0; i < valueType.NumField(); i++ {
		val := value.Field(i)

		if !val.CanInterface() {
			continue
		}

		field := valueType.Field(i)
		tag, ok := field.Tag.Lookup("validate")

		if !ok {
			continue
		}

		constraints := parseConstraints(tag)

		for _, validator := range validators {
			errs := validator.Validate(val, field, constraints)

			if len(errs) > 0 {
				valErrors = append(valErrors, errs...)
			}
		}
	}

	if len(valErrors) > 0 {
		return valErrors
	}

	return nil
}

func parseConstraints(tag string) Constraints {
	cons := strings.Split(tag, "|")
	constraints := make(Constraints, 0)

	for _, c := range cons {
		kv := strings.SplitN(c, ":", 2)

		c := Constraint{Name: kv[0], Value: ""}

		if len(kv) == 2 {
			c.Value = kv[1]
		}

		constraints = append(constraints, c)
	}

	return constraints
}

func assert(validator Validator, fieldName string, constraints Constraints) ValidationErrors {
	errs := make(ValidationErrors, 0)

	v := reflect.ValueOf(validator)

	for _, c := range constraints {
		m := v.MethodByName(getAssertMethod(c.Name))

		if !m.IsValid() {
			errs = append(errs, ValidationError{fieldName, ErrConstraintInvalid, c})
			break
		}

		fn, ok := m.Interface().(func(string) error)
		if ok {
			if err := fn(c.Value); err != nil {
				errs = append(errs, ValidationError{fieldName, err, c})
				break
			}
		}
	}

	return errs
}

func getAssertMethod(s string) string {
	n := "Assert"

	for _, r := range s {
		u := string(unicode.ToUpper(r))

		return n + u + s[len(u):]
	}

	return ""
}
