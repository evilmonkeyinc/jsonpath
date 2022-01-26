package standard

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/evilmonkeyinc/jsonpath/script"
	"github.com/evilmonkeyinc/jsonpath/token"
)

type regexOperator struct {
	arg1, arg2 interface{}
}

func (op *regexOperator) Evaluate(parameters map[string]interface{}) (interface{}, error) {
	b, err := getString(op.arg1, parameters)
	if err != nil {
		return nil, err
	}

	if len(b) > 1 && strings.HasPrefix(b, "'") && strings.HasSuffix(b, "'") {
		b = b[1 : len(b)-1]
	}

	pattern, err := getString(op.arg2, parameters)
	if err != nil {
		return nil, err
	}

	if len(pattern) > 1 && strings.HasPrefix(pattern, "'") && strings.HasSuffix(pattern, "'") {
		pattern = pattern[1 : len(pattern)-1]
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errInvalidArgumentExpectedRegex
	}

	return regex.Match([]byte(b)), nil
}

func newSelectorOperator(selector string, engine script.Engine, options *option.QueryOptions) (*selectorOperator, error) {
	tokens := make([]token.Token, 0)

	split, err := token.Tokenize(selector)
	if err != nil {
		return nil, err
	}

	for _, str := range split {
		token, err := token.Parse(str, engine, options)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return &selectorOperator{
		selector: selector,
		tokens:   tokens,
	}, nil
}

type selectorOperator struct {
	selector string
	tokens   []token.Token
}

func (op *selectorOperator) Evaluate(parameters map[string]interface{}) (interface{}, error) {
	root := parameters["$"]
	current := parameters["@"]

	var next []token.Token
	if len(op.tokens) > 1 {
		next = op.tokens[1:]
	}

	value, err := op.tokens[0].Apply(root, current, next)
	if err != nil {
		return nil, err
	}
	if strValue, ok := value.(string); ok {
		if len(strValue) > 1 && strings.HasPrefix(strValue, "'") && strings.HasSuffix(strValue, "'") {
			return strValue, nil
		}
		return fmt.Sprintf("'%s'", strValue), nil
	}
	return value, nil
}
