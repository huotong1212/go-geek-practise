package ch3

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	ref := reflect.TypeOf(&Student{Name: "Tom"})
	//ref := reflect.TypeOf(111) //
	if ref.Kind() == reflect.Struct {
		fmt.Println("结构体")
	} else if ref.Kind() == reflect.Ptr {
		fmt.Println("指针")
		// 通过指针取值
		ref = ref.Elem()
	} else {
		fmt.Println(ref.Kind())
		return
	}
	fmt.Println(ref.NumField()) // 该方法只能对结构体使用
}

type Student struct {
	Name string
}
