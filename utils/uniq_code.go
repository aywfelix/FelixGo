package utils

import "sync"

type uniqcode struct {
	key uint
	index uint
	timeSec uint64
	genMutex sync.Mutex
}

func NewUniqCode(key uint) *uniqcode {
	code := &uniqcode{
		key: key,
		index: 0,
		timeSec : uint64(TimeSecond()),
	}
	return code
}

func (u *uniqcode) Gen() uint64 {
	u.genMutex.Lock()
	defer u.genMutex.Unlock()
	code := u.timeSec << 32 | uint64(u.key << 16) | uint64(u.index & 0xFFFF)
	u.index++
	if u.index > 0xFFFF {
		u.index = 0
		u.timeSec++
	}
	return code
}

