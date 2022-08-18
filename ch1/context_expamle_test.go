package ch1

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

/**
使用 context 控制超时
*/

func SomeBusiness(ctx context.Context) {
	fmt.Println("Do some business")
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		fmt.Println("业务关闭...")
		//return
	}
}

func TestContextDoneClose(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithCancel(ctx)
	t.Log("begin:", runtime.NumGoroutine())

	go func() {
		// 启动一个协程来跑业务
		SomeBusiness(childCtx)
	}()
	t.Log("run goroutine:", runtime.NumGoroutine())
	// 设置3秒后业务超时
	time.Sleep(time.Second * 3)
	cancel()
	t.Log("context cancel:", runtime.NumGoroutine())
}
