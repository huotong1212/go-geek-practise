package ch1

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContextTimeOut(t *testing.T) {
	// 定义context的起点
	root := context.Background() // empty ctx
	// 当我们不知道context从哪里来，但是将来会有一个context(不是起点，却需要一个context，就使用TODO)  很少的情况才会使用
	//todoCtx := context.TODO()
	timeoutCtx, cancel := context.WithTimeout(root, time.Millisecond) // timerCtx
	defer cancel()                                                    // cancel 大部分情况下cancel会在方法中使用,因为如果手动调用的话，运行出现了panic或者return就会调用不成功，而defer可以保证一定会执行。

	time.Sleep(time.Second)
	err := timeoutCtx.Err()
	fmt.Println(err)

	switch err {
	case context.DeadlineExceeded:
		fmt.Println("context time out!")
	case context.Canceled:
		fmt.Println("context has been canceled")
	}
	/** console
	context deadline exceeded
	context time out!
	*/
}

func TestContextCancel(t *testing.T) {
	// 定义context的起点
	root := context.Background() // empty ctx
	// 当我们不知道context从哪里来，但是将来会有一个context(不是起点，却需要一个context，就使用TODO)  很少的情况才会使用
	//todoCtx := context.TODO()
	timeoutCtx, cancel := context.WithTimeout(root, time.Millisecond) // timerCtx
	cancel()                                                          // cancel 大部分情况下cancel会在方法中使用,因为如果手动调用的话，运行出现了panic或者return就会调用不成功，而defer可以保证一定会执行。

	err := timeoutCtx.Err()
	fmt.Println(err)

	switch err {
	case context.DeadlineExceeded:
		fmt.Println("context time out!")
	case context.Canceled:
		fmt.Println("context has been canceled")
	}
	/** console
	context canceled
	context has been canceled
	*/
}

func TestContextDeadline(t *testing.T) {
	root := context.Background()                                 // empty ctx
	timeoutCtx, cancel := context.WithTimeout(root, time.Second) // timerCtx
	defer cancel()

	dl, ok := timeoutCtx.Deadline()
	fmt.Println("dl:", dl, "|ok:", ok) // 如果设置了截止时间， dl表示截止时间 ok为true

	dl, ok = root.Deadline()
	fmt.Println("dl:", dl, "|ok:", ok) // 如果未设置截止时间， dl表示公元时间的起点 ok为false
}

func TestContextWithValue(t *testing.T) {
	root := context.Background()                         // empty ctx
	valCtx := context.WithValue(root, "clearlove", "厂长") // 这里会返回一个新的context
	val := valCtx.Value("clearlove")                     // Value返回的类型是 interface{} 在go1.8中别名是是any
	fmt.Println(val)

	/** console
	厂长
	*/
}
