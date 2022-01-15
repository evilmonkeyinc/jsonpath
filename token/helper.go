package token

import (
	"reflect"
	"strings"
)

func isInteger(obj interface{}) (int64, bool) {
	objType, objVal := getTypeAndValue(obj)
	if objType == nil {
		return 0, false
	}

	switch objType.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return objVal.Int(), true
	}

	return 0, false
}

func getTypeAndValue(obj interface{}) (reflect.Type, reflect.Value) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, reflect.ValueOf(nil)
	}

	objVal := reflect.ValueOf(obj)

	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
		objVal = objVal.Elem()
	}

	return objType, objVal
}

func getStructFields(obj reflect.Value) map[string]reflect.StructField {
	objType := obj.Type()
	if objType.Kind() != reflect.Struct {
		return nil
	}

	fields := make(map[string]reflect.StructField)

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldName := field.Name

		switch jsonTag := field.Tag.Get("json"); jsonTag {
		case "-":
		case "":
			if _, exists := fields[fieldName]; !exists {
				// Do not want to override one set with json tag
				fields[fieldName] = field
			}
		default:
			parts := strings.Split(jsonTag, ",")
			name := parts[0]
			if name == "" {
				name = fieldName
			}
			fields[name] = field
		}
	}

	return fields
}
