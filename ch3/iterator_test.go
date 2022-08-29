package ch3

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
构造 tdd 单元测试
*/

func TestTdd(t *testing.T) {
	tdd := []struct {
		name string
		// 输入
		value any
		// 输出
		wantRes map[string]any
		wantErr error
	}{
		{
			"nil",
			nil,
			nil,
			errors.New("val不能为nil"),
		},
		{
			"struct",
			User{Name: "Tom"},
			map[string]any{
				"Name": "Tom",
			},
			nil,
		},
		{
			"pointer",
			&User{Name: "Tom"},
			map[string]any{
				"Name": "Tom",
			},
			nil,
		},
		{
			"multiple pointer",
			func() **User {
				user := &User{Name: "Tom"}
				return &user
			}(),
			map[string]any{
				"Name": "Tom",
			},
			nil,
		},
		{
			// 私有字段不可读
			"unexport field",
			&User{Name: "Tom", age: 10},
			map[string]any{
				"Name": "Tom",
			},
			nil,
		},
	}

	for _, td := range tdd {
		t.Run(td.name, func(t *testing.T) {
			res, err := iterateFields(td.value)
			assert.Equal(t, td.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, td.wantRes, res)
		})
	}
}
