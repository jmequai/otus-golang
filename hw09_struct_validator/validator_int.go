//nolint:exhaustive
package hw09structvalidator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrAssertIntMin   = errors.New("value must not be less than")
	ErrAssertIntMax   = errors.New("value must not be greater than")
	ErrAssertIntRange = errors.New("value not in range")
)

type IntValidator struct {
	Value int64
}

func (v *IntValidator) Validate(
	value reflect.Value,
	field reflect.StructField,
	constraints Constraints,
) ValidationErrors {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.Value = value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.Value = int64(value.Uint())
	default:
		return make(ValidationErrors, 0)
	}

	return assert(v, field.Name, constraints)
}

func (v *IntValidator) AssertMin(constraint string) error {
	i, err := strconv.Atoi(constraint)
	if err != nil {
		return errors.New("invalid constraint (min), int required")
	}

	if v.Value < int64(i) {
		return ErrAssertIntMin
	}

	return nil
}

func (v *IntValidator) AssertMax(constraint string) error {
	i, err := strconv.Atoi(constraint)
	if err != nil {
		return errors.New("invalid constraint (max), int required")
	}

	if v.Value > int64(i) {
		return ErrAssertIntMax
	}

	return nil
}

func (v *IntValidator) AssertIn(constraint string) error {
	for _, s := range strings.Split(constraint, ",") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return errors.New("invalid constraint (in), int required")
		}

		if v.Value == int64(i) {
			return nil
		}
	}

	return ErrAssertIntRange
}

func NewIntValidator() Validator {
	return &IntValidator{}
}
