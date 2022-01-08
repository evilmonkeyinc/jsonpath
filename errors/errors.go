package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrGetElementsFromNilObject error returned when requesting elements from nil object
	ErrGetElementsFromNilObject error = fmt.Errorf("cannot get elements from nil object")
	// ErrGetIndexFromNilArray error returned when requesting index from nil array
	ErrGetIndexFromNilArray error = fmt.Errorf("cannot get index from nil array")
	// ErrGetKeyFromNilMap error returned when requesting key from nil map
	ErrGetKeyFromNilMap error = fmt.Errorf("cannot get key from nil map")
	// ErrGetRangeFromNilArray error returned when requesting range from nil array
	ErrGetRangeFromNilArray error = fmt.Errorf("cannot get range from nil array")
	// ErrGetUnionFromNilObject error returned when requesting values from nil object
	ErrGetUnionFromNilObject error = fmt.Errorf("cannot get union from nil object")
	// ErrIllegalCharacterAtPositionOne error returned during parsing if '.' or '[' doesn't follow the initial token.
	ErrIllegalCharacterAtPositionOne = fmt.Errorf("expected '.' or '[' after initial token")
	// ErrIndexOutOfRange returned when an invalid index is requested of an array
	ErrIndexOutOfRange error = fmt.Errorf("index out of range")
	// ErrInvalidInitialToken error returned during parsing if an invalid initial token is found.
	ErrInvalidInitialToken error = fmt.Errorf("query must start with '$'")
	// ErrInvalidObject error returned when an invalid or unexpected object is found
	ErrInvalidObject error = fmt.Errorf("invalid object")
	// ErrInvalidParameterInteger error returned when an invalid parameter is passed to a function when an integer is expected.
	ErrInvalidParameterInteger error = fmt.Errorf("%w. expected integer", ErrInvalidParameter)
	// ErrInvalidObjectMap error returned when an invalid object is found when expecting a map
	ErrInvalidObjectMap error = fmt.Errorf("%w. expected map", ErrInvalidObject)
	// ErrInvalidObjectArray error returned when an invalid object is found when expecting an array
	ErrInvalidObjectArray error = fmt.Errorf("%w. expected array", ErrInvalidObject)
	// ErrInvalidObjectArrayOrMap error returned when an invalid object is found when expecting an array or map
	ErrInvalidObjectArrayOrMap error = fmt.Errorf("%w. expected array or map", ErrInvalidObject)
	// ErrInvalidObjectArrayMapOrString error returned when an invalid object is found when expecting an array, map, or string
	ErrInvalidObjectArrayMapOrString error = fmt.Errorf("%w. expected array, map, or string", ErrInvalidObject)
	// ErrInvalidParameter error returned when an invalid parameter is passed to a function
	ErrInvalidParameter error = fmt.Errorf("invalid parameter")
	// ErrInvalidParameterRangeNegativeStep error returned with the range token is given a step value that is less than one
	ErrInvalidParameterRangeNegativeStep error = fmt.Errorf("%w. step should be greater than or equal to 1", ErrInvalidParameter)
	// ErrInvalidParameterUnionExpectedString error returned when union token is expecting string keys but did not
	ErrInvalidParameterUnionExpectedString error = fmt.Errorf("%w. expected string keys", ErrInvalidParameter)
	// ErrInvalidParameterUnionExpectedInteger error returned when union token is expecting integer keys but did not
	ErrInvalidParameterUnionExpectedInteger error = fmt.Errorf("%w. expected integer keys", ErrInvalidParameter)
	// ErrInvalidParameterScriptExpectedToReturnString error returned when expected to get a string value from an expression but did not
	ErrInvalidParameterScriptExpectedToReturnString error = fmt.Errorf("%w. expected script to return string", ErrInvalidParameter)
	// ErrInvalidParameterScriptExpectedToReturnInteger error returned when expected to get an integer value from an expression but did not
	ErrInvalidParameterScriptExpectedToReturnInteger error = fmt.Errorf("%w. expected script to return integer", ErrInvalidParameter)
	// ErrInvalidParameterExpressionEmpty error returned when trying to evaluate empty expression
	ErrInvalidParameterExpressionEmpty error = fmt.Errorf("%w. expression is empty", ErrInvalidParameter)
	// ErrInvalidToken error returned while parsing a token which is invalid
	ErrInvalidToken error = fmt.Errorf("invalid token")
	// ErrInvalidTokenEmpty error returned when an empty token is parsed
	ErrInvalidTokenEmpty error = fmt.Errorf("%w. token can not be empty", ErrInvalidToken)
	// ErrInvalidTokenNoRangeInUnion error returned when attempting to combine unions and range subscripts
	ErrInvalidTokenNoRangeInUnion error = fmt.Errorf("%w. cannot specify a range in a union", ErrInvalidToken)
	// ErrInvalidTokenInvalidRangeArguments error returned when using invalid arguments with a range subscript token
	ErrInvalidTokenInvalidRangeArguments error = fmt.Errorf("%w. only integer or scripts allowed in range arguments", ErrInvalidToken)
	// ErrInvalidTokenIncorrectNumberOfRangeArguments error returned when using the wrong number of arguments with a range subscript token
	ErrInvalidTokenIncorrectNumberOfRangeArguments error = fmt.Errorf("%w. incorrect number of arguments in range", ErrInvalidToken)
	// ErrInvalidTokenUnexpectedUnionArguments error returned when an unexpected argument used with a union subscript token
	ErrInvalidTokenUnexpectedUnionArguments error = fmt.Errorf("%w. unexpected union argument", ErrInvalidToken)
	// ErrInvalidTokenEmptyUnionArguments error returned when an empty argument used with a union subscript token
	ErrInvalidTokenEmptyUnionArguments error = fmt.Errorf("%w. empty argument in union", ErrInvalidToken)
	// ErrInvalidTokenUnexpectedIndex error returned when an integer is specified as a map key
	ErrInvalidTokenUnexpectedIndex error = fmt.Errorf("%w. index specified as key", ErrInvalidToken)
	// ErrInvalidTokenMissingSubscriptClose error returned when missing a subscript close bracket
	ErrInvalidTokenMissingSubscriptClose error = fmt.Errorf("%w. missing subscript close", ErrInvalidToken)
	// ErrInvalidTokenEmptySubscript error returned when a empty subscript is parsed
	ErrInvalidTokenEmptySubscript error = fmt.Errorf("%w. empty subscript", ErrInvalidToken)
	// ErrInvalidTokenInvalidFilterFormat error returned when an invalid filter is parsed
	ErrInvalidTokenInvalidFilterFormat error = fmt.Errorf("%w. invalid filter format", ErrInvalidToken)
	// ErrInvalidTokenInvalidScriptFormat error returned when an invalid script is parsed
	ErrInvalidTokenInvalidScriptFormat error = fmt.Errorf("%w. invalid script format", ErrInvalidToken)
	// ErrInvalidTokenInvalidKeyFormat error returned when an invalid key format is parsed
	ErrInvalidTokenInvalidKeyFormat error = fmt.Errorf("%w. invalid key format", ErrInvalidToken)
	// ErrInvalidTokenUnexpectedString error returned when an unexpected string is parsed
	ErrInvalidTokenUnexpectedString error = fmt.Errorf("%w. unexpected string", ErrInvalidToken)
	// ErrInvalidTokenUnexpectedSpace error returned when an unexpected space is parsed
	ErrInvalidTokenUnexpectedSpace error = fmt.Errorf("%w. unexpected space", ErrInvalidToken)
	// ErrInvalidTokenUnexpectedQuote error returned when an unexpected single quote is parsed
	ErrInvalidTokenUnexpectedQuote error = fmt.Errorf("%w. unexpected single quote", ErrInvalidToken)
	// ErrInvalidTokenInvalidIndex error returned when an invalid index is parsed
	ErrInvalidTokenInvalidIndex error = fmt.Errorf("%w. invalid index", ErrInvalidToken)
	// ErrInvalidQuery error returned when parsing an invalid query string
	ErrInvalidQuery error = fmt.Errorf("invalid query")
	// ErrInvalidQueryNoTokens error returned when no tokens are found when parsing a query string
	ErrInvalidQueryNoTokens error = fmt.Errorf("%w. no tokens", ErrInvalidQuery)
	// ErrInvalidQueryUnexpectedTokens error returned when parsing a query string and finding unexpected tokens
	ErrInvalidQueryUnexpectedTokens error = fmt.Errorf("%w. unexpected tokens", ErrInvalidQuery)
	// ErrQueryNotSpecified error returned when a query is required or parsed and it has not be specified or is empty
	ErrQueryNotSpecified error = fmt.Errorf("no valid JSONPath query has been specified")
	// ErrUnexpectedScriptResultInteger error returned when receive unexpected result from script when expecting an integer
	ErrUnexpectedScriptResultInteger error = fmt.Errorf("%w. expected integer", errUnexpectedScriptResult)

	errFailedToParseExpression error = fmt.Errorf("failed to parse expression")
	errInvalidExpression       error = fmt.Errorf("invalid expression")
	errJSONMarshalFailed       error = fmt.Errorf("json marshal failed")
	errKeyNotFound             error = fmt.Errorf("key not found in object")
	errUnexpectedScriptResult  error = fmt.Errorf("unexpected script result")
)

// IsKeyNotFoundError returns true if the specified error includes a key not found error
func IsKeyNotFoundError(err error) bool {
	return errors.Is(err, errKeyNotFound)
}

// GetKeyNotFoundError retruns a key not found error
func GetKeyNotFoundError(key string) error {
	return fmt.Errorf("'%s' %w", key, errKeyNotFound)
}

// IsFailedToParseExpressionError returns true if specified error is failed to parse expression error
func IsFailedToParseExpressionError(err error) bool {
	return errors.Is(err, errFailedToParseExpression)
}

// GetFailedToParseExpressionError returns a failed to parse expression error
func GetFailedToParseExpressionError(reason error) error {
	return fmt.Errorf("%w. %s", errFailedToParseExpression, reason.Error())
}

// IsJSONMarshalFailedError returns true if specified error is JSON marshaling error
func IsJSONMarshalFailedError(err error) bool {
	return errors.Is(err, errJSONMarshalFailed)
}

// GetJSONMarshalFailedError returns a JSON marshal failed error
func GetJSONMarshalFailedError(reason error) error {
	return fmt.Errorf("%w. %s", errJSONMarshalFailed, reason.Error())
}

// IsInvalidExpressionError returns true if specified error is an invalid expression error
func IsInvalidExpressionError(err error) bool {
	return errors.Is(err, errInvalidExpression)
}

// GetInvalidExpressionError returns an invalid expression error
func GetInvalidExpressionError(reason error) error {
	return fmt.Errorf("%w. %s", errInvalidExpression, reason.Error())
}
