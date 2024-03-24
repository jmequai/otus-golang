//nolint:lll
package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type ExpectedError struct {
	Field string
	Error error
}

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	SliceOfInt struct {
		Ints []int `validate:"min:0|max:100"`
	}

	Failed struct {
		Code int `validate:"length:10"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected []ExpectedError
	}{
		{
			App{Version: "1.0.0"},
			[]ExpectedError{},
		},
		{
			App{Version: "1.0.0-beta"},
			[]ExpectedError{newErr("Version", ErrAssertStrLen)},
		},

		{
			User{ID: "4f15371d-7f86-4b40-aaa7-8419208db1c5", Age: 25, Email: "email@email.com", Role: "admin"},
			[]ExpectedError{},
		},
		{
			User{ID: "4f15371d-7f86-4b40-aaa7-8419208db1c5", Age: 25, Email: "email@email.com", Role: "admin", Phones: []string{"+0123456789"}},
			[]ExpectedError{},
		},
		{
			User{ID: "123-456", Age: 10, Email: "email_email.com", Role: "guest"},
			[]ExpectedError{newErr("ID", ErrAssertStrLen), newErr("Age", ErrAssertIntMin), newErr("Email", ErrAssertStrRegex), newErr("Role", ErrAssertStrRange)},
		},
		{
			User{ID: "4f15371d-7f86-4b40-aaa7-8419208db1c5", Age: 100, Email: "email@email.com", Role: "admin", Phones: []string{"123-456", "789-000"}},
			[]ExpectedError{newErr("Age", ErrAssertIntMax), newErr("Phones", ErrAssertStrLen), newErr("Phones", ErrAssertStrLen)},
		},

		{
			Token{},
			[]ExpectedError{},
		},

		{
			Response{Code: 200},
			[]ExpectedError{},
		},
		{
			Response{Code: 100},
			[]ExpectedError{newErr("Code", ErrAssertIntRange)},
		},

		{
			SliceOfInt{Ints: []int{0, 1, 10, 50, 100}},
			[]ExpectedError{},
		},
		{
			SliceOfInt{Ints: []int{200, 300}},
			[]ExpectedError{newErr("Ints", ErrAssertIntMax), newErr("Ints", ErrAssertIntMax)},
		},

		{
			Failed{Code: 100},
			[]ExpectedError{newErr("Code", ErrConstraintInvalid)},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if len(tt.expected) == 0 {
				require.NoError(t, err)
			} else {
				var errs ValidationErrors

				if errors.As(err, &errs) {
					for i, e := range tt.expected {
						require.Equal(t, e.Field, errs[i].Field)
						require.ErrorIs(t, e.Error, errs[i].Err)
					}
				}
			}
		})
	}
}

func TestNotStruct(t *testing.T) {
	var errs ValidationErrors

	err := Validate(1)

	if errors.As(err, &errs) && len(errs) > 0 {
		require.ErrorIs(t, ErrStructRequired, errs[0].Err)
	}
}

func newErr(f string, e error) ExpectedError {
	return ExpectedError{Field: f, Error: e}
}
