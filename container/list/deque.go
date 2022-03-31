package container

import (
	"container/list"
	"fmt"
	"sync"
)

type Deque struct {
	sync.RWMutex
	container *list.List // 双向链表
	capacity  int
}

/*
Deque使用双向链表实现，每个操作的时间复杂度为O（1）。
并且是同步&&并发安全的
*/
func NewDeque() *Deque {
	return NewCappedDeque(-1)
}

/*
newCappedRequest创建具有指定容量限制的Deque。
 -1 表示不限制容量
*/
func NewCappedDeque(capacity int) *Deque {
	return &Deque{
		RWMutex:   sync.RWMutex{},
		container: list.New(),
		capacity:  capacity,
	}
}

/*
PushBack： 在Deque尾部以O（1）时间复杂度插入元素，
如果成功，则返回true；如果Deque已满容量，则返回false。
*/
func (s *Deque) PushBack(item interface{}) bool {
	s.Lock()
	defer s.Unlock()

	// 如果不限制容量 || 当前容量 < 限制容量 ， 则尾插
	if s.capacity < 0 || s.container.Len() < s.capacity {
		s.container.PushBack(item)
		return true
	}

	return false
}

/*
PushFront： 在Deque头部以O（1）时间复杂度插入元素，
如果成功，则返回true；如果Deque已满容量，则返回false。
*/
func (s *Deque) PushFront(item interface{}) bool {
	s.Lock()
	defer s.Unlock()

	// 如果不限制容量 || 当前容量 < 限制容量 ， 则头插
	if s.capacity < 0 || s.container.Len() < s.capacity {
		s.container.PushFront(item)
		return true
	}

	return false
}

// popBack： 时间复杂度O（1）, 移除末尾元素
// 返回移除的元素。 如果当前队列为空，则返回nil
func (s *Deque) popBack() interface{} {
	s.Lock()
	defer s.Unlock()

	var item interface{} = nil
	var lastContainerItem *list.Element = nil

	lastContainerItem = s.container.Back() // 获取最末尾的元素
	if lastContainerItem != nil {
		item = s.container.Remove(lastContainerItem)
	}

	return item
}

// PopFront： 时间复杂度O（1）, 移除队首元素
// 返回移除的元素。 如果当前队列为空，则返回nil
func (s *Deque) PopFront() interface{} {
	s.Lock()
	defer s.Unlock()

	var item interface{} = nil
	var lastContainerItem *list.Element = nil

	lastContainerItem = s.container.Front() // 获取队首元素
	if lastContainerItem != nil {
		item = s.container.Remove(lastContainerItem)
	}

	return item
}

// Front()： 时间复杂度O（1）, 获取队首元素
func (s *Deque) Front() interface{} {
	s.Lock()
	defer s.Unlock()

	item := s.container.Front() // 获取队首元素
	if item != nil {
		return item.Value
	}

	return nil
}

// Back()： 时间复杂度O（1）, 获取队尾元素
func (s *Deque) Back() interface{} {
	s.Lock()
	defer s.Unlock()

	item := s.container.Back() // 获取队尾元素
	if item != nil {
		return item.Value
	}

	return nil
}

// Size()： 时间复杂度O（1）, 获取队列长度
func (s *Deque) Size() int {
	s.Lock()
	defer s.Unlock()

	return s.container.Len()
}

// Capacity()： 时间复杂度O（1）, 获取队列容量， -1表示不限制容量
func (s *Deque) Capacity() int {
	s.Lock()
	defer s.Unlock()

	return s.capacity
}

// Empty()： 时间复杂度O（1）,检查队列是否为空
func (s *Deque) Empty() bool {
	s.Lock()
	defer s.Unlock()

	return s.container.Len() == 0
}

// Empty()： 时间复杂度O（1）,检查队列是否满了
func (s *Deque) Full() bool {
	s.Lock()
	defer s.Unlock()

	return s.capacity > 0 && s.container.Len() >= s.capacity
}

func (q *Deque) Dump() {
	for iter := q.container.Front(); iter != nil; iter = iter.Next() {
		fmt.Println("item:", iter.Value)
	}
}
