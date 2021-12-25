package jsonpath

import (
	"reflect"

	"github.com/evilmokeyinc/jsonpath/errors"
)

func getKey(obj interface{}, key string) (interface{}, error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, errors.ErrGetKeyFromNilMap
	}
	switch objType.Kind() {
	case reflect.Map:
		keys := reflect.ValueOf(obj).MapKeys()
		for _, kv := range keys {
			if kv.String() == key {
				return reflect.ValueOf(obj).MapIndex(kv).Interface(), nil
			}
		}
		return nil, errors.GetKeyNotFoundError(key)
	default:
		return nil, errors.ErrInvalidObjectMap
	}
}

func getIndex(obj interface{}, idx int) (interface{}, error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, errors.ErrGetIndexFromNilSlice
	}
	switch objType.Kind() {
	case reflect.Slice:
		objVal := reflect.ValueOf(obj)
		length := objVal.Len()
		if idx >= 0 {
			if idx >= length {
				return nil, errors.ErrIndexOutOfRange
			}
		} else {
			idx = length + idx
			if idx < 0 {
				return nil, errors.ErrIndexOutOfRange
			}
		}
		return objVal.Index(idx).Interface(), nil
	default:
		return nil, errors.ErrInvalidObjectSlice
	}
}

func getRange(obj interface{}, start, end, step *int) (interface{}, error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, errors.ErrGetRangeFromNilSlice
	}

	switch objType.Kind() {
	case reflect.Slice:
		objValue := reflect.ValueOf(obj)
		length := objValue.Len()

		from := 0
		if start != nil {
			from = *start
			if from < 0 {
				from = length + from
			}

			if from < 0 || from >= length {
				return nil, errors.ErrIndexOutOfRange
			}
		}

		to := length
		if end != nil {
			to = *end
			if to < 0 {
				to = length + to
			}

			if to < 0 || to >= length {
				return nil, errors.ErrIndexOutOfRange
			}
		}

		stp := 1
		if step != nil {
			stp = *step
			if stp < 1 {
				return nil, errors.GetInvalidParameterError("step should be greater than 1")
			}
		}
		sliceValue := objValue.Slice(from, to)

		if stp == 1 {
			return sliceValue.Interface(), nil
		}

		newLength := sliceValue.Len()
		stepedSlice := make([]interface{}, 0)

		for i := 0; i < newLength; i += stp {
			obj := sliceValue.Index(i).Interface()
			stepedSlice = append(stepedSlice, obj)
		}
		return stepedSlice, nil
	default:
		return nil, errors.ErrInvalidObjectSlice
	}
}
