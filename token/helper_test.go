package token

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sampleStruct struct {
	One   string `json:"one"`
	Two   string `json:"two,omitempty"`
	Three int64  `json:"-"`
	Four  int64  `json:"three"`
	Five  string
	Six   string `json:",omitempty"`
}

func Test_isInteger(t *testing.T) {

	type expected struct {
		value int64
		ok    bool
	}

	tests := []struct {
		input    interface{}
		expected expected
	}{
		{
			input: nil,
			expected: expected{
				value: 0,
				ok:    false,
			},
		},
		{
			input: []string{},
			expected: expected{
				value: 0,
				ok:    false,
			},
		},
		{
			input: "string",
			expected: expected{
				value: 0,
				ok:    false,
			},
		},
		{
			input: map[string]string{},
			expected: expected{
				value: 0,
				ok:    false,
			},
		},
		{
			input: 100,
			expected: expected{
				value: 100,
				ok:    true,
			},
		},
		{
			input: 0,
			expected: expected{
				value: 0,
				ok:    true,
			},
		},
		{
			input: int64(100),
			expected: expected{
				value: 100,
				ok:    true,
			},
		},
		{
			input: int(100),
			expected: expected{
				value: 100,
				ok:    true,
			},
		},
		{
			input: int32(100),
			expected: expected{
				value: 100,
				ok:    true,
			},
		},
		{
			input: float64(100),
			expected: expected{
				value: 100,
				ok:    true,
			},
		},
		{
			input: float64(3.14),
			expected: expected{
				value: 0,
				ok:    false,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			value, ok := isInteger(test.input)

			assert.Equal(t, test.expected.ok, ok)
			assert.EqualValues(t, test.expected.value, value)
		})
	}
}

func Test_getStructFields(t *testing.T) {

	type input struct {
		value     reflect.Value
		omitempty bool
	}

	tests := []struct {
		input    input
		expected []string
	}{
		{
			input: input{
				value: reflect.ValueOf(""),
			},
			expected: nil,
		},
		{
			input: input{
				value: reflect.ValueOf(sampleStruct{}),
			},
			expected: []string{
				"one",
				"two",
				"three",
				"Five",
				"Six",
			},
		},
		{
			input: input{
				value:     reflect.ValueOf(sampleStruct{}),
				omitempty: true,
			},
			expected: []string{
				"one",
				"three",
				"Five",
			},
		},
		{
			input: input{
				value: reflect.ValueOf(sampleStruct{
					One:   "one",
					Two:   "two",
					Three: 3,
					Four:  4,
					Five:  "five",
					Six:   "six",
				}),
				omitempty: true,
			},
			expected: []string{
				"one",
				"two",
				"three",
				"Five",
				"Six",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getStructFields(test.input.value, test.input.omitempty)
			if test.expected == nil {
				assert.Nil(t, actual)
			} else {
				keys := make([]string, 0)

				for k := range actual {
					keys = append(keys, k)
				}
				assert.ElementsMatch(t, test.expected, keys)
			}

		})
	}
}

func Test_getTypeAndValue(t *testing.T) {

	getNilPointer := func() *sampleStruct {
		return nil
	}

	sampleString := "sample"

	type expected struct {
		kind  reflect.Kind
		value interface{}
	}

	tests := []struct {
		input    interface{}
		expected expected
	}{
		{
			input: nil,
			expected: expected{
				kind:  reflect.Invalid,
				value: nil,
			},
		},
		{
			input: sampleString,
			expected: expected{
				kind:  reflect.String,
				value: sampleString,
			},
		},
		{
			input: &sampleString,
			expected: expected{
				kind:  reflect.String,
				value: sampleString,
			},
		},
		{
			input: sampleStruct{},
			expected: expected{
				kind:  reflect.Struct,
				value: sampleStruct{},
			},
		},
		{
			input: &sampleStruct{},
			expected: expected{
				kind:  reflect.Struct,
				value: sampleStruct{},
			},
		},
		{
			input: getNilPointer(),
			expected: expected{
				kind:  reflect.Invalid,
				value: nil,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			objType, objValue := getTypeAndValue(test.input)

			if test.expected.kind == reflect.Invalid {
				assert.Nil(t, objType, "expected nil but got %v", objType)
				return
			}

			assert.Equal(t, test.expected.kind, objType.Kind())

			actual := objValue.Interface()
			assert.Equal(t, test.expected.value, actual)
		})
	}

}
