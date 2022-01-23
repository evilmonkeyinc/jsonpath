package standard

import "testing"

func Test_plusOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &plusOperator{
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
				operator: &plusOperator{
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
				operator: &plusOperator{
					arg1: "1",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: float64(3),
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_subtractOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &subtractOperator{
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
				operator: &subtractOperator{
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
				operator: &subtractOperator{
					arg1: "1",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: float64(-1),
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_multiplyOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &multiplyOperator{
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
				operator: &multiplyOperator{
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
				operator: &multiplyOperator{
					arg1: "2",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: float64(4),
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_divideOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &divideOperator{
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
				operator: &divideOperator{
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
				operator: &divideOperator{
					arg1: "4",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: float64(2),
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_modulusOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &modulusOperator{
					arg1: "",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected integer",
			},
		},
		{
			input: operatorTestInput{
				operator: &modulusOperator{
					arg1: "1",
					arg2: "",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected integer",
			},
		},
		{
			input: operatorTestInput{
				operator: &modulusOperator{
					arg1: "3",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: int64(1),
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_powerOfOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &powerOfOperator{
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
				operator: &powerOfOperator{
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
				operator: &powerOfOperator{
					arg1: "3",
					arg2: "2",
				},
			},
			expected: operatorTestExpected{
				value: float64(9),
			},
		},
	}
	batchOperatorTests(t, tests)
}
