package ch1

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithCancel(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithCancel(ctx)
	cancel()
	// 父取消，子取消
	fmt.Println(childCtx.Err())
}

func TestWithTimeout(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	time.Sleep(time.Second)
	// 父超时，子超时
	fmt.Println(childCtx.Err())
}

func TestWithDeadline(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Millisecond))
	defer cancel()
	time.Sleep(time.Second)
	// 父过期，子过期
	fmt.Println(childCtx.Err())
}

func TestResetDeadline(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Millisecond))
	fmt.Println(childCtx.Deadline()) // 2022-08-18 17:55:22.2209609 +0800 CST m=+0.003014201 true
	ccCtx, cancel := context.WithDeadline(childCtx, time.Now().Add(time.Second))
	fmt.Println(ccCtx.Deadline()) // 2022-08-18 17:55:22.2209609 +0800 CST m=+0.003014201 true
	defer cancel()
	fmt.Println(ccCtx.Err())
}

func TestWithDone(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithCancel(ctx)
	//ch := ctx.Done()
	cancel()
	//fmt.Println(<-ch)

	select {
	case <-ctx.Done():
		fmt.Println(childCtx.Err())
	}
}

func TestWithTimeOut(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	ch := childCtx.Done()
	select {
	case <-ch:
		fmt.Println(childCtx.Err())
	}
}
