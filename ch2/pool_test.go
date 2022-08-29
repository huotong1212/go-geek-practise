package ch2

import (
	"fmt"
	"sync"
	"testing"
)

func TestPoolBasicUsage(t *testing.T) {
	// New方法只有在池中没有实例时才会执行
	pool := sync.Pool{New: func() any {
		fmt.Println("only init when no instance in pool")
		return ""
	}}

	for i := 0; i < 100; i++ {
		v := pool.Get()
		pool.Put(v)
	}

}
