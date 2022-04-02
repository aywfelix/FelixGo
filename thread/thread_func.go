package utils

import (
	"runtime/debug"
	"time"

	. "github.com/aywfelix/felixgo/logger"
)

type Func func()
type FuncArgs func(args ...interface{})

func TimeLoopArgs(idleTime time.Duration, f FuncArgs, args ...interface{}) {
	loop := func() {
		defer func() {
			err := recover()
			if err != nil {
				LogError("func TimeLoop error, %s\n stack:%s", err, debug.Stack())
			}
		}()
		f(args...)
	}
	go func() {
		for {
			loop()
			time.Sleep(idleTime)
		}
	}()
}

func TimeLoop(idleTime time.Duration, f Func) {
	loop := func() {
		defer func() {
			err := recover()
			if err != nil {
				LogError("func TimeLoop error, %s\n stack:%s", err, debug.Stack())
			}
		}()
		f()
	}
	go func() {
		for {
			loop()
			time.Sleep(idleTime)
		}
	}()
}

func GoRunArgs(f FuncArgs, args ...interface{}) {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				LogError("func GoRunArgs error, %s\n stack:%s", err, debug.Stack())
				return
			}
		}()
		f(args...)
	}()
}

func GoRun(f Func) {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				LogError("func GoRun error, %s\n stack:%s", err, debug.Stack())
				return
			}
		}()
		f()
	}()
}

func init() {

}
