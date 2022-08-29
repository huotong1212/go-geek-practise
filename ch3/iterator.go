package ch3

import (
	"errors"
	"fmt"
	"reflect"
)

// 测试方法 该方法打印所传入val中的所有字段
func IterateFields(val any) {
	ret, err := iterateFields(val)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range ret {
		fmt.Println(k, v)
	}
}

func iterateFields(val any) (map[string]any, error) {
	if val == nil {
		return nil, errors.New("val不能为nil")
	}

	refTyp := reflect.TypeOf(val)
	refVal := reflect.ValueOf(val)

	// 如果是指针
	for refTyp.Kind() == reflect.Ptr {
		// 指针type转type
		refTyp = refTyp.Elem()
		// 指针Value转value
		refVal = refVal.Elem()
	}

	data := make(map[string]any, 16)
	numField := refVal.NumField()
	for i := 0; i < numField; i++ {
		fieldTyp := refTyp.Field(i)
		fieldVal := refVal.Field(i)
		fieldName := fieldTyp.Name
		if fieldTyp.IsExported() {
			data[fieldName] = fieldVal.Interface()
		} else {
			data[fieldName] = reflect.Zero(fieldTyp.Type).Interface()
		}
	}

	return data, nil
}
