package ch3

import (
	"errors"
	"reflect"
)

func IterateFunc(input any) (map[string]*FuncInfo, error) {
	if input == nil {
		return nil, errors.New("nil输入")
	}

	// 判断如果不是结构体或者结构体的一级指针 -> error
	refTyp := reflect.TypeOf(input)
	if refTyp.Kind() != reflect.Struct && !(refTyp.Kind() == reflect.Ptr && refTyp.Elem().Kind() == reflect.Struct) {
		return nil, errors.New("非法输入")
	}
	//refVal := reflect.ValueOf(input)

	// 获取结构体的方法数量
	numMd := refTyp.NumMethod()
	data := make(map[string]*FuncInfo, numMd)
	for i := 0; i < numMd; i++ {
		md := refTyp.Method(i)

		numIn := md.Type.NumIn()
		in := make([]reflect.Type, 0, numIn)
		for j := 0; j < numIn; j++ {
			in = append(in, md.Type.In(j))
		}

		numOut := md.Type.NumOut()
		out := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			out = append(out, md.Type.Out(j))
		}

		result := md.Func.Call([]reflect.Value{reflect.ValueOf(input)})

		resBs := make([]any, 0, len(result))
		for n := 0; n < len(result); n++ {
			resBs = append(resBs, result[n].Interface())
		}

		funcInfo := &FuncInfo{
			Name:   md.Name,
			In:     in,
			Out:    out,
			Result: resBs,
		}

		data[md.Name] = funcInfo
	}

	return data, nil
}

type FuncInfo struct {
	Name string
	In   []reflect.Type
	Out  []reflect.Type

	// 反射调用得到的结果
	Result []any
}
