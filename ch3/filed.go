package ch3

import (
	"errors"
	"reflect"
)

func SetField(obj any, fName string, newVal any) (any, error) {

	if obj == nil {
		return nil, errors.New("nil类型非法")
	}

	refTyp := reflect.TypeOf(obj)
	if refTyp.Kind() == reflect.Struct {
		return nil, errors.New("结构体不可修改值")
	}

	refVal := reflect.ValueOf(obj)
	// 如果是指针或者是多级指针
	for refTyp.Kind() == reflect.Ptr {
		refTyp = refTyp.Elem()
		refVal = refVal.Elem()
	}
	_, ok := refTyp.FieldByName(fName)
	if !ok {
		return nil, errors.New("字段不存在")
	}
	field := refVal.FieldByName(fName)
	//fieldTyp := reflect.TypeOf(field)
	//fieldTyp.F
	if field.CanSet() {
		// It panics if CanSet returns false.
		field.Set(reflect.ValueOf(newVal))
		return field.Interface(), nil
	} else {
		return nil, errors.New("字段不可修改值")
	}
}
