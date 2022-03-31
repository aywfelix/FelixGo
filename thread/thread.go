package utils

import (
	"runtime/debug"
	"sync"
	
	. "github.com/felix/felixgo/logger"
)

type taskFunc func(args ...interface{}) error

type Task struct {
	tf    taskFunc
	param []interface{}
}

var taskPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return new(Task)
	},
}

func NewTask(f taskFunc, args ...interface{}) *Task {
	task := taskPool.Get().(*Task)
	task.tf = f
	task.param = args
	return task
}

func (t *Task) Execute() error {
	return t.tf(t.param...)
}

//=====================================================================
type Thread struct {
	threadId      int
	taskQueueSize int
	taskQueue     chan *Task
}

func newThread(threadId, taskQueueSize int) *Thread {
	return &Thread{
		threadId:      threadId,
		taskQueueSize: taskQueueSize,
		taskQueue:     make(chan *Task, taskQueueSize),
	}
}

func (t *Thread) submit(task *Task) {
	t.taskQueue <- task
}

func (t *Thread) worker() {
	for task := range t.taskQueue {
		task.Execute()

		task.tf = nil
		task.param = nil
		taskPool.Put(task)
	}
}

func (t *Thread) startThread() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				LogError("the thread of threadpool go work error, %s\n stack:%s", err, debug.Stack())
			}
		}()
		t.worker()
	}()
}

func (t *Thread) stopThread() {
	close(t.taskQueue)
	t.taskQueue = nil
}

// 开启多个线程，每个线程都有一个独立的任务队列，所有队列并行执行
type ThreadPool struct {
	threadPoolSize int
	threadPool     []*Thread
	threadWorkSize int
}

func NewThreadPool(poolSize, threadWorkSize int) *ThreadPool {
	return &ThreadPool{
		threadPoolSize: poolSize,
		threadWorkSize: threadWorkSize,
		threadPool:     make([]*Thread, poolSize),
	}
}

func (d *ThreadPool) Start() {
	for i := 0; i < d.threadPoolSize; i++ {
		d.threadPool[i] = newThread(i, d.threadWorkSize)

		d.threadPool[i].startThread()
	}
}

func (d *ThreadPool) Stop() {
	for i := 0; i < d.threadPoolSize; i++ {
		d.threadPool[i].stopThread()
	}
}

func (d *ThreadPool) Submit(threadId int, task *Task) {
	d.threadPool[threadId].submit(task)
}

//================================================================
// 多个线程，一个任务队列，多个线程竞争从队列中获取任务并执行
type TaskPool struct {
	taskQueue     chan *Task
	taskQueueSize int
	threadSize    int
}

func NewTaskPool(threadSize, taskQueueSize int) *TaskPool {
	t := &TaskPool{
		threadSize:    threadSize,
		taskQueueSize: taskQueueSize,
		taskQueue:     make(chan *Task, taskQueueSize),
	}
	return t
}

func (t *TaskPool) worker() {
	for task := range t.taskQueue {
		task.Execute()

		task.tf = nil
		task.param = nil
		taskPool.Put(task)
	}
}

func (t *TaskPool) Start() {
	for i := 0; i < t.threadSize; i++ {
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					LogError("TaskPool go work error, %s\n stack:%s", err, debug.Stack())
				}
			}()
			t.worker()
		}()
	}
}

func (t *TaskPool) Submit(task *Task) {
	t.taskQueue <- task
}

func (t *TaskPool) Stop() {
	close(t.taskQueue)
	t.taskQueue = nil
}

//===============================================================
