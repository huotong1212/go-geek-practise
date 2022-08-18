package ch1

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContextParentCancel(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second))

	ccCtx := context.WithValue(childCtx, "clearlove", "厂长")

	// 调用父级context的Cancel方法，子Context也会呗cancel掉
	cancel()

	switch ccCtx.Err() {
	case context.DeadlineExceeded:
		fmt.Println("context deadline error")
	case context.Canceled:
		fmt.Println("context canceled error")
	}
	/** console
	context canceled error
	*/
}

func TestContextParentDeadline(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Millisecond))

	ccCtx := context.WithValue(childCtx, "clearlove", "厂长")

	defer cancel()
	// 父级context的DeadlineError ，子Context也会deadline

	time.Sleep(time.Second)
	switch ccCtx.Err() {
	case context.DeadlineExceeded:
		fmt.Println("context deadline error")
	case context.Canceled:
		fmt.Println("context canceled error")
	}
	/** console
	context deadline error
	*/
}

func TestContextParentValue(t *testing.T) {
	ctx := context.Background()
	childCtx := context.WithValue(ctx, "uzi", "乌兹")
	ccCtx := context.WithValue(childCtx, "clearlove", "厂长")

	// 子context可以拿到父Context中的Value
	uzi := ccCtx.Value("uzi")
	fmt.Println("uzi:", uzi)

	// 父Context不可拿到子Context中的Value
	clearLove := childCtx.Value("clearlove")
	fmt.Println("clearlove", clearLove)

	/** console
	uzi: 乌兹
	clearlove <nil>
	*/
}
