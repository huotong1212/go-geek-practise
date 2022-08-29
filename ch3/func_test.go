package ch3

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	var tests = []struct {
		name    string
		input   any
		wantVal map[string]*FuncInfo
		wantErr error
	}{
		{
			name:    "nil",
			input:   nil,
			wantErr: errors.New("nil输入"),
		},
		{ // OK
			name:  "struct receiver func input struct",
			input: User{age: 10},
			wantVal: map[string]*FuncInfo{
				"GetAge": &FuncInfo{
					Name:   "GetAge",
					In:     []reflect.Type{reflect.TypeOf(User{})},
					Out:    []reflect.Type{reflect.TypeOf(int(0))},
					Result: []any{int(10)},
				},
			},
		},
		{ // OK
			name:  "pointer receiver func input pointer",
			input: &UserV1{age: 10},
			wantVal: map[string]*FuncInfo{
				"GetAge": &FuncInfo{
					Name:   "GetAge",
					In:     []reflect.Type{reflect.TypeOf(&UserV1{})},
					Out:    []reflect.Type{reflect.TypeOf(int(0))},
					Result: []any{int(10)},
				},
			},
		},
		{ // Failed
			name:  "pointer receiver func input struct",
			input: UserV1{age: 10},
			wantVal: map[string]*FuncInfo{
				"GetAge": &FuncInfo{
					Name:   "GetAge",
					In:     []reflect.Type{reflect.TypeOf(UserV1{})},
					Out:    []reflect.Type{reflect.TypeOf(int(0))},
					Result: []any{int(10)},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualVal, err := IterateFunc(test.input)
			if err != nil {
				assert.Equal(t, test.wantErr, err)
			}
			assert.Equal(t, test.wantVal, actualVal)
		})
	}
}
