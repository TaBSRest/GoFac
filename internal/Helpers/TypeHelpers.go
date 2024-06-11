package Helpers

import "reflect"

func IsArrayOrSlice(info reflect.Type) bool {
	return info.Kind() == reflect.Slice || info.Kind() == reflect.Array
}

func IsValueArrayOrSlice(value reflect.Value) bool {
	return value.Kind() == reflect.Slice || value.Kind() == reflect.Array
}
