//nolint:exhaustive
package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrAssertStrLen   = errors.New("length must be equal to")
	ErrAssertStrRegex = errors.New("value must match with")
	ErrAssertStrRange = errors.New("value not in range")
)

type StrValidator struct {
	Value string
}

func (v *StrValidator) Validate(
	value reflect.Value,
	field reflect.StructField,
	constraints Constraints,
) ValidationErrors {
	switch value.Kind() {
	case reflect.String:
		v.Value = value.String()
	default:
		return make(ValidationErrors, 0)
	}

	return assert(v, field.Name, constraints)
}

func (v *StrValidator) AssertLen(constraint string) error {
	i, err := strconv.Atoi(constraint)
	if err != nil {
		return errors.New("invalid constraint (len), int required")
	}

	if len(v.Value) != i {
		return ErrAssertStrLen
	}

	return nil
}

func (v *StrValidator) AssertRegexp(constraint string) error {
	rgx, err := regexp.Compile(constraint)
	if err != nil {
		return errors.New("invalid constraint (regexp)")
	}

	if !rgx.MatchString(v.Value) {
		return ErrAssertStrRegex
	}

	return nil
}

func (v *StrValidator) AssertIn(constraint string) error {
	for _, s := range strings.Split(constraint, ",") {
		if v.Value == s {
			return nil
		}
	}

	return ErrAssertStrRange
}

func NewStrValidator() Validator {
	return &StrValidator{}
}
