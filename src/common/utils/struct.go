package utils

import (
	"reflect"
	"strings"
	"time"
)

func WalkField(obj interface{}, fn func(name string, val interface{})) {
	walkField("", obj, fn)
}

func walkField(name string, obj interface{}, fn func(name string, val interface{})) {
	if obj == nil {
		return
	}
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		typ = typ.Elem()
		val = val.Elem()
	}
	switch typ.Kind() {
	case reflect.Struct:
		structWalk(typ, val, fn)
	case reflect.Map:
	case reflect.Slice:
	default:
		fn(name, val.Interface())
	}
}

func structWalk(typ reflect.Type, val reflect.Value, fn func(name string, val interface{})) {
	num := typ.NumField()
	for i := 0; i < num; i++ {
		field := typ.Field(i)
		fieldTyp := field.Type
		fieldVal := val.Field(i)
		if fieldTyp.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}
		if fieldTyp.Kind() == reflect.Struct && strings.Contains(fieldTyp.String(), "time.Time") {
			// 如果两个都为空就直接退出
			if fieldVal.IsZero() {
				fn(field.Name, "")
				continue
			}

			if strings.Contains(fieldVal.Type().String(), "time.Time") {
				tmA := fieldVal.Interface().(time.Time).Format("2006-01-02 15:04:05")
				fn(field.Name, tmA)
			}
		} else {
			walkField(field.Name, fieldVal.Interface(), fn)
		}
	}
}

func mapWalk(obj interface{}, fn func(name string, val interface{})) {

}

func sliceWalk(obj interface{}, fn func(name string, val interface{})) {
}
