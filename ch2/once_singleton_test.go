package ch2

import (
	"fmt"
	"sync"
)

type singleton struct {
}

var instance *singleton
var once sync.Once

func GetSingleInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

type initSource struct {
}

func (s *initSource) init() {
	once.Do(func() {
		fmt.Println("只会执行一次。。。")
	})
}
