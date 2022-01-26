package standard

import (
	"encoding/json"
	"strings"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/evilmonkeyinc/jsonpath/script"
)

// TODO : add tests for what is in readme
// TODO : update readme to give more details, maybe add readme to this package and link from main
// TODO : add support for bitwise operators | &^ ^ &  << >> after + and -
var defaultTokens []string = []string{
	"||", "&&",
	"==", "!=", "<=", ">=", "<", ">", "=~",
	"+", "-",
	"**", "*", "/", "%",
	"@", "$",
}

// ScriptEngine standard implementation of the script engine interface
type ScriptEngine struct {
}

// Compile returns a compiled expression that can be evaluated multiple times
func (engine *ScriptEngine) Compile(expression string, options *option.QueryOptions) (script.CompiledExpression, error) {
	tokens := defaultTokens
	operator, err := engine.buildOperators(expression, tokens, options)
	if err != nil {
		return nil, err
	}

	return &compiledExpression{
		expression:   expression,
		rootOperator: operator,
		engine:       engine,
		options:      options,
	}, nil
}

// Evaluate return the result of the expression evaluation
func (engine *ScriptEngine) Evaluate(root, current interface{}, expression string, options *option.QueryOptions) (interface{}, error) {
	compiled, err := engine.Compile(expression, options)
	if err != nil {
		return nil, err
	}
	evaluation, err := compiled.Evaluate(root, current)
	if err != nil {
		return nil, err
	}
	return evaluation, nil
}

func (engine *ScriptEngine) buildOperators(expression string, tokens []string, options *option.QueryOptions) (operator, error) {
	nextToken := tokens[0]
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return nil, nil
	}
	if strings.HasPrefix(expression, "(") && strings.HasSuffix(expression, ")") {
		expression = strings.TrimSpace(expression[1 : len(expression)-1])
		// since we were in brackets, we need to try all the tokens again
		tokens = append([]string{""}, defaultTokens...)
	}

	idx := findUnquotedOperators(expression, nextToken)
	if idx < 0 {
		if len(tokens) == 1 {
			return nil, nil
		}
		// if none of these tokens, move onto next token
		return engine.buildOperators(expression, tokens[1:], options)
	}

	// check right for more tokens, or use raw string as input
	// right check needs done first as some expressions just have left sides
	arg2Str := strings.TrimSpace(expression[idx+(len(nextToken)):])
	if arg2Str == "" && (nextToken != "$" && nextToken != "@") {
		if len(tokens) == 1 {
			return nil, nil
		}
		// if none of these tokens, move onto next token
		return engine.buildOperators(expression, tokens[1:], options)
	}

	arg2, err := engine.parseArgument(arg2Str, tokens, options)
	if err != nil {
		return nil, err
	}

	// check left for more tokens, or use raw string as input
	arg1Str := strings.TrimSpace(expression[0:idx])
	if arg1Str == "" && (nextToken != "$" && nextToken != "@") {
		if len(tokens) == 1 {
			return nil, nil
		}
		// if none of these tokens, move onto next token
		return engine.buildOperators(expression, tokens[1:], options)
	}

	arg1, err := engine.parseArgument(arg1Str, tokens, options)
	if err != nil {
		return nil, err
	}

	switch nextToken {
	case "@", "$":
		selector, err := newSelectorOperator(expression, engine, options)
		if err != nil {
			return nil, err
		}
		return selector, nil
	case "=~":
		return &regexOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "||":
		return &orOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "&&":
		return &andOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "==":
		return &equalsOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "!=":
		return &notEqualsOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "<=":
		return &lessThanOrEqualOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case ">=":
		return &greaterThanOrEqualOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "<":
		return &lessThanOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case ">":
		return &greaterThanOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "+":
		return &plusOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "-":
		return &subtractOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "**":
		return &powerOfOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "*":
		return &multiplyOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "/":
		return &divideOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "%":
		return &modulusOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	}

	// will cover when we add a new operator token
	// but forget to update the switch/case
	return nil, errUnsupportedOperator
}

func (engine *ScriptEngine) parseArgument(argument string, tokens []string, options *option.QueryOptions) (interface{}, error) {
	if op, err := engine.buildOperators(argument, tokens, options); err != nil {
		return nil, err
	} else if op != nil {
		return op, nil
	}

	argument = strings.TrimSpace(argument)
	var arg interface{} = argument
	if strings.HasPrefix(argument, "[") && strings.HasSuffix(argument, "]") {
		val := make([]interface{}, 0)
		if err := json.Unmarshal([]byte(argument), &val); err == nil {
			arg = val
		}
	} else if strings.HasPrefix(argument, "{") && strings.HasSuffix(argument, "}") {
		val := make(map[string]interface{})
		if err := json.Unmarshal([]byte(argument), &val); err == nil {
			arg = val
		}
	}
	return arg, nil
}
