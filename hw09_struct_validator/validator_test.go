package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

type UserRole string

// Test the function on different structures and other types.
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:    "Lorem ipsum dolor sit amet accumsan.",
				Age:   19,
				Email: "mail@mail.ru",
				Role:  "stuff",
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "Lorem ipsum dolor sit amet accumsan.",
				Age:   17,
				Email: "mail@mail.ru",
				Role:  "stufff",
			},
			expectedErr: ErrLessThanMinimum,
		},
		{
			in: User{
				ID:    "Lorem ipsum dolor sit amet accumsan.",
				Age:   51,
				Email: "mail@mail.ru",
				Role:  "stufff",
			},
			expectedErr: ErrGreaterThanMaximum,
		},
		{
			in: User{
				ID:    "Lorem ipsum dolor",
				Age:   18,
				Email: "mail@mail.ru",
				Role:  "stuff",
			},
			expectedErr: ErrLengthDoesNotMatch,
		},
		{
			in: User{
				ID:    "Lorem ipsum dolor sit amet accumsan.",
				Age:   18,
				Email: "mailmail.ru",
				Role:  "stuff",
			},
			expectedErr: ErrRegexpDoesNotMatch,
		},
		{
			in:          5,
			expectedErr: ErrNotAStructure,
		},
		{
			in: User{
				ID:    "Lorem ipsum dolor sit amet accumsan.",
				Age:   18,
				Email: "mail@mail.ru",
				Role:  "moderator",
			},
			expectedErr: ErrNotInRange,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr != nil || err != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			}
			_ = tt
		})
	}
}
