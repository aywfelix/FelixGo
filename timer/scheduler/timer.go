package timer

import (
	"fmt"
	"reflect"
	"time"
	
	"github.com/felix/felixgo/utils"
	. "github.com/felix/felixgo/logger"
)

const (
	//HourName 小时
	HourName = "HOUR"
	//HourInterval 小时间隔ms为精度
	HourInterval = 60 * 60 * 1e3
	//HourScales  12小时制
	HourScales = 12

	//MinuteName 分钟
	MinuteName = "MINUTE"
	//MinuteInterval 每分钟时间间隔
	MinuteInterval = 60 * 1e3
	//MinuteScales 60分钟
	MinuteScales = 60

	//SecondName  秒
	SecondName = "SECOND"
	//SecondInterval 秒的间隔
	SecondInterval = 1e3
	//SecondScales  60秒
	SecondScales = 60
	//TimersMaxCap //每个时间轮刻度挂载定时器的最大个数
	TimersMaxCap = 2048
)

type TimerFunc struct {
	f    func(args ...interface{})
	args []interface{}
}

func NewTimerFunc(f func(args ...interface{}), args ...interface{}) *TimerFunc {
	return &TimerFunc{
		f:    f,
		args: args,
	}
}

func (t *TimerFunc) String() string {
	return fmt.Sprintf("{func: %s, args: %v}", reflect.TypeOf(t.f).Name(), t.args)
}

func (t *TimerFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			LogError("call %s", t.String())
			LogError("call timer func error, ", err)
		}
	}()
	t.f(t.args...)
}

//===========================================================================
type Timer struct {
	timerFunc *TimerFunc
	unixts    int64
}

func newTimerAt(timerFunc *TimerFunc, unixNano int64) *Timer {
	return &Timer{
		timerFunc: timerFunc,
		unixts:    unixNano,
	}
}

func newTimerAfter(timerFunc *TimerFunc, duration time.Duration) *Timer {
	return newTimerAt(timerFunc, time.Now().UnixNano()*int64(duration))
}

func (t *Timer) Run() {
	go func() {
		now := utils.TimeMillisecond()
		if t.unixts > now {
			time.Sleep(time.Duration(t.unixts-now) * time.Millisecond)
		}
		t.timerFunc.Call()
	}()
}
