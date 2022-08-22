package ch2

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func add(i int) int {
	//time.Sleep(time.)
	return i + 1
}

func read(i int) {
	fmt.Println("read some data", i)
}

func write(i int) {
	fmt.Println("write some data:", i)
}

var wg sync.WaitGroup
var rwLock sync.RWMutex
var lock sync.Mutex
var num int

// 互斥锁并发增加
func TestLockAdd(t *testing.T) {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			lock.Lock()
			num = add(num)
			defer lock.Unlock()
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("num:", num)
	/** console
	num: 10
	*/
}

// 读锁并发增加
func TestRLockAdd(t *testing.T) {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			rwLock.RLock()
			num = add(num)
			defer rwLock.RUnlock()
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("num:", num)
	/** console
	num: 981  // 结论：读锁是可以重入的，一个协程在读的时候，另一个协程也可以访问资源
	*/
}

// 互斥锁顺序输出
func TestLockSort(t *testing.T) {
	for i := 0; i < 1000; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}

// 互斥锁顺序输出
func TestLockSort02(t *testing.T) {
	for i := 0; i < 1000; i++ {
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
}

// 互斥锁顺序输出:如果要保证顺序输出，则应该在 +1 的地方加锁
func TestLockSort03(t *testing.T) {
	var n int

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			lock.Lock()
			defer lock.Unlock()
			defer wg.Done()

			n++
			fmt.Println(n)
		}()
	}
	wg.Wait()
}

// 互斥锁顺序输出:如果要保证顺序输出，则应该在 +1 的地方加锁
func TestLockSort04(t *testing.T) {
	var n int
	count := make(chan int)
	for {
		wg.Add(1)

		go func() {
			lock.Lock()
			defer lock.Unlock()
			defer wg.Done()
			n++
			count <- n
			fmt.Println(n)
		}()

		if s := <-count; s == 1000 {
			break
		}

		//if n == 1000 {
		//	break
		//}
	}
	wg.Wait()
}

// 读锁会阻塞住写锁，但读锁之间不阻塞
func TestRWLock(t *testing.T) {
	var n int
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Println("read start")
			rwLock.RLock()
			fmt.Println("read", n)
			time.Sleep(time.Microsecond)
			rwLock.RUnlock()
			fmt.Println("read over")
		}()

		go func() {
			fmt.Println("write start")
			rwLock.Lock()
			n++
			fmt.Println("write", n)
			time.Sleep(time.Microsecond)
			rwLock.Unlock()
			fmt.Println("write over")
		}()
	}

	time.Sleep(3 * time.Second)
}

// 写锁会阻塞住读锁
func TestRWLock03(t *testing.T) {
	var n int
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Println("write start")
			rwLock.Lock()
			n++
			fmt.Println("write", n)
			time.Sleep(time.Microsecond)
			rwLock.Unlock()
			fmt.Println("write over")
		}()
	}

	for i := 0; i < 100; i++ {
		go func() {
			fmt.Println("read start")
			rwLock.RLock()
			fmt.Println("read", n)
			time.Sleep(time.Microsecond)
			rwLock.RUnlock()
			fmt.Println("read over")
		}()
	}

	time.Sleep(3 * time.Second)
}

// 读锁会阻塞住写锁，但读锁之间不阻塞
func TestRLock02(t *testing.T) {
	var n int
	//for i := 0; i < 100; i++ {

	//ch := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			//time.Sleep(time.Millisecond)
			//fmt.Println("read start")
			rwLock.RLock()
			fmt.Println("reading :", n)
			rwLock.RUnlock()
			//fmt.Println("read over")
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond)
			rwLock.Lock()
			n++
			fmt.Println("writing:", n)
			rwLock.Unlock()
		}
	}()

	time.Sleep(time.Second * 10)
}

func TestRLock04(t *testing.T) {
	var n int
	//for i := 0; i < 100; i++ {

	//ch := make(chan int)

	// 先读，看看读写能否阻塞住写锁
	for i := 0; i < 10; i++ {
		//time.Sleep(time.Millisecond)
		go func() {
			time.Sleep(time.Millisecond)
			rwLock.RLock()
			fmt.Println("reading :", n)
			// 等待2秒，验证写操作只会在读操作都解锁完才会不阻塞
			time.Sleep(time.Second * 2)
			rwLock.RUnlock()
		}()
	}
	// 验证思路，让读锁拿到，让读锁的执行时间变长，而后RUnlock，而写锁则必须在所有RUnlock之后才会执行。
	for i := 0; i < 10; i++ {
		//time.Sleep(time.Millisecond)
		go func() {
			// 确保读操作锁先拿到
			time.Sleep(time.Second)
			// 如果正确，所有的writing都会在2s后的所有读锁解锁后才会执行
			rwLock.Lock()
			n++
			fmt.Println("writing:", n)
			rwLock.Unlock()
		}()
	}

	time.Sleep(time.Second * 10)
}
