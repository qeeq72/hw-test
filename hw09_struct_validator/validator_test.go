package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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

	CustomStruct struct {
		UIVar        uint        `validate:"in:12,112,1112"`
		UI16Slice    []uint16    `validate:"min:10"`
		I32Slice     []int32     `validate:"max:22|in:3,13,23"`
		ByVar        byte        `validate:"max:254"`
		ItfVar       interface{} `validate:"in:text,123"`
		PtrIVar      *int        `validate:"min:0|max:10"`
		PtrStringVar *string     `validate:"regexp:\\d+"`
	}

	StructWithInvalidTag struct {
		Name    string `validate:"len:5:10"`
		Surname string `validate:"in:Ivanov,Petrov"`
	}
)

func TestValidate(t *testing.T) {
	var iVar int
	var sVar string

	iVar = 115
	sVar = "abc123"

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:    "1234",
				Name:  "John",
				Age:   1,
				Email: "smth@gmail.com",
				Role:  "administrator",
				Phones: []string{
					"79058220300",
					"+79058220301",
				},
				meta: json.RawMessage{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			expectedErr: ValidationErrors{
				{"ID", ErrStringHasWrongLength},
				{"Age", ErrValueIsUnderRange},
				{"Role", ErrStringIsNotInList},
				{"Phones[1]", ErrStringHasWrongLength},
			},
		},
		{
			in: &App{
				Version: "PrgVers",
			},
			expectedErr: ValidationErrors{
				{"Version", ErrStringHasWrongLength},
			},
		},
		{
			in: Response{
				Code: 300,
				Body: "Something happens!",
			},
			expectedErr: ValidationErrors{
				{"Code", ErrValueIsNotInList},
			},
		},
		{
			in: CustomStruct{
				UIVar:        111,
				UI16Slice:    []uint16{11, 12, 5},
				I32Slice:     []int32{1, 3, 13, 23, 22},
				ByVar:        255,
				ItfVar:       []interface{}{"text"},
				PtrIVar:      &iVar,
				PtrStringVar: &sVar,
			},
			expectedErr: ValidationErrors{
				{"UIVar", ErrUnsupportedType},
				{"UI16Slice", ErrUnsupportedType},
				{"I32Slice[0]", ErrValueIsNotInList},
				{"I32Slice[3]", ErrValueIsAboveRange},
				{"I32Slice[4]", ErrValueIsNotInList},
				{"ByVar", ErrUnsupportedType},
				{"ItfVar", ErrUnsupportedType},
				{"PtrIVar", ErrUnsupportedType},
				{"PtrStringVar", ErrUnsupportedType},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.True(t, errors.As(err, &ValidationErrors{}))
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNonValidationCases(t *testing.T) {
	t.Run("invalid validator tag", func(t *testing.T) {
		test := StructWithInvalidTag{
			Name:    "Ivan",
			Surname: "Ivanov",
		}
		err := Validate(test)
		require.True(t, errors.As(err, &ErrInvalidValidatorTag))
	})

	t.Run("not a structure value", func(t *testing.T) {
		var test UserRole
		err := Validate(test)
		require.True(t, errors.Is(err, ErrNotStructureValue))
	})

	t.Run("nil pointer", func(t *testing.T) {
		var test *UserRole
		err := Validate(test)
		require.False(t, errors.Is(err, ErrNotStructureValue))
		require.True(t, errors.Is(err, ErrNilPointer))
	})
}
