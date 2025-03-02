package rest

import (
	"net/url"
	"reflect"
	"strconv"
)

func decodeQueryParams(params url.Values, dest interface{}, tagName string) error {
	val := reflect.ValueOf(dest).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get(tagName)
		if tag == "" {
			tag = fieldType.Name
		}

		if param, exists := params[tag]; exists {
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intVal, err := strconv.ParseInt(param[0], 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(intVal)
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(param[0])
				if err != nil {
					return err
				}
				field.SetBool(boolVal)
			case reflect.String:
				field.SetString(param[0])
			case reflect.Float32, reflect.Float64:
				floatVal, err := strconv.ParseFloat(param[0], 64)
				if err != nil {
					return err
				}
				field.SetFloat(floatVal)
			case reflect.Slice:
				sliceType := fieldType.Type.Elem()
				slice := reflect.MakeSlice(fieldType.Type, len(param), len(param))
				for j, v := range param {
					elem := reflect.New(sliceType).Elem()
					switch sliceType.Kind() {
					case reflect.Int:
						intVal, err := strconv.Atoi(v)
						if err != nil {
							return err
						}
						elem.SetInt(int64(intVal))
					case reflect.String:
						elem.SetString(v)
					}
					slice.Index(j).Set(elem)
				}
				field.Set(slice)
			}
		}
	}
	return nil
}
