package timer

import (
	"runtime/debug"

	. "github.com/aywfelix/felixgo/logger"
)

const (
	HIGH_FRAME_INTERVAL = 30
	HIGH_FRAME_LIMIT    = 100
	LOW_FRAME_INTERVAL  = 100
	LOW_FRAME_LIMIT     = 1000
)

type WorkerFunc func(nowTime int64)

type ILoopWorker interface {
	RegHighWorker(worker *HighFrameWorker)
	RegLowWorker(worker *LowFrameWorker)
	RegFixWorker(worker *FixWorker)
	Stop()
}

type TimerLoop struct {
	highFrameInterval int64
	lowFrameInterval  int64
	highWorkers       []*HighFrameWorker
	lowWorkers        []*LowFrameWorker
	fixWorkers        map[uint64]*FixWorker
	lastTime          int64
	isStop            bool
}

func NewTimerLoop(highInterval, lowInterval int64) *TimerLoop {
	timerLoop := &TimerLoop{
		highFrameInterval: highInterval,
		lowFrameInterval:  lowInterval,
		highWorkers:       make([]*HighFrameWorker, 0),
		lowWorkers:        make([]*LowFrameWorker, 0),
		fixWorkers:        make(map[uint64]*FixWorker),
		lastTime:          0,
		isStop:            false,
	}
	return timerLoop
}

// call this function after init all resoure over
func (t *TimerLoop) TimeLoop(nowTime int64) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("timer loop run error, %s\nStack trace: %s\n", err, debug.Stack())
		}
	}()

	if t.isStop {
		return
	}
	if nowTime >= (t.highFrameInterval + t.lastTime) {
		for _, worker := range t.highWorkers {
			worker.Handle(nowTime)
		}
	}
	if nowTime >= (t.lowFrameInterval + t.lastTime) {
		for _, worker := range t.lowWorkers {
			worker.Handle(nowTime)
		}
	}

	for workerID, worker := range t.fixWorkers {
		if worker.count > 0 {
			worker.workerFunc(nowTime)
			worker.count--
			if worker.count == 0 {
				delete(t.fixWorkers, workerID)
			}
		} else {
			worker.workerFunc(nowTime)
		}
	}
	t.lastTime = nowTime
}

func (t *TimerLoop) Stop() {
	t.isStop = true
	t.highWorkers = nil
	t.lowWorkers = nil
	t.fixWorkers = nil
}

func (t *TimerLoop) RegHighWorker(worker *HighFrameWorker) {
	if worker == nil {
		return
	}
	t.highWorkers = append(t.highWorkers, worker)
}
func (t *TimerLoop) RegLowWorker(worker *LowFrameWorker) {
	if worker == nil {
		return
	}
	t.lowWorkers = append(t.lowWorkers, worker)
}
func (t *TimerLoop) RegFixWorker(worker *FixWorker) {
	if worker == nil {
		return
	}
	if _, ok := t.fixWorkers[worker.workerID]; ok {
		return
	}
	t.fixWorkers[worker.workerID] = worker
}

type IWorker interface {
	Handle(nowTime int64)
}

type Worker struct {
	interval   int64
	lastTime   int64
	workerFunc WorkerFunc
}

func NewWoker(interval int64, workerFunc WorkerFunc) *Worker {
	worker := &Worker{
		interval:   interval,
		lastTime:   0,
		workerFunc: workerFunc,
	}
	return worker
}

func (w *Worker) Handle(nowTime int64) {
	if w.workerFunc != nil {
		if nowTime >= (w.lastTime + w.interval) {
			w.workerFunc(nowTime)
		}
	}
	w.lastTime = nowTime
}

type HighFrameWorker struct {
	*Worker
}

func NewHighFrameWorker(interval int64, workerFunc WorkerFunc) *HighFrameWorker {
	if interval < HIGH_FRAME_INTERVAL || interval >= HIGH_FRAME_LIMIT {
		return nil
	}
	worker := &HighFrameWorker{
		Worker: NewWoker(interval, workerFunc),
	}
	return worker
}

type LowFrameWorker struct {
	*Worker
}

func NewLowFrameWorker(interval int64, workerFunc WorkerFunc) *LowFrameWorker {
	if interval < LOW_FRAME_INTERVAL || interval >= LOW_FRAME_LIMIT {
		return nil
	}
	worker := &LowFrameWorker{
		Worker: NewWoker(interval, workerFunc),
	}
	return worker
}

type FixWorker struct {
	*Worker
	workerID uint64
	count    int // <=0 代表无限制
}

func NewFixWorker(workerID uint64, interval int64, workerFunc WorkerFunc, count int) *FixWorker {
	worker := &FixWorker{
		Worker:   NewWoker(interval, workerFunc),
		count:    count,
		workerID: workerID,
	}
	return worker
}
