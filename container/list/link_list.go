// 此文件实现了线程安全一个队列和线程安全一个双链表
package container

import (
	"fmt"
	"sync"
)

type LinkList struct {
	Head  *ListNode
	Tail  *ListNode
	dic   map[interface{}]*ListNode
	count int
	sync.Mutex
}

func NewLinkList() *LinkList {
	return &LinkList{
		Head:  nil,
		Tail:  nil,
		dic:   make(map[interface{}]*ListNode),
		count: 0,
	}
}

// 头部插入
func (l *LinkList) PushFront(Data interface{}) {
	if Data == nil {
		return
	}
	l.Lock()
	defer l.Unlock()
	if _, ok := l.dic[Data]; ok {
		return
	}
	node := newlinkNode(Data)
	if l.Head == nil {
		l.Head = node
		l.Tail = node
	} else {
		l.Head.Prev = node
		node.Next = l.Head
		l.Head = node
	}
	l.dic[Data] = node
	l.count++
}

// 尾部插入
func (l *LinkList) PushTail(Data interface{}) {
	if Data == nil {
		return
	}
	l.Lock()
	defer l.Unlock()
	if _, ok := l.dic[Data]; ok {
		return
	}
	node := newlinkNode(Data)
	if l.Head == nil {
		l.Head = node
		l.Tail = node
	} else {
		l.Tail.Next = node
		node.Prev = l.Tail
		l.Tail = node
	}
	l.dic[Data] = node
	l.count++
}

// Data1前插入Data2
func (l *LinkList) PushBefore(node *ListNode, inData interface{}) {
	if node == nil || inData == nil || l.count == 0 {
		return
	}
	l.Lock()
	defer l.Unlock()
	if _, ok := l.dic[node.Data]; !ok {
		return
	}
	if _, ok := l.dic[inData]; ok {
		return
	}
	inNode := newlinkNode(inData)

	if node.Prev != nil {
		node.Prev.Next = inNode
	}
	inNode.Prev = node.Prev
	node.Prev = inNode
	inNode.Next = node

	l.dic[inData] = inNode
	l.count++
}

// Data1后插入Data2
func (l *LinkList) PushAfter(node *ListNode, inData interface{}) {
	if node == nil || inData == nil || l.count == 0 {
		return
	}
	l.Lock()
	defer l.Unlock()
	if _, ok := l.dic[node.Data]; !ok {
		return
	}
	if _, ok := l.dic[inData]; ok {
		return
	}
	inNode := newlinkNode(inData)

	if node.Next != nil {
		node.Next.Prev = inNode
	}
	inNode.Next = node.Next
	node.Next = inNode
	inNode.Prev = node

	l.dic[inData] = inNode
	l.count++
}

// 删除头部
func (l *LinkList) PopFront() interface{} {
	l.Lock()
	defer l.Unlock()

	if l.count == 0 {
		return nil
	}

	retNode := l.Head
	l.Head = l.Head.Next
	l.Head.Prev = nil

	Data := retNode.Data
	delete(l.dic, Data)
	l.count--
	return Data
}

// 删除尾部
func (l *LinkList) PopTail() interface{} {
	l.Lock()
	defer l.Unlock()

	if l.count == 0 {
		return nil
	}

	retNode := l.Tail
	l.Tail = l.Tail.Prev
	l.Tail.Next = nil

	Data := retNode.Data
	delete(l.dic, Data)
	l.count--
	return Data
}

// 删除某个元素
func (l *LinkList) Pop(Data interface{}) {
	if Data == nil {
		return
	}
	l.Lock()
	defer l.Unlock()

	if l.count == 0 {
		return
	}
	if node, ok := l.dic[Data]; ok {
		if node.Prev != nil {
			node.Prev.Next = node.Next
		}
		if node.Next != nil {
			node.Next.Prev = node.Prev
		}
		delete(l.dic, Data)
		l.count--
	}
}

// 获取某个元素的前一个节点
func (l *LinkList) GetBefore(Data interface{}) *ListNode {
	if Data == nil {
		return nil
	}
	l.Lock()
	defer l.Unlock()
	if l.count == 0 {
		return nil
	}
	if node, ok := l.dic[Data]; ok {
		return node.Prev
	}
	return nil
}

// 获取某个元素的后一个节点
func (l *LinkList) GetAfter(Data interface{}) *ListNode {
	if Data == nil {
		return nil
	}
	l.Lock()
	defer l.Unlock()
	if l.count == 0 {
		return nil
	}
	if node, ok := l.dic[Data]; ok {
		return node.Next
	}
	return nil
}

func (l *LinkList) Clear() {
	l.Lock()
	defer l.Unlock()

	l.Head = nil
	l.Tail = nil
	l.dic = nil
	l.count = 0
}

func (l *LinkList) Len() int {
	l.Lock()
	defer l.Unlock()

	return l.count
}

func (l *LinkList) IsEmpty() bool {
	return l.Len() == 0
}

func (l *LinkList) ForEach(f func(node *ListNode) bool) {
	l.Lock()
	defer l.Unlock()

	head := l.Head
	for head != nil {
		if !f(head) {
			break
		}
		head = head.Next
	}
}

func (l *LinkList) Snapshot() []*ListNode {
	l.Lock()
	defer l.Unlock()

	if l.count == 0 {
		return nil
	}
	var array []*ListNode
	head := l.Head
	tail := l.Tail
	for head != nil && tail != nil {
		if head == tail {
			array = append(array, head)
			break
		}
		array = append(array, head)
		array = append(array, tail)
		head = head.Next
		tail = tail.Prev
	}
	return array
}

func (l *LinkList) String() string {
	if l.Head == nil {
		return ""
	}
	node := l.Head
	s := "print link list:\n"
	for node != nil {
		s += fmt.Sprintf("%v ", node.Data)
		node = node.Next
	}
	return s
}
