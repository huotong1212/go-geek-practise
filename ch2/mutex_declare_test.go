package ch2

import "sync"

/**
介绍一般是如何使用mutex的
*/

// 绝对不可以
var PublicSource map[string]string
var PublicLock sync.Mutex

// 非必要不使用
var privateSource map[string]string
var privateLock sync.Mutex

// 建议使用
type source struct {
	safeSource map[string]string
	lock       sync.Mutex
}

func (s *source) Add(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.safeSource[key] = val
}

var SafeSource source
