package ch3

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type FieldAccessor interface {
	Field(field string) (int, error)
	SetField(field string, val any) error
}

type UnsafeAccessor struct {
	fields     map[string]FieldMeta
	entityAddr unsafe.Pointer
}

func NewUnsafeAccessor(entity interface{}) (*UnsafeAccessor, error) {
	if entity == nil {
		return nil, errors.New("nil非法输入")
	}

	refTyp := reflect.TypeOf(entity)
	// 只能是一级指针 ,否则 refVal.UnsafeAddr() 值传递无法寻址
	if refTyp.Kind() != reflect.Pointer || refTyp.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid entity")
	}
	refVal := reflect.ValueOf(entity)
	if refTyp.Kind() == reflect.Ptr {
		refTyp = refTyp.Elem()
		refVal = refVal.Elem()
	}

	numField := refTyp.NumField()
	fields := make(map[string]FieldMeta, numField)

	for i := 0; i < numField; i++ {
		field := refTyp.Field(i)
		fields[field.Name] = FieldMeta{
			offset: field.Offset,
			typ:    field.Type,
		}
	}

	accessor := &UnsafeAccessor{
		fields:     fields,
		entityAddr: unsafe.Pointer(refVal.UnsafeAddr()),
	}

	return accessor, nil
}

func (u *UnsafeAccessor) Field(field string) (any, error) {
	fdMeta, ok := u.fields[field]
	// 判断字段是否存在
	if !ok {
		return 0, errors.New("字段不存在")
	}

	// 将内存中的字节转换为类型值  字节=(起始位置+字段偏移量)
	ptr := unsafe.Pointer(uintptr(u.entityAddr) + fdMeta.offset)
	if ptr == nil {
		return 0, fmt.Errorf("invalid address of the field: %s", field)
	}
	// 如果提前知道了类型，可以这样设置值
	//fdv := *(*int)(ptr)
	// 如果不知道类型，可以用reflect.Type，但是效率没有上面的快  取出来时是个指针，所以还需要再Elem
	fdv := reflect.NewAt(fdMeta.typ, ptr).Elem().Interface()
	return fdv, nil
}

func (u *UnsafeAccessor) SetField(field string, val int) error {
	fdMeta, ok := u.fields[field]
	// 判断字段是否存在
	if !ok {
		return errors.New("字段不存在")
	}

	// 将内存中的字节转换为类型值  字节=(起始位置+字段偏移量)
	ptr := unsafe.Pointer(uintptr(u.entityAddr) + fdMeta.offset)
	if ptr == nil {
		return fmt.Errorf("invalid address of the field: %s", field)
	}

	// 如果已知类型，可以这样操作
	//*(*int)(ptr) = val
	// 如果未知，这样操作
	fdv := reflect.NewAt(fdMeta.typ, ptr)
	if fdv.CanSet() {
		fdv.Set(reflect.ValueOf(val))
	}
	return nil
}

type FieldMeta struct {
	// offset 后期在我们考虑组合，或者复杂类型字段的时候，它的含义衍生为表达相当于最外层的结构体的偏移量
	offset uintptr
	typ    reflect.Type
}
