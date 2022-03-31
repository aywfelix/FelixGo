package container

import (
	"fmt"
	"sync"
)

type Queue struct {
	head  *ListNode
	tail  *ListNode
	count int
	sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		head:  nil,
		tail:  nil,
		count: 0,
	}
}

// 头部插入
func (q *Queue) PushFront(Data interface{}) {
	if Data == nil {
		return
	}
	q.Lock()
	defer q.Unlock()
	node := newlinkNode(Data)
	if q.head == nil {
		q.head = node
		q.tail = node
		q.count++
		return
	}
	q.head.Prev = node
	node.Next = q.head
	q.head = node
	q.count++
}

// 尾部插入
func (q *Queue) PushTail(Data interface{}) {
	if Data == nil {
		return
	}
	q.Lock()
	defer q.Unlock()
	node := newlinkNode(Data)
	if q.head == nil {
		q.head = node
		q.tail = node
		q.count++
		return
	}
	q.tail.Next = node
	node.Prev = q.tail
	q.tail = node
	q.count++
}

// 删除头部
func (q *Queue) PopFront() interface{} {
	q.Lock()
	defer q.Unlock()

	if q.count == 0 {
		return nil
	}
	retNode := q.head
	q.head = q.head.Next
	q.head.Prev = nil
	q.count--
	return retNode.Data
}

// 删除尾部
func (q *Queue) PopTail() interface{} {
	q.Lock()
	defer q.Unlock()

	if q.count == 0 {
		return nil
	}
	retNode := q.tail
	q.tail = q.tail.Prev
	q.tail.Next = nil
	q.count--
	return retNode.Data
}

func (q *Queue) Clear() {
	q.Lock()
	defer q.Unlock()

	q.head = nil
	q.tail = nil
	q.count = 0
}

func (q *Queue) Len() int {
	q.Lock()
	defer q.Unlock()

	return q.count
}

func (q *Queue) IsEmpty() bool {
	return q.Len() == 0
}

func (q *Queue) ForEach(f func(node *ListNode) bool) {
	node := q.head
	for node != nil {
		if !f(node) {
			break
		}
		node = node.Next
	}
}

func (q *Queue) Snapshot() []*ListNode {
	if q.count == 0 {
		return nil
	}
	var nodes []*ListNode
	head := q.head
	tail := q.tail
	for head != nil && tail != nil {
		if head == tail {
			nodes = append(nodes, head)
			break
		}
		nodes = append(nodes, head)
		nodes = append(nodes, tail)
		head = head.Next
		tail = tail.Prev
	}
	return nodes
}

func (q *Queue) String() string {
	if q.head == nil {
		return ""
	}
	node := q.head
	s := "print link list:\n"
	for node != nil {
		s += fmt.Sprintf("%v ", node.Data)
		node = node.Next
	}
	return s
}
