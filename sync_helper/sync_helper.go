package syncHelper

import (
	"sync/atomic"
)


type AtomicLock struct {
	flag uint32
}


func NewAtomicLock() *AtomicLock{
	return &AtomicLock{flag: 0}
}

func (s *AtomicLock)TryLock() bool {
	return atomic.CompareAndSwapUint32(&s.flag, 0, 1)
}

func (s *AtomicLock)TryUnlock() bool {
	return atomic.CompareAndSwapUint32(&s.flag, 1, 0)
}

func (s *AtomicLock)Unlock() {
	atomic.StoreUint32(&s.flag, 0)
}
