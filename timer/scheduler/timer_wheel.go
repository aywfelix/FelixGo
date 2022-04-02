package timer

import (
	"context"
	"sync"
	"time"

	. "github.com/aywfelix/felixgo/logger"
	"github.com/aywfelix/felixgo/utils"
)

type TimerWheel struct {
	//TimeWheel的名称
	name string
	//刻度的时间间隔，单位ms
	interval int64
	//每个时间轮上的刻度数
	scales int
	//当前时间指针的指向
	curIndex int
	//每个刻度所存放的timer定时器的最大容量
	maxCap int
	//当前时间轮上的所有timer
	timerQueue map[int]map[uint32]*Timer //map[int] VALUE  其中int表示当前时间轮的刻度,
	// map[int] map[uint32] *Timer, uint32表示Timer的ID号
	//下一层时间轮
	nextTimeWheel *TimerWheel
	//互斥锁（继承RWMutex的 RWLock,UnLock 等方法）
	sync.RWMutex
	ctx context.Context
}

func NewTimerWheel(name string, interval int64, scales int, maxCap int, ctx context.Context) *TimerWheel {
	tw := &TimerWheel{
		name:       name,
		interval:   interval,
		scales:     scales,
		maxCap:     maxCap,
		timerQueue: make(map[int]map[uint32]*Timer, scales),
		ctx:        ctx,
	}

	for i := 0; i < scales; i++ {
		tw.timerQueue[i] = make(map[uint32]*Timer, maxCap)
	}

	return tw
}

func (t *TimerWheel) SetNextTimerWheel(tw *TimerWheel) {
	t.nextTimeWheel = tw
}

func (t *TimerWheel) addTimer(timerId uint32, timer *Timer, forceNext bool) {
	delayInterval := timer.unixts - utils.TimeMillisecond()

	if delayInterval >= t.interval {
		dn := delayInterval / t.interval
		t.timerQueue[(t.curIndex+int(dn))%t.scales][timerId] = timer
		return
	}

	if delayInterval < t.interval && t.nextTimeWheel == nil {
		if forceNext {
			t.timerQueue[(t.curIndex+1)%t.scales][timerId] = timer
		} else {
			t.timerQueue[t.curIndex][timerId] = timer
		}
		return
	}

	t.nextTimeWheel.addTimer(timerId, timer, forceNext)
}

func (t *TimerWheel) AddTimer(timerId uint32, timer *Timer) {
	t.Lock()
	defer t.Unlock()

	t.addTimer(timerId, timer, false)
}

func (t *TimerWheel) Run() {
	go func() {
		for {
			select {
			case <-t.ctx.Done():
				return
			default:
				t.Lock()
				defer t.Unlock()

				curTimers := t.timerQueue[t.curIndex]
				t.timerQueue[t.curIndex] = make(map[uint32]*Timer, t.maxCap)
				for timerId, timer := range curTimers {
					t.addTimer(timerId, timer, true)
				}

				t.curIndex = (t.curIndex + 1) % t.scales
			}
			time.Sleep(time.Duration(t.interval) * time.Millisecond)
		}
	}()

	LogInfo("timer wheel stop...")
}

func (t *TimerWheel) GetTimersWithIn(duration time.Duration) map[uint32]*Timer {
	leafTw := t
	for leafTw != nil {
		leafTw = leafTw.nextTimeWheel
	}

	leafTw.Lock()
	defer leafTw.Unlock()

	now := utils.TimeMillisecond()

	timerList := make(map[uint32]*Timer)
	for timerId, timer := range leafTw.timerQueue[leafTw.curIndex] {
		if timer.unixts-now < int64(duration/1e6) {
			timerList[timerId] = timer
			delete(leafTw.timerQueue[leafTw.curIndex], timerId)
		}
	}
	return timerList
}
