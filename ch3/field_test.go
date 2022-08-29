package ch3

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetField(t *testing.T) {

	testCase := []struct {
		name      string
		obj       any
		fieldName string
		newVal    any
		wantErr   error
	}{
		{
			"nil",
			nil,
			"",
			nil,
			errors.New("nil类型非法"),
		},
		{
			"struct",
			User{
				Name: "Tom",
			},
			"Name",
			nil,
			errors.New("结构体不可修改值"),
		},
		{
			"pointer",
			&User{
				Name: "Tom",
			},
			"Name",
			"Jack",
			nil,
		},
		{
			"multiple pointer",
			func() **User {
				user := &User{Name: "Tom"}
				return &user
			}(),
			"Name",
			"Jack",
			nil,
		},
		{
			"invalid field",
			&User{
				Name: "Tom",
			},
			"invalid",
			"Jack",
			errors.New("字段不存在"),
		},
		{
			"unexport field",
			&User{
				Name: "Tom",
				age:  10,
			},
			"age",
			11,
			errors.New("字段不可修改值"), //cannot return value obtained from unexported field or method
		},
	}

	for _, c := range testCase {
		t.Run(c.name, func(t *testing.T) {
			val, err := SetField(c.obj, c.fieldName, c.newVal)
			if err != nil {
				assert.Equal(t, err, c.wantErr)
				return
			}
			assert.Equal(t, c.newVal, val)
		})
	}

}
