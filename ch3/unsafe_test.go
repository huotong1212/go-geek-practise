package ch3

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnsafeAccessor_Field(t *testing.T) {
	var tests = []struct {
		name    string
		input   any
		field   string
		wantVal any
		wantErr error
	}{
		{
			"nil",
			nil,
			"",
			nil,
			errors.New("nil非法输入"),
		},
		{
			"ptr",
			&User{"Tom", 18},
			"age",
			any(18),
			nil,
		},
		{ // Failed panic: reflect.Value.UnsafeAddr of unaddressable value
			"struct",
			User{"Tom", 18},
			"age",
			18,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			accessor, err := NewUnsafeAccessor(test.input)
			if err != nil {
				assert.Equal(t, test.wantErr, err)
				return
			}
			val, err := accessor.Field(test.field)
			if err != nil {
				assert.Equal(t, test.wantErr, err)
				return
			}
			//anyval := any(val)
			assert.Equal(t, test.wantVal, val)
		})
	}
}
