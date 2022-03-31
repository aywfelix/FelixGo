package container

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

const (
	PROBABILITY = 0.25
)

// 定义跳表节点
type SkipNode struct {
	index     uint64
	value     interface{}
	nextNodes []*SkipNode
}

func newSkipNode(index uint64, value interface{}, level int) *SkipNode {
	return &SkipNode{
		index:     index,
		value:     value,
		nextNodes: make([]*SkipNode, level, level),
	}
}

func (s *SkipNode) Value() interface{} {
	return s.value
}

func (s *SkipNode) Index() uint64 {
	return s.index
}

//==========================================================================
// 定义跳表
type SkipList struct {
	level  int
	length int32
	head   *SkipNode
	tail   *SkipNode
	mutex  sync.RWMutex
}

func NewSkipList(level int) *SkipList {
	head := newSkipNode(0, nil, level)
	var tail *SkipNode
	for i := 0; i < len(head.nextNodes); i++ {
		head.nextNodes[i] = tail
	}
	return &SkipList{
		level:  level,
		length: 0,
		head:   head,
		tail:   tail,
	}
}

func (s *SkipList) Insert(index uint64, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	previousNodes, _ := s.searchWithPreviousNodes(index)
	// 默认可以插入相同index的节点
	// if currentNode != s.head && currentNode.index == index {
	// 	currentNode.value = value
	// 	return
	// }

	newNode := newSkipNode(index, value, s.randLevel())
	for i := len(newNode.nextNodes) - 1; i >= 0; i-- {
		newNode.nextNodes[i] = previousNodes[i].nextNodes[i]
		previousNodes[i].nextNodes[i] = newNode

		previousNodes[i] = nil
	}
	atomic.AddInt32(&s.length, 1)

	for i := len(newNode.nextNodes); i < len(previousNodes); i++ {
		previousNodes[i] = nil
	}
}

func (s *SkipList) Delete(index uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	previousNodes, currentNode := s.searchWithPreviousNodes(index)
	if currentNode != s.head && currentNode.index == index {
		for i := 0; i < len(currentNode.nextNodes); i++ {
			previousNodes[i].nextNodes[i] = currentNode.nextNodes[i]
			currentNode.nextNodes[i] = nil
			previousNodes[i] = nil
		}
		atomic.AddInt32(&s.length, -1)
	}

	for i := len(currentNode.nextNodes); i < len(previousNodes); i++ {
		previousNodes[i] = nil
	}
}

func (s *SkipList) Search(index uint64) (*SkipNode, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := s.searchWithoutPreviousNodes(index)
	return result, result != nil
}

func (s *SkipList) SearchRagne(minIndex, maxIndex uint64) ([]*SkipNode, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if maxIndex <= minIndex {
		minIndex, maxIndex = maxIndex, minIndex
	}
	node := s.searchWithoutPreviousNodes(minIndex)
	if node == nil {
		return nil, false
	}
	var nodes []*SkipNode
	nodes = append(nodes, node)
	for {
		node = node.nextNodes[0]
		if node == nil || node.index > maxIndex {
			break
		}
		nodes = append(nodes, node)
	}
	return nodes, true
}

func (s *SkipList) ForEach(f func(node *SkipNode) bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	nodes := s.Snapshot()
	for _, node := range nodes {
		if !f(node) {
			break
		}
	}
}

func (s *SkipList) Length() int32 {
	return atomic.LoadInt32(&s.length)
}

func (s *SkipList) Level() int {
	return s.level
}

func (s *SkipList) searchWithPreviousNodes(index uint64) ([]*SkipNode, *SkipNode) {
	previousNodes := make([]*SkipNode, s.level)
	currentNode := s.head

	for l := s.level - 1; l >= 0; l-- {
		for currentNode.nextNodes[l] != s.tail && currentNode.nextNodes[l].index < index {
			currentNode = currentNode.nextNodes[l]
		}
		previousNodes[l] = currentNode
	}

	if currentNode.nextNodes[0] != s.tail {
		currentNode = currentNode.nextNodes[0]
	}

	return previousNodes, currentNode
}

func (s *SkipList) searchWithoutPreviousNodes(index uint64) *SkipNode {
	currentNode := s.head
	for l := s.level - 1; l >= 0; l-- {
		for currentNode.nextNodes[l] != s.tail && currentNode.nextNodes[l].index < index {
			currentNode = currentNode.nextNodes[l]
		}
	}
	currentNode = currentNode.nextNodes[0]
	if currentNode != s.tail && currentNode.index >= index {
		return currentNode
	}
	return nil
}

func (s *SkipList) Snapshot() []*SkipNode {
	result := make([]*SkipNode, s.length)
	i := 0

	currentNode := s.head.nextNodes[0]
	for currentNode != s.tail {
		node := &SkipNode{
			index:     currentNode.index,
			value:     currentNode.value,
			nextNodes: nil,
		}

		result[i] = node
		currentNode = currentNode.nextNodes[0]
		i++
	}

	return result
}

func (s *SkipList) randLevel() int {
	level := 1
	for rand.Float64() <= PROBABILITY && level < s.level {
		level++
	}
	return level
}
