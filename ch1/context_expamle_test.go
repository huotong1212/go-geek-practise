package ch1

import (
	"context"
	"fmt"
	"net/http"
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

func doSomeBusiness() {
	fmt.Println("some business begin")
	time.Sleep(time.Second * 10)
}

// context 超时控制
func TestContextTimeout(t *testing.T) {
	ctx := context.Background()
	tmCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	end := make(chan int)

	go func() {
		doSomeBusiness()
		end <- 1
	}()

	select {
	case <-tmCtx.Done():
		// 这里是通过主线程结束来控制 goroutine结束的
		fmt.Println(tmCtx.Err())
	case <-end:
		fmt.Println("business over")
	}

	//time.Sleep(time.Second * 10)
	//
	//select {
	//case <-tmCtx.Done():
	//	// 这里是通过主线程结束来控制 goroutine结束的
	//	fmt.Println(tmCtx.Err())
	//case <-end:
	//	fmt.Println("business over")
	//}
}

// time.After 超时控制
func TestTimeAfter(t *testing.T) {
	ch := make(chan int)
	go func() {
		doSomeBusiness()
		ch <- 1
	}()

	select {
	case <-ch:
		fmt.Println("business over")
	case <-time.After(time.Second):
		fmt.Println("out of time")
	}
}

// time.AfterFunc 超时控制
func TestTimeAfterFunc(t *testing.T) {
	http.Request{}
	ch := make(chan int)
	go func() {
		doSomeBusiness()
		ch <- 1
	}()

	timer := time.AfterFunc(time.Second, func() {
		fmt.Println("time after func is bound to execute")
		ch <- 1
	})

	<-ch
	timer.Stop()
	fmt.Println("business end")
}
