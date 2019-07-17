package syncHelper

import (
	"sync/atomic"
	"runtime"
	"sync"
)


type AtomicLock struct {
	_    sync.Mutex // for copy protection compiler warning
	flag uint32
}

func NewAtomicLock() *AtomicLock{
	return &AtomicLock{flag: 0}
}

func (s *AtomicLock) Lock() {
	for !atomic.CompareAndSwapUint32(&s.flag, 0, 1) {
		runtime.Gosched()
	}
}

func (s *AtomicLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&s.flag, 0, 1)
}

func (s *AtomicLock) TryUnlock() bool {
	return atomic.CompareAndSwapUint32(&s.flag, 1, 0)
}

func (s *AtomicLock) Unlock() {
	atomic.StoreUint32(&s.flag, 0)
}
