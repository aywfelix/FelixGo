package timer

import (
	"context"
	"math"
	"sync"
	"time"

	. "github.com/aywfelix/felixgo/logger"
	. "github.com/aywfelix/felixgo/utils"
)

const (
	MaxChanBuff  = 2048
	MaxTimeDelay = 50
)

type TimerScheduler struct {
	tw           *TimerWheel
	timerId      uint32
	triggerQueue chan *TimerFunc
	sync.RWMutex
	timerIds map[uint32]bool

	ctx    context.Context
	cancel context.CancelFunc
}

func NewTimerScheduler() *TimerScheduler {
	ctx, cancel := context.WithCancel(context.Background())

	secondTimerWheel := NewTimerWheel(SecondName, SecondInterval, SecondScales, TimersMaxCap, ctx)
	minuteTimerWheel := NewTimerWheel(MinuteName, MinuteInterval, MinuteScales, TimersMaxCap, ctx)
	hourTimerWheel := NewTimerWheel(HourName, HourInterval, HourScales, TimersMaxCap, ctx)

	hourTimerWheel.SetNextTimerWheel(minuteTimerWheel)
	minuteTimerWheel.SetNextTimerWheel(secondTimerWheel)

	secondTimerWheel.Run()
	minuteTimerWheel.Run()
	hourTimerWheel.Run()

	timerScheduler := &TimerScheduler{
		tw:           hourTimerWheel,
		triggerQueue: make(chan *TimerFunc, MaxChanBuff),
		timerIds:     make(map[uint32]bool),
		ctx:          ctx,
		cancel:       cancel,
	}

	return timerScheduler
}

func (t *TimerScheduler) CreateTimerAt(timerFunc *TimerFunc, unixNano int64) {
	t.Lock()
	defer t.Unlock()

	t.timerId++
	t.timerIds[t.timerId] = true
	t.tw.AddTimer(t.timerId, newTimerAt(timerFunc, unixNano))
}

func (t *TimerScheduler) CreateTimerAfter(timerFunc *TimerFunc, duration time.Duration) {
	t.Lock()
	defer t.Unlock()

	t.timerId++
	t.timerIds[t.timerId] = true
	t.tw.AddTimer(t.timerId, newTimerAfter(timerFunc, duration))
}

func (t *TimerScheduler) CancelTimer(timerId uint32) {
	if _, ok := t.timerIds[timerId]; ok {
		delete(t.timerIds, timerId)
	}
}

func (t *TimerScheduler) HasTimer(timerId uint32) bool {
	_, ok := t.timerIds[timerId]
	return ok
}

func (t *TimerScheduler) Start() {
	go func() {
		for {
			select {
			case <-t.ctx.Done():
				return
			default:
				now := TimeMillisecond()
				timerList := t.tw.GetTimersWithIn(MaxTimeDelay * time.Millisecond)
				for timerId, timer := range timerList {
					if math.Abs(float64(now-timer.unixts)) > MaxTimeDelay {
						LogDebug("want call at ", timer.unixts, "; real call at", now, "; delay ", now-timer.unixts)
					}
					if t.HasTimer(timerId) {
						t.triggerQueue <- timer.timerFunc
					}
				}
			}
			time.Sleep(MaxTimeDelay / 2 * time.Millisecond)
		}
	}()
}

func (t *TimerScheduler) Stop() {
	close(t.triggerQueue)
	t.cancel()
	LogInfo("timer scheduler stop...")
}

func NewAutoTimerScheduler() *TimerScheduler {
	timerScheduler := NewTimerScheduler()
	timerScheduler.Start()

	go func() {
		for timerFunc := range timerScheduler.triggerQueue {
			go timerFunc.Call()
		}
	}()

	return timerScheduler
}
