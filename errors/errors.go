package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidInitialToken error returned during parsing if an invalid initial token is found.
	ErrInvalidInitialToken error = fmt.Errorf("query must start with '$'")
	// ErrIllegalCharacterAtPositionOne error returned during parsing if '.' or '[' doesn't follow the initial token.
	ErrIllegalCharacterAtPositionOne = fmt.Errorf("expected '.' or '[' after initial token")
	// ErrQueryNotSpecified error returned when a query is required or parsed and it has not be specified or is empty
	ErrQueryNotSpecified = fmt.Errorf("no valid JSONPath query has been specified")

	// ErrGetKeyFromNilMap error returned when requesting key from nil map
	ErrGetKeyFromNilMap error = fmt.Errorf("cannot get key from nil map")
	// ErrGetIndexFromNilArray error returned when requesting index from nil array
	ErrGetIndexFromNilArray error = fmt.Errorf("cannot get index from nil array")
	// ErrGetRangeFromNilArray error returned when requesting range from nil array
	ErrGetRangeFromNilArray error = fmt.Errorf("cannot get range from nil array")
	// ErrGetUnionFromNilObject error returned when requesting values from nil object
	ErrGetUnionFromNilObject error = fmt.Errorf("cannot get union from nil object")
	// ErrGetElementsFromNilObject error returned when requesting elements from nil object
	ErrGetElementsFromNilObject error = fmt.Errorf("cannot get elements from nil object")
	// ErrIndexOutOfRange returned when an invalid index is requested of an array
	ErrIndexOutOfRange error = fmt.Errorf("index out of range")
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

	ErrUnexpectedScriptResult        error = fmt.Errorf("unexpected script result")
	ErrUnexpectedScriptResultInteger error = fmt.Errorf("%w. expected integer", ErrUnexpectedScriptResult)

	errKeyNotFound error = fmt.Errorf("key not found in object")

	errFailedToParseExpression error = fmt.Errorf("failed to parse expression")

	errJSONMarshalFailed error = fmt.Errorf("json marshal failed")

	ErrInvalidQuery                 error = fmt.Errorf("invalid query")
	ErrInvalidQueryNoTokens         error = fmt.Errorf("%w. no tokens", ErrInvalidQuery)
	ErrInvalidQueryUnexpectedTokens error = fmt.Errorf("%w. unexpected tokens", ErrInvalidQuery)

	ErrInvalidExpression error = fmt.Errorf("invalid expression")

	// ErrInvalidToken err returned while parsing a token which is invalid
	ErrInvalidToken                                      = fmt.Errorf("invalid token")
	ErrInvalidTokenEmpty                           error = fmt.Errorf("%w. token can not be empty", ErrInvalidToken)
	ErrInvalidTokenNoRangeInUnion                  error = fmt.Errorf("%w. cannot specify a range in a union", ErrInvalidToken)
	ErrInvalidTokenInvalidRangeArguments           error = fmt.Errorf("%w. only integer or scripts allowed in range arguments", ErrInvalidToken)
	ErrInvalidTokenIncorrectNumberOfRangeArguments error = fmt.Errorf("%w. incorrect number of arguments in range", ErrInvalidToken)
	ErrInvalidTokenUnexpectedUnionArguments        error = fmt.Errorf("%w. unexpected union argument", ErrInvalidToken)
	ErrInvalidTokenEmptyUnionArguments             error = fmt.Errorf("%w. empty argument in union", ErrInvalidToken)
	ErrInvalidTokenUnexpectedIndex                 error = fmt.Errorf("%w. index specified as key", ErrInvalidToken)

	ErrInvalidTokenMissingSubscriptClose error = fmt.Errorf("%w. missing subscript close", ErrInvalidToken)
	ErrInvalidTokenEmptySubscript        error = fmt.Errorf("%w. empty subscript", ErrInvalidToken)
	ErrInvalidTokenInvalidFilterFormat   error = fmt.Errorf("%w. invalid filter format", ErrInvalidToken)
	ErrInvalidTokenInvalidScriptFormat   error = fmt.Errorf("%w. invalid script format", ErrInvalidToken)
	ErrInvalidTokenInvalidKeyFormat      error = fmt.Errorf("%w. invalid key format", ErrInvalidToken)
	ErrInvalidTokenUnexpectedString      error = fmt.Errorf("%w. unexpected string", ErrInvalidToken)
	ErrInvalidTokenUnexpectedSpace       error = fmt.Errorf("%w. unexpected space", ErrInvalidToken)
	ErrInvalidTokenUnexpectedQuote       error = fmt.Errorf("%w. unexpected single quote", ErrInvalidToken)
	ErrInvalidTokenInvalidIndex          error = fmt.Errorf("%w. invalid index", ErrInvalidToken)

	// ErrInvalidParameter error returned when an invalid parameter is passed to a function
	ErrInvalidParameter                              error = fmt.Errorf("invalid parameter")
	ErrInvalidParameterRangeNegativeStep             error = fmt.Errorf("%w. step should be greater than 1", ErrInvalidParameter)
	ErrInvalidParameterUnionExpectedString           error = fmt.Errorf("%w. expected string keys", ErrInvalidParameter)
	ErrInvalidParameterUnionExpectedInteger          error = fmt.Errorf("%w. expected integer keys", ErrInvalidParameter)
	ErrInvalidParameterScriptExpectedToReturnString  error = fmt.Errorf("%w. expected script to return string", ErrInvalidParameter)
	ErrInvalidParameterScriptExpectedToReturnInteger error = fmt.Errorf("%w. expected script to return integer", ErrInvalidParameter)
	ErrInvalidParameterExpressionEmpty               error = fmt.Errorf("%w. expression is empty", ErrInvalidParameter)
)

// IsKeyNotFoundError returns true if the specified error includes a key not found error
func IsKeyNotFoundError(err error) bool {
	return errors.Is(err, errKeyNotFound)
}

// GetKeyNotFoundError retruns a key not found error
func GetKeyNotFoundError(key string) error {
	return fmt.Errorf("'%s' %w", key, errKeyNotFound)
}

// GetFailedToParseExpressionError returns a failed to parse expression error.
func GetFailedToParseExpressionError(reason error) error {
	return fmt.Errorf("%w. %s", errFailedToParseExpression, reason.Error())
}

func GetJSONMarshalFailedError(reason error) error {
	return fmt.Errorf("%w. %s", errJSONMarshalFailed, reason.Error())
}

func GetInvalidExpressionError(reason error) error {
	return fmt.Errorf("%w. %s", ErrInvalidExpression, reason.Error())
}
