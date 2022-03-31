package utils

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type spinLock uint32

const maxBackOff = 64

func (s *spinLock) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapUint32((*uint32)(s), 0, 1) {
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < maxBackOff {
			backoff <<= 1
		}
	}
}

func (s *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(s), 0)
}

func NewSpinLock() sync.Locker {
	return new(spinLock)
}
