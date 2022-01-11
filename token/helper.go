package token

import "reflect"

func isInteger(obj interface{}) (int64, bool) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return 0, false
	}

	switch objType.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(obj).Int(), true
	}

	return 0, false
}
