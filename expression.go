package jsonpath

import (
	"fmt"
	"strings"
)

func newExpression(raw string) (*expression, error) {
	script := new(expression)

	err := script.parse(raw)
	if err != nil {
		return nil, err
	}

	return script, nil
}

// TODO: should scripts and filters be handled together?
// TODO : if so rename this as expression or script expression
type expression struct {
	tokens []string
}

// TODO : parsing logic can be shared but execution should be different
func (expression *expression) parse(raw string) error {
	expression.tokens = make([]string, 0)

	if strings.HasPrefix(raw, "?") {
		// TODO : this is a filter expression
		// This changes the behaviour of @ symbols
		// The result of a filter would be an array
	}
	// the result of an expression is a key, index, or union

	if !strings.HasPrefix(raw, "(") || !strings.HasSuffix(raw, ")") {
		return fmt.Errorf("not a valid expression format") // TODO : make var in errors
	}

	//TODO
	/**
	remove prefix and suffix
	split on operators, full stops
	look for quoted text '', numbers, regex //
	**/

	return nil
}
