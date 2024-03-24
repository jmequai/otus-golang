package hw09structvalidator

import (
	"reflect"
)

type SliceValidator struct{}

func (v *SliceValidator) Validate(
	value reflect.Value,
	field reflect.StructField,
	constraints Constraints,
) ValidationErrors {
	if value.Kind() == reflect.Slice {
		switch value.Interface().(type) {
		case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64:
			validator := NewIntValidator()

			return v.loop(validator, value, field, constraints)
		case []string:
			validator := NewStrValidator()

			return v.loop(validator, value, field, constraints)
		}
	}

	return make(ValidationErrors, 0)
}

func (v *SliceValidator) loop(
	validator Validator,
	value reflect.Value,
	field reflect.StructField,
	constraints Constraints,
) ValidationErrors {
	valErrors := make(ValidationErrors, 0)

	for i := 0; i < value.Len(); i++ {
		errs := validator.Validate(value.Index(i), field, constraints)

		if len(errs) > 0 {
			valErrors = append(valErrors, errs...)
		}
	}

	return valErrors
}

func NewSliceValidator() Validator {
	return &SliceValidator{}
}
