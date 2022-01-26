package standard

import "testing"

func Test_andOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &andOperator{
					arg1: "",
					arg2: "true",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected boolean",
			},
		},
		{
			input: operatorTestInput{
				operator: &andOperator{
					arg1: "true",
					arg2: "",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected boolean",
			},
		},
		{
			input: operatorTestInput{
				operator: &andOperator{
					arg1: "true",
					arg2: "true",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &andOperator{
					arg1: "true",
					arg2: "false",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_orOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &orOperator{
					arg1: "",
					arg2: "true",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected boolean",
			},
		},
		{
			input: operatorTestInput{
				operator: &orOperator{
					arg1: "true",
					arg2: "",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected boolean",
			},
		},
		{
			input: operatorTestInput{
				operator: &orOperator{
					arg1: "true",
					arg2: "true",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &orOperator{
					arg1: "true",
					arg2: "false",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &orOperator{
					arg1: "false",
					arg2: "false",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_lessThanOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &lessThanOperator{
					arg1: "",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected number",
			},
		},
		{
			input: operatorTestInput{
				operator: &lessThanOperator{
					arg1: "1",
					arg2: "",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected number",
			},
		},
		{
			input: operatorTestInput{
				operator: &lessThanOperator{
					arg1: "1",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &lessThanOperator{
					arg1: "2",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_lessThanOrEqualOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &lessThanOrEqualOperator{
					arg1: "",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected number",
			},
		},
		{
			input: operatorTestInput{
				operator: &lessThanOrEqualOperator{
					arg1: "1",
					arg2: "",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected number",
			},
		},
		{
			input: operatorTestInput{
				operator: &lessThanOrEqualOperator{
					arg1: "1",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &lessThanOrEqualOperator{
					arg1: "2",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &lessThanOrEqualOperator{
					arg1: "2",
					arg2: "1",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_greaterThanOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &greaterThanOperator{
					arg1: "",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected number",
			},
		},
		{
			input: operatorTestInput{
				operator: &greaterThanOperator{
					arg1: "1",
					arg2: "",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected number",
			},
		},
		{
			input: operatorTestInput{
				operator: &greaterThanOperator{
					arg1: "3",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &greaterThanOperator{
					arg1: "2",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_greaterThanOrEqualOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &greaterThanOrEqualOperator{
					arg1: "",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected number",
			},
		},
		{
			input: operatorTestInput{
				operator: &greaterThanOrEqualOperator{
					arg1: "1",
					arg2: "",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected number",
			},
		},
		{
			input: operatorTestInput{
				operator: &greaterThanOrEqualOperator{
					arg1: "3",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &greaterThanOrEqualOperator{
					arg1: "2",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &greaterThanOrEqualOperator{
					arg1: "1",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_equalsOperator(t *testing.T) {
	currentKeySelector, _ := newSelectorOperator("@.key", &ScriptEngine{}, nil)

	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &equalsOperator{
					arg1: nil,
					arg2: "value",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: operatorTestInput{
				operator: &equalsOperator{
					arg1: "value",
					arg2: nil,
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: operatorTestInput{
				operator: &equalsOperator{
					arg1: "value",
					arg2: "value",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &equalsOperator{
					arg1: "value",
					arg2: "other",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator: &equalsOperator{
					arg1: "1",
					arg2: "1.0",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &equalsOperator{
					arg1: 2,
					arg2: "2.0",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &equalsOperator{
					arg1: currentKeySelector,
					arg2: "true",
				},
				paramters: map[string]interface{}{
					"@": map[string]interface{}{
						"key": true,
					},
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &equalsOperator{
					arg1: currentKeySelector,
					arg2: "true",
				},
				paramters: map[string]interface{}{
					"@": map[string]interface{}{
						"key": "true",
					},
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_notEqualsOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &notEqualsOperator{
					arg1: nil,
					arg2: "value",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: operatorTestInput{
				operator: &notEqualsOperator{
					arg1: "value",
					arg2: nil,
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: operatorTestInput{
				operator: &notEqualsOperator{
					arg1: "value",
					arg2: "value",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator: &notEqualsOperator{
					arg1: "value",
					arg2: "other",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &notEqualsOperator{
					arg1: "1",
					arg2: "1.0",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator: &notEqualsOperator{
					arg1: 2,
					arg2: "2.0",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_notOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &notOperator{
					arg: nil,
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &notOperator{
					arg: "value",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected boolean",
			},
		},
		{
			input: operatorTestInput{
				operator: &notOperator{
					arg: "1",
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator: &notOperator{
					arg: 2,
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected boolean",
			},
		},
		{
			input: operatorTestInput{
				operator: &notOperator{
					arg: true,
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator: &notOperator{
					arg: false,
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &notOperator{
					arg: &equalsOperator{arg1: "1", arg2: "1"},
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator: &notOperator{
					arg: &equalsOperator{arg1: "2", arg2: "1"},
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
	}
	batchOperatorTests(t, tests)
}
