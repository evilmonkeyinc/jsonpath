// the test package has tests detailed by https://cburgmer.github.io/json-path-comparison/
// as showing the community consensus on the expected response
package test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/evilmonkeyinc/jsonpath"
	"github.com/stretchr/testify/assert"
)

const consensusNone string = "none"

type testData struct {
	selector      string
	data          string
	expected      interface{}
	consensus     interface{}
	expectedError string
}

func batchTest(t *testing.T, tests []testData) {
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := jsonpath.QueryString(test.selector, test.data)
			if test.expectedError == "" {
				assert.Nil(t, err, fmt.Sprintf("%s error should be nil", test.selector))
			} else {
				assert.EqualError(t, err, test.expectedError, fmt.Sprintf("%s invalid error", test.selector))
			}

			if expectedArray, ok := test.expected.([]interface{}); ok {
				actualArray, ok := actual.([]interface{})
				assert.True(t, ok, "expected array response")
				if ok {
					assert.ElementsMatch(t, expectedArray, actualArray, fmt.Sprintf("%s unexpected value", test.selector))
				}
			} else {
				assert.EqualValues(t, test.expected, actual, fmt.Sprintf("%s unexpected value", test.selector))
			}
		})
	}
}

func batchBenchmark(b *testing.B, tests []testData) {
	for idx, test := range tests {
		b.Run(fmt.Sprintf("%d", idx), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				jsonpath.QueryString(test.selector, test.data)
			}
		})
	}
}

func Test_generateReadme(t *testing.T) {

	header := `# Consensus Tests

This test package has tests detailed by https://cburgmer.github.io/json-path-comparison/ comparison matrix which details the community consensus on the expected response from multiple JSONPath queries by various implementations.
	
This implementation would be closer to the 'Scalar consensus' as it does not always return an array, but instead can return a single item when that is expected.
	`

	type section struct {
		title    string
		testData []testData
	}

	sections := []section{
		{
			title:    "Array Test",
			testData: arrayTests,
		},
		{
			title:    "Bracket Test",
			testData: bracketTests,
		},
		{
			title:    "Dot Test",
			testData: dotTests,
		},
		{
			title:    "Filter Test",
			testData: filterTests,
		},
		{
			title:    "Misc Test",
			testData: miscTests,
		},
		{
			title:    "Union Test",
			testData: unionTests,
		},
	}

	file, err := os.OpenFile("README.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		t.FailNow()
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "%s\n", header)

	for _, section := range sections {
		fmt.Fprintf(writer, "\n## %s\n\n", section.title)
		printConsensusMatrix(writer, section.testData)
	}
}

func printConsensusMatrix(writer io.Writer, tests []testData) {
	fmt.Fprintf(writer, "|match|selector|data|consensus|actual|\n")
	fmt.Fprintf(writer, "|---|---|---|---|---|\n")
	for _, test := range tests {

		selector := test.selector
		// escape | so format doesnt break
		selector = strings.ReplaceAll(selector, "|", "\\|")

		expected := test.expected
		if expected == nil {
			expected = "null"
		} else {
			bytes, _ := json.Marshal(expected)
			expected = string(bytes)
		}

		if test.consensus == consensusNone {
			fmt.Fprintf(writer, "|%s|`%s`|`%v`|%s|`%v`|\n", ":question:", selector, test.data, "none", expected)
			continue
		}

		consensus := test.consensus
		if consensus == nil {
			consensus = "nil"
		} else if consensus != consensusNone {
			bytes, _ := json.Marshal(consensus)
			consensus = string(bytes)
		}

		symbol := ":no_entry:"
		if isEqual := reflect.DeepEqual(test.consensus, test.expected); isEqual {
			symbol = ":white_check_mark:"
		}

		fmt.Fprintf(writer, "|%s|`%s`|`%v`|`%v`|`%v`|\n", symbol, selector, test.data, consensus, expected)
	}
}
