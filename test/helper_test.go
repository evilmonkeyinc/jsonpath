// the test package has tests detailed by https://cburgmer.github.io/json-path-comparison/
// as showing the community consensus on the expected response
package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/evilmonkeyinc/jsonpath"
	"github.com/stretchr/testify/assert"
)

const consensusNone string = "none"

type testData struct {
	query         string
	data          string
	expected      interface{}
	consensus     interface{}
	expectedError string
}

func batchTest(t *testing.T, tests []testData) {
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := jsonpath.QueryString(test.query, test.data)
			if test.expectedError == "" {
				assert.Nil(t, err, fmt.Sprintf("%s error should be nil", test.query))
			} else {
				assert.EqualError(t, err, test.expectedError, fmt.Sprintf("%s invalid error", test.query))
			}
			assert.EqualValues(t, test.expected, actual, fmt.Sprintf("%s unexpected value", test.query))
		})
	}
}

func printConsensusMatrix(tests []testData) {
	fmt.Println("|query|data|consensus|actual|match|")
	fmt.Println("|---|---|---|---|---|")
	for _, test := range tests {
		if test.consensus == consensusNone {
			fmt.Printf("|%s|%v|%s|%v|%s|\n", test.query, test.data, "none", test.expected, ":question:")
			continue
		}

		symbol := ":no_entry:"
		if isEqual := reflect.DeepEqual(test.consensus, test.expected); isEqual {
			symbol = ":white_check_mark:"
		}

		fmt.Printf("|%s|%v|%v|%v|%s|\n", test.query, test.data, test.consensus, test.expected, symbol)
	}
}
