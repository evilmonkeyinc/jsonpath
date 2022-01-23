package standard

import (
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

	pattern, err := getString(op.arg2, parameters)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(pattern, "/") {
		end := strings.LastIndex(pattern, "/")
		pattern = pattern[1:end]
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

	return op.tokens[0].Apply(root, current, next)
}
