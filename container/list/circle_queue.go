package container

import (
	"fmt"
	"sync"
)

const CQUEUE_SIZE = 5

//SqQueue 结构体定义
type CirCleQueue struct {
	data   [CQUEUE_SIZE]interface{}
	front  int
	rear   int
	isFull bool
	sync.RWMutex
}

//New 新建空队列
func NewCircleQueue() *CirCleQueue {
	return &CirCleQueue{
		front:  0,
		rear:   0,
		isFull: false,
	}
}

// Length 队列长度
func (c *CirCleQueue) Len() interface{} {
	c.RLock()
	defer c.RUnlock()
	if c.isFull {
		return CQUEUE_SIZE
	} else {
		return (c.rear - c.front + CQUEUE_SIZE) % CQUEUE_SIZE
	}
}

// Enqueue 入队
func (c *CirCleQueue) Enqueue(e interface{}) error {
	c.Lock()
	defer c.Unlock()
	if c.rear%CQUEUE_SIZE == c.front && c.isFull {
		return fmt.Errorf("quque is full")
	}
	c.data[c.rear] = e
	c.rear = (c.rear + 1) % CQUEUE_SIZE
	if c.rear == c.front {
		c.isFull = true
	}
	return nil
}

// Dequeue 出队
func (c *CirCleQueue) Dequeue() (e interface{}, err error) {
	c.Lock()
	defer c.Unlock()
	if c.rear == c.front && c.isFull == false {
		return e, fmt.Errorf("quque is empty")
	}
	e = c.data[c.front]
	c.front = (c.front + 1) % CQUEUE_SIZE
	if c.isFull == true {
		c.isFull = false
	}
	return e, nil
}

func (c *CirCleQueue) ForEach(f func(data interface{}) bool) {
	c.RLock()
	defer c.RUnlock()

	front := c.front
	for (front % CQUEUE_SIZE) != c.rear {
		data := c.data[front]
		if data != nil {
			if !f(data) {
				break
			}
		}
		front++
	}
}

func (c *CirCleQueue) Snapshot() []interface{} {
	if c.rear < c.front {
		data := make([]interface{}, 0, CQUEUE_SIZE)
		data = append(data, c.data[c.front:]...)
		data = append(data, c.data[:c.rear]...)
		return data
	} else if c.isFull {
		return c.data[:]
	}
	return c.data[c.front:c.rear]
}
