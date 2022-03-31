package db

import (
	"context"
	"fmt"
	. "github.com/felix/felixgo/logger"
	"os"
	"time"
	"github.com/felix/felixgo/utils"
)

const (
	GLOCK_TIMEOUT = 5
)

type IGLocker interface {
	Lock() bool
	UnLock()
	Release()
}

type GLocker struct {
	key        string
	token      string
	lockedChan chan bool
	ctx        context.Context
	cancel     func()
}

func NewGLocker(key string) *GLocker {
	g := &GLocker{
		key:        key,
		lockedChan: make(chan bool, 1),
	}
	g.token = fmt.Sprintf("%s-%s-%d-%p", key, utils.GetLocalIPV4(), os.Getpid(), g)
	g.ctx, g.cancel = context.WithCancel(context.Background())
	g.goLock()
	return g
}

func (g *GLocker) Lock() bool {
	for {
		select {
		case <-g.ctx.Done():
			LogInfo("GLocker Lock exit...")
			return false
		case isLocked := <-g.lockedChan:
			if isLocked {
				return true
			}
		}
	}
}

func (g *GLocker) UnLock() {
	g.unlock()
}

func (g *GLocker) Release() {
	g.cancel()
	g.unlock()
}

func (g *GLocker) goLock() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				LogError("GLocker:Lock err, %v, token: %v\n", err, g.token)
			}
		}()
		ticker := time.NewTicker(time.Millisecond * 500)
		for {
			select {
			case <-g.ctx.Done():
				LogInfo("GLocker goroutine exit...")
				return
			case <-ticker.C:
				if g.lock() {
					g.lockedChan <- true
				}
			}
		}
	}()
}

func (g *GLocker) lock() bool {
	if ret, err := DbRedis.DoScript(REDIS_SCRIPT_GLOCK, g.key, g.token, GLOCK_TIMEOUT); err == nil {
		if bytes, ok := ret.([]byte); ok {
			return string(bytes) == "ok"
		}
	} else {
		LogError("GLocker:script err, %v, token: %s\n", err.Error(), g.token)
		return false
	}
	return false
}

func (g *GLocker) unlock() {
	_, err := DbRedis.DoScript(REDIS_SCRIPT_GUNLOCK, g.key, g.token)
	if err != nil {
		LogError("GLocker: unlock failed, %v, token: %s\n", err, g.token)
	}
}
