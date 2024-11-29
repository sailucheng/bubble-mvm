package mvm

import "reflect"

func MakeReflectValues(args ...any) []reflect.Value {
	var values []reflect.Value
	for _, arg := range args {
		values = append(values, reflect.ValueOf(arg))
	}
	return values
}
