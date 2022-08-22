package ch2

import (
	"fmt"
	"sync"
	"testing"
)

type SafeMap[T string, V any] struct {
	values map[T]V
	rwLock sync.RWMutex
}

// 如果有key，返回value,如果没有，设置k,v并返回 val,并保证并发安全
// 假如有两个goroutine
// g1 => "a","b"
// g2 => "a","c"
func (s *SafeMap[T, V]) LoadOrStore(key T, val V) (V, bool) {
	// 加读锁，check第一次  加读锁是为了并发效率
	s.rwLock.RLock()
	v, ok := s.values[key]
	s.rwLock.RUnlock() // 这里不可以使用 defer ，因为defer在return中执行，对这个goroutine同时加两个锁会直接报错
	if ok {
		return v, true
	}

	// 加写锁，double check
	// 如果这里不加写锁，则g2中的a，可能会覆盖g1中的a
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	v, ok = s.values[key]
	if ok {
		return v, true
	}
	s.values[key] = val
	return val, false
}

func TestSafeMap(t *testing.T) {
	sMap := SafeMap[string, string]{
		values: map[string]string{},
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			k := "a"
			v := fmt.Sprintf("b%v", i)
			sMap.LoadOrStore(k, v)
		}(i)
	}

	v, _ := sMap.LoadOrStore("a", "c1")
	fmt.Println(v)
}

type SafeList[T any] struct {
	List[T]
	rwLock sync.RWMutex
}

func (s *SafeList[T]) Get(i int) (any, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	return s.List.Get(i), nil
}

func (s *SafeList[T]) Add(i int, el any) (any, error) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	return s.List.Add(i, el), nil
}
